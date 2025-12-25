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
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	// 引入第三方包
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	// Swagger 文档
	_ "github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/docs" // Swagger docs

	"github.com/lwmacct/251117-go-ddd-template/internal/config"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// 引入处理器和中间件包
	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/handler"
	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/middleware"
	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/response"

	// 引入应用层包
	"github.com/lwmacct/251117-go-ddd-template/internal/application/auditlog"

	// 引入基础设施包
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/auth"
)

// RouterDependencies 路由依赖项（参数对象模式）
// 将所有依赖项聚合为单一结构体，减少函数参数数量
type RouterDependencies struct {
	// Config
	Config      *config.Config
	RedisClient *redis.Client

	// Application Handlers (for middleware)
	CreateLogHandler *auditlog.CreateLogHandler

	// Infrastructure Services
	JWTManager             *auth.JWTManager
	PATService             *auth.PATService
	PermissionCacheService *auth.PermissionCacheService

	// HTTP Handlers
	HealthHandler      *handler.HealthHandler
	AuthHandler        *handler.AuthHandler
	CaptchaHandler     *handler.CaptchaHandler
	RoleHandler        *handler.RoleHandler
	MenuHandler        *handler.MenuHandler
	SettingHandler     *handler.SettingHandler
	PATHandler         *handler.PATHandler
	AuditLogHandler    *handler.AuditLogHandler
	AdminUserHandler   *handler.AdminUserHandler
	UserProfileHandler *handler.UserProfileHandler
	OverviewHandler    *handler.OverviewHandler
	TwoFAHandler       *handler.TwoFAHandler
	CacheHandler       *handler.CacheHandler
}

// SetupRouterWithDeps 使用依赖对象配置路由（推荐方式）
// 通过参数对象模式，将所有依赖聚合为单一结构体，简化函数签名
func SetupRouterWithDeps(deps *RouterDependencies) *gin.Engine {
	cfg := deps.Config

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
	// 自定义 Recovery，输出 panic 到 slog，生产环境隐藏详细错误
	r.Use(gin.CustomRecovery(func(c *gin.Context, recovered any) {
		slog.Error("PANIC recovered", "error", recovered, "path", c.Request.URL.Path, "method", c.Request.Method)
		if cfg.Server.Env != "production" {
			response.InternalError(c, fmt.Sprintf("%v", recovered))
		} else {
			response.InternalError(c)
		}
		c.Abort()
	}))
	r.Use(middleware.CORS())
	// 使用基于 slog 的日志中间件，跳过健康检查端点
	r.Use(middleware.LoggerSkipPaths("/health"))

	// 健康检查
	r.GET("/health", deps.HealthHandler.Check)

	// Swagger API 文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API 路由组
	setupAPIRoutes(r, deps)

	// 静态文件服务
	setupStaticRoutes(r, cfg)

	return r
}

