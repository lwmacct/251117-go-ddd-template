package container

import (
	"go.uber.org/fx"

	"github.com/lwmacct/251117-go-ddd-template/internal/application/auditlog"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/auth"
	app_captcha "github.com/lwmacct/251117-go-ddd-template/internal/application/captcha"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/menu"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/pat"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/role"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/setting"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/stats"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/twofa"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/user"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/event"
	infra_auth "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/auth"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/validation"
)

// --- 用例模块结构体 ---

// UseCasesModule 聚合所有用例处理器。
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

type AuthUseCases struct {
	Login        *auth.LoginHandler
	Login2FA     *auth.Login2FAHandler
	Register     *auth.RegisterHandler
	RefreshToken *auth.RefreshTokenHandler
}

type UserUseCases struct {
	Create         *user.CreateHandler
	Update         *user.UpdateHandler
	Delete         *user.DeleteHandler
	AssignRoles    *user.AssignRolesHandler
	ChangePassword *user.ChangePasswordHandler
	BatchCreate    *user.BatchCreateHandler
	Get            *user.GetHandler
	List           *user.ListHandler
}

type RoleUseCases struct {
	Create          *role.CreateHandler
	Update          *role.UpdateHandler
	Delete          *role.DeleteHandler
	SetPermissions  *role.SetPermissionsHandler
	Get             *role.GetHandler
	List            *role.ListHandler
	ListPermissions *role.ListPermissionsHandler
}

type MenuUseCases struct {
	Create  *menu.CreateHandler
	Update  *menu.UpdateHandler
	Delete  *menu.DeleteHandler
	Reorder *menu.ReorderHandler
	Get     *menu.GetHandler
	List    *menu.ListHandler
}

type SettingUseCases struct {
	Create         *setting.CreateHandler
	Update         *setting.UpdateHandler
	Delete         *setting.DeleteHandler
	BatchUpdate    *setting.BatchUpdateHandler
	Get            *setting.GetHandler
	List           *setting.ListHandler
	ListSchema     *setting.ListSchemaHandler
	CreateCategory *setting.CreateCategoryHandler
	UpdateCategory *setting.UpdateCategoryHandler
	DeleteCategory *setting.DeleteCategoryHandler
	GetCategory    *setting.GetCategoryHandler
	ListCategories *setting.ListCategoriesHandler
}

type UserSettingUseCases struct {
	Set            *setting.UserSetHandler
	BatchSet       *setting.UserBatchSetHandler
	Reset          *setting.UserResetHandler
	ResetAll       *setting.UserResetAllHandler
	Get            *setting.UserGetHandler
	List           *setting.UserListHandler
	ListSchema     *setting.UserListSchemaHandler
	ListCategories *setting.UserListCategoriesHandler
}

type PATUseCases struct {
	Create  *pat.CreateHandler
	Delete  *pat.DeleteHandler
	Disable *pat.DisableHandler
	Enable  *pat.EnableHandler
	Get     *pat.GetHandler
	List    *pat.ListHandler
}

type AuditLogUseCases struct {
	CreateLog *auditlog.CreateHandler
	Get       *auditlog.GetHandler
	List      *auditlog.ListHandler
}

type StatsUseCases struct {
	GetStats *stats.GetStatsHandler
}

type CaptchaUseCases struct {
	Generate *app_captcha.GenerateHandler
}

type TwoFAUseCases struct {
	Setup        *twofa.SetupHandler
	VerifyEnable *twofa.VerifyEnableHandler
	Disable      *twofa.DisableHandler
	GetStatus    *twofa.GetStatusHandler
}

// --- Fx 模块 ---

// UseCaseModule 提供按领域组织的所有用例处理器。
//
// 关键依赖：AuditLog 最先创建，因为 Auth 需要它来记录登录日志。
var UseCaseModule = fx.Module("usecase",
	fx.Provide(
		// AuditLog 优先（Auth 依赖它）
		newAuditLogUseCases,

		// 领域用例
		newAuthUseCases,
		newUserUseCases,
		newRoleUseCases,
		newMenuUseCases,
		newSettingUseCases,
		newUserSettingUseCases,
		newPATUseCases,
		newStatsUseCases,
		newCaptchaUseCases,
		newTwoFAUseCases,

		// 聚合模块
		newUseCasesModule,
	),
)

// --- 构造函数 ---

func newAuditLogUseCases(repos *RepositoriesModule) *AuditLogUseCases {
	return &AuditLogUseCases{
		CreateLog: auditlog.NewCreateHandler(repos.AuditLog.Command),
		Get:       auditlog.NewGetHandler(repos.AuditLog.Query),
		List:      auditlog.NewListHandler(repos.AuditLog.Query),
	}
}

func newAuthUseCases(
	repos *RepositoriesModule,
	services *ServicesModule,
	auditLog *AuditLogUseCases,
) *AuthUseCases {
	return &AuthUseCases{
		Login:        auth.NewLoginHandler(repos.User.Query, repos.CaptchaCommand, repos.TwoFA.Query, services.Auth, services.LoginSession, auditLog.CreateLog),
		Login2FA:     auth.NewLogin2FAHandler(repos.User.Query, services.Auth, services.LoginSession, services.TwoFA, auditLog.CreateLog),
		Register:     auth.NewRegisterHandler(repos.User.Command, repos.User.Query, services.Auth),
		RefreshToken: auth.NewRefreshTokenHandler(repos.User.Query, services.Auth),
	}
}

