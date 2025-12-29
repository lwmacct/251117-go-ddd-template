package middleware

import (
	"slices"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/response"
	op "github.com/lwmacct/251117-go-ddd-template/internal/domain/operation"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/role"
)

// RequireOperation 检查用户是否有执行指定 Operation 的权限。
// 新 RBAC 模型：权限为 Operation Pattern + Resource Pattern 组合。
// 对于 user: 域的操作，自动检查 user/self 资源匹配。
func RequireOperation(operation op.Operation) gin.HandlerFunc {
	return func(c *gin.Context) {
		permissions, exists := c.Get("permissions")
		if !exists {
			response.Unauthorized(c, "No permissions found")
			c.Abort()
			return
		}

		permList, ok := permissions.([]role.Permission)
		if !ok {
			response.InternalError(c, "Invalid permissions format")
			c.Abort()
			return
		}

		// 检查是否有匹配的权限
		hasPermission := false
		operationStr := string(operation)

		// 确定要检查的资源
		// 对于 user: 域的操作，检查 user/self 和 *
		resourcesToCheck := []string{"*"}
		if operation.Domain() == "user" {
			resourcesToCheck = appendUserResources(c, resourcesToCheck)
		}

		for _, p := range permList {
			// 匹配 Operation Pattern
			if op.MatchOperation(p.OperationPattern, operationStr) {
				// 检查资源是否匹配任一候选资源
				for _, res := range resourcesToCheck {
					if op.MatchResource(p.ResourcePattern, res) {
						hasPermission = true
						break
					}
				}
				if hasPermission {
					break
				}
			}
		}

		if !hasPermission {
			response.Forbidden(c, "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireOperationWithResource 检查用户是否有对指定资源执行指定 Operation 的权限。
// 支持细粒度资源控制，如 user/123、role/*。
func RequireOperationWithResource(operation op.Operation, resource op.Resource) gin.HandlerFunc {
	return func(c *gin.Context) {
		permissions, exists := c.Get("permissions")
		if !exists {
			response.Unauthorized(c, "No permissions found")
			c.Abort()
			return
		}

		permList, ok := permissions.([]role.Permission)
		if !ok {
			response.InternalError(c, "Invalid permissions format")
			c.Abort()
			return
		}

		hasPermission := false
		operationStr := string(operation)
		resourceStr := string(resource)
		for _, p := range permList {
			if op.MatchOperation(p.OperationPattern, operationStr) &&
				op.MatchResource(p.ResourcePattern, resourceStr) {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			response.Forbidden(c, "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireRole creates a middleware that checks if the user has a specific role
func RequireRole(roleName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, exists := c.Get("roles")
		if !exists {
			response.Unauthorized(c, "No roles found")
			c.Abort()
			return
		}

		rolesList, ok := roles.([]string)
		if !ok {
			response.InternalError(c, "Invalid roles format")
			c.Abort()
			return
		}

		if !slices.Contains(rolesList, roleName) {
			response.Forbidden(c, "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyRole creates a middleware that checks if the user has any of the specified roles
func RequireAnyRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRoles, exists := c.Get("roles")
		if !exists {
			response.Unauthorized(c, "No roles found")
			c.Abort()
			return
		}

		rolesList, ok := userRoles.([]string)
		if !ok {
			response.InternalError(c, "Invalid roles format")
			c.Abort()
			return
		}

		hasRole := slices.ContainsFunc(roles, func(requiredRole string) bool {
			return slices.Contains(rolesList, requiredRole)
		})

		if !hasRole {
			response.Forbidden(c, "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireOwnership creates a middleware that checks if the user is accessing their own resource
// The resource ID should be in the URL parameter specified by paramName (default: "id")
func RequireOwnership(paramName ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			response.Unauthorized(c, "No user ID found")
			c.Abort()
			return
		}

		uid, ok := userID.(uint)
		if !ok {
			response.InternalError(c, "Invalid user ID format")
			c.Abort()
			return
		}

		param := "id"
		if len(paramName) > 0 {
			param = paramName[0]
		}

		resourceIDStr := c.Param(param)
		resourceID, err := strconv.ParseUint(resourceIDStr, 10, 32)
		if err != nil {
			response.BadRequest(c, "Invalid resource ID")
			c.Abort()
			return
		}

		if uint(resourceID) != uid {
			response.Forbidden(c, "Can only access own resources")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAdminOrOwnership combines admin role check with ownership check
// Allows access if user is admin OR owns the resource
func RequireAdminOrOwnership(paramName ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if user is admin
		if isAdmin(c) {
			c.Next()
			return
		}

		// If not admin, check ownership
		userID, exists := c.Get("user_id")
		if !exists {
			response.Unauthorized(c, "No user ID found")
			c.Abort()
			return
		}

		uid, ok := userID.(uint)
		if !ok {
			response.InternalError(c, "Invalid user ID format")
			c.Abort()
			return
		}

		param := "id"
		if len(paramName) > 0 {
			param = paramName[0]
		}

		resourceIDStr := c.Param(param)
		resourceID, err := strconv.ParseUint(resourceIDStr, 10, 32)
		if err != nil {
			response.BadRequest(c, "Invalid resource ID")
			c.Abort()
			return
		}

		if uint(resourceID) != uid {
			response.Forbidden(c, "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

// isAdmin 检查当前用户是否具有 admin 角色
func isAdmin(c *gin.Context) bool {
	roles, exists := c.Get("roles")
	if !exists {
		return false
	}
	rolesList, ok := roles.([]string)
	if !ok {
		return false
	}
	return slices.Contains(rolesList, "admin")
}

// appendUserResources 为 user: 域操作添加 user/self 和 user/{id} 资源
func appendUserResources(c *gin.Context, resources []string) []string {
	userID, exists := c.Get("user_id")
	if !exists {
		return resources
	}
	uid, ok := userID.(uint)
	if !ok {
		return resources
	}
	return append(resources,
		"user/self",
		"user/"+strconv.FormatUint(uint64(uid), 10),
	)
}
