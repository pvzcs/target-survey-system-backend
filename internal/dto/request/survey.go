package request

// CreateSurveyRequest represents the request to create a survey
type CreateSurveyRequest struct {
	Title       string `json:"title" binding:"required,max=200"`
	Description string `json:"description" binding:"max=5000"`
}

// UpdateSurveyRequest represents the request to update a survey
type UpdateSurveyRequest struct {
	Title       string `json:"title" binding:"required,max=200"`
	Description string `json:"description" binding:"max=5000"`
}
