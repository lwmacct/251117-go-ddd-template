package container

import (
	"go.uber.org/fx"

	"github.com/lwmacct/251117-go-ddd-template/internal/application/audit"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/auth"
	app_captcha "github.com/lwmacct/251117-go-ddd-template/internal/application/captcha"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/org"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/pat"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/product"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/role"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/setting"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/stats"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/task"
	app_twofa "github.com/lwmacct/251117-go-ddd-template/internal/application/twofa"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/user"
	domain_auth "github.com/lwmacct/251117-go-ddd-template/internal/domain/auth"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/captcha"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/event"
	domain_stats "github.com/lwmacct/251117-go-ddd-template/internal/domain/stats"
	domain_twofa "github.com/lwmacct/251117-go-ddd-template/internal/domain/twofa"
	infra_auth "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/auth"
	infra_captcha "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/captcha"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/persistence"
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
	Create         *role.CreateHandler
	Update         *role.UpdateHandler
	Delete         *role.DeleteHandler
	SetPermissions *role.SetPermissionsHandler
	Get            *role.GetHandler
	List           *role.ListHandler
}

type SettingUseCases struct {
	Create         *setting.CreateHandler
	Update         *setting.UpdateHandler
	Delete         *setting.DeleteHandler
	BatchUpdate    *setting.BatchUpdateHandler
	Get            *setting.GetHandler
	List           *setting.ListHandler
	ListSettings   *setting.ListSettingsHandler
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
	ListSettings   *setting.UserListSettingsHandler
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

type AuditUseCases struct {
	CreateLog *audit.CreateHandler
	Get       *audit.GetHandler
	List      *audit.ListHandler
}

type StatsUseCases struct {
	GetStats *stats.GetStatsHandler
}

type CaptchaUseCases struct {
	Generate *app_captcha.GenerateHandler
}

type TwoFAUseCases struct {
	Setup        *app_twofa.SetupHandler
	VerifyEnable *app_twofa.VerifyEnableHandler
	Disable      *app_twofa.DisableHandler
	GetStatus    *app_twofa.GetStatusHandler
}

// OrganizationUseCases 组织相关用例处理器
type OrganizationUseCases struct {
	// Organization
	Create *org.CreateHandler
	Update *org.UpdateHandler
	Delete *org.DeleteHandler
	Get    *org.GetHandler
	List   *org.ListHandler

	// Member
	MemberAdd        *org.MemberAddHandler
	MemberRemove     *org.MemberRemoveHandler
	MemberUpdateRole *org.MemberUpdateRoleHandler
	MemberList       *org.MemberListHandler

	// Team
	TeamCreate *org.TeamCreateHandler
	TeamUpdate *org.TeamUpdateHandler
	TeamDelete *org.TeamDeleteHandler
	TeamGet    *org.TeamGetHandler
	TeamList   *org.TeamListHandler

	// Team Member
	TeamMemberAdd    *org.TeamMemberAddHandler
	TeamMemberRemove *org.TeamMemberRemoveHandler
	TeamMemberList   *org.TeamMemberListHandler

	// User View
	UserOrgs  *org.UserOrgsHandler
	UserTeams *org.UserTeamsHandler
}

// ProductUseCases 产品相关用例处理器
type ProductUseCases struct {
	Create *product.CreateHandler
	Update *product.UpdateHandler
	Delete *product.DeleteHandler
	Get    *product.GetHandler
	List   *product.ListHandler
}

// TaskUseCases 团队任务相关用例处理器
type TaskUseCases struct {
	Create *task.CreateHandler
	Update *task.UpdateHandler
	Delete *task.DeleteHandler
	Get    *task.GetHandler
	List   *task.ListHandler
}

// --- Fx 模块 ---

// UseCaseModule 提供按领域组织的所有用例处理器。
var UseCaseModule = fx.Module("usecase",
	fx.Provide(
		newAuditUseCases,
		newAuthUseCases,
		newUserUseCases,
		newRoleUseCases,
		newSettingUseCases,
		newUserSettingUseCases,
		newPATUseCases,
		newStatsUseCases,
		newCaptchaUseCases,
		newTwoFAUseCases,
		newOrganizationUseCases,
		newProductUseCases,
		newTaskUseCases,
	),
)

// --- 构造函数 ---

func newAuditUseCases(repos persistence.AuditRepositories) *AuditUseCases {
	return &AuditUseCases{
		CreateLog: audit.NewCreateHandler(repos.Command),
		Get:       audit.NewGetHandler(repos.Query),
		List:      audit.NewListHandler(repos.Query),
	}
}

// authUseCasesParams 聚合 Auth 用例所需的依赖。
type authUseCasesParams struct {
	fx.In

	UserRepos      persistence.UserRepositories
	CaptchaCommand captcha.CommandRepository
	TwoFARepos     persistence.TwoFARepositories
	AuthSvc        domain_auth.Service
	LoginSession   domain_auth.SessionService
	TwoFASvc       domain_twofa.Service
	Audit          *AuditUseCases
}

func newAuthUseCases(p authUseCasesParams) *AuthUseCases {
	return &AuthUseCases{
		Login:        auth.NewLoginHandler(p.UserRepos.Query, p.CaptchaCommand, p.TwoFARepos.Query, p.AuthSvc, p.LoginSession, p.Audit.CreateLog),
		Login2FA:     auth.NewLogin2FAHandler(p.UserRepos.Query, p.AuthSvc, p.LoginSession, p.TwoFASvc, p.Audit.CreateLog),
		Register:     auth.NewRegisterHandler(p.UserRepos.Command, p.UserRepos.Query, p.AuthSvc),
		RefreshToken: auth.NewRefreshTokenHandler(p.UserRepos.Query, p.AuthSvc, p.Audit.CreateLog),
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

	RoleRepos persistence.RoleRepositories
	EventBus  event.EventBus
}

func newRoleUseCases(p roleUseCasesParams) *RoleUseCases {
	return &RoleUseCases{
		Create:         role.NewCreateHandler(p.RoleRepos.Command, p.RoleRepos.Query),
		Update:         role.NewUpdateHandler(p.RoleRepos.Command, p.RoleRepos.Query),
		Delete:         role.NewDeleteHandler(p.RoleRepos.Command, p.RoleRepos.Query),
		SetPermissions: role.NewSetPermissionsHandler(p.RoleRepos.Command, p.RoleRepos.Query, p.EventBus),
		Get:            role.NewGetHandler(p.RoleRepos.Query),
		List:           role.NewListHandler(p.RoleRepos.Query),
	}
}

// settingUseCasesParams 聚合 Setting 用例所需的依赖。
type settingUseCasesParams struct {
	fx.In

	SettingRepos  persistence.SettingRepositories
	SettingsCache setting.SettingsCacheService
}

func newSettingUseCases(p settingUseCasesParams) *SettingUseCases {
	validator := validation.NewJSONLogicValidator()

	return &SettingUseCases{
		Create:         setting.NewCreateHandler(p.SettingRepos.Command, p.SettingRepos.Query, p.SettingsCache),
		Update:         setting.NewUpdateHandler(p.SettingRepos.Command, p.SettingRepos.Query, validator, p.SettingsCache),
		Delete:         setting.NewDeleteHandler(p.SettingRepos.Command, p.SettingRepos.Query, p.SettingsCache),
		BatchUpdate:    setting.NewBatchUpdateHandler(p.SettingRepos.Command, p.SettingRepos.Query, validator, p.SettingsCache),
		Get:            setting.NewGetHandler(p.SettingRepos.Query),
		List:           setting.NewListHandler(p.SettingRepos.Query),
		ListSettings:   setting.NewListSettingsHandler(p.SettingRepos.Query, p.SettingRepos.CategoryQuery, p.SettingsCache),
		CreateCategory: setting.NewCreateCategoryHandler(p.SettingRepos.CategoryCommand, p.SettingRepos.CategoryQuery, p.SettingsCache),
		UpdateCategory: setting.NewUpdateCategoryHandler(p.SettingRepos.CategoryCommand, p.SettingRepos.CategoryQuery, p.SettingsCache),
		DeleteCategory: setting.NewDeleteCategoryHandler(p.SettingRepos.CategoryCommand, p.SettingRepos.CategoryQuery, p.SettingRepos.Query, p.SettingsCache),
		GetCategory:    setting.NewGetCategoryHandler(p.SettingRepos.CategoryQuery),
		ListCategories: setting.NewListCategoriesHandler(p.SettingRepos.CategoryQuery),
	}
}

// userSettingUseCasesParams 聚合 UserSetting 用例所需的依赖。
type userSettingUseCasesParams struct {
	fx.In

	SettingRepos     persistence.SettingRepositories
	UserSettingRepos persistence.UserSettingRepositories
	SettingsCache    setting.SettingsCacheService
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
		ListSettings:   setting.NewUserListSettingsHandler(p.SettingRepos.Query, p.UserSettingRepos.Query, p.SettingRepos.CategoryQuery, p.SettingsCache),
		ListCategories: setting.NewUserListCategoriesHandler(p.SettingRepos.Query, p.SettingRepos.CategoryQuery, p.SettingsCache),
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

func newTwoFAUseCases(twofaSvc domain_twofa.Service) *TwoFAUseCases {
	return &TwoFAUseCases{
		Setup:        app_twofa.NewSetupHandler(twofaSvc),
		VerifyEnable: app_twofa.NewVerifyEnableHandler(twofaSvc),
		Disable:      app_twofa.NewDisableHandler(twofaSvc),
		GetStatus:    app_twofa.NewGetStatusHandler(twofaSvc),
	}
}

// organizationUseCasesParams 聚合 Organization 用例所需的依赖。
type organizationUseCasesParams struct {
	fx.In

	OrgRepos persistence.OrganizationRepositories
}

func newOrganizationUseCases(p organizationUseCasesParams) *OrganizationUseCases {
	return &OrganizationUseCases{
		// Organization
		Create: org.NewCreateHandler(p.OrgRepos.Command, p.OrgRepos.Query, p.OrgRepos.MemberCommand),
		Update: org.NewUpdateHandler(p.OrgRepos.Command, p.OrgRepos.Query),
		Delete: org.NewDeleteHandler(
			p.OrgRepos.Command,
			p.OrgRepos.Query,
			p.OrgRepos.MemberQuery,
			p.OrgRepos.MemberCommand,
			p.OrgRepos.TeamQuery,
			p.OrgRepos.TeamCommand,
			p.OrgRepos.TeamMemberQuery,
			p.OrgRepos.TeamMemberCommand,
		),
		Get:  org.NewGetHandler(p.OrgRepos.Query),
		List: org.NewListHandler(p.OrgRepos.Query),

		// Member
		MemberAdd:        org.NewMemberAddHandler(p.OrgRepos.MemberCommand, p.OrgRepos.MemberQuery, p.OrgRepos.Query),
		MemberRemove:     org.NewMemberRemoveHandler(p.OrgRepos.MemberCommand, p.OrgRepos.MemberQuery),
		MemberUpdateRole: org.NewMemberUpdateRoleHandler(p.OrgRepos.MemberCommand, p.OrgRepos.MemberQuery),
		MemberList:       org.NewMemberListHandler(p.OrgRepos.MemberQuery),

		// Team
		TeamCreate: org.NewTeamCreateHandler(p.OrgRepos.TeamCommand, p.OrgRepos.TeamQuery, p.OrgRepos.Query, p.OrgRepos.TeamMemberCommand),
		TeamUpdate: org.NewTeamUpdateHandler(p.OrgRepos.TeamCommand, p.OrgRepos.TeamQuery),
		TeamDelete: org.NewTeamDeleteHandler(
			p.OrgRepos.TeamCommand,
			p.OrgRepos.TeamQuery,
			p.OrgRepos.TeamMemberQuery,
			p.OrgRepos.TeamMemberCommand,
		),
		TeamGet:  org.NewTeamGetHandler(p.OrgRepos.TeamQuery),
		TeamList: org.NewTeamListHandler(p.OrgRepos.TeamQuery),

		// Team Member
		TeamMemberAdd:    org.NewTeamMemberAddHandler(p.OrgRepos.TeamMemberCommand, p.OrgRepos.TeamMemberQuery, p.OrgRepos.TeamQuery, p.OrgRepos.MemberQuery),
		TeamMemberRemove: org.NewTeamMemberRemoveHandler(p.OrgRepos.TeamMemberCommand, p.OrgRepos.TeamQuery),
		TeamMemberList:   org.NewTeamMemberListHandler(p.OrgRepos.TeamMemberQuery, p.OrgRepos.TeamQuery),

		// User View
		UserOrgs:  org.NewUserOrgsHandler(p.OrgRepos.MemberQuery, p.OrgRepos.Query),
		UserTeams: org.NewUserTeamsHandler(p.OrgRepos.TeamMemberQuery, p.OrgRepos.TeamQuery, p.OrgRepos.Query),
	}
}

func newProductUseCases(repos persistence.ProductRepositories) *ProductUseCases {
	return &ProductUseCases{
		Create: product.NewCreateHandler(repos.Command, repos.Query),
		Update: product.NewUpdateHandler(repos.Command, repos.Query),
		Delete: product.NewDeleteHandler(repos.Command),
		Get:    product.NewGetHandler(repos.Query),
		List:   product.NewListHandler(repos.Query),
	}
}

func newTaskUseCases(repos persistence.TaskRepositories) *TaskUseCases {
	return &TaskUseCases{
		Create: task.NewCreateHandler(repos.Command),
		Update: task.NewUpdateHandler(repos.Command, repos.Query),
		Delete: task.NewDeleteHandler(repos.Command, repos.Query),
		Get:    task.NewGetHandler(repos.Query),
		List:   task.NewListHandler(repos.Query),
	}
}
