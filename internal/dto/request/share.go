package request

import "time"

// GenerateShareLinkRequest represents the request to generate a share link
type GenerateShareLinkRequest struct {
	PrefillData map[string]interface{} `json:"prefill_data"` // Map of prefill_key to value
	ExpiresAt   *time.Time             `json:"expires_at"`   // Optional expiration time
}
