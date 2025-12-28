package container

import (
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"

	"github.com/lwmacct/251117-go-ddd-template/internal/application/setting"
	"github.com/lwmacct/251117-go-ddd-template/internal/config"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/cache"

	infracache "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/cache"
)

// CacheServicesModule 聚合所有缓存服务。
// 定义在 bootstrap 包以保持兼容性。
type CacheServicesModule struct {
	Setting          cache.SettingCacheService
	SettingCategory  cache.SettingCategoryCacheService
	UserSettingQuery cache.UserSettingQueryCacheService
	UserSetting      cache.UserSettingCacheService
	UserWithRoles    cache.UserWithRolesCacheService
	Permission       cache.PermissionCacheService
	Schema           setting.SchemaCacheService
}

// CacheModule 提供所有缓存服务。
//
// 每个服务独立提供以支持细粒度依赖注入，
// 同时聚合到 CacheServicesModule 以保持向后兼容。
var CacheModule = fx.Module("cache",
	fx.Provide(
		// 独立缓存服务
		newSettingCacheService,
		newSettingCategoryCacheService,
		newUserSettingQueryCacheService,
		newUserSettingCacheService,
		newUserWithRolesCacheService,
		newPermissionCacheService,
		newSchemaCacheService,

		// 聚合模块，用于向后兼容
		newCacheServicesModule,
	),
)

// getKeyPrefix 从配置获取 Redis Key 前缀。
func getKeyPrefix(cfg *config.Config) string {
	return cfg.Data.RedisKeyPrefix
}

func newSettingCacheService(client *redis.Client, cfg *config.Config) cache.SettingCacheService {
	return infracache.NewSettingCacheService(client, getKeyPrefix(cfg))
}

func newSettingCategoryCacheService(client *redis.Client, cfg *config.Config) cache.SettingCategoryCacheService {
	return infracache.NewSettingCategoryCacheService(client, getKeyPrefix(cfg))
}

func newUserSettingQueryCacheService(client *redis.Client, cfg *config.Config) cache.UserSettingQueryCacheService {
	return infracache.NewUserSettingQueryCacheService(client, getKeyPrefix(cfg))
}

func newUserSettingCacheService(client *redis.Client, cfg *config.Config) cache.UserSettingCacheService {
	return infracache.NewUserSettingCacheService(client, getKeyPrefix(cfg))
}

func newUserWithRolesCacheService(client *redis.Client, cfg *config.Config) cache.UserWithRolesCacheService {
	return infracache.NewUserWithRolesCacheService(client, getKeyPrefix(cfg))
}

func newPermissionCacheService(client *redis.Client, cfg *config.Config) cache.PermissionCacheService {
	return infracache.NewPermissionCacheService(client, getKeyPrefix(cfg))
}

func newSchemaCacheService(client *redis.Client, cfg *config.Config) setting.SchemaCacheService {
	return infracache.NewSchemaCacheService(client, getKeyPrefix(cfg))
}

// newCacheServicesModule 创建聚合的缓存服务模块。
// 用于向后兼容依赖该结构体的现有代码。
func newCacheServicesModule(
	setting cache.SettingCacheService,
	settingCategory cache.SettingCategoryCacheService,
	userSettingQuery cache.UserSettingQueryCacheService,
	userSetting cache.UserSettingCacheService,
	userWithRoles cache.UserWithRolesCacheService,
	permission cache.PermissionCacheService,
	schema setting.SchemaCacheService,
) *CacheServicesModule {
	return &CacheServicesModule{
		Setting:          setting,
		SettingCategory:  settingCategory,
		UserSettingQuery: userSettingQuery,
		UserSetting:      userSetting,
		UserWithRoles:    userWithRoles,
		Permission:       permission,
		Schema:           schema,
	}
}
