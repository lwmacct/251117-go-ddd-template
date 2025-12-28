package container

import (
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/persistence"
)

// ContainerOptions holds options for container initialization.
type ContainerOptions struct {
	AutoMigrate bool // Whether to auto-migrate database (recommended only for development)
}

// DefaultOptions returns default container options.
func DefaultOptions() *ContainerOptions {
	return &ContainerOptions{
		AutoMigrate: false, // Production default: no auto-migration
	}
}

// GetAllModels returns all domain models that need migration.
// When adding new domain models, register them here.
func GetAllModels() []any {
	return []any{
		&persistence.UserModel{},
		&persistence.RoleModel{},
		&persistence.PermissionModel{},
		&persistence.PersonalAccessTokenModel{},
		&persistence.AuditLogModel{},
		&persistence.TwoFAModel{},
		&persistence.MenuModel{},
		&persistence.SettingModel{},
		&persistence.SettingCategoryModel{},
		&persistence.UserSettingModel{},
	}
}
