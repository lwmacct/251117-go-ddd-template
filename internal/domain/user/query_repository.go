// Package user 定义用户查询仓储接口（读操作）
package user

import (
	"context"
)

// QueryRepository 用户查询仓储接口
// 负责所有查询操作（Get, List, Count）
type QueryRepository interface {
	// GetByID 根据 ID 获取用户
	GetByID(ctx context.Context, id uint) (*User, error)

	// GetByUsername 根据用户名获取用户
	GetByUsername(ctx context.Context, username string) (*User, error)

	// GetByEmail 根据邮箱获取用户
	GetByEmail(ctx context.Context, email string) (*User, error)

	// GetByIDWithRoles 根据 ID 获取用户（包含角色和权限信息）
	GetByIDWithRoles(ctx context.Context, id uint) (*User, error)

	// GetByUsernameWithRoles 根据用户名获取用户（包含角色和权限信息）
	GetByUsernameWithRoles(ctx context.Context, username string) (*User, error)

	// GetByEmailWithRoles 根据邮箱获取用户（包含角色和权限信息）
	GetByEmailWithRoles(ctx context.Context, email string) (*User, error)

	// List 获取用户列表 (分页)
	List(ctx context.Context, offset, limit int) ([]*User, error)

	// Count 统计用户数量
	Count(ctx context.Context) (int64, error)

	// GetRoles 获取用户的所有角色 ID
	GetRoles(ctx context.Context, userID uint) ([]uint, error)

	// Search 搜索用户（支持用户名、邮箱、全名模糊匹配）
	Search(ctx context.Context, keyword string, offset, limit int) ([]*User, error)

	// Exists 检查用户是否存在
	Exists(ctx context.Context, id uint) (bool, error)

	// ExistsByUsername 检查用户名是否存在
	ExistsByUsername(ctx context.Context, username string) (bool, error)

	// ExistsByEmail 检查邮箱是否存在
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}
