package handler

import (
	"net/http"
	"survey-system/internal/dto/request"
	"survey-system/internal/dto/response"
	"survey-system/internal/service"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler creates a new auth handler instance
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login handles user login requests
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request.LoginRequest true "Login credentials"
// @Success 200 {object} response.LoginResponse
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "VALIDATION_FAILED",
				"message": "请求参数验证失败",
				"details": err.Error(),
			},
		})
		return
	}

	// Call auth service to login
	loginResp, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		// Check if it's an authentication error
		if err.Error() == "invalid username or password" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "INVALID_CREDENTIALS",
					"message": "用户名或密码错误",
				},
			})
			return
		}

		// Internal server error
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "服务器内部错误",
			},
		})
		return
	}

	// Convert to response DTO
	resp := &response.LoginResponse{
		Token: loginResp.Token,
		User: response.UserResponse{
			ID:        loginResp.User.ID,
			Username:  loginResp.User.Username,
			Email:     loginResp.User.Email,
			Role:      loginResp.User.Role,
			CreatedAt: loginResp.User.CreatedAt,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    resp,
	})
}


