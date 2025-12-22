package middleware

import (
	"slices"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/response"
)

// RequireRole creates a middleware that checks if the user has a specific role
func RequireRole(role string) gin.HandlerFunc {
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

		if !slices.Contains(rolesList, role) {
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

// RequirePermission creates a middleware that checks if the user has a specific permission
// Supports three-part permission format: domain:resource:action
// Also supports wildcard matching: admin:users:*, admin:*:create, *:*:*
func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		permissions, exists := c.Get("permissions")
		if !exists {
			response.Unauthorized(c, "No permissions found")
			c.Abort()
			return
		}

		permissionsList, ok := permissions.([]string)
		if !ok {
			response.InternalError(c, "Invalid permissions format")
			c.Abort()
			return
		}

		hasPermission := false
		for _, p := range permissionsList {
			if matchPermission(p, permission) {
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

// RequireAnyPermission creates a middleware that checks if the user has any of the specified permissions
// Supports wildcard matching for three-part permission format
func RequireAnyPermission(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userPermissions, exists := c.Get("permissions")
		if !exists {
			response.Unauthorized(c, "No permissions found")
			c.Abort()
			return
		}

		permissionsList, ok := userPermissions.([]string)
		if !ok {
			response.InternalError(c, "Invalid permissions format")
			c.Abort()
			return
		}

		hasPermission := false
		for _, requiredPermission := range permissions {
			for _, userPermission := range permissionsList {
				if matchPermission(userPermission, requiredPermission) {
					hasPermission = true
					break
				}
			}
			if hasPermission {
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

// matchPermission checks if a user permission matches a required permission
// Supports three-part wildcard matching: domain:resource:action
//
// Examples:
//   - matchPermission("admin:users:create", "admin:users:create") -> true (exact match)
//   - matchPermission("admin:users:*", "admin:users:create") -> true (action wildcard)
//   - matchPermission("admin:*:create", "admin:users:create") -> true (resource wildcard)
//   - matchPermission("*:users:create", "admin:users:create") -> true (domain wildcard)
//   - matchPermission("admin:*:*", "admin:users:create") -> true (all admin permissions)
//   - matchPermission("*:*:*", "admin:users:create") -> true (super admin)
func matchPermission(userPerm, requiredPerm string) bool {
	// Exact match
	if userPerm == requiredPerm {
		return true
	}

	// Split both permissions into parts
	userParts := strings.Split(userPerm, ":")
	requiredParts := strings.Split(requiredPerm, ":")

	// Must be three-part format
	if len(userParts) != 3 || len(requiredParts) != 3 {
		return userPerm == requiredPerm // Fallback to exact match for non-standard format
	}

	// Check each part: domain, resource, action
	for i := range 3 {
		// Wildcard in user permission matches anything
		if userParts[i] == "*" {
			continue
		}

		// Parts must match exactly
		if userParts[i] != requiredParts[i] {
			return false
		}
	}

	return true
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
