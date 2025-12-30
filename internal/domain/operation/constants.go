package operation

// ============================================================================
// 操作常量
// ============================================================================
//
// 格式：{scope}:{type}:{action}
//
// Scope 划分：
//   - public: 公开域（无需认证）
//   - sys:    系统管理域（需管理员权限）
//   - self:   用户自服务域（当前用户权限）

// Public 域（公开接口）
const (
	PublicAuthRegister Operation = "public:auth:register"
	PublicAuthLogin    Operation = "public:auth:login"
	PublicAuthLogin2FA Operation = "public:auth:login2fa"
	PublicAuthRefresh  Operation = "public:auth:refresh"
	PublicAuthCaptcha  Operation = "public:auth:captcha"
)

// Self 域 - 2FA（需认证）
const (
	Self2FASetup   Operation = "self:2fa:setup"
	Self2FAVerify  Operation = "self:2fa:verify"
	Self2FADisable Operation = "self:2fa:disable"
	Self2FAStatus  Operation = "self:2fa:status"
)

// Sys 域 - 用户管理
const (
	SysUsersCreate      Operation = "sys:users:create"
	SysUsersBatchCreate Operation = "sys:users:batch-create"
	SysUsersList        Operation = "sys:users:list"
	SysUsersGet         Operation = "sys:users:get"
	SysUsersUpdate      Operation = "sys:users:update"
	SysUsersDelete      Operation = "sys:users:delete"
	SysUsersAssignRoles Operation = "sys:users:assign-roles"
)

// Sys 域 - 角色管理
const (
	SysRolesCreate         Operation = "sys:roles:create"
	SysRolesList           Operation = "sys:roles:list"
	SysRolesGet            Operation = "sys:roles:get"
	SysRolesUpdate         Operation = "sys:roles:update"
	SysRolesDelete         Operation = "sys:roles:delete"
	SysRolesSetPermissions Operation = "sys:roles:set-permissions"
)

// Sys 域 - 操作列表（供前端权限配置使用）
const (
	SysOperationsList Operation = "sys:operations:list"
)

// Sys 域 - 审计日志
const (
	SysAuditLogsList    Operation = "sys:auditlogs:list"
	SysAuditLogsGet     Operation = "sys:auditlogs:get"
	SysAuditLogsActions Operation = "sys:auditlogs:actions"
)

// Sys 域 - 菜单管理
const (
	SysMenusCreate  Operation = "sys:menus:create"
	SysMenusList    Operation = "sys:menus:list"
	SysMenusGet     Operation = "sys:menus:get"
	SysMenusUpdate  Operation = "sys:menus:update"
	SysMenusDelete  Operation = "sys:menus:delete"
	SysMenusReorder Operation = "sys:menus:reorder"
)

// Sys 域 - 系统概览
const (
	SysOverviewStats Operation = "sys:overview:stats"
)

// Sys 域 - 系统配置
const (
	SysSettingsCreate      Operation = "sys:settings:create"
	SysSettingsList        Operation = "sys:settings:list"
	SysSettingsGet         Operation = "sys:settings:get"
	SysSettingsUpdate      Operation = "sys:settings:update"
	SysSettingsDelete      Operation = "sys:settings:delete"
	SysSettingsBatchUpdate Operation = "sys:settings:batch-update"
)

// Sys 域 - 配置分类
const (
	SysSettingCategoriesList   Operation = "sys:settings-categories:list"
	SysSettingCategoriesGet    Operation = "sys:settings-categories:get"
	SysSettingCategoriesCreate Operation = "sys:settings-categories:create"
	SysSettingCategoriesUpdate Operation = "sys:settings-categories:update"
	SysSettingCategoriesDelete Operation = "sys:settings-categories:delete"
)

// Sys 域 - 缓存管理
const (
	SysCacheInfo          Operation = "sys:cache:info"
	SysCacheScanKeys      Operation = "sys:cache:scan-keys"
	SysCacheGetKey        Operation = "sys:cache:get-key"
	SysCacheDeleteKey     Operation = "sys:cache:delete-key"
	SysCacheDeletePattern Operation = "sys:cache:delete-pattern"
)

// Self 域 - 个人资料
const (
	SelfProfileGet     Operation = "self:profile:get"
	SelfProfileUpdate  Operation = "self:profile:update"
	SelfPasswordUpdate Operation = "self:password:update"
	SelfAccountDelete  Operation = "self:account:delete"
)

