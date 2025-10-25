package errors

import "fmt"

// AppError represents an application error with code, message and HTTP status
type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

// NewAppError creates a new AppError
func NewAppError(code, message string, status int) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
	}
}

// Predefined errors
var (
	ErrUnauthorized       = &AppError{"UNAUTHORIZED", "未授权访问", 401}
	ErrForbidden          = &AppError{"FORBIDDEN", "禁止访问", 403}
	ErrNotFound           = &AppError{"NOT_FOUND", "资源不存在", 404}
	ErrInvalidToken       = &AppError{"INVALID_TOKEN", "无效的令牌", 400}
	ErrTokenExpired       = &AppError{"TOKEN_EXPIRED", "令牌已过期", 403}
	ErrLinkUsed           = &AppError{"LINK_USED", "链接已被使用", 403}
	ErrValidationFailed   = &AppError{"VALIDATION_FAILED", "数据验证失败", 400}
	ErrSurveyNotPublished = &AppError{"SURVEY_NOT_PUBLISHED", "问卷未发布", 400}
	ErrInternalServer     = &AppError{"INTERNAL_ERROR", "服务器内部错误", 500}
	ErrBadRequest         = &AppError{"BAD_REQUEST", "请求参数错误", 400}
)

// WrapError wraps an error with additional context
func WrapError(err error, message string) error {
	return fmt.Errorf("%s: %w", message, err)
}

// NewValidationError creates a validation error with field and reason
func NewValidationError(field, reason string) *AppError {
	return &AppError{
		Code:    "VALIDATION_FAILED",
		Message: fmt.Sprintf("validation failed for field '%s': %s", field, reason),
		Status:  400,
	}
}
