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

// RepositoriesModule 聚合所有 CQRS 仓储。
type RepositoriesModule struct {
	// CQRS 仓储（数据库）
	User        persistence.UserRepositories
	AuditLog    persistence.AuditLogRepositories
	Role        persistence.RoleRepositories
	Permission  persistence.PermissionRepositories
	PAT         persistence.PATRepositories
	Menu        persistence.MenuRepositories
	Setting     persistence.SettingRepositories
	UserSetting persistence.UserSettingRepositories
	TwoFA       persistence.TwoFARepositories

	// 特殊仓储（内存）
	CaptchaCommand captcha.CommandRepository
	CaptchaQuery   captcha.QueryRepository

	// 只读仓储
	StatsQuery stats.QueryRepository
}

// RepositoryModule 提供所有仓储实现。
//
// 适当的仓储已装饰缓存层：
//   - User: 缓存查询（GetByIDWithRoles）
//   - Setting: 缓存查询 + 命令，支持多级失效
//   - UserSetting: 缓存查询 + 命令，支持三级失效
var RepositoryModule = fx.Module("repository",
	fx.Provide(
		// 聚合模块
		newRepositoriesModule,

		// 独立仓储，用于细粒度依赖注入
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

// newRepositoriesModule 创建聚合的仓储模块。
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

// --- 独立仓储构造函数 ---

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

// CaptchaRepositoryResult 从单个仓储提供 Command 和 Query 两个接口。
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
