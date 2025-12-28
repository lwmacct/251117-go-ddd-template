//nolint:dupl // 预热器遵循统一模式，与 SettingCacheWarmer 结构相同是设计意图
package cache

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	domaincache "github.com/lwmacct/251117-go-ddd-template/internal/domain/cache"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

const (
	settingCategoryWarmupTimeout     = 30 * time.Second
	settingCategoryWaitWarmupTimeout = 30 * time.Second
)

// SettingCategoryCacheWarmer 配置分类缓存预热服务。
//
// 在系统启动时将所有 SettingCategory 加载到 Redis 缓存，确保后续查询可以直接从缓存获取。
//
// 多实例安全：
//   - 使用分布式锁防止多实例同时预热
//   - 未获取锁的实例会等待预热完成
//   - 双重检查避免不必要的数据库查询
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
//  1. 尝试获取分布式锁
//  2. 双重检查是否已预热
//  3. 从数据库加载所有 SettingCategory
//  4. 批量写入缓存
//  5. 标记预热完成
//
// 如果未获取到锁，会等待其他实例完成预热。
// 预热失败不会阻塞服务启动，会降级为惰性加载。
func (w *SettingCategoryCacheWarmer) WarmUp(ctx context.Context) error {
	// 1. 尝试获取分布式锁
	acquired, release := w.cacheService.TryAcquireWarmupLock(ctx)
	if !acquired {
		slog.Info("SettingCategory cache warmup lock not acquired, waiting for other instance")
		return w.waitForWarmUp(ctx)
	}
	defer release()

	// 2. 双重检查（获取锁后再次检查）
	if w.cacheService.IsWarmedUp(ctx) {
		slog.Info("SettingCategory cache already warmed up by another instance")
		return nil
	}

	// 3. 执行预热
	start := time.Now()
	slog.Info("Starting SettingCategory cache warmup")

	categories, err := w.categoryQueryRepo.FindAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to load setting categories for warmup: %w", err)
	}

	if len(categories) == 0 {
		slog.Warn("No setting categories found for cache warmup")
		// 即使没有数据也标记预热完成
		if err := w.cacheService.SetWarmedUp(ctx); err != nil {
			return fmt.Errorf("failed to set warmup flag: %w", err)
		}
		return nil
	}

	// 4. 批量写入缓存
	if err := w.cacheService.SetAll(ctx, categories); err != nil {
		return fmt.Errorf("failed to warmup cache: %w", err)
	}

	// 5. 标记预热完成
	if err := w.cacheService.SetWarmedUp(ctx); err != nil {
		return fmt.Errorf("failed to set warmup flag: %w", err)
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

// waitForWarmUp 等待缓存预热完成。
func (w *SettingCategoryCacheWarmer) waitForWarmUp(ctx context.Context) error {
	waitCtx, cancel := context.WithTimeout(ctx, settingCategoryWaitWarmupTimeout)
	defer cancel()

	if err := w.cacheService.WaitForWarmup(waitCtx); err != nil {
		return fmt.Errorf("timeout waiting for setting category cache warmup: %w", err)
	}

	slog.Info("SettingCategory cache warmup completed by another instance")
	return nil
}
