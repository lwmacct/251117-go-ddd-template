package seeds

import "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/database"

// DefaultSeeders returns the default ordered seeders that bootstrap the system.
// Keep RBAC first because it provisions permissions/roles required by other seeders.
func DefaultSeeders() []database.Seeder {
	return []database.Seeder{
		&RBACSeeder{},
		&UserSeeder{},
	}
}
