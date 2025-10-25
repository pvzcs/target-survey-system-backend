package model

import "time"

// User represents a user in the system
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Password  string    `gorm:"size:255;not null" json:"-"` // bcrypt hashed, never expose in JSON
	Email     string    `gorm:"uniqueIndex;size:100" json:"email"`
	Role      string    `gorm:"size:20;default:'admin'" json:"role"` // admin
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "users"
}
