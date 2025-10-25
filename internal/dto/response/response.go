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
	ID          uint      `json:"id"`
	SurveyID    uint      `json:"survey_id"`
	IPAddress   string    `json:"ip_address"`
	SubmittedAt time.Time `json:"submitted_at"`
}

// PaginatedResponseResponse represents paginated response list
type PaginatedResponseResponse struct {
	Responses []ResponseListItem `json:"responses"`
	Total     int64              `json:"total"`
	Page      int                `json:"page"`
	PageSize  int                `json:"page_size"`
}

// StatisticsResponse represents survey statistics
type StatisticsResponse struct {
	SurveyID       uint    `json:"survey_id"`
	TotalResponses int64   `json:"total_responses"`
	CompletionRate float64 `json:"completion_rate"`
}
