package response

import (
	"survey-system/internal/model"
	"time"
)

// QuestionResponse represents a question in API responses
type QuestionResponse struct {
	ID          uint                 `json:"id"`
	SurveyID    uint                 `json:"survey_id"`
	Type        string               `json:"type"`
	Title       string               `json:"title"`
	Description string               `json:"description"`
	Required    bool                 `json:"required"`
	Order       int                  `json:"order"`
	Config      model.QuestionConfig `json:"config"`
	PrefillKey  string               `json:"prefill_key,omitempty"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
}

// ToQuestionResponse converts a Question model to QuestionResponse
func ToQuestionResponse(question *model.Question) *QuestionResponse {
	return &QuestionResponse{
		ID:          question.ID,
		SurveyID:    question.SurveyID,
		Type:        question.Type,
		Title:       question.Title,
		Description: question.Description,
		Required:    question.Required,
		Order:       question.Order,
		Config:      question.Config,
		PrefillKey:  question.PrefillKey,
		CreatedAt:   question.CreatedAt,
		UpdatedAt:   question.UpdatedAt,
	}
}
