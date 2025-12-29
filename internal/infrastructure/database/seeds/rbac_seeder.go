package seeds

import (
	"context"
	"log/slog"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/operation"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/persistence"
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

	// 从统一操作注册表获取所有权限定义
	defs := operation.AllPermissions()
	permissions := make([]persistence.PermissionModel, 0, len(defs))
	for _, def := range defs {
		permissions = append(permissions, persistence.PermissionModel{
			Domain:      def.Domain,
			Resource:    def.Resource,
			Action:      def.Action,
			Code:        def.Code,
			Description: def.Description,
		})
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
	var allPermissions []persistence.PermissionModel
	if err := db.Find(&allPermissions).Error; err != nil {
		return err
	}

	// 按 domain 分组
	permByDomain := make(map[string][]persistence.PermissionModel)
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
		{"user", "Regular User", "Standard user with limited permissions", true, []string{"user", "auth"}},
	}

	// 3. 批量创建角色并分配权限
	for _, r := range roles {
		role := persistence.RoleModel{
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
		var perms []persistence.PermissionModel
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
	var adminRole persistence.RoleModel
	if err := db.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		return err
	}

	// Hash password for default user creation
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	adminUser := persistence.UserModel{
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
