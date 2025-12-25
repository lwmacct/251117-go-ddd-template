// Package middleware 提供 HTTP 中间件。
//
// 本包实现了 Gin 框架的中间件，用于请求处理管道：
//
// 认证中间件：
//   - Auth: 统一认证（支持 JWT 和 PAT 双模式）
//   - JWTAuth: 仅 JWT 认证（已废弃，保留向后兼容）
//
// 授权中间件：
//   - RequireRole: 角色检查（如 RequireRole("admin")）
//   - RequirePermission: 权限检查（如 RequirePermission("admin:users:read")）
//
// 通用中间件：
//   - CORS: 跨域资源共享配置
//   - Logger: 基于 slog 的请求日志
//   - AuditMiddleware: 审计日志记录
//
// 权限缓存机制：
// 新架构中，JWT/PAT 仅存储 user_id，权限信息从 PermissionCacheService
// 实时查询，支持权限变更后立即生效。
package middleware

import (
	"context"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/response"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/auth"
)

// Auth 统一认证中间件 - 支持 JWT 和 PAT
// 新架构：权限信息统一从 PermissionCacheService 查询
func Auth(jwtManager *auth.JWTManager, patService *auth.PATService, permCacheService *auth.PermissionCacheService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取 Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "Authorization header is required")
			c.Abort()
			return
		}

		// 验证格式：Bearer <token>
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "Authorization header format must be Bearer {token}")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 判断 token 类型: PAT 以 "pat_" 开头，使用相应的认证方式
		ctx := c.Request.Context()
		var authErr error
		if strings.HasPrefix(tokenString, "pat_") {
			authErr = authenticateWithPAT(ctx, c, patService, permCacheService, tokenString)
		} else {
			authErr = authenticateWithJWT(ctx, c, jwtManager, permCacheService, tokenString)
		}

		if authErr != nil {
			response.Unauthorized(c, authErr.Error())
			c.Abort()
			return
		}

		c.Next()
	}
}

// JWTAuth JWT 认证中间件 (向后兼容，已废弃)
//
// Deprecated: 使用 Auth(jwtManager, patService, permissionCacheService) 代替
func JWTAuth(jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "Authorization header is required")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "Authorization header format must be Bearer {token}")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 旧方法：仅从 token 中读取权限（不支持实时权限查询）
		claims, err := jwtManager.ValidateToken(tokenString)
		if err != nil {
			response.Unauthorized(c, err.Error())
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("roles", claims.Roles)
		c.Set("permissions", claims.Permissions)
		c.Set("auth_type", "jwt")

		c.Next()
	}
}

// authenticateWithJWT 使用 JWT 进行认证
// 新架构：从 token 获取 user_id，权限信息从缓存实时查询
func authenticateWithJWT(ctx context.Context, c *gin.Context, jwtManager *auth.JWTManager, permCacheService *auth.PermissionCacheService, tokenString string) error {
	claims, err := jwtManager.ValidateToken(tokenString)
	if err != nil {
		return err
	}

	// 从缓存查询权限信息（向后兼容：优先使用 token 中的权限，如果为空则查询缓存）
	var roles, permissions []string
	if len(claims.Roles) > 0 || len(claims.Permissions) > 0 {
		// 旧 token 包含权限信息，直接使用（向后兼容）
		roles = claims.Roles
		permissions = claims.Permissions
	} else {
		// 新 token 不包含权限信息，从缓存查询
		roles, permissions, err = permCacheService.GetUserPermissions(ctx, claims.UserID)
		if err != nil {
			return fmt.Errorf("failed to get user permissions: %w", err)
		}
	}

	// 将用户信息存入上下文
	c.Set("user_id", claims.UserID)
	c.Set("username", claims.Username)
	c.Set("email", claims.Email)
	c.Set("roles", roles)
	c.Set("permissions", permissions)
	c.Set("auth_type", "jwt")

	return nil
}

// authenticateWithPAT 使用 Personal Access Token 进行认证
// 新架构：PAT 自动继承用户全部权限，从缓存实时查询
func authenticateWithPAT(ctx context.Context, c *gin.Context, patService *auth.PATService, permCacheService *auth.PermissionCacheService, tokenString string) error {
	// 验证 PAT (包含 IP 白名单检查)
	clientIP := c.ClientIP()
	pat, err := patService.ValidateTokenWithIP(ctx, tokenString, clientIP)
	if err != nil {
		return err
	}

	// 从缓存查询用户权限（PAT 自动继承用户全部权限）
	roles, permissions, err := permCacheService.GetUserPermissions(ctx, pat.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user permissions: %w", err)
	}

	// 将用户信息存入上下文
	c.Set("user_id", pat.UserID)
	c.Set("username", "") // PAT 不存储 username，可从用户表查询
	c.Set("email", "")
	c.Set("roles", roles)
	c.Set("permissions", permissions) // 从缓存查询的完整权限
	c.Set("auth_type", "pat")
	c.Set("pat_id", pat.ID) // 额外存储 PAT ID，用于审计

	return nil
}