// setupAPIRoutes 配置 API 路由组
func setupAPIRoutes(r *gin.Engine, deps *RouterDependencies) {
	api := r.Group("/api")

	// 认证路由 (公开)
	auth := api.Group("/auth")
	{
		auth.POST("/register", deps.AuthHandler.Register)
		auth.POST("/login", deps.AuthHandler.Login)
		auth.POST("/login/2fa", deps.AuthHandler.Login2FA)
		auth.POST("/refresh", deps.AuthHandler.RefreshToken)
		auth.GET("/captcha", deps.CaptchaHandler.GetCaptcha)
	}

	// 2FA 路由（需要认证）
	twofa := api.Group("/auth/2fa")
	twofa.Use(middleware.Auth(deps.JWTManager, deps.PATService, deps.PermissionCacheService))
	{
		twofa.POST("/setup", deps.TwoFAHandler.Setup)            // 设置 2FA
		twofa.POST("/verify", deps.TwoFAHandler.VerifyAndEnable) // 验证并启用 2FA
		twofa.POST("/disable", deps.TwoFAHandler.Disable)        // 禁用 2FA
		twofa.GET("/status", deps.TwoFAHandler.GetStatus)        // 获取 2FA 状态
	}

	// 管理员路由 (/api/admin/*) - 使用三段式权限控制
	admin := api.Group("/admin")
	admin.Use(middleware.Auth(deps.JWTManager, deps.PATService, deps.PermissionCacheService))
	admin.Use(middleware.AuditMiddleware(deps.CreateLogHandler))
	admin.Use(middleware.RequireRole("admin"))
	{
		// 用户管理
		admin.POST("/users", middleware.RequirePermission("admin:users:create"), deps.AdminUserHandler.CreateUser)
		admin.POST("/users/batch", middleware.RequirePermission("admin:users:create"), deps.AdminUserHandler.BatchCreateUsers)
		admin.GET("/users", middleware.RequirePermission("admin:users:read"), deps.AdminUserHandler.ListUsers)
		admin.GET("/users/:id", middleware.RequirePermission("admin:users:read"), deps.AdminUserHandler.GetUser)
		admin.PUT("/users/:id", middleware.RequirePermission("admin:users:update"), deps.AdminUserHandler.UpdateUser)
		admin.DELETE("/users/:id", middleware.RequirePermission("admin:users:delete"), deps.AdminUserHandler.DeleteUser)
		admin.PUT("/users/:id/roles", middleware.RequirePermission("admin:users:update"), deps.AdminUserHandler.AssignRoles)

		// 角色管理
		admin.POST("/roles", middleware.RequirePermission("admin:roles:create"), deps.RoleHandler.CreateRole)
		admin.GET("/roles", middleware.RequirePermission("admin:roles:read"), deps.RoleHandler.ListRoles)
		admin.GET("/roles/:id", middleware.RequirePermission("admin:roles:read"), deps.RoleHandler.GetRole)
		admin.PUT("/roles/:id", middleware.RequirePermission("admin:roles:update"), deps.RoleHandler.UpdateRole)
		admin.DELETE("/roles/:id", middleware.RequirePermission("admin:roles:delete"), deps.RoleHandler.DeleteRole)
		admin.PUT("/roles/:id/permissions", middleware.RequirePermission("admin:roles:update"), deps.RoleHandler.SetPermissions)

		// 权限列表
		admin.GET("/permissions", middleware.RequirePermission("admin:permissions:read"), deps.RoleHandler.ListPermissions)

		// 审计日志
		admin.GET("/auditlogs", middleware.RequirePermission("admin:audit_logs:read"), deps.AuditLogHandler.ListLogs)
		admin.GET("/auditlogs/:id", middleware.RequirePermission("admin:audit_logs:read"), deps.AuditLogHandler.GetLog)

		// 菜单管理
		admin.POST("/menus", middleware.RequirePermission("admin:menus:create"), deps.MenuHandler.Create)
		admin.GET("/menus", middleware.RequirePermission("admin:menus:read"), deps.MenuHandler.List)
		admin.GET("/menus/:id", middleware.RequirePermission("admin:menus:read"), deps.MenuHandler.Get)
		admin.PUT("/menus/:id", middleware.RequirePermission("admin:menus:update"), deps.MenuHandler.Update)
		admin.DELETE("/menus/:id", middleware.RequirePermission("admin:menus:delete"), deps.MenuHandler.Delete)
		admin.POST("/menus/reorder", middleware.RequirePermission("admin:menus:update"), deps.MenuHandler.Reorder)

		// 系统概览
		admin.GET("/overview/stats", middleware.RequirePermission("admin:overview:read"), deps.OverviewHandler.GetStats)

		// 系统配置
		admin.GET("/settings", middleware.RequirePermission("admin:settings:read"), deps.SettingHandler.GetSettings)
		admin.GET("/settings/:key", middleware.RequirePermission("admin:settings:read"), deps.SettingHandler.GetSetting)
		admin.POST("/settings", middleware.RequirePermission("admin:settings:create"), deps.SettingHandler.CreateSetting)
		admin.PUT("/settings/:key", middleware.RequirePermission("admin:settings:update"), deps.SettingHandler.UpdateSetting)
		admin.DELETE("/settings/:key", middleware.RequirePermission("admin:settings:delete"), deps.SettingHandler.DeleteSetting)
		admin.POST("/settings/batch", middleware.RequirePermission("admin:settings:update"), deps.SettingHandler.BatchUpdateSettings)
	}

	// 用户路由 (/api/user/*) - 使用三段式权限控制
	userGroup := api.Group("/user")
	userGroup.Use(middleware.Auth(deps.JWTManager, deps.PATService, deps.PermissionCacheService))
	{
		// 个人资料管理
		userGroup.GET("/profile", middleware.RequirePermission("user:profile:read"), deps.UserProfileHandler.GetProfile)
		userGroup.PUT("/profile", middleware.RequirePermission("user:profile:update"), deps.UserProfileHandler.UpdateProfile)
		userGroup.PUT("/password", middleware.RequirePermission("user:password:update"), deps.UserProfileHandler.ChangePassword)
		userGroup.DELETE("/account", middleware.RequirePermission("user:profile:delete"), deps.UserProfileHandler.DeleteAccount)

		// Personal Access Token 管理
		userGroup.POST("/tokens", middleware.RequirePermission("user:tokens:create"), deps.PATHandler.CreateToken)
		userGroup.GET("/tokens", middleware.RequirePermission("user:tokens:read"), deps.PATHandler.ListTokens)
		userGroup.GET("/tokens/:id", middleware.RequirePermission("user:tokens:read"), deps.PATHandler.GetToken)
		userGroup.DELETE("/tokens/:id", middleware.RequirePermission("user:tokens:delete"), deps.PATHandler.DeleteToken)
		userGroup.PATCH("/tokens/:id/disable", middleware.RequirePermission("user:tokens:disable"), deps.PATHandler.DisableToken)
		userGroup.PATCH("/tokens/:id/enable", middleware.RequirePermission("user:tokens:enable"), deps.PATHandler.EnableToken)
	}

	// 缓存操作示例 (公开，仅用于演示)
	api.POST("/cache", deps.CacheHandler.SetCache)
	api.GET("/cache/:key", deps.CacheHandler.GetCache)
	api.DELETE("/cache/:key", deps.CacheHandler.DeleteCache)
}

