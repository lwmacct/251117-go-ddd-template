package bootstrap

import (
	goredis "github.com/redis/go-redis/v9"

	infracache "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/cache"
)

// newCacheServicesModule 初始化缓存服务模块
// 依赖：Redis 客户端和 key 前缀
func newCacheServicesModule(redisClient *goredis.Client, keyPrefix string) *CacheServicesModule {
	return &CacheServicesModule{
		Setting:          infracache.NewSettingCacheService(redisClient, keyPrefix),
		SettingCategory:  infracache.NewSettingCategoryCacheService(redisClient, keyPrefix),
		UserSettingQuery: infracache.NewUserSettingQueryCacheService(redisClient, keyPrefix),
		UserSetting:      infracache.NewUserSettingCacheService(redisClient, keyPrefix),
		UserWithRoles:    infracache.NewUserWithRolesCacheService(redisClient, keyPrefix),
		Permission:       infracache.NewPermissionCacheService(redisClient, keyPrefix),
		Schema:           infracache.NewSchemaCacheService(redisClient, keyPrefix),
	}
}
