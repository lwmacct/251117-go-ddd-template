package container

import (
	"go.uber.org/fx"
	"gorm.io/gorm"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/captcha"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/stats"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/persistence"

	infracaptcha "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/captcha"
	infrastats "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/stats"
)

// RepositoriesModule aggregates all CQRS repositories.
type RepositoriesModule struct {
	// CQRS repositories (database)
	User        persistence.UserRepositories
	AuditLog    persistence.AuditLogRepositories
	Role        persistence.RoleRepositories
	Permission  persistence.PermissionRepositories
	PAT         persistence.PATRepositories
	Menu        persistence.MenuRepositories
	Setting     persistence.SettingRepositories
	UserSetting persistence.UserSettingRepositories
	TwoFA       persistence.TwoFARepositories

	// Special repositories (in-memory)
	CaptchaCommand captcha.CommandRepository
	CaptchaQuery   captcha.QueryRepository

	// Read-only repositories
	StatsQuery stats.QueryRepository
}

// RepositoryModule provides all repository implementations.
//
// Repositories are decorated with cache layers where appropriate:
//   - User: cached query (GetByIDWithRoles)
//   - Setting: cached query + command with multi-level invalidation
//   - UserSetting: cached query + command with three-level invalidation
var RepositoryModule = fx.Module("repository",
	fx.Provide(
		// Aggregated module
		newRepositoriesModule,

		// Individual repositories for direct injection
		newUserRepositoriesWithCache,
		newAuditLogRepositories,
		newRoleRepositories,
		newPermissionRepositories,
		newPATRepositories,
		newMenuRepositories,
		newSettingRepositoriesWithCache,
		newUserSettingRepositoriesWithCache,
		newTwoFARepositories,
		newCaptchaRepository,
		newStatsQueryRepository,
	),
)

// newRepositoriesModule creates the aggregated repositories module.
func newRepositoriesModule(
	user persistence.UserRepositories,
	auditLog persistence.AuditLogRepositories,
	role persistence.RoleRepositories,
	permission persistence.PermissionRepositories,
	pat persistence.PATRepositories,
	menu persistence.MenuRepositories,
	setting persistence.SettingRepositories,
	userSetting persistence.UserSettingRepositories,
	twoFA persistence.TwoFARepositories,
	captchaCommand captcha.CommandRepository,
	captchaQuery captcha.QueryRepository,
	statsQuery stats.QueryRepository,
) *RepositoriesModule {
	return &RepositoriesModule{
		User:           user,
		AuditLog:       auditLog,
		Role:           role,
		Permission:     permission,
		PAT:            pat,
		Menu:           menu,
		Setting:        setting,
		UserSetting:    userSetting,
		TwoFA:          twoFA,
		CaptchaCommand: captchaCommand,
		CaptchaQuery:   captchaQuery,
		StatsQuery:     statsQuery,
	}
}

// --- Individual repository constructors ---

func newUserRepositoriesWithCache(db *gorm.DB, cache *CacheServicesModule) persistence.UserRepositories {
	rawRepos := persistence.NewUserRepositories(db)
	cachedQuery := persistence.NewCachedUserQueryRepository(rawRepos.Query, cache.UserWithRoles)
	return persistence.UserRepositories{
		Command: rawRepos.Command,
		Query:   cachedQuery,
	}
}

func newAuditLogRepositories(db *gorm.DB) persistence.AuditLogRepositories {
	return persistence.NewAuditLogRepositories(db)
}

func newRoleRepositories(db *gorm.DB) persistence.RoleRepositories {
	return persistence.NewRoleRepositories(db)
}

func newPermissionRepositories(db *gorm.DB) persistence.PermissionRepositories {
	return persistence.NewPermissionRepositories(db)
}

func newPATRepositories(db *gorm.DB) persistence.PATRepositories {
	return persistence.NewPATRepositories(db)
}

func newMenuRepositories(db *gorm.DB) persistence.MenuRepositories {
	return persistence.NewMenuRepositories(db)
}

func newSettingRepositoriesWithCache(db *gorm.DB, cache *CacheServicesModule) persistence.SettingRepositories {
	rawRepos := persistence.NewSettingRepositories(db)

	cachedQuery := persistence.NewCachedSettingQueryRepository(rawRepos.Query, cache.Setting)
	cachedCommand := persistence.NewCachedSettingCommandRepositoryWithUserCache(
		rawRepos.Command,
		cache.Setting,
		cache.UserSetting,
	)
	cachedCategoryQuery := persistence.NewCachedSettingCategoryQueryRepository(
		rawRepos.CategoryQuery,
		cache.SettingCategory,
	)

	return persistence.SettingRepositories{
		Command:         cachedCommand,
		Query:           cachedQuery,
		CategoryQuery:   cachedCategoryQuery,
		CategoryCommand: rawRepos.CategoryCommand,
	}
}

func newUserSettingRepositoriesWithCache(db *gorm.DB, cache *CacheServicesModule) persistence.UserSettingRepositories {
	rawRepos := persistence.NewUserSettingRepositories(db)

	cachedQuery := persistence.NewCachedUserSettingQueryRepository(
		rawRepos.Query,
		cache.UserSettingQuery,
	)
	cachedCommand := persistence.NewCachedUserSettingCommandRepository(
		rawRepos.Command,
		cache.UserSettingQuery,
		cache.UserSetting,
		cache.Schema,
	)

	return persistence.UserSettingRepositories{
		Command: cachedCommand,
		Query:   cachedQuery,
	}
}

func newTwoFARepositories(db *gorm.DB) persistence.TwoFARepositories {
	return persistence.NewTwoFARepositories(db)
}

// CaptchaRepositoryResult provides both Command and Query interfaces from a single repository.
type CaptchaRepositoryResult struct {
	fx.Out

	Command captcha.CommandRepository
	Query   captcha.QueryRepository
}

func newCaptchaRepository() CaptchaRepositoryResult {
	repo := infracaptcha.NewRepository()
	return CaptchaRepositoryResult{
		Command: repo,
		Query:   repo,
	}
}

func newStatsQueryRepository(db *gorm.DB) stats.QueryRepository {
	return infrastats.NewQueryRepository(db)
}
