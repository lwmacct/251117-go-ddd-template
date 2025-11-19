// Package persistence 提供权限查询仓储的 GORM 实现
package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/role"
	"gorm.io/gorm"
)

// permissionQueryRepository 权限查询仓储的 GORM 实现
type permissionQueryRepository struct {
	db *gorm.DB
}

// NewPermissionQueryRepository 创建权限查询仓储实例
func NewPermissionQueryRepository(db *gorm.DB) role.PermissionQueryRepository {
	return &permissionQueryRepository{db: db}
}

// FindByID 根据 ID 查找权限
func (p *permissionQueryRepository) FindByID(ctx context.Context, id uint) (*role.Permission, error) {
	var model PermissionModel
	err := p.db.WithContext(ctx).First(&model, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find permission by id: %w", err)
	}
	return model.toEntity(), nil
}

// FindByCode 根据代码查找权限
func (p *permissionQueryRepository) FindByCode(ctx context.Context, code string) (*role.Permission, error) {
	var model PermissionModel
	err := p.db.WithContext(ctx).Where("code = ?", code).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find permission by code: %w", err)
	}
	return model.toEntity(), nil
}

// FindByIDs 根据 ID 列表查找多个权限
func (p *permissionQueryRepository) FindByIDs(ctx context.Context, ids []uint) ([]role.Permission, error) {
	var models []PermissionModel
	err := p.db.WithContext(ctx).Find(&models, ids).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find permissions by ids: %w", err)
	}
	return mapPermissionModelsToEntities(models), nil
}

// List 获取权限列表 (分页)
func (p *permissionQueryRepository) List(ctx context.Context, page, limit int) ([]role.Permission, int64, error) {
	var models []PermissionModel
	var total int64

	query := p.db.WithContext(ctx).Model(&PermissionModel{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count permissions: %w", err)
	}

	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Find(&models).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list permissions: %w", err)
	}

	return mapPermissionModelsToEntities(models), total, nil
}

// ListByResource 根据资源获取权限列表
func (p *permissionQueryRepository) ListByResource(ctx context.Context, resource string) ([]role.Permission, error) {
	var models []PermissionModel
	err := p.db.WithContext(ctx).Where("resource = ?", resource).Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list permissions by resource: %w", err)
	}
	return mapPermissionModelsToEntities(models), nil
}

// Exists 检查权限是否存在
func (p *permissionQueryRepository) Exists(ctx context.Context, id uint) (bool, error) {
	var count int64
	if err := p.db.WithContext(ctx).Model(&PermissionModel{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check permission existence: %w", err)
	}
	return count > 0, nil
}

// ExistsByCode 检查权限代码是否存在
func (p *permissionQueryRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	var count int64
	if err := p.db.WithContext(ctx).Model(&PermissionModel{}).Where("code = ?", code).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check permission code existence: %w", err)
	}
	return count > 0, nil
}
