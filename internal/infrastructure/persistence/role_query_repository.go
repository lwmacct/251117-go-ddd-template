package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/role"
	"gorm.io/gorm"
)

// roleQueryRepository 角色查询仓储的 GORM 实现
type roleQueryRepository struct {
	db *gorm.DB
}

// NewRoleQueryRepository 创建角色查询仓储实例
func NewRoleQueryRepository(db *gorm.DB) role.QueryRepository {
	return &roleQueryRepository{db: db}
}

// FindByID 根据 ID 查找角色
func (r *roleQueryRepository) FindByID(ctx context.Context, id uint) (*role.Role, error) {
	var model RoleModel
	err := r.db.WithContext(ctx).
		First(&model, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil //nolint:nilnil // returns nil for not found, valid pattern
		}
		return nil, fmt.Errorf("failed to find role by id: %w", err)
	}
	return model.ToEntity(), nil
}

// FindByName 根据名称查找角色
func (r *roleQueryRepository) FindByName(ctx context.Context, name string) (*role.Role, error) {
	var model RoleModel
	err := r.db.WithContext(ctx).
		Where("name = ?", name).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil //nolint:nilnil // returns nil for not found, valid pattern
		}
		return nil, fmt.Errorf("failed to find role by name: %w", err)
	}
	return model.ToEntity(), nil
}

// FindByIDWithPermissions 根据 ID 查找角色 (包含权限)
func (r *roleQueryRepository) FindByIDWithPermissions(ctx context.Context, id uint) (*role.Role, error) {
	var model RoleModel
	err := r.db.WithContext(ctx).
		Preload("Permissions").
		First(&model, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil //nolint:nilnil // returns nil for not found, valid pattern
		}
		return nil, fmt.Errorf("failed to find role with permissions: %w", err)
	}
	return model.ToEntity(), nil
}

// List 获取角色列表 (分页)
func (r *roleQueryRepository) List(ctx context.Context, page, limit int) ([]role.Role, int64, error) {
	var models []RoleModel
	var total int64

	query := r.db.WithContext(ctx).Model(&RoleModel{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count roles: %w", err)
	}

	offset := (page - 1) * limit
	err := query.Preload("Permissions").
		Offset(offset).
		Limit(limit).
		Find(&models).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list roles: %w", err)
	}

	return mapRoleModelsToEntities(models), total, nil
}

// GetPermissions 获取角色的所有权限
func (r *roleQueryRepository) GetPermissions(ctx context.Context, roleID uint) ([]role.Permission, error) {
	var roleModel RoleModel
	err := r.db.WithContext(ctx).
		Preload("Permissions").
		First(&roleModel, roleID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("role not found with id: %d", roleID)
		}
		return nil, fmt.Errorf("failed to get role permissions: %w", err)
	}
	return mapPermissionModelsToEntities(roleModel.Permissions), nil
}

// Exists 检查角色是否存在
func (r *roleQueryRepository) Exists(ctx context.Context, id uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&RoleModel{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check role existence: %w", err)
	}
	return count > 0, nil
}

// ExistsByName 检查角色名称是否存在
func (r *roleQueryRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&RoleModel{}).Where("name = ?", name).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check role name existence: %w", err)
	}
	return count > 0, nil
}
