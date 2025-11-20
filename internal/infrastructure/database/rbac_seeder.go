package database

import (
	"context"
	"log/slog"

	_persistence "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/persistence"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// RBACSeeder seeds roles, permissions, and admin user
type RBACSeeder struct{}

// Seed implements Seeder interface
func (s *RBACSeeder) Seed(ctx context.Context, db *gorm.DB) error {
	if err := s.seedPermissions(ctx, db); err != nil {
		return err
	}

	if err := s.seedRoles(ctx, db); err != nil {
		return err
	}

	if err := s.seedAdminUser(ctx, db); err != nil {
		return err
	}

	return nil
}

// seedPermissions seeds initial permissions with three-part format: domain:resource:action
func (s *RBACSeeder) seedPermissions(ctx context.Context, db *gorm.DB) error {
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
		{Domain: "user", Resource: "tokens", Action: "delete", Code: "user:tokens:delete", Description: "Revoke own tokens"},

		// API domain - Cache management (example for API endpoints)
		{Domain: "api", Resource: "cache", Action: "read", Code: "api:cache:read", Description: "Read cache data"},
		{Domain: "api", Resource: "cache", Action: "write", Code: "api:cache:write", Description: "Write cache data"},
		{Domain: "api", Resource: "cache", Action: "delete", Code: "api:cache:delete", Description: "Delete cache data"},
	}

	for _, perm := range permissions {
		var existing _persistence.PermissionModel
		result := db.WithContext(ctx).Where("code = ?", perm.Code).First(&existing)
		if result.Error == gorm.ErrRecordNotFound {
			if err := db.WithContext(ctx).Create(&perm).Error; err != nil {
				return err
			}
			slog.Info("Created permission", "code", perm.Code)
		} else {
			slog.Info("Permission already exists", "code", perm.Code)
		}
	}

	return nil
}

// seedRoles seeds initial roles with permissions

func (s *RBACSeeder) seedRoles(ctx context.Context, db *gorm.DB) error {
	// Get all permissions
	var allPermissions []_persistence.PermissionModel
	if err := db.WithContext(ctx).Find(&allPermissions).Error; err != nil {
		return err
	}

	// Find user permissions (user domain permissions only)
	var userPermissions []_persistence.PermissionModel
	if err := db.WithContext(ctx).Where("domain = ?", "user").Find(&userPermissions).Error; err != nil {
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
		var existing _persistence.RoleModel
		result := db.WithContext(ctx).Where("name = ?", r.role.Name).First(&existing)
		if result.Error == gorm.ErrRecordNotFound {
			// Create role
			if err := db.WithContext(ctx).Create(&r.role).Error; err != nil {
				return err
			}
			slog.Info("Created role", "name", r.role.Name)

			// Assign permissions
			if len(r.permissions) > 0 {
				if err := db.WithContext(ctx).Model(&r.role).Association("Permissions").Append(r.permissions); err != nil {
					return err
				}
				slog.Info("Assigned permissions to role", "role", r.role.Name, "count", len(r.permissions))
			}
		} else {
			slog.Info("Role already exists", "name", r.role.Name)
		}
	}

	return nil
}

// seedAdminUser seeds an initial admin user
func (s *RBACSeeder) seedAdminUser(ctx context.Context, db *gorm.DB) error {
	// Get admin role
	var adminRole _persistence.RoleModel
	if err := db.WithContext(ctx).Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		return err
	}

	// Check if admin user already exists
	var adminUser _persistence.UserModel
	result := db.WithContext(ctx).Where("username = ?", "admin").First(&adminUser)

	if result.Error == gorm.ErrRecordNotFound {
		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		// Create admin user
		adminUser = _persistence.UserModel{
			Username: "admin",
			Email:    "admin@example.com",
			Password: string(hashedPassword),
			FullName: "System Administrator",
			Status:   "active",
		}

		if err := db.WithContext(ctx).Create(&adminUser).Error; err != nil {
			return err
		}
		slog.Info("Created admin user", "username", "admin")
		slog.Warn("Default admin credentials", "username", "admin", "password", "admin123", "warning", "PLEASE CHANGE THIS PASSWORD IMMEDIATELY")
	} else {
		slog.Info("Admin user already exists", "username", "admin")
	}

	// Ensure admin role is assigned (idempotent)
	var roleCount int64
	db.WithContext(ctx).Model(&adminUser).Association("Roles").Count()
	db.WithContext(ctx).Table("user_roles").
		Where("user_model_id = ? AND role_model_id = ?", adminUser.ID, adminRole.ID).
		Count(&roleCount)

	if roleCount == 0 {
		if err := db.WithContext(ctx).Model(&adminUser).Association("Roles").Append(&adminRole); err != nil {
			return err
		}
		slog.Info("Assigned admin role to admin user", "username", "admin", "role", "admin")
	} else {
		slog.Info("Admin user already has admin role", "username", "admin")
	}

	return nil
}
