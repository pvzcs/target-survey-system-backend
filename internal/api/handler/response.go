package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"survey-system/internal/dto/request"
	"survey-system/internal/service"
	"survey-system/pkg/errors"

	"github.com/gin-gonic/gin"
)

// ResponseHandler handles response-related HTTP requests
type ResponseHandler struct {
	responseSvc *service.ResponseService
}

// NewResponseHandler creates a new ResponseHandler
func NewResponseHandler(responseSvc *service.ResponseService) *ResponseHandler {
	return &ResponseHandler{
		responseSvc: responseSvc,
	}
}

// SubmitResponse handles POST /api/v1/public/responses
func (h *ResponseHandler) SubmitResponse(c *gin.Context) {
	var req request.SubmitResponseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "BAD_REQUEST",
				"message": "请求参数错误: " + err.Error(),
			},
		})
		return
	}

	// Get IP address
	ipAddress := c.ClientIP()

	// Get User-Agent
	userAgent := c.GetHeader("User-Agent")

	// Submit response
	resp, err := h.responseSvc.SubmitResponse(&req, ipAddress, userAgent)
	if err != nil {
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

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "服务器内部错误",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    resp,
	})
}

// GetResponses handles GET /api/v1/surveys/:id/responses
func (h *ResponseHandler) GetResponses(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "UNAUTHORIZED",
				"message": "未授权访问",
			},
		})
		return
	}

	// Get survey ID from URL parameter
	surveyID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "无效的问卷 ID",
			},
		})
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// Get responses
	responseList, meta, err := h.responseSvc.GetResponses(userID.(uint), uint(surveyID), page, pageSize)
	if err != nil {
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

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "服务器内部错误",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    responseList,
		"meta":    meta,
	})
}

// GetStatistics handles GET /api/v1/surveys/:id/statistics
func (h *ResponseHandler) GetStatistics(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "UNAUTHORIZED",
				"message": "未授权访问",
			},
		})
		return
	}

	// Get survey ID from URL parameter
	surveyID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "无效的问卷 ID",
			},
		})
		return
	}

	// Get statistics
	resp, err := h.responseSvc.GetStatistics(userID.(uint), uint(surveyID))
	if err != nil {
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

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "服务器内部错误",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    resp,
	})
}

// ExportResponses handles GET /api/v1/surveys/:id/export
func (h *ResponseHandler) ExportResponses(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "UNAUTHORIZED",
				"message": "未授权访问",
			},
		})
		return
	}

	// Get survey ID from URL parameter
	surveyID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "无效的问卷 ID",
			},
		})
		return
	}

	// Get format parameter (default to csv)
	format := c.DefaultQuery("format", "csv")
	if format != "csv" && format != "excel" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_FORMAT",
				"message": "不支持的导出格式，请使用 csv 或 excel",
			},
		})
		return
	}

	// Export responses
	data, filename, err := h.responseSvc.ExportResponses(userID.(uint), uint(surveyID), format)
	if err != nil {
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

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "服务器内部错误",
			},
		})
		return
	}

	// Set appropriate headers based on format
	var contentType string
	if format == "csv" {
		contentType = "text/csv; charset=utf-8"
	} else {
		contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	}

	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Header("Content-Length", strconv.Itoa(len(data)))

	c.Data(http.StatusOK, contentType, data)
}
