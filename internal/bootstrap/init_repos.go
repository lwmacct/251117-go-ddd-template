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

	return &RepositoriesModule{
		// CQRS 仓储（数据库实现）
		User:        persistence.NewUserRepositories(db),
		AuditLog:    persistence.NewAuditLogRepositories(db),
		Role:        persistence.NewRoleRepositories(db),
		Permission:  persistence.NewPermissionRepositories(db),
		PAT:         persistence.NewPATRepositories(db),
		Menu:        persistence.NewMenuRepositories(db),
		Setting:     settingRepos,
		UserSetting: persistence.NewUserSettingRepositories(db),
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
	cachedCommand := persistence.NewCachedSettingCommandRepository(rawRepos.Command, cacheServices.Setting)

	return persistence.SettingRepositories{
		Command:         cachedCommand,
		Query:           cachedQuery,
		CategoryQuery:   rawRepos.CategoryQuery, // Category 暂不缓存
		CategoryCommand: rawRepos.CategoryCommand,
	}
}
