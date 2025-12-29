package container

import (
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"

	"github.com/lwmacct/251117-go-ddd-template/internal/application/setting"
	"github.com/lwmacct/251117-go-ddd-template/internal/config"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/cache"

	infracache "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/cache"
)

// CacheServicesResult 使用 fx.Out 批量返回所有缓存服务。
type CacheServicesResult struct {
	fx.Out

	SettingCategory  cache.SettingCategoryCacheService
	UserSettingQuery cache.UserSettingQueryCacheService
	UserSetting      cache.UserSettingCacheService
	UserWithRoles    cache.UserWithRolesCacheService
	Permission       cache.PermissionCacheService
	Schema           setting.SchemaCacheService
}

// CacheModule 提供所有缓存服务。
var CacheModule = fx.Module("cache",
	fx.Provide(NewAllCacheServices),
)

// NewAllCacheServices 创建所有缓存服务。
func NewAllCacheServices(client *redis.Client, cfg *config.Config) CacheServicesResult {
	prefix := cfg.Data.RedisKeyPrefix
	return CacheServicesResult{
		SettingCategory:  infracache.NewSettingCategoryCacheService(client, prefix),
		UserSettingQuery: infracache.NewUserSettingQueryCacheService(client, prefix),
		UserSetting:      infracache.NewUserSettingCacheService(client, prefix),
		UserWithRoles:    infracache.NewUserWithRolesCacheService(client, prefix),
		Permission:       infracache.NewPermissionCacheService(client, prefix),
		Schema:           infracache.NewSchemaCacheService(client, prefix),
	}
}
