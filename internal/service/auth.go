package service

import (
	"errors"
	"survey-system/internal/model"
	"survey-system/internal/repository"
	"survey-system/pkg/utils"

	"gorm.io/gorm"
)

// AuthService defines the interface for authentication operations
type AuthService interface {
	Login(username, password string) (*LoginResponse, error)
	Register(username, password, email string) error
	ValidateToken(token string) (*utils.JWTClaims, error)
}

// LoginResponse represents the response after successful login
type LoginResponse struct {
	Token string       `json:"token"`
	User  *model.User  `json:"user"`
}

// authService implements AuthService interface
type authService struct {
	userRepo repository.UserRepository
	jwtUtil  *utils.JWTUtil
}

// NewAuthService creates a new auth service instance
func NewAuthService(userRepo repository.UserRepository, jwtUtil *utils.JWTUtil) AuthService {
	return &authService{
		userRepo: userRepo,
		jwtUtil:  jwtUtil,
	}
}

// Login authenticates a user and returns a JWT token
func (s *authService) Login(username, password string) (*LoginResponse, error) {
	// Find user by username
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid username or password")
		}
		return nil, err
	}

	// Verify password
	if err := s.userRepo.ComparePassword(user.Password, password); err != nil {
		return nil, errors.New("invalid username or password")
	}

	// Generate JWT token
	token, err := s.jwtUtil.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

// Register creates a new user account
func (s *authService) Register(username, password, email string) error {
	// Check if username already exists
	existingUser, err := s.userRepo.FindByUsername(username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existingUser != nil {
		return errors.New("username already exists")
	}

	// Create new user
	user := &model.User{
		Username: username,
		Password: password, // Will be hashed by repository
		Email:    email,
		Role:     "admin", // Default role
	}

	return s.userRepo.Create(user)
}

// ValidateToken validates a JWT token and returns the claims
func (s *authService) ValidateToken(token string) (*utils.JWTClaims, error) {
	return s.jwtUtil.ValidateToken(token)
}
