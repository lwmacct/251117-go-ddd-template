package role

import "context"

// QueryRepository 角色查询仓储接口（读操作）
type QueryRepository interface {
	// FindByID finds a role by ID
	FindByID(ctx context.Context, id uint) (*Role, error)

	// FindByName finds a role by name
	FindByName(ctx context.Context, name string) (*Role, error)

	// FindByIDWithPermissions finds a role by ID with permissions
	FindByIDWithPermissions(ctx context.Context, id uint) (*Role, error)

	// List returns all roles with pagination
	List(ctx context.Context, page, limit int) ([]Role, int64, error)

	// GetPermissions retrieves all permissions for a role
	GetPermissions(ctx context.Context, roleID uint) ([]Permission, error)

	// Exists checks if a role exists by ID
	Exists(ctx context.Context, id uint) (bool, error)

	// ExistsByName checks if a role exists by name
	ExistsByName(ctx context.Context, name string) (bool, error)
}

// PermissionQueryRepository 权限查询仓储接口（读操作）
type PermissionQueryRepository interface {
	// FindByID finds a permission by ID
	FindByID(ctx context.Context, id uint) (*Permission, error)

	// FindByCode finds a permission by code
	FindByCode(ctx context.Context, code string) (*Permission, error)

	// FindByIDs finds multiple permissions by their IDs
	FindByIDs(ctx context.Context, ids []uint) ([]Permission, error)

	// List returns all permissions with pagination
	List(ctx context.Context, page, limit int) ([]Permission, int64, error)

	// ListByResource returns all permissions for a specific resource
	ListByResource(ctx context.Context, resource string) ([]Permission, error)

	// Exists checks if a permission exists by ID
	Exists(ctx context.Context, id uint) (bool, error)

	// ExistsByCode checks if a permission exists by code
	ExistsByCode(ctx context.Context, code string) (bool, error)
}
