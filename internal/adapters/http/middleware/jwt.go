package middleware

import (
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

		// 判断 token 类型: PAT 以 "pat_" 开头
		if strings.HasPrefix(tokenString, "pat_") {
			// Personal Access Token 认证
			if err := authenticateWithPAT(c, patService, permCacheService, tokenString); err != nil {
				response.Unauthorized(c, err.Error())
				c.Abort()
				return
			}
		} else {
			// JWT 认证
			if err := authenticateWithJWT(c, jwtManager, permCacheService, tokenString); err != nil {
				response.Unauthorized(c, err.Error())
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// JWTAuth JWT 认证中间件 (向后兼容，已废弃)
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
func authenticateWithJWT(c *gin.Context, jwtManager *auth.JWTManager, permCacheService *auth.PermissionCacheService, tokenString string) error {
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
		roles, permissions, err = permCacheService.GetUserPermissions(c.Request.Context(), claims.UserID)
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
func authenticateWithPAT(c *gin.Context, patService *auth.PATService, permCacheService *auth.PermissionCacheService, tokenString string) error {
	// 验证 PAT (包含 IP 白名单检查)
	clientIP := c.ClientIP()
	pat, err := patService.ValidateTokenWithIP(c.Request.Context(), tokenString, clientIP)
	if err != nil {
		return err
	}

	// 从缓存查询用户权限（PAT 自动继承用户全部权限）
	roles, permissions, err := permCacheService.GetUserPermissions(c.Request.Context(), pat.UserID)
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
