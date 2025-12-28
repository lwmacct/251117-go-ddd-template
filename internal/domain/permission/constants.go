package permission

// Admin 域 - 用户管理
const (
	AdminUsersCreate = "admin:users:create"
	AdminUsersRead   = "admin:users:read"
	AdminUsersUpdate = "admin:users:update"
	AdminUsersDelete = "admin:users:delete"
)

// Admin 域 - 角色管理
const (
	AdminRolesCreate = "admin:roles:create"
	AdminRolesRead   = "admin:roles:read"
	AdminRolesUpdate = "admin:roles:update"
	AdminRolesDelete = "admin:roles:delete"
)

// Admin 域 - 权限管理
const (
	AdminPermissionsRead = "admin:permissions:read"
)

// Admin 域 - 系统概览
const (
	AdminOverviewRead = "admin:overview:read"
)

// Admin 域 - 菜单管理
const (
	AdminMenusCreate = "admin:menus:create"
	AdminMenusRead   = "admin:menus:read"
	AdminMenusUpdate = "admin:menus:update"
	AdminMenusDelete = "admin:menus:delete"
)

// Admin 域 - 系统配置
const (
	AdminSettingsCreate = "admin:settings:create"
	AdminSettingsRead   = "admin:settings:read"
	AdminSettingsUpdate = "admin:settings:update"
	AdminSettingsDelete = "admin:settings:delete"
)

// Admin 域 - 审计日志
const (
	AdminAuditLogsRead = "admin:audit_logs:read"
)

// Admin 域 - 缓存管理
const (
	AdminCacheRead   = "admin:cache:read"
	AdminCacheDelete = "admin:cache:delete"
)

// User 域 - 个人资料
const (
	UserProfileRead   = "user:profile:read"
	UserProfileUpdate = "user:profile:update"
	UserProfileDelete = "user:profile:delete"
)

// User 域 - 密码管理
const (
	UserPasswordUpdate = "user:password:update"
)

// User 域 - 邮箱管理
const (
	UserEmailUpdate = "user:email:update"
)

// User 域 - 访问令牌管理
//
//nolint:gosec // G101: 这些是权限常量名，不是硬编码凭证
const (
	UserTokensCreate  = "user:tokens:create"
	UserTokensRead    = "user:tokens:read"
	UserTokensDisable = "user:tokens:disable"
	UserTokensEnable  = "user:tokens:enable"
	UserTokensDelete  = "user:tokens:delete"
)

// User 域 - 用户配置
const (
	UserSettingsRead   = "user:settings:read"
	UserSettingsUpdate = "user:settings:update"
)

// API 域 - 缓存操作
const (
	APICacheRead   = "api:cache:read"
	APICacheWrite  = "api:cache:write"
	APICacheDelete = "api:cache:delete"
)

// Definition 权限定义结构，用于种子数据生成。
type Definition struct {
	Code        string
	Domain      string
	Resource    string
	Action      string
	Description string
}

// AllDefinitions 返回所有权限定义，供 Seeder 使用。
func AllDefinitions() []Definition {
	return []Definition{
		// Admin 域 - 用户管理
		{AdminUsersCreate, "admin", "users", "create", "Create users"},
		{AdminUsersRead, "admin", "users", "read", "Read all users"},
		{AdminUsersUpdate, "admin", "users", "update", "Update any user"},
		{AdminUsersDelete, "admin", "users", "delete", "Delete users"},

		// Admin 域 - 角色管理
		{AdminRolesCreate, "admin", "roles", "create", "Create roles"},
		{AdminRolesRead, "admin", "roles", "read", "Read all roles"},
		{AdminRolesUpdate, "admin", "roles", "update", "Update roles"},
		{AdminRolesDelete, "admin", "roles", "delete", "Delete roles"},

		// Admin 域 - 权限管理
		{AdminPermissionsRead, "admin", "permissions", "read", "Read all permissions"},

		// Admin 域 - 系统概览
		{AdminOverviewRead, "admin", "overview", "read", "View system overview stats"},

		// Admin 域 - 菜单管理
		{AdminMenusCreate, "admin", "menus", "create", "Create menus"},
		{AdminMenusRead, "admin", "menus", "read", "Read menus"},
		{AdminMenusUpdate, "admin", "menus", "update", "Update menus"},
		{AdminMenusDelete, "admin", "menus", "delete", "Delete menus"},

		// Admin 域 - 系统配置
		{AdminSettingsCreate, "admin", "settings", "create", "Create settings"},
		{AdminSettingsRead, "admin", "settings", "read", "Read settings"},
		{AdminSettingsUpdate, "admin", "settings", "update", "Update settings"},
		{AdminSettingsDelete, "admin", "settings", "delete", "Delete settings"},

		// Admin 域 - 审计日志
		{AdminAuditLogsRead, "admin", "audit_logs", "read", "Read audit logs"},

		// Admin 域 - 缓存管理
		{AdminCacheRead, "admin", "cache", "read", "Read cache status and keys"},
		{AdminCacheDelete, "admin", "cache", "delete", "Delete cache keys"},

		// User 域 - 个人资料
		{UserProfileRead, "user", "profile", "read", "Read own profile"},
		{UserProfileUpdate, "user", "profile", "update", "Update own profile"},
		{UserProfileDelete, "user", "profile", "delete", "Delete own account"},

		// User 域 - 密码管理
		{UserPasswordUpdate, "user", "password", "update", "Change own password"},

		// User 域 - 邮箱管理
		{UserEmailUpdate, "user", "email", "update", "Change own email"},

		// User 域 - 访问令牌管理
		{UserTokensCreate, "user", "tokens", "create", "Create personal access tokens"},
		{UserTokensRead, "user", "tokens", "read", "List own tokens"},
		{UserTokensDisable, "user", "tokens", "update", "Disable own tokens"},
		{UserTokensEnable, "user", "tokens", "update", "Enable own tokens"},
		{UserTokensDelete, "user", "tokens", "delete", "Delete own tokens"},

		// User 域 - 用户配置
		{UserSettingsRead, "user", "settings", "read", "Read own user settings"},
		{UserSettingsUpdate, "user", "settings", "update", "Update own user settings"},

		// API 域 - 缓存操作
		{APICacheRead, "api", "cache", "read", "Read cache data"},
		{APICacheWrite, "api", "cache", "write", "Write cache data"},
		{APICacheDelete, "api", "cache", "delete", "Delete cache data"},
	}
}
