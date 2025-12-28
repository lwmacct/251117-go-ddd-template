package bootstrap

import (
	goredis "github.com/redis/go-redis/v9"

	infraredis "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/redis"
)

// newCacheServicesModule 初始化缓存服务模块
// 依赖：Redis 客户端和 key 前缀
func newCacheServicesModule(redisClient *goredis.Client, keyPrefix string) *CacheServicesModule {
	return &CacheServicesModule{
		Setting:          infraredis.NewSettingCacheService(redisClient, keyPrefix),
		UserSettingQuery: infraredis.NewUserSettingQueryCacheService(redisClient, keyPrefix),
		UserSetting:      infraredis.NewUserSettingCacheService(redisClient, keyPrefix),
		Permission:       infraredis.NewPermissionCacheService(redisClient, keyPrefix),
		Schema:           infraredis.NewSchemaCacheService(redisClient, keyPrefix),
	}
}
