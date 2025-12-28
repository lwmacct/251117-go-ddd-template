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
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/health"
)

// HandlersModule aggregates all HTTP handlers.
type HandlersModule struct {
	Health      *handler.HealthHandler
	Auth        *handler.AuthHandler
	Captcha     *handler.CaptchaHandler
	AdminUser   *handler.AdminUserHandler
	User        *handler.UserHandler
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

// HTTPModule provides HTTP handlers and router.
var HTTPModule = fx.Module("http",
	fx.Provide(
		newHealthChecker,
		newHandlersModule,
		newRouter,
	),
)

func newHealthChecker(db *gorm.DB, redisClient *redis.Client) *health.SystemChecker {
	return health.NewSystemChecker(db, redisClient)
}

func newHandlersModule(
	cfg *config.Config,
	redisClient *redis.Client,
	healthChecker *health.SystemChecker,
	useCases *UseCasesModule,
) *HandlersModule {
	keyPrefix := cfg.Data.RedisKeyPrefix

	return &HandlersModule{
		Health: handler.NewHealthHandler(healthChecker),
		Auth: handler.NewAuthHandler(
			useCases.Auth.Login,
			useCases.Auth.Login2FA,
			useCases.Auth.Register,
			useCases.Auth.RefreshToken,
		),
		Captcha: handler.NewCaptchaHandler(useCases.Captcha.Generate, cfg.Auth.DevSecret),
		AdminUser: handler.NewAdminUserHandler(
			useCases.User.Create,
			useCases.User.Update,
			useCases.User.Delete,
			useCases.User.AssignRoles,
			useCases.User.BatchCreate,
			useCases.User.Get,
			useCases.User.List,
		),
		User: handler.NewUserHandler(
			useCases.User.Create,
			useCases.User.Update,
			useCases.User.Delete,
			useCases.User.Get,
			useCases.User.List,
		),
		UserProfile: handler.NewUserProfileHandler(
			useCases.User.Get,
			useCases.User.Update,
			useCases.User.ChangePassword,
			useCases.User.Delete,
		),
		Role: handler.NewRoleHandler(
			useCases.Role.Create,
			useCases.Role.Update,
			useCases.Role.Delete,
			useCases.Role.SetPermissions,
			useCases.Role.Get,
			useCases.Role.List,
			useCases.Role.ListPermissions,
		),
		Menu: handler.NewMenuHandler(
			useCases.Menu.Create,
			useCases.Menu.Update,
			useCases.Menu.Delete,
			useCases.Menu.Reorder,
			useCases.Menu.Get,
			useCases.Menu.List,
		),
		Setting: handler.NewSettingHandler(
			useCases.Setting.Create,
			useCases.Setting.Update,
			useCases.Setting.Delete,
			useCases.Setting.BatchUpdate,
			useCases.Setting.Get,
			useCases.Setting.List,
			useCases.Setting.ListSchema,
			useCases.Setting.CreateCategory,
			useCases.Setting.UpdateCategory,
			useCases.Setting.DeleteCategory,
			useCases.Setting.GetCategory,
			useCases.Setting.ListCategories,
		),
		UserSetting: handler.NewUserSettingHandler(
			useCases.UserSetting.Set,
			useCases.UserSetting.BatchSet,
			useCases.UserSetting.Reset,
			useCases.UserSetting.ResetAll,
			useCases.UserSetting.Get,
			useCases.UserSetting.List,
			useCases.UserSetting.ListSchema,
			useCases.UserSetting.ListCategories,
		),
		PAT: handler.NewPATHandler(
			useCases.PAT.Create,
			useCases.PAT.Delete,
			useCases.PAT.Disable,
			useCases.PAT.Enable,
			useCases.PAT.Get,
			useCases.PAT.List,
		),
		AuditLog: handler.NewAuditLogHandler(
			useCases.AuditLog.List,
			useCases.AuditLog.Get,
		),
		Overview: handler.NewOverviewHandler(useCases.Stats.GetStats),
		TwoFA: handler.NewTwoFAHandler(
			useCases.TwoFA.Setup,
			useCases.TwoFA.VerifyEnable,
			useCases.TwoFA.Disable,
			useCases.TwoFA.GetStatus,
		),
		Cache: handler.NewCacheHandler(
			cache.NewInfoHandler(redisClient, keyPrefix),
			cache.NewScanKeysHandler(redisClient, keyPrefix),
			cache.NewGetKeyHandler(redisClient, keyPrefix),
			cache.NewDeleteHandler(redisClient, keyPrefix),
		),
	}
}

func newRouter(
	cfg *config.Config,
	redisClient *redis.Client,
	services *ServicesModule,
	useCases *UseCasesModule,
	handlers *HandlersModule,
) *gin.Engine {
	deps := &adapthttp.RouterDependencies{
		Config:                 cfg,
		RedisClient:            redisClient,
		CreateLogHandler:       useCases.AuditLog.CreateLog,
		JWTManager:             services.JWT,
		PATService:             services.PAT,
		PermissionCacheService: services.PermissionCache,
		HealthHandler:          handlers.Health,
		AuthHandler:            handlers.Auth,
		CaptchaHandler:         handlers.Captcha,
		RoleHandler:            handlers.Role,
		MenuHandler:            handlers.Menu,
		SettingHandler:         handlers.Setting,
		UserSettingHandler:     handlers.UserSetting,
		PATHandler:             handlers.PAT,
		AuditLogHandler:        handlers.AuditLog,
		AdminUserHandler:       handlers.AdminUser,
		UserHandler:            handlers.User,
		UserProfileHandler:     handlers.UserProfile,
		OverviewHandler:        handlers.Overview,
		TwoFAHandler:           handlers.TwoFA,
		CacheHandler:           handlers.Cache,
	}

	return adapthttp.SetupRouterWithDeps(deps)
}

// StartHTTPServer starts the HTTP server with lifecycle management.
// Use this with fx.Invoke for applications that want Fx to manage the server.
func StartHTTPServer(lc fx.Lifecycle, cfg *config.Config, router *gin.Engine) {
	srv := &http.Server{
		Addr:              cfg.Server.Addr,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second, // 防止 Slowloris 攻击
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
