package role

import "context"

// CommandRepository 角色命令仓储接口（写操作）
type CommandRepository interface {
	// Create creates a new role
	Create(ctx context.Context, role *Role) error

	// Update updates a role
	Update(ctx context.Context, role *Role) error

	// Delete deletes a role (soft delete)
	Delete(ctx context.Context, id uint) error

	// SetPermissions sets permissions for a role
	SetPermissions(ctx context.Context, roleID uint, permissionIDs []uint) error

	// AddPermissions adds permissions to a role
	AddPermissions(ctx context.Context, roleID uint, permissionIDs []uint) error

	// RemovePermissions removes permissions from a role
	RemovePermissions(ctx context.Context, roleID uint, permissionIDs []uint) error
}

// PermissionCommandRepository 权限命令仓储接口（写操作）
type PermissionCommandRepository interface {
	// Create creates a new permission
	Create(ctx context.Context, permission *Permission) error

	// Update updates a permission
	Update(ctx context.Context, permission *Permission) error

	// Delete deletes a permission (soft delete)
	Delete(ctx context.Context, id uint) error
}
