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

// UpdateProfile handles user profile update requests
// @Summary Update user profile
// @Description Update username, email, and/or password
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.UpdateProfileRequest true "Profile update data"
// @Success 200 {object} response.UpdateProfileResponse
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Router /api/v1/auth/profile [put]
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "UNAUTHORIZED",
				"message": "用户未认证",
			},
		})
		return
	}

	var req request.UpdateProfileRequest
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

	// Validate that at least one field is being updated
	if req.Username == "" && req.Email == "" && req.NewPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "VALIDATION_FAILED",
				"message": "至少需要提供一个要更新的字段",
			},
		})
		return
	}

	// If password is being changed, old password is required
	if req.NewPassword != "" && req.OldPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "VALIDATION_FAILED",
				"message": "修改密码需要提供旧密码",
			},
		})
		return
	}

	// Call auth service to update profile
	updatedUser, err := h.authService.UpdateProfile(
		userID.(uint),
		req.Username,
		req.Email,
		req.OldPassword,
		req.NewPassword,
	)
	if err != nil {
		// Check specific error types
		switch err.Error() {
		case "user not found":
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "USER_NOT_FOUND",
					"message": "用户不存在",
				},
			})
			return
		case "username already exists":
			c.JSON(http.StatusConflict, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "USERNAME_EXISTS",
					"message": "用户名已存在",
				},
			})
			return
		case "old password is incorrect":
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "INVALID_PASSWORD",
					"message": "旧密码不正确",
				},
			})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "INTERNAL_ERROR",
					"message": "服务器内部错误",
				},
			})
			return
		}
	}

	// Convert to response DTO
	resp := &response.UpdateProfileResponse{
		Message: "个人信息更新成功",
		User: response.UserResponse{
			ID:        updatedUser.ID,
			Username:  updatedUser.Username,
			Email:     updatedUser.Email,
			Role:      updatedUser.Role,
			CreatedAt: updatedUser.CreatedAt,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    resp,
	})
}
