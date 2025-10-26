package model

import "time"

// Survey represents a survey/questionnaire
type Survey struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"index;not null" json:"user_id"`
	Title       string    `gorm:"size:200;not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	Status      string    `gorm:"size:20;default:'draft';index" json:"status"` // draft, published
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Associations
	User      User       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Questions []Question `gorm:"foreignKey:SurveyID;constraint:OnDelete:CASCADE" json:"questions,omitempty"`
	OneLinks  []OneLink  `gorm:"foreignKey:SurveyID;constraint:OnDelete:CASCADE" json:"one_links,omitempty"`
	Responses []Response `gorm:"foreignKey:SurveyID;constraint:OnDelete:CASCADE" json:"responses,omitempty"`
}

// TableName specifies the table name for Survey model
func (Survey) TableName() string {
	return "surveys"
}

// Survey status constants
const (
	SurveyStatusDraft     = "draft"
	SurveyStatusPublished = "published"
)
