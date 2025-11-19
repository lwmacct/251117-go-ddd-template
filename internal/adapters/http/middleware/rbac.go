package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// RequireRole creates a middleware that checks if the user has a specific role
func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, exists := c.Get("roles")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized: no roles found",
			})
			c.Abort()
			return
		}

		rolesList, ok := roles.([]string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal error: invalid roles format",
			})
			c.Abort()
			return
		}

		hasRole := false
		for _, r := range rolesList {
			if r == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "forbidden: insufficient permissions",
			})
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
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized: no roles found",
			})
			c.Abort()
			return
		}

		rolesList, ok := userRoles.([]string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal error: invalid roles format",
			})
			c.Abort()
			return
		}

		hasRole := false
		for _, requiredRole := range roles {
			for _, userRole := range rolesList {
				if userRole == requiredRole {
					hasRole = true
					break
				}
			}
			if hasRole {
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "forbidden: insufficient permissions",
			})
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
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized: no permissions found",
			})
			c.Abort()
			return
		}

		permissionsList, ok := permissions.([]string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal error: invalid permissions format",
			})
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
			c.JSON(http.StatusForbidden, gin.H{
				"error": "forbidden: insufficient permissions",
			})
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
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized: no permissions found",
			})
			c.Abort()
			return
		}

		permissionsList, ok := userPermissions.([]string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal error: invalid permissions format",
			})
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
			c.JSON(http.StatusForbidden, gin.H{
				"error": "forbidden: insufficient permissions",
			})
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
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized: no user ID found",
			})
			c.Abort()
			return
		}

		uid, ok := userID.(uint)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal error: invalid user ID format",
			})
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
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid resource ID",
			})
			c.Abort()
			return
		}

		if uint(resourceID) != uid {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "forbidden: can only access own resources",
			})
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
		roles, exists := c.Get("roles")
		if exists {
			if rolesList, ok := roles.([]string); ok {
				for _, role := range rolesList {
					if role == "admin" {
						c.Next()
						return
					}
				}
			}
		}

		// If not admin, check ownership
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized: no user ID found",
			})
			c.Abort()
			return
		}

		uid, ok := userID.(uint)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal error: invalid user ID format",
			})
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
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid resource ID",
			})
			c.Abort()
			return
		}

		if uint(resourceID) != uid {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "forbidden: insufficient permissions",
			})
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
	for i := 0; i < 3; i++ {
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
