package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"survey-system/internal/api/handler"
	"survey-system/internal/api/router"
	"survey-system/internal/cache"
	"survey-system/internal/config"
	"survey-system/internal/repository"
	"survey-system/internal/service"
	"survey-system/pkg/database"
	pkgRedis "survey-system/pkg/redis"
	"survey-system/pkg/utils"
)

func main() {
	// Parse command line flags
	configPath := flag.String("config", "./config/config.yaml", "path to config file")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Log successful configuration load
	log.Printf("Configuration loaded successfully")
	log.Printf("Server will run on port: %d", cfg.Server.Port)
	log.Printf("Server mode: %s", cfg.Server.Mode)
	log.Printf("Database: %s@%s:%d/%s", cfg.Database.Username, cfg.Database.Host, cfg.Database.Port, cfg.Database.Database)
	log.Printf("Redis: %s:%d", cfg.Redis.Host, cfg.Redis.Port)

	// Initialize database connection
	db, err := database.InitDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Run auto-migration
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to run database migration: %v", err)
	}

	// Initialize default admin account
	if err := database.InitializeDefaultAdmin(db); err != nil {
		log.Fatalf("Failed to initialize default admin: %v", err)
	}

	// Initialize Redis connection
	redisClient, err := pkgRedis.NewClient(&cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}
	log.Printf("Redis connection established successfully")

	// Create cache instance
	cacheInstance := cache.NewRedisCache(redisClient.GetClient())

	// Initialize encryption service
	encryptionSvc, err := service.NewEncryptionService(cfg.Encryption.Key)
	if err != nil {
		log.Fatalf("Failed to initialize encryption service: %v", err)
	}

	// Initialize repositories
	surveyRepo := repository.NewSurveyRepository(db)
	questionRepo := repository.NewQuestionRepository(db)
	oneLinkRepo := repository.NewOneLinkRepository(db)
	userRepo := repository.NewUserRepository(db)
	responseRepo := repository.NewResponseRepository(db)

	// Initialize JWT util
	jwtUtil := utils.NewJWTUtil(cfg.JWT.Secret, cfg.JWT.Expiration)

	// Initialize services
	surveyService := service.NewSurveyService(surveyRepo, cacheInstance)
	questionService := service.NewQuestionService(questionRepo, surveyRepo, cacheInstance)
	shareService := service.NewShareService(
		surveyRepo,
		questionRepo,
		oneLinkRepo,
		encryptionSvc,
		cacheInstance,
		cfg.OneLink.BaseURL,
		cfg.OneLink.DefaultExpiration,
		cfg.OneLink.MaxExpiration,
	)
	exportService := service.NewExportService(surveyRepo, questionRepo, responseRepo)
	responseService := service.NewResponseService(
		responseRepo,
		surveyRepo,
		questionRepo,
		oneLinkRepo,
		encryptionSvc,
		cacheInstance,
		exportService,
	)
	authService := service.NewAuthService(userRepo, jwtUtil)

	// Initialize handlers
	surveyHandler := handler.NewSurveyHandler(surveyService)
	questionHandler := handler.NewQuestionHandler(questionService)
	shareHandler := handler.NewShareHandler(shareService)
	responseHandler := handler.NewResponseHandler(responseService)
	authHandler := handler.NewAuthHandler(authService)

	// Setup router
	r := router.SetupRouter(
		surveyHandler,
		questionHandler,
		shareHandler,
		responseHandler,
		authHandler,
		jwtUtil,
		cfg,
		redisClient.GetClient(),
	)

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	// SIGINT handles Ctrl+C, SIGTERM handles termination signal
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown HTTP server gracefully
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	// Close database connection
	if err := database.Close(); err != nil {
		log.Printf("Error closing database connection: %v", err)
	}

	// Close Redis connection
	if err := redisClient.Close(); err != nil {
		log.Printf("Error closing Redis connection: %v", err)
	}

	log.Println("Server exited successfully")
}
