package seeds

import (
	"context"
	"log/slog"

	_persistence "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/persistence"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// RBACSeeder seeds roles, permissions, and admin user
type RBACSeeder struct{}

// Seed implements Seeder interface
func (s *RBACSeeder) Seed(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.seedPermissions(ctx, tx); err != nil {
			return err
		}

		if err := s.seedRoles(ctx, tx); err != nil {
			return err
		}

		return s.seedAdminUser(ctx, tx)
	})
}

// seedPermissions seeds initial permissions with three-part format: domain:resource:action
func (s *RBACSeeder) seedPermissions(ctx context.Context, db *gorm.DB) error {
	db = db.WithContext(ctx)

	permissions := []_persistence.PermissionModel{
		// Admin domain - User management
		{Domain: "admin", Resource: "users", Action: "create", Code: "admin:users:create", Description: "Create users"},
		{Domain: "admin", Resource: "users", Action: "read", Code: "admin:users:read", Description: "Read all users"},
		{Domain: "admin", Resource: "users", Action: "update", Code: "admin:users:update", Description: "Update any user"},
		{Domain: "admin", Resource: "users", Action: "delete", Code: "admin:users:delete", Description: "Delete users"},

		// Admin domain - Role management
		{Domain: "admin", Resource: "roles", Action: "create", Code: "admin:roles:create", Description: "Create roles"},
		{Domain: "admin", Resource: "roles", Action: "read", Code: "admin:roles:read", Description: "Read all roles"},
		{Domain: "admin", Resource: "roles", Action: "update", Code: "admin:roles:update", Description: "Update roles"},
		{Domain: "admin", Resource: "roles", Action: "delete", Code: "admin:roles:delete", Description: "Delete roles"},

		// Admin domain - Permission management
		{Domain: "admin", Resource: "permissions", Action: "read", Code: "admin:permissions:read", Description: "Read all permissions"},

		// Admin domain - Overview dashboard
		{Domain: "admin", Resource: "overview", Action: "read", Code: "admin:overview:read", Description: "View system overview stats"},

		// Admin domain - Menu management
		{Domain: "admin", Resource: "menus", Action: "create", Code: "admin:menus:create", Description: "Create menus"},
		{Domain: "admin", Resource: "menus", Action: "read", Code: "admin:menus:read", Description: "Read menus"},
		{Domain: "admin", Resource: "menus", Action: "update", Code: "admin:menus:update", Description: "Update menus"},
		{Domain: "admin", Resource: "menus", Action: "delete", Code: "admin:menus:delete", Description: "Delete menus"},

		// Admin domain - Settings management
		{Domain: "admin", Resource: "settings", Action: "create", Code: "admin:settings:create", Description: "Create settings"},
		{Domain: "admin", Resource: "settings", Action: "read", Code: "admin:settings:read", Description: "Read settings"},
		{Domain: "admin", Resource: "settings", Action: "update", Code: "admin:settings:update", Description: "Update settings"},
		{Domain: "admin", Resource: "settings", Action: "delete", Code: "admin:settings:delete", Description: "Delete settings"},

		// Admin domain - Audit log management
		{Domain: "admin", Resource: "audit_logs", Action: "read", Code: "admin:audit_logs:read", Description: "Read audit logs"},

		// User domain - Profile management
		{Domain: "user", Resource: "profile", Action: "read", Code: "user:profile:read", Description: "Read own profile"},
		{Domain: "user", Resource: "profile", Action: "update", Code: "user:profile:update", Description: "Update own profile"},
		{Domain: "user", Resource: "profile", Action: "delete", Code: "user:profile:delete", Description: "Delete own account"},

		// User domain - Password management
		{Domain: "user", Resource: "password", Action: "update", Code: "user:password:update", Description: "Change own password"},

		// User domain - Email management
		{Domain: "user", Resource: "email", Action: "update", Code: "user:email:update", Description: "Change own email"},

		// User domain - Token management
		{Domain: "user", Resource: "tokens", Action: "create", Code: "user:tokens:create", Description: "Create personal access tokens"},
		{Domain: "user", Resource: "tokens", Action: "read", Code: "user:tokens:read", Description: "List own tokens"},
		{Domain: "user", Resource: "tokens", Action: "update", Code: "user:tokens:disable", Description: "Disable own tokens"},
		{Domain: "user", Resource: "tokens", Action: "update", Code: "user:tokens:enable", Description: "Enable own tokens"},
		{Domain: "user", Resource: "tokens", Action: "delete", Code: "user:tokens:delete", Description: "Delete own tokens"},

		// API domain - Cache management (example for API endpoints)
		{Domain: "api", Resource: "cache", Action: "read", Code: "api:cache:read", Description: "Read cache data"},
		{Domain: "api", Resource: "cache", Action: "write", Code: "api:cache:write", Description: "Write cache data"},
		{Domain: "api", Resource: "cache", Action: "delete", Code: "api:cache:delete", Description: "Delete cache data"},
	}

	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "code"}},
		DoNothing: true,
	}).Create(&permissions)
	if result.Error != nil {
		return result.Error
	}

	slog.Info("Permissions ensured", "total", len(permissions), "inserted", result.RowsAffected)

	return nil
}

