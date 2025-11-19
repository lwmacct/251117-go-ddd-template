package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/auth"
)

// Auth 统一认证中间件 - 支持 JWT 和 PAT
func Auth(jwtManager *auth.JWTManager, patService *auth.PATService, userQueryRepo user.QueryRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取 Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "authorization header is required"})
			c.Abort()
			return
		}

		// 验证格式：Bearer <token>
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(401, gin.H{"error": "authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 判断 token 类型: PAT 以 "pat_" 开头
		if strings.HasPrefix(tokenString, "pat_") {
			// Personal Access Token 认证
			if err := authenticateWithPAT(c, patService, userQueryRepo, tokenString); err != nil{
				c.JSON(401, gin.H{"error": err.Error()})
				c.Abort()
				return
			}
		} else {
			// JWT 认证
			if err := authenticateWithJWT(c, jwtManager, tokenString); err != nil {
				c.JSON(401, gin.H{"error": err.Error()})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// JWTAuth JWT 认证中间件 (向后兼容)
func JWTAuth(jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(401, gin.H{"error": "authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		if err := authenticateWithJWT(c, jwtManager, tokenString); err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Next()
	}
}

// authenticateWithJWT 使用 JWT 进行认证
func authenticateWithJWT(c *gin.Context, jwtManager *auth.JWTManager, tokenString string) error {
	claims, err := jwtManager.ValidateToken(tokenString)
	if err != nil {
		return err
	}

	// 将用户信息存入上下文
	c.Set("user_id", claims.UserID)
	c.Set("username", claims.Username)
	c.Set("email", claims.Email)
	c.Set("roles", claims.Roles)
	c.Set("permissions", claims.Permissions)
	c.Set("auth_type", "jwt")

	return nil
}

// authenticateWithPAT 使用 Personal Access Token 进行认证
func authenticateWithPAT(c *gin.Context, patService *auth.PATService, userQueryRepo user.QueryRepository, tokenString string) error {
	// 验证 PAT (包含 IP 白名单检查)
	clientIP := c.ClientIP()
	pat, err := patService.ValidateTokenWithIP(c.Request.Context(), tokenString, clientIP)
	if err != nil {
		return err
	}

	// 获取用户信息 (包含角色)
	u, err := userQueryRepo.GetByIDWithRoles(c.Request.Context(), pat.UserID)
	if err != nil {
		return err
	}

	// 将用户信息存入上下文
	c.Set("user_id", u.ID)
	c.Set("username", u.Username)
	c.Set("email", u.Email)
	c.Set("roles", u.GetRoleNames())
	c.Set("permissions", pat.Permissions) // PAT 权限（用户权限的子集）
	c.Set("auth_type", "pat")
	c.Set("pat_id", pat.ID) // 额外存储 PAT ID，用于审计

	return nil
}
