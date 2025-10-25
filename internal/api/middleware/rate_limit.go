package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RateLimit returns a middleware that limits requests per IP address
func RateLimit(redisClient *redis.Client, requestsPerMinute int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get client IP
		clientIP := c.ClientIP()
		
		// Create Redis key
		key := fmt.Sprintf("ratelimit:ip:%s", clientIP)
		
		ctx := context.Background()
		
		// Increment counter
		count, err := redisClient.Incr(ctx, key).Result()
		if err != nil {
			// If Redis fails, allow the request but log the error
			c.Next()
			return
		}
		
		// Set expiration on first request
		if count == 1 {
			redisClient.Expire(ctx, key, time.Minute)
		}
		
		// Check if limit exceeded
		if count > int64(requestsPerMinute) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "RATE_LIMIT_EXCEEDED",
					"message": "请求过于频繁，请稍后再试",
				},
			})
			c.Abort()
			return
		}
		
		// Add rate limit headers
		c.Writer.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", requestsPerMinute))
		c.Writer.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", requestsPerMinute-int(count)))
		
		c.Next()
	}
}
