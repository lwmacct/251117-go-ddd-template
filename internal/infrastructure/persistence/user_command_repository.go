// Package persistence 提供用户命令仓储的 GORM 实现
package persistence

import (
	"context"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
	"gorm.io/gorm"
)

// userCommandRepository 用户命令仓储的 GORM 实现
type userCommandRepository struct {
	db *gorm.DB
}

// NewUserCommandRepository 创建用户命令仓储实例
func NewUserCommandRepository(db *gorm.DB) user.CommandRepository {
	return &userCommandRepository{db: db}
}

// Create 创建用户
func (r *userCommandRepository) Create(ctx context.Context, u *user.User) error {
	model := newUserModelFromEntity(u)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	if saved := model.toEntity(); saved != nil {
		*u = *saved
	}
	return nil
}

// Update 更新用户
func (r *userCommandRepository) Update(ctx context.Context, u *user.User) error {
	model := newUserModelFromEntity(u)
	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	if saved := model.toEntity(); saved != nil {
		*u = *saved
	}
	return nil
}

// Delete 删除用户 (软删除)
func (r *userCommandRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&UserModel{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// AssignRoles 为用户分配角色（替换现有角色）
func (r *userCommandRepository) AssignRoles(ctx context.Context, userID uint, roleIDs []uint) error {
	var u UserModel
	if err := r.db.WithContext(ctx).First(&u, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return user.ErrUserNotFound
		}
		return fmt.Errorf("failed to find user: %w", err)
	}

	var roles []RoleModel
	if err := r.db.WithContext(ctx).Find(&roles, roleIDs).Error; err != nil {
		return fmt.Errorf("failed to find roles: %w", err)
	}

	if err := r.db.WithContext(ctx).Model(&u).Association("Roles").Replace(roles); err != nil {
		return fmt.Errorf("failed to assign roles: %w", err)
	}

	return nil
}

// RemoveRoles 移除用户的角色
func (r *userCommandRepository) RemoveRoles(ctx context.Context, userID uint, roleIDs []uint) error {
	var u UserModel
	if err := r.db.WithContext(ctx).First(&u, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return user.ErrUserNotFound
		}
		return fmt.Errorf("failed to find user: %w", err)
	}

	var roles []RoleModel
	if err := r.db.WithContext(ctx).Find(&roles, roleIDs).Error; err != nil {
		return fmt.Errorf("failed to find roles: %w", err)
	}

	if err := r.db.WithContext(ctx).Model(&u).Association("Roles").Delete(roles); err != nil {
		return fmt.Errorf("failed to remove roles: %w", err)
	}

	return nil
}

// UpdatePassword 更新用户密码
func (r *userCommandRepository) UpdatePassword(ctx context.Context, userID uint, hashedPassword string) error {
	if err := r.db.WithContext(ctx).Model(&UserModel{}).
		Where("id = ?", userID).
		Update("password", hashedPassword).Error; err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}
	return nil
}

// UpdateStatus 更新用户状态
func (r *userCommandRepository) UpdateStatus(ctx context.Context, userID uint, status string) error {
	if err := r.db.WithContext(ctx).Model(&UserModel{}).
		Where("id = ?", userID).
		Update("status", status).Error; err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}
	return nil
}
