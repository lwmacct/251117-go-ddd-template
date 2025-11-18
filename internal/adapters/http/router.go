package http

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/handler"
	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/middleware"
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
	jwtManager *infraauth.JWTManager,
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

		// 需要认证的路由
		authenticated := api.Group("")
		authenticated.Use(middleware.JWTAuth(jwtManager))
		{
			// 当前用户信息
			authenticated.GET("/auth/me", authHandler.Me)

			// 用户管理 (需要认证)
			userHandler := handler.NewUserHandler(userRepo)
			authenticated.GET("/users", userHandler.List)
			authenticated.GET("/users/:id", userHandler.GetByID)
			authenticated.PUT("/users/:id", userHandler.Update)
			authenticated.DELETE("/users/:id", userHandler.Delete)
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
			// 获取请求的文件路径 (已经移除了 /docs 前缀)
			reqPath := c.Param("filepath")
			if reqPath == "/" || reqPath == "" {
				reqPath = "/index.html"
			}

			// 构建完整文件路径
			fullPath := filepath.Join(cfg.Server.DocsDir, reqPath)

			// 检查文件是否存在
			if _, err := os.Stat(fullPath); err == nil {
				c.File(fullPath)
				return
			}

			// 如果路径不存在，尝试添加 .html 扩展名 (VitePress 清洁 URL)
			if !strings.HasSuffix(reqPath, ".html") && !strings.Contains(reqPath, ".") {
				htmlPath := filepath.Join(cfg.Server.DocsDir, reqPath+".html")
				if _, err := os.Stat(htmlPath); err == nil {
					c.File(htmlPath)
					return
				}
			}

			// 文件不存在，返回 index.html (用于 SPA 路由)
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
			// 构建文件路径
			path := filepath.Join(cfg.Server.StaticDir, c.Request.URL.Path)

			// 检查文件是否存在
			if _, err := os.Stat(path); err == nil {
				c.File(path)
				return
			}

			// 文件不存在，返回 index.html (用于 SPA 路由)
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
