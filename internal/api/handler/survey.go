package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"survey-system/internal/dto/request"
	"survey-system/internal/service"
	"survey-system/pkg/errors"
)

// SurveyHandler handles survey-related HTTP requests
type SurveyHandler struct {
	surveyService service.SurveyService
}

// NewSurveyHandler creates a new survey handler instance
func NewSurveyHandler(surveyService service.SurveyService) *SurveyHandler {
	return &SurveyHandler{
		surveyService: surveyService,
	}
}

// CreateSurvey handles POST /api/v1/surveys
func (h *SurveyHandler) CreateSurvey(c *gin.Context) {
	var req request.CreateSurveyRequest
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

	survey, err := h.surveyService.CreateSurvey(c.Request.Context(), userID.(uint), &req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    survey,
	})
}

// UpdateSurvey handles PUT /api/v1/surveys/:id
func (h *SurveyHandler) UpdateSurvey(c *gin.Context) {
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

	var req request.UpdateSurveyRequest
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

	survey, err := h.surveyService.UpdateSurvey(c.Request.Context(), userID.(uint), uint(surveyID), &req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    survey,
	})
}

// DeleteSurvey handles DELETE /api/v1/surveys/:id
func (h *SurveyHandler) DeleteSurvey(c *gin.Context) {
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

	if err := h.surveyService.DeleteSurvey(c.Request.Context(), userID.(uint), uint(surveyID)); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Survey deleted successfully",
	})
}

// GetSurvey handles GET /api/v1/surveys/:id
func (h *SurveyHandler) GetSurvey(c *gin.Context) {
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

	survey, err := h.surveyService.GetSurvey(c.Request.Context(), uint(surveyID))
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    survey,
	})
}

// ListSurveys handles GET /api/v1/surveys
func (h *SurveyHandler) ListSurveys(c *gin.Context) {
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

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	surveys, err := h.surveyService.ListSurveys(c.Request.Context(), userID.(uint), page, pageSize)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    surveys.Data,
		"meta":    surveys.Meta,
	})
}

// PublishSurvey handles POST /api/v1/surveys/:id/publish
func (h *SurveyHandler) PublishSurvey(c *gin.Context) {
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

	if err := h.surveyService.PublishSurvey(c.Request.Context(), userID.(uint), uint(surveyID)); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Survey published successfully",
	})
}

// handleError handles errors and returns appropriate HTTP responses
func handleError(c *gin.Context, err error) {
	if appErr, ok := err.(*errors.AppError); ok {
		c.JSON(appErr.Status, gin.H{
			"success": false,
			"error": gin.H{
				"code":    appErr.Code,
				"message": appErr.Message,
			},
		})
		return
	}

	// Default to internal server error
	c.JSON(http.StatusInternalServerError, gin.H{
		"success": false,
		"error": gin.H{
			"code":    errors.ErrInternalServer.Code,
			"message": err.Error(),
		},
	})
}