// setupStaticRoutes 配置静态文件服务路由
func setupStaticRoutes(r *gin.Engine, cfg *config.Config) {
	// 提供 VitePress 文档服务 (通过 /docs 路由访问)
	if cfg.Server.DocsDist != "" {
		docs := r.Group("/docs")
		docs.GET("/*filepath", serveVitePressHandler(cfg.Server.DocsDist))
	}

	// 提供静态文件服务 (使用 NoRoute 避免与 API 路由冲突)
	if cfg.Server.WebDist != "" {
		r.NoRoute(serveSPAHandler(cfg.Server.WebDist))
	}
}

// serveVitePressHandler 返回 VitePress 文档服务处理函数
func serveVitePressHandler(distDocs string) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqPath := c.Param("filepath")
		if reqPath == "/" || reqPath == "" {
			reqPath = "/index.html"
		}

		// 尝试直接提供文件
		fullPath := filepath.Join(distDocs, reqPath)
		if fileExists(fullPath) {
			c.File(fullPath)
			return
		}

		// VitePress clean URL: 尝试 .html 后缀
		if !strings.HasSuffix(reqPath, ".html") && !strings.Contains(reqPath, ".") {
			htmlPath := filepath.Join(distDocs, reqPath+".html")
			if fileExists(htmlPath) {
				c.File(htmlPath)
				return
			}
		}

		// fallback 到 index.html 或 404
		serveIndexOrNotFound(c, distDocs)
	}
}

// serveSPAHandler 返回 SPA 静态文件服务处理函数
func serveSPAHandler(distWeb string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// API 路由返回 JSON 404，避免 SPA fallback 干扰
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			response.NotFound(c, "endpoint")
			return
		}

		// 非 API 路径使用 SPA fallback
		path := filepath.Join(distWeb, c.Request.URL.Path)
		if fileExists(path) {
			c.File(path)
			return
		}

		// fallback 到 index.html 或 404
		serveIndexOrNotFound(c, distWeb)
	}
}

// fileExists 检查文件是否存在
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// serveIndexOrNotFound 提供 index.html 或返回 404
func serveIndexOrNotFound(c *gin.Context, dir string) {
	indexPath := filepath.Join(dir, "index.html")
	if fileExists(indexPath) {
		c.File(indexPath)
		return
	}
	c.Status(404)
}
