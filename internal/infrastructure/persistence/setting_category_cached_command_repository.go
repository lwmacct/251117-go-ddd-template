package persistence

import (
	"context"
	"log/slog"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/cache"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

// cachedSettingCategoryCommandRepository 带缓存失效的 SettingCategory 命令仓储装饰器。
//
// 装饰 [setting.SettingCategoryCommandRepository]，写操作后自动失效相关缓存。
// 缓存失效失败不阻塞业务，仅记录警告日志。
type cachedSettingCategoryCommandRepository struct {
	delegate setting.SettingCategoryCommandRepository
	cache    cache.SettingCategoryCacheService
}

// NewCachedSettingCategoryCommandRepository 创建带缓存失效的 SettingCategory 命令仓储。
func NewCachedSettingCategoryCommandRepository(
	delegate setting.SettingCategoryCommandRepository,
	cacheService cache.SettingCategoryCacheService,
) setting.SettingCategoryCommandRepository {
	return &cachedSettingCategoryCommandRepository{
		delegate: delegate,
		cache:    cacheService,
	}
}

// Create 创建配置分类。
//
// 创建成功后同步写入单条缓存并失效全量缓存，确保缓存一致性。
func (r *cachedSettingCategoryCommandRepository) Create(ctx context.Context, category *setting.SettingCategory) error {
	if err := r.delegate.Create(ctx, category); err != nil {
		return err
	}

	// 同步写入单条缓存
	if err := r.cache.Set(ctx, category); err != nil {
		slog.Warn("cache set failed after create", "id", category.ID, "key", category.Key, "err", err)
	}

	// 失效全量缓存（因为 Set 只写入 ID 和 Key 索引，不更新 all 列表）
	if err := r.cache.DeleteAll(ctx); err != nil {
		slog.Warn("cache delete all failed after create", "id", category.ID, "err", err)
	}

	return nil
}

// Update 更新配置分类。
//
// 更新成功后同步更新缓存并失效全量缓存。
func (r *cachedSettingCategoryCommandRepository) Update(ctx context.Context, category *setting.SettingCategory) error {
	if err := r.delegate.Update(ctx, category); err != nil {
		return err
	}

	// 同步更新单条缓存
	if err := r.cache.Set(ctx, category); err != nil {
		slog.Warn("cache set failed after update", "id", category.ID, "key", category.Key, "err", err)
	}

	// 失效全量缓存
	if err := r.cache.DeleteAll(ctx); err != nil {
		slog.Warn("cache delete all failed after update", "id", category.ID, "err", err)
	}

	return nil
}

// Delete 删除配置分类。
//
// 删除成功后失效所有缓存，确保一致性。
// 由于只有 ID 参数，无法精确删除 Key 索引，故采用 DeleteAll。
func (r *cachedSettingCategoryCommandRepository) Delete(ctx context.Context, id uint) error {
	if err := r.delegate.Delete(ctx, id); err != nil {
		return err
	}

	// 失效所有缓存（包括 ID 索引、Key 索引、全量缓存）
	if err := r.cache.DeleteAll(ctx); err != nil {
		slog.Warn("cache delete all failed after delete", "id", id, "err", err)
	}

	return nil
}

var _ setting.SettingCategoryCommandRepository = (*cachedSettingCategoryCommandRepository)(nil)
