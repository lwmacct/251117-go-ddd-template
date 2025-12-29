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

	// SetPermissions sets permissions for a role (replaces all existing permissions)
	SetPermissions(ctx context.Context, roleID uint, permissions []Permission) error
}
