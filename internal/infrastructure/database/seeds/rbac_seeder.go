package seeds

import (
	"context"
	"log/slog"

	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/persistence"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// RBACSeeder seeds roles and admin user
// 使用 URN-Centric RBAC：角色直接关联 Operation/Resource URN 模式
type RBACSeeder struct{}

// Seed implements Seeder interface
func (s *RBACSeeder) Seed(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.seedRoles(ctx, tx); err != nil {
			return err
		}

		return s.seedAdminUser(ctx, tx)
	})
}

// seedRoles seeds initial roles with operation patterns
func (s *RBACSeeder) seedRoles(ctx context.Context, db *gorm.DB) error {
	db = db.WithContext(ctx)

	// 定义角色及其权限模式
	// URN-Centric RBAC: 使用 URN 格式 {scope}:{type}:{identifier}
	type permissionConfig struct {
		operationPattern string // URN 格式，如 sys:users:*, self:*:*
		resourcePattern  string // URN 格式，如 *:*:*, self:user:@me
	}

	type roleConfig struct {
		name        string
		displayName string
		description string
		isSystem    bool
		permissions []permissionConfig
	}

	roles := []roleConfig{
		{
			name:        "admin",
			displayName: "系统管理员",
			description: "完整系统访问权限",
			isSystem:    true,
			permissions: []permissionConfig{
				// 超级管理员：所有操作对所有资源
				{operationPattern: "*:*:*", resourcePattern: "*:*:*"},
			},
		},
		{
			name:        "user",
			displayName: "普通用户",
			description: "标准用户权限",
			isSystem:    true,
			permissions: []permissionConfig{
				// self 域所有操作（对自己的资源）
				{operationPattern: "self:*:*", resourcePattern: "self:user:@me"},
			},
		},
	}

	// 创建角色并分配权限
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

		// 删除现有权限（确保幂等）
		if err := db.Where("role_id = ?", role.ID).Delete(&persistence.RolePermissionModel{}).Error; err != nil {
			return err
		}

		// 插入新权限
		if len(r.permissions) > 0 {
			perms := make([]persistence.RolePermissionModel, len(r.permissions))
			for i, p := range r.permissions {
				resPattern := p.resourcePattern
				if resPattern == "" {
					resPattern = "*:*:*"
				}
				perms[i] = persistence.RolePermissionModel{
					RoleID:           role.ID,
					OperationPattern: p.operationPattern,
					ResourcePattern:  resPattern,
				}
			}

			if err := db.Create(&perms).Error; err != nil {
				return err
			}
		}

		slog.Info("Role ensured", "name", role.Name, "permissions", len(r.permissions))
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
