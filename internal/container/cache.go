package container

import (
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"

	"github.com/lwmacct/251117-go-ddd-template/internal/config"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/cache"

	"github.com/lwmacct/251117-go-ddd-template/internal/application/auth"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/setting"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/user"
)

// CacheServicesResult 使用 fx.Out 批量返回所有缓存服务。
type CacheServicesResult struct {
	fx.Out

	UserSettingQuery setting.UserSettingQueryCacheService
	UserSetting      setting.UserSettingCacheService
	UserWithRoles    user.UserWithRolesCacheService
	Permission       auth.PermissionCacheService
	Settings         setting.SettingsCacheService
}

// CacheModule 提供所有缓存服务。
var CacheModule = fx.Module("cache",
	fx.Provide(NewAllCacheServices),
)

// NewAllCacheServices 创建所有缓存服务。
func NewAllCacheServices(client *redis.Client, cfg *config.Config) CacheServicesResult {
	prefix := cfg.Data.RedisKeyPrefix
	return CacheServicesResult{
		UserSettingQuery: cache.NewUserSettingQueryCacheService(client, prefix),
		UserSetting:      cache.NewUserSettingCacheService(client, prefix),
		UserWithRoles:    cache.NewUserWithRolesCacheService(client, prefix),
		Permission:       cache.NewPermissionCacheService(client, prefix),
		Settings:         cache.NewSettingsCacheService(client, prefix),
	}
}
