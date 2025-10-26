package repository

import (
	"survey-system/internal/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(user *model.User) error
	FindByID(id uint) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
	Update(user *model.User) error
	UpdatePassword(userID uint, newPassword string) error
	HashPassword(password string) (string, error)
	ComparePassword(hashedPassword, password string) error
}

// userRepository implements UserRepository interface
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create creates a new user with hashed password
func (r *userRepository) Create(user *model.User) error {
	// Hash the password before storing
	hashedPassword, err := r.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	return r.db.Create(user).Error
}

// FindByID finds a user by ID
func (r *userRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsername finds a user by username
func (r *userRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// HashPassword hashes a plain text password using bcrypt
func (r *userRepository) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// ComparePassword compares a hashed password with a plain text password
func (r *userRepository) ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// Update updates user information (excluding password)
func (r *userRepository) Update(user *model.User) error {
	return r.db.Model(user).Updates(map[string]interface{}{
		"username": user.Username,
		"email":    user.Email,
	}).Error
}

// UpdatePassword updates user password with hashing
func (r *userRepository) UpdatePassword(userID uint, newPassword string) error {
	hashedPassword, err := r.HashPassword(newPassword)
	if err != nil {
		return err
	}

	return r.db.Model(&model.User{}).Where("id = ?", userID).Update("password", hashedPassword).Error
}
