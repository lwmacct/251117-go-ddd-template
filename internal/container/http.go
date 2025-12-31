package container

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"

	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http"
	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/handler"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/cache"
	"github.com/lwmacct/251117-go-ddd-template/internal/config"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/auth"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/health"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/persistence"
)

// HandlersResult 使用 fx.Out 批量返回所有 HTTP 处理器。
type HandlersResult struct {
	fx.Out

	Health           *handler.HealthHandler
	Auth             *handler.AuthHandler
	Captcha          *handler.CaptchaHandler
	AdminUser        *handler.AdminUserHandler
	UserProfile      *handler.UserProfileHandler
	Role             *handler.RoleHandler
	Setting          *handler.SettingHandler
	UserSetting      *handler.UserSettingHandler
	PAT              *handler.PATHandler
	AuditLog         *handler.AuditLogHandler
	Overview         *handler.OverviewHandler
	TwoFA            *handler.TwoFAHandler
	Cache            *handler.CacheHandler
	Operation        *handler.OperationHandler
	Organization     *handler.OrgHandler
	OrgMember        *handler.OrgMemberHandler
	Team             *handler.TeamHandler
	TeamMember       *handler.TeamMemberHandler
	UserOrganization *handler.UserOrgHandler
}

// HTTPModule 提供 HTTP 处理器和路由。
var HTTPModule = fx.Module("http",
	fx.Provide(
		health.NewSystemChecker,
		newAllHandlers,
		newRouter,
	),
)

// handlersParams 聚合创建 Handler 所需的依赖。
type handlersParams struct {
	fx.In

	Config        *config.Config
	AdminCacheSvc cache.AdminCacheService
	HealthChecker *health.SystemChecker
	Auth          *AuthUseCases
	User          *UserUseCases
	Role          *RoleUseCases
	Setting       *SettingUseCases
	UserSetting   *UserSettingUseCases
	PAT           *PATUseCases
	AuditLog      *AuditLogUseCases
	Stats         *StatsUseCases
	Captcha       *CaptchaUseCases
	TwoFA         *TwoFAUseCases
	Organization  *OrganizationUseCases
}

func newAllHandlers(p handlersParams) HandlersResult {
	return HandlersResult{
		Health: handler.NewHealthHandler(p.HealthChecker),
		Auth: handler.NewAuthHandler(
			p.Auth.Login,
			p.Auth.Login2FA,
			p.Auth.Register,
			p.Auth.RefreshToken,
		),
		Captcha: handler.NewCaptchaHandler(p.Captcha.Generate, p.Config.Auth.DevSecret),
		AdminUser: handler.NewAdminUserHandler(
			p.User.Create,
			p.User.Update,
			p.User.Delete,
			p.User.AssignRoles,
			p.User.BatchCreate,
			p.User.Get,
			p.User.List,
		),
		UserProfile: handler.NewUserProfileHandler(
			p.User.Get,
			p.User.Update,
			p.User.ChangePassword,
			p.User.Delete,
		),
		Role: handler.NewRoleHandler(
			p.Role.Create,
			p.Role.Update,
			p.Role.Delete,
			p.Role.SetPermissions,
			p.Role.Get,
			p.Role.List,
		),
		Setting: handler.NewSettingHandler(
			p.Setting.Create,
			p.Setting.Update,
			p.Setting.Delete,
			p.Setting.BatchUpdate,
			p.Setting.Get,
			p.Setting.List,
			p.Setting.ListSettings,
			p.Setting.CreateCategory,
			p.Setting.UpdateCategory,
			p.Setting.DeleteCategory,
			p.Setting.GetCategory,
			p.Setting.ListCategories,
		),
		UserSetting: handler.NewUserSettingHandler(
			p.UserSetting.Set,
			p.UserSetting.BatchSet,
			p.UserSetting.Reset,
			p.UserSetting.ResetAll,
			p.UserSetting.Get,
			p.UserSetting.List,
			p.UserSetting.ListSettings,
			p.UserSetting.ListCategories,
		),
		PAT: handler.NewPATHandler(
			p.PAT.Create,
			p.PAT.Delete,
			p.PAT.Disable,
			p.PAT.Enable,
			p.PAT.Get,
			p.PAT.List,
		),
		AuditLog: handler.NewAuditLogHandler(
			p.AuditLog.List,
			p.AuditLog.Get,
		),
		Overview: handler.NewOverviewHandler(p.Stats.GetStats),
		TwoFA: handler.NewTwoFAHandler(
			p.TwoFA.Setup,
			p.TwoFA.VerifyEnable,
			p.TwoFA.Disable,
			p.TwoFA.GetStatus,
		),
		Cache: handler.NewCacheHandler(
			cache.NewInfoHandler(p.AdminCacheSvc),
			cache.NewScanKeysHandler(p.AdminCacheSvc),
			cache.NewGetKeyHandler(p.AdminCacheSvc),
			cache.NewDeleteHandler(p.AdminCacheSvc),
		),
		Operation: handler.NewOperationHandler(),
		Organization: handler.NewOrgHandler(
			p.Organization.Create,
			p.Organization.Update,
			p.Organization.Delete,
			p.Organization.Get,
			p.Organization.List,
		),
		OrgMember: handler.NewOrgMemberHandler(
			p.Organization.MemberAdd,
			p.Organization.MemberRemove,
			p.Organization.MemberUpdateRole,
			p.Organization.MemberList,
		),
		Team: handler.NewTeamHandler(
			p.Organization.TeamCreate,
			p.Organization.TeamUpdate,
			p.Organization.TeamDelete,
			p.Organization.TeamGet,
			p.Organization.TeamList,
		),
		TeamMember: handler.NewTeamMemberHandler(
			p.Organization.TeamMemberAdd,
			p.Organization.TeamMemberRemove,
			p.Organization.TeamMemberList,
		),
		UserOrganization: handler.NewUserOrgHandler(
			p.Organization.UserOrgs,
			p.Organization.UserTeams,
		),
	}
}

