package request

// LoginRequest represents the request to login
type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
}

// RegisterRequest represents the request to register a new user
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"omitempty,email,max=100"`
}

// UpdateProfileRequest represents the request to update user profile
type UpdateProfileRequest struct {
	Username    string `json:"username" binding:"omitempty,min=3,max=50"`
	Email       string `json:"email" binding:"omitempty,email,max=100"`
	OldPassword string `json:"old_password" binding:"omitempty,min=6"`
	NewPassword string `json:"new_password" binding:"omitempty,min=6"`
}
