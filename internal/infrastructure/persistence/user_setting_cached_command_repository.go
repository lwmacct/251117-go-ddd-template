package persistence

import (
	"context"
	"log/slog"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/cache"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

// cachedUserSettingCommandRepository 带缓存失效的 UserSetting 命令仓储装饰器。
//
// 装饰 [setting.UserSettingCommandRepository]，在写操作后自动失效相关缓存。
// 不负责缓存填充（缓存填充在 Application 层完成，因为需要合并 Setting + UserSetting）。
type cachedUserSettingCommandRepository struct {
	delegate         setting.UserSettingCommandRepository
	userSettingCache cache.UserSettingCacheService
}

// NewCachedUserSettingCommandRepository 创建带缓存失效的 UserSetting 命令仓储。
func NewCachedUserSettingCommandRepository(
	delegate setting.UserSettingCommandRepository,
	userSettingCache cache.UserSettingCacheService,
) setting.UserSettingCommandRepository {
	return &cachedUserSettingCommandRepository{
		delegate:         delegate,
		userSettingCache: userSettingCache,
	}
}

// Upsert 插入或更新用户配置，并失效缓存。
func (r *cachedUserSettingCommandRepository) Upsert(ctx context.Context, us *setting.UserSetting) error {
	if err := r.delegate.Upsert(ctx, us); err != nil {
		return err
	}

	// 失效该用户的该 key 缓存
	if err := r.userSettingCache.Delete(ctx, us.UserID, us.SettingKey); err != nil {
		slog.Warn("failed to invalidate user setting cache after upsert",
			"userID", us.UserID, "key", us.SettingKey, "err", err)
	}

	return nil
}

// Delete 删除用户配置，并失效缓存。
func (r *cachedUserSettingCommandRepository) Delete(ctx context.Context, userID uint, key string) error {
	if err := r.delegate.Delete(ctx, userID, key); err != nil {
		return err
	}

	// 失效该用户的该 key 缓存
	if err := r.userSettingCache.Delete(ctx, userID, key); err != nil {
		slog.Warn("failed to invalidate user setting cache after delete",
			"userID", userID, "key", key, "err", err)
	}

	return nil
}

// DeleteByUser 删除用户的所有配置，并失效缓存。
func (r *cachedUserSettingCommandRepository) DeleteByUser(ctx context.Context, userID uint) error {
	if err := r.delegate.DeleteByUser(ctx, userID); err != nil {
		return err
	}

	// 失效该用户的所有缓存
	if err := r.userSettingCache.DeleteByUser(ctx, userID); err != nil {
		slog.Warn("failed to invalidate all user setting cache after delete by user",
			"userID", userID, "err", err)
	}

	return nil
}

// BatchUpsert 批量插入或更新用户配置，并失效缓存。
func (r *cachedUserSettingCommandRepository) BatchUpsert(ctx context.Context, settings []*setting.UserSetting) error {
	if err := r.delegate.BatchUpsert(ctx, settings); err != nil {
		return err
	}

	// 收集需要失效的 key
	if len(settings) > 0 {
		userID := settings[0].UserID
		keys := make([]string, 0, len(settings))
		for _, s := range settings {
			keys = append(keys, s.SettingKey)
		}

		// 批量失效缓存
		if err := r.userSettingCache.DeleteByKeys(ctx, userID, keys); err != nil {
			slog.Warn("failed to invalidate user setting cache after batch upsert",
				"userID", userID, "count", len(keys), "err", err)
		}
	}

	return nil
}

var _ setting.UserSettingCommandRepository = (*cachedUserSettingCommandRepository)(nil)
