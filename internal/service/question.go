package service

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"survey-system/internal/cache"
	"survey-system/internal/dto/request"
	"survey-system/internal/dto/response"
	"survey-system/internal/model"
	"survey-system/internal/repository"
	"survey-system/pkg/errors"
)

// QuestionService defines the interface for question business logic
type QuestionService interface {
	CreateQuestion(ctx context.Context, userID uint, req *request.CreateQuestionRequest) (*response.QuestionResponse, error)
	UpdateQuestion(ctx context.Context, userID, questionID uint, req *request.UpdateQuestionRequest) (*response.QuestionResponse, error)
	DeleteQuestion(ctx context.Context, userID, questionID uint) error
	ReorderQuestions(ctx context.Context, userID, surveyID uint, questionIDs []uint) error
}

// questionService implements QuestionService interface
type questionService struct {
	questionRepo repository.QuestionRepository
	surveyRepo   repository.SurveyRepository
	cache        cache.Cache
}

// NewQuestionService creates a new question service instance
func NewQuestionService(
	questionRepo repository.QuestionRepository,
	surveyRepo repository.SurveyRepository,
	cache cache.Cache,
) QuestionService {
	return &questionService{
		questionRepo: questionRepo,
		surveyRepo:   surveyRepo,
		cache:        cache,
	}
}

// CreateQuestion creates a new question after verifying survey ownership and validating configuration
func (s *questionService) CreateQuestion(ctx context.Context, userID uint, req *request.CreateQuestionRequest) (*response.QuestionResponse, error) {
	// Verify survey exists and user owns it
	survey, err := s.surveyRepo.FindByID(req.SurveyID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, errors.WrapError(err, "failed to find survey")
	}

	if survey.UserID != userID {
		return nil, errors.ErrForbidden
	}

	// Validate question configuration based on type
	if err := s.validateQuestionConfig(req.Type, &req.Config); err != nil {
		return nil, err
	}

	// Create the question
	question := &model.Question{
		SurveyID:    req.SurveyID,
		Type:        req.Type,
		Title:       req.Title,
		Description: req.Description,
		Required:    req.Required,
		Order:       req.Order,
		Config:      req.Config,
		PrefillKey:  req.PrefillKey,
	}

	if err := s.questionRepo.Create(question); err != nil {
		return nil, errors.WrapError(err, "failed to create question")
	}

	// Invalidate survey cache since questions changed
	if err := s.cache.DeleteSurvey(ctx, req.SurveyID); err != nil {
		fmt.Printf("failed to invalidate survey cache: %v\n", err)
	}

	return response.ToQuestionResponse(question), nil
}

// UpdateQuestion updates an existing question after verifying ownership and validating configuration
func (s *questionService) UpdateQuestion(ctx context.Context, userID, questionID uint, req *request.UpdateQuestionRequest) (*response.QuestionResponse, error) {
	// Find the question
	question, err := s.questionRepo.FindByID(questionID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, errors.WrapError(err, "failed to find question")
	}

	// Verify survey ownership
	survey, err := s.surveyRepo.FindByID(question.SurveyID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, errors.WrapError(err, "failed to find survey")
	}

	if survey.UserID != userID {
		return nil, errors.ErrForbidden
	}

	// Validate question configuration based on type
	if err := s.validateQuestionConfig(req.Type, &req.Config); err != nil {
		return nil, err
	}

	// Update fields
	question.Type = req.Type
	question.Title = req.Title
	question.Description = req.Description
	question.Required = req.Required
	question.Order = req.Order
	question.Config = req.Config
	question.PrefillKey = req.PrefillKey

	if err := s.questionRepo.Update(question); err != nil {
		return nil, errors.WrapError(err, "failed to update question")
	}

	// Invalidate survey cache
	if err := s.cache.DeleteSurvey(ctx, question.SurveyID); err != nil {
		fmt.Printf("failed to invalidate survey cache: %v\n", err)
	}

	return response.ToQuestionResponse(question), nil
}

