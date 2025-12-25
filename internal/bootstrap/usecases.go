package bootstrap

import (
	"github.com/lwmacct/251117-go-ddd-template/internal/application/auditlog"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/auth"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/cache"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/captcha"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/menu"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/pat"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/role"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/setting"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/stats"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/twofa"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/user"
)

// UseCasesModule 用例模块
// 聚合所有 Application 层的 Use Case Handlers
type UseCasesModule struct {
	Auth     *AuthUseCases
	User     *UserUseCases
	Role     *RoleUseCases
	Menu     *MenuUseCases
	Setting  *SettingUseCases
	PAT      *PATUseCases
	AuditLog *AuditLogUseCases
	Stats    *StatsUseCases
	Captcha  *CaptchaUseCases
	TwoFA    *TwoFAUseCases
	Cache    *CacheUseCases
}

// AuthUseCases 认证相关用例
type AuthUseCases struct {
	Login        *auth.LoginHandler
	Login2FA     *auth.Login2FAHandler
	Register     *auth.RegisterHandler
	RefreshToken *auth.RefreshTokenHandler
}

// UserUseCases 用户管理用例
type UserUseCases struct {
	// Commands
	Create         *user.CreateUserHandler
	Update         *user.UpdateUserHandler
	Delete         *user.DeleteUserHandler
	AssignRoles    *user.AssignRolesHandler
	ChangePassword *user.ChangePasswordHandler
	BatchCreate    *user.BatchCreateUsersHandler

	// Queries
	Get  *user.GetUserHandler
	List *user.ListUsersHandler
}

// RoleUseCases 角色管理用例
type RoleUseCases struct {
	// Commands
	Create         *role.CreateRoleHandler
	Update         *role.UpdateRoleHandler
	Delete         *role.DeleteRoleHandler
	SetPermissions *role.SetPermissionsHandler

	// Queries
	Get             *role.GetRoleHandler
	List            *role.ListRolesHandler
	ListPermissions *role.ListPermissionsHandler
}

// MenuUseCases 菜单管理用例
type MenuUseCases struct {
	// Commands
	Create  *menu.CreateMenuHandler
	Update  *menu.UpdateMenuHandler
	Delete  *menu.DeleteMenuHandler
	Reorder *menu.ReorderMenusHandler

	// Queries
	Get  *menu.GetMenuHandler
	List *menu.ListMenusHandler
}

// SettingUseCases 系统配置用例
type SettingUseCases struct {
	// Commands
	Create      *setting.CreateSettingHandler
	Update      *setting.UpdateSettingHandler
	Delete      *setting.DeleteSettingHandler
	BatchUpdate *setting.BatchUpdateSettingsHandler

	// Queries
	Get  *setting.GetSettingHandler
	List *setting.ListSettingsHandler
}

// PATUseCases 个人访问令牌用例
type PATUseCases struct {
	// Commands
	Create  *pat.CreateTokenHandler
	Delete  *pat.DeleteTokenHandler
	Disable *pat.DisableTokenHandler
	Enable  *pat.EnableTokenHandler

	// Queries
	Get  *pat.GetTokenHandler
	List *pat.ListTokensHandler
}

// AuditLogUseCases 审计日志用例
type AuditLogUseCases struct {
	// Commands
	CreateLog *auditlog.CreateLogHandler

	// Queries
	Get  *auditlog.GetLogHandler
	List *auditlog.ListLogsHandler
}

// StatsUseCases 统计用例（只读）
type StatsUseCases struct {
	GetStats *stats.GetStatsHandler
}

// CaptchaUseCases 验证码用例
type CaptchaUseCases struct {
	Generate *captcha.GenerateCaptchaHandler
}

// TwoFAUseCases 双因素认证用例
type TwoFAUseCases struct {
	// Commands
	Setup        *twofa.SetupHandler
	VerifyEnable *twofa.VerifyEnableHandler
	Disable      *twofa.DisableHandler

	// Queries
	GetStatus *twofa.GetStatusHandler
}

// CacheUseCases 缓存用例（演示用）
type CacheUseCases struct {
	// Commands
	Set    *cache.SetCacheHandler
	Delete *cache.DeleteCacheHandler

	// Queries
	Get *cache.GetCacheHandler
}
