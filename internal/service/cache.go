package service

import (
	"context"
	"time"

	"survey-system/internal/model"
)

// Cache defines the interface for cache operations used by services
type Cache interface {
	// Survey cache operations
	GetSurvey(ctx context.Context, surveyID uint) (*model.Survey, error)
	SetSurvey(ctx context.Context, survey *model.Survey, expiration time.Duration) error
	DeleteSurvey(ctx context.Context, surveyID uint) error

	// OneLink status cache operations
	GetOneLinkStatus(ctx context.Context, token string) (bool, error)
	SetOneLinkStatus(ctx context.Context, token string, used bool, expiration time.Duration) error

	// Distributed lock operations
	AcquireLock(ctx context.Context, key string, expiration time.Duration) (bool, error)
	ReleaseLock(ctx context.Context, key string) error
}
