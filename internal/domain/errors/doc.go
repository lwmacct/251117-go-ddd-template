// Package errors 提供统一的领域错误处理体系。
//
// 本包定义了分层错误类型，使 HTTP Handler 能自动映射错误到正确的状态码。
//
// # 错误接口
//
// [DomainError] 是所有领域错误的基础接口：
//
//	type DomainError interface {
//	    error
//	    Code() string        // 错误码，如 "user.not_found"
//	    HTTPStatus() int     // 建议的 HTTP 状态码
//	    IsRetryable() bool   // 是否可重试
//	}
//
// # 错误类型
//
// 根据 HTTP 状态码分类的错误类型：
//   - [NotFoundError]: 404 资源不存在
//   - [ValidationError]: 400 验证失败
//   - [ConflictError]: 409 资源冲突
//   - [UnauthorizedError]: 401 未认证
//   - [ForbiddenError]: 403 无权限
//   - [InternalError]: 500 内部错误
//
// # 使用示例
//
// 定义领域错误：
//
//	var ErrUserNotFound = errors.NewNotFoundError("user.not_found", "user not found")
//	var ErrEmailExists = errors.NewConflictError("user.email_exists", "email already exists")
//
// HTTP Handler 错误处理：
//
//	if err := h.handler.Handle(ctx, cmd); err != nil {
//	    errors.HandleError(c, err)
//	    return
//	}
//
// # 错误包装
//
// 使用 [Wrap] 和 [Wrapf] 添加上下文：
//
//	err := errors.Wrap(ErrUserNotFound, "failed to get user by ID")
//	err := errors.Wrapf(ErrUserNotFound, "user %d not found", userID)
//
// # 错误码规范
//
// 错误码格式：{domain}.{error_type}
//   - user.not_found
//   - user.email_exists
//   - auth.invalid_password
//   - role.system_protected
package errors
