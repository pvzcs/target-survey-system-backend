package repository

import (
	"survey-system/internal/model"
	"time"

	"gorm.io/gorm"
)

// OneLinkRepository defines the interface for one-time link data operations
type OneLinkRepository interface {
	Create(oneLink *model.OneLink) error
	FindByToken(token string) (*model.OneLink, error)
	MarkAsUsed(id uint) error
	MarkAsAccessed(id uint) error
	DeleteExpired() error
}

// oneLinkRepository implements OneLinkRepository interface
type oneLinkRepository struct {
	db *gorm.DB
}

// NewOneLinkRepository creates a new one-time link repository instance
func NewOneLinkRepository(db *gorm.DB) OneLinkRepository {
	return &oneLinkRepository{db: db}
}

// Create creates a new one-time link record
func (r *oneLinkRepository) Create(oneLink *model.OneLink) error {
	return r.db.Create(oneLink).Error
}

// FindByToken finds a one-time link by its token
func (r *oneLinkRepository) FindByToken(token string) (*model.OneLink, error) {
	var oneLink model.OneLink
	err := r.db.Where("token = ?", token).First(&oneLink).Error
	if err != nil {
		return nil, err
	}
	return &oneLink, nil
}

// MarkAsUsed marks a one-time link as used
func (r *oneLinkRepository) MarkAsUsed(id uint) error {
	now := time.Now()
	return r.db.Model(&model.OneLink{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"used":    true,
			"used_at": now,
		}).Error
}

// MarkAsAccessed marks a one-time link as accessed (first time viewing)
func (r *oneLinkRepository) MarkAsAccessed(id uint) error {
	now := time.Now()
	return r.db.Model(&model.OneLink{}).
		Where("id = ? AND accessed_at IS NULL", id).
		Update("accessed_at", now).Error
}

// DeleteExpired deletes all expired one-time links
func (r *oneLinkRepository) DeleteExpired() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&model.OneLink{}).Error
}
