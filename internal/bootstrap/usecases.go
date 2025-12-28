package bootstrap

import (
	"github.com/lwmacct/251117-go-ddd-template/internal/application/auditlog"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/auth"
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
	Auth        *AuthUseCases
	User        *UserUseCases
	Role        *RoleUseCases
	Menu        *MenuUseCases
	Setting     *SettingUseCases
	UserSetting *UserSettingUseCases
	PAT         *PATUseCases
	AuditLog    *AuditLogUseCases
	Stats       *StatsUseCases
	Captcha     *CaptchaUseCases
	TwoFA       *TwoFAUseCases
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
	Create         *user.CreateHandler
	Update         *user.UpdateHandler
	Delete         *user.DeleteHandler
	AssignRoles    *user.AssignRolesHandler
	ChangePassword *user.ChangePasswordHandler
	BatchCreate    *user.BatchCreateHandler

	// Queries
	Get  *user.GetHandler
	List *user.ListHandler
}

// RoleUseCases 角色管理用例
type RoleUseCases struct {
	// Commands
	Create         *role.CreateHandler
	Update         *role.UpdateHandler
	Delete         *role.DeleteHandler
	SetPermissions *role.SetPermissionsHandler

	// Queries
	Get             *role.GetHandler
	List            *role.ListHandler
	ListPermissions *role.ListPermissionsHandler
}

// MenuUseCases 菜单管理用例
type MenuUseCases struct {
	// Commands
	Create  *menu.CreateHandler
	Update  *menu.UpdateHandler
	Delete  *menu.DeleteHandler
	Reorder *menu.ReorderHandler

	// Queries
	Get  *menu.GetHandler
	List *menu.ListHandler
}

// SettingUseCases 系统配置用例
type SettingUseCases struct {
	// Setting Commands
	Create      *setting.CreateHandler
	Update      *setting.UpdateHandler
	Delete      *setting.DeleteHandler
	BatchUpdate *setting.BatchUpdateHandler

	// Setting Queries
	Get        *setting.GetHandler
	List       *setting.ListHandler
	ListSchema *setting.ListSchemaHandler

	// Category Commands
	CreateCategory *setting.CreateCategoryHandler
	UpdateCategory *setting.UpdateCategoryHandler
	DeleteCategory *setting.DeleteCategoryHandler

	// Category Queries
	GetCategory    *setting.GetCategoryHandler
	ListCategories *setting.ListCategoriesHandler
}

// UserSettingUseCases 用户配置用例
type UserSettingUseCases struct {
	// Commands
	Set      *setting.UserSetHandler
	BatchSet *setting.UserBatchSetHandler
	Reset    *setting.UserResetHandler
	ResetAll *setting.UserResetAllHandler

	// Queries
	Get            *setting.UserGetHandler
	List           *setting.UserListHandler
	ListSchema     *setting.UserListSchemaHandler
	ListCategories *setting.UserListCategoriesHandler
}

// PATUseCases 个人访问令牌用例
type PATUseCases struct {
	// Commands
	Create  *pat.CreateHandler
	Delete  *pat.DeleteHandler
	Disable *pat.DisableHandler
	Enable  *pat.EnableHandler

	// Queries
	Get  *pat.GetHandler
	List *pat.ListHandler
}

// AuditLogUseCases 审计日志用例
type AuditLogUseCases struct {
	// Commands
	CreateLog *auditlog.CreateHandler

	// Queries
	Get  *auditlog.GetHandler
	List *auditlog.ListHandler
}

// StatsUseCases 统计用例（只读）
type StatsUseCases struct {
	GetStats *stats.GetStatsHandler
}

// CaptchaUseCases 验证码用例
type CaptchaUseCases struct {
	Generate *captcha.GenerateHandler
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
