package registry

import (
	"strings"
	"sync"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/permission"
)

// Endpoint API 端点定义。
type Endpoint struct {
	// OperationID 唯一操作标识符，格式：{domain}.{resource}.{action}
	// 例如：admin.users.create, user.profile.get
	OperationID string

	// Method HTTP 方法
	Method string

	// Path 路由路径（Gin 格式，如 /api/admin/users/:id）
	Path string

	// Permission 所需权限（可选，为空则不检查权限）
	Permission string

	// Description 端点描述
	Description string
}

// endpoints 存储所有端点定义
var endpoints = []Endpoint{
	// ==================== Auth API ====================
	{OperationID: "auth.register", Method: "POST", Path: "/api/auth/register", Permission: "", Description: "User registration"},
	{OperationID: "auth.login", Method: "POST", Path: "/api/auth/login", Permission: "", Description: "User login"},
	{OperationID: "auth.login2fa", Method: "POST", Path: "/api/auth/login/2fa", Permission: "", Description: "Two-factor authentication login"},
	{OperationID: "auth.refresh", Method: "POST", Path: "/api/auth/refresh", Permission: "", Description: "Refresh access token"},
	{OperationID: "auth.captcha", Method: "GET", Path: "/api/auth/captcha", Permission: "", Description: "Get captcha image"},

	// ==================== 2FA API ====================
	{OperationID: "auth.2fa.setup", Method: "POST", Path: "/api/auth/2fa/setup", Permission: "", Description: "Setup 2FA"},
	{OperationID: "auth.2fa.verify", Method: "POST", Path: "/api/auth/2fa/verify", Permission: "", Description: "Verify and enable 2FA"},
	{OperationID: "auth.2fa.disable", Method: "POST", Path: "/api/auth/2fa/disable", Permission: "", Description: "Disable 2FA"},
	{OperationID: "auth.2fa.status", Method: "GET", Path: "/api/auth/2fa/status", Permission: "", Description: "Get 2FA status"},

	// ==================== Admin - Users ====================
	{OperationID: "admin.users.create", Method: "POST", Path: "/api/admin/users", Permission: permission.AdminUsersCreate, Description: "Create user"},
	{OperationID: "admin.users.batchCreate", Method: "POST", Path: "/api/admin/users/batch", Permission: permission.AdminUsersCreate, Description: "Batch create users"},
	{OperationID: "admin.users.list", Method: "GET", Path: "/api/admin/users", Permission: permission.AdminUsersRead, Description: "List users"},
	{OperationID: "admin.users.get", Method: "GET", Path: "/api/admin/users/:id", Permission: permission.AdminUsersRead, Description: "Get user by ID"},
	{OperationID: "admin.users.update", Method: "PUT", Path: "/api/admin/users/:id", Permission: permission.AdminUsersUpdate, Description: "Update user"},
	{OperationID: "admin.users.delete", Method: "DELETE", Path: "/api/admin/users/:id", Permission: permission.AdminUsersDelete, Description: "Delete user"},
	{OperationID: "admin.users.assignRoles", Method: "PUT", Path: "/api/admin/users/:id/roles", Permission: permission.AdminUsersUpdate, Description: "Assign roles to user"},

	// ==================== Admin - Roles ====================
	{OperationID: "admin.roles.create", Method: "POST", Path: "/api/admin/roles", Permission: permission.AdminRolesCreate, Description: "Create role"},
	{OperationID: "admin.roles.list", Method: "GET", Path: "/api/admin/roles", Permission: permission.AdminRolesRead, Description: "List roles"},
	{OperationID: "admin.roles.get", Method: "GET", Path: "/api/admin/roles/:id", Permission: permission.AdminRolesRead, Description: "Get role by ID"},
	{OperationID: "admin.roles.update", Method: "PUT", Path: "/api/admin/roles/:id", Permission: permission.AdminRolesUpdate, Description: "Update role"},
	{OperationID: "admin.roles.delete", Method: "DELETE", Path: "/api/admin/roles/:id", Permission: permission.AdminRolesDelete, Description: "Delete role"},
	{OperationID: "admin.roles.setPermissions", Method: "PUT", Path: "/api/admin/roles/:id/permissions", Permission: permission.AdminRolesUpdate, Description: "Set role permissions"},

	// ==================== Admin - Permissions ====================
	{OperationID: "admin.permissions.list", Method: "GET", Path: "/api/admin/permissions", Permission: permission.AdminPermissionsRead, Description: "List permissions"},

	// ==================== Admin - Audit Logs ====================
	{OperationID: "admin.auditlogs.list", Method: "GET", Path: "/api/admin/auditlogs", Permission: permission.AdminAuditLogsRead, Description: "List audit logs"},
	{OperationID: "admin.auditlogs.get", Method: "GET", Path: "/api/admin/auditlogs/:id", Permission: permission.AdminAuditLogsRead, Description: "Get audit log by ID"},

	// ==================== Admin - Menus ====================
	{OperationID: "admin.menus.create", Method: "POST", Path: "/api/admin/menus", Permission: permission.AdminMenusCreate, Description: "Create menu"},
	{OperationID: "admin.menus.list", Method: "GET", Path: "/api/admin/menus", Permission: permission.AdminMenusRead, Description: "List menus"},
	{OperationID: "admin.menus.get", Method: "GET", Path: "/api/admin/menus/:id", Permission: permission.AdminMenusRead, Description: "Get menu by ID"},
	{OperationID: "admin.menus.update", Method: "PUT", Path: "/api/admin/menus/:id", Permission: permission.AdminMenusUpdate, Description: "Update menu"},
	{OperationID: "admin.menus.delete", Method: "DELETE", Path: "/api/admin/menus/:id", Permission: permission.AdminMenusDelete, Description: "Delete menu"},
	{OperationID: "admin.menus.reorder", Method: "POST", Path: "/api/admin/menus/reorder", Permission: permission.AdminMenusUpdate, Description: "Reorder menus"},

	// ==================== Admin - Overview ====================
	{OperationID: "admin.overview.stats", Method: "GET", Path: "/api/admin/overview/stats", Permission: permission.AdminOverviewRead, Description: "Get system overview stats"},

	// ==================== Admin - Settings ====================
	{OperationID: "admin.settings.list", Method: "GET", Path: "/api/admin/settings", Permission: permission.AdminSettingsRead, Description: "List settings"},
	{OperationID: "admin.settings.create", Method: "POST", Path: "/api/admin/settings", Permission: permission.AdminSettingsCreate, Description: "Create setting"},
	{OperationID: "admin.settings.batchUpdate", Method: "POST", Path: "/api/admin/settings/batch", Permission: permission.AdminSettingsUpdate, Description: "Batch update settings"},
	{OperationID: "admin.settings.get", Method: "GET", Path: "/api/admin/settings/:key", Permission: permission.AdminSettingsRead, Description: "Get setting by key"},
	{OperationID: "admin.settings.update", Method: "PUT", Path: "/api/admin/settings/:key", Permission: permission.AdminSettingsUpdate, Description: "Update setting"},
	{OperationID: "admin.settings.delete", Method: "DELETE", Path: "/api/admin/settings/:key", Permission: permission.AdminSettingsDelete, Description: "Delete setting"},

	// ==================== Admin - Setting Categories ====================
	{OperationID: "admin.settings.categories.list", Method: "GET", Path: "/api/admin/settings/categories", Permission: permission.AdminSettingsRead, Description: "List setting categories"},
	{OperationID: "admin.settings.categories.get", Method: "GET", Path: "/api/admin/settings/categories/:id", Permission: permission.AdminSettingsRead, Description: "Get setting category by ID"},
	{OperationID: "admin.settings.categories.create", Method: "POST", Path: "/api/admin/settings/categories", Permission: permission.AdminSettingsCreate, Description: "Create setting category"},
	{OperationID: "admin.settings.categories.update", Method: "PUT", Path: "/api/admin/settings/categories/:id", Permission: permission.AdminSettingsUpdate, Description: "Update setting category"},
	{OperationID: "admin.settings.categories.delete", Method: "DELETE", Path: "/api/admin/settings/categories/:id", Permission: permission.AdminSettingsDelete, Description: "Delete setting category"},

	// ==================== Admin - Cache ====================
	{OperationID: "admin.cache.info", Method: "GET", Path: "/api/admin/cache/info", Permission: permission.AdminCacheRead, Description: "Get cache info"},
	{OperationID: "admin.cache.scanKeys", Method: "GET", Path: "/api/admin/cache/keys", Permission: permission.AdminCacheRead, Description: "Scan cache keys"},
	{OperationID: "admin.cache.getKey", Method: "GET", Path: "/api/admin/cache/key", Permission: permission.AdminCacheRead, Description: "Get cache key value"},
	{OperationID: "admin.cache.deleteKey", Method: "DELETE", Path: "/api/admin/cache/key", Permission: permission.AdminCacheDelete, Description: "Delete cache key"},
	{OperationID: "admin.cache.deleteByPattern", Method: "DELETE", Path: "/api/admin/cache/keys", Permission: permission.AdminCacheDelete, Description: "Delete cache keys by pattern"},

	// ==================== User - Profile ====================
	{OperationID: "user.profile.get", Method: "GET", Path: "/api/user/profile", Permission: permission.UserProfileRead, Description: "Get current user profile"},
	{OperationID: "user.profile.update", Method: "PUT", Path: "/api/user/profile", Permission: permission.UserProfileUpdate, Description: "Update current user profile"},
	{OperationID: "user.password.update", Method: "PUT", Path: "/api/user/password", Permission: permission.UserPasswordUpdate, Description: "Change password"},
	{OperationID: "user.account.delete", Method: "DELETE", Path: "/api/user/account", Permission: permission.UserProfileDelete, Description: "Delete own account"},

	// ==================== User - Tokens ====================
	{OperationID: "user.tokens.create", Method: "POST", Path: "/api/user/tokens", Permission: permission.UserTokensCreate, Description: "Create personal access token"},
	{OperationID: "user.tokens.list", Method: "GET", Path: "/api/user/tokens", Permission: permission.UserTokensRead, Description: "List personal access tokens"},
	{OperationID: "user.tokens.get", Method: "GET", Path: "/api/user/tokens/:id", Permission: permission.UserTokensRead, Description: "Get token by ID"},
	{OperationID: "user.tokens.delete", Method: "DELETE", Path: "/api/user/tokens/:id", Permission: permission.UserTokensDelete, Description: "Delete token"},
	{OperationID: "user.tokens.disable", Method: "PATCH", Path: "/api/user/tokens/:id/disable", Permission: permission.UserTokensDisable, Description: "Disable token"},
	{OperationID: "user.tokens.enable", Method: "PATCH", Path: "/api/user/tokens/:id/enable", Permission: permission.UserTokensEnable, Description: "Enable token"},

	// ==================== User - Settings ====================
	{OperationID: "user.settings.categories.list", Method: "GET", Path: "/api/user/settings/categories", Permission: permission.UserSettingsRead, Description: "List user setting categories"},
	{OperationID: "user.settings.list", Method: "GET", Path: "/api/user/settings", Permission: permission.UserSettingsRead, Description: "Get all user settings"},
	{OperationID: "user.settings.get", Method: "GET", Path: "/api/user/settings/:key", Permission: permission.UserSettingsRead, Description: "Get user setting by key"},
	{OperationID: "user.settings.set", Method: "PUT", Path: "/api/user/settings/:key", Permission: permission.UserSettingsUpdate, Description: "Set user setting"},
	{OperationID: "user.settings.reset", Method: "DELETE", Path: "/api/user/settings/:key", Permission: permission.UserSettingsUpdate, Description: "Reset user setting to default"},
	{OperationID: "user.settings.batchSet", Method: "POST", Path: "/api/user/settings/batch", Permission: permission.UserSettingsUpdate, Description: "Batch set user settings"},
}

