package persistence

import (
	"context"
	"errors"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/role"
	"gorm.io/gorm"
)

// RoleRepositoryImpl implements role.RoleRepository interface
type RoleRepositoryImpl struct {
	db *gorm.DB
}

// NewRoleRepository creates a new instance of RoleRepositoryImpl
func NewRoleRepository(db *gorm.DB) role.RoleRepository {
	return &RoleRepositoryImpl{db: db}
}

// Create creates a new role
func (r *RoleRepositoryImpl) Create(ctx context.Context, roleEntity *role.Role) error {
	return r.db.WithContext(ctx).Create(roleEntity).Error
}

// FindByID finds a role by ID with permissions preloaded
func (r *RoleRepositoryImpl) FindByID(ctx context.Context, id uint) (*role.Role, error) {
	var roleEntity role.Role
	err := r.db.WithContext(ctx).
		Preload("Permissions").
		First(&roleEntity, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &roleEntity, nil
}

// FindByName finds a role by name with permissions preloaded
func (r *RoleRepositoryImpl) FindByName(ctx context.Context, name string) (*role.Role, error) {
	var roleEntity role.Role
	err := r.db.WithContext(ctx).
		Preload("Permissions").
		Where("name = ?", name).
		First(&roleEntity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &roleEntity, nil
}

// List returns all roles with pagination
func (r *RoleRepositoryImpl) List(ctx context.Context, page, limit int) ([]role.Role, int64, error) {
	var roles []role.Role
	var total int64

	query := r.db.WithContext(ctx).Model(&role.Role{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := query.Preload("Permissions").
		Offset(offset).
		Limit(limit).
		Find(&roles).Error

	return roles, total, err
}

// Update updates a role
func (r *RoleRepositoryImpl) Update(ctx context.Context, roleEntity *role.Role) error {
	return r.db.WithContext(ctx).Save(roleEntity).Error
}

// Delete deletes a role (soft delete)
func (r *RoleRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&role.Role{}, id).Error
}

// GetPermissions retrieves all permissions for a role
func (r *RoleRepositoryImpl) GetPermissions(ctx context.Context, roleID uint) ([]role.Permission, error) {
	var roleEntity role.Role
	err := r.db.WithContext(ctx).
		Preload("Permissions").
		First(&roleEntity, roleID).Error
	if err != nil {
		return nil, err
	}
	return roleEntity.Permissions, nil
}

// SetPermissions sets permissions for a role (replaces existing permissions)
func (r *RoleRepositoryImpl) SetPermissions(ctx context.Context, roleID uint, permissionIDs []uint) error {
	var roleEntity role.Role
	if err := r.db.WithContext(ctx).First(&roleEntity, roleID).Error; err != nil {
		return err
	}

	var permissions []role.Permission
	if err := r.db.WithContext(ctx).Find(&permissions, permissionIDs).Error; err != nil {
		return err
	}

	return r.db.WithContext(ctx).Model(&roleEntity).Association("Permissions").Replace(permissions)
}

// AddPermissions adds permissions to a role
func (r *RoleRepositoryImpl) AddPermissions(ctx context.Context, roleID uint, permissionIDs []uint) error {
	var roleEntity role.Role
	if err := r.db.WithContext(ctx).First(&roleEntity, roleID).Error; err != nil {
		return err
	}

	var permissions []role.Permission
	if err := r.db.WithContext(ctx).Find(&permissions, permissionIDs).Error; err != nil {
		return err
	}

	return r.db.WithContext(ctx).Model(&roleEntity).Association("Permissions").Append(permissions)
}

// RemovePermissions removes permissions from a role
func (r *RoleRepositoryImpl) RemovePermissions(ctx context.Context, roleID uint, permissionIDs []uint) error {
	var roleEntity role.Role
	if err := r.db.WithContext(ctx).First(&roleEntity, roleID).Error; err != nil {
		return err
	}

	var permissions []role.Permission
	if err := r.db.WithContext(ctx).Find(&permissions, permissionIDs).Error; err != nil {
		return err
	}

	return r.db.WithContext(ctx).Model(&roleEntity).Association("Permissions").Delete(permissions)
}
