package bootstrap

import (
	"github.com/gin-gonic/gin"

	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http"
	"github.com/lwmacct/251117-go-ddd-template/internal/config"
)

// newRouter 初始化路由
// 使用 RouterDependencies 参数对象模式，简化依赖传递
func newRouter(cfg *config.Config, infra *InfrastructureModule, services *ServicesModule, usecases *UseCasesModule, handlers *HandlersModule) *gin.Engine {
	deps := &http.RouterDependencies{
		Config:                 cfg,
		RedisClient:            infra.RedisClient,
		CreateLogHandler:       usecases.AuditLog.CreateLog,
		JWTManager:             services.JWT,
		PATService:             services.PAT,
		PermissionCacheService: services.PermissionCache,
		HealthHandler:          handlers.Health,
		AuthHandler:            handlers.Auth,
		CaptchaHandler:         handlers.Captcha,
		RoleHandler:            handlers.Role,
		MenuHandler:            handlers.Menu,
		SettingHandler:         handlers.Setting,
		PATHandler:             handlers.PAT,
		AuditLogHandler:        handlers.AuditLog,
		AdminUserHandler:       handlers.AdminUser,
		UserProfileHandler:     handlers.UserProfile,
		OverviewHandler:        handlers.Overview,
		TwoFAHandler:           handlers.TwoFA,
		CacheHandler:           handlers.Cache,
	}

	return http.SetupRouterWithDeps(deps)
}
