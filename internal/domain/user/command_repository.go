// Package user 定义用户命令仓储接口（写操作）
package user

import (
	"context"
)

// CommandRepository 用户命令仓储接口
// 负责所有修改状态的操作（Create, Update, Delete）
type CommandRepository interface {
	// Create 创建用户
	Create(ctx context.Context, user *User) error

	// Update 更新用户
	Update(ctx context.Context, user *User) error

	// Delete 删除用户 (软删除)
	Delete(ctx context.Context, id uint) error

	// AssignRoles 为用户分配角色
	AssignRoles(ctx context.Context, userID uint, roleIDs []uint) error

	// RemoveRoles 移除用户的角色
	RemoveRoles(ctx context.Context, userID uint, roleIDs []uint) error

	// UpdatePassword 更新用户密码
	UpdatePassword(ctx context.Context, userID uint, hashedPassword string) error

	// UpdateStatus 更新用户状态
	UpdateStatus(ctx context.Context, userID uint, status string) error
}
