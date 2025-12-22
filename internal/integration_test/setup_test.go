package integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	auditlogQuery "github.com/lwmacct/251117-go-ddd-template/internal/application/auditlog/query"
	authCommand "github.com/lwmacct/251117-go-ddd-template/internal/application/auth/command"
	roleCommand "github.com/lwmacct/251117-go-ddd-template/internal/application/role/command"
	roleQuery "github.com/lwmacct/251117-go-ddd-template/internal/application/role/query"
	userCommand "github.com/lwmacct/251117-go-ddd-template/internal/application/user/command"
	userQuery "github.com/lwmacct/251117-go-ddd-template/internal/application/user/query"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/auditlog"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/auth"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/role"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
	infraAuth "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/auth"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/persistence"
)

// TestEnv 集成测试环境
type TestEnv struct {
	DB *gorm.DB

	// Repositories
	UserCommandRepo     user.CommandRepository
	UserQueryRepo       user.QueryRepository
	RoleCommandRepo     role.CommandRepository
	RoleQueryRepo       role.QueryRepository
	AuditLogCommandRepo auditlog.CommandRepository
	AuditLogQueryRepo   auditlog.QueryRepository

	// Services
	AuthService auth.Service

	// User Handlers
	RegisterHandler    *authCommand.RegisterHandler
	CreateUserHandler  *userCommand.CreateUserHandler
	GetUserHandler     *userQuery.GetUserHandler
	ListUsersHandler   *userQuery.ListUsersHandler
	AssignRolesHandler *userCommand.AssignRolesHandler

	// Role Handlers
	CreateRoleHandler *roleCommand.CreateRoleHandler
	UpdateRoleHandler *roleCommand.UpdateRoleHandler
	DeleteRoleHandler *roleCommand.DeleteRoleHandler
	GetRoleHandler    *roleQuery.GetRoleHandler
	ListRolesHandler  *roleQuery.ListRolesHandler

	// AuditLog Handlers
	ListLogsHandler *auditlogQuery.ListLogsHandler
}

// SetupTestEnv 创建集成测试环境
func SetupTestEnv(t *testing.T) *TestEnv {
	t.Helper()

	// 1. 创建内存数据库
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err, "failed to create test database")

	// 2. 迁移所有表
	err = db.AutoMigrate(
		&persistence.UserModel{},
		&persistence.RoleModel{},
		&persistence.PermissionModel{},
		&persistence.PersonalAccessTokenModel{},
		&persistence.AuditLogModel{},
		&persistence.TwoFAModel{},
		&persistence.MenuModel{},
		&persistence.SettingModel{},
	)
	require.NoError(t, err, "failed to migrate database")

	// 3. 初始化 Repositories
	userRepos := persistence.NewUserRepositories(db)
	roleRepos := persistence.NewRoleRepositories(db)
	auditLogRepos := persistence.NewAuditLogRepositories(db)

	// 4. 初始化 Services
	jwtManager := infraAuth.NewJWTManager("test-secret", time.Hour, 24*time.Hour)
	tokenGenerator := infraAuth.NewTokenGenerator()
	passwordPolicy := auth.DefaultPasswordPolicy()
	authService := infraAuth.NewAuthService(jwtManager, tokenGenerator, passwordPolicy)

	// 5. 初始化 User Handlers
	registerHandler := authCommand.NewRegisterHandler(userRepos.Command, userRepos.Query, authService)
	createUserHandler := userCommand.NewCreateUserHandler(userRepos.Command, userRepos.Query, authService)
	getUserHandler := userQuery.NewGetUserHandler(userRepos.Query)
	listUsersHandler := userQuery.NewListUsersHandler(userRepos.Query)
	assignRolesHandler := userCommand.NewAssignRolesHandler(userRepos.Command, userRepos.Query)

	// 6. 初始化 Role Handlers
	createRoleHandler := roleCommand.NewCreateRoleHandler(roleRepos.Command, roleRepos.Query)
	updateRoleHandler := roleCommand.NewUpdateRoleHandler(roleRepos.Command, roleRepos.Query)
	deleteRoleHandler := roleCommand.NewDeleteRoleHandler(roleRepos.Command, roleRepos.Query)
	getRoleHandler := roleQuery.NewGetRoleHandler(roleRepos.Query)
	listRolesHandler := roleQuery.NewListRolesHandler(roleRepos.Query)

	// 7. 初始化 AuditLog Handlers
	listLogsHandler := auditlogQuery.NewListLogsHandler(auditLogRepos.Query)

	return &TestEnv{
		DB:                  db,
		UserCommandRepo:     userRepos.Command,
		UserQueryRepo:       userRepos.Query,
		RoleCommandRepo:     roleRepos.Command,
		RoleQueryRepo:       roleRepos.Query,
		AuditLogCommandRepo: auditLogRepos.Command,
		AuditLogQueryRepo:   auditLogRepos.Query,
		AuthService:         authService,
		RegisterHandler:     registerHandler,
		CreateUserHandler:   createUserHandler,
		GetUserHandler:      getUserHandler,
		ListUsersHandler:    listUsersHandler,
		AssignRolesHandler:  assignRolesHandler,
		CreateRoleHandler:   createRoleHandler,
		UpdateRoleHandler:   updateRoleHandler,
		DeleteRoleHandler:   deleteRoleHandler,
		GetRoleHandler:      getRoleHandler,
		ListRolesHandler:    listRolesHandler,
		ListLogsHandler:     listLogsHandler,
	}
}

// Cleanup 清理测试环境
func (e *TestEnv) Cleanup() {
	if e.DB != nil {
		sqlDB, _ := e.DB.DB()
		if sqlDB != nil {
			_ = sqlDB.Close() // 显式忽略错误，测试清理阶段
		}
	}
}

// CreateTestRole 创建测试角色
func (e *TestEnv) CreateTestRole(t *testing.T, name string, permissions []string) *role.Role {
	t.Helper()

	r := &role.Role{
		Name:        name,
		DisplayName: name + " Display",
		Description: "Test role: " + name,
	}

	for i, code := range permissions {
		r.Permissions = append(r.Permissions, role.Permission{
			ID:   uint(i + 1), //nolint:gosec // test data with small indices
			Code: code,
		})
	}

	err := e.RoleCommandRepo.Create(context.Background(), r)
	require.NoError(t, err)
	return r
}

// CreateTestUser 创建测试用户
func (e *TestEnv) CreateTestUser(ctx context.Context, t *testing.T, username, email, password string) *user.User {
	t.Helper()

	result, err := e.RegisterHandler.Handle(ctx, authCommand.RegisterCommand{
		Username: username,
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)

	u, err := e.UserQueryRepo.GetByID(ctx, result.UserID)
	require.NoError(t, err)

	return u
}
