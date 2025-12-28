//nolint:dupl // 预热器遵循统一模式，与 SettingCacheWarmer 结构相同是设计意图
package warmup

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	domaincache "github.com/lwmacct/251117-go-ddd-template/internal/domain/cache"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

const settingCategoryWarmupTimeout = 30 * time.Second

// SettingCategoryCacheWarmer 配置分类缓存预热服务。
//
// 在系统启动时将所有 SettingCategory 加载到 Redis 缓存，确保后续查询可以直接从缓存获取。
//
// 设计简化：
//   - 无分布式锁：多实例同时预热是幂等操作，数据量小，重复写入可接受
//   - 无预热标记：每次启动都执行预热，确保缓存数据最新
type SettingCategoryCacheWarmer struct {
	categoryQueryRepo setting.SettingCategoryQueryRepository // 原始仓储（非缓存装饰器）
	cacheService      domaincache.SettingCategoryCacheService
}

// NewSettingCategoryCacheWarmer 创建配置分类缓存预热服务。
//
// 注意：categoryQueryRepo 应传入原始仓储，而非缓存装饰器，避免循环依赖。
func NewSettingCategoryCacheWarmer(
	categoryQueryRepo setting.SettingCategoryQueryRepository,
	cacheService domaincache.SettingCategoryCacheService,
) *SettingCategoryCacheWarmer {
	return &SettingCategoryCacheWarmer{
		categoryQueryRepo: categoryQueryRepo,
		cacheService:      cacheService,
	}
}

// WarmUp 执行缓存预热。
//
// 预热流程：
//  1. 从数据库加载所有 SettingCategory
//  2. 批量写入缓存
//
// 预热失败不会阻塞服务启动，会降级为惰性加载。
func (w *SettingCategoryCacheWarmer) WarmUp(ctx context.Context) error {
	start := time.Now()
	slog.Info("Starting SettingCategory cache warmup")

	categories, err := w.categoryQueryRepo.FindAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to load setting categories for warmup: %w", err)
	}

	if len(categories) == 0 {
		slog.Warn("No setting categories found for cache warmup")
		return nil
	}

	if err := w.cacheService.SetAll(ctx, categories); err != nil {
		return fmt.Errorf("failed to warmup cache: %w", err)
	}

	slog.Info("SettingCategory cache warmup completed",
		"count", len(categories),
		"duration", time.Since(start),
	)

	return nil
}

// WarmUpWithTimeout 带超时的预热。
//
// 这是 WarmUp 的便捷包装，使用默认超时时间。
func (w *SettingCategoryCacheWarmer) WarmUpWithTimeout(ctx context.Context) error {
	warmupCtx, cancel := context.WithTimeout(ctx, settingCategoryWarmupTimeout)
	defer cancel()

	return w.WarmUp(warmupCtx)
}
