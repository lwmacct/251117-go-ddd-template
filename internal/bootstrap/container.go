// Package bootstrap 提供应用程序的依赖注入容器。
//
// 本包是 DDD+CQRS 架构的核心组装点，负责：
//   - 初始化所有基础设施组件（数据库、Redis、JWT）
//   - 创建 CQRS Repositories（Command/Query 分离）
//   - 实例化 Use Case Handlers（业务逻辑编排）
//   - 装配 HTTP Handlers（适配器层）
//   - 配置路由和中间件
//
// 依赖注入顺序：
//  1. 基础设施 → 2. Repositories → 3. Domain Services
//  4. Use Case Handlers → 5. HTTP Handlers → 6. Router
//
// 使用方式：
//
//	container, err := bootstrap.NewContainer(cfg, nil)
//	if err != nil { ... }
//	defer container.Close()
//	container.Router.Run(":8080")
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
	"github.com/lwmacct/251117-go-ddd-template/internal/config"

	// Application imports - Use Case Handlers
	auditlogQuery "github.com/lwmacct/251117-go-ddd-template/internal/application/auditlog/query"
	authCommand "github.com/lwmacct/251117-go-ddd-template/internal/application/auth/command"
	captchaCommand "github.com/lwmacct/251117-go-ddd-template/internal/application/captcha/command"
	menuCommand "github.com/lwmacct/251117-go-ddd-template/internal/application/menu/command"
	menuQuery "github.com/lwmacct/251117-go-ddd-template/internal/application/menu/query"
	patCommand "github.com/lwmacct/251117-go-ddd-template/internal/application/pat/command"
	patQuery "github.com/lwmacct/251117-go-ddd-template/internal/application/pat/query"
	roleCommand "github.com/lwmacct/251117-go-ddd-template/internal/application/role/command"
	roleQuery "github.com/lwmacct/251117-go-ddd-template/internal/application/role/query"
	settingCommand "github.com/lwmacct/251117-go-ddd-template/internal/application/setting/command"
	settingQuery "github.com/lwmacct/251117-go-ddd-template/internal/application/setting/query"
	statsQuery "github.com/lwmacct/251117-go-ddd-template/internal/application/stats/query"
	userCommand "github.com/lwmacct/251117-go-ddd-template/internal/application/user/command"
	userQuery "github.com/lwmacct/251117-go-ddd-template/internal/application/user/query"

	// Domain imports
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/auth"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/captcha"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/stats"

	// Infrastructure imports, 统一使用前缀 _

	_auth "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/auth"
	_captcha "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/captcha"
	_database "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/database"
	_health "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/health"
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
	Config      *config.Config
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

	// Stats（只读统计）
	StatsQueryRepo stats.QueryRepository

	// Domain Services
	AuthService auth.Service

	// Infrastructure Services
	JWTManager             *_auth.JWTManager
	TokenGenerator         auth.TokenGenerator // Domain interface, implemented by _auth.TokenGenerator
	LoginSessionService    *_auth.LoginSessionService
	PermissionCacheService *_auth.PermissionCacheService
	PATService             *_auth.PATService
	CaptchaService         *_captcha.Service
	TwoFAService           *_twofa.Service

	// Use Case Handlers - Auth
	LoginHandler        *authCommand.LoginHandler
	RegisterHandler     *authCommand.RegisterHandler
	RefreshTokenHandler *authCommand.RefreshTokenHandler

	// Use Case Handlers - User
	CreateUserHandler      *userCommand.CreateUserHandler
	UpdateUserHandler      *userCommand.UpdateUserHandler
	DeleteUserHandler      *userCommand.DeleteUserHandler
	AssignRolesHandler     *userCommand.AssignRolesHandler
	ChangePasswordHandler  *userCommand.ChangePasswordHandler
	BatchCreateUserHandler *userCommand.BatchCreateUsersHandler
	GetUserHandler         *userQuery.GetUserHandler
	ListUsersHandler       *userQuery.ListUsersHandler

	// HTTP Handlers
	AuthHandler        *handler.AuthHandler
	UserHandler        *handler.UserHandler
	AdminUserHandler   *handler.AdminUserHandler
	UserProfileHandler *handler.UserProfileHandler

	Router *gin.Engine

	// 内部使用的临时变量（用于分层初始化）
	captchaRepo    captchaRepository // 组合接口
	captchaService *_captcha.Service
	tokenGenerator *_auth.TokenGenerator
}

