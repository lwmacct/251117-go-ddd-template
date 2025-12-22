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
		GenericCommandRepository: NewGenericCommandRepository[role.Role, *RoleModel](
			db, newRoleModelFromEntity,
		),
	}
}

// Create、Update、Delete 方法由 GenericCommandRepository 提供

// SetPermissions 设置角色权限 (替换现有权限)
func (r *roleCommandRepository) SetPermissions(ctx context.Context, roleID uint, permissionIDs []uint) error {
	return r.modifyPermissions(ctx, roleID, permissionIDs,
		func(assoc *gorm.Association, perms []PermissionModel) error { return assoc.Replace(perms) },
		"failed to set permissions")
}

// AddPermissions 为角色添加权限
func (r *roleCommandRepository) AddPermissions(ctx context.Context, roleID uint, permissionIDs []uint) error {
	return r.modifyPermissions(ctx, roleID, permissionIDs,
		func(assoc *gorm.Association, perms []PermissionModel) error { return assoc.Append(perms) },
		"failed to add permissions")
}

// RemovePermissions 移除角色权限
func (r *roleCommandRepository) RemovePermissions(ctx context.Context, roleID uint, permissionIDs []uint) error {
	return r.modifyPermissions(ctx, roleID, permissionIDs,
		func(assoc *gorm.Association, perms []PermissionModel) error { return assoc.Delete(perms) },
		"failed to remove permissions")
}

// modifyPermissions 通用权限操作方法，减少重复代码
func (r *roleCommandRepository) modifyPermissions(ctx context.Context, roleID uint, permissionIDs []uint, operation func(*gorm.Association, []PermissionModel) error, errMsg string) error {
	var roleModel RoleModel
	if err := r.DB().WithContext(ctx).First(&roleModel, roleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("role not found with id: %d", roleID)
		}
		return fmt.Errorf("failed to find role: %w", err)
	}

	var permissions []PermissionModel
	if err := r.DB().WithContext(ctx).Find(&permissions, permissionIDs).Error; err != nil {
		return fmt.Errorf("failed to find permissions: %w", err)
	}

	assoc := r.DB().WithContext(ctx).Model(&roleModel).Association("Permissions")
	if err := operation(assoc, permissions); err != nil {
		return fmt.Errorf("%s: %w", errMsg, err)
	}

	return nil
}
