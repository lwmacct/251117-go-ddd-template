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

		// User domain - Settings management
		{Domain: "user", Resource: "settings", Action: "read", Code: "user:settings:read", Description: "Read own user settings"},
		{Domain: "user", Resource: "settings", Action: "update", Code: "user:settings:update", Description: "Update own user settings"},

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
// 优化：直接操作关联表，避免 N+1 问题
func (s *RBACSeeder) seedRoles(ctx context.Context, db *gorm.DB) error {
	db = db.WithContext(ctx)

	// 1. 获取所有权限（一次查询）
	var allPermissions []_persistence.PermissionModel
	if err := db.Find(&allPermissions).Error; err != nil {
		return err
	}

	// 按 domain 分组
	permByDomain := make(map[string][]_persistence.PermissionModel)
	for _, p := range allPermissions {
		permByDomain[p.Domain] = append(permByDomain[p.Domain], p)
	}

	// 2. 定义角色及其权限
	type roleConfig struct {
		name        string
		displayName string
		description string
		isSystem    bool
		permDomains []string // nil 表示所有权限
	}

	roles := []roleConfig{
		{"admin", "Administrator", "Full system access with all permissions", true, nil},
		{"user", "Regular User", "Standard user with limited permissions", true, []string{"user"}},
	}

	// 3. 批量创建角色并分配权限
	for _, r := range roles {
		role := _persistence.RoleModel{
			Name:        r.name,
			DisplayName: r.displayName,
			Description: r.description,
			IsSystem:    r.isSystem,
		}

		// Upsert 角色
		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}},
			DoUpdates: clause.AssignmentColumns([]string{"display_name", "description"}),
		}).Create(&role).Error; err != nil {
			return err
		}

		// 确定该角色需要的权限
		var perms []_persistence.PermissionModel
		if r.permDomains == nil {
			perms = allPermissions
		} else {
			for _, domain := range r.permDomains {
				perms = append(perms, permByDomain[domain]...)
			}
		}

		// 4. 直接批量插入关联表（跳过 Association API）
		if len(perms) > 0 {
			records := make([]map[string]any, 0, len(perms))
			for _, p := range perms {
				records = append(records, map[string]any{
					"role_model_id":       role.ID,
					"permission_model_id": p.ID,
				})
			}
			// 一次性插入所有关联，冲突时跳过
			if err := db.Table("role_permissions").Clauses(clause.OnConflict{DoNothing: true}).Create(&records).Error; err != nil {
				return err
			}
		}

		slog.Info("Role ensured", "name", role.Name, "permissions", len(perms))
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

	adminUser := _persistence.UserModel{
		Username: "admin",
		Email:    "admin@example.com",
		Password: string(hashedPassword),
		FullName: "System Administrator",
		Status:   "active",
	}

	// Upsert admin user
	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "username"}},
		DoNothing: true,
	}).Create(&adminUser)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected > 0 {
		slog.Info("Created admin user", "username", adminUser.Username)
		slog.Warn("Default admin credentials", "username", "admin", "password", "admin123", "warning", "PLEASE CHANGE THIS PASSWORD IMMEDIATELY")
	} else {
		// 如果用户已存在，需要获取其 ID
		if err := db.Where("username = ?", "admin").First(&adminUser).Error; err != nil {
			return err
		}
		slog.Info("Admin user ensured", "username", adminUser.Username)
	}

	// 直接插入用户-角色关联（跳过 Association API）
	userRole := map[string]any{
		"user_id": adminUser.ID,
		"role_id": adminRole.ID,
	}
	if err := db.Table("user_roles").Clauses(clause.OnConflict{DoNothing: true}).Create(&userRole).Error; err != nil {
		return err
	}

	slog.Info("Assigned admin role to admin user", "username", adminUser.Username, "role", adminRole.Name)
	return nil
}

// Name implements Seeder interface
func (s *RBACSeeder) Name() string {
	return "RBACSeeder"
}
