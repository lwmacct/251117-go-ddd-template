package bootstrap

import (
	"context"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/auditlog"
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

// Container 依赖注入容器
type Container struct {
	Config               *config.Config
	DB                   *gorm.DB
	RedisClient          *redis.Client
	UserRepository       user.Repository
	RoleRepository       role.RoleRepository
	PermissionRepository role.PermissionRepository
	AuditLogRepository   auditlog.Repository
	PATRepository        pat.Repository
	CaptchaRepository    captcha.Repository
	TwoFARepository      twofa.Repository
	MenuRepository       menu.Repository
	SettingRepository    setting.Repository
	JWTManager           *infraauth.JWTManager
	TokenGenerator       *infraauth.TokenGenerator
	LoginSessionService  *infraauth.LoginSessionService
	PATService           *infraauth.PATService
	AuthService          *infraauth.Service
	CaptchaService       *infracaptcha.Service
	TwoFAService         *infratwofa.Service
	Router               *gin.Engine
}

// NewContainer 创建并初始化依赖注入容器
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

	// 4. 初始化仓储
	userRepo := persistence.NewUserRepository(db)
	roleRepo := persistence.NewRoleRepository(db)
	permissionRepo := persistence.NewPermissionRepository(db)
	auditLogRepo := persistence.NewAuditLogRepository(db)
	patRepo := persistence.NewPATRepository(db)
	captchaRepo := persistence.NewCaptchaMemoryRepository()
	twofaRepo := persistence.NewTwoFARepository(db)
	menuRepo := persistence.NewMenuRepository(db)
	settingRepo := persistence.NewSettingRepository(db)

	// 5. 初始化 JWT 管理器和 Token 生成器
	jwtManager := infraauth.NewJWTManager(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenExpiry,
		cfg.JWT.RefreshTokenExpiry,
	)
	tokenGenerator := infraauth.NewTokenGenerator()
	loginSessionService := infraauth.NewLoginSessionService()

	// 6. 初始化 PAT 服务
	patService := infraauth.NewPATService(patRepo, userRepo, tokenGenerator)

	// 7. 初始化认证服务（集成验证码和2FA）
	authService := infraauth.NewService(
		userRepo,
		twofaRepo,
		captchaRepo,
		jwtManager,
		loginSessionService,
	)

	// 8. 初始化验证码和2FA服务
	captchaService := infracaptcha.NewService()
	twofaService := infratwofa.NewService(twofaRepo, userRepo, cfg.Auth.TwoFAIssuer)

	// 9. 初始化路由 (传入依赖)
	router := http.SetupRouter(
		cfg,
		db,
		redisClient,
		userRepo,
		roleRepo,
		permissionRepo,
		auditLogRepo,
		patRepo,
		captchaRepo,
		menuRepo,
		settingRepo,
		jwtManager,
		tokenGenerator,
		patService,
		authService,
		captchaService,
		twofaService,
	)

	return &Container{
		Config:               cfg,
		DB:                   db,
		RedisClient:          redisClient,
		UserRepository:       userRepo,
		RoleRepository:       roleRepo,
		PermissionRepository: permissionRepo,
		AuditLogRepository:   auditLogRepo,
		PATRepository:        patRepo,
		CaptchaRepository:    captchaRepo,
		TwoFARepository:      twofaRepo,
		MenuRepository:       menuRepo,
		SettingRepository:    settingRepo,
		JWTManager:           jwtManager,
		TokenGenerator:       tokenGenerator,
		LoginSessionService:  loginSessionService,
		PATService:           patService,
		AuthService:          authService,
		CaptchaService:       captchaService,
		TwoFAService:         twofaService,
		Router:               router,
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
		&twofa.TwoFA{},   // 2FA 配置表
		&menu.Menu{},     // 菜单表
		&setting.Setting{}, // 系统配置表
	}
}
