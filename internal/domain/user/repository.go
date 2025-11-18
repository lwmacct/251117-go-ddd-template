// Package user 定义用户仓储接口
package user

import (
	"context"
)

// Repository 用户仓储接口
type Repository interface {
	// Create 创建用户
	Create(ctx context.Context, user *User) error

	// GetByID 根据 ID 获取用户
	GetByID(ctx context.Context, id uint) (*User, error)

	// GetByUsername 根据用户名获取用户
	GetByUsername(ctx context.Context, username string) (*User, error)

	// GetByEmail 根据邮箱获取用户
	GetByEmail(ctx context.Context, email string) (*User, error)

	// List 获取用户列表 (分页)
	List(ctx context.Context, offset, limit int) ([]*User, error)

	// Update 更新用户
	Update(ctx context.Context, user *User) error

	// Delete 删除用户 (软删除)
	Delete(ctx context.Context, id uint) error

	// Count 统计用户数量
	Count(ctx context.Context) (int64, error)
}
