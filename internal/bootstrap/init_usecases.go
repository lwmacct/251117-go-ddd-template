package bootstrap

import (
	"github.com/lwmacct/251117-go-ddd-template/internal/config"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/event"

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

	authInfra "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/auth"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/redis"
)

// newUseCasesModule 初始化用例模块
// 依赖：RepositoriesModule, ServicesModule, InfrastructureModule, EventBus, Config
func newUseCasesModule(cfg *config.Config, infra *InfrastructureModule, repos *RepositoriesModule, services *ServicesModule, eventBus event.EventBus) *UseCasesModule {
	return &UseCasesModule{
		Auth:     newAuthUseCases(repos, services),
		User:     newUserUseCases(repos, services, eventBus),
		Role:     newRoleUseCases(repos, eventBus),
		Menu:     newMenuUseCases(repos),
		Setting:  newSettingUseCases(repos),
		PAT:      newPATUseCases(repos, services),
		AuditLog: newAuditLogUseCases(repos),
		Stats:    newStatsUseCases(repos),
		Captcha:  newCaptchaUseCases(repos, services),
		TwoFA:    newTwoFAUseCases(services),
		Cache:    newCacheUseCases(infra, cfg),
	}
}

// newAuthUseCases 初始化认证用例
func newAuthUseCases(repos *RepositoriesModule, services *ServicesModule) *AuthUseCases {
	return &AuthUseCases{
		Login:        auth.NewLoginHandler(repos.User.Query, repos.CaptchaCommand, repos.TwoFA.Query, services.Auth, services.LoginSession),
		Login2FA:     auth.NewLogin2FAHandler(repos.User.Query, services.Auth, services.LoginSession, services.TwoFA),
		Register:     auth.NewRegisterHandler(repos.User.Command, repos.User.Query, services.Auth),
		RefreshToken: auth.NewRefreshTokenHandler(repos.User.Query, services.Auth),
	}
}

// newUserUseCases 初始化用户管理用例
func newUserUseCases(repos *RepositoriesModule, services *ServicesModule, eventBus event.EventBus) *UserUseCases {
	return &UserUseCases{
		Create:         user.NewCreateUserHandler(repos.User.Command, repos.User.Query, services.Auth),
		Update:         user.NewUpdateUserHandler(repos.User.Command, repos.User.Query),
		Delete:         user.NewDeleteUserHandler(repos.User.Command, repos.User.Query, eventBus),
		AssignRoles:    user.NewAssignRolesHandler(repos.User.Command, repos.User.Query, eventBus),
		ChangePassword: user.NewChangePasswordHandler(repos.User.Command, repos.User.Query, services.Auth),
		BatchCreate:    user.NewBatchCreateUsersHandler(repos.User.Command, repos.User.Query, services.Auth),
		Get:            user.NewGetUserHandler(repos.User.Query),
		List:           user.NewListUsersHandler(repos.User.Query),
	}
}

// newRoleUseCases 初始化角色管理用例
func newRoleUseCases(repos *RepositoriesModule, eventBus event.EventBus) *RoleUseCases {
	return &RoleUseCases{
		Create:          role.NewCreateRoleHandler(repos.Role.Command, repos.Role.Query),
		Update:          role.NewUpdateRoleHandler(repos.Role.Command, repos.Role.Query),
		Delete:          role.NewDeleteRoleHandler(repos.Role.Command, repos.Role.Query),
		SetPermissions:  role.NewSetPermissionsHandler(repos.Role.Command, repos.Role.Query, repos.Permission.Query, eventBus),
		Get:             role.NewGetRoleHandler(repos.Role.Query),
		List:            role.NewListRolesHandler(repos.Role.Query),
		ListPermissions: role.NewListPermissionsHandler(repos.Permission.Query),
	}
}

// newMenuUseCases 初始化菜单管理用例
func newMenuUseCases(repos *RepositoriesModule) *MenuUseCases {
	return &MenuUseCases{
		Create:  menu.NewCreateMenuHandler(repos.Menu.Command, repos.Menu.Query),
		Update:  menu.NewUpdateMenuHandler(repos.Menu.Command, repos.Menu.Query),
		Delete:  menu.NewDeleteMenuHandler(repos.Menu.Command, repos.Menu.Query),
		Reorder: menu.NewReorderMenusHandler(repos.Menu.Command, repos.Menu.Query),
		Get:     menu.NewGetMenuHandler(repos.Menu.Query),
		List:    menu.NewListMenusHandler(repos.Menu.Query),
	}
}

// newSettingUseCases 初始化系统配置用例
func newSettingUseCases(repos *RepositoriesModule) *SettingUseCases {
	return &SettingUseCases{
		Create:      setting.NewCreateSettingHandler(repos.Setting.Command, repos.Setting.Query),
		Update:      setting.NewUpdateSettingHandler(repos.Setting.Command, repos.Setting.Query),
		Delete:      setting.NewDeleteSettingHandler(repos.Setting.Command, repos.Setting.Query),
		BatchUpdate: setting.NewBatchUpdateSettingsHandler(repos.Setting.Command, repos.Setting.Query),
		Get:         setting.NewGetSettingHandler(repos.Setting.Query),
		List:        setting.NewListSettingsHandler(repos.Setting.Query),
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
		Create:  pat.NewCreateTokenHandler(repos.PAT.Command, repos.User.Query, tokenGenerator),
		Delete:  pat.NewDeleteTokenHandler(repos.PAT.Command, repos.PAT.Query),
		Disable: pat.NewDisableTokenHandler(repos.PAT.Command, repos.PAT.Query),
		Enable:  pat.NewEnableTokenHandler(repos.PAT.Command, repos.PAT.Query),
		Get:     pat.NewGetTokenHandler(repos.PAT.Query),
		List:    pat.NewListTokensHandler(repos.PAT.Query),
	}
}

// newAuditLogUseCases 初始化审计日志用例
func newAuditLogUseCases(repos *RepositoriesModule) *AuditLogUseCases {
	return &AuditLogUseCases{
		CreateLog: auditlog.NewCreateLogHandler(repos.AuditLog.Command),
		Get:       auditlog.NewGetLogHandler(repos.AuditLog.Query),
		List:      auditlog.NewListLogsHandler(repos.AuditLog.Query),
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
		Generate: captcha.NewGenerateCaptchaHandler(repos.CaptchaCommand, services.Captcha),
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

// newCacheUseCases 初始化缓存用例（演示用）
func newCacheUseCases(infra *InfrastructureModule, cfg *config.Config) *CacheUseCases {
	// 创建缓存仓储（CQRS 分离）
	cacheCommandRepo := redis.NewCacheCommandRepository(infra.RedisClient, cfg.Data.RedisKeyPrefix)
	cacheQueryRepo := redis.NewCacheQueryRepository(infra.RedisClient, cfg.Data.RedisKeyPrefix)

	return &CacheUseCases{
		Set:    cache.NewSetCacheHandler(cacheCommandRepo),
		Delete: cache.NewDeleteCacheHandler(cacheCommandRepo),
		Get:    cache.NewGetCacheHandler(cacheQueryRepo),
	}
}
