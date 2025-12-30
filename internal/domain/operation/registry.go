package operation

// ============================================================================
// Swagger 分组常量
// ============================================================================

//nolint:gochecknoglobals // 分组常量是只读配置
var (
	groupAuth         = "认证 (Auth)"
	group2FA          = "认证 - 2FA (Auth - 2FA)"
	groupSysUser      = "系统管理 - 用户管理 (Sys - User)"
	groupSysRole      = "系统管理 - 角色管理 (Sys - Role)"
	groupSysOperation = "系统管理 - 操作管理 (Sys - Operation)"
	groupSysAuditLog  = "系统管理 - 审计日志 (Sys - Audit Log)"
	groupSysMenu      = "系统管理 - 菜单管理 (Sys - Menu)"
	groupSysOverview  = "系统管理 - 系统概览 (Sys - Overview)"
	groupSysSetting   = "系统管理 - 系统配置 (Sys - Setting)"
	groupSysCache     = "系统管理 - 缓存管理 (Sys - Cache)"
	groupUserProfile  = "用户中心 - 个人资料 (User - Profile)"
	groupUserToken    = "用户中心 - 访问令牌 (User - Token)" //nolint:gosec // 非凭据
	groupUserSetting  = "用户中心 - 用户配置 (User - Setting)"
)

// ============================================================================
// 预定义基类（路径前缀 + 分组 + 审计分类）
// ============================================================================

//nolint:gochecknoglobals // 基类配置是只读配置
var (
	baseAuth         = B().Path("/api/auth").Group(groupAuth).Cat(AuditCatAuth)
	base2FA          = B().Path("/api/auth/2fa").Group(group2FA).Cat(AuditCatAuth)
	baseSysUser      = B().Path("/api/system/users").Group(groupSysUser).Cat(AuditCatUser)
	baseSysRole      = B().Path("/api/system/roles").Group(groupSysRole).Cat(AuditCatRole)
	baseSysOperation = B().Path("/api/system/operations").Group(groupSysOperation)
	baseSysAuditLog  = B().Path("/api/system/auditlogs").Group(groupSysAuditLog)
	baseSysMenu      = B().Path("/api/system/menus").Group(groupSysMenu).Cat(AuditCatMenu)
	baseSysOverview  = B().Path("/api/system/overview").Group(groupSysOverview)
	baseSysSetting   = B().Path("/api/system/settings").Group(groupSysSetting).Cat(AuditCatSetting)
	baseSysCache     = B().Path("/api/system/cache").Group(groupSysCache).Cat(AuditCatCache)
	baseSelfProfile  = B().Path("/api/user").Group(groupUserProfile).Cat(AuditCatProfile)
	baseSelfToken    = B().Path("/api/user/tokens").Group(groupUserToken).Cat(AuditCatToken)
	baseSelfSetting  = B().Path("/api/user/settings").Group(groupUserSetting).Cat(AuditCatUserSetting)
)

// ============================================================================
// 操作注册表
// ============================================================================

