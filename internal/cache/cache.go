package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"survey-system/internal/model"
)

// Cache defines the interface for cache operations
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

	// Health check
	HealthCheck(ctx context.Context) error
}

// RedisCache implements the Cache interface using Redis
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache creates a new Redis cache instance
func NewRedisCache(client *redis.Client) Cache {
	return &RedisCache{
		client: client,
	}
}

// GetSurvey retrieves a survey from cache
func (c *RedisCache) GetSurvey(ctx context.Context, surveyID uint) (*model.Survey, error) {
	key := fmt.Sprintf("survey:%d", surveyID)
	
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, fmt.Errorf("failed to get survey from cache: %w", err)
	}

	var survey model.Survey
	if err := json.Unmarshal(data, &survey); err != nil {
		return nil, fmt.Errorf("failed to unmarshal survey: %w", err)
	}

	return &survey, nil
}

// SetSurvey stores a survey in cache
func (c *RedisCache) SetSurvey(ctx context.Context, survey *model.Survey, expiration time.Duration) error {
	key := fmt.Sprintf("survey:%d", survey.ID)
	
	data, err := json.Marshal(survey)
	if err != nil {
		return fmt.Errorf("failed to marshal survey: %w", err)
	}

	if err := c.client.Set(ctx, key, data, expiration).Err(); err != nil {
		return fmt.Errorf("failed to set survey in cache: %w", err)
	}

	return nil
}

// DeleteSurvey removes a survey from cache
func (c *RedisCache) DeleteSurvey(ctx context.Context, surveyID uint) error {
	key := fmt.Sprintf("survey:%d", surveyID)
	
	if err := c.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete survey from cache: %w", err)
	}

	return nil
}

// GetOneLinkStatus retrieves the used status of a one-time link from cache
func (c *RedisCache) GetOneLinkStatus(ctx context.Context, token string) (bool, error) {
	key := fmt.Sprintf("onelink:status:%s", token)
	
	status, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil // Cache miss, assume not used
		}
		return false, fmt.Errorf("failed to get onelink status from cache: %w", err)
	}

	return status == "used", nil
}

// SetOneLinkStatus stores the used status of a one-time link in cache
func (c *RedisCache) SetOneLinkStatus(ctx context.Context, token string, used bool, expiration time.Duration) error {
	key := fmt.Sprintf("onelink:status:%s", token)
	
	status := "unused"
	if used {
		status = "used"
	}

	if err := c.client.Set(ctx, key, status, expiration).Err(); err != nil {
		return fmt.Errorf("failed to set onelink status in cache: %w", err)
	}

	return nil
}

// AcquireLock attempts to acquire a distributed lock
func (c *RedisCache) AcquireLock(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	lockKey := fmt.Sprintf("lock:%s", key)
	
	// Use SET NX (set if not exists) with expiration
	success, err := c.client.SetNX(ctx, lockKey, "1", expiration).Result()
	if err != nil {
		return false, fmt.Errorf("failed to acquire lock: %w", err)
	}

	return success, nil
}

// ReleaseLock releases a distributed lock
func (c *RedisCache) ReleaseLock(ctx context.Context, key string) error {
	lockKey := fmt.Sprintf("lock:%s", key)
	
	if err := c.client.Del(ctx, lockKey).Err(); err != nil {
		return fmt.Errorf("failed to release lock: %w", err)
	}

	return nil
}

// HealthCheck performs a health check on the Redis connection
func (c *RedisCache) HealthCheck(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}
