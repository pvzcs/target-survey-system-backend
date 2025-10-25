package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// Question represents a question in a survey
type Question struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	SurveyID    uint           `gorm:"index;not null" json:"survey_id"`
	Type        string         `gorm:"size:20;not null" json:"type"` // text, single, multiple, table
	Title       string         `gorm:"size:500;not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	Required    bool           `gorm:"default:false" json:"required"`
	Order       int            `gorm:"not null" json:"order"`
	Config      QuestionConfig `gorm:"type:json" json:"config"`
	PrefillKey  string         `gorm:"size:100" json:"prefill_key"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	
	// Associations
	Survey Survey `gorm:"foreignKey:SurveyID" json:"survey,omitempty"`
}

// TableName specifies the table name for Question model
func (Question) TableName() string {
	return "questions"
}

// Question type constants
const (
	QuestionTypeText     = "text"
	QuestionTypeSingle   = "single"
	QuestionTypeMultiple = "multiple"
	QuestionTypeTable    = "table"
)

// QuestionConfig holds the configuration for different question types
type QuestionConfig struct {
	// For single/multiple choice questions
	Options []string `json:"options,omitempty"`
	
	// For table questions
	Columns   []TableColumn `json:"columns,omitempty"`
	MinRows   int           `json:"min_rows,omitempty"`
	MaxRows   int           `json:"max_rows,omitempty"`
	CanAddRow bool          `json:"can_add_row,omitempty"`
}

// TableColumn represents a column in a table question
type TableColumn struct {
	ID      string   `json:"id"`
	Type    string   `json:"type"`    // text, number, select
	Label   string   `json:"label"`
	Options []string `json:"options,omitempty"` // for select type
}

// Scan implements the sql.Scanner interface for QuestionConfig
func (c *QuestionConfig) Scan(value interface{}) error {
	if value == nil {
		*c = QuestionConfig{}
		return nil
	}
	
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal QuestionConfig value: %v", value)
	}
	
	return json.Unmarshal(bytes, c)
}

// Value implements the driver.Valuer interface for QuestionConfig
func (c QuestionConfig) Value() (driver.Value, error) {
	if c.Options == nil && c.Columns == nil {
		return nil, nil
	}
	return json.Marshal(c)
}
