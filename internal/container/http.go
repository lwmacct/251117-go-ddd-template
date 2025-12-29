package container

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"gorm.io/gorm"

	adapthttp "github.com/lwmacct/251117-go-ddd-template/internal/adapters/http"
	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/handler"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/cache"
	"github.com/lwmacct/251117-go-ddd-template/internal/config"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/auth"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/health"
)

// HandlersResult 使用 fx.Out 批量返回所有 HTTP 处理器。
type HandlersResult struct {
	fx.Out

	Health      *handler.HealthHandler
	Auth        *handler.AuthHandler
	Captcha     *handler.CaptchaHandler
	AdminUser   *handler.AdminUserHandler
	UserProfile *handler.UserProfileHandler
	Role        *handler.RoleHandler
	Menu        *handler.MenuHandler
	Setting     *handler.SettingHandler
	UserSetting *handler.UserSettingHandler
	PAT         *handler.PATHandler
	AuditLog    *handler.AuditLogHandler
	Overview    *handler.OverviewHandler
	TwoFA       *handler.TwoFAHandler
	Cache       *handler.CacheHandler
}

// HTTPModule 提供 HTTP 处理器和路由。
var HTTPModule = fx.Module("http",
	fx.Provide(
		newHealthChecker,
		newAllHandlers,
		newRouter,
	),
)

func newHealthChecker(db *gorm.DB, redisClient *redis.Client) *health.SystemChecker {
	return health.NewSystemChecker(db, redisClient)
}

// handlersParams 聚合创建 Handler 所需的依赖。
type handlersParams struct {
	fx.In

	Config        *config.Config
	RedisClient   *redis.Client
	HealthChecker *health.SystemChecker
	Auth          *AuthUseCases
	User          *UserUseCases
	Role          *RoleUseCases
	Menu          *MenuUseCases
	Setting       *SettingUseCases
	UserSetting   *UserSettingUseCases
	PAT           *PATUseCases
	AuditLog      *AuditLogUseCases
	Stats         *StatsUseCases
	Captcha       *CaptchaUseCases
	TwoFA         *TwoFAUseCases
}

func newAllHandlers(p handlersParams) HandlersResult {
	keyPrefix := p.Config.Data.RedisKeyPrefix

	return HandlersResult{
		Health: handler.NewHealthHandler(p.HealthChecker),
		Auth: handler.NewAuthHandler(
			p.Auth.Login,
			p.Auth.Login2FA,
			p.Auth.Register,
			p.Auth.RefreshToken,
		),
		Captcha: handler.NewCaptchaHandler(p.Captcha.Generate, p.Config.Auth.DevSecret),
		AdminUser: handler.NewAdminUserHandler(
			p.User.Create,
			p.User.Update,
			p.User.Delete,
			p.User.AssignRoles,
			p.User.BatchCreate,
			p.User.Get,
			p.User.List,
		),
		UserProfile: handler.NewUserProfileHandler(
			p.User.Get,
			p.User.Update,
			p.User.ChangePassword,
			p.User.Delete,
		),
		Role: handler.NewRoleHandler(
			p.Role.Create,
			p.Role.Update,
			p.Role.Delete,
			p.Role.SetPermissions,
			p.Role.Get,
			p.Role.List,
			p.Role.ListPermissions,
		),
		Menu: handler.NewMenuHandler(
			p.Menu.Create,
			p.Menu.Update,
			p.Menu.Delete,
			p.Menu.Reorder,
			p.Menu.Get,
			p.Menu.List,
		),
		Setting: handler.NewSettingHandler(
			p.Setting.Create,
			p.Setting.Update,
			p.Setting.Delete,
			p.Setting.BatchUpdate,
			p.Setting.Get,
			p.Setting.List,
			p.Setting.ListSettings,
			p.Setting.CreateCategory,
			p.Setting.UpdateCategory,
			p.Setting.DeleteCategory,
			p.Setting.GetCategory,
			p.Setting.ListCategories,
		),
		UserSetting: handler.NewUserSettingHandler(
			p.UserSetting.Set,
			p.UserSetting.BatchSet,
			p.UserSetting.Reset,
			p.UserSetting.ResetAll,
			p.UserSetting.Get,
			p.UserSetting.List,
			p.UserSetting.ListSettings,
			p.UserSetting.ListCategories,
		),
		PAT: handler.NewPATHandler(
			p.PAT.Create,
			p.PAT.Delete,
			p.PAT.Disable,
			p.PAT.Enable,
			p.PAT.Get,
			p.PAT.List,
		),
		AuditLog: handler.NewAuditLogHandler(
			p.AuditLog.List,
			p.AuditLog.Get,
		),
		Overview: handler.NewOverviewHandler(p.Stats.GetStats),
		TwoFA: handler.NewTwoFAHandler(
			p.TwoFA.Setup,
			p.TwoFA.VerifyEnable,
			p.TwoFA.Disable,
			p.TwoFA.GetStatus,
		),
		Cache: handler.NewCacheHandler(
			cache.NewInfoHandler(p.RedisClient, keyPrefix),
			cache.NewScanKeysHandler(p.RedisClient, keyPrefix),
			cache.NewGetKeyHandler(p.RedisClient, keyPrefix),
			cache.NewDeleteHandler(p.RedisClient, keyPrefix),
		),
	}
}

