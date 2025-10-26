package request

import "survey-system/internal/model"

// CreateQuestionRequest represents the request to create a question
type CreateQuestionRequest struct {
	SurveyID    uint                 `json:"survey_id" binding:"required"`
	Type        string               `json:"type" binding:"required,oneof=text single multiple table"`
	Title       string               `json:"title" binding:"required,max=500"`
	Description string               `json:"description" binding:"max=5000"`
	Required    bool                 `json:"required"`
	Order       *int                 `json:"order" binding:"required,min=0"`
	Config      model.QuestionConfig `json:"config"`
	PrefillKey  string               `json:"prefill_key" binding:"max=100"`
}

// UpdateQuestionRequest represents the request to update a question
type UpdateQuestionRequest struct {
	Type        string               `json:"type" binding:"required,oneof=text single multiple table"`
	Title       string               `json:"title" binding:"required,max=500"`
	Description string               `json:"description" binding:"max=5000"`
	Required    bool                 `json:"required"`
	Order       *int                 `json:"order" binding:"required,min=0"`
	Config      model.QuestionConfig `json:"config"`
	PrefillKey  string               `json:"prefill_key" binding:"max=100"`
}

// ReorderQuestionsRequest represents the request to reorder questions
type ReorderQuestionsRequest struct {
	QuestionIDs []uint `json:"question_ids" binding:"required,min=1"`
}
