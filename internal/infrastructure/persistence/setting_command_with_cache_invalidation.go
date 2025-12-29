package persistence

import (
	"context"
	"log/slog"
	"time"

	appsetting "github.com/lwmacct/251117-go-ddd-template/internal/application/setting"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/cache"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

// settingCommandWithCacheInvalidation 写操作后失效下游缓存的装饰器。
//
// 简化设计：
//   - 不再缓存 Setting 实体本身（由 Application 层 Schema 缓存覆盖）
//   - 只负责写操作后失效下游缓存（Schema + UserSetting）
//
// 失效策略：
//   - Create/Update/Delete/BatchUpsert 后异步失效 Schema 缓存
//   - Update/Delete user scope 设置后异步失效 UserSetting 缓存
type settingCommandWithCacheInvalidation struct {
	delegate         setting.CommandRepository
	userSettingCache cache.UserSettingCacheService
	schemaCache      appsetting.SchemaCacheService
}

// NewSettingCommandWithCacheInvalidation 创建带缓存失效的 Setting 命令仓储。
func NewSettingCommandWithCacheInvalidation(
	delegate setting.CommandRepository,
	userSettingCache cache.UserSettingCacheService,
	schemaCache appsetting.SchemaCacheService,
) setting.CommandRepository {
	return &settingCommandWithCacheInvalidation{
		delegate:         delegate,
		userSettingCache: userSettingCache,
		schemaCache:      schemaCache,
	}
}

// Create 创建配置定义。
func (r *settingCommandWithCacheInvalidation) Create(ctx context.Context, s *setting.Setting) error {
	if err := r.delegate.Create(ctx, s); err != nil {
		return err
	}
	r.invalidateSchemaCacheAsync() //nolint:contextcheck // 故意使用独立 context 进行异步失效
	return nil
}

// Update 更新配置定义。
func (r *settingCommandWithCacheInvalidation) Update(ctx context.Context, s *setting.Setting) error {
	if err := r.delegate.Update(ctx, s); err != nil {
		return err
	}
	r.invalidateSchemaCacheAsync() //nolint:contextcheck // 故意使用独立 context 进行异步失效
	if s.IsUserScope() {
		r.invalidateUserSettingCacheAsync(s.Key) //nolint:contextcheck // 故意使用独立 context 进行异步失效
	}
	return nil
}

// Delete 删除配置定义。
func (r *settingCommandWithCacheInvalidation) Delete(ctx context.Context, key string) error {
	if err := r.delegate.Delete(ctx, key); err != nil {
		return err
	}
	r.invalidateSchemaCacheAsync()         //nolint:contextcheck // 故意使用独立 context 进行异步失效
	r.invalidateUserSettingCacheAsync(key) //nolint:contextcheck // 故意使用独立 context 进行异步失效
	return nil
}

// BatchUpsert 批量插入或更新配置定义。
func (r *settingCommandWithCacheInvalidation) BatchUpsert(ctx context.Context, settings []*setting.Setting) error {
	if err := r.delegate.BatchUpsert(ctx, settings); err != nil {
		return err
	}
	r.invalidateSchemaCacheAsync() //nolint:contextcheck // 故意使用独立 context 进行异步失效
	for _, s := range settings {
		if s.IsUserScope() {
			r.invalidateUserSettingCacheAsync(s.Key) //nolint:contextcheck // 故意使用独立 context 进行异步失效
		}
	}
	return nil
}

// invalidateSchemaCacheAsync 异步失效 Schema 缓存。
func (r *settingCommandWithCacheInvalidation) invalidateSchemaCacheAsync() {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := r.schemaCache.DeleteAll(ctx); err != nil {
			slog.Warn("failed to invalidate schema cache", "err", err)
		}
	}()
}

// invalidateUserSettingCacheAsync 异步失效 UserSetting 缓存。
func (r *settingCommandWithCacheInvalidation) invalidateUserSettingCacheAsync(key string) {
	if r.userSettingCache == nil {
		return
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := r.userSettingCache.DeleteBySettingKey(ctx, key); err != nil {
			slog.Warn("failed to invalidate user setting cache", "key", key, "err", err)
		}
	}()
}

var _ setting.CommandRepository = (*settingCommandWithCacheInvalidation)(nil)
