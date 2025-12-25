package errors

import (
	"errors"
	"fmt"
)

// ============================================================================
// 错误包装
// ============================================================================

// Wrap 包装错误并添加上下文信息
// 保留原始错误的类型信息
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}

	// 如果是领域错误，保留其类型
	if domainErr := GetDomainError(err); domainErr != nil {
		return &wrappedDomainError{
			DomainError: domainErr,
			message:     message,
			cause:       err,
		}
	}

	return fmt.Errorf("%s: %w", message, err)
}

// Wrapf 包装错误并添加格式化的上下文信息
func Wrapf(err error, format string, args ...any) error {
	return Wrap(err, fmt.Sprintf(format, args...))
}

// wrappedDomainError 包装后的领域错误
// 保留原始领域错误的接口信息
type wrappedDomainError struct {
	DomainError

	message string
	cause   error
}

func (e *wrappedDomainError) Error() string {
	return fmt.Sprintf("%s: %v", e.message, e.cause)
}

func (e *wrappedDomainError) Unwrap() error {
	return e.cause
}

// ============================================================================
// 错误链工具
// ============================================================================

// WithCause 为错误添加原因
func WithCause(err, cause error) error {
	if err == nil {
		return nil
	}

	var notFound *NotFoundError
	var validation *ValidationError
	var conflict *ConflictError
	var unauthorized *UnauthorizedError
	var forbidden *ForbiddenError
	var internal *InternalError

	switch {
	case errors.As(err, &notFound):
		notFound.cause = cause
		return notFound
	case errors.As(err, &validation):
		validation.cause = cause
		return validation
	case errors.As(err, &conflict):
		conflict.cause = cause
		return conflict
	case errors.As(err, &unauthorized):
		unauthorized.cause = cause
		return unauthorized
	case errors.As(err, &forbidden):
		forbidden.cause = cause
		return forbidden
	case errors.As(err, &internal):
		internal.cause = cause
		return internal
	default:
		return fmt.Errorf("%w: %w", err, cause)
	}
}

// RootCause 获取错误链的根因
func RootCause(err error) error {
	for {
		unwrapper, ok := err.(interface{ Unwrap() error })
		if !ok {
			return err
		}
		unwrapped := unwrapper.Unwrap()
		if unwrapped == nil {
			return err
		}
		err = unwrapped
	}
}