// Self 域 - 访问令牌
//
//nolint:gosec // G101: 这些是操作标识符，非硬编码凭证
const (
	SelfTokensCreate  Operation = "self:tokens:create"
	SelfTokensList    Operation = "self:tokens:list"
	SelfTokensGet     Operation = "self:tokens:get"
	SelfTokensDelete  Operation = "self:tokens:delete"
	SelfTokensDisable Operation = "self:tokens:disable"
	SelfTokensEnable  Operation = "self:tokens:enable"
)

// Self 域 - 用户配置
const (
	SelfSettingsCategoriesList Operation = "self:settings-categories:list"
	SelfSettingsList           Operation = "self:settings:list"
	SelfSettingsGet            Operation = "self:settings:get"
	SelfSettingsSet            Operation = "self:settings:set"
	SelfSettingsReset          Operation = "self:settings:reset"
	SelfSettingsBatchSet       Operation = "self:settings:batch-set"
)

// ============================================================================
// 旧常量别名（废弃，保持向后兼容）
// ============================================================================
// Deprecated: 使用新的 URN 风格常量。以下别名将在下一版本移除。

const (
	// AuthRegister 是 PublicAuthRegister 的别名。
	AuthRegister = PublicAuthRegister
	// AuthLogin 是 PublicAuthLogin 的别名。
	AuthLogin = PublicAuthLogin
	// AuthLogin2FA 是 PublicAuthLogin2FA 的别名。
	AuthLogin2FA = PublicAuthLogin2FA
	// AuthRefresh 是 PublicAuthRefresh 的别名。
	AuthRefresh = PublicAuthRefresh
	// AuthCaptcha 是 PublicAuthCaptcha 的别名。
	AuthCaptcha = PublicAuthCaptcha

	// Auth2FASetup 是 Self2FASetup 的别名。
	Auth2FASetup = Self2FASetup
	// Auth2FAVerify 是 Self2FAVerify 的别名。
	Auth2FAVerify = Self2FAVerify
	// Auth2FADisable 是 Self2FADisable 的别名。
	Auth2FADisable = Self2FADisable
	// Auth2FAStatus 是 Self2FAStatus 的别名。
	Auth2FAStatus = Self2FAStatus

	// UserProfileGet 是 SelfProfileGet 的别名。
	UserProfileGet = SelfProfileGet
	// UserProfileUpdate 是 SelfProfileUpdate 的别名。
	UserProfileUpdate = SelfProfileUpdate
	// UserPasswordUpdate 是 SelfPasswordUpdate 的别名。
	UserPasswordUpdate = SelfPasswordUpdate
	// UserAccountDelete 是 SelfAccountDelete 的别名。
	UserAccountDelete = SelfAccountDelete

	// UserTokensCreate 是 SelfTokensCreate 的别名。
	UserTokensCreate = SelfTokensCreate
	// UserTokensList 是 SelfTokensList 的别名。
	UserTokensList = SelfTokensList
	// UserTokensGet 是 SelfTokensGet 的别名。
	UserTokensGet = SelfTokensGet
	// UserTokensDelete 是 SelfTokensDelete 的别名。
	UserTokensDelete = SelfTokensDelete
	// UserTokensDisable 是 SelfTokensDisable 的别名。
	UserTokensDisable = SelfTokensDisable
	// UserTokensEnable 是 SelfTokensEnable 的别名。
	UserTokensEnable = SelfTokensEnable

	// UserSettingsCategoriesList 是 SelfSettingsCategoriesList 的别名。
	UserSettingsCategoriesList = SelfSettingsCategoriesList
	// UserSettingsList 是 SelfSettingsList 的别名。
	UserSettingsList = SelfSettingsList
	// UserSettingsGet 是 SelfSettingsGet 的别名。
	UserSettingsGet = SelfSettingsGet
	// UserSettingsSet 是 SelfSettingsSet 的别名。
	UserSettingsSet = SelfSettingsSet
	// UserSettingsReset 是 SelfSettingsReset 的别名。
	UserSettingsReset = SelfSettingsReset
	// UserSettingsBatchSet 是 SelfSettingsBatchSet 的别名。
	UserSettingsBatchSet = SelfSettingsBatchSet
)
