package persistence

import (
	"context"
	"log/slog"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/cache"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

// cachedSettingCommandRepository 带缓存失效的 Setting 命令仓储装饰器。
//
// 装饰 [setting.CommandRepository]，写操作后自动失效相关缓存。
// 缓存失效失败不阻塞业务，仅记录警告日志。
type cachedSettingCommandRepository struct {
	delegate setting.CommandRepository
	cache    cache.SettingCacheService
}

// NewCachedSettingCommandRepository 创建带缓存失效的 Setting 命令仓储。
func NewCachedSettingCommandRepository(
	delegate setting.CommandRepository,
	cacheService cache.SettingCacheService,
) setting.CommandRepository {
	return &cachedSettingCommandRepository{
		delegate: delegate,
		cache:    cacheService,
	}
}

// Create 创建配置定义。
// 创建新记录无需失效缓存（新 key 不在缓存中）。
func (r *cachedSettingCommandRepository) Create(ctx context.Context, s *setting.Setting) error {
	return r.delegate.Create(ctx, s)
}

// Update 更新配置定义。
// 更新后失效对应 key 的缓存。
func (r *cachedSettingCommandRepository) Update(ctx context.Context, s *setting.Setting) error {
	if err := r.delegate.Update(ctx, s); err != nil {
		return err
	}

	// 失效缓存（失败不阻塞业务）
	if err := r.cache.Delete(ctx, s.Key); err != nil {
		slog.Warn("cache delete failed after update", "key", s.Key, "err", err)
	}

	return nil
}

// Delete 删除配置定义。
// 删除后失效对应 key 的缓存。
func (r *cachedSettingCommandRepository) Delete(ctx context.Context, key string) error {
	if err := r.delegate.Delete(ctx, key); err != nil {
		return err
	}

	// 失效缓存（失败不阻塞业务）
	if err := r.cache.Delete(ctx, key); err != nil {
		slog.Warn("cache delete failed after delete", "key", key, "err", err)
	}

	return nil
}

// BatchUpsert 批量插入或更新配置定义。
// 批量操作后失效所有相关 key 的缓存。
func (r *cachedSettingCommandRepository) BatchUpsert(ctx context.Context, settings []*setting.Setting) error {
	if err := r.delegate.BatchUpsert(ctx, settings); err != nil {
		return err
	}

	// 批量失效缓存
	keys := make([]string, 0, len(settings))
	for _, s := range settings {
		keys = append(keys, s.Key)
	}

	if err := r.cache.DeleteByKeys(ctx, keys); err != nil {
		slog.Warn("cache batch delete failed", "count", len(keys), "err", err)
	}

	return nil
}

var _ setting.CommandRepository = (*cachedSettingCommandRepository)(nil)
