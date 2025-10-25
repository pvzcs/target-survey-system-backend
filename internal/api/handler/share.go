package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"survey-system/internal/dto/request"
	"survey-system/internal/service"
	"survey-system/pkg/errors"
)

// ShareHandler handles share link related HTTP requests
type ShareHandler struct {
	shareService service.ShareService
}

// NewShareHandler creates a new share handler instance
func NewShareHandler(shareService service.ShareService) *ShareHandler {
	return &ShareHandler{
		shareService: shareService,
	}
}

// GenerateShareLink handles POST /api/v1/surveys/:id/share
func (h *ShareHandler) GenerateShareLink(c *gin.Context) {
	surveyID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "Invalid survey ID",
			},
		})
		return
	}

	var req request.GenerateShareLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error": gin.H{
				"code":    errors.ErrUnauthorized.Code,
				"message": errors.ErrUnauthorized.Message,
			},
		})
		return
	}

	shareLink, err := h.shareService.GenerateShareLink(c.Request.Context(), userID.(uint), uint(surveyID), &req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    shareLink,
	})
}

// GetSurveyByToken handles GET /api/v1/public/surveys/:id (with token query parameter)
func (h *ShareHandler) GetSurveyByToken(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_TOKEN",
				"message": "Token parameter is required",
			},
		})
		return
	}

	survey, err := h.shareService.ValidateAndGetSurvey(c.Request.Context(), token)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    survey,
	})
}
