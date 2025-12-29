package operation

// ============================================================================
// 操作常量
// ============================================================================
//
// 格式：{domain}:{module}.{action}
//
// 域划分：
//   - auth: 认证域（公开访问）
//   - sys:  系统管理域（需管理员权限）
//   - user: 用户自服务域（当前用户权限）

// Auth 域（公开接口）
const (
	AuthRegister Operation = "auth:register"
	AuthLogin    Operation = "auth:login"
	AuthLogin2FA Operation = "auth:login2fa"
	AuthRefresh  Operation = "auth:refresh"
	AuthCaptcha  Operation = "auth:captcha"
)

// Auth 域 - 2FA（需认证）
const (
	Auth2FASetup   Operation = "auth:2fa.setup"
	Auth2FAVerify  Operation = "auth:2fa.verify"
	Auth2FADisable Operation = "auth:2fa.disable"
	Auth2FAStatus  Operation = "auth:2fa.status"
)

// Sys 域 - 用户管理
const (
	SysUsersCreate      Operation = "sys:users.create"
	SysUsersBatchCreate Operation = "sys:users.batchCreate"
	SysUsersList        Operation = "sys:users.list"
	SysUsersGet         Operation = "sys:users.get"
	SysUsersUpdate      Operation = "sys:users.update"
	SysUsersDelete      Operation = "sys:users.delete"
	SysUsersAssignRoles Operation = "sys:users.assignRoles"
)

// Sys 域 - 角色管理
const (
	SysRolesCreate         Operation = "sys:roles.create"
	SysRolesList           Operation = "sys:roles.list"
	SysRolesGet            Operation = "sys:roles.get"
	SysRolesUpdate         Operation = "sys:roles.update"
	SysRolesDelete         Operation = "sys:roles.delete"
	SysRolesSetPermissions Operation = "sys:roles.setPermissions"
)

// Sys 域 - 操作列表（供前端权限配置使用）
const (
	SysOperationsList Operation = "sys:operations.list"
)

// Sys 域 - 审计日志
const (
	SysAuditLogsList    Operation = "sys:auditlogs.list"
	SysAuditLogsGet     Operation = "sys:auditlogs.get"
	SysAuditLogsActions Operation = "sys:auditlogs.actions"
)

// Sys 域 - 菜单管理
const (
	SysMenusCreate  Operation = "sys:menus.create"
	SysMenusList    Operation = "sys:menus.list"
	SysMenusGet     Operation = "sys:menus.get"
	SysMenusUpdate  Operation = "sys:menus.update"
	SysMenusDelete  Operation = "sys:menus.delete"
	SysMenusReorder Operation = "sys:menus.reorder"
)

// Sys 域 - 系统概览
const (
	SysOverviewStats Operation = "sys:overview.stats"
)

// Sys 域 - 系统配置
const (
	SysSettingsCreate      Operation = "sys:settings.create"
	SysSettingsList        Operation = "sys:settings.list"
	SysSettingsGet         Operation = "sys:settings.get"
	SysSettingsUpdate      Operation = "sys:settings.update"
	SysSettingsDelete      Operation = "sys:settings.delete"
	SysSettingsBatchUpdate Operation = "sys:settings.batchUpdate"
)

// Sys 域 - 配置分类
const (
	SysSettingCategoriesList   Operation = "sys:settings-categories.list"
	SysSettingCategoriesGet    Operation = "sys:settings-categories.get"
	SysSettingCategoriesCreate Operation = "sys:settings-categories.create"
	SysSettingCategoriesUpdate Operation = "sys:settings-categories.update"
	SysSettingCategoriesDelete Operation = "sys:settings-categories.delete"
)

// Sys 域 - 缓存管理
const (
	SysCacheInfo          Operation = "sys:cache.info"
	SysCacheScanKeys      Operation = "sys:cache.scanKeys"
	SysCacheGetKey        Operation = "sys:cache.getKey"
	SysCacheDeleteKey     Operation = "sys:cache.deleteKey"
	SysCacheDeletePattern Operation = "sys:cache.deletePattern"
)

// User 域 - 个人资料
const (
	UserProfileGet     Operation = "user:profile.get"
	UserProfileUpdate  Operation = "user:profile.update"
	UserPasswordUpdate Operation = "user:password.update"
	UserAccountDelete  Operation = "user:account.delete"
)

// User 域 - 访问令牌
//
//nolint:gosec // G101: 这些是操作标识符，非硬编码凭证
const (
	UserTokensCreate  Operation = "user:tokens.create"
	UserTokensList    Operation = "user:tokens.list"
	UserTokensGet     Operation = "user:tokens.get"
	UserTokensDelete  Operation = "user:tokens.delete"
	UserTokensDisable Operation = "user:tokens.disable"
	UserTokensEnable  Operation = "user:tokens.enable"
)

// User 域 - 用户配置
const (
	UserSettingsCategoriesList Operation = "user:settings-categories.list"
	UserSettingsList           Operation = "user:settings.list"
	UserSettingsGet            Operation = "user:settings.get"
	UserSettingsSet            Operation = "user:settings.set"
	UserSettingsReset          Operation = "user:settings.reset"
	UserSettingsBatchSet       Operation = "user:settings.batchSet"
)
