package container

import (
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"

	"github.com/lwmacct/251117-go-ddd-template/internal/application/setting"
	"github.com/lwmacct/251117-go-ddd-template/internal/config"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/cache"

	infracache "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/cache"
)

// CacheServicesModule aggregates all cache services.
// Defined in bootstrap package for compatibility.
type CacheServicesModule struct {
	Setting          cache.SettingCacheService
	SettingCategory  cache.SettingCategoryCacheService
	UserSettingQuery cache.UserSettingQueryCacheService
	UserSetting      cache.UserSettingCacheService
	UserWithRoles    cache.UserWithRolesCacheService
	Permission       cache.PermissionCacheService
	Schema           setting.SchemaCacheService
}

// CacheModule provides all cache services.
//
// Each service is provided individually for fine-grained dependency injection,
// and also aggregated into CacheServicesModule for backward compatibility.
var CacheModule = fx.Module("cache",
	fx.Provide(
		// Individual cache services
		newSettingCacheService,
		newSettingCategoryCacheService,
		newUserSettingQueryCacheService,
		newUserSettingCacheService,
		newUserWithRolesCacheService,
		newPermissionCacheService,
		newSchemaCacheService,

		// Aggregated module for backward compatibility
		newCacheServicesModule,
	),
)

// Helper to get key prefix from config
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

// newCacheServicesModule creates the aggregated cache services module.
// This is for backward compatibility with existing code that depends on the struct.
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
