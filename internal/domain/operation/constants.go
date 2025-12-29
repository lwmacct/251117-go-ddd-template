package operation

// ============================================================================
// 操作常量
// ============================================================================

// Auth 域（公开接口）
const (
	AuthRegister Operation = "auth.register"
	AuthLogin    Operation = "auth.login"
	AuthLogin2FA Operation = "auth.login2fa"
	AuthRefresh  Operation = "auth.refresh"
	AuthCaptcha  Operation = "auth.captcha"
)

// Auth 域 - 2FA（需认证）
const (
	Auth2FASetup   Operation = "auth.2fa.setup"
	Auth2FAVerify  Operation = "auth.2fa.verify"
	Auth2FADisable Operation = "auth.2fa.disable"
	Auth2FAStatus  Operation = "auth.2fa.status"
)

// Admin 域 - 用户管理
const (
	AdminUsersCreate      Operation = "admin.users.create"
	AdminUsersBatchCreate Operation = "admin.users.batchCreate"
	AdminUsersList        Operation = "admin.users.list"
	AdminUsersGet         Operation = "admin.users.get"
	AdminUsersUpdate      Operation = "admin.users.update"
	AdminUsersDelete      Operation = "admin.users.delete"
	AdminUsersAssignRoles Operation = "admin.users.assignRoles"
)

// Admin 域 - 角色管理
const (
	AdminRolesCreate         Operation = "admin.roles.create"
	AdminRolesList           Operation = "admin.roles.list"
	AdminRolesGet            Operation = "admin.roles.get"
	AdminRolesUpdate         Operation = "admin.roles.update"
	AdminRolesDelete         Operation = "admin.roles.delete"
	AdminRolesSetPermissions Operation = "admin.roles.setPermissions"
)

// Admin 域 - 权限管理
const (
	AdminPermissionsList Operation = "admin.permissions.list"
)

// Admin 域 - 审计日志
const (
	AdminAuditLogsList    Operation = "admin.auditlogs.list"
	AdminAuditLogsGet     Operation = "admin.auditlogs.get"
	AdminAuditLogsActions Operation = "admin.auditlogs.actions"
)

// Admin 域 - 菜单管理
const (
	AdminMenusCreate  Operation = "admin.menus.create"
	AdminMenusList    Operation = "admin.menus.list"
	AdminMenusGet     Operation = "admin.menus.get"
	AdminMenusUpdate  Operation = "admin.menus.update"
	AdminMenusDelete  Operation = "admin.menus.delete"
	AdminMenusReorder Operation = "admin.menus.reorder"
)

// Admin 域 - 系统概览
const (
	AdminOverviewStats Operation = "admin.overview.stats"
)

// Admin 域 - 系统配置
const (
	AdminSettingsCreate      Operation = "admin.settings.create"
	AdminSettingsList        Operation = "admin.settings.list"
	AdminSettingsGet         Operation = "admin.settings.get"
	AdminSettingsUpdate      Operation = "admin.settings.update"
	AdminSettingsDelete      Operation = "admin.settings.delete"
	AdminSettingsBatchUpdate Operation = "admin.settings.batchUpdate"
)

// Admin 域 - 配置分类
const (
	AdminSettingCategoriesList   Operation = "admin.settings.categories.list"
	AdminSettingCategoriesGet    Operation = "admin.settings.categories.get"
	AdminSettingCategoriesCreate Operation = "admin.settings.categories.create"
	AdminSettingCategoriesUpdate Operation = "admin.settings.categories.update"
	AdminSettingCategoriesDelete Operation = "admin.settings.categories.delete"
)

// Admin 域 - 缓存管理
const (
	AdminCacheInfo          Operation = "admin.cache.info"
	AdminCacheScanKeys      Operation = "admin.cache.scanKeys"
	AdminCacheGetKey        Operation = "admin.cache.getKey"
	AdminCacheDeleteKey     Operation = "admin.cache.deleteKey"
	AdminCacheDeletePattern Operation = "admin.cache.deletePattern"
)

// User 域 - 个人资料
const (
	UserProfileGet     Operation = "user.profile.get"
	UserProfileUpdate  Operation = "user.profile.update"
	UserPasswordUpdate Operation = "user.password.update"
	UserAccountDelete  Operation = "user.account.delete"
)

// User 域 - 访问令牌
//
//nolint:gosec // G101: 这些是操作标识符，非硬编码凭证
const (
	UserTokensCreate  Operation = "user.tokens.create"
	UserTokensList    Operation = "user.tokens.list"
	UserTokensGet     Operation = "user.tokens.get"
	UserTokensDelete  Operation = "user.tokens.delete"
	UserTokensDisable Operation = "user.tokens.disable"
	UserTokensEnable  Operation = "user.tokens.enable"
)

// User 域 - 用户配置
const (
	UserSettingsCategoriesList Operation = "user.settings.categories.list"
	UserSettingsList           Operation = "user.settings.list"
	UserSettingsGet            Operation = "user.settings.get"
	UserSettingsSet            Operation = "user.settings.set"
	UserSettingsReset          Operation = "user.settings.reset"
	UserSettingsBatchSet       Operation = "user.settings.batchSet"
)
