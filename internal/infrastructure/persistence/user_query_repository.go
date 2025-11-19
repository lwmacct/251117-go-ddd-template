// Package persistence 提供用户查询仓储的 GORM 实现
package persistence

import (
	"context"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
	"gorm.io/gorm"
)

// userQueryRepository 用户查询仓储的 GORM 实现
type userQueryRepository struct {
	db *gorm.DB
}

// NewUserQueryRepository 创建用户查询仓储实例
func NewUserQueryRepository(db *gorm.DB) user.QueryRepository {
	return &userQueryRepository{db: db}
}

// GetByID 根据 ID 获取用户
func (r *userQueryRepository) GetByID(ctx context.Context, id uint) (*user.User, error) {
	var u user.User
	if err := r.db.WithContext(ctx).First(&u, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, user.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return &u, nil
}

// GetByUsername 根据用户名获取用户
func (r *userQueryRepository) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	var u user.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, user.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}
	return &u, nil
}

// GetByEmail 根据邮箱获取用户
func (r *userQueryRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	var u user.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, user.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &u, nil
}

// GetByIDWithRoles 根据 ID 获取用户（包含角色和权限信息）
func (r *userQueryRepository) GetByIDWithRoles(ctx context.Context, id uint) (*user.User, error) {
	var u user.User
	if err := r.db.WithContext(ctx).
		Preload("Roles.Permissions").
		First(&u, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, user.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by id with roles: %w", err)
	}
	return &u, nil
}

// GetByUsernameWithRoles 根据用户名获取用户（包含角色和权限信息）
func (r *userQueryRepository) GetByUsernameWithRoles(ctx context.Context, username string) (*user.User, error) {
	var u user.User
	if err := r.db.WithContext(ctx).
		Preload("Roles.Permissions").
		Where("username = ?", username).
		First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, user.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by username with roles: %w", err)
	}
	return &u, nil
}

// GetByEmailWithRoles 根据邮箱获取用户（包含角色和权限信息）
func (r *userQueryRepository) GetByEmailWithRoles(ctx context.Context, email string) (*user.User, error) {
	var u user.User
	if err := r.db.WithContext(ctx).
		Preload("Roles.Permissions").
		Where("email = ?", email).
		First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, user.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by email with roles: %w", err)
	}
	return &u, nil
}

// List 获取用户列表 (分页)
func (r *userQueryRepository) List(ctx context.Context, offset, limit int) ([]*user.User, error) {
	var users []*user.User
	query := r.db.WithContext(ctx).Offset(offset).Limit(limit)

	if err := query.Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return users, nil
}

// Count 统计用户数量
func (r *userQueryRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&user.User{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}
	return count, nil
}

// GetRoles 获取用户的所有角色 ID
func (r *userQueryRepository) GetRoles(ctx context.Context, userID uint) ([]uint, error) {
	var u user.User
	if err := r.db.WithContext(ctx).Preload("Roles").First(&u, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, user.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	roleIDs := make([]uint, 0, len(u.Roles))
	for _, r := range u.Roles {
		roleIDs = append(roleIDs, r.ID)
	}

	return roleIDs, nil
}

// Search 搜索用户（支持用户名、邮箱、全名模糊匹配）
func (r *userQueryRepository) Search(ctx context.Context, keyword string, offset, limit int) ([]*user.User, error) {
	var users []*user.User
	query := r.db.WithContext(ctx).
		Where("username LIKE ? OR email LIKE ? OR full_name LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%").
		Offset(offset).
		Limit(limit)

	if err := query.Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}

	return users, nil
}

// Exists 检查用户是否存在
func (r *userQueryRepository) Exists(ctx context.Context, id uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&user.User{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}
	return count > 0, nil
}

// ExistsByUsername 检查用户名是否存在
func (r *userQueryRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&user.User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check username existence: %w", err)
	}
	return count > 0, nil
}

// ExistsByEmail 检查邮箱是否存在
func (r *userQueryRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&user.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}
	return count > 0, nil
}
