package response

import "time"

// LoginResponse represents the response after successful login
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// UserResponse represents user information in responses
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// RegisterResponse represents the response after successful registration
type RegisterResponse struct {
	Message string `json:"message"`
}
