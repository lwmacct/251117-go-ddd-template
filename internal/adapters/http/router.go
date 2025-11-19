package http

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/handler"
	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/middleware"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/auditlog"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/pat"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/role"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
	infraauth "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/auth"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/config"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// SetupRouter 配置路由
func SetupRouter(
	cfg *config.Config,
	db *gorm.DB,
	redisClient *redis.Client,
	userRepo user.Repository,
	roleRepo role.RoleRepository,
	permissionRepo role.PermissionRepository,
	auditLogRepo auditlog.Repository,
	patRepo pat.Repository,
	jwtManager *infraauth.JWTManager,
	tokenGenerator *infraauth.TokenGenerator,
	patService *infraauth.PATService,
	authService *infraauth.Service,
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
		authHandler := handler.NewAuthHandler(authService)
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// 管理员路由 (/api/admin/*) - 使用三段式权限控制
		admin := api.Group("/admin")
		admin.Use(middleware.Auth(jwtManager, patService, userRepo))
		admin.Use(middleware.AuditMiddleware(auditLogRepo))
		admin.Use(middleware.RequireRole("admin"))
		{
			// 用户管理
			adminUserHandler := handler.NewAdminUserHandler(userRepo)
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
		}

		// 用户路由 (/api/user/*) - 使用三段式权限控制
		userGroup := api.Group("/user")
		userGroup.Use(middleware.Auth(jwtManager, patService, userRepo))
		{
			// 个人资料管理
			userProfileHandler := handler.NewUserProfileHandler(userRepo)
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
