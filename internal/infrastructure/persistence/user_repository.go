// Package persistence 提供用户仓储的 GORM 实现
package persistence

import (
	"context"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/role"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
	"gorm.io/gorm"
)

// userRepository 用户仓储的 GORM 实现
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储实例
func NewUserRepository(db *gorm.DB) user.Repository {
	return &userRepository{db: db}
}

// Create 创建用户
func (r *userRepository) Create(ctx context.Context, u *user.User) error {
	if err := r.db.WithContext(ctx).Create(u).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// GetByID 根据 ID 获取用户
func (r *userRepository) GetByID(ctx context.Context, id uint) (*user.User, error) {
	var u user.User
	if err := r.db.WithContext(ctx).First(&u, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found with id: %d", id)
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return &u, nil
}

// GetByUsername 根据用户名获取用户
func (r *userRepository) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	var u user.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found with username: %s", username)
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}
	return &u, nil
}

// GetByEmail 根据邮箱获取用户
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	var u user.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found with email: %s", email)
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &u, nil
}

// List 获取用户列表 (分页)
func (r *userRepository) List(ctx context.Context, offset, limit int) ([]*user.User, error) {
	var users []*user.User
	query := r.db.WithContext(ctx).Offset(offset).Limit(limit)

	if err := query.Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return users, nil
}

// Update 更新用户
func (r *userRepository) Update(ctx context.Context, u *user.User) error {
	if err := r.db.WithContext(ctx).Save(u).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// Delete 删除用户 (软删除)
func (r *userRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&user.User{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// Count 统计用户数量
func (r *userRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&user.User{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}
	return count, nil
}

// GetByIDWithRoles 根据 ID 获取用户（包含角色和权限信息）
func (r *userRepository) GetByIDWithRoles(ctx context.Context, id uint) (*user.User, error) {
	var u user.User
	if err := r.db.WithContext(ctx).
		Preload("Roles.Permissions").
		First(&u, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found with id: %d", id)
		}
		return nil, fmt.Errorf("failed to get user by id with roles: %w", err)
	}
	return &u, nil
}

// GetByUsernameWithRoles 根据用户名获取用户（包含角色和权限信息）
func (r *userRepository) GetByUsernameWithRoles(ctx context.Context, username string) (*user.User, error) {
	var u user.User
	if err := r.db.WithContext(ctx).
		Preload("Roles.Permissions").
		Where("username = ?", username).
		First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found with username: %s", username)
		}
		return nil, fmt.Errorf("failed to get user by username with roles: %w", err)
	}
	return &u, nil
}

// AssignRoles 为用户分配角色（替换现有角色）
func (r *userRepository) AssignRoles(ctx context.Context, userID uint, roleIDs []uint) error {
	var u user.User
	if err := r.db.WithContext(ctx).First(&u, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("user not found with id: %d", userID)
		}
		return fmt.Errorf("failed to find user: %w", err)
	}

	var roles []role.Role
	if err := r.db.WithContext(ctx).Find(&roles, roleIDs).Error; err != nil {
		return fmt.Errorf("failed to find roles: %w", err)
	}

	if err := r.db.WithContext(ctx).Model(&u).Association("Roles").Replace(roles); err != nil {
		return fmt.Errorf("failed to assign roles: %w", err)
	}

	return nil
}

// RemoveRoles 移除用户的角色
func (r *userRepository) RemoveRoles(ctx context.Context, userID uint, roleIDs []uint) error {
	var u user.User
	if err := r.db.WithContext(ctx).First(&u, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("user not found with id: %d", userID)
		}
		return fmt.Errorf("failed to find user: %w", err)
	}

	var roles []role.Role
	if err := r.db.WithContext(ctx).Find(&roles, roleIDs).Error; err != nil {
		return fmt.Errorf("failed to find roles: %w", err)
	}

	if err := r.db.WithContext(ctx).Model(&u).Association("Roles").Delete(roles); err != nil {
		return fmt.Errorf("failed to remove roles: %w", err)
	}

	return nil
}

// GetRoles 获取用户的所有角色ID
func (r *userRepository) GetRoles(ctx context.Context, userID uint) ([]uint, error) {
	var u user.User
	if err := r.db.WithContext(ctx).Preload("Roles").First(&u, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found with id: %d", userID)
		}
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	roleIDs := make([]uint, 0, len(u.Roles))
	for _, r := range u.Roles {
		roleIDs = append(roleIDs, r.ID)
	}

	return roleIDs, nil
}
