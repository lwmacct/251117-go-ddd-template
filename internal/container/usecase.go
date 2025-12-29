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
	domain_auth "github.com/lwmacct/251117-go-ddd-template/internal/domain/auth"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/captcha"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/event"
	domain_stats "github.com/lwmacct/251117-go-ddd-template/internal/domain/stats"
	infra_auth "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/auth"
	infra_captcha "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/captcha"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/persistence"
	infra_twofa "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/twofa"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/validation"
)

// --- 用例模块结构体 ---

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
var UseCaseModule = fx.Module("usecase",
	fx.Provide(
		newAuditLogUseCases,
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
	),
)

// --- 构造函数 ---

func newAuditLogUseCases(repos persistence.AuditLogRepositories) *AuditLogUseCases {
	return &AuditLogUseCases{
		CreateLog: auditlog.NewCreateHandler(repos.Command),
		Get:       auditlog.NewGetHandler(repos.Query),
		List:      auditlog.NewListHandler(repos.Query),
	}
}

// authUseCasesParams 聚合 Auth 用例所需的依赖。
type authUseCasesParams struct {
	fx.In

	UserRepos      persistence.UserRepositories
	CaptchaCommand captcha.CommandRepository
	TwoFARepos     persistence.TwoFARepositories
	AuthSvc        domain_auth.Service
	LoginSession   *infra_auth.LoginSessionService
	TwoFASvc       *infra_twofa.Service
	AuditLog       *AuditLogUseCases
}

func newAuthUseCases(p authUseCasesParams) *AuthUseCases {
	return &AuthUseCases{
		Login:        auth.NewLoginHandler(p.UserRepos.Query, p.CaptchaCommand, p.TwoFARepos.Query, p.AuthSvc, p.LoginSession, p.AuditLog.CreateLog),
		Login2FA:     auth.NewLogin2FAHandler(p.UserRepos.Query, p.AuthSvc, p.LoginSession, p.TwoFASvc, p.AuditLog.CreateLog),
		Register:     auth.NewRegisterHandler(p.UserRepos.Command, p.UserRepos.Query, p.AuthSvc),
		RefreshToken: auth.NewRefreshTokenHandler(p.UserRepos.Query, p.AuthSvc, p.AuditLog.CreateLog),
	}
}

// userUseCasesParams 聚合 User 用例所需的依赖。
type userUseCasesParams struct {
	fx.In

	UserRepos persistence.UserRepositories
	AuthSvc   domain_auth.Service
	EventBus  event.EventBus
}

func newUserUseCases(p userUseCasesParams) *UserUseCases {
	return &UserUseCases{
		Create:         user.NewCreateHandler(p.UserRepos.Command, p.UserRepos.Query, p.AuthSvc),
		Update:         user.NewUpdateHandler(p.UserRepos.Command, p.UserRepos.Query),
		Delete:         user.NewDeleteHandler(p.UserRepos.Command, p.UserRepos.Query, p.EventBus),
		AssignRoles:    user.NewAssignRolesHandler(p.UserRepos.Command, p.UserRepos.Query, p.EventBus),
		ChangePassword: user.NewChangePasswordHandler(p.UserRepos.Command, p.UserRepos.Query, p.AuthSvc),
		BatchCreate:    user.NewBatchCreateHandler(p.UserRepos.Command, p.UserRepos.Query, p.AuthSvc),
		Get:            user.NewGetHandler(p.UserRepos.Query),
		List:           user.NewListHandler(p.UserRepos.Query),
	}
}

// roleUseCasesParams 聚合 Role 用例所需的依赖。
type roleUseCasesParams struct {
	fx.In

	RoleRepos       persistence.RoleRepositories
	PermissionRepos persistence.PermissionRepositories
	EventBus        event.EventBus
}

func newRoleUseCases(p roleUseCasesParams) *RoleUseCases {
	return &RoleUseCases{
		Create:          role.NewCreateHandler(p.RoleRepos.Command, p.RoleRepos.Query),
		Update:          role.NewUpdateHandler(p.RoleRepos.Command, p.RoleRepos.Query),
		Delete:          role.NewDeleteHandler(p.RoleRepos.Command, p.RoleRepos.Query),
		SetPermissions:  role.NewSetPermissionsHandler(p.RoleRepos.Command, p.RoleRepos.Query, p.PermissionRepos.Query, p.EventBus),
		Get:             role.NewGetHandler(p.RoleRepos.Query),
		List:            role.NewListHandler(p.RoleRepos.Query),
		ListPermissions: role.NewListPermissionsHandler(p.PermissionRepos.Query),
	}
}

func newMenuUseCases(repos persistence.MenuRepositories) *MenuUseCases {
	return &MenuUseCases{
		Create:  menu.NewCreateHandler(repos.Command, repos.Query),
		Update:  menu.NewUpdateHandler(repos.Command, repos.Query),
		Delete:  menu.NewDeleteHandler(repos.Command, repos.Query),
		Reorder: menu.NewReorderHandler(repos.Command, repos.Query),
		Get:     menu.NewGetHandler(repos.Query),
		List:    menu.NewListHandler(repos.Query),
	}
}