// operationRegistry 操作注册表（单一数据源）
//
// URN-Centric RBAC:
//   - Operation code 本身即权限标识符（URN 格式）
//   - scope 为 "public" 的操作无需权限检查
//   - 非公开操作默认需要权限检查
//
//nolint:gochecknoglobals // 注册表是只读全局配置
var operationRegistry = Build(
	// ==================== Public 域（公开） ====================
	D(PublicAuthRegister).Use(baseAuth).
		Post("/register").Audit("auth.register", AuditOpCreate).
		I18n("注册", "User registration"),

	D(PublicAuthLogin).Use(baseAuth).
		Post("/login").Audit("auth.login", AuditOpAuthenticate).
		I18n("登录", "User login"),

	D(PublicAuthLogin2FA).Use(baseAuth).
		Post("/login/2fa").Audit("auth.login_2fa", AuditOpAuthenticate).
		I18n("2FA 登录", "Two-factor authentication login"),

	D(PublicAuthRefresh).Use(baseAuth).
		Post("/refresh").Audit("auth.refresh", AuditOpAuthenticate).
		I18n("刷新令牌", "Refresh access token"),

	D(PublicAuthCaptcha).Use(baseAuth).
		Get("/captcha").
		I18n("获取验证码", "Get captcha image"),

	// ==================== Self 域 - 2FA（需认证） ====================
	D(Self2FASetup).Use(base2FA).
		Post("/setup").Audit("auth.2fa_setup", AuditOpUpdate).
		I18n("设置 2FA", "Setup two-factor authentication"),

	D(Self2FAVerify).Use(base2FA).
		Post("/verify").Audit("auth.2fa_enable", AuditOpUpdate).
		I18n("启用 2FA", "Verify and enable 2FA"),

	D(Self2FADisable).Use(base2FA).
		Post("/disable").Audit("auth.2fa_disable", AuditOpUpdate).
		I18n("禁用 2FA", "Disable two-factor authentication"),

	D(Self2FAStatus).Use(base2FA).
		Get("/status").
		I18n("2FA 状态", "Get 2FA status"),

	// ==================== Sys 域 - 用户管理 ====================
	D(SysUsersCreate).Use(baseSysUser).
		Post("").AuditCreate().
		I18n("创建用户", "Create new user"),

	D(SysUsersBatchCreate).Use(baseSysUser).
		Post("/batch").Audit("user.batch_create", AuditOpCreate).
		I18n("批量创建用户", "Batch create users"),

	D(SysUsersList).
		Use(baseSysUser).
		Get("").
		I18n("用户列表", "List users"),

	D(SysUsersGet).Use(baseSysUser).
		Get("/:id").
		I18n("用户详情", "Get user by ID"),

	D(SysUsersUpdate).Use(baseSysUser).
		Put("/:id").AuditUpdate().
		I18n("更新用户", "Update user"),

	D(SysUsersDelete).Use(baseSysUser).
		Delete("/:id").AuditDelete().
		I18n("删除用户", "Delete user"),

	D(SysUsersAssignRoles).Use(baseSysUser).
		Put("/:id/roles").Audit("user.assign_roles", AuditOpUpdate).
		I18n("分配角色", "Assign roles to user"),

	// ==================== Sys 域 - 角色管理 ====================
	D(SysRolesCreate).Use(baseSysRole).
		Post("").AuditCreate().
		I18n("创建角色", "Create role"),

	D(SysRolesList).Use(baseSysRole).
		Get("").
		I18n("角色列表", "List roles"),

	D(SysRolesGet).Use(baseSysRole).
		Get("/:id").
		I18n("角色详情", "Get role by ID"),

	D(SysRolesUpdate).Use(baseSysRole).
		Put("/:id").AuditUpdate().
		I18n("更新角色", "Update role"),

	D(SysRolesDelete).Use(baseSysRole).
		Delete("/:id").AuditDelete().
		I18n("删除角色", "Delete role"),

	D(SysRolesSetPermissions).Use(baseSysRole).
		Put("/:id/permissions").Audit("role.set_permissions", AuditOpUpdate).
		I18n("设置权限", "Set role permissions"),

	// ==================== Sys 域 - 操作列表 ====================
	D(SysOperationsList).Use(baseSysOperation).
		Get("").
		I18n("操作列表", "List available operations"),

	// ==================== Sys 域 - 审计日志 ====================
	D(SysAuditLogsList).Use(baseSysAuditLog).
		Get("").
		I18n("审计日志列表", "List audit logs"),

	D(SysAuditLogsGet).Use(baseSysAuditLog).
		Get("/:id").
		I18n("审计日志详情", "Get audit log by ID"),

	D(SysAuditLogsActions).Use(baseSysAuditLog).
		Get("/actions").
		I18n("审计操作定义", "Get audit action definitions"),

	// ==================== Sys 域 - 菜单管理 ====================
	D(SysMenusCreate).Use(baseSysMenu).
		Post("").AuditCreate().
		I18n("创建菜单", "Create menu"),

	D(SysMenusList).Use(baseSysMenu).
		Get("").
		I18n("菜单列表", "List menus"),

	D(SysMenusGet).Use(baseSysMenu).
		Get("/:id").
		I18n("菜单详情", "Get menu by ID"),

	D(SysMenusUpdate).Use(baseSysMenu).
		Put("/:id").AuditUpdate().
		I18n("更新菜单", "Update menu"),

	D(SysMenusDelete).Use(baseSysMenu).
		Delete("/:id").AuditDelete().
		I18n("删除菜单", "Delete menu"),

	D(SysMenusReorder).Use(baseSysMenu).
		Post("/reorder").Audit("menu.reorder", AuditOpUpdate).
		I18n("重排菜单", "Reorder menus"),

	// ==================== Sys 域 - 系统概览 ====================
	D(SysOverviewStats).Use(baseSysOverview).
		Get("/stats").
		I18n("系统概览", "Get system overview stats"),

	// ==================== Sys 域 - 系统配置 ====================
	D(SysSettingsCreate).Use(baseSysSetting).
		Post("").AuditCreate().
		I18n("创建配置", "Create setting"),

	D(SysSettingsList).Use(baseSysSetting).
		Get("").
		I18n("配置列表", "List settings"),

	D(SysSettingsGet).Use(baseSysSetting).
		Get("/:key").
		I18n("配置详情", "Get setting by key"),

	D(SysSettingsUpdate).Use(baseSysSetting).
		Put("/:key").AuditUpdate().
		I18n("更新配置", "Update setting"),

	D(SysSettingsDelete).Use(baseSysSetting).
		Delete("/:key").AuditDelete().
		I18n("删除配置", "Delete setting"),

	D(SysSettingsBatchUpdate).Use(baseSysSetting).
		Post("/batch").Audit("setting.batch_update", AuditOpUpdate).
		I18n("批量更新配置", "Batch update settings"),

	// ==================== Sys 域 - 配置分类 ====================
	D(SysSettingCategoriesList).Use(baseSysSetting).
		Get("/categories").
		I18n("配置分类列表", "List setting categories"),

	D(SysSettingCategoriesGet).Use(baseSysSetting).
		Get("/categories/:id").
		I18n("配置分类详情", "Get setting category by ID"),

	D(SysSettingCategoriesCreate).Use(baseSysSetting).
		Post("/categories").Audit("setting_category.create", AuditOpCreate).
		I18n("创建配置分类", "Create setting category"),

	D(SysSettingCategoriesUpdate).Use(baseSysSetting).
		Put("/categories/:id").Audit("setting_category.update", AuditOpUpdate).
		I18n("更新配置分类", "Update setting category"),

	D(SysSettingCategoriesDelete).Use(baseSysSetting).
		Delete("/categories/:id").Audit("setting_category.delete", AuditOpDelete).
		I18n("删除配置分类", "Delete setting category"),

	// ==================== Sys 域 - 缓存管理 ====================
	D(SysCacheInfo).Use(baseSysCache).
		Get("/info").
		I18n("缓存信息", "Get cache info"),

	D(SysCacheScanKeys).Use(baseSysCache).
		Get("/keys").
		I18n("扫描缓存键", "Scan cache keys"),

	D(SysCacheGetKey).Use(baseSysCache).
		Get("/key").
		I18n("获取缓存值", "Get cache key value"),

	D(SysCacheDeleteKey).Use(baseSysCache).
		Delete("/key").AuditDelete().
		I18n("删除缓存键", "Delete cache key"),

	D(SysCacheDeletePattern).Use(baseSysCache).
		Delete("/keys").Audit("cache.delete_pattern", AuditOpDelete).
		I18n("批量删除缓存", "Delete cache keys by pattern"),

	// ==================== Self 域 - 个人资料 ====================
	D(SelfProfileGet).Use(baseSelfProfile).
		Get("/profile").
		I18n("获取资料", "Get current user profile"),

	D(SelfProfileUpdate).Use(baseSelfProfile).
		Put("/profile").AuditUpdate().
		I18n("更新资料", "Update current user profile"),

	D(SelfPasswordUpdate).Use(baseSelfProfile).
		Put("/password").Audit("password.update", AuditOpUpdate).
		I18n("修改密码", "Change password"),

	D(SelfAccountDelete).Use(baseSelfProfile).
		Delete("/account").Audit("account.delete", AuditOpDelete).
		I18n("注销账户", "Delete own account"),

	// ==================== Self 域 - 访问令牌 ====================
	D(SelfTokensCreate).Use(baseSelfToken).
		Post("").AuditCreate().
		I18n("创建令牌", "Create personal access token"),

	D(SelfTokensList).Use(baseSelfToken).
		Get("").
		I18n("令牌列表", "List personal access tokens"),

	D(SelfTokensGet).Use(baseSelfToken).
		Get("/:id").
		I18n("令牌详情", "Get token by ID"),

	D(SelfTokensDelete).Use(baseSelfToken).
		Delete("/:id").AuditDelete().
		I18n("删除令牌", "Delete token"),

	D(SelfTokensDisable).Use(baseSelfToken).
		Patch("/:id/disable").Audit("token.disable", AuditOpUpdate).
		I18n("禁用令牌", "Disable token"),

	D(SelfTokensEnable).Use(baseSelfToken).
		Patch("/:id/enable").Audit("token.enable", AuditOpUpdate).
		I18n("启用令牌", "Enable token"),

	// ==================== Self 域 - 用户配置 ====================
	D(SelfSettingsCategoriesList).Use(baseSelfSetting).
		Get("/categories").
		I18n("配置分类列表", "List user setting categories"),

	D(SelfSettingsList).Use(baseSelfSetting).
		Get("").
		I18n("用户配置列表", "Get all user settings"),

	D(SelfSettingsGet).Use(baseSelfSetting).
		Get("/:key").
		I18n("用户配置详情", "Get user setting by key"),

	D(SelfSettingsSet).Use(baseSelfSetting).
		Put("/:key").Audit("user_setting.set", AuditOpUpdate).
		I18n("设置用户配置", "Set user setting"),

	D(SelfSettingsReset).Use(baseSelfSetting).
		Delete("/:key").Audit("user_setting.reset", AuditOpUpdate).
		I18n("重置用户配置", "Reset user setting to default"),

	D(SelfSettingsBatchSet).Use(baseSelfSetting).
		Post("/batch").Audit("user_setting.batch_set", AuditOpUpdate).
		I18n("批量设置配置", "Batch set user settings"),
)
