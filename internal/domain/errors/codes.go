package errors

// ============================================================================
// 用户模块错误码
// ============================================================================

// 用户相关错误
var (
	// ErrUserNotFound 用户不存在
	ErrUserNotFound = NewNotFoundError("user.not_found", "user not found")

	// ErrUserAlreadyExists 用户已存在
	ErrUserAlreadyExists = NewConflictError("user.already_exists", "user already exists")

	// ErrUsernameAlreadyExists 用户名已存在
	ErrUsernameAlreadyExists = NewConflictError("user.username_exists", "username already exists")

	// ErrEmailAlreadyExists 邮箱已存在
	ErrEmailAlreadyExists = NewConflictError("user.email_exists", "email already exists")

	// ErrRoleAlreadyAssigned 角色已分配
	ErrRoleAlreadyAssigned = NewConflictError("user.role_already_assigned", "role already assigned to user")

	// ErrInvalidUserStatus 无效的用户状态
	ErrInvalidUserStatus = NewValidationError("user.invalid_status", "invalid user status")

	// ErrCannotDeleteSelf 不能删除自己
	ErrCannotDeleteSelf = NewForbiddenError("user.cannot_delete_self", "cannot delete yourself")

	// ErrCannotModifyAdmin 不能修改管理员
	ErrCannotModifyAdmin = NewForbiddenError("user.cannot_modify_admin", "cannot modify admin user")
)

// ============================================================================
// 认证模块错误码
// ============================================================================

// 认证相关错误
var (
	// ErrInvalidCredentials 凭据无效
	ErrInvalidCredentials = NewUnauthorizedError("auth.invalid_credentials", "invalid username or password")

	// ErrInvalidPassword 密码错误
	ErrInvalidPassword = NewUnauthorizedError("auth.invalid_password", "invalid password")

	// ErrInvalidToken 无效的令牌
	ErrInvalidToken = NewUnauthorizedError("auth.invalid_token", "invalid token")

	// ErrTokenExpired 令牌已过期
	ErrTokenExpired = NewUnauthorizedError("auth.token_expired", "token has expired")

	// ErrUserInactive 用户未激活
	ErrUserInactive = NewForbiddenError("auth.user_inactive", "user account is inactive")

	// ErrUserBanned 用户已禁用
	ErrUserBanned = NewForbiddenError("auth.user_banned", "user account is banned")

	// ErrPasswordPolicyViolation 密码策略违规
	ErrPasswordPolicyViolation = NewValidationError("auth.password_policy", "password does not meet policy requirements")
)

// ============================================================================
// 角色模块错误码
// ============================================================================

// 角色相关错误
var (
	// ErrRoleNotFound 角色不存在
	ErrRoleNotFound = NewNotFoundError("role.not_found", "role not found")

	// ErrRoleAlreadyExists 角色已存在
	ErrRoleAlreadyExists = NewConflictError("role.already_exists", "role already exists")

	// ErrSystemRoleProtected 系统角色受保护
	ErrSystemRoleProtected = NewForbiddenError("role.system_protected", "system role cannot be modified or deleted")

	// ErrPermissionNotFound 权限不存在
	ErrPermissionNotFound = NewNotFoundError("permission.not_found", "permission not found")
)

// ============================================================================
// 菜单模块错误码
// ============================================================================

// 菜单相关错误
var (
	// ErrMenuNotFound 菜单不存在
	ErrMenuNotFound = NewNotFoundError("menu.not_found", "menu not found")

	// ErrMenuAlreadyExists 菜单已存在
	ErrMenuAlreadyExists = NewConflictError("menu.already_exists", "menu already exists")

	// ErrMenuHasChildren 菜单有子菜单
	ErrMenuHasChildren = NewConflictError("menu.has_children", "menu has children and cannot be deleted")

	// ErrInvalidMenuParent 无效的父菜单
	ErrInvalidMenuParent = NewValidationError("menu.invalid_parent", "invalid parent menu")
)

// ============================================================================
// PAT 模块错误码
// ============================================================================

// PAT 相关错误
var (
	// ErrTokenNotFound 令牌不存在
	ErrTokenNotFound = NewNotFoundError("pat.not_found", "personal access token not found")

	// ErrTokenDisabled 令牌已禁用
	ErrTokenDisabled = NewForbiddenError("pat.disabled", "personal access token is disabled")

	// ErrIPNotAllowed IP 不在白名单
	ErrIPNotAllowed = NewForbiddenError("pat.ip_not_allowed", "IP address not in whitelist")
)

// ============================================================================
// 设置模块错误码
// ============================================================================

// 设置相关错误
var (
	// ErrSettingNotFound 设置不存在
	ErrSettingNotFound = NewNotFoundError("setting.not_found", "setting not found")

	// ErrSettingKeyExists 设置键已存在
	ErrSettingKeyExists = NewConflictError("setting.key_exists", "setting key already exists")

	// ErrInvalidSettingValue 无效的设置值
	ErrInvalidSettingValue = NewValidationError("setting.invalid_value", "invalid setting value")
)

// ============================================================================
// 验证码模块错误码
// ============================================================================

// 验证码相关错误
var (
	// ErrCaptchaNotFound 验证码不存在
	ErrCaptchaNotFound = NewNotFoundError("captcha.not_found", "captcha not found")

	// ErrCaptchaExpired 验证码已过期
	ErrCaptchaExpired = NewValidationError("captcha.expired", "captcha has expired")

	// ErrCaptchaInvalid 验证码无效
	ErrCaptchaInvalid = NewValidationError("captcha.invalid", "invalid captcha code")
)

// ============================================================================
// 2FA 模块错误码
// ============================================================================

// 2FA 相关错误
var (
	// ErrTwoFANotEnabled 2FA 未启用
	ErrTwoFANotEnabled = NewValidationError("twofa.not_enabled", "two-factor authentication is not enabled")

	// ErrTwoFAAlreadyEnabled 2FA 已启用
	ErrTwoFAAlreadyEnabled = NewConflictError("twofa.already_enabled", "two-factor authentication is already enabled")

	// ErrInvalidTwoFACode 无效的 2FA 代码
	ErrInvalidTwoFACode = NewValidationError("twofa.invalid_code", "invalid two-factor authentication code")

	// ErrInvalidRecoveryCode 无效的恢复码
	ErrInvalidRecoveryCode = NewValidationError("twofa.invalid_recovery", "invalid recovery code")
)

// ============================================================================
// 通用错误码
// ============================================================================

// 通用错误
var (
	// ErrInternalError 内部错误
	ErrInternalError = NewInternalError("internal.error", "internal server error")

	// ErrDatabaseError 数据库错误
	ErrDatabaseError = NewInternalError("internal.database", "database operation failed")

	// ErrCacheError 缓存错误
	ErrCacheError = NewInternalError("internal.cache", "cache operation failed")
)