func newUserUseCases(
	repos *RepositoriesModule,
	services *ServicesModule,
	eventBus event.EventBus,
) *UserUseCases {
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

func newSettingUseCases(repos *RepositoriesModule, cache *CacheServicesModule) *SettingUseCases {
	validator := validation.NewJSONLogicValidator()

	return &SettingUseCases{
		Create:         setting.NewCreateHandler(repos.Setting.Command, repos.Setting.Query, cache.Schema),
		Update:         setting.NewUpdateHandler(repos.Setting.Command, repos.Setting.Query, validator, cache.Schema),
		Delete:         setting.NewDeleteHandler(repos.Setting.Command, repos.Setting.Query, cache.Schema),
		BatchUpdate:    setting.NewBatchUpdateHandler(repos.Setting.Command, repos.Setting.Query, validator, cache.Schema),
		Get:            setting.NewGetHandler(repos.Setting.Query),
		List:           setting.NewListHandler(repos.Setting.Query),
		ListSchema:     setting.NewListSchemaHandler(repos.Setting.Query, repos.Setting.CategoryQuery, cache.Schema),
		CreateCategory: setting.NewCreateCategoryHandler(repos.Setting.CategoryCommand, repos.Setting.CategoryQuery),
		UpdateCategory: setting.NewUpdateCategoryHandler(repos.Setting.CategoryCommand, repos.Setting.CategoryQuery),
		DeleteCategory: setting.NewDeleteCategoryHandler(repos.Setting.CategoryCommand, repos.Setting.CategoryQuery, repos.Setting.Query),
		GetCategory:    setting.NewGetCategoryHandler(repos.Setting.CategoryQuery),
		ListCategories: setting.NewListCategoriesHandler(repos.Setting.CategoryQuery),
	}
}

func newUserSettingUseCases(repos *RepositoriesModule, cache *CacheServicesModule) *UserSettingUseCases {
	validator := validation.NewJSONLogicValidator()

	return &UserSettingUseCases{
		Set:            setting.NewUserSetHandler(repos.Setting.Query, repos.UserSetting.Command, validator),
		BatchSet:       setting.NewUserBatchSetHandler(repos.Setting.Query, repos.UserSetting.Command, validator),
		Reset:          setting.NewUserResetHandler(repos.UserSetting.Command),
		ResetAll:       setting.NewUserResetAllHandler(repos.UserSetting.Command),
		Get:            setting.NewUserGetHandler(repos.Setting.Query, repos.UserSetting.Query),
		List:           setting.NewUserListHandler(repos.Setting.Query, repos.UserSetting.Query),
		ListSchema:     setting.NewUserListSchemaHandler(repos.Setting.Query, repos.UserSetting.Query, repos.Setting.CategoryQuery, cache.Schema),
		ListCategories: setting.NewUserListCategoriesHandler(repos.Setting.Query, repos.Setting.CategoryQuery, cache.Schema),
	}
}

func newPATUseCases(repos *RepositoriesModule, tokenGen *infra_auth.TokenGenerator) *PATUseCases {
	return &PATUseCases{
		Create:  pat.NewCreateHandler(repos.PAT.Command, repos.User.Query, tokenGen),
		Delete:  pat.NewDeleteHandler(repos.PAT.Command, repos.PAT.Query),
		Disable: pat.NewDisableHandler(repos.PAT.Command, repos.PAT.Query),
		Enable:  pat.NewEnableHandler(repos.PAT.Command, repos.PAT.Query),
		Get:     pat.NewGetHandler(repos.PAT.Query),
		List:    pat.NewListHandler(repos.PAT.Query),
	}
}

func newStatsUseCases(repos *RepositoriesModule) *StatsUseCases {
	return &StatsUseCases{
		GetStats: stats.NewGetStatsHandler(repos.StatsQuery),
	}
}

func newCaptchaUseCases(repos *RepositoriesModule, services *ServicesModule) *CaptchaUseCases {
	return &CaptchaUseCases{
		Generate: app_captcha.NewGenerateHandler(repos.CaptchaCommand, services.Captcha),
	}
}

func newTwoFAUseCases(services *ServicesModule) *TwoFAUseCases {
	return &TwoFAUseCases{
		Setup:        twofa.NewSetupHandler(services.TwoFA),
		VerifyEnable: twofa.NewVerifyEnableHandler(services.TwoFA),
		Disable:      twofa.NewDisableHandler(services.TwoFA),
		GetStatus:    twofa.NewGetStatusHandler(services.TwoFA),
	}
}

// newUseCasesModule 创建聚合的用例模块。
func newUseCasesModule(
	authUC *AuthUseCases,
	userUC *UserUseCases,
	roleUC *RoleUseCases,
	menuUC *MenuUseCases,
	settingUC *SettingUseCases,
	userSettingUC *UserSettingUseCases,
	patUC *PATUseCases,
	auditLogUC *AuditLogUseCases,
	statsUC *StatsUseCases,
	captchaUC *CaptchaUseCases,
	twofaUC *TwoFAUseCases,
) *UseCasesModule {
	return &UseCasesModule{
		Auth:        authUC,
		User:        userUC,
		Role:        roleUC,
		Menu:        menuUC,
		Setting:     settingUC,
		UserSetting: userSettingUC,
		PAT:         patUC,
		AuditLog:    auditLogUC,
		Stats:       statsUC,
		Captcha:     captchaUC,
		TwoFA:       twofaUC,
	}
}
