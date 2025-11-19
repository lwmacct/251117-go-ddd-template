package http

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/handler"
	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/middleware"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/auditlog"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/captcha"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/menu"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/role"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
	infraauth "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/auth"
	infracaptcha "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/captcha"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/config"
	infratwofa "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/twofa"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// SetupRouter 配置路由
func SetupRouter(
	cfg *config.Config,
	db *gorm.DB,
	redisClient *redis.Client,
	userCommandRepo user.CommandRepository,
	userQueryRepo user.QueryRepository,
	roleRepo role.RoleRepository,
	permissionRepo role.PermissionRepository,
	auditLogRepo auditlog.Repository,
	captchaRepo captcha.Repository,
	menuCommandRepo menu.CommandRepository,
	menuQueryRepo menu.QueryRepository,
	settingCommandRepo setting.CommandRepository,
	settingQueryRepo setting.QueryRepository,
	jwtManager *infraauth.JWTManager,
	tokenGenerator *infraauth.TokenGenerator,
	patService *infraauth.PATService,
	authService *infraauth.Service,
	captchaService *infracaptcha.Service,
	twofaService *infratwofa.Service,
	authHandler *handler.AuthHandler,
) *gin.Engine {
	r := gin.New()

	// 全局中间件
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())

	// 健康检查 (包含数据库和 Redis 连接检查)
	healthHandler := handler.NewHealthHandler(db, redisClient)
	r.GET("/health", healthHandler.Check)

	// API 路由组
	api := r.Group("/api")
	{
		// 认证路由 (公开)
		captchaHandler := handler.NewCaptchaHandler(captchaRepo, captchaService, cfg.Auth.DevSecret)
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.GET("/captcha", captchaHandler.GetCaptcha) // 获取验证码（公开）
		}

		// 认证用户路由
		authUser := api.Group("/auth")
		authUser.Use(middleware.Auth(jwtManager, patService, userQueryRepo))
		{
			authUser.GET("/me", authHandler.Me) // 获取当前用户信息
		}

		// 2FA 路由（需要认证）
		twofaHandler := handler.NewTwoFAHandler(twofaService)
		twofa := api.Group("/auth/2fa")
		twofa.Use(middleware.Auth(jwtManager, patService, userQueryRepo))
		{
			twofa.POST("/setup", twofaHandler.Setup)               // 设置 2FA
			twofa.POST("/verify", twofaHandler.VerifyAndEnable)    // 验证并启用 2FA
			twofa.POST("/disable", twofaHandler.Disable)           // 禁用 2FA
			twofa.GET("/status", twofaHandler.GetStatus)           // 获取 2FA 状态
		}

		// 管理员路由 (/api/admin/*) - 使用三段式权限控制
		admin := api.Group("/admin")
		admin.Use(middleware.Auth(jwtManager, patService, userQueryRepo))
		admin.Use(middleware.AuditMiddleware(auditLogRepo))
		admin.Use(middleware.RequireRole("admin"))
		{
			// 用户管理
			adminUserHandler := handler.NewAdminUserHandler(userCommandRepo, userQueryRepo)
			admin.POST("/users", middleware.RequirePermission("admin:users:create"), adminUserHandler.CreateUser)
			admin.GET("/users", middleware.RequirePermission("admin:users:read"), adminUserHandler.ListUsers)
			admin.GET("/users/:id", middleware.RequirePermission("admin:users:read"), adminUserHandler.GetUser)
			admin.PUT("/users/:id", middleware.RequirePermission("admin:users:update"), adminUserHandler.UpdateUser)
			admin.DELETE("/users/:id", middleware.RequirePermission("admin:users:delete"), adminUserHandler.DeleteUser)
			admin.PUT("/users/:id/roles", middleware.RequirePermission("admin:users:update"), adminUserHandler.AssignRoles)

			// 角色管理
			roleHandler := handler.NewRoleHandler(roleRepo, permissionRepo)
			admin.POST("/roles", middleware.RequirePermission("admin:roles:create"), roleHandler.CreateRole)
			admin.GET("/roles", middleware.RequirePermission("admin:roles:read"), roleHandler.ListRoles)
			admin.GET("/roles/:id", middleware.RequirePermission("admin:roles:read"), roleHandler.GetRole)
			admin.PUT("/roles/:id", middleware.RequirePermission("admin:roles:update"), roleHandler.UpdateRole)
			admin.DELETE("/roles/:id", middleware.RequirePermission("admin:roles:delete"), roleHandler.DeleteRole)
			admin.PUT("/roles/:id/permissions", middleware.RequirePermission("admin:roles:update"), roleHandler.SetPermissions)

			// 权限列表
			admin.GET("/permissions", middleware.RequirePermission("admin:permissions:read"), roleHandler.ListPermissions)

			// 审计日志
			auditLogHandler := handler.NewAuditLogHandler(auditLogRepo)
			admin.GET("/audit-logs", middleware.RequirePermission("admin:audit_logs:read"), auditLogHandler.ListLogs)
			admin.GET("/audit-logs/:id", middleware.RequirePermission("admin:audit_logs:read"), auditLogHandler.GetLog)

			// 菜单管理
			menuHandler := handler.NewMenuHandler(menuCommandRepo, menuQueryRepo)
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
			settingHandler := handler.NewSettingHandler(settingCommandRepo, settingQueryRepo)
			admin.GET("/settings", middleware.RequirePermission("admin:settings:read"), settingHandler.GetSettings)
			admin.GET("/settings/:key", middleware.RequirePermission("admin:settings:read"), settingHandler.GetSetting)
			admin.POST("/settings", middleware.RequirePermission("admin:settings:create"), settingHandler.CreateSetting)
			admin.PUT("/settings/:key", middleware.RequirePermission("admin:settings:update"), settingHandler.UpdateSetting)
			admin.DELETE("/settings/:key", middleware.RequirePermission("admin:settings:delete"), settingHandler.DeleteSetting)
			admin.PUT("/settings/batch", middleware.RequirePermission("admin:settings:update"), settingHandler.BatchUpdateSettings)
		}

		// 用户路由 (/api/user/*) - 使用三段式权限控制
		userGroup := api.Group("/user")
		userGroup.Use(middleware.Auth(jwtManager, patService, userQueryRepo))
		{
			// 个人资料管理
			userProfileHandler := handler.NewUserProfileHandler(userCommandRepo, userQueryRepo)
			userGroup.GET("/me", middleware.RequirePermission("user:profile:read"), userProfileHandler.GetProfile)
			userGroup.PUT("/me", middleware.RequirePermission("user:profile:update"), userProfileHandler.UpdateProfile)
			userGroup.PUT("/me/password", middleware.RequirePermission("user:password:update"), userProfileHandler.ChangePassword)
			userGroup.DELETE("/me", middleware.RequirePermission("user:profile:delete"), userProfileHandler.DeleteAccount)

			// Personal Access Token 管理
			patHandler := handler.NewPATHandler(patService)
			userGroup.POST("/tokens", middleware.RequirePermission("user:tokens:create"), patHandler.CreateToken)
			userGroup.GET("/tokens", middleware.RequirePermission("user:tokens:read"), patHandler.ListTokens)
			userGroup.GET("/tokens/:id", middleware.RequirePermission("user:tokens:read"), patHandler.GetToken)
			userGroup.DELETE("/tokens/:id", middleware.RequirePermission("user:tokens:delete"), patHandler.RevokeToken)
		}

		// 缓存操作示例 (公开，仅用于演示)
		cacheHandler := handler.NewCacheHandler(redisClient, cfg.Data.RedisKeyPrefix)
		api.POST("/cache", cacheHandler.SetCache)
		api.GET("/cache/:key", cacheHandler.GetCache)
		api.DELETE("/cache/:key", cacheHandler.DeleteCache)
	}

	// 提供 VitePress 文档服务 (通过 /docs 路由访问)
	if cfg.Server.DocsDir != "" {
		docs := r.Group("/docs")
		docs.GET("/*filepath", func(c *gin.Context) {
			reqPath := c.Param("filepath")
			if reqPath == "/" || reqPath == "" {
				reqPath = "/index.html"
			}

			fullPath := filepath.Join(cfg.Server.DocsDir, reqPath)

			if _, err := os.Stat(fullPath); err == nil {
				c.File(fullPath)
				return
			}

			if !strings.HasSuffix(reqPath, ".html") && !strings.Contains(reqPath, ".") {
				htmlPath := filepath.Join(cfg.Server.DocsDir, reqPath+".html")
				if _, err := os.Stat(htmlPath); err == nil {
					c.File(htmlPath)
					return
				}
			}

			indexPath := filepath.Join(cfg.Server.DocsDir, "index.html")
			if _, err := os.Stat(indexPath); err == nil {
				c.File(indexPath)
			} else {
				c.Status(404)
			}
		})
	}

	// 提供静态文件服务 (使用 NoRoute 避免与 API 路由冲突)
	if cfg.Server.StaticDir != "" {
		r.NoRoute(func(c *gin.Context) {
			path := filepath.Join(cfg.Server.StaticDir, c.Request.URL.Path)

			if _, err := os.Stat(path); err == nil {
				c.File(path)
				return
			}

			indexPath := filepath.Join(cfg.Server.StaticDir, "index.html")
			if _, err := os.Stat(indexPath); err == nil {
				c.File(indexPath)
			} else {
				c.Status(404)
			}
		})
	}

	return r
}
