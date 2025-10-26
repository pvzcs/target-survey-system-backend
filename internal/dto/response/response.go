package response

import "time"

// SubmitResponseResponse represents the response after submitting a survey response
type SubmitResponseResponse struct {
	ID          uint      `json:"id"`
	SurveyID    uint      `json:"survey_id"`
	SubmittedAt time.Time `json:"submitted_at"`
	Message     string    `json:"message"`
}

// ResponseListItem represents a single response in the list
type ResponseListItem struct {
	ID          uint                   `json:"id"`
	SurveyID    uint                   `json:"survey_id"`
	Data        map[string]interface{} `json:"data"`
	IPAddress   string                 `json:"ip_address"`
	UserAgent   string                 `json:"user_agent"`
	SubmittedAt time.Time              `json:"submitted_at"`
	CreatedAt   time.Time              `json:"created_at"`
}

// PaginatedResponseMeta represents pagination metadata
type PaginatedResponseMeta struct {
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
}

// StatisticsResponse represents survey statistics
type StatisticsResponse struct {
	SurveyID       uint    `json:"survey_id"`
	TotalResponses int64   `json:"total_responses"`
	CompletionRate float64 `json:"completion_rate"`
}
