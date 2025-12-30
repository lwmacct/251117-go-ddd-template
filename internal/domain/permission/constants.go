package permission

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
// 常量名查找（供代码生成工具使用）
// ============================================================================

// constNameRegistry 常量名到 Operation 的映射。
// 用于从源代码解析的常量名反查 Operation。
//
//nolint:gochecknoglobals // 只读注册表
var constNameRegistry = map[string]Operation{
	// Public 域
	"PublicAuthRegister": PublicAuthRegister,
	"PublicAuthLogin":    PublicAuthLogin,
	"PublicAuthLogin2FA": PublicAuthLogin2FA,
	"PublicAuthRefresh":  PublicAuthRefresh,
	"PublicAuthCaptcha":  PublicAuthCaptcha,

	// Self 域 - 2FA
	"Self2FASetup":   Self2FASetup,
	"Self2FAVerify":  Self2FAVerify,
	"Self2FADisable": Self2FADisable,
	"Self2FAStatus":  Self2FAStatus,

	// Sys 域 - 用户管理
	"SysUsersCreate":      SysUsersCreate,
	"SysUsersBatchCreate": SysUsersBatchCreate,
	"SysUsersList":        SysUsersList,
	"SysUsersGet":         SysUsersGet,
	"SysUsersUpdate":      SysUsersUpdate,
	"SysUsersDelete":      SysUsersDelete,
	"SysUsersAssignRoles": SysUsersAssignRoles,

	// Sys 域 - 角色管理
	"SysRolesCreate":         SysRolesCreate,
	"SysRolesList":           SysRolesList,
	"SysRolesGet":            SysRolesGet,
	"SysRolesUpdate":         SysRolesUpdate,
	"SysRolesDelete":         SysRolesDelete,
	"SysRolesSetPermissions": SysRolesSetPermissions,

	// Sys 域 - 操作列表
	"SysOperationsList": SysOperationsList,

	// Sys 域 - 审计日志
	"SysAuditLogsList":    SysAuditLogsList,
	"SysAuditLogsGet":     SysAuditLogsGet,
	"SysAuditLogsActions": SysAuditLogsActions,

	// Sys 域 - 菜单管理
	"SysMenusCreate":  SysMenusCreate,
	"SysMenusList":    SysMenusList,
	"SysMenusGet":     SysMenusGet,
	"SysMenusUpdate":  SysMenusUpdate,
	"SysMenusDelete":  SysMenusDelete,
	"SysMenusReorder": SysMenusReorder,

	// Sys 域 - 系统概览
	"SysOverviewStats": SysOverviewStats,

	// Sys 域 - 系统配置
	"SysSettingsCreate":      SysSettingsCreate,
	"SysSettingsList":        SysSettingsList,
	"SysSettingsGet":         SysSettingsGet,
	"SysSettingsUpdate":      SysSettingsUpdate,
	"SysSettingsDelete":      SysSettingsDelete,
	"SysSettingsBatchUpdate": SysSettingsBatchUpdate,

	// Sys 域 - 配置分类
	"SysSettingCategoriesList":   SysSettingCategoriesList,
	"SysSettingCategoriesGet":    SysSettingCategoriesGet,
	"SysSettingCategoriesCreate": SysSettingCategoriesCreate,
	"SysSettingCategoriesUpdate": SysSettingCategoriesUpdate,
	"SysSettingCategoriesDelete": SysSettingCategoriesDelete,

	// Sys 域 - 缓存管理
	"SysCacheInfo":          SysCacheInfo,
	"SysCacheScanKeys":      SysCacheScanKeys,
	"SysCacheGetKey":        SysCacheGetKey,
	"SysCacheDeleteKey":     SysCacheDeleteKey,
	"SysCacheDeletePattern": SysCacheDeletePattern,

	// Self 域 - 个人资料
	"SelfProfileGet":     SelfProfileGet,
	"SelfProfileUpdate":  SelfProfileUpdate,
	"SelfPasswordUpdate": SelfPasswordUpdate,
	"SelfAccountDelete":  SelfAccountDelete,

	// Self 域 - 访问令牌
	"SelfTokensCreate":  SelfTokensCreate,
	"SelfTokensList":    SelfTokensList,
	"SelfTokensGet":     SelfTokensGet,
	"SelfTokensDelete":  SelfTokensDelete,
	"SelfTokensDisable": SelfTokensDisable,
	"SelfTokensEnable":  SelfTokensEnable,

	// Self 域 - 用户配置
	"SelfSettingsCategoriesList": SelfSettingsCategoriesList,
	"SelfSettingsList":           SelfSettingsList,
	"SelfSettingsGet":            SelfSettingsGet,
	"SelfSettingsSet":            SelfSettingsSet,
	"SelfSettingsReset":          SelfSettingsReset,
	"SelfSettingsBatchSet":       SelfSettingsBatchSet,
}

// ByConstName 通过常量名查找操作。
// 如果未找到返回空 Operation。
func ByConstName(name string) Operation {
	return constNameRegistry[name]
}