// routerParams 聚合创建路由所需的依赖。
type routerParams struct {
	fx.In

	Config      *config.Config
	RedisClient *redis.Client

	// Services
	JWTManager      *auth.JWTManager
	PATService      *auth.PATService
	PermissionCache *auth.PermissionCacheService

	// UseCases
	AuditLog *AuditLogUseCases

	// Handlers
	Health      *handler.HealthHandler
	Auth        *handler.AuthHandler
	Captcha     *handler.CaptchaHandler
	AdminUser   *handler.AdminUserHandler
	UserProfile *handler.UserProfileHandler
	Role        *handler.RoleHandler
	Menu        *handler.MenuHandler
	Setting     *handler.SettingHandler
	UserSetting *handler.UserSettingHandler
	PAT         *handler.PATHandler
	AuditLogH   *handler.AuditLogHandler
	Overview    *handler.OverviewHandler
	TwoFA       *handler.TwoFAHandler
	Cache       *handler.CacheHandler
}

func newRouter(p routerParams) *gin.Engine {
	deps := &adapthttp.RouterDependencies{
		Config:                 p.Config,
		RedisClient:            p.RedisClient,
		CreateLogHandler:       p.AuditLog.CreateLog,
		JWTManager:             p.JWTManager,
		PATService:             p.PATService,
		PermissionCacheService: p.PermissionCache,
		HealthHandler:          p.Health,
		AuthHandler:            p.Auth,
		CaptchaHandler:         p.Captcha,
		RoleHandler:            p.Role,
		MenuHandler:            p.Menu,
		SettingHandler:         p.Setting,
		UserSettingHandler:     p.UserSetting,
		PATHandler:             p.PAT,
		AuditLogHandler:        p.AuditLogH,
		AdminUserHandler:       p.AdminUser,
		UserProfileHandler:     p.UserProfile,
		OverviewHandler:        p.Overview,
		TwoFAHandler:           p.TwoFA,
		CacheHandler:           p.Cache,
	}

	return adapthttp.SetupRouterWithDeps(deps)
}

// StartHTTPServer 启动 HTTP 服务器并管理生命周期。
func StartHTTPServer(lc fx.Lifecycle, cfg *config.Config, router *gin.Engine) {
	srv := &http.Server{
		Addr:              cfg.Server.Addr,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			lc := &net.ListenConfig{}
			ln, err := lc.Listen(ctx, "tcp", srv.Addr)
			if err != nil {
				return err
			}
			slog.Info("Starting HTTP server", "address", srv.Addr)
			go func() {
				if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
					slog.Error("Server error", "error", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			slog.Info("Stopping HTTP server")
			return srv.Shutdown(ctx)
		},
	})
}
