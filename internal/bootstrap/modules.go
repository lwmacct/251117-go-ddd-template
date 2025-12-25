package bootstrap

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/handler"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/auth"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/captcha"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/event"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/stats"

	_auth "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/auth"
	_captcha "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/captcha"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/persistence"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/twofa"
)

// InfrastructureModule 基础设施模块
// 包含数据库、缓存、事件总线等底层组件
type InfrastructureModule struct {
	DB          *gorm.DB
	RedisClient *redis.Client
	EventBus    event.EventBus
}

// RepositoriesModule 仓储模块
// 聚合所有 CQRS 仓储，按领域划分
type RepositoriesModule struct {
	User       persistence.UserRepositories
	AuditLog   persistence.AuditLogRepositories
	Role       persistence.RoleRepositories
	Permission persistence.PermissionRepositories
	PAT        persistence.PATRepositories
	Menu       persistence.MenuRepositories
	Setting    persistence.SettingRepositories
	TwoFA      persistence.TwoFARepositories

	// 特殊仓储（内存实现）
	CaptchaCommand captcha.CommandRepository
	CaptchaQuery   captcha.QueryRepository

	// 只读仓储
	StatsQuery stats.QueryRepository
}

// ServicesModule 服务模块
// 聚合所有领域服务和基础设施服务
type ServicesModule struct {
	// Domain Services
	Auth auth.Service

	// Infrastructure Services
	JWT             *_auth.JWTManager
	TokenGenerator  auth.TokenGenerator
	LoginSession    *_auth.LoginSessionService
	PermissionCache *_auth.PermissionCacheService
	PAT             *_auth.PATService
	Captcha         *_captcha.Service
	TwoFA           *twofa.Service
}

// HandlersModule HTTP Handler 模块
// 聚合所有 HTTP Handler，统一初始化入口
type HandlersModule struct {
	Health      *handler.HealthHandler
	Auth        *handler.AuthHandler
	Captcha     *handler.CaptchaHandler
	AdminUser   *handler.AdminUserHandler
	UserProfile *handler.UserProfileHandler
	Role        *handler.RoleHandler
	Menu        *handler.MenuHandler
	Setting     *handler.SettingHandler
	PAT         *handler.PATHandler
	AuditLog    *handler.AuditLogHandler
	Overview    *handler.OverviewHandler
	TwoFA       *handler.TwoFAHandler
	Cache       *handler.CacheHandler
}

// RouterModule 路由模块
type RouterModule struct {
	Engine *gin.Engine
}