// routerParams 聚合创建路由所需的依赖。
type routerParams struct {
	fx.In

	Config      *config.Config
	RedisClient *redis.Client

	// Services
	JWTManager      *auth.JWTManager
	PATService      *auth.PATService
	PermissionCache *auth.PermissionCacheService

	// UseCases
	AuditLog *AuditLogUseCases

	// Repositories (for middleware)
	MemberRepos     persistence.OrgMemberRepositories
	TeamRepos       persistence.TeamRepositories
	TeamMemberRepos persistence.TeamMemberRepositories

	// Handlers
	Health      *handler.HealthHandler
	Auth        *handler.AuthHandler
	Captcha     *handler.CaptchaHandler
	AdminUser   *handler.AdminUserHandler
	UserProfile *handler.UserProfileHandler
	Role        *handler.RoleHandler
	Setting     *handler.SettingHandler
	UserSetting *handler.UserSettingHandler
	PAT         *handler.PATHandler
	AuditLogH   *handler.AuditLogHandler
	Overview    *handler.OverviewHandler
	TwoFA       *handler.TwoFAHandler
	Cache       *handler.CacheHandler
	Operation   *handler.OperationHandler
	Org         *handler.OrgHandler
	OrgMember   *handler.OrgMemberHandler
	Team        *handler.TeamHandler
	TeamMember  *handler.TeamMemberHandler
	UserOrg     *handler.UserOrgHandler
}

func newRouter(p routerParams) *gin.Engine {
	deps := &http.RouterDependencies{
		Config:                 p.Config,
		RedisClient:            p.RedisClient,
		CreateLogHandler:       p.AuditLog.CreateLog,
		JWTManager:             p.JWTManager,
		PATService:             p.PATService,
		PermissionCacheService: p.PermissionCache,
		OrgMemberQuery:         p.MemberRepos.Query,
		TeamQuery:              p.TeamRepos.Query,
		TeamMemberQuery:        p.TeamMemberRepos.Query,
		HealthHandler:          p.Health,
		AuthHandler:            p.Auth,
		CaptchaHandler:         p.Captcha,
		RoleHandler:            p.Role,
		SettingHandler:         p.Setting,
		UserSettingHandler:     p.UserSetting,
		PATHandler:             p.PAT,
		AuditLogHandler:        p.AuditLogH,
		AdminUserHandler:       p.AdminUser,
		UserProfileHandler:     p.UserProfile,
		OverviewHandler:        p.Overview,
		TwoFAHandler:           p.TwoFA,
		CacheHandler:           p.Cache,
		OperationHandler:       p.Operation,
		OrgHandler:             p.Org,
		OrgMemberHandler:       p.OrgMember,
		TeamHandler:            p.Team,
		TeamMemberHandler:      p.TeamMember,
		UserOrgHandler:         p.UserOrg,
	}

	return http.SetupRouterWithDeps(deps)
}
