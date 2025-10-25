package repository

import (
	"survey-system/internal/model"

	"gorm.io/gorm"
)

// ResponseRepository defines the interface for response data operations
type ResponseRepository interface {
	Create(response *model.Response) error
	FindByID(id uint) (*model.Response, error)
	FindBySurveyID(surveyID uint, page, pageSize int) ([]model.Response, int64, error)
	CountBySurveyID(surveyID uint) (int64, error)
}

// responseRepository implements ResponseRepository interface
type responseRepository struct {
	db *gorm.DB
}

// NewResponseRepository creates a new response repository instance
func NewResponseRepository(db *gorm.DB) ResponseRepository {
	return &responseRepository{db: db}
}

// Create creates a new response record
func (r *responseRepository) Create(response *model.Response) error {
	return r.db.Create(response).Error
}

// FindByID finds a response by ID
func (r *responseRepository) FindByID(id uint) (*model.Response, error) {
	var response model.Response
	err := r.db.First(&response, id).Error
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// FindBySurveyID finds all responses for a survey with pagination
func (r *responseRepository) FindBySurveyID(surveyID uint, page, pageSize int) ([]model.Response, int64, error) {
	var responses []model.Response
	var total int64

	// Count total records
	if err := r.db.Model(&model.Response{}).Where("survey_id = ?", surveyID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Query with pagination
	err := r.db.Where("survey_id = ?", surveyID).
		Order("submitted_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&responses).Error

	if err != nil {
		return nil, 0, err
	}

	return responses, total, nil
}

// CountBySurveyID counts the total number of responses for a survey
func (r *responseRepository) CountBySurveyID(surveyID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Response{}).Where("survey_id = ?", surveyID).Count(&count).Error
	return count, err
}
