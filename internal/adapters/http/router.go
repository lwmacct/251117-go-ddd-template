// Package http 提供 HTTP 适配器层的实现。
//
// 本包是 DDD 架构的适配器层入口，负责：
//   - 路由配置：基于 Gin 框架的 RESTful API 路由定义
//   - 中间件集成：认证、授权、日志、CORS 等中间件
//   - 静态文件服务：前端 SPA 和文档服务
//
// 路由结构：
//   - /api/auth/*: 认证相关（登录、注册、刷新令牌）
//   - /api/admin/*: 管理后台（用户、角色、权限、菜单管理）
//   - /api/user/*: 用户中心（个人资料、PAT 管理）
//   - /swagger/*: API 文档
//   - /docs/*: VitePress 文档
//   - /health: 健康检查
//
// 权限控制采用三段式格式：domain:resource:action
// 例如：admin:users:create, user:profile:read
package http

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	// 引入第三方包
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	// Swagger 文档
	_ "github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/docs" // Swagger docs

	"github.com/lwmacct/251117-go-ddd-template/internal/config"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// 引入处理器和中间件包
	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/handler"
	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/middleware"

	// 引入领域包
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/auditlog"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/captcha"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"

	// 引入基础设施包
	infra_auth "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/auth"
	infra_captcha "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/captcha"
	infra_twofa "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/twofa"
)

