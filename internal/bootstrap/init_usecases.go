package bootstrap

import (
	"github.com/lwmacct/251117-go-ddd-template/internal/config"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/event"

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

	authInfra "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/auth"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/validation"
)

// newUseCasesModule 初始化用例模块
// 依赖：RepositoriesModule, ServicesModule, InfrastructureModule, EventBus, Config
func newUseCasesModule(_ *config.Config, _ *InfrastructureModule, repos *RepositoriesModule, services *ServicesModule, eventBus event.EventBus) *UseCasesModule {
	// 先创建 AuditLog（Auth 依赖它记录登录日志）
	auditLogUseCases := newAuditLogUseCases(repos)

	return &UseCasesModule{
		Auth:        newAuthUseCases(repos, services, auditLogUseCases.CreateLog),
		User:        newUserUseCases(repos, services, eventBus),
		Role:        newRoleUseCases(repos, eventBus),
		Menu:        newMenuUseCases(repos),
		Setting:     newSettingUseCases(repos),
		UserSetting: newUserSettingUseCases(repos),
		PAT:         newPATUseCases(repos, services),
		AuditLog:    auditLogUseCases,
		Stats:       newStatsUseCases(repos),
		Captcha:     newCaptchaUseCases(repos, services),
		TwoFA:       newTwoFAUseCases(services),
	}
}

// newAuthUseCases 初始化认证用例
func newAuthUseCases(repos *RepositoriesModule, services *ServicesModule, auditLogHandler *auditlog.CreateHandler) *AuthUseCases {
	return &AuthUseCases{
		Login:        auth.NewLoginHandler(repos.User.Query, repos.CaptchaCommand, repos.TwoFA.Query, services.Auth, services.LoginSession, auditLogHandler),
		Login2FA:     auth.NewLogin2FAHandler(repos.User.Query, services.Auth, services.LoginSession, services.TwoFA, auditLogHandler),
		Register:     auth.NewRegisterHandler(repos.User.Command, repos.User.Query, services.Auth),
		RefreshToken: auth.NewRefreshTokenHandler(repos.User.Query, services.Auth),
	}
}

// newUserUseCases 初始化用户管理用例
func newUserUseCases(repos *RepositoriesModule, services *ServicesModule, eventBus event.EventBus) *UserUseCases {
	return &UserUseCases{
		Create:         user.NewCreateHandler(repos.User.Command, repos.User.Query, services.Auth),
		Update:         user.NewUpdateHandler(repos.User.Command, repos.User.Query),
		Delete:         user.NewDeleteHandler(repos.User.Command, repos.User.Query, eventBus),
		AssignRoles:    user.NewAssignRolesHandler(repos.User.Command, repos.User.Query, eventBus),
		ChangePassword: user.NewChangePasswordHandler(repos.User.Command, repos.User.Query, services.Auth),
		BatchCreate:    user.NewBatchCreateHandler(repos.User.Command, repos.User.Query, services.Auth),
		Get:            user.NewGetHandler(repos.User.Query),
		List:           user.NewListHandler(repos.User.Query),
	}
}

// newRoleUseCases 初始化角色管理用例
func newRoleUseCases(repos *RepositoriesModule, eventBus event.EventBus) *RoleUseCases {
	return &RoleUseCases{
		Create:          role.NewCreateHandler(repos.Role.Command, repos.Role.Query),
		Update:          role.NewUpdateHandler(repos.Role.Command, repos.Role.Query),
		Delete:          role.NewDeleteHandler(repos.Role.Command, repos.Role.Query),
		SetPermissions:  role.NewSetPermissionsHandler(repos.Role.Command, repos.Role.Query, repos.Permission.Query, eventBus),
		Get:             role.NewGetHandler(repos.Role.Query),
		List:            role.NewListHandler(repos.Role.Query),
		ListPermissions: role.NewListPermissionsHandler(repos.Permission.Query),
	}
}

// newMenuUseCases 初始化菜单管理用例
func newMenuUseCases(repos *RepositoriesModule) *MenuUseCases {
	return &MenuUseCases{
		Create:  menu.NewCreateHandler(repos.Menu.Command, repos.Menu.Query),
		Update:  menu.NewUpdateHandler(repos.Menu.Command, repos.Menu.Query),
		Delete:  menu.NewDeleteHandler(repos.Menu.Command, repos.Menu.Query),
		Reorder: menu.NewReorderHandler(repos.Menu.Command, repos.Menu.Query),
		Get:     menu.NewGetHandler(repos.Menu.Query),
		List:    menu.NewListHandler(repos.Menu.Query),
	}
}

