package bootstrap

import (
	"context"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http"
	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/handler"
	authCommand "github.com/lwmacct/251117-go-ddd-template/internal/application/auth/command"
	menuCommand "github.com/lwmacct/251117-go-ddd-template/internal/application/menu/command"
	menuQuery "github.com/lwmacct/251117-go-ddd-template/internal/application/menu/query"
	patCommand "github.com/lwmacct/251117-go-ddd-template/internal/application/pat/command"
	patQuery "github.com/lwmacct/251117-go-ddd-template/internal/application/pat/query"
	roleCommand "github.com/lwmacct/251117-go-ddd-template/internal/application/role/command"
	roleQuery "github.com/lwmacct/251117-go-ddd-template/internal/application/role/query"
	settingCommand "github.com/lwmacct/251117-go-ddd-template/internal/application/setting/command"
	settingQuery "github.com/lwmacct/251117-go-ddd-template/internal/application/setting/query"
	userCommand "github.com/lwmacct/251117-go-ddd-template/internal/application/user/command"
	userQuery "github.com/lwmacct/251117-go-ddd-template/internal/application/user/query"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/auditlog"
	domainAuth "github.com/lwmacct/251117-go-ddd-template/internal/domain/auth"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/captcha"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/menu"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/pat"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/role"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/twofa"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
	infraauth "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/auth"
	infracaptcha "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/captcha"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/config"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/database"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/persistence"
	redisinfra "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/redis"
	infratwofa "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/twofa"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// ContainerOptions 容器初始化选项
type ContainerOptions struct {
	AutoMigrate bool // 是否自动执行数据库迁移 (仅开发环境推荐)
}

// DefaultOptions 返回默认选项
func DefaultOptions() *ContainerOptions {
	return &ContainerOptions{
		AutoMigrate: false, // 生产环境默认不自动迁移
	}
}

// Container DDD+CQRS 架构的依赖注入容器
type Container struct {
	Config      *config.Config
	DB          *gorm.DB
	RedisClient *redis.Client

	// CQRS Repositories
	UserCommandRepo       user.CommandRepository
	UserQueryRepo         user.QueryRepository
	AuditLogCommandRepo   auditlog.CommandRepository
	AuditLogQueryRepo     auditlog.QueryRepository
	RoleCommandRepo       role.CommandRepository
	RoleQueryRepo         role.QueryRepository
	PermissionCommandRepo role.PermissionCommandRepository
	PermissionQueryRepo   role.PermissionQueryRepository
	PATCommandRepo        pat.CommandRepository
	PATQueryRepo          pat.QueryRepository
	MenuCommandRepo       menu.CommandRepository
	MenuQueryRepo         menu.QueryRepository
	SettingCommandRepo    setting.CommandRepository
	SettingQueryRepo      setting.QueryRepository
	TwoFACommandRepo      twofa.CommandRepository
	TwoFAQueryRepo        twofa.QueryRepository

	// Special Repository (内存存储，不需要 CQRS 分离)
	CaptchaRepository captcha.Repository

	// Domain Services
	AuthService domainAuth.Service

	// Infrastructure Services
	JWTManager          *infraauth.JWTManager
	TokenGenerator      *infraauth.TokenGenerator
	LoginSessionService *infraauth.LoginSessionService
	PATService          *infraauth.PATService
	CaptchaService      *infracaptcha.Service
	TwoFAService        *infratwofa.Service

	// Use Case Handlers - Auth
	LoginHandler        *authCommand.LoginHandler
	RegisterHandler     *authCommand.RegisterHandler
	RefreshTokenHandler *authCommand.RefreshTokenHandler

	// Use Case Handlers - User
	CreateUserHandler *userCommand.CreateUserHandler
	UpdateUserHandler *userCommand.UpdateUserHandler
	DeleteUserHandler *userCommand.DeleteUserHandler
	GetUserHandler    *userQuery.GetUserHandler
	ListUsersHandler  *userQuery.ListUsersHandler

	// HTTP Handlers
	AuthHandler *handler.AuthHandler
	UserHandler *handler.UserHandler

	Router *gin.Engine
}

// NewContainerNew 创建并初始化新架构的依赖注入容器
func NewContainer(cfg *config.Config, opts *ContainerOptions) (*Container, error) {
	if opts == nil {
		opts = DefaultOptions()
	}

	ctx := context.Background()

	// 1. 初始化数据库连接
	dbConfig := database.DefaultConfig(cfg.Data.PgsqlURL)
	db, err := database.NewConnection(ctx, dbConfig)
	if err != nil {
		return nil, err
	}

	// 2. 条件性执行自动迁移
	if opts.AutoMigrate {
		slog.Info("Auto-migration enabled, migrating database...")
		migrator := database.NewMigrator(db)
		if err := migrator.AutoMigrate(GetAllModels()...); err != nil {
			return nil, err
		}
		slog.Info("Database migration completed")
	} else {
		slog.Info("Auto-migration disabled, skipping database migration")
	}

	// 3. 初始化 Redis 客户端
	redisClient, err := redisinfra.NewClient(ctx, cfg.Data.RedisURL)
	if err != nil {
		return nil, err
	}

	// =================================================================
	// 4. 初始化 CQRS Repositories（完全符合 DDD+CQRS 架构）
	// =================================================================
	userCommandRepo := persistence.NewUserCommandRepository(db)
	userQueryRepo := persistence.NewUserQueryRepository(db)
	auditLogCommandRepo := persistence.NewAuditLogCommandRepository(db)
	auditLogQueryRepo := persistence.NewAuditLogQueryRepository(db)
	roleCommandRepo := persistence.NewRoleCommandRepository(db)
	roleQueryRepo := persistence.NewRoleQueryRepository(db)
	permissionCommandRepo := persistence.NewPermissionCommandRepository(db)
	permissionQueryRepo := persistence.NewPermissionQueryRepository(db)
	patCommandRepo := persistence.NewPATCommandRepository(db)
	patQueryRepo := persistence.NewPATQueryRepository(db)
	menuCommandRepo := persistence.NewMenuCommandRepository(db)
	menuQueryRepo := persistence.NewMenuQueryRepository(db)
	settingCommandRepo := persistence.NewSettingCommandRepository(db)
	settingQueryRepo := persistence.NewSettingQueryRepository(db)
	twofaCommandRepo := persistence.NewTwoFACommandRepository(db)
	twofaQueryRepo := persistence.NewTwoFAQueryRepository(db)

	// Special Repository (内存存储，不需要 CQRS 分离)
	captchaRepo := persistence.NewCaptchaMemoryRepository()

	// =================================================================
	// 5. 初始化 Infrastructure 组件（技术实现）
	// =================================================================
	jwtManager := infraauth.NewJWTManager(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenExpiry,
		cfg.JWT.RefreshTokenExpiry,
	)
	tokenGenerator := infraauth.NewTokenGenerator()
	loginSessionService := infraauth.NewLoginSessionService()

	// =================================================================
	// 6. 初始化 Domain Services（领域服务）
	// =================================================================
	passwordPolicy := domainAuth.DefaultPasswordPolicy()
	authService := infraauth.NewAuthService(jwtManager, tokenGenerator, passwordPolicy)

	// =================================================================
	// 7. 初始化 Use Case Handlers - Auth
	// =================================================================
	loginHandler := authCommand.NewLoginHandler(
		userQueryRepo,
		captchaRepo,
		twofaQueryRepo,
		authService,
	)

	registerHandler := authCommand.NewRegisterHandler(
		userCommandRepo,
		userQueryRepo,
		authService,
	)

	refreshTokenHandler := authCommand.NewRefreshTokenHandler(
		userQueryRepo,
		authService,
	)

	// =================================================================
	// 8. 初始化 Use Case Handlers - User
	// =================================================================
	createUserHandler := userCommand.NewCreateUserHandler(
		userCommandRepo,
		userQueryRepo,
		authService,
	)

	updateUserHandler := userCommand.NewUpdateUserHandler(
		userCommandRepo,
		userQueryRepo,
	)

	deleteUserHandler := userCommand.NewDeleteUserHandler(
		userCommandRepo,
		userQueryRepo,
	)

	getUserHandler := userQuery.NewGetUserHandler(userQueryRepo)
	listUsersHandler := userQuery.NewListUsersHandler(userQueryRepo)

	// =================================================================
	// 8.5. 初始化 Use Case Handlers - Role
	// =================================================================
	createRoleHandler := roleCommand.NewCreateRoleHandler(
		roleCommandRepo,
		roleQueryRepo,
	)

	updateRoleHandler := roleCommand.NewUpdateRoleHandler(
		roleCommandRepo,
		roleQueryRepo,
	)

	deleteRoleHandler := roleCommand.NewDeleteRoleHandler(
		roleCommandRepo,
		roleQueryRepo,
	)

	setPermissionsHandler := roleCommand.NewSetPermissionsHandler(
		roleCommandRepo,
		roleQueryRepo,
		permissionQueryRepo,
	)

	getRoleHandler := roleQuery.NewGetRoleHandler(roleQueryRepo)
	listRolesHandler := roleQuery.NewListRolesHandler(roleQueryRepo)
	listPermissionsHandler := roleQuery.NewListPermissionsHandler(permissionQueryRepo)

	// =================================================================
	// 8.6. 初始化 Use Case Handlers - Menu
	// =================================================================
	createMenuHandler := menuCommand.NewCreateMenuHandler(
		menuCommandRepo,
		menuQueryRepo,
	)

	updateMenuHandler := menuCommand.NewUpdateMenuHandler(
		menuCommandRepo,
		menuQueryRepo,
	)

	deleteMenuHandler := menuCommand.NewDeleteMenuHandler(
		menuCommandRepo,
		menuQueryRepo,
	)

	reorderMenusHandler := menuCommand.NewReorderMenusHandler(
		menuCommandRepo,
		menuQueryRepo,
	)

	getMenuHandler := menuQuery.NewGetMenuHandler(menuQueryRepo)
	listMenusHandler := menuQuery.NewListMenusHandler(menuQueryRepo)

	// =================================================================
	// 8.7. 初始化 Use Case Handlers - Setting
	// =================================================================
	createSettingHandler := settingCommand.NewCreateSettingHandler(
		settingCommandRepo,
		settingQueryRepo,
	)

	updateSettingHandler := settingCommand.NewUpdateSettingHandler(
		settingCommandRepo,
		settingQueryRepo,
	)

	deleteSettingHandler := settingCommand.NewDeleteSettingHandler(
		settingCommandRepo,
		settingQueryRepo,
	)

	batchUpdateSettingsHandler := settingCommand.NewBatchUpdateSettingsHandler(
		settingCommandRepo,
		settingQueryRepo,
	)

	getSettingHandler := settingQuery.NewGetSettingHandler(settingQueryRepo)
	listSettingsHandler := settingQuery.NewListSettingsHandler(settingQueryRepo)

	// =================================================================
	// 8.8. 初始化 Use Case Handlers - PAT
	// =================================================================
	createTokenHandler := patCommand.NewCreateTokenHandler(
		patCommandRepo,
		patQueryRepo,
		tokenGenerator,
	)

	revokeTokenHandler := patCommand.NewRevokeTokenHandler(
		patCommandRepo,
		patQueryRepo,
	)

	getTokenHandler := patQuery.NewGetTokenHandler(patQueryRepo)
	listTokensHandler := patQuery.NewListTokensHandler(patQueryRepo)

	// =================================================================
	// 9. 初始化 HTTP Handlers（适配器层）
	// =================================================================
	authHandler := handler.NewAuthHandler(
		loginHandler,
		registerHandler,
		refreshTokenHandler,
		getUserHandler,
	)

	userHandler := handler.NewUserHandler(
		createUserHandler,
		updateUserHandler,
		deleteUserHandler,
		getUserHandler,
		listUsersHandler,
	)

	roleHandler := handler.NewRoleHandler(
		createRoleHandler,
		updateRoleHandler,
		deleteRoleHandler,
		setPermissionsHandler,
		getRoleHandler,
		listRolesHandler,
		listPermissionsHandler,
	)

	menuHandler := handler.NewMenuHandler(
		createMenuHandler,
		updateMenuHandler,
		deleteMenuHandler,
		reorderMenusHandler,
		getMenuHandler,
		listMenusHandler,
	)

	settingHandler := handler.NewSettingHandler(
		createSettingHandler,
		updateSettingHandler,
		deleteSettingHandler,
		batchUpdateSettingsHandler,
		getSettingHandler,
		listSettingsHandler,
	)

	patHandler := handler.NewPATHandler(
		createTokenHandler,
		revokeTokenHandler,
		getTokenHandler,
		listTokensHandler,
	)

	// =================================================================
	// 10. 初始化 Infrastructure Services（基础设施服务）
	// =================================================================
	patService := infraauth.NewPATService(patCommandRepo, patQueryRepo, userQueryRepo, tokenGenerator)
	captchaService := infracaptcha.NewService()
	twofaService := infratwofa.NewService(twofaCommandRepo, twofaQueryRepo, userQueryRepo, cfg.Auth.TwoFAIssuer)

	// =================================================================
	// 11. 初始化路由（使用 DDD+CQRS 架构）
	// =================================================================
	authServiceForRouter := infraauth.NewService(
		userCommandRepo,
		userQueryRepo,
		twofaCommandRepo,
		twofaQueryRepo,
		captchaRepo,
		jwtManager,
		loginSessionService,
	)

	router := http.SetupRouter(
		cfg,
		db,
		redisClient,
		userCommandRepo,
		userQueryRepo,
		auditLogCommandRepo,
		auditLogQueryRepo,
		captchaRepo,
		jwtManager,
		patService,
		authServiceForRouter,
		captchaService,
		twofaService,
		authHandler,     // 使用新的 DDD+CQRS AuthHandler
		roleHandler,     // 使用新的 DDD+CQRS RoleHandler
		menuHandler,     // 使用新的 DDD+CQRS MenuHandler
		settingHandler,  // 使用新的 DDD+CQRS SettingHandler
		patHandler,      // 使用新的 DDD+CQRS PATHandler
	)

	return &Container{
		Config:      cfg,
		DB:          db,
		RedisClient: redisClient,

		// CQRS Repositories
		UserCommandRepo:       userCommandRepo,
		UserQueryRepo:         userQueryRepo,
		AuditLogCommandRepo:   auditLogCommandRepo,
		AuditLogQueryRepo:     auditLogQueryRepo,
		RoleCommandRepo:       roleCommandRepo,
		RoleQueryRepo:         roleQueryRepo,
		PermissionCommandRepo: permissionCommandRepo,
		PermissionQueryRepo:   permissionQueryRepo,
		PATCommandRepo:        patCommandRepo,
		PATQueryRepo:          patQueryRepo,
		MenuCommandRepo:       menuCommandRepo,
		MenuQueryRepo:         menuQueryRepo,
		SettingCommandRepo:    settingCommandRepo,
		SettingQueryRepo:      settingQueryRepo,
		TwoFACommandRepo:      twofaCommandRepo,
		TwoFAQueryRepo:        twofaQueryRepo,

		// Special Repository
		CaptchaRepository: captchaRepo,

		// Domain Services
		AuthService: authService,

		// Infrastructure Services
		JWTManager:          jwtManager,
		TokenGenerator:      tokenGenerator,
		LoginSessionService: loginSessionService,
		PATService:          patService,
		CaptchaService:      captchaService,
		TwoFAService:        twofaService,

		// Use Case Handlers - Auth
		LoginHandler:        loginHandler,
		RegisterHandler:     registerHandler,
		RefreshTokenHandler: refreshTokenHandler,

		// Use Case Handlers - User
		CreateUserHandler: createUserHandler,
		UpdateUserHandler: updateUserHandler,
		DeleteUserHandler: deleteUserHandler,
		GetUserHandler:    getUserHandler,
		ListUsersHandler:  listUsersHandler,

		// HTTP Handlers
		AuthHandler: authHandler,
		UserHandler: userHandler,

		Router: router,
	}, nil
}

// Close 关闭容器中的所有资源
func (c *Container) Close() error {
	// 关闭 Redis 连接
	if err := redisinfra.Close(c.RedisClient); err != nil {
		return err
	}

	// 关闭数据库连接
	if err := database.Close(c.DB); err != nil {
		return err
	}

	return nil
}

// GetAllModels 返回所有需要迁移的领域模型
// 当添加新的领域模型时，需要在这里注册
func GetAllModels() []any {
	return []any{
		&user.User{},
		&role.Role{},
		&role.Permission{},
		&auditlog.AuditLog{},
		&pat.PersonalAccessToken{},
		&twofa.TwoFA{},     // 2FA 配置表
		&menu.Menu{},       // 菜单表
		&setting.Setting{}, // 系统配置表
	}
}
