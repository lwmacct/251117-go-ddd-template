package container

import (
	"go.uber.org/fx"
	"gorm.io/gorm"

	"github.com/lwmacct/251117-go-ddd-template/internal/application/setting"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/cache"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/captcha"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/stats"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/persistence"

	infracaptcha "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/captcha"
	infrastats "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/stats"
)

// RepositoryModule 提供所有仓储实现。
//
// 适当的仓储已装饰缓存层：
//   - User: 缓存查询（GetByIDWithRoles）
//   - Setting: 缓存查询 + 命令，支持多级失效
//   - UserSetting: 缓存查询 + 命令，支持三级失效
var RepositoryModule = fx.Module("repository",
	fx.Provide(
		// 直接使用 persistence 构造函数（无需包装）
		persistence.NewAuditLogRepositories,
		persistence.NewRoleRepositories,
		persistence.NewPermissionRepositories,
		persistence.NewPATRepositories,
		persistence.NewMenuRepositories,
		persistence.NewTwoFARepositories,

		// 带缓存装饰的仓储
		newUserRepositoriesWithCache,
		newSettingRepositoriesWithCache,
		newUserSettingRepositoriesWithCache,

		// 特殊仓储
		newCaptchaRepository,
		newStatsQueryRepository,
	),
)

// --- 带缓存装饰的仓储构造函数 ---

func newUserRepositoriesWithCache(
	db *gorm.DB,
	userWithRolesCache cache.UserWithRolesCacheService,
) persistence.UserRepositories {
	rawRepos := persistence.NewUserRepositories(db)
	cachedQuery := persistence.NewCachedUserQueryRepository(rawRepos.Query, userWithRolesCache)
	return persistence.UserRepositories{
		Command: rawRepos.Command,
		Query:   cachedQuery,
	}
}

// settingRepositoriesParams 聚合 Setting 仓储所需的缓存服务。
type settingRepositoriesParams struct {
	fx.In

	DB               *gorm.DB
	SettingCache     cache.SettingCacheService
	CategoryCache    cache.SettingCategoryCacheService
	UserSettingCache cache.UserSettingCacheService
}

func newSettingRepositoriesWithCache(p settingRepositoriesParams) persistence.SettingRepositories {
	rawRepos := persistence.NewSettingRepositories(p.DB)

	cachedQuery := persistence.NewCachedSettingQueryRepository(rawRepos.Query, p.SettingCache)
	cachedCommand := persistence.NewCachedSettingCommandRepositoryWithUserCache(
		rawRepos.Command,
		p.SettingCache,
		p.UserSettingCache,
	)
	cachedCategoryQuery := persistence.NewCachedSettingCategoryQueryRepository(
		rawRepos.CategoryQuery,
		p.CategoryCache,
	)
	cachedCategoryCommand := persistence.NewCachedSettingCategoryCommandRepository(
		rawRepos.CategoryCommand,
		p.CategoryCache,
	)

	return persistence.SettingRepositories{
		Command:         cachedCommand,
		Query:           cachedQuery,
		CategoryQuery:   cachedCategoryQuery,
		CategoryCommand: cachedCategoryCommand,
	}
}

// userSettingRepositoriesParams 聚合 UserSetting 仓储所需的缓存服务。
type userSettingRepositoriesParams struct {
	fx.In

	DB                    *gorm.DB
	UserSettingQueryCache cache.UserSettingQueryCacheService
	UserSettingCache      cache.UserSettingCacheService
	SchemaCache           setting.SchemaCacheService
}

func newUserSettingRepositoriesWithCache(p userSettingRepositoriesParams) persistence.UserSettingRepositories {
	rawRepos := persistence.NewUserSettingRepositories(p.DB)

	cachedQuery := persistence.NewCachedUserSettingQueryRepository(
		rawRepos.Query,
		p.UserSettingQueryCache,
	)
	cachedCommand := persistence.NewCachedUserSettingCommandRepository(
		rawRepos.Command,
		p.UserSettingQueryCache,
		p.UserSettingCache,
		p.SchemaCache,
	)

	return persistence.UserSettingRepositories{
		Command: cachedCommand,
		Query:   cachedQuery,
	}
}

// --- 特殊仓储 ---

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
