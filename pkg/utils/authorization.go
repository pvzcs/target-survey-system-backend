package utils

import (
	"errors"
	"survey-system/internal/model"
	"survey-system/internal/repository"

	"gorm.io/gorm"
)

var (
	// ErrUnauthorized is returned when user is not authenticated
	ErrUnauthorized = errors.New("未授权访问")
	
	// ErrForbidden is returned when user doesn't have permission
	ErrForbidden = errors.New("禁止访问：您没有权限访问此资源")
	
	// ErrSurveyNotFound is returned when survey is not found
	ErrSurveyNotFound = errors.New("问卷不存在")
)

// AuthorizationUtil provides authorization checking utilities
type AuthorizationUtil struct {
	surveyRepo   repository.SurveyRepository
	questionRepo repository.QuestionRepository
}

// NewAuthorizationUtil creates a new authorization utility instance
func NewAuthorizationUtil(surveyRepo repository.SurveyRepository, questionRepo repository.QuestionRepository) *AuthorizationUtil {
	return &AuthorizationUtil{
		surveyRepo:   surveyRepo,
		questionRepo: questionRepo,
	}
}

// CheckSurveyOwnership verifies that the user owns the specified survey
func (a *AuthorizationUtil) CheckSurveyOwnership(userID, surveyID uint) error {
	survey, err := a.surveyRepo.FindByID(surveyID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrSurveyNotFound
		}
		return err
	}

	if survey.UserID != userID {
		return ErrForbidden
	}

	return nil
}

// CheckQuestionOwnership verifies that the user owns the survey containing the question
func (a *AuthorizationUtil) CheckQuestionOwnership(userID, questionID uint) (*model.Question, error) {
	question, err := a.questionRepo.FindByID(questionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("题目不存在")
		}
		return nil, err
	}

	// Check if user owns the survey
	if err := a.CheckSurveyOwnership(userID, question.SurveyID); err != nil {
		return nil, err
	}

	return question, nil
}

// GetSurveyIfOwned retrieves a survey only if the user owns it
func (a *AuthorizationUtil) GetSurveyIfOwned(userID, surveyID uint) (*model.Survey, error) {
	survey, err := a.surveyRepo.FindByID(surveyID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSurveyNotFound
		}
		return nil, err
	}

	if survey.UserID != userID {
		return nil, ErrForbidden
	}

	return survey, nil
}

// GetSurveyWithQuestionsIfOwned retrieves a survey with questions only if the user owns it
func (a *AuthorizationUtil) GetSurveyWithQuestionsIfOwned(userID, surveyID uint) (*model.Survey, error) {
	survey, err := a.surveyRepo.FindByIDWithQuestions(surveyID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSurveyNotFound
		}
		return nil, err
	}

	if survey.UserID != userID {
		return nil, ErrForbidden
	}

	return survey, nil
}