// DeleteQuestion deletes a question after verifying ownership
func (s *questionService) DeleteQuestion(ctx context.Context, userID, questionID uint) error {
	// Find the question
	question, err := s.questionRepo.FindByID(questionID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.ErrNotFound
		}
		return errors.WrapError(err, "failed to find question")
	}

	// Verify survey ownership
	survey, err := s.surveyRepo.FindByID(question.SurveyID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.ErrNotFound
		}
		return errors.WrapError(err, "failed to find survey")
	}

	if survey.UserID != userID {
		return errors.ErrForbidden
	}

	// Delete the question
	if err := s.questionRepo.Delete(questionID); err != nil {
		return errors.WrapError(err, "failed to delete question")
	}

	// Invalidate survey cache
	if err := s.cache.DeleteSurvey(ctx, question.SurveyID); err != nil {
		fmt.Printf("failed to invalidate survey cache: %v\n", err)
	}

	return nil
}

// ReorderQuestions updates the order of questions in a survey
func (s *questionService) ReorderQuestions(ctx context.Context, userID, surveyID uint, questionIDs []uint) error {
	// Verify survey ownership
	survey, err := s.surveyRepo.FindByID(surveyID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.ErrNotFound
		}
		return errors.WrapError(err, "failed to find survey")
	}

	if survey.UserID != userID {
		return errors.ErrForbidden
	}

	// Get all questions for this survey
	questions, err := s.questionRepo.FindBySurveyID(surveyID)
	if err != nil {
		return errors.WrapError(err, "failed to find questions")
	}

	// Validate that all question IDs belong to this survey
	questionMap := make(map[uint]*model.Question)
	for i := range questions {
		questionMap[questions[i].ID] = &questions[i]
	}

	// Build the list of questions to update with new order
	questionsToUpdate := make([]model.Question, 0, len(questionIDs))
	for order, questionID := range questionIDs {
		question, exists := questionMap[questionID]
		if !exists {
			return errors.NewValidationError("question_id", fmt.Sprintf("question %d does not belong to survey %d", questionID, surveyID))
		}
		
		// Create a copy with updated order
		updatedQuestion := *question
		updatedQuestion.Order = order
		questionsToUpdate = append(questionsToUpdate, updatedQuestion)
	}

	// Batch update the order
	if err := s.questionRepo.BatchUpdateOrder(questionsToUpdate); err != nil {
		return errors.WrapError(err, "failed to reorder questions")
	}

	// Invalidate survey cache
	if err := s.cache.DeleteSurvey(ctx, surveyID); err != nil {
		fmt.Printf("failed to invalidate survey cache: %v\n", err)
	}

	return nil
}

// validateQuestionConfig validates the question configuration based on question type
func (s *questionService) validateQuestionConfig(questionType string, config *model.QuestionConfig) error {
	switch questionType {
	case model.QuestionTypeText:
		// Text questions don't need special configuration
		return nil

	case model.QuestionTypeSingle, model.QuestionTypeMultiple:
		// Single and multiple choice questions must have options
		if len(config.Options) == 0 {
			return errors.NewValidationError("config.options", "single and multiple choice questions must have at least one option")
		}
		return nil

	case model.QuestionTypeTable:
		// Table questions must have column definitions
		if len(config.Columns) == 0 {
			return errors.NewValidationError("config.columns", "table questions must have at least one column")
		}

		// Validate each column
		for i, col := range config.Columns {
			if col.ID == "" {
				return errors.NewValidationError(fmt.Sprintf("config.columns[%d].id", i), "column ID is required")
			}
			if col.Type == "" {
				return errors.NewValidationError(fmt.Sprintf("config.columns[%d].type", i), "column type is required")
			}
			if col.Type != "text" && col.Type != "number" && col.Type != "select" {
				return errors.NewValidationError(fmt.Sprintf("config.columns[%d].type", i), "column type must be text, number, or select")
			}
			if col.Label == "" {
				return errors.NewValidationError(fmt.Sprintf("config.columns[%d].label", i), "column label is required")
			}
			// If column type is select, it must have options
			if col.Type == "select" && len(col.Options) == 0 {
				return errors.NewValidationError(fmt.Sprintf("config.columns[%d].options", i), "select columns must have at least one option")
			}
		}

		// Validate row constraints
		if config.MinRows < 0 {
			return errors.NewValidationError("config.min_rows", "min_rows cannot be negative")
		}
		if config.MaxRows < 0 {
			return errors.NewValidationError("config.max_rows", "max_rows cannot be negative")
		}
		if config.MinRows > 0 && config.MaxRows > 0 && config.MinRows > config.MaxRows {
			return errors.NewValidationError("config.min_rows", "min_rows cannot be greater than max_rows")
		}

		return nil

	default:
		return errors.NewValidationError("type", fmt.Sprintf("invalid question type: %s", questionType))
	}
}
