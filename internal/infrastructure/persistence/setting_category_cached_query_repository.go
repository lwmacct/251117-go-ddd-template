package persistence

import (
	"context"
	"log/slog"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/cache"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

// cachedSettingCategoryQueryRepository 带缓存的 SettingCategory 查询仓储装饰器。
//
// 装饰 [setting.SettingCategoryQueryRepository]，在查询前检查缓存，未命中再查数据库。
// 采用 Cache-Aside 模式，同步回写缓存。
//
// 缓存策略：
//   - [FindByID]: 先查缓存，未命中查库后同步回写
//   - [FindByKey]: 先查缓存，未命中查库后同步回写
//   - [FindByIDs]: 批量从缓存获取，未命中的查库后批量回写
//   - [FindAll]: 预热后从缓存获取全量，否则查库
type cachedSettingCategoryQueryRepository struct {
	delegate setting.SettingCategoryQueryRepository // 被装饰的原始仓储
	cache    cache.SettingCategoryCacheService      // 缓存服务
}

// NewCachedSettingCategoryQueryRepository 创建带缓存的 SettingCategory 查询仓储。
func NewCachedSettingCategoryQueryRepository(
	delegate setting.SettingCategoryQueryRepository,
	cacheService cache.SettingCategoryCacheService,
) setting.SettingCategoryQueryRepository {
	return &cachedSettingCategoryQueryRepository{
		delegate: delegate,
		cache:    cacheService,
	}
}

// FindByID 根据 ID 查找配置分类（带缓存）。
func (r *cachedSettingCategoryQueryRepository) FindByID(ctx context.Context, id uint) (*setting.SettingCategory, error) {
	// 1. 查缓存
	cached, err := r.cache.GetByID(ctx, id)
	if err != nil {
		slog.Warn("cache get by id failed, fallback to db", "id", id, "err", err)
	}
	if cached != nil {
		return cached, nil
	}

	// 2. 查数据库
	result, err := r.delegate.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil //nolint:nilnil // not found
	}

	// 3. 同步回写缓存
	if err := r.cache.Set(ctx, result); err != nil {
		slog.Warn("cache set failed", "id", result.ID, "key", result.Key, "err", err)
	}

	return result, nil
}

// FindByKey 根据 Key 查找配置分类（带缓存）。
func (r *cachedSettingCategoryQueryRepository) FindByKey(ctx context.Context, key string) (*setting.SettingCategory, error) {
	// 1. 查缓存
	cached, err := r.cache.GetByKey(ctx, key)
	if err != nil {
		slog.Warn("cache get by key failed, fallback to db", "key", key, "err", err)
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
		slog.Warn("cache set failed", "id", result.ID, "key", result.Key, "err", err)
	}

	return result, nil
}

// FindByIDs 批量查找配置分类（带缓存）。
func (r *cachedSettingCategoryQueryRepository) FindByIDs(ctx context.Context, ids []uint) ([]*setting.SettingCategory, error) {
	if len(ids) == 0 {
		return []*setting.SettingCategory{}, nil
	}

	// 1. 批量从缓存获取
	cachedList, err := r.cache.GetByIDs(ctx, ids)
	if err != nil {
		slog.Warn("cache batch get failed, fallback to db", "count", len(ids), "err", err)
		return r.delegate.FindByIDs(ctx, ids)
	}

	// 构建已缓存的 ID 集合
	cachedMap := make(map[uint]*setting.SettingCategory, len(cachedList))
	for _, c := range cachedList {
		cachedMap[c.ID] = c
	}

	// 2. 找出未命中的 ID
	var missedIDs []uint
	for _, id := range ids {
		if _, ok := cachedMap[id]; !ok {
			missedIDs = append(missedIDs, id)
		}
	}

	// 3. 查询未命中的
	if len(missedIDs) > 0 { //nolint:nestif // 缓存回写逻辑需要嵌套判断
		dbResults, err := r.delegate.FindByIDs(ctx, missedIDs)
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
		for _, c := range dbResults {
			cachedMap[c.ID] = c
		}
	}

	// 4. 按原始顺序返回
	result := make([]*setting.SettingCategory, 0, len(ids))
	for _, id := range ids {
		if c, ok := cachedMap[id]; ok {
			result = append(result, c)
		}
	}

	return result, nil
}

// FindAll 查找所有配置分类。
//
// 缓存策略：
//   - 缓存非空时，直接返回缓存结果
//   - 缓存为空时，回退到数据库查询
func (r *cachedSettingCategoryQueryRepository) FindAll(ctx context.Context) ([]*setting.SettingCategory, error) {
	// 尝试从缓存获取
	cachedList, err := r.cache.GetAll(ctx)
	if err != nil {
		slog.Warn("cache get all failed, fallback to db", "err", err)
		return r.delegate.FindAll(ctx)
	}

	// 缓存非空时返回
	if len(cachedList) > 0 {
		return cachedList, nil
	}

	// 空缓存总是回退到数据库
	return r.delegate.FindAll(ctx)
}

// ExistsByKey 检查 Key 是否已存在。
// 先查缓存，未命中再查数据库。
func (r *cachedSettingCategoryQueryRepository) ExistsByKey(ctx context.Context, key string) (bool, error) {
	// 先尝试从缓存判断
	cached, err := r.cache.GetByKey(ctx, key)
	if err == nil && cached != nil {
		return true, nil
	}

	// 缓存未命中，查数据库
	return r.delegate.ExistsByKey(ctx, key)
}

var _ setting.SettingCategoryQueryRepository = (*cachedSettingCategoryQueryRepository)(nil)
