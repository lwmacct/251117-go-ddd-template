package bootstrap

import (
	"gorm.io/gorm"

	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/captcha"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/persistence"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/stats"
)

// newRepositoriesModule 初始化仓储模块
// 依赖：InfrastructureModule.DB, CacheServicesModule
func newRepositoriesModule(db *gorm.DB, cacheServices *CacheServicesModule) *RepositoriesModule {
	// Captcha Repository（内存实现，组合接口）
	captchaRepo := captcha.NewRepository()

	// Setting 仓储（带缓存装饰器）
	settingRepos := newSettingRepositoriesWithCache(db, cacheServices)

	// UserSetting 仓储（带缓存装饰器）
	userSettingRepos := newUserSettingRepositoriesWithCache(db, cacheServices)

	return &RepositoriesModule{
		// CQRS 仓储（数据库实现）
		User:        persistence.NewUserRepositories(db),
		AuditLog:    persistence.NewAuditLogRepositories(db),
		Role:        persistence.NewRoleRepositories(db),
		Permission:  persistence.NewPermissionRepositories(db),
		PAT:         persistence.NewPATRepositories(db),
		Menu:        persistence.NewMenuRepositories(db),
		Setting:     settingRepos,
		UserSetting: userSettingRepos,
		TwoFA:       persistence.NewTwoFARepositories(db),

		// 特殊仓储（内存实现）
		CaptchaCommand: captchaRepo,
		CaptchaQuery:   captchaRepo,

		// 只读仓储
		StatsQuery: stats.NewQueryRepository(db),
	}
}

// newSettingRepositoriesWithCache 创建带缓存的 Setting 仓储
func newSettingRepositoriesWithCache(db *gorm.DB, cacheServices *CacheServicesModule) persistence.SettingRepositories {
	// 原始仓储
	rawRepos := persistence.NewSettingRepositories(db)

	// 用装饰器包装 Command 和 Query 仓储
	cachedQuery := persistence.NewCachedSettingQueryRepository(rawRepos.Query, cacheServices.Setting)
	// 使用带用户缓存失效的 Command 仓储
	cachedCommand := persistence.NewCachedSettingCommandRepositoryWithUserCache(
		rawRepos.Command,
		cacheServices.Setting,
		cacheServices.UserSetting,
	)

	// 用装饰器包装 Category Query 仓储
	cachedCategoryQuery := persistence.NewCachedSettingCategoryQueryRepository(
		rawRepos.CategoryQuery,
		cacheServices.SettingCategory,
	)

	return persistence.SettingRepositories{
		Command:         cachedCommand,
		Query:           cachedQuery,
		CategoryQuery:   cachedCategoryQuery,
		CategoryCommand: rawRepos.CategoryCommand,
	}
}

// newUserSettingRepositoriesWithCache 创建带缓存的 UserSetting 仓储
func newUserSettingRepositoriesWithCache(db *gorm.DB, cacheServices *CacheServicesModule) persistence.UserSettingRepositories {
	// 原始仓储
	rawRepos := persistence.NewUserSettingRepositories(db)

	// 用装饰器包装 Query 仓储（读操作时填充缓存）
	cachedQuery := persistence.NewCachedUserSettingQueryRepository(
		rawRepos.Query,
		cacheServices.UserSettingQuery,
	)

	// 用装饰器包装 Command 仓储（写操作时失效三层缓存）
	cachedCommand := persistence.NewCachedUserSettingCommandRepository(
		rawRepos.Command,
		cacheServices.UserSettingQuery,
		cacheServices.UserSetting,
		cacheServices.Schema,
	)

	return persistence.UserSettingRepositories{
		Command: cachedCommand,
		Query:   cachedQuery,
	}
}
