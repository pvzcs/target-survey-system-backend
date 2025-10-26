package service

import (
	"context"
	"fmt"
	"time"

	"survey-system/internal/cache"
	"survey-system/internal/dto/request"
	"survey-system/internal/dto/response"
	"survey-system/internal/model"
	"survey-system/internal/repository"
	"survey-system/pkg/errors"

	"gorm.io/gorm"
)

// SurveyService defines the interface for survey business logic
type SurveyService interface {
	CreateSurvey(ctx context.Context, userID uint, req *request.CreateSurveyRequest) (*response.SurveyResponse, error)
	UpdateSurvey(ctx context.Context, userID, surveyID uint, req *request.UpdateSurveyRequest) (*response.SurveyResponse, error)
	DeleteSurvey(ctx context.Context, userID, surveyID uint) error
	GetSurvey(ctx context.Context, surveyID uint) (*response.SurveyDetailResponse, error)
	ListSurveys(ctx context.Context, userID uint, page, pageSize int) (*response.PaginatedSurveyResponse, error)
	PublishSurvey(ctx context.Context, userID, surveyID uint) error
}

// surveyService implements SurveyService interface
type surveyService struct {
	surveyRepo repository.SurveyRepository
	cache      cache.Cache
}

// NewSurveyService creates a new survey service instance
func NewSurveyService(surveyRepo repository.SurveyRepository, cache cache.Cache) SurveyService {
	return &surveyService{
		surveyRepo: surveyRepo,
		cache:      cache,
	}
}

// CreateSurvey creates a new survey with draft status
func (s *surveyService) CreateSurvey(ctx context.Context, userID uint, req *request.CreateSurveyRequest) (*response.SurveyResponse, error) {
	survey := &model.Survey{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Status:      model.SurveyStatusDraft,
	}

	if err := s.surveyRepo.Create(survey); err != nil {
		return nil, errors.WrapError(err, "failed to create survey")
	}

	return response.ToSurveyResponse(survey), nil
}

// UpdateSurvey updates an existing survey after verifying ownership
func (s *surveyService) UpdateSurvey(ctx context.Context, userID, surveyID uint, req *request.UpdateSurveyRequest) (*response.SurveyResponse, error) {
	// Find the survey
	survey, err := s.surveyRepo.FindByID(surveyID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, errors.WrapError(err, "failed to find survey")
	}

	// Verify ownership
	if survey.UserID != userID {
		return nil, errors.ErrForbidden
	}

	// Update fields
	survey.Title = req.Title
	survey.Description = req.Description

	if err := s.surveyRepo.Update(survey); err != nil {
		return nil, errors.WrapError(err, "failed to update survey")
	}

	// Invalidate cache
	if err := s.cache.DeleteSurvey(ctx, surveyID); err != nil {
		// Log error but don't fail the request
		fmt.Printf("failed to invalidate survey cache: %v\n", err)
	}

	return response.ToSurveyResponse(survey), nil
}

// DeleteSurvey deletes a survey after verifying ownership
// If cascade delete fails due to foreign key constraints, manually deletes associated data
func (s *surveyService) DeleteSurvey(ctx context.Context, userID, surveyID uint) error {
	// Find the survey
	survey, err := s.surveyRepo.FindByID(surveyID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.ErrNotFound
		}
		return errors.WrapError(err, "failed to find survey")
	}

	// Verify ownership
	if survey.UserID != userID {
		return errors.ErrForbidden
	}

	// Delete the survey (cascade delete handled by database)
	if err := s.surveyRepo.Delete(surveyID); err != nil {
		return errors.WrapError(err, "failed to delete survey")
	}

	// Invalidate cache
	if err := s.cache.DeleteSurvey(ctx, surveyID); err != nil {
		// Log error but don't fail the request
		fmt.Printf("failed to invalidate survey cache: %v\n", err)
	}

	return nil
}

// GetSurvey retrieves survey details with questions, using cache when available
func (s *surveyService) GetSurvey(ctx context.Context, surveyID uint) (*response.SurveyDetailResponse, error) {
	// Try to get from cache first
	cachedSurvey, err := s.cache.GetSurvey(ctx, surveyID)
	if err != nil {
		// Log error but continue to database
		fmt.Printf("failed to get survey from cache: %v\n", err)
	}

	if cachedSurvey != nil {
		return response.ToSurveyDetailResponse(cachedSurvey), nil
	}

	// Cache miss, get from database
	survey, err := s.surveyRepo.FindByIDWithQuestions(surveyID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, errors.WrapError(err, "failed to find survey")
	}

	// Cache the survey for 1 hour
	if err := s.cache.SetSurvey(ctx, survey, time.Hour); err != nil {
		// Log error but don't fail the request
		fmt.Printf("failed to cache survey: %v\n", err)
	}

	return response.ToSurveyDetailResponse(survey), nil
}

// ListSurveys retrieves a paginated list of surveys for a user
func (s *surveyService) ListSurveys(ctx context.Context, userID uint, page, pageSize int) (*response.PaginatedSurveyResponse, error) {
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	surveys, total, err := s.surveyRepo.FindByUserID(userID, page, pageSize)
	if err != nil {
		return nil, errors.WrapError(err, "failed to list surveys")
	}

	// Convert to response DTOs
	surveyResponses := make([]response.SurveyResponse, len(surveys))
	for i, survey := range surveys {
		surveyResponses[i] = *response.ToSurveyResponse(&survey)
	}

	// Calculate total pages
	totalPage := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPage++
	}

	return &response.PaginatedSurveyResponse{
		Data: surveyResponses,
		Meta: response.PaginationMeta{
			Page:      page,
			PageSize:  pageSize,
			Total:     total,
			TotalPage: totalPage,
		},
	}, nil
}

// PublishSurvey publishes a survey after verifying ownership
func (s *surveyService) PublishSurvey(ctx context.Context, userID, surveyID uint) error {
	// Find the survey
	survey, err := s.surveyRepo.FindByID(surveyID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.ErrNotFound
		}
		return errors.WrapError(err, "failed to find survey")
	}

	// Verify ownership
	if survey.UserID != userID {
		return errors.ErrForbidden
	}

	// Update status to published
	if err := s.surveyRepo.UpdateStatus(surveyID, model.SurveyStatusPublished); err != nil {
		return errors.WrapError(err, "failed to publish survey")
	}

	// Invalidate cache
	if err := s.cache.DeleteSurvey(ctx, surveyID); err != nil {
		// Log error but don't fail the request
		fmt.Printf("failed to invalidate survey cache: %v\n", err)
	}

	return nil
}