// settingUseCasesParams 聚合 Setting 用例所需的依赖。
type settingUseCasesParams struct {
	fx.In

	SettingRepos persistence.SettingRepositories
	SchemaCache  setting.SchemaCacheService
}

func newSettingUseCases(p settingUseCasesParams) *SettingUseCases {
	validator := validation.NewJSONLogicValidator()

	return &SettingUseCases{
		Create:         setting.NewCreateHandler(p.SettingRepos.Command, p.SettingRepos.Query, p.SchemaCache),
		Update:         setting.NewUpdateHandler(p.SettingRepos.Command, p.SettingRepos.Query, validator, p.SchemaCache),
		Delete:         setting.NewDeleteHandler(p.SettingRepos.Command, p.SettingRepos.Query, p.SchemaCache),
		BatchUpdate:    setting.NewBatchUpdateHandler(p.SettingRepos.Command, p.SettingRepos.Query, validator, p.SchemaCache),
		Get:            setting.NewGetHandler(p.SettingRepos.Query),
		List:           setting.NewListHandler(p.SettingRepos.Query),
		ListSchema:     setting.NewListSchemaHandler(p.SettingRepos.Query, p.SettingRepos.CategoryQuery, p.SchemaCache),
		CreateCategory: setting.NewCreateCategoryHandler(p.SettingRepos.CategoryCommand, p.SettingRepos.CategoryQuery),
		UpdateCategory: setting.NewUpdateCategoryHandler(p.SettingRepos.CategoryCommand, p.SettingRepos.CategoryQuery),
		DeleteCategory: setting.NewDeleteCategoryHandler(p.SettingRepos.CategoryCommand, p.SettingRepos.CategoryQuery, p.SettingRepos.Query),
		GetCategory:    setting.NewGetCategoryHandler(p.SettingRepos.CategoryQuery),
		ListCategories: setting.NewListCategoriesHandler(p.SettingRepos.CategoryQuery),
	}
}

// userSettingUseCasesParams 聚合 UserSetting 用例所需的依赖。
type userSettingUseCasesParams struct {
	fx.In

	SettingRepos     persistence.SettingRepositories
	UserSettingRepos persistence.UserSettingRepositories
	SchemaCache      setting.SchemaCacheService
}

func newUserSettingUseCases(p userSettingUseCasesParams) *UserSettingUseCases {
	validator := validation.NewJSONLogicValidator()

	return &UserSettingUseCases{
		Set:            setting.NewUserSetHandler(p.SettingRepos.Query, p.UserSettingRepos.Command, validator),
		BatchSet:       setting.NewUserBatchSetHandler(p.SettingRepos.Query, p.UserSettingRepos.Command, validator),
		Reset:          setting.NewUserResetHandler(p.UserSettingRepos.Command),
		ResetAll:       setting.NewUserResetAllHandler(p.UserSettingRepos.Command),
		Get:            setting.NewUserGetHandler(p.SettingRepos.Query, p.UserSettingRepos.Query),
		List:           setting.NewUserListHandler(p.SettingRepos.Query, p.UserSettingRepos.Query),
		ListSchema:     setting.NewUserListSchemaHandler(p.SettingRepos.Query, p.UserSettingRepos.Query, p.SettingRepos.CategoryQuery, p.SchemaCache),
		ListCategories: setting.NewUserListCategoriesHandler(p.SettingRepos.Query, p.SettingRepos.CategoryQuery, p.SchemaCache),
	}
}

// patUseCasesParams 聚合 PAT 用例所需的依赖。
type patUseCasesParams struct {
	fx.In

	PATRepos  persistence.PATRepositories
	UserRepos persistence.UserRepositories
	TokenGen  *infra_auth.TokenGenerator
}

func newPATUseCases(p patUseCasesParams) *PATUseCases {
	return &PATUseCases{
		Create:  pat.NewCreateHandler(p.PATRepos.Command, p.UserRepos.Query, p.TokenGen),
		Delete:  pat.NewDeleteHandler(p.PATRepos.Command, p.PATRepos.Query),
		Disable: pat.NewDisableHandler(p.PATRepos.Command, p.PATRepos.Query),
		Enable:  pat.NewEnableHandler(p.PATRepos.Command, p.PATRepos.Query),
		Get:     pat.NewGetHandler(p.PATRepos.Query),
		List:    pat.NewListHandler(p.PATRepos.Query),
	}
}

func newStatsUseCases(statsQuery domain_stats.QueryRepository) *StatsUseCases {
	return &StatsUseCases{
		GetStats: stats.NewGetStatsHandler(statsQuery),
	}
}

func newCaptchaUseCases(
	captchaCommand captcha.CommandRepository,
	captchaSvc *infra_captcha.Service,
) *CaptchaUseCases {
	return &CaptchaUseCases{
		Generate: app_captcha.NewGenerateHandler(captchaCommand, captchaSvc),
	}
}

func newTwoFAUseCases(twofaSvc *infra_twofa.Service) *TwoFAUseCases {
	return &TwoFAUseCases{
		Setup:        twofa.NewSetupHandler(twofaSvc),
		VerifyEnable: twofa.NewVerifyEnableHandler(twofaSvc),
		Disable:      twofa.NewDisableHandler(twofaSvc),
		GetStatus:    twofa.NewGetStatusHandler(twofaSvc),
	}
}
