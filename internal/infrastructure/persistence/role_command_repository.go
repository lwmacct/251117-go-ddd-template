package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/role"
	"gorm.io/gorm"
)

// roleCommandRepository 角色命令仓储的 GORM 实现
// 嵌入 GenericCommandRepository 以复用基础 CRUD 操作
type roleCommandRepository struct {
	*GenericCommandRepository[role.Role, *RoleModel]
}

// NewRoleCommandRepository 创建角色命令仓储实例
func NewRoleCommandRepository(db *gorm.DB) role.CommandRepository {
	return &roleCommandRepository{
		GenericCommandRepository: NewGenericCommandRepository(
			db, newRoleModelFromEntity,
		),
	}
}

// Create、Update、Delete 方法由 GenericCommandRepository 提供

// SetPermissions 设置角色权限 (替换现有权限)
func (r *roleCommandRepository) SetPermissions(ctx context.Context, roleID uint, permissions []role.Permission) error {
	// 验证角色存在
	var roleModel RoleModel
	if err := r.DB().WithContext(ctx).First(&roleModel, roleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("role not found with id: %d", roleID)
		}
		return fmt.Errorf("failed to find role: %w", err)
	}

	// 使用事务处理
	return r.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除现有权限
		if err := tx.Where("role_id = ?", roleID).Delete(&RolePermissionModel{}).Error; err != nil {
			return fmt.Errorf("failed to delete existing permissions: %w", err)
		}

		// 如果没有新权限，直接返回
		if len(permissions) == 0 {
			return nil
		}

		// 插入新权限
		models := mapPermissionEntitiesToRolePermissionModels(roleID, permissions)
		if err := tx.Create(&models).Error; err != nil {
			return fmt.Errorf("failed to create permissions: %w", err)
		}

		return nil
	})
}
