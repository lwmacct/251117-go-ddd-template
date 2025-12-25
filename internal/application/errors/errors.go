// Package errors 重新导出领域错误供 Adapters 层使用
// 遵循 DDD 依赖方向: Adapters → Application → Domain
package errors

import "github.com/lwmacct/251117-go-ddd-template/internal/domain/errors"

// 重新导出领域错误类型
type (
	DomainError       = errors.DomainError
	NotFoundError     = errors.NotFoundError
	ValidationError   = errors.ValidationError
	ConflictError     = errors.ConflictError
	UnauthorizedError = errors.UnauthorizedError
	ForbiddenError    = errors.ForbiddenError
	InternalError     = errors.InternalError
)

// 重新导出错误工厂函数
var (
	NewNotFoundError     = errors.NewNotFoundError
	NewValidationError   = errors.NewValidationError
	NewConflictError     = errors.NewConflictError
	NewUnauthorizedError = errors.NewUnauthorizedError
	NewForbiddenError    = errors.NewForbiddenError
	NewInternalError     = errors.NewInternalError
)

// 重新导出快捷方式
var (
	NotFound      = errors.NotFound
	AlreadyExists = errors.AlreadyExists
	InvalidInput  = errors.InvalidInput
	Unauthorized  = errors.Unauthorized
	Forbidden     = errors.Forbidden
	Internal      = errors.Internal
)

// 重新导出错误判断函数
var (
	Is             = errors.Is
	As             = errors.As
	IsDomainError  = errors.IsDomainError
	GetDomainError = errors.GetDomainError
)
