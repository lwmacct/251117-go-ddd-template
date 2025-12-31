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
	SelfTokensScopes  Operation = "self:tokens:scopes"
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

// Self 域 - 用户组织/团队
const (
	SelfOrgsList  Operation = "self:orgs:list"
	SelfTeamsList Operation = "self:teams:list"
)

// Sys 域 - 组织管理
const (
	SysOrgsCreate Operation = "sys:orgs:create"
	SysOrgsList   Operation = "sys:orgs:list"
	SysOrgsGet    Operation = "sys:orgs:get"
	SysOrgsUpdate Operation = "sys:orgs:update"
	SysOrgsDelete Operation = "sys:orgs:delete"
)

// Org 域 - 组织成员管理
const (
	OrgMembersList       Operation = "org:members:list"
	OrgMembersAdd        Operation = "org:members:add"
	OrgMembersRemove     Operation = "org:members:remove"
	OrgMembersUpdateRole Operation = "org:members:update-role"
)

// Org 域 - 团队管理
const (
	OrgTeamsCreate Operation = "org:teams:create"
	OrgTeamsList   Operation = "org:teams:list"
	OrgTeamsGet    Operation = "org:teams:get"
	OrgTeamsUpdate Operation = "org:teams:update"
	OrgTeamsDelete Operation = "org:teams:delete"
)

// Org 域 - 团队成员管理
const (
	OrgTeamMembersList   Operation = "org:team-members:list"
	OrgTeamMembersAdd    Operation = "org:team-members:add"
	OrgTeamMembersRemove Operation = "org:team-members:remove"
)
