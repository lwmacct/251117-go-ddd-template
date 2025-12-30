package seeds

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
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

// permissionJSON 权限 JSONB 序列化结构
type permissionJSON struct {
	OperationPattern string `json:"operation_pattern"`
	ResourcePattern  string `json:"resource_pattern"`
}

// seedRoles seeds initial roles with operation patterns
func (s *RBACSeeder) seedRoles(ctx context.Context, db *gorm.DB) error {
	db = db.WithContext(ctx)

	// 定义角色及其权限模式
	// URN-Centric RBAC: 使用 URN 格式 {scope}:{type}:{identifier}
	type roleConfig struct {
		name        string
		displayName string
		description string
		isSystem    bool
		permissions []permissionJSON
	}

	roles := []roleConfig{
		{
			name:        "admin",
			displayName: "系统管理员",
			description: "完整系统访问权限",
			isSystem:    true,
			permissions: []permissionJSON{
				// 超级管理员：所有操作对所有资源
				{OperationPattern: "*:*:*", ResourcePattern: "*:*:*"},
			},
		},
		{
			name:        "user",
			displayName: "普通用户",
			description: "标准用户权限",
			isSystem:    true,
			permissions: []permissionJSON{
				// self 域所有操作（对自己的资源）
				{OperationPattern: "self:*:*", ResourcePattern: "self:user:@me"},
			},
		},
	}

	// 创建角色并分配权限
	for _, r := range roles {
		// 序列化权限为 JSONB
		permsData, err := json.Marshal(r.permissions)
		if err != nil {
			return err
		}

		role := persistence.RoleModel{
			Name:        r.name,
			DisplayName: r.displayName,
			Description: r.description,
			IsSystem:    r.isSystem,
			Permissions: permsData,
		}

		// Upsert 角色（包含权限 JSONB）
		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}},
			DoUpdates: clause.AssignmentColumns([]string{"display_name", "description", "permissions"}),
		}).Create(&role).Error; err != nil {
			return err
		}

		slog.Info("Role ensured", "name", role.Name, "permissions", len(r.permissions))
	}

	return nil
}

// seedAdminUser seeds initial system users (root and admin) and demo users
func (s *RBACSeeder) seedAdminUser(ctx context.Context, db *gorm.DB) error {
	db = db.WithContext(ctx)

	// Get admin role
	var adminRole persistence.RoleModel
	if err := db.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		return err
	}

	// Get user role
	var userRole persistence.RoleModel
	if err := db.Where("name = ?", "user").First(&userRole).Error; err != nil {
		return err
	}

	// Hash password for default user creation
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 定义用户配置
	type userConfig struct {
		username string
		email    string
		fullName string
		isSystem bool
		role     *persistence.RoleModel // nil 表示不分配角色（root 用户硬编码权限）
	}

	users := []userConfig{
		// 系统用户
		{
			username: user.RootUsername,
			email:    "root@localhost",
			fullName: "Root Administrator",
			isSystem: true,
			role:     nil, // root 用户权限硬编码，不需要角色
		},
		{
			username: user.AdminUsername,
			email:    "admin@example.com",
			fullName: "System Administrator",
			isSystem: true,
			role:     &adminRole,
		},
		// 普通用户（用于测试）
		{
			username: "human",
			email:    "human@example.com",
			fullName: "Human User",
			isSystem: false,
			role:     &userRole,
		},
	}

	for _, cfg := range users {
		userModel := persistence.UserModel{
			Username: cfg.username,
			Email:    cfg.email,
			Password: string(hashedPassword),
			FullName: cfg.fullName,
			Status:   "active",
			Type:     string(user.UserTypeHuman),
			IsSystem: cfg.isSystem,
		}

		// Upsert 用户
		result := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "username"}},
			DoUpdates: clause.AssignmentColumns([]string{"is_system", "type"}), // 确保标记被更新
		}).Create(&userModel)
		if result.Error != nil {
			return result.Error
		}

		// 用户已存在时，获取其 ID
		if result.RowsAffected == 0 {
			if err := db.Where("username = ?", cfg.username).First(&userModel).Error; err != nil {
				return err
			}
			slog.Info("User ensured", "username", userModel.Username, "is_system", userModel.IsSystem)
		} else {
			slog.Info("Created user", "username", userModel.Username, "is_system", cfg.isSystem)
		}

		// 系统用户提示修改默认密码
		if result.RowsAffected > 0 && cfg.isSystem {
			slog.Warn("Default credentials", "username", cfg.username, "password", "admin123", "warning", "PLEASE CHANGE THIS PASSWORD IMMEDIATELY")
		}

		// 分配角色（如果配置了角色）
		if cfg.role != nil {
			userRole := map[string]any{
				"user_id": userModel.ID,
				"role_id": cfg.role.ID,
			}
			if err := db.Table("user_roles").Clauses(clause.OnConflict{DoNothing: true}).Create(&userRole).Error; err != nil {
				return err
			}
			slog.Info("Assigned role to user", "username", userModel.Username, "role", cfg.role.Name)
		}
	}

	return nil
}

// Name implements Seeder interface
func (s *RBACSeeder) Name() string {
	return "RBACSeeder"
}
