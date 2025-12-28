package persistence

import (
	"context"
	"log/slog"
	"time"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/cache"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

// cachedSettingQueryRepository 带缓存的 Setting 查询仓储装饰器。
//
// 装饰 [setting.QueryRepository]，在查询前检查缓存，未命中再查数据库。
// 采用 Cache-Aside 模式，异步回写缓存不阻塞请求。
type cachedSettingQueryRepository struct {
	delegate setting.QueryRepository   // 被装饰的原始仓储
	cache    cache.SettingCacheService // 缓存服务
}

// NewCachedSettingQueryRepository 创建带缓存的 Setting 查询仓储。
func NewCachedSettingQueryRepository(
	delegate setting.QueryRepository,
	cacheService cache.SettingCacheService,
) setting.QueryRepository {
	return &cachedSettingQueryRepository{
		delegate: delegate,
		cache:    cacheService,
	}
}

// FindByKey 根据 Key 查找配置定义（带缓存）。
func (r *cachedSettingQueryRepository) FindByKey(ctx context.Context, key string) (*setting.Setting, error) {
	// 1. 查缓存
	cached, err := r.cache.Get(ctx, key)
	if err != nil {
		slog.Warn("cache get failed, fallback to db", "key", key, "err", err)
	}
	if cached != nil {
		return cached, nil
	}

	// 2. 查数据库
	result, err := r.delegate.FindByKey(ctx, key)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil //nolint:nilnil // not found
	}

	// 3. 异步回写缓存
	go func(s *setting.Setting) {
		cacheCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 2*time.Second)
		defer cancel()
		if err := r.cache.Set(cacheCtx, s); err != nil {
			slog.Warn("cache set failed", "key", s.Key, "err", err)
		}
	}(result)

	return result, nil
}

// FindByKeys 批量查找配置定义。
// 批量查询场景缓存收益有限，直接查库。
func (r *cachedSettingQueryRepository) FindByKeys(ctx context.Context, keys []string) ([]*setting.Setting, error) {
	return r.delegate.FindByKeys(ctx, keys)
}

// FindByCategoryID 根据分类 ID 查找配置定义列表。
// 列表查询不缓存，直接委托。
func (r *cachedSettingQueryRepository) FindByCategoryID(ctx context.Context, categoryID uint) ([]*setting.Setting, error) {
	return r.delegate.FindByCategoryID(ctx, categoryID)
}

// FindByScope 根据作用域查找配置定义列表。
// 列表查询不缓存，直接委托。
func (r *cachedSettingQueryRepository) FindByScope(ctx context.Context, scope string) ([]*setting.Setting, error) {
	return r.delegate.FindByScope(ctx, scope)
}

// FindAll 查找所有配置定义。
// 列表查询不缓存，直接委托。
func (r *cachedSettingQueryRepository) FindAll(ctx context.Context) ([]*setting.Setting, error) {
	return r.delegate.FindAll(ctx)
}

// ExistsByKey 检查 Key 是否已存在。
// 存在性检查不缓存，直接委托。
func (r *cachedSettingQueryRepository) ExistsByKey(ctx context.Context, key string) (bool, error) {
	return r.delegate.ExistsByKey(ctx, key)
}

var _ setting.QueryRepository = (*cachedSettingQueryRepository)(nil)
