package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// DomainError 领域错误接口
// 所有领域错误都应实现此接口，便于 HTTP Handler 统一处理
type DomainError interface {
	error
	Code() string      // 错误码，如 "user.not_found"
	HTTPStatus() int   // 建议的 HTTP 状态码
	IsRetryable() bool // 是否可重试
}

// ============================================================================
// 基础错误结构
// ============================================================================

// baseError 基础错误结构
type baseError struct {
	code       string
	message    string
	httpStatus int
	retryable  bool
	cause      error
}

func (e *baseError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %v", e.message, e.cause)
	}
	return e.message
}

func (e *baseError) Code() string      { return e.code }
func (e *baseError) HTTPStatus() int   { return e.httpStatus }
func (e *baseError) IsRetryable() bool { return e.retryable }
func (e *baseError) Unwrap() error     { return e.cause }

// ============================================================================
// 具体错误类型
// ============================================================================

// NotFoundError 404 资源不存在错误
type NotFoundError struct{ baseError }

// NewNotFoundError 创建资源不存在错误
func NewNotFoundError(code, message string) *NotFoundError {
	return &NotFoundError{
		baseError: baseError{
			code:       code,
			message:    message,
			httpStatus: http.StatusNotFound,
			retryable:  false,
		},
	}
}

// ValidationError 400 验证错误
type ValidationError struct {
	baseError

	Details any // 验证详情（如字段错误列表）
}

// NewValidationError 创建验证错误
func NewValidationError(code, message string) *ValidationError {
	return &ValidationError{
		baseError: baseError{
			code:       code,
			message:    message,
			httpStatus: http.StatusBadRequest,
			retryable:  false,
		},
	}
}

// WithDetails 添加验证详情
func (e *ValidationError) WithDetails(details any) *ValidationError {
	e.Details = details
	return e
}

// ConflictError 409 资源冲突错误
type ConflictError struct{ baseError }

// NewConflictError 创建资源冲突错误
func NewConflictError(code, message string) *ConflictError {
	return &ConflictError{
		baseError: baseError{
			code:       code,
			message:    message,
			httpStatus: http.StatusConflict,
			retryable:  false,
		},
	}
}

// UnauthorizedError 401 未认证错误
type UnauthorizedError struct{ baseError }

// NewUnauthorizedError 创建未认证错误
func NewUnauthorizedError(code, message string) *UnauthorizedError {
	return &UnauthorizedError{
		baseError: baseError{
			code:       code,
			message:    message,
			httpStatus: http.StatusUnauthorized,
			retryable:  false,
		},
	}
}

// ForbiddenError 403 无权限错误
type ForbiddenError struct{ baseError }

// NewForbiddenError 创建无权限错误
func NewForbiddenError(code, message string) *ForbiddenError {
	return &ForbiddenError{
		baseError: baseError{
			code:       code,
			message:    message,
			httpStatus: http.StatusForbidden,
			retryable:  false,
		},
	}
}

// InternalError 500 内部错误
type InternalError struct{ baseError }

// NewInternalError 创建内部错误
func NewInternalError(code, message string) *InternalError {
	return &InternalError{
		baseError: baseError{
			code:       code,
			message:    message,
			httpStatus: http.StatusInternalServerError,
			retryable:  true, // 内部错误通常可重试
		},
	}
}

// ============================================================================
// 错误判断函数
// ============================================================================

// Is 检查错误是否为指定类型
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As 将错误转换为指定类型
func As(err error, target any) bool {
	return errors.As(err, target)
}

// IsDomainError 检查是否为领域错误
func IsDomainError(err error) bool {
	var domainErr DomainError
	return errors.As(err, &domainErr)
}

// GetDomainError 获取领域错误，如果不是则返回 nil
func GetDomainError(err error) DomainError {
	var domainErr DomainError
	if errors.As(err, &domainErr) {
		return domainErr
	}
	return nil
}

// ============================================================================
// 错误工厂函数（快捷方式）
// ============================================================================

// NotFound 创建资源不存在错误（快捷方式）
func NotFound(resource string) *NotFoundError {
	return NewNotFoundError(
		resource+".not_found",
		resource+" 不存在",
	)
}

// AlreadyExists 创建资源已存在错误（快捷方式）
func AlreadyExists(resource string) *ConflictError {
	return NewConflictError(
		resource+".already_exists",
		resource+" 已存在",
	)
}

// InvalidInput 创建无效输入错误（快捷方式）
func InvalidInput(message string) *ValidationError {
	return NewValidationError("validation.failed", message)
}

// Unauthorized 创建未认证错误（快捷方式）
func Unauthorized(message string) *UnauthorizedError {
	if message == "" {
		message = "需要认证"
	}
	return NewUnauthorizedError("auth.unauthorized", message)
}

// Forbidden 创建无权限错误（快捷方式）
func Forbidden(message string) *ForbiddenError {
	if message == "" {
		message = "访问被禁止"
	}
	return NewForbiddenError("auth.forbidden", message)
}

// Internal 创建内部错误（快捷方式）
func Internal(message string) *InternalError {
	if message == "" {
		message = "内部服务器错误"
	}
	return NewInternalError("internal.error", message)
}
