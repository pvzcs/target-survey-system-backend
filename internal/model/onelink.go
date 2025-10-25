package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// OneLink represents a one-time access link for a survey
type OneLink struct {
	ID          uint                   `gorm:"primaryKey" json:"id"`
	SurveyID    uint                   `gorm:"index;not null" json:"survey_id"`
	Token       string                 `gorm:"uniqueIndex;size:500;not null" json:"token"` // Encrypted token
	PrefillData map[string]interface{} `gorm:"type:json" json:"prefill_data"`              // JSON prefill values
	ExpiresAt   time.Time              `gorm:"index;not null" json:"expires_at"`
	Used        bool                   `gorm:"default:false;index" json:"used"`
	UsedAt      *time.Time             `json:"used_at"`
	AccessedAt  *time.Time             `json:"accessed_at"`
	CreatedAt   time.Time              `json:"created_at"`
	
	// Associations
	Survey Survey `gorm:"foreignKey:SurveyID" json:"survey,omitempty"`
}

// TableName specifies the table name for OneLink model
func (OneLink) TableName() string {
	return "one_links"
}

// IsExpired checks if the link has expired
func (o *OneLink) IsExpired() bool {
	return time.Now().After(o.ExpiresAt)
}

// IsValid checks if the link is valid (not used and not expired)
func (o *OneLink) IsValid() bool {
	return !o.Used && !o.IsExpired()
}

// PrefillDataType is a custom type for handling JSON prefill data
type PrefillDataType map[string]interface{}

// Scan implements the sql.Scanner interface for PrefillDataType
func (p *PrefillDataType) Scan(value interface{}) error {
	if value == nil {
		*p = make(map[string]interface{})
		return nil
	}
	
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal PrefillDataType value: %v", value)
	}
	
	return json.Unmarshal(bytes, p)
}

// Value implements the driver.Valuer interface for PrefillDataType
func (p PrefillDataType) Value() (driver.Value, error) {
	if p == nil || len(p) == 0 {
		return nil, nil
	}
	return json.Marshal(p)
}
