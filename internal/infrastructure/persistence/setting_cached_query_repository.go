package persistence

import (
	"context"
	"log/slog"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/cache"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

// cachedSettingQueryRepository 带缓存的 Setting 查询仓储装饰器。
//
// 装饰 [setting.QueryRepository]，采用真正的 Cache-Aside 模式：
//   - 所有查询先查缓存，未命中再查数据库
//   - 查库后同步回写缓存
//   - 不依赖预热标记，支持惰性加载
//
// 缓存策略：
//   - [FindByKey]: 先查缓存，未命中查库后同步回写
//   - [FindByKeys]: 批量从缓存获取，未命中的查库后批量回写
//   - [FindByCategoryID]: 从全量缓存过滤，未命中查库后回写
//   - [FindAll]: 从缓存获取全量，未命中查库后回写
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

	// 3. 同步回写缓存
	if err := r.cache.Set(ctx, result); err != nil {
		slog.Warn("cache set failed", "key", result.Key, "err", err)
	}

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

		// 同步批量回写缓存
		if len(dbResults) > 0 {
			if err := r.cache.SetAll(ctx, dbResults); err != nil {
				slog.Warn("cache batch set failed", "count", len(dbResults), "err", err)
			}
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

// FindByCategoryID 根据分类 ID 查找配置定义列表（带缓存）。
//
// 采用 Cache-Aside 模式：
//  1. 尝试从全量缓存中过滤
//  2. 缓存未命中则查数据库
//  3. 查库后同步回写全量缓存
func (r *cachedSettingQueryRepository) FindByCategoryID(ctx context.Context, categoryID uint) ([]*setting.Setting, error) {
	// 1. 尝试从缓存获取
	if result, ok := r.tryFilterFromCache(ctx, func(s *setting.Setting) bool {
		return s.CategoryID == categoryID
	}); ok {
		return result, nil
	}

	// 2. 缓存未命中，查数据库
	result, err := r.delegate.FindByCategoryID(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	// 3. 同步回写缓存
	r.trySetAllCache(ctx, result, "categoryID", categoryID)

	return result, nil
}

// FindByScope 根据作用域查找配置定义列表（带缓存）。
//
// 采用 Cache-Aside 模式：
//  1. 尝试从全量缓存中过滤
//  2. 缓存未命中则查数据库
//  3. 查库后同步回写全量缓存
func (r *cachedSettingQueryRepository) FindByScope(ctx context.Context, scope string) ([]*setting.Setting, error) {
	// 1. 尝试从缓存获取
	if result, ok := r.tryFilterFromCache(ctx, func(s *setting.Setting) bool {
		return s.Scope == scope
	}); ok {
		return result, nil
	}

	// 2. 缓存未命中，查数据库
	result, err := r.delegate.FindByScope(ctx, scope)
	if err != nil {
		return nil, err
	}

	// 3. 同步回写缓存
	r.trySetAllCache(ctx, result, "scope", scope)

	return result, nil
}

// FindAll 查找所有配置定义（带缓存）。
//
// 采用 Cache-Aside 模式：
//  1. 尝试从缓存获取全量
//  2. 缓存未命中则查数据库
//  3. 查库后同步回写全量缓存
func (r *cachedSettingQueryRepository) FindAll(ctx context.Context) ([]*setting.Setting, error) {
	// 1. 尝试从缓存获取（不过滤，返回全部）
	if result, ok := r.tryFilterFromCache(ctx, func(*setting.Setting) bool {
		return true
	}); ok {
		return result, nil
	}

	// 2. 缓存未命中，查数据库
	result, err := r.delegate.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	// 3. 同步回写缓存
	r.trySetAllCache(ctx, result, "operation", "findAll")

	return result, nil
}

// ExistsByKey 检查 Key 是否已存在（带缓存）。
func (r *cachedSettingQueryRepository) ExistsByKey(ctx context.Context, key string) (bool, error) {
	// 1. 尝试从缓存判断
	cached, err := r.cache.Get(ctx, key)
	if err == nil && cached != nil {
		return true, nil
	}

	// 2. 缓存未命中，查数据库
	return r.delegate.ExistsByKey(ctx, key)
}

// tryFilterFromCache 尝试从全量缓存中过滤数据。
// 返回 (result, true) 表示缓存命中，(nil, false) 表示需要查库。
func (r *cachedSettingQueryRepository) tryFilterFromCache(
	ctx context.Context,
	filter func(*setting.Setting) bool,
) ([]*setting.Setting, bool) {
	cachedMap, err := r.cache.GetAll(ctx)
	if err != nil {
		slog.Warn("cache get all failed", "err", err)
		return nil, false
	}

	if len(cachedMap) == 0 {
		// 缓存为空，检查是否已预热
		if r.cache.IsWarmedUp(ctx) {
			return []*setting.Setting{}, true // 已预热，空结果是有效的
		}
		return nil, false // 未预热，需要查库
	}

	// 从缓存过滤
	result := make([]*setting.Setting, 0)
	for _, s := range cachedMap {
		if filter(s) {
			result = append(result, s)
		}
	}

	// 有结果直接返回
	if len(result) > 0 {
		return result, true
	}

	// 无结果但已预热，说明确实没有匹配的数据
	if r.cache.IsWarmedUp(ctx) {
		return result, true
	}

	// 无结果且未预热，需要查库确认
	return nil, false
}

// trySetAllCache 尝试回写缓存，失败时记录警告日志。
func (r *cachedSettingQueryRepository) trySetAllCache(ctx context.Context, result []*setting.Setting, keyName string, keyValue any) {
	if len(result) > 0 {
		if err := r.cache.SetAll(ctx, result); err != nil {
			slog.Warn("cache set all failed", keyName, keyValue, "err", err)
		}
	}
}

var _ setting.QueryRepository = (*cachedSettingQueryRepository)(nil)
