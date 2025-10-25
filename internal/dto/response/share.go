package response

import "time"

// ShareLinkResponse represents the response for a generated share link
type ShareLinkResponse struct {
	Token     string    `json:"token"`
	URL       string    `json:"url"`
	ExpiresAt time.Time `json:"expires_at"`
}

// SurveyWithPrefillResponse represents a survey with prefilled values
type SurveyWithPrefillResponse struct {
	ID          uint                   `json:"id"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Questions   []QuestionWithPrefill  `json:"questions"`
	PrefillData map[string]interface{} `json:"prefill_data"`
}

// QuestionWithPrefill represents a question with optional prefilled value
type QuestionWithPrefill struct {
	QuestionResponse
	PrefillValue interface{} `json:"prefill_value,omitempty"`
}
