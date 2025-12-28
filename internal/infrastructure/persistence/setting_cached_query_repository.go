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
// 采用 Cache-Aside 模式，版本化异步回写缓存防止竞争。
//
// 缓存策略：
//   - [FindByKey]: 先查缓存，未命中查库后版本化回写
//   - [FindByKeys]: 批量从缓存获取，未命中的查库后批量回写
//   - [FindByCategoryID]: 预热后从全量缓存过滤，否则查库
//   - [FindAll]: 预热后从缓存获取全量，否则查库
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

	// 3. 异步版本化回写缓存（防止回写竞争）
	go func(s *setting.Setting) {
		cacheCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 2*time.Second)
		defer cancel()

		version := s.UpdatedAt.UnixNano()
		if written, err := r.cache.SetWithVersion(cacheCtx, s, version); err != nil {
			slog.Warn("cache set with version failed", "key", s.Key, "err", err)
		} else if !written {
			slog.Debug("cache set skipped (newer version exists)", "key", s.Key)
		}
	}(result)

	return result, nil
}

// FindByKeys 批量查找配置定义（带缓存）。
func (r *cachedSettingQueryRepository) FindByKeys(ctx context.Context, keys []string) ([]*setting.Setting, error) {
	if len(keys) == 0 {
		return []*setting.Setting{}, nil
	}

	// 1. 批量从缓存获取
	cachedMap, err := r.cache.GetByKeys(ctx, keys)
	if err != nil {
		slog.Warn("cache batch get failed, fallback to db", "count", len(keys), "err", err)
		return r.delegate.FindByKeys(ctx, keys)
	}

	// 2. 找出未命中的 key
	var missedKeys []string
	for _, k := range keys {
		if _, ok := cachedMap[k]; !ok {
			missedKeys = append(missedKeys, k)
		}
	}

	// 3. 查询未命中的
	if len(missedKeys) > 0 { //nolint:nestif // 缓存回写逻辑需要嵌套判断
		dbResults, err := r.delegate.FindByKeys(ctx, missedKeys)
		if err != nil {
			return nil, err
		}

		// 异步批量回写缓存
		if len(dbResults) > 0 {
			go func(settings []*setting.Setting) {
				cacheCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
				defer cancel()
				if err := r.cache.SetAll(cacheCtx, settings); err != nil {
					slog.Warn("cache batch set failed", "count", len(settings), "err", err)
				}
			}(dbResults)
		}

		// 合并结果
		for _, s := range dbResults {
			cachedMap[s.Key] = s
		}
	}

	// 4. 按原始顺序返回
	result := make([]*setting.Setting, 0, len(keys))
	for _, k := range keys {
		if s, ok := cachedMap[k]; ok {
			result = append(result, s)
		}
	}

	return result, nil
}

// FindByCategoryID 根据分类 ID 查找配置定义列表。
//
// 如果缓存已预热，从全量缓存中过滤；否则直接查库。
// 若缓存过滤结果为空，回退数据库查询（防止缓存部分过期）。
func (r *cachedSettingQueryRepository) FindByCategoryID(ctx context.Context, categoryID uint) ([]*setting.Setting, error) {
	// 如果已预热，从缓存中过滤
	if r.cache.IsWarmedUp(ctx) { //nolint:nestif // 缓存查询逻辑需要嵌套判断
		cachedMap, err := r.cache.GetAll(ctx)
		if err == nil && len(cachedMap) > 0 {
			result := make([]*setting.Setting, 0)
			for _, s := range cachedMap {
				if s.CategoryID == categoryID {
					result = append(result, s)
				}
			}
			// 仅当过滤结果非空时使用缓存，否则回退数据库
			// 这可以防止缓存部分过期导致的空结果
			if len(result) > 0 {
				return result, nil
			}
			slog.Debug("cache filter returned empty, fallback to db", "categoryID", categoryID)
		}
		// 缓存获取失败，回退到数据库
		if err != nil {
			slog.Warn("cache get all failed, fallback to db", "categoryID", categoryID, "err", err)
		}
	}

	return r.delegate.FindByCategoryID(ctx, categoryID)
}

// FindByScope 根据作用域查找配置定义列表。
//
// 如果缓存已预热，从全量缓存中过滤；否则直接查库。
func (r *cachedSettingQueryRepository) FindByScope(ctx context.Context, scope string) ([]*setting.Setting, error) {
	// 如果已预热，从缓存中过滤
	if r.cache.IsWarmedUp(ctx) { //nolint:nestif // 缓存查询逻辑需要嵌套判断
		cachedMap, err := r.cache.GetAll(ctx)
		if err == nil && len(cachedMap) > 0 {
			result := make([]*setting.Setting, 0)
			for _, s := range cachedMap {
				if s.Scope == scope {
					result = append(result, s)
				}
			}
			return result, nil
		}
		if err != nil {
			slog.Warn("cache get all failed, fallback to db", "scope", scope, "err", err)
		}
	}

	return r.delegate.FindByScope(ctx, scope)
}

// FindAll 查找所有配置定义。
//
// 如果缓存已预热，从缓存获取全量；否则直接查库。
func (r *cachedSettingQueryRepository) FindAll(ctx context.Context) ([]*setting.Setting, error) {
	// 如果已预热，从缓存获取
	if r.cache.IsWarmedUp(ctx) {
		cachedMap, err := r.cache.GetAll(ctx)
		if err == nil && len(cachedMap) > 0 {
			result := make([]*setting.Setting, 0, len(cachedMap))
			for _, s := range cachedMap {
				result = append(result, s)
			}
			return result, nil
		}
		if err != nil {
			slog.Warn("cache get all failed, fallback to db", "err", err)
		}
	}

	return r.delegate.FindAll(ctx)
}

// ExistsByKey 检查 Key 是否已存在。
// 存在性检查直接委托，不缓存。
func (r *cachedSettingQueryRepository) ExistsByKey(ctx context.Context, key string) (bool, error) {
	// 如果已预热，可以从缓存判断
	if r.cache.IsWarmedUp(ctx) {
		cached, err := r.cache.Get(ctx, key)
		if err == nil && cached != nil {
			return true, nil
		}
		// 缓存未命中，还需要查库（可能是新增的）
	}

	return r.delegate.ExistsByKey(ctx, key)
}

var _ setting.QueryRepository = (*cachedSettingQueryRepository)(nil)
