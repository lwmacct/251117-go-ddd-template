package bootstrap

import (
	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/handler"
	"github.com/lwmacct/251117-go-ddd-template/internal/config"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/health"
)

// newHandlersModule 初始化 HTTP Handler 模块
// 依赖：UseCasesModule, InfrastructureModule, config.Config
func newHandlersModule(cfg *config.Config, infra *InfrastructureModule, useCases *UseCasesModule) *HandlersModule {
	m := &HandlersModule{}

	// Health Handler
	healthChecker := health.NewSystemChecker(infra.DB, infra.RedisClient)
	m.Health = handler.NewHealthHandler(healthChecker)

	// Auth Handler
	m.Auth = handler.NewAuthHandler(
		useCases.Auth.Login,
		useCases.Auth.Login2FA,
		useCases.Auth.Register,
		useCases.Auth.RefreshToken,
	)

	// Captcha Handler
	m.Captcha = handler.NewCaptchaHandler(useCases.Captcha.Generate, cfg.Auth.DevSecret)

	// Admin User Handler
	m.AdminUser = handler.NewAdminUserHandler(
		useCases.User.Create,
		useCases.User.Update,
		useCases.User.Delete,
		useCases.User.AssignRoles,
		useCases.User.BatchCreate,
		useCases.User.Get,
		useCases.User.List,
	)

	// User Profile Handler
	m.UserProfile = handler.NewUserProfileHandler(
		useCases.User.Get,
		useCases.User.Update,
		useCases.User.ChangePassword,
		useCases.User.Delete,
	)

	// Role Handler
	m.Role = handler.NewRoleHandler(
		useCases.Role.Create,
		useCases.Role.Update,
		useCases.Role.Delete,
		useCases.Role.SetPermissions,
		useCases.Role.Get,
		useCases.Role.List,
		useCases.Role.ListPermissions,
	)

	// Menu Handler
	m.Menu = handler.NewMenuHandler(
		useCases.Menu.Create,
		useCases.Menu.Update,
		useCases.Menu.Delete,
		useCases.Menu.Reorder,
		useCases.Menu.Get,
		useCases.Menu.List,
	)

	// Setting Handler
	m.Setting = handler.NewSettingHandler(
		useCases.Setting.Create,
		useCases.Setting.Update,
		useCases.Setting.Delete,
		useCases.Setting.BatchUpdate,
		useCases.Setting.Get,
		useCases.Setting.List,
	)

	// PAT Handler
	m.PAT = handler.NewPATHandler(
		useCases.PAT.Create,
		useCases.PAT.Delete,
		useCases.PAT.Disable,
		useCases.PAT.Enable,
		useCases.PAT.Get,
		useCases.PAT.List,
	)

	// AuditLog Handler
	m.AuditLog = handler.NewAuditLogHandler(
		useCases.AuditLog.List,
		useCases.AuditLog.Get,
	)

	// Overview Handler
	m.Overview = handler.NewOverviewHandler(useCases.Stats.GetStats)

	// TwoFA Handler
	m.TwoFA = handler.NewTwoFAHandler(
		useCases.TwoFA.Setup,
		useCases.TwoFA.VerifyEnable,
		useCases.TwoFA.Disable,
		useCases.TwoFA.GetStatus,
	)

	// Cache Handler (for demo)
	m.Cache = handler.NewCacheHandler(
		useCases.Cache.Set,
		useCases.Cache.Get,
		useCases.Cache.Delete,
	)

	return m
}
