package bootstrap

import (
	"context"
	"log/slog"

	// Standard library imports
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	// Adapter imports
	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http"
	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/handler"

	// Application imports - Use Case Handlers
	auditlogQuery "github.com/lwmacct/251117-go-ddd-template/internal/application/auditlog/query"
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

	// Domain imports
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/auth"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/captcha"

	// Infrastructure imports, 统一使用前缀 _

	_auth "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/auth"
	_captcha "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/captcha"
	_config "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/config"
	_database "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/database"
	_persistence "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/persistence"
	_redis "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/redis"
	_twofa "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/twofa"
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
	Config      *_config.Config
	DB          *gorm.DB
	RedisClient *redis.Client

	// CQRS Repositories
	UserRepos       _persistence.UserRepositories
	AuditLogRepos   _persistence.AuditLogRepositories
	RoleRepos       _persistence.RoleRepositories
	PermissionRepos _persistence.PermissionRepositories
	PATRepos        _persistence.PATRepositories
	MenuRepos       _persistence.MenuRepositories
	SettingRepos    _persistence.SettingRepositories
	TwoFARepos      _persistence.TwoFARepositories

	// Captcha（内存实现，同样遵循 CQRS）
	CaptchaCommandRepo captcha.CommandRepository
	CaptchaQueryRepo   captcha.QueryRepository

	// Domain Services
	AuthService auth.Service

	// Infrastructure Services
	JWTManager          *_auth.JWTManager
	TokenGenerator      *_auth.TokenGenerator
	LoginSessionService *_auth.LoginSessionService
	PATService          *_auth.PATService
	CaptchaService      *_captcha.Service
	TwoFAService        *_twofa.Service

	// Use Case Handlers - Auth
	LoginHandler        *authCommand.LoginHandler
	RegisterHandler     *authCommand.RegisterHandler
	RefreshTokenHandler *authCommand.RefreshTokenHandler

	// Use Case Handlers - User
	CreateUserHandler     *userCommand.CreateUserHandler
	UpdateUserHandler     *userCommand.UpdateUserHandler
	DeleteUserHandler     *userCommand.DeleteUserHandler
	AssignRolesHandler    *userCommand.AssignRolesHandler
	ChangePasswordHandler *userCommand.ChangePasswordHandler
	GetUserHandler        *userQuery.GetUserHandler
	ListUsersHandler      *userQuery.ListUsersHandler

	// HTTP Handlers
	AuthHandler        *handler.AuthHandler
	UserHandler        *handler.UserHandler
	AdminUserHandler   *handler.AdminUserHandler
	UserProfileHandler *handler.UserProfileHandler

	Router *gin.Engine
}