// SetupRouter 配置路由 (完全符合 DDD+CQRS 架构)
func SetupRouter(
	cfg *config.Config,
	db *gorm.DB,
	redisClient *redis.Client,
	userQueryRepo user.QueryRepository,
	auditLogCommandRepo auditlog.CommandRepository,
	captchaCommandRepo captcha.CommandRepository,
	jwtManager *infra_auth.JWTManager,
	patService *infra_auth.PATService,
	permissionCacheService *infra_auth.PermissionCacheService, // 新增：权限缓存服务
	authService *infra_auth.Service,
	captchaService *infra_captcha.Service,
	twofaService *infra_twofa.Service,
	authHandler *handler.AuthHandler,
	roleHandler *handler.RoleHandler,
	menuHandler *handler.MenuHandler,
	settingHandler *handler.SettingHandler,
	patHandler *handler.PATHandler,
	auditLogHandler *handler.AuditLogHandler,
	adminUserHandler *handler.AdminUserHandler,
	userProfileHandler *handler.UserProfileHandler,
) *gin.Engine {
	// 配置 Gin 模式和日志输出
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	// 禁用 Gin 的默认调试输出（路由注册信息等），我们使用 slog
	gin.DefaultWriter = io.Discard
	gin.DefaultErrorWriter = io.Discard

	r := gin.New()

	// 全局中间件
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())
	// 使用基于 slog 的日志中间件，跳过健康检查端点
	r.Use(middleware.LoggerSkipPaths("/health"))

	// 健康检查 (包含数据库和 Redis 连接检查)
	healthHandler := handler.NewHealthHandler(db, redisClient)
	r.GET("/health", healthHandler.Check)

	// Swagger API 文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API 路由组
	api := r.Group("/api")
	{
		// 认证路由 (公开)
		captchaHandler := handler.NewCaptchaHandler(captchaCommandRepo, captchaService, cfg.Auth.DevSecret)
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.GET("/captcha", captchaHandler.GetCaptcha) // 获取验证码（公开）
		}

		// 认证用户路由
		authUser := api.Group("/auth")
		authUser.Use(middleware.Auth(jwtManager, patService, permissionCacheService))
		{
			authUser.GET("/me", authHandler.Me) // 获取当前用户信息
		}

		// 2FA 路由（需要认证）
		twofaHandler := handler.NewTwoFAHandler(twofaService)
		twofa := api.Group("/auth/2fa")
		twofa.Use(middleware.Auth(jwtManager, patService, permissionCacheService))
		{
			twofa.POST("/setup", twofaHandler.Setup)            // 设置 2FA
			twofa.POST("/verify", twofaHandler.VerifyAndEnable) // 验证并启用 2FA
			twofa.POST("/disable", twofaHandler.Disable)        // 禁用 2FA
			twofa.GET("/status", twofaHandler.GetStatus)        // 获取 2FA 状态
		}

		// 管理员路由 (/api/admin/*) - 使用三段式权限控制
		admin := api.Group("/admin")
		admin.Use(middleware.Auth(jwtManager, patService, permissionCacheService))
		admin.Use(middleware.AuditMiddleware(auditLogCommandRepo))
		admin.Use(middleware.RequireRole("admin"))
		{
			// 用户管理
			admin.POST("/users", middleware.RequirePermission("admin:users:create"), adminUserHandler.CreateUser)
			admin.POST("/users/batch", middleware.RequirePermission("admin:users:create"), adminUserHandler.BatchCreateUsers)
			admin.GET("/users", middleware.RequirePermission("admin:users:read"), adminUserHandler.ListUsers)
			admin.GET("/users/:id", middleware.RequirePermission("admin:users:read"), adminUserHandler.GetUser)
			admin.PUT("/users/:id", middleware.RequirePermission("admin:users:update"), adminUserHandler.UpdateUser)
			admin.DELETE("/users/:id", middleware.RequirePermission("admin:users:delete"), adminUserHandler.DeleteUser)
			admin.PUT("/users/:id/roles", middleware.RequirePermission("admin:users:update"), adminUserHandler.AssignRoles)

			// 角色管理
			admin.POST("/roles", middleware.RequirePermission("admin:roles:create"), roleHandler.CreateRole)
			admin.GET("/roles", middleware.RequirePermission("admin:roles:read"), roleHandler.ListRoles)
			admin.GET("/roles/:id", middleware.RequirePermission("admin:roles:read"), roleHandler.GetRole)
			admin.PUT("/roles/:id", middleware.RequirePermission("admin:roles:update"), roleHandler.UpdateRole)
			admin.DELETE("/roles/:id", middleware.RequirePermission("admin:roles:delete"), roleHandler.DeleteRole)
			admin.PUT("/roles/:id/permissions", middleware.RequirePermission("admin:roles:update"), roleHandler.SetPermissions)

			// 权限列表
			admin.GET("/permissions", middleware.RequirePermission("admin:permissions:read"), roleHandler.ListPermissions)

			// 审计日志
			admin.GET("/audit-logs", middleware.RequirePermission("admin:audit_logs:read"), auditLogHandler.ListLogs)
			admin.GET("/audit-logs/:id", middleware.RequirePermission("admin:audit_logs:read"), auditLogHandler.GetLog)

			// 菜单管理
			admin.POST("/menus", middleware.RequirePermission("admin:menus:create"), menuHandler.Create)
			admin.GET("/menus", middleware.RequirePermission("admin:menus:read"), menuHandler.List)
			admin.GET("/menus/:id", middleware.RequirePermission("admin:menus:read"), menuHandler.Get)
			admin.PUT("/menus/:id", middleware.RequirePermission("admin:menus:update"), menuHandler.Update)
			admin.DELETE("/menus/:id", middleware.RequirePermission("admin:menus:delete"), menuHandler.Delete)
			admin.PUT("/menus/reorder", middleware.RequirePermission("admin:menus:update"), menuHandler.Reorder)

			// 系统概览
			overviewHandler := handler.NewOverviewHandler(db)
			admin.GET("/overview/stats", middleware.RequirePermission("admin:overview:read"), overviewHandler.GetStats)

			// 系统配置
			admin.GET("/settings", middleware.RequirePermission("admin:settings:read"), settingHandler.GetSettings)
			admin.GET("/settings/:key", middleware.RequirePermission("admin:settings:read"), settingHandler.GetSetting)
			admin.POST("/settings", middleware.RequirePermission("admin:settings:create"), settingHandler.CreateSetting)
			admin.PUT("/settings/:key", middleware.RequirePermission("admin:settings:update"), settingHandler.UpdateSetting)
			admin.DELETE("/settings/:key", middleware.RequirePermission("admin:settings:delete"), settingHandler.DeleteSetting)
			admin.PUT("/settings/batch", middleware.RequirePermission("admin:settings:update"), settingHandler.BatchUpdateSettings)
		}

		// 用户路由 (/api/user/*) - 使用三段式权限控制
		userGroup := api.Group("/user")
		userGroup.Use(middleware.Auth(jwtManager, patService, permissionCacheService))
		{
			// 个人资料管理
			userGroup.GET("/me", middleware.RequirePermission("user:profile:read"), userProfileHandler.GetProfile)
			userGroup.PUT("/me", middleware.RequirePermission("user:profile:update"), userProfileHandler.UpdateProfile)
			userGroup.PUT("/me/password", middleware.RequirePermission("user:password:update"), userProfileHandler.ChangePassword)
			userGroup.DELETE("/me", middleware.RequirePermission("user:profile:delete"), userProfileHandler.DeleteAccount)

			// Personal Access Token 管理
			userGroup.POST("/tokens", middleware.RequirePermission("user:tokens:create"), patHandler.CreateToken)
			userGroup.GET("/tokens", middleware.RequirePermission("user:tokens:read"), patHandler.ListTokens)
			userGroup.GET("/tokens/:id", middleware.RequirePermission("user:tokens:read"), patHandler.GetToken)
			userGroup.DELETE("/tokens/:id", middleware.RequirePermission("user:tokens:delete"), patHandler.DeleteToken)
			userGroup.PATCH("/tokens/:id/disable", middleware.RequirePermission("user:tokens:disable"), patHandler.DisableToken)
			userGroup.PATCH("/tokens/:id/enable", middleware.RequirePermission("user:tokens:enable"), patHandler.EnableToken)
		}

		// 缓存操作示例 (公开，仅用于演示)
		cacheHandler := handler.NewCacheHandler(redisClient, cfg.Data.RedisKeyPrefix)
		api.POST("/cache", cacheHandler.SetCache)
		api.GET("/cache/:key", cacheHandler.GetCache)
		api.DELETE("/cache/:key", cacheHandler.DeleteCache)
	}

	// 提供 VitePress 文档服务 (通过 /docs 路由访问)
	if cfg.Server.DistDocs != "" {
		docs := r.Group("/docs")
		docs.GET("/*filepath", func(c *gin.Context) {
			reqPath := c.Param("filepath")
			if reqPath == "/" || reqPath == "" {
				reqPath = "/index.html"
			}

			fullPath := filepath.Join(cfg.Server.DistDocs, reqPath)

			if _, err := os.Stat(fullPath); err == nil {
				c.File(fullPath)
				return
			}

			if !strings.HasSuffix(reqPath, ".html") && !strings.Contains(reqPath, ".") {
				htmlPath := filepath.Join(cfg.Server.DistDocs, reqPath+".html")
				if _, err := os.Stat(htmlPath); err == nil {
					c.File(htmlPath)
					return
				}
			}

			indexPath := filepath.Join(cfg.Server.DistDocs, "index.html")
			if _, err := os.Stat(indexPath); err == nil {
				c.File(indexPath)
			} else {
				c.Status(404)
			}
		})
	}

	// 提供静态文件服务 (使用 NoRoute 避免与 API 路由冲突)
	if cfg.Server.DistWeb != "" {
		r.NoRoute(func(c *gin.Context) {
			// API 路由返回 JSON 404，避免 SPA fallback 干扰
			if strings.HasPrefix(c.Request.URL.Path, "/api/") {
				c.JSON(404, gin.H{
					"code":    404,
					"message": "endpoint not found",
					"error":   "route does not exist",
				})
				return
			}

			// 非 API 路径使用 SPA fallback
			path := filepath.Join(cfg.Server.DistWeb, c.Request.URL.Path)

			if _, err := os.Stat(path); err == nil {
				c.File(path)
				return
			}

			indexPath := filepath.Join(cfg.Server.DistWeb, "index.html")
			if _, err := os.Stat(indexPath); err == nil {
				c.File(indexPath)
			} else {
				c.Status(404)
			}
		})
	}

	return r
}