// seedRoles seeds initial roles with permissions

func (s *RBACSeeder) seedRoles(ctx context.Context, db *gorm.DB) error {
	db = db.WithContext(ctx)

	// Get all permissions
	var allPermissions []_persistence.PermissionModel
	if err := db.Find(&allPermissions).Error; err != nil {
		return err
	}

	// Find user permissions (user domain permissions only)
	var userPermissions []_persistence.PermissionModel
	if err := db.Where("domain = ?", "user").Find(&userPermissions).Error; err != nil {
		return err
	}

	roles := []struct {
		role        _persistence.RoleModel
		permissions []_persistence.PermissionModel
	}{
		{
			role: _persistence.RoleModel{
				Name:        "admin",
				DisplayName: "Administrator",
				Description: "Full system access with all permissions",
				IsSystem:    true,
			},
			permissions: allPermissions,
		},
		{
			role: _persistence.RoleModel{
				Name:        "user",
				DisplayName: "Regular User",
				Description: "Standard user with limited permissions",
				IsSystem:    true,
			},
			permissions: userPermissions,
		},
	}

	for _, r := range roles {
		role := r.role
		result := db.Where("name = ?", r.role.Name).
			Attrs(_persistence.RoleModel{
				DisplayName: r.role.DisplayName,
				Description: r.role.Description,
				IsSystem:    r.role.IsSystem,
			}).
			FirstOrCreate(&role)
		if result.Error != nil {
			return result.Error
		}

		association := db.Model(&role).Association("Permissions")
		var existingPermissions []_persistence.PermissionModel
		if err := association.Find(&existingPermissions); err != nil {
			return err
		}

		existing := make(map[uint]struct{}, len(existingPermissions))
		for _, perm := range existingPermissions {
			existing[perm.ID] = struct{}{}
		}

		var permsToAdd []*_persistence.PermissionModel
		for i := range r.permissions {
			if _, ok := existing[r.permissions[i].ID]; ok {
				continue
			}
			permsToAdd = append(permsToAdd, &r.permissions[i])
		}

		for _, perm := range permsToAdd {
			if err := association.Append(perm); err != nil {
				return err
			}
		}

		slog.Info("Role ensured", "name", role.Name, "added_permissions", len(permsToAdd), "required_total", len(r.permissions))
		if result.RowsAffected > 0 {
			slog.Info("Role created", "name", role.Name)
		}
	}

	return nil
}

// seedAdminUser seeds an initial admin user
func (s *RBACSeeder) seedAdminUser(ctx context.Context, db *gorm.DB) error {
	db = db.WithContext(ctx)

	// Get admin role
	var adminRole _persistence.RoleModel
	if err := db.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		return err
	}

	// Hash password for default user creation
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	adminSearch := _persistence.UserModel{Username: "admin"}
	result := db.Where(adminSearch).
		Attrs(_persistence.UserModel{
			Email:    "admin@example.com",
			Password: string(hashedPassword),
			FullName: "System Administrator",
			Status:   "active",
		}).
		FirstOrCreate(&adminSearch)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected > 0 {
		slog.Info("Created admin user", "username", adminSearch.Username)
		slog.Warn("Default admin credentials", "username", "admin", "password", "admin123", "warning", "PLEASE CHANGE THIS PASSWORD IMMEDIATELY")
	} else {
		slog.Info("Admin user ensured", "username", adminSearch.Username)
	}

	association := db.Model(&adminSearch).Association("Roles")
	var existingRoles []_persistence.RoleModel
	if err := association.Find(&existingRoles); err != nil {
		return err
	}

	hasAdminRole := false
	for _, role := range existingRoles {
		if role.ID == adminRole.ID {
			hasAdminRole = true
			break
		}
	}

	if !hasAdminRole {
		if err := association.Append(&adminRole); err != nil {
			return err
		}
		slog.Info("Assigned admin role to admin user", "username", adminSearch.Username, "role", adminRole.Name)
	} else {
		slog.Info("Admin user already has admin role", "username", adminSearch.Username)
	}

	return nil
}
