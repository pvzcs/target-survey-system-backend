package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"survey-system/internal/dto/request"
	"survey-system/internal/service"
	"survey-system/pkg/errors"
)

// QuestionHandler handles question-related HTTP requests
type QuestionHandler struct {
	questionService service.QuestionService
}

// NewQuestionHandler creates a new question handler instance
func NewQuestionHandler(questionService service.QuestionService) *QuestionHandler {
	return &QuestionHandler{
		questionService: questionService,
	}
}

// CreateQuestion handles POST /api/v1/questions
func (h *QuestionHandler) CreateQuestion(c *gin.Context) {
	var req request.CreateQuestionRequest
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

	question, err := h.questionService.CreateQuestion(c.Request.Context(), userID.(uint), &req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    question,
	})
}

// UpdateQuestion handles PUT /api/v1/questions/:id
func (h *QuestionHandler) UpdateQuestion(c *gin.Context) {
	questionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "Invalid question ID",
			},
		})
		return
	}

	var req request.UpdateQuestionRequest
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

	question, err := h.questionService.UpdateQuestion(c.Request.Context(), userID.(uint), uint(questionID), &req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    question,
	})
}

// DeleteQuestion handles DELETE /api/v1/questions/:id
func (h *QuestionHandler) DeleteQuestion(c *gin.Context) {
	questionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "Invalid question ID",
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

	if err := h.questionService.DeleteQuestion(c.Request.Context(), userID.(uint), uint(questionID)); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Question deleted successfully",
	})
}

// ReorderQuestions handles PUT /api/v1/surveys/:id/questions/reorder
func (h *QuestionHandler) ReorderQuestions(c *gin.Context) {
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

	var req request.ReorderQuestionsRequest
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

	if err := h.questionService.ReorderQuestions(c.Request.Context(), userID.(uint), uint(surveyID), req.QuestionIDs); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Questions reordered successfully",
	})
}
