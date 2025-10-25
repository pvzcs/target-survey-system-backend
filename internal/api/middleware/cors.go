package middleware

import (
	"survey-system/internal/config"

	"github.com/gin-gonic/gin"
)

// CORS returns a middleware that handles CORS
func CORS(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// Check if origin is allowed
		allowed := false
		for _, allowedOrigin := range cfg.CORS.AllowedOrigins {
			if allowedOrigin == "*" || allowedOrigin == origin {
				allowed = true
				break
			}
		}

		if allowed {
			if origin != "" {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			} else if len(cfg.CORS.AllowedOrigins) > 0 {
				c.Writer.Header().Set("Access-Control-Allow-Origin", cfg.CORS.AllowedOrigins[0])
			}
			
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Methods", joinStrings(cfg.CORS.AllowedMethods, ", "))
			c.Writer.Header().Set("Access-Control-Allow-Headers", joinStrings(cfg.CORS.AllowedHeaders, ", "))
			c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}
