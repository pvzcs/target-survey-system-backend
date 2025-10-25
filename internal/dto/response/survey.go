package response

import (
	"survey-system/internal/model"
	"time"
)

// SurveyResponse represents a basic survey response
type SurveyResponse struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// SurveyDetailResponse represents a detailed survey response with questions
type SurveyDetailResponse struct {
	ID          uint               `json:"id"`
	UserID      uint               `json:"user_id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Status      string             `json:"status"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	Questions   []QuestionResponse `json:"questions"`
}

// PaginatedSurveyResponse represents a paginated list of surveys
type PaginatedSurveyResponse struct {
	Data []SurveyResponse `json:"data"`
	Meta PaginationMeta   `json:"meta"`
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	Page      int   `json:"page"`
	PageSize  int   `json:"page_size"`
	Total     int64 `json:"total"`
	TotalPage int   `json:"total_page"`
}

// ToSurveyResponse converts a model.Survey to SurveyResponse
func ToSurveyResponse(survey *model.Survey) *SurveyResponse {
	return &SurveyResponse{
		ID:          survey.ID,
		UserID:      survey.UserID,
		Title:       survey.Title,
		Description: survey.Description,
		Status:      survey.Status,
		CreatedAt:   survey.CreatedAt,
		UpdatedAt:   survey.UpdatedAt,
	}
}

// ToSurveyDetailResponse converts a model.Survey to SurveyDetailResponse
func ToSurveyDetailResponse(survey *model.Survey) *SurveyDetailResponse {
	questions := make([]QuestionResponse, len(survey.Questions))
	for i, q := range survey.Questions {
		questions[i] = *ToQuestionResponse(&q)
	}

	return &SurveyDetailResponse{
		ID:          survey.ID,
		UserID:      survey.UserID,
		Title:       survey.Title,
		Description: survey.Description,
		Status:      survey.Status,
		CreatedAt:   survey.CreatedAt,
		UpdatedAt:   survey.UpdatedAt,
		Questions:   questions,
	}
}
