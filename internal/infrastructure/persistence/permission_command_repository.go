// Package persistence 提供权限命令仓储的 GORM 实现
package persistence

import (
	"context"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/role"
	"gorm.io/gorm"
)

// permissionCommandRepository 权限命令仓储的 GORM 实现
type permissionCommandRepository struct {
	db *gorm.DB
}

// NewPermissionCommandRepository 创建权限命令仓储实例
func NewPermissionCommandRepository(db *gorm.DB) role.PermissionCommandRepository {
	return &permissionCommandRepository{db: db}
}

// Create 创建新权限
func (p *permissionCommandRepository) Create(ctx context.Context, permission *role.Permission) error {
	model := newPermissionModelFromEntity(permission)
	if err := p.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to create permission: %w", err)
	}
	if saved := model.toEntity(); saved != nil {
		*permission = *saved
	}
	return nil
}

// Update 更新权限
func (p *permissionCommandRepository) Update(ctx context.Context, permission *role.Permission) error {
	model := newPermissionModelFromEntity(permission)
	if err := p.db.WithContext(ctx).Save(model).Error; err != nil {
		return fmt.Errorf("failed to update permission: %w", err)
	}
	if saved := model.toEntity(); saved != nil {
		*permission = *saved
	}
	return nil
}

// Delete 删除权限 (软删除)
func (p *permissionCommandRepository) Delete(ctx context.Context, id uint) error {
	if err := p.db.WithContext(ctx).Delete(&PermissionModel{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete permission: %w", err)
	}
	return nil
}
