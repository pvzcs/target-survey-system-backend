package database

import (
	"fmt"
	"log"

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
