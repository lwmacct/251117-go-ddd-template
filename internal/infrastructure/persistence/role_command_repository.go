// Package persistence 提供角色命令仓储的 GORM 实现
package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/role"
	"gorm.io/gorm"
)

// roleCommandRepository 角色命令仓储的 GORM 实现
type roleCommandRepository struct {
	db *gorm.DB
}

// NewRoleCommandRepository 创建角色命令仓储实例
func NewRoleCommandRepository(db *gorm.DB) role.CommandRepository {
	return &roleCommandRepository{db: db}
}

// Create 创建新角色
func (r *roleCommandRepository) Create(ctx context.Context, roleEntity *role.Role) error {
	if err := r.db.WithContext(ctx).Create(roleEntity).Error; err != nil {
		return fmt.Errorf("failed to create role: %w", err)
	}
	return nil
}

// Update 更新角色
func (r *roleCommandRepository) Update(ctx context.Context, roleEntity *role.Role) error {
	if err := r.db.WithContext(ctx).Save(roleEntity).Error; err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}
	return nil
}

// Delete 删除角色 (软删除)
func (r *roleCommandRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&role.Role{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}
	return nil
}

// SetPermissions 设置角色权限 (替换现有权限)
func (r *roleCommandRepository) SetPermissions(ctx context.Context, roleID uint, permissionIDs []uint) error {
	var roleEntity role.Role
	if err := r.db.WithContext(ctx).First(&roleEntity, roleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("role not found with id: %d", roleID)
		}
		return fmt.Errorf("failed to find role: %w", err)
	}

	var permissions []role.Permission
	if err := r.db.WithContext(ctx).Find(&permissions, permissionIDs).Error; err != nil {
		return fmt.Errorf("failed to find permissions: %w", err)
	}

	if err := r.db.WithContext(ctx).Model(&roleEntity).Association("Permissions").Replace(permissions); err != nil {
		return fmt.Errorf("failed to set permissions: %w", err)
	}

	return nil
}

// AddPermissions 为角色添加权限
func (r *roleCommandRepository) AddPermissions(ctx context.Context, roleID uint, permissionIDs []uint) error {
	var roleEntity role.Role
	if err := r.db.WithContext(ctx).First(&roleEntity, roleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("role not found with id: %d", roleID)
		}
		return fmt.Errorf("failed to find role: %w", err)
	}

	var permissions []role.Permission
	if err := r.db.WithContext(ctx).Find(&permissions, permissionIDs).Error; err != nil {
		return fmt.Errorf("failed to find permissions: %w", err)
	}

	if err := r.db.WithContext(ctx).Model(&roleEntity).Association("Permissions").Append(permissions); err != nil {
		return fmt.Errorf("failed to add permissions: %w", err)
	}

	return nil
}

// RemovePermissions 移除角色权限
func (r *roleCommandRepository) RemovePermissions(ctx context.Context, roleID uint, permissionIDs []uint) error {
	var roleEntity role.Role
	if err := r.db.WithContext(ctx).First(&roleEntity, roleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("role not found with id: %d", roleID)
		}
		return fmt.Errorf("failed to find role: %w", err)
	}

	var permissions []role.Permission
	if err := r.db.WithContext(ctx).Find(&permissions, permissionIDs).Error; err != nil {
		return fmt.Errorf("failed to find permissions: %w", err)
	}

	if err := r.db.WithContext(ctx).Model(&roleEntity).Association("Permissions").Delete(permissions); err != nil {
		return fmt.Errorf("failed to remove permissions: %w", err)
	}

	return nil
}
