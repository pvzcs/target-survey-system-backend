package router

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"survey-system/internal/api/handler"
	"survey-system/internal/api/middleware"
	"survey-system/internal/config"
	"survey-system/pkg/utils"
)

// SetupRouter configures all routes for the application
func SetupRouter(
	surveyHandler *handler.SurveyHandler,
	questionHandler *handler.QuestionHandler,
	shareHandler *handler.ShareHandler,
	responseHandler *handler.ResponseHandler,
	authHandler *handler.AuthHandler,
	jwtUtil *utils.JWTUtil,
	cfg *config.Config,
	redisClient *redis.Client,
) *gin.Engine {
	router := gin.Default()

	// Apply global middleware
	router.Use(middleware.CORS(cfg))
	router.Use(middleware.RateLimit(redisClient, cfg.RateLimit.RequestsPerMinute))

	// Create auth middleware
	authMiddleware := middleware.AuthMiddleware(jwtUtil)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes (public, no authentication required)
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
		}
		// Survey routes (protected)
		surveys := v1.Group("/surveys")
		surveys.Use(authMiddleware)
		{
			surveys.POST("", surveyHandler.CreateSurvey)
			surveys.GET("", surveyHandler.ListSurveys)
			surveys.GET("/:id", surveyHandler.GetSurvey)
			surveys.PUT("/:id", surveyHandler.UpdateSurvey)
			surveys.DELETE("/:id", surveyHandler.DeleteSurvey)
			surveys.POST("/:id/publish", surveyHandler.PublishSurvey)
			
			// Share link generation (protected)
			surveys.POST("/:id/share", shareHandler.GenerateShareLink)
			
			// Response management routes (protected)
			surveys.GET("/:id/responses", responseHandler.GetResponses)
			surveys.GET("/:id/statistics", responseHandler.GetStatistics)
			surveys.GET("/:id/export", responseHandler.ExportResponses)
			
			// Question reorder route (nested under surveys)
			surveys.PUT("/:id/questions/reorder", questionHandler.ReorderQuestions)
		}

		// Question routes (protected)
		questions := v1.Group("/questions")
		questions.Use(authMiddleware)
		{
			questions.POST("", questionHandler.CreateQuestion)
			questions.PUT("/:id", questionHandler.UpdateQuestion)
			questions.DELETE("/:id", questionHandler.DeleteQuestion)
		}

		// Public routes (no authentication required)
		public := v1.Group("/public")
		{
			// Get survey by token (public access for respondents)
			public.GET("/surveys/:id", shareHandler.GetSurveyByToken)
			
			// Submit response (public access for respondents)
			public.POST("/responses", responseHandler.SubmitResponse)
		}
	}

	return router
}
