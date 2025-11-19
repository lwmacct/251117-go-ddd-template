package role

import "context"

// RoleRepository defines the interface for role data operations
type RoleRepository interface {
	// Create creates a new role
	Create(ctx context.Context, role *Role) error

	// FindByID finds a role by ID
	FindByID(ctx context.Context, id uint) (*Role, error)

	// FindByName finds a role by name
	FindByName(ctx context.Context, name string) (*Role, error)

	// List returns all roles with pagination
	List(ctx context.Context, page, limit int) ([]Role, int64, error)

	// Update updates a role
	Update(ctx context.Context, role *Role) error

	// Delete deletes a role (soft delete)
	Delete(ctx context.Context, id uint) error

	// GetPermissions retrieves all permissions for a role
	GetPermissions(ctx context.Context, roleID uint) ([]Permission, error)

	// SetPermissions sets permissions for a role
	SetPermissions(ctx context.Context, roleID uint, permissionIDs []uint) error

	// AddPermissions adds permissions to a role
	AddPermissions(ctx context.Context, roleID uint, permissionIDs []uint) error

	// RemovePermissions removes permissions from a role
	RemovePermissions(ctx context.Context, roleID uint, permissionIDs []uint) error
}

// PermissionRepository defines the interface for permission data operations
type PermissionRepository interface {
	// Create creates a new permission
	Create(ctx context.Context, permission *Permission) error

	// FindByID finds a permission by ID
	FindByID(ctx context.Context, id uint) (*Permission, error)

	// FindByCode finds a permission by code
	FindByCode(ctx context.Context, code string) (*Permission, error)

	// List returns all permissions with pagination
	List(ctx context.Context, page, limit int) ([]Permission, int64, error)

	// ListByResource returns all permissions for a specific resource
	ListByResource(ctx context.Context, resource string) ([]Permission, error)

	// Update updates a permission
	Update(ctx context.Context, permission *Permission) error

	// Delete deletes a permission (soft delete)
	Delete(ctx context.Context, id uint) error

	// FindByIDs finds multiple permissions by their IDs
	FindByIDs(ctx context.Context, ids []uint) ([]Permission, error)
}
