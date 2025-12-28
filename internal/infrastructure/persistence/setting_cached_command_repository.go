package persistence

import (
	"context"
	"log/slog"
	"time"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/cache"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

// cachedSettingCommandRepository 带缓存失效的 Setting 命令仓储装饰器。
//
// 装饰 [setting.CommandRepository]，写操作后自动失效相关缓存。
// 缓存失效失败不阻塞业务，仅记录警告日志。
//
// 跨层失效策略：
//   - 系统设置更新时，同时失效系统缓存和所有用户的对应缓存
//   - 用户设置缓存失效异步执行，避免阻塞管理员操作
type cachedSettingCommandRepository struct {
	delegate         setting.CommandRepository
	systemCache      cache.SettingCacheService
	userSettingCache cache.UserSettingCacheService // 可选，用于跨层失效
}

// NewCachedSettingCommandRepository 创建带缓存失效的 Setting 命令仓储。
func NewCachedSettingCommandRepository(
	delegate setting.CommandRepository,
	systemCacheService cache.SettingCacheService,
) setting.CommandRepository {
	return &cachedSettingCommandRepository{
		delegate:    delegate,
		systemCache: systemCacheService,
	}
}

// NewCachedSettingCommandRepositoryWithUserCache 创建带用户缓存失效的 Setting 命令仓储。
//
// 当系统设置默认值变更时，会异步失效所有用户的对应缓存。
func NewCachedSettingCommandRepositoryWithUserCache(
	delegate setting.CommandRepository,
	systemCacheService cache.SettingCacheService,
	userSettingCacheService cache.UserSettingCacheService,
) setting.CommandRepository {
	return &cachedSettingCommandRepository{
		delegate:         delegate,
		systemCache:      systemCacheService,
		userSettingCache: userSettingCacheService,
	}
}

// Create 创建配置定义。
// 创建新记录无需失效缓存（新 key 不在缓存中）。
func (r *cachedSettingCommandRepository) Create(ctx context.Context, s *setting.Setting) error {
	return r.delegate.Create(ctx, s)
}

// Update 更新配置定义。
// 更新后失效系统缓存和所有用户的对应缓存。
func (r *cachedSettingCommandRepository) Update(ctx context.Context, s *setting.Setting) error {
	if err := r.delegate.Update(ctx, s); err != nil {
		return err
	}

	// 1. 同步失效系统缓存
	if err := r.systemCache.Delete(ctx, s.Key); err != nil {
		slog.Warn("system cache delete failed after update", "key", s.Key, "err", err)
	}

	// 2. 异步失效所有用户的该 key 缓存（仅对 user scope 设置）
	// 使用独立 context 确保失效操作不受调用方 context 取消影响
	if s.IsUserScope() && r.userSettingCache != nil {
		go func(key string) { //nolint:contextcheck // 故意使用独立 context
			asyncCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if err := r.userSettingCache.DeleteBySettingKey(asyncCtx, key); err != nil {
				slog.Warn("user setting cache invalidation failed after update",
					"key", key, "err", err)
			} else {
				slog.Debug("user setting cache invalidated after update", "key", key)
			}
		}(s.Key)
	}

	return nil
}

// Delete 删除配置定义。
// 删除后失效系统缓存和所有用户的对应缓存。
func (r *cachedSettingCommandRepository) Delete(ctx context.Context, key string) error {
	// 先获取设置信息（判断是否需要失效用户缓存）
	// 注意：这里简化处理，直接尝试失效用户缓存
	if err := r.delegate.Delete(ctx, key); err != nil {
		return err
	}

	// 1. 同步失效系统缓存
	if err := r.systemCache.Delete(ctx, key); err != nil {
		slog.Warn("system cache delete failed after delete", "key", key, "err", err)
	}

	// 2. 异步失效所有用户的该 key 缓存
	// 使用独立 context 确保失效操作不受调用方 context 取消影响
	if r.userSettingCache != nil {
		go func(settingKey string) { //nolint:contextcheck // 故意使用独立 context
			asyncCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if err := r.userSettingCache.DeleteBySettingKey(asyncCtx, settingKey); err != nil {
				slog.Warn("user setting cache invalidation failed after delete",
					"key", settingKey, "err", err)
			}
		}(key)
	}

	return nil
}

// BatchUpsert 批量插入或更新配置定义。
// 批量操作后失效所有相关 key 的系统缓存和用户缓存。
func (r *cachedSettingCommandRepository) BatchUpsert(ctx context.Context, settings []*setting.Setting) error {
	if err := r.delegate.BatchUpsert(ctx, settings); err != nil {
		return err
	}

	// 收集需要失效的 key
	keys := make([]string, 0, len(settings))
	userScopeKeys := make([]string, 0)
	for _, s := range settings {
		keys = append(keys, s.Key)
		if s.IsUserScope() {
			userScopeKeys = append(userScopeKeys, s.Key)
		}
	}

	// 1. 同步批量失效系统缓存
	if err := r.systemCache.DeleteByKeys(ctx, keys); err != nil {
		slog.Warn("system cache batch delete failed", "count", len(keys), "err", err)
	}

	// 2. 异步批量失效用户缓存
	// 使用独立 context 确保失效操作不受调用方 context 取消影响
	if len(userScopeKeys) > 0 && r.userSettingCache != nil {
		go func(settingKeys []string) { //nolint:contextcheck // 故意使用独立 context
			asyncCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			if err := r.userSettingCache.DeleteBySettingKeys(asyncCtx, settingKeys); err != nil {
				slog.Warn("user setting cache batch invalidation failed",
					"count", len(settingKeys), "err", err)
			}
		}(userScopeKeys)
	}

	return nil
}

var _ setting.CommandRepository = (*cachedSettingCommandRepository)(nil)
