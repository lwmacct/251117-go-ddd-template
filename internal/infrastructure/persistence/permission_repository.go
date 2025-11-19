package persistence

import (
	"context"
	"errors"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/role"
	"gorm.io/gorm"
)

// PermissionRepositoryImpl implements role.PermissionRepository interface
type PermissionRepositoryImpl struct {
	db *gorm.DB
}

// NewPermissionRepository creates a new instance of PermissionRepositoryImpl
func NewPermissionRepository(db *gorm.DB) role.PermissionRepository {
	return &PermissionRepositoryImpl{db: db}
}

// Create creates a new permission
func (p *PermissionRepositoryImpl) Create(ctx context.Context, permission *role.Permission) error {
	return p.db.WithContext(ctx).Create(permission).Error
}

// FindByID finds a permission by ID
func (p *PermissionRepositoryImpl) FindByID(ctx context.Context, id uint) (*role.Permission, error) {
	var permission role.Permission
	err := p.db.WithContext(ctx).First(&permission, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &permission, nil
}

// FindByCode finds a permission by code
func (p *PermissionRepositoryImpl) FindByCode(ctx context.Context, code string) (*role.Permission, error) {
	var permission role.Permission
	err := p.db.WithContext(ctx).Where("code = ?", code).First(&permission).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &permission, nil
}

// List returns all permissions with pagination
func (p *PermissionRepositoryImpl) List(ctx context.Context, page, limit int) ([]role.Permission, int64, error) {
	var permissions []role.Permission
	var total int64

	query := p.db.WithContext(ctx).Model(&role.Permission{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Find(&permissions).Error

	return permissions, total, err
}

// ListByResource returns all permissions for a specific resource
func (p *PermissionRepositoryImpl) ListByResource(ctx context.Context, resource string) ([]role.Permission, error) {
	var permissions []role.Permission
	err := p.db.WithContext(ctx).Where("resource = ?", resource).Find(&permissions).Error
	return permissions, err
}

// Update updates a permission
func (p *PermissionRepositoryImpl) Update(ctx context.Context, permission *role.Permission) error {
	return p.db.WithContext(ctx).Save(permission).Error
}

// Delete deletes a permission (soft delete)
func (p *PermissionRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return p.db.WithContext(ctx).Delete(&role.Permission{}, id).Error
}

// FindByIDs finds multiple permissions by their IDs
func (p *PermissionRepositoryImpl) FindByIDs(ctx context.Context, ids []uint) ([]role.Permission, error) {
	var permissions []role.Permission
	err := p.db.WithContext(ctx).Find(&permissions, ids).Error
	return permissions, err
}
