// Package user 定义用户领域的错误
package user

import "errors"

// 用户相关错误
var (
	// ErrUserNotFound 用户不存在
	ErrUserNotFound = errors.New("user not found")

	// ErrUserAlreadyExists 用户已存在
	ErrUserAlreadyExists = errors.New("user already exists")

	// ErrUsernameAlreadyExists 用户名已存在
	ErrUsernameAlreadyExists = errors.New("username already exists")

	// ErrEmailAlreadyExists 邮箱已存在
	ErrEmailAlreadyExists = errors.New("email already exists")

	// ErrRoleAlreadyAssigned 角色已分配
	ErrRoleAlreadyAssigned = errors.New("role already assigned to user")

	// ErrRoleNotFound 角色不存在
	ErrRoleNotFound = errors.New("role not found")

	// ErrInvalidUserStatus 无效的用户状态
	ErrInvalidUserStatus = errors.New("invalid user status")

	// ErrCannotDeleteSelf 不能删除自己
	ErrCannotDeleteSelf = errors.New("cannot delete yourself")

	// ErrCannotModifyAdmin 不能修改管理员
	ErrCannotModifyAdmin = errors.New("cannot modify admin user")

	// ErrInvalidPassword 密码错误
	ErrInvalidPassword = errors.New("invalid password")
)
