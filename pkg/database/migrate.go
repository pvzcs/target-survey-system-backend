package database

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"survey-system/internal/model"
)

// AutoMigrate runs automatic migration for all models
func AutoMigrate(db *gorm.DB) error {
	log.Println("Starting database auto-migration...")

	// List of all models to migrate
	models := []interface{}{
		&model.User{},
		&model.Survey{},
		&model.Question{},
		&model.Response{},
		&model.OneLink{},
	}

	// Run auto-migration for each model
	for _, m := range models {
		if err := db.AutoMigrate(m); err != nil {
			return fmt.Errorf("failed to migrate model %T: %w", m, err)
		}
		log.Printf("Successfully migrated model: %T", m)
	}

	log.Println("Database auto-migration completed successfully")
	return nil
}

// DropAllTables drops all tables (use with caution, mainly for testing)
func DropAllTables(db *gorm.DB) error {
	log.Println("Dropping all tables...")

	// Drop tables in reverse order to respect foreign key constraints
	models := []interface{}{
		&model.OneLink{},
		&model.Response{},
		&model.Question{},
		&model.Survey{},
		&model.User{},
	}

	for _, m := range models {
		if err := db.Migrator().DropTable(m); err != nil {
			return fmt.Errorf("failed to drop table for model %T: %w", m, err)
		}
		log.Printf("Successfully dropped table for model: %T", m)
	}

	log.Println("All tables dropped successfully")
	return nil
}

// CreateIndexes creates additional indexes that might not be created by AutoMigrate
func CreateIndexes(db *gorm.DB) error {
	log.Println("Creating additional indexes...")

	// Create composite indexes if needed
	// Example: composite index on survey_id and order for questions
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_questions_survey_order ON questions(survey_id, `order`)").Error; err != nil {
		log.Printf("Warning: failed to create composite index on questions: %v", err)
	}

	log.Println("Additional indexes created successfully")
	return nil
}

// InitializeDefaultAdmin creates a default admin account if no users exist
func InitializeDefaultAdmin(db *gorm.DB) error {
	log.Println("Checking for existing users...")

	// Check if any users exist
	var count int64
	if err := db.Model(&model.User{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to count users: %w", err)
	}

	// If users already exist, skip initialization
	if count > 0 {
		log.Printf("Found %d existing user(s), skipping default admin creation", count)
		return nil
	}

	log.Println("No users found, creating default admin account...")

	// Hash the default password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Create default admin user
	defaultAdmin := &model.User{
		Username: "admin",
		Password: string(hashedPassword),
		Email:    "admin@example.com",
		Role:     "admin",
	}

	if err := db.Create(defaultAdmin).Error; err != nil {
		return fmt.Errorf("failed to create default admin: %w", err)
	}

	log.Println("✓ Default admin account created successfully")
	log.Println("  Username: admin")
	log.Println("  Password: admin123")
	log.Println("  Email: admin@example.com")
	log.Println("  ⚠️  Please change the default password after first login!")

	return nil
}
