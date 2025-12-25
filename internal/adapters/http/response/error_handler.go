package response

import (
	"errors"
	"net/http"

	appErrors "github.com/lwmacct/251117-go-ddd-template/internal/application/errors"

	"github.com/gin-gonic/gin"
)

// HandleError 统一处理领域错误并返回适当的 HTTP 响应
//
// 使用示例：
//
//	if err := h.handler.Handle(ctx, cmd); err != nil {
//	    response.HandleError(c, err)
//	    return
//	}
func HandleError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	// 尝试获取领域错误
	var domainErr appErrors.DomainError
	if errors.As(err, &domainErr) {
		handleDomainError(c, domainErr)
		return
	}

	// 未知错误视为内部错误
	InternalError(c)
}

// handleDomainError 处理领域错误
func handleDomainError(c *gin.Context, err appErrors.DomainError) {
	statusCode := err.HTTPStatus()
	errorCode := err.Code()
	message := err.Error()

	// 构建错误详情
	errorDetail := ErrorDetail{
		Code:    errorCode,
		Message: message,
	}

	// 如果是验证错误，尝试获取详情
	var validationErr *appErrors.ValidationError
	if errors.As(err, &validationErr) && validationErr.Details != nil {
		errorDetail.Details = validationErr.Details
	}

	c.JSON(statusCode, UnifiedResponse{
		Code:    statusCode,
		Message: message,
		Error:   errorDetail,
	})
}

// HandleErrorWithMessage 处理错误并使用自定义消息
func HandleErrorWithMessage(c *gin.Context, err error, message string) {
	if err == nil {
		return
	}

	var domainErr appErrors.DomainError
	if errors.As(err, &domainErr) {
		statusCode := domainErr.HTTPStatus()
		c.JSON(statusCode, UnifiedResponse{
			Code:    statusCode,
			Message: message,
			Error: ErrorDetail{
				Code:    domainErr.Code(),
				Message: domainErr.Error(),
			},
		})
		return
	}

	InternalError(c)
}

// ============================================================================
// 错误类型判断辅助函数
// ============================================================================

// IsNotFoundError 检查是否为资源不存在错误
func IsNotFoundError(err error) bool {
	var notFoundErr *appErrors.NotFoundError
	return errors.As(err, &notFoundErr)
}

// IsConflictError 检查是否为资源冲突错误
func IsConflictError(err error) bool {
	var conflictErr *appErrors.ConflictError
	return errors.As(err, &conflictErr)
}

// IsValidationError 检查是否为验证错误
func IsValidationError(err error) bool {
	var validationErr *appErrors.ValidationError
	return errors.As(err, &validationErr)
}

// IsUnauthorizedError 检查是否为未认证错误
func IsUnauthorizedError(err error) bool {
	var unauthorizedErr *appErrors.UnauthorizedError
	return errors.As(err, &unauthorizedErr)
}

// IsForbiddenError 检查是否为无权限错误
func IsForbiddenError(err error) bool {
	var forbiddenErr *appErrors.ForbiddenError
	return errors.As(err, &forbiddenErr)
}

// ============================================================================
// 快捷错误响应函数
// ============================================================================

// DomainNotFound 返回领域资源不存在错误
func DomainNotFound(c *gin.Context, err error) {
	if err == nil {
		NotFound(c, "")
		return
	}
	HandleError(c, err)
}

// DomainConflict 返回领域资源冲突错误
func DomainConflict(c *gin.Context, err error) {
	if err == nil {
		Conflict(c, "")
		return
	}
	HandleError(c, err)
}

// DomainValidationError 返回领域验证错误
func DomainValidationError(c *gin.Context, err error) {
	if err == nil {
		BadRequest(c, "validation failed")
		return
	}
	HandleError(c, err)
}

// ============================================================================
// 错误映射（用于向后兼容）
// ============================================================================

// MapErrorToStatus 将错误映射到 HTTP 状态码
// 用于需要手动处理错误状态码的场景
func MapErrorToStatus(err error) int {
	if err == nil {
		return http.StatusOK
	}

	var domainErr appErrors.DomainError
	if errors.As(err, &domainErr) {
		return domainErr.HTTPStatus()
	}

	return http.StatusInternalServerError
}
