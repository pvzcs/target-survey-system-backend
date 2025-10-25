package request

// SubmitResponseRequest represents the request to submit a survey response
type SubmitResponseRequest struct {
	Token   string                   `json:"token" binding:"required"`
	Answers []AnswerRequest          `json:"answers" binding:"required,min=1"`
}

// AnswerRequest represents an answer to a single question
type AnswerRequest struct {
	QuestionID uint        `json:"question_id" binding:"required"`
	Value      interface{} `json:"value" binding:"required"`
}
