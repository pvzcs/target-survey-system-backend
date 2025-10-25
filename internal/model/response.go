package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// Response represents a survey response/submission
type Response struct {
	ID          uint         `gorm:"primaryKey" json:"id"`
	SurveyID    uint         `gorm:"index;not null" json:"survey_id"`
	OneLinkID   uint         `gorm:"index" json:"one_link_id"`
	Data        ResponseData `gorm:"type:json;not null" json:"data"`
	IPAddress   string       `gorm:"size:45" json:"ip_address"`
	UserAgent   string       `gorm:"size:500" json:"user_agent"`
	SubmittedAt time.Time    `gorm:"not null;index" json:"submitted_at"`
	CreatedAt   time.Time    `json:"created_at"`
	
	// Associations
	Survey  Survey  `gorm:"foreignKey:SurveyID" json:"survey,omitempty"`
	OneLink OneLink `gorm:"foreignKey:OneLinkID" json:"one_link,omitempty"`
}

// TableName specifies the table name for Response model
func (Response) TableName() string {
	return "responses"
}

// ResponseData holds the actual response data
type ResponseData struct {
	Answers []Answer `json:"answers"`
}

// Answer represents an answer to a single question
type Answer struct {
	QuestionID uint        `json:"question_id"`
	Value      interface{} `json:"value"` // string, []string, or []map[string]interface{} for table
}

// Scan implements the sql.Scanner interface for ResponseData
func (r *ResponseData) Scan(value interface{}) error {
	if value == nil {
		*r = ResponseData{}
		return nil
	}
	
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal ResponseData value: %v", value)
	}
	
	return json.Unmarshal(bytes, r)
}

// Value implements the driver.Valuer interface for ResponseData
func (r ResponseData) Value() (driver.Value, error) {
	if r.Answers == nil {
		return json.Marshal(ResponseData{Answers: []Answer{}})
	}
	return json.Marshal(r)
}
