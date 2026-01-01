package permission

// ============================================================================
// 操作常量
// ============================================================================
//
// 格式：{scope}:{type}:{action}
//
// Scope 划分：
//   - public: 公开域（无需认证）
//   - admin:  系统管理域（需管理员权限）
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

// Admin 域 - 用户管理
const (
	AdminUsersCreate      Operation = "admin:users:create"
	AdminUsersBatchCreate Operation = "admin:users:batch-create"
	AdminUsersList        Operation = "admin:users:list"
	AdminUsersGet         Operation = "admin:users:get"
	AdminUsersUpdate      Operation = "admin:users:update"
	AdminUsersDelete      Operation = "admin:users:delete"
	AdminUsersAssignRoles Operation = "admin:users:assign-roles"
)

// Admin 域 - 角色管理
const (
	AdminRolesCreate         Operation = "admin:roles:create"
	AdminRolesList           Operation = "admin:roles:list"
	AdminRolesGet            Operation = "admin:roles:get"
	AdminRolesUpdate         Operation = "admin:roles:update"
	AdminRolesDelete         Operation = "admin:roles:delete"
	AdminRolesSetPermissions Operation = "admin:roles:set-permissions"
)

// Admin 域 - 操作列表（供前端权限配置使用）
const (
	AdminOperationsList Operation = "admin:operations:list"
)

// Admin 域 - 审计日志
const (
	AdminAuditList    Operation = "admin:audit:list"
	AdminAuditGet     Operation = "admin:audit:get"
	AdminAuditActions Operation = "admin:audit:actions"
)

// Admin 域 - 系统概览
const (
	AdminOverviewStats Operation = "admin:overview:stats"
)

// Admin 域 - 系统配置
const (
	AdminSettingsCreate      Operation = "admin:settings:create"
	AdminSettingsList        Operation = "admin:settings:list"
	AdminSettingsGet         Operation = "admin:settings:get"
	AdminSettingsUpdate      Operation = "admin:settings:update"
	AdminSettingsDelete      Operation = "admin:settings:delete"
	AdminSettingsBatchUpdate Operation = "admin:settings:batch-update"
)

// Admin 域 - 配置分类
const (
	AdminSettingCategoriesList   Operation = "admin:settings-categories:list"
	AdminSettingCategoriesGet    Operation = "admin:settings-categories:get"
	AdminSettingCategoriesCreate Operation = "admin:settings-categories:create"
	AdminSettingCategoriesUpdate Operation = "admin:settings-categories:update"
	AdminSettingCategoriesDelete Operation = "admin:settings-categories:delete"
)

// Admin 域 - 缓存管理
const (
	AdminCacheInfo          Operation = "admin:cache:info"
	AdminCacheScanKeys      Operation = "admin:cache:scan-keys"
	AdminCacheGetKey        Operation = "admin:cache:get-key"
	AdminCacheDeleteKey     Operation = "admin:cache:delete-key"
	AdminCacheDeletePattern Operation = "admin:cache:delete-pattern"
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

// Admin 域 - 组织管理
const (
	AdminOrgsCreate Operation = "admin:orgs:create"
	AdminOrgsList   Operation = "admin:orgs:list"
	AdminOrgsGet    Operation = "admin:orgs:get"
	AdminOrgsUpdate Operation = "admin:orgs:update"
	AdminOrgsDelete Operation = "admin:orgs:delete"
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

// Admin 域 - 产品管理
const (
	AdminProductsCreate Operation = "admin:products:create"
	AdminProductsList   Operation = "admin:products:list"
	AdminProductsGet    Operation = "admin:products:get"
	AdminProductsUpdate Operation = "admin:products:update"
	AdminProductsDelete Operation = "admin:products:delete"
)
