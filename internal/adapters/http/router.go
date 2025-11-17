package http

import (
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251117-bd-vmalert/internal/adapters/http/middleware"
	"github.com/lwmacct/251117-bd-vmalert/internal/infrastructure/config"
)

// SetupRouter 配置路由
func SetupRouter(cfg *config.Config) *gin.Engine {
	r := gin.New()

	// 全局中间件
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 提供静态文件服务（使用 NoRoute 避免与 API 路由冲突）
	if cfg.Server.StaticDir != "" {
		r.NoRoute(func(c *gin.Context) {
			// 构建文件路径
			path := filepath.Join(cfg.Server.StaticDir, c.Request.URL.Path)

			// 检查文件是否存在
			if _, err := os.Stat(path); err == nil {
				c.File(path)
				return
			}

			// 文件不存在，返回 index.html（用于 SPA 路由）
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
