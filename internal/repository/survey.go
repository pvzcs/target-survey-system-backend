package repository

import (
	"survey-system/internal/model"

	"gorm.io/gorm"
)

// SurveyRepository defines the interface for survey data operations
type SurveyRepository interface {
	Create(survey *model.Survey) error
	Update(survey *model.Survey) error
	Delete(id uint) error
	FindByID(id uint) (*model.Survey, error)
	FindByIDWithQuestions(id uint) (*model.Survey, error)
	FindByUserID(userID uint, page, pageSize int) ([]model.Survey, int64, error)
	UpdateStatus(id uint, status string) error
}

// surveyRepository implements SurveyRepository interface
type surveyRepository struct {
	db *gorm.DB
}

// NewSurveyRepository creates a new survey repository instance
func NewSurveyRepository(db *gorm.DB) SurveyRepository {
	return &surveyRepository{db: db}
}

// Create creates a new survey
func (r *surveyRepository) Create(survey *model.Survey) error {
	return r.db.Create(survey).Error
}

// Update updates an existing survey
func (r *surveyRepository) Update(survey *model.Survey) error {
	return r.db.Save(survey).Error
}

// Delete deletes a survey by ID (cascade delete handled by database)
func (r *surveyRepository) Delete(id uint) error {
	return r.db.Delete(&model.Survey{}, id).Error
}

// FindByID finds a survey by ID without preloading questions
func (r *surveyRepository) FindByID(id uint) (*model.Survey, error) {
	var survey model.Survey
	err := r.db.First(&survey, id).Error
	if err != nil {
		return nil, err
	}
	return &survey, nil
}

// FindByIDWithQuestions finds a survey by ID with preloaded questions
func (r *surveyRepository) FindByIDWithQuestions(id uint) (*model.Survey, error) {
	var survey model.Survey
	err := r.db.Preload("Questions", func(db *gorm.DB) *gorm.DB {
		return db.Order("questions.order ASC")
	}).First(&survey, id).Error
	if err != nil {
		return nil, err
	}
	return &survey, nil
}

// FindByUserID finds surveys by user ID with pagination
func (r *surveyRepository) FindByUserID(userID uint, page, pageSize int) ([]model.Survey, int64, error) {
	var surveys []model.Survey
	var total int64

	// Count total records
	if err := r.db.Model(&model.Survey{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Query with pagination
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&surveys).Error

	if err != nil {
		return nil, 0, err
	}

	return surveys, total, nil
}

// UpdateStatus updates the status of a survey
func (r *surveyRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&model.Survey{}).Where("id = ?", id).Update("status", status).Error
}
