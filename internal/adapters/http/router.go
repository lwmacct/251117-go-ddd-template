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
//
// @title           Go DDD Template API
// @version         1.0
// @description     基于 DDD + CQRS 架构的 Go Web 应用模板
//
// @contact.name    API Support
// @contact.url     https://github.com/lwmacct/251117-go-ddd-template
//
// @license.name    MIT
// @license.url     https://opensource.org/licenses/MIT
//
// @host            localhost:8080
// @BasePath        /
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Bearer token authentication
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
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

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

	// 引入领域层包
	op "github.com/lwmacct/251117-go-ddd-template/internal/domain/operation"

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
	CreateLogHandler *auditlog.CreateHandler

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
	UserSettingHandler *handler.UserSettingHandler
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
	// OpenTelemetry 追踪中间件（如果启用）
	if cfg.Telemetry.Enabled {
		r.Use(otelgin.Middleware("go-ddd-template"))
	}

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

	// 声明式路由注册
	registerRoutes(r, deps)

	// 静态文件服务
	setupStaticRoutes(r, cfg)

	return r
}

// registerRoutes 自动注册所有路由
// 根据 Operation 元数据自动构建中间件链
func registerRoutes(r *gin.Engine, deps *RouterDependencies) {
	bindings := deps.AllRouteBindings()

	for _, b := range bindings {
		middlewares := buildMiddlewares(deps, b.Op)
		middlewares = append(middlewares, b.Handler)

		switch op.Method(b.Op) {
		case op.HttpGET:
			r.GET(op.Path(b.Op), middlewares...)
		case op.HttpPOST:
			r.POST(op.Path(b.Op), middlewares...)
		case op.HttpPUT:
			r.PUT(op.Path(b.Op), middlewares...)
		case op.HttpDELETE:
			r.DELETE(op.Path(b.Op), middlewares...)
		case op.HttpPATCH:
			r.PATCH(op.Path(b.Op), middlewares...)
		default:
			slog.Warn("unknown HTTP method", "operation", b.Op, "method", op.Method(b.Op))
		}
	}
}

// buildMiddlewares 根据 Operation 自动构建中间件链
// 中间件顺序：RequestID → OperationID → Auth → Permission → Audit
func buildMiddlewares(deps *RouterDependencies, o op.Operation) []gin.HandlerFunc {
	var mws []gin.HandlerFunc

	// 1. Request ID（所有请求）
	mws = append(mws, middleware.RequestID())

	// 2. Operation ID（所有请求）
	mws = append(mws, middleware.SetOperationID(o.String()))

	// 3. Auth + Permission（非公开操作需要认证和权限检查）
	if !op.IsPublic(o) {
		mws = append(mws, middleware.Auth(deps.JWTManager, deps.PATService, deps.PermissionCacheService))
		// 新 RBAC 模型：使用 Operation 本身作为权限标识符
		mws = append(mws, middleware.RequireOperation(o))
	}

	// 4. Audit（需要审计的操作）
	if op.NeedsAudit(o) {
		mws = append(mws, middleware.AuditMiddleware(deps.CreateLogHandler))
	}

	return mws
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