// newSettingUseCases 初始化系统配置用例
func newSettingUseCases(repos *RepositoriesModule) *SettingUseCases {
	// 创建 JSON Logic 验证器
	validator := validation.NewJSONLogicValidator()

	return &SettingUseCases{
		// Setting handlers
		Create:      setting.NewCreateHandler(repos.Setting.Command, repos.Setting.Query),
		Update:      setting.NewUpdateHandler(repos.Setting.Command, repos.Setting.Query, validator),
		Delete:      setting.NewDeleteHandler(repos.Setting.Command, repos.Setting.Query),
		BatchUpdate: setting.NewBatchUpdateHandler(repos.Setting.Command, repos.Setting.Query, validator),
		Get:         setting.NewGetHandler(repos.Setting.Query),
		List:        setting.NewListHandler(repos.Setting.Query),
		ListSchema:  setting.NewListSchemaHandler(repos.Setting.Query, repos.Setting.CategoryQuery),

		// Category handlers
		CreateCategory: setting.NewCreateCategoryHandler(repos.Setting.CategoryCommand, repos.Setting.CategoryQuery),
		UpdateCategory: setting.NewUpdateCategoryHandler(repos.Setting.CategoryCommand, repos.Setting.CategoryQuery),
		DeleteCategory: setting.NewDeleteCategoryHandler(repos.Setting.CategoryCommand, repos.Setting.CategoryQuery, repos.Setting.Query),
		GetCategory:    setting.NewGetCategoryHandler(repos.Setting.CategoryQuery),
		ListCategories: setting.NewListCategoriesHandler(repos.Setting.CategoryQuery),
	}
}

// newUserSettingUseCases 初始化用户配置用例
func newUserSettingUseCases(repos *RepositoriesModule) *UserSettingUseCases {
	// 创建 JSON Logic 验证器
	validator := validation.NewJSONLogicValidator()

	return &UserSettingUseCases{
		Set:            setting.NewUserSetHandler(repos.Setting.Query, repos.UserSetting.Command, validator),
		BatchSet:       setting.NewUserBatchSetHandler(repos.Setting.Query, repos.UserSetting.Command, validator),
		Reset:          setting.NewUserResetHandler(repos.UserSetting.Command),
		ResetAll:       setting.NewUserResetAllHandler(repos.UserSetting.Command),
		Get:            setting.NewUserGetHandler(repos.Setting.Query, repos.UserSetting.Query),
		List:           setting.NewUserListHandler(repos.Setting.Query, repos.UserSetting.Query),
		ListSchema:     setting.NewUserListSchemaHandler(repos.Setting.Query, repos.UserSetting.Query, repos.Setting.CategoryQuery),
		ListCategories: setting.NewUserListCategoriesHandler(repos.Setting.Query, repos.Setting.CategoryQuery),
	}
}

// newPATUseCases 初始化个人访问令牌用例
func newPATUseCases(repos *RepositoriesModule, services *ServicesModule) *PATUseCases {
	// 获取内部 tokenGenerator（用于 PAT 生成）
	tokenGenerator, ok := services.TokenGenerator.(*authInfra.TokenGenerator)
	if !ok {
		// fallback: 使用 PAT Service 的默认实现
		tokenGenerator = authInfra.NewTokenGenerator()
	}

	return &PATUseCases{
		Create:  pat.NewCreateHandler(repos.PAT.Command, repos.User.Query, tokenGenerator),
		Delete:  pat.NewDeleteHandler(repos.PAT.Command, repos.PAT.Query),
		Disable: pat.NewDisableHandler(repos.PAT.Command, repos.PAT.Query),
		Enable:  pat.NewEnableHandler(repos.PAT.Command, repos.PAT.Query),
		Get:     pat.NewGetHandler(repos.PAT.Query),
		List:    pat.NewListHandler(repos.PAT.Query),
	}
}

// newAuditLogUseCases 初始化审计日志用例
func newAuditLogUseCases(repos *RepositoriesModule) *AuditLogUseCases {
	return &AuditLogUseCases{
		CreateLog: auditlog.NewCreateHandler(repos.AuditLog.Command),
		Get:       auditlog.NewGetHandler(repos.AuditLog.Query),
		List:      auditlog.NewListHandler(repos.AuditLog.Query),
	}
}

// newStatsUseCases 初始化统计用例
func newStatsUseCases(repos *RepositoriesModule) *StatsUseCases {
	return &StatsUseCases{
		GetStats: stats.NewGetStatsHandler(repos.StatsQuery),
	}
}

// newCaptchaUseCases 初始化验证码用例
func newCaptchaUseCases(repos *RepositoriesModule, services *ServicesModule) *CaptchaUseCases {
	return &CaptchaUseCases{
		Generate: captcha.NewGenerateHandler(repos.CaptchaCommand, services.Captcha),
	}
}

// newTwoFAUseCases 初始化双因素认证用例
func newTwoFAUseCases(services *ServicesModule) *TwoFAUseCases {
	return &TwoFAUseCases{
		Setup:        twofa.NewSetupHandler(services.TwoFA),
		VerifyEnable: twofa.NewVerifyEnableHandler(services.TwoFA),
		Disable:      twofa.NewDisableHandler(services.TwoFA),
		GetStatus:    twofa.NewGetStatusHandler(services.TwoFA),
	}
}