// NewContainerNew 创建并初始化新架构的依赖注入容器
func NewContainer(cfg *_config.Config, opts *ContainerOptions) (*Container, error) {
	if opts == nil {
		opts = DefaultOptions()
	}

	ctx := context.Background()

	// 1. 初始化数据库连接
	dbConfig := _database.DefaultConfig(cfg.Data.PgsqlURL)
	db, err := _database.NewConnection(ctx, dbConfig)
	if err != nil {
		return nil, err
	}

	// 2. 条件性执行自动迁移
	if opts.AutoMigrate {
		slog.Info("Auto-migration enabled, migrating database...")
		migrator := _database.NewMigrator(db)
		if err := migrator.AutoMigrate(GetAllModels()...); err != nil {
			return nil, err
		}
		slog.Info("Database migration completed")
	} else {
		slog.Info("Auto-migration disabled, skipping database migration")
	}

	// 3. 初始化 Redis 客户端
	redisClient, err := _redis.NewClient(ctx, cfg.Data.RedisURL)
	if err != nil {
		return nil, err
	}

	// =================================================================
	// 4. 初始化 CQRS Repositories（完全符合 DDD+CQRS 架构）
	// =================================================================
	userRepos := _persistence.NewUserRepositories(db)
	auditLogRepos := _persistence.NewAuditLogRepositories(db)
	roleRepos := _persistence.NewRoleRepositories(db)
	permissionRepos := _persistence.NewPermissionRepositories(db)
	patRepos := _persistence.NewPATRepositories(db)
	menuRepos := _persistence.NewMenuRepositories(db)
	settingRepos := _persistence.NewSettingRepositories(db)
	twofaRepos := _persistence.NewTwoFARepositories(db)

	// Captcha Repository（内存实现，同样提供 Command/Query 能力）
	captchaRepo := _persistence.NewCaptchaMemoryRepository()

	// =================================================================
	// 5. 初始化 Infrastructure 组件（技术实现）
	// =================================================================
	jwtManager := _auth.NewJWTManager(cfg.JWT.Secret, cfg.JWT.AccessTokenExpiry, cfg.JWT.RefreshTokenExpiry)
	tokenGenerator := _auth.NewTokenGenerator()
	loginSessionService := _auth.NewLoginSessionService()

	// =================================================================
	// 6. 初始化 Domain Services（领域服务）
	// =================================================================
	passwordPolicy := auth.DefaultPasswordPolicy()
	authService := _auth.NewAuthService(jwtManager, tokenGenerator, passwordPolicy)

	// =================================================================
	// 7. 初始化 Use Case Handlers - Auth
	// =================================================================
	loginHandler := authCommand.NewLoginHandler(userRepos.Query, captchaRepo, twofaRepos.Query, authService)
	registerHandler := authCommand.NewRegisterHandler(userRepos.Command, userRepos.Query, authService)
	refreshTokenHandler := authCommand.NewRefreshTokenHandler(userRepos.Query, authService)

	// =================================================================
	// 8. 初始化 Use Case Handlers - User
	// =================================================================
	createUserHandler := userCommand.NewCreateUserHandler(userRepos.Command, userRepos.Query, authService)
	updateUserHandler := userCommand.NewUpdateUserHandler(userRepos.Command, userRepos.Query)
	deleteUserHandler := userCommand.NewDeleteUserHandler(userRepos.Command, userRepos.Query)
	assignRolesHandler := userCommand.NewAssignRolesHandler(userRepos.Command, userRepos.Query)
	changePasswordHandler := userCommand.NewChangePasswordHandler(userRepos.Command, userRepos.Query, authService)
	getUserHandler := userQuery.NewGetUserHandler(userRepos.Query)
	listUsersHandler := userQuery.NewListUsersHandler(userRepos.Query)

	// =================================================================
	// 8.5. 初始化 Use Case Handlers - Role
	// =================================================================
	createRoleHandler := roleCommand.NewCreateRoleHandler(roleRepos.Command, roleRepos.Query)
	updateRoleHandler := roleCommand.NewUpdateRoleHandler(roleRepos.Command, roleRepos.Query)
	deleteRoleHandler := roleCommand.NewDeleteRoleHandler(roleRepos.Command, roleRepos.Query)
	setPermissionsHandler := roleCommand.NewSetPermissionsHandler(roleRepos.Command, roleRepos.Query, permissionRepos.Query)
	getRoleHandler := roleQuery.NewGetRoleHandler(roleRepos.Query)
	listRolesHandler := roleQuery.NewListRolesHandler(roleRepos.Query)
	listPermissionsHandler := roleQuery.NewListPermissionsHandler(permissionRepos.Query)

	// =================================================================
	// 8.6. 初始化 Use Case Handlers - Menu
	// =================================================================
	createMenuHandler := menuCommand.NewCreateMenuHandler(menuRepos.Command, menuRepos.Query)
	updateMenuHandler := menuCommand.NewUpdateMenuHandler(menuRepos.Command, menuRepos.Query)
	deleteMenuHandler := menuCommand.NewDeleteMenuHandler(menuRepos.Command, menuRepos.Query)
	reorderMenusHandler := menuCommand.NewReorderMenusHandler(menuRepos.Command, menuRepos.Query)
	getMenuHandler := menuQuery.NewGetMenuHandler(menuRepos.Query)
	listMenusHandler := menuQuery.NewListMenusHandler(menuRepos.Query)

	// =================================================================
	// 8.7. 初始化 Use Case Handlers - Setting
	// =================================================================
	createSettingHandler := settingCommand.NewCreateSettingHandler(settingRepos.Command, settingRepos.Query)
	updateSettingHandler := settingCommand.NewUpdateSettingHandler(settingRepos.Command, settingRepos.Query)
	deleteSettingHandler := settingCommand.NewDeleteSettingHandler(settingRepos.Command, settingRepos.Query)
	batchUpdateSettingsHandler := settingCommand.NewBatchUpdateSettingsHandler(settingRepos.Command, settingRepos.Query)
	getSettingHandler := settingQuery.NewGetSettingHandler(settingRepos.Query)
	listSettingsHandler := settingQuery.NewListSettingsHandler(settingRepos.Query)

	// =================================================================
	// 8.8. 初始化 Use Case Handlers - PAT
	// =================================================================
	createTokenHandler := patCommand.NewCreateTokenHandler(patRepos.Command, patRepos.Query, tokenGenerator)
	revokeTokenHandler := patCommand.NewRevokeTokenHandler(patRepos.Command, patRepos.Query)
	getTokenHandler := patQuery.NewGetTokenHandler(patRepos.Query)
	listTokensHandler := patQuery.NewListTokensHandler(patRepos.Query)

	// =================================================================
	// 8.9. 初始化 Use Case Handlers - AuditLog
	// =================================================================
	listLogsHandler := auditlogQuery.NewListLogsHandler(auditLogRepos.Query)
	getLogHandler := auditlogQuery.NewGetLogHandler(auditLogRepos.Query)

	// =================================================================
	// 9. 初始化 HTTP Handlers（适配器层）
	// =================================================================
	authHandler := handler.NewAuthHandler(loginHandler, registerHandler, refreshTokenHandler, getUserHandler)
	userHandler := handler.NewUserHandler(createUserHandler, updateUserHandler, deleteUserHandler, getUserHandler, listUsersHandler)
	adminUserHandler := handler.NewAdminUserHandler(createUserHandler, updateUserHandler, deleteUserHandler, assignRolesHandler, getUserHandler, listUsersHandler)
	userProfileHandler := handler.NewUserProfileHandler(getUserHandler, updateUserHandler, changePasswordHandler, deleteUserHandler)
	roleHandler := handler.NewRoleHandler(createRoleHandler, updateRoleHandler, deleteRoleHandler, setPermissionsHandler, getRoleHandler, listRolesHandler, listPermissionsHandler)
	menuHandler := handler.NewMenuHandler(createMenuHandler, updateMenuHandler, deleteMenuHandler, reorderMenusHandler, getMenuHandler, listMenusHandler)
	settingHandler := handler.NewSettingHandler(createSettingHandler, updateSettingHandler, deleteSettingHandler, batchUpdateSettingsHandler, getSettingHandler, listSettingsHandler)
	patHandler := handler.NewPATHandler(createTokenHandler, revokeTokenHandler, getTokenHandler, listTokensHandler)
	auditLogHandler := handler.NewAuditLogHandler(listLogsHandler, getLogHandler)

	// =================================================================
	// 10. 初始化 Infrastructure Services（基础设施服务）
	// =================================================================
	patService := _auth.NewPATService(patRepos.Command, patRepos.Query, userRepos.Query, tokenGenerator)
	captchaService := _captcha.NewService()
	twofaService := _twofa.NewService(twofaRepos.Command, twofaRepos.Query, userRepos.Query, cfg.Auth.TwoFAIssuer)

	// =================================================================
	// 11. 初始化路由（使用 DDD+CQRS 架构）
	// =================================================================
	authServiceForRouter := _auth.NewService(userRepos.Command, userRepos.Query, twofaRepos.Command, twofaRepos.Query, captchaRepo, jwtManager, loginSessionService)

	router := http.SetupRouter(
		cfg,
		db,
		redisClient,
		userRepos.Query,
		auditLogRepos.Command, // 审计中间件需要
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
		auditLogHandler, // 使用新的 DDD+CQRS AuditLogHandler
		adminUserHandler,
		userProfileHandler,
	)

	return &Container{
		Config:      cfg,
		DB:          db,
		RedisClient: redisClient,

		// CQRS Repositories
		UserRepos:       userRepos,
		AuditLogRepos:   auditLogRepos,
		RoleRepos:       roleRepos,
		PermissionRepos: permissionRepos,
		PATRepos:        patRepos,
		MenuRepos:       menuRepos,
		SettingRepos:    settingRepos,
		TwoFARepos:      twofaRepos,

		// Captcha Repository（单实例实现 Command/Query）
		CaptchaCommandRepo: captchaRepo,
		CaptchaQueryRepo:   captchaRepo,

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
		CreateUserHandler:     createUserHandler,
		UpdateUserHandler:     updateUserHandler,
		DeleteUserHandler:     deleteUserHandler,
		AssignRolesHandler:    assignRolesHandler,
		ChangePasswordHandler: changePasswordHandler,
		GetUserHandler:        getUserHandler,
		ListUsersHandler:      listUsersHandler,

		// HTTP Handlers
		AuthHandler:        authHandler,
		UserHandler:        userHandler,
		AdminUserHandler:   adminUserHandler,
		UserProfileHandler: userProfileHandler,

		Router: router,
	}, nil
}

// Close 关闭容器中的所有资源
func (c *Container) Close() error {
	// 关闭 Redis 连接
	if err := _redis.Close(c.RedisClient); err != nil {
		return err
	}

	// 关闭数据库连接
	if err := _database.Close(c.DB); err != nil {
		return err
	}

	return nil
}

// GetAllModels 返回所有需要迁移的领域模型
// 当添加新的领域模型时，需要在这里注册
func GetAllModels() []any {
	return []any{
		&_persistence.UserModel{},
		&_persistence.RoleModel{},
		&_persistence.PermissionModel{},
		&_persistence.PersonalAccessTokenModel{},
		&_persistence.AuditLogModel{},
		&_persistence.TwoFAModel{},
		&_persistence.MenuModel{},
		&_persistence.SettingModel{},
	}
}