// captchaRepository 内部组合接口，同时满足 Command 和 Query
type captchaRepository interface {
	captcha.CommandRepository
	captcha.QueryRepository
}

// NewContainer 创建并初始化新架构的依赖注入容器
func NewContainer(ctx context.Context, cfg *config.Config, opts *ContainerOptions) (*Container, error) {
	if opts == nil {
		opts = DefaultOptions()
	}

	c := &Container{Config: cfg}

	// 分层初始化
	if err := c.initInfrastructure(ctx, cfg, opts); err != nil {
		return nil, err
	}

	c.initRepositories()
	c.initDomainServices(cfg)
	c.initUseCaseHandlers()
	c.initRouter(cfg)

	return c, nil
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

// initInfrastructure 初始化基础设施（数据库、Redis）
func (c *Container) initInfrastructure(ctx context.Context, cfg *config.Config, opts *ContainerOptions) error {
	// 初始化数据库连接
	dbConfig := _database.DefaultConfig(cfg.Data.PgsqlURL)
	db, err := _database.NewConnection(ctx, dbConfig)
	if err != nil {
		return err
	}
	c.DB = db

	// 条件性执行自动迁移
	if opts.AutoMigrate {
		slog.Info("Auto-migration enabled, migrating database...")
		migrator := _database.NewMigrator(db)
		if err = migrator.AutoMigrate(GetAllModels()...); err != nil {
			return err
		}
		slog.Info("Database migration completed")
	} else {
		slog.Info("Auto-migration disabled, skipping database migration")
	}

	// 初始化 Redis 客户端
	redisClient, err := _redis.NewClient(ctx, cfg.Data.RedisURL)
	if err != nil {
		return err
	}
	c.RedisClient = redisClient

	return nil
}

// initRepositories 初始化所有 CQRS Repositories
func (c *Container) initRepositories() {
	c.UserRepos = _persistence.NewUserRepositories(c.DB)
	c.AuditLogRepos = _persistence.NewAuditLogRepositories(c.DB)
	c.RoleRepos = _persistence.NewRoleRepositories(c.DB)
	c.PermissionRepos = _persistence.NewPermissionRepositories(c.DB)
	c.PATRepos = _persistence.NewPATRepositories(c.DB)
	c.MenuRepos = _persistence.NewMenuRepositories(c.DB)
	c.SettingRepos = _persistence.NewSettingRepositories(c.DB)
	c.TwoFARepos = _persistence.NewTwoFARepositories(c.DB)

	// Captcha Repository（内存实现）
	c.captchaRepo = _persistence.NewCaptchaMemoryRepository()
	c.CaptchaCommandRepo = c.captchaRepo
	c.CaptchaQueryRepo = c.captchaRepo

	// Stats Repository（只读统计）
	c.StatsQueryRepo = _persistence.NewStatsQueryRepository(c.DB)
}

// initDomainServices 初始化领域服务和基础设施服务
func (c *Container) initDomainServices(cfg *config.Config) {
	// Infrastructure 组件
	c.JWTManager = _auth.NewJWTManager(cfg.JWT.Secret, cfg.JWT.AccessTokenExpiry, cfg.JWT.RefreshTokenExpiry)
	c.tokenGenerator = _auth.NewTokenGenerator()
	c.TokenGenerator = c.tokenGenerator
	c.LoginSessionService = _auth.NewLoginSessionService()
	c.PermissionCacheService = _auth.NewPermissionCacheService(c.RedisClient, c.UserRepos.Query, cfg.Data.RedisKeyPrefix)

	// Domain Services
	passwordPolicy := auth.DefaultPasswordPolicy()
	c.AuthService = _auth.NewAuthService(c.JWTManager, c.tokenGenerator, passwordPolicy)
	c.captchaService = _captcha.NewService()
	c.CaptchaService = c.captchaService
}

// initUseCaseHandlers 初始化所有 Use Case Handlers
func (c *Container) initUseCaseHandlers() {
	c.initAuthHandlers()
	c.initUserHandlers()
}

// initAuthHandlers 初始化认证相关的 Use Case Handlers
func (c *Container) initAuthHandlers() {
	c.LoginHandler = authCommand.NewLoginHandler(c.UserRepos.Query, c.captchaRepo, c.TwoFARepos.Query, c.AuthService)
	c.RegisterHandler = authCommand.NewRegisterHandler(c.UserRepos.Command, c.UserRepos.Query, c.AuthService)
	c.RefreshTokenHandler = authCommand.NewRefreshTokenHandler(c.UserRepos.Query, c.AuthService)
}

// initUserHandlers 初始化用户相关的 Use Case Handlers
func (c *Container) initUserHandlers() {
	c.CreateUserHandler = userCommand.NewCreateUserHandler(c.UserRepos.Command, c.UserRepos.Query, c.AuthService)
	c.UpdateUserHandler = userCommand.NewUpdateUserHandler(c.UserRepos.Command, c.UserRepos.Query)
	c.DeleteUserHandler = userCommand.NewDeleteUserHandler(c.UserRepos.Command, c.UserRepos.Query)
	c.AssignRolesHandler = userCommand.NewAssignRolesHandler(c.UserRepos.Command, c.UserRepos.Query)
	c.ChangePasswordHandler = userCommand.NewChangePasswordHandler(c.UserRepos.Command, c.UserRepos.Query, c.AuthService)
	c.BatchCreateUserHandler = userCommand.NewBatchCreateUsersHandler(c.UserRepos.Command, c.UserRepos.Query, c.AuthService)
	c.GetUserHandler = userQuery.NewGetUserHandler(c.UserRepos.Query)
	c.ListUsersHandler = userQuery.NewListUsersHandler(c.UserRepos.Query)
}

// initRouter 初始化路由和 HTTP Handlers
func (c *Container) initRouter(cfg *config.Config) {
	// Role handlers
	createRoleHandler := roleCommand.NewCreateRoleHandler(c.RoleRepos.Command, c.RoleRepos.Query)
	updateRoleHandler := roleCommand.NewUpdateRoleHandler(c.RoleRepos.Command, c.RoleRepos.Query)
	deleteRoleHandler := roleCommand.NewDeleteRoleHandler(c.RoleRepos.Command, c.RoleRepos.Query)
	setPermissionsHandler := roleCommand.NewSetPermissionsHandler(c.RoleRepos.Command, c.RoleRepos.Query, c.PermissionRepos.Query)
	getRoleHandler := roleQuery.NewGetRoleHandler(c.RoleRepos.Query)
	listRolesHandler := roleQuery.NewListRolesHandler(c.RoleRepos.Query)
	listPermissionsHandler := roleQuery.NewListPermissionsHandler(c.PermissionRepos.Query)

	// Menu handlers
	createMenuHandler := menuCommand.NewCreateMenuHandler(c.MenuRepos.Command, c.MenuRepos.Query)
	updateMenuHandler := menuCommand.NewUpdateMenuHandler(c.MenuRepos.Command, c.MenuRepos.Query)
	deleteMenuHandler := menuCommand.NewDeleteMenuHandler(c.MenuRepos.Command, c.MenuRepos.Query)
	reorderMenusHandler := menuCommand.NewReorderMenusHandler(c.MenuRepos.Command, c.MenuRepos.Query)
	getMenuHandler := menuQuery.NewGetMenuHandler(c.MenuRepos.Query)
	listMenusHandler := menuQuery.NewListMenusHandler(c.MenuRepos.Query)

	// Setting handlers
	createSettingHandler := settingCommand.NewCreateSettingHandler(c.SettingRepos.Command, c.SettingRepos.Query)
	updateSettingHandler := settingCommand.NewUpdateSettingHandler(c.SettingRepos.Command, c.SettingRepos.Query)
	deleteSettingHandler := settingCommand.NewDeleteSettingHandler(c.SettingRepos.Command, c.SettingRepos.Query)
	batchUpdateSettingsHandler := settingCommand.NewBatchUpdateSettingsHandler(c.SettingRepos.Command, c.SettingRepos.Query)
	getSettingHandler := settingQuery.NewGetSettingHandler(c.SettingRepos.Query)
	listSettingsHandler := settingQuery.NewListSettingsHandler(c.SettingRepos.Query)

	// PAT handlers
	createTokenHandler := patCommand.NewCreateTokenHandler(c.PATRepos.Command, c.UserRepos.Query, c.tokenGenerator)
	deleteTokenHandler := patCommand.NewDeleteTokenHandler(c.PATRepos.Command, c.PATRepos.Query)
	disableTokenHandler := patCommand.NewDisableTokenHandler(c.PATRepos.Command, c.PATRepos.Query)
	enableTokenHandler := patCommand.NewEnableTokenHandler(c.PATRepos.Command, c.PATRepos.Query)
	getTokenHandler := patQuery.NewGetTokenHandler(c.PATRepos.Query)
	listTokensHandler := patQuery.NewListTokensHandler(c.PATRepos.Query)

	// AuditLog handlers
	listLogsHandler := auditlogQuery.NewListLogsHandler(c.AuditLogRepos.Query)
	getLogHandler := auditlogQuery.NewGetLogHandler(c.AuditLogRepos.Query)

	// Stats handlers
	getStatsHandler := statsQuery.NewGetStatsHandler(c.StatsQueryRepo)

	// Captcha handler
	generateCaptchaHandler := captchaCommand.NewGenerateCaptchaHandler(c.captchaRepo, c.captchaService)

	// Health checker
	healthChecker := _health.NewSystemChecker(c.DB, c.RedisClient)

	// HTTP Handlers for router
	roleHandler := handler.NewRoleHandler(createRoleHandler, updateRoleHandler, deleteRoleHandler, setPermissionsHandler, getRoleHandler, listRolesHandler, listPermissionsHandler)
	menuHandler := handler.NewMenuHandler(createMenuHandler, updateMenuHandler, deleteMenuHandler, reorderMenusHandler, getMenuHandler, listMenusHandler)
	settingHandler := handler.NewSettingHandler(createSettingHandler, updateSettingHandler, deleteSettingHandler, batchUpdateSettingsHandler, getSettingHandler, listSettingsHandler)
	patHandler := handler.NewPATHandler(createTokenHandler, deleteTokenHandler, disableTokenHandler, enableTokenHandler, getTokenHandler, listTokensHandler)
	auditLogHandler := handler.NewAuditLogHandler(listLogsHandler, getLogHandler)
	overviewHandler := handler.NewOverviewHandler(getStatsHandler)
	captchaHandler := handler.NewCaptchaHandler(generateCaptchaHandler, cfg.Auth.DevSecret)
	healthHandler := handler.NewHealthHandler(healthChecker)

	// Infrastructure Services
	c.PATService = _auth.NewPATService(c.PATRepos.Command, c.PATRepos.Query, c.UserRepos.Query, c.tokenGenerator)
	c.TwoFAService = _twofa.NewService(c.TwoFARepos.Command, c.TwoFARepos.Query, c.UserRepos.Query, cfg.Auth.TwoFAIssuer)

	// Auth service for router
	authServiceForRouter := _auth.NewService(c.UserRepos.Command, c.UserRepos.Query, c.TwoFARepos.Command, c.TwoFARepos.Query, c.captchaRepo, c.JWTManager, c.LoginSessionService)

	c.Router = http.SetupRouter(
		cfg,
		c.RedisClient,
		c.UserRepos.Query,
		c.AuditLogRepos.Command,
		c.JWTManager,
		c.PATService,
		c.PermissionCacheService,
		authServiceForRouter,
		c.TwoFAService,
		healthHandler,
		c.AuthHandler,
		captchaHandler,
		roleHandler,
		menuHandler,
		settingHandler,
		patHandler,
		auditLogHandler,
		c.AdminUserHandler,
		c.UserProfileHandler,
		overviewHandler,
	)
}