// pathIndex 路径索引，用于快速查找
var pathIndex map[string]*Endpoint
var pathIndexOnce sync.Once

// buildPathIndex 构建路径索引
func buildPathIndex() {
	pathIndex = make(map[string]*Endpoint, len(endpoints))
	for i := range endpoints {
		key := endpoints[i].Method + " " + endpoints[i].Path
		pathIndex[key] = &endpoints[i]
	}
}

// ByPath 通过 HTTP 方法和路径查找端点。
// path 应为 Gin 路由模式（如 /api/admin/users/:id）。
// 如果未找到返回 nil。
func ByPath(method, path string) *Endpoint {
	pathIndexOnce.Do(buildPathIndex)
	return pathIndex[method+" "+path]
}

// ByOperationID 通过 Operation ID 查找端点。
// 如果未找到返回 nil。
func ByOperationID(operationID string) *Endpoint {
	for i := range endpoints {
		if endpoints[i].OperationID == operationID {
			return &endpoints[i]
		}
	}
	return nil
}

// All 返回所有端点定义的副本。
func All() []Endpoint {
	result := make([]Endpoint, len(endpoints))
	copy(result, endpoints)
	return result
}

// NormalizePath 将实际请求路径转换为 Gin 路由模式。
// 例如：/api/admin/users/123 -> /api/admin/users/:id
func NormalizePath(path string) string {
	parts := strings.Split(path, "/")
	for i, part := range parts {
		// 如果是数字，替换为 :id
		if isNumeric(part) {
			parts[i] = ":id"
		}
		// 如果是 UUID 格式，替换为 :id
		if isUUID(part) {
			parts[i] = ":id"
		}
	}
	return strings.Join(parts, "/")
}

// isNumeric 检查字符串是否全是数字
func isNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// isUUID 检查字符串是否是 UUID 格式
func isUUID(s string) bool {
	// UUID 格式：xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	if len(s) != 36 {
		return false
	}
	for i, c := range s {
		switch i {
		case 8, 13, 18, 23:
			if c != '-' {
				return false
			}
		default:
			if !isHexDigit(c) {
				return false
			}
		}
	}
	return true
}

// isHexDigit 检查字符是否是十六进制数字
func isHexDigit(c rune) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}
