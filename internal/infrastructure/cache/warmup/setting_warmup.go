//nolint:dupl // 预热器遵循统一模式，与 SettingCategoryCacheWarmer 结构相同是设计意图
package warmup

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	domaincache "github.com/lwmacct/251117-go-ddd-template/internal/domain/cache"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

const warmupTimeout = 30 * time.Second

// SettingCacheWarmer 系统设置缓存预热服务。
//
// 在系统启动时将所有 Setting 加载到 Redis 缓存，确保后续查询可以直接从缓存获取。
//
// 设计简化：
//   - 无分布式锁：多实例同时预热是幂等操作，数据量小，重复写入可接受
//   - 无预热标记：每次启动都执行预热，确保缓存数据最新
type SettingCacheWarmer struct {
	settingQueryRepo setting.QueryRepository // 原始仓储（非缓存装饰器）
	cacheService     domaincache.SettingCacheService
}

// NewSettingCacheWarmer 创建设置缓存预热服务。
//
// 注意：settingQueryRepo 应传入原始仓储，而非缓存装饰器，避免循环依赖。
func NewSettingCacheWarmer(
	settingQueryRepo setting.QueryRepository,
	cacheService domaincache.SettingCacheService,
) *SettingCacheWarmer {
	return &SettingCacheWarmer{
		settingQueryRepo: settingQueryRepo,
		cacheService:     cacheService,
	}
}

// WarmUp 执行缓存预热。
//
// 预热流程：
//  1. 从数据库加载所有 Setting
//  2. 批量写入缓存
//
// 预热失败不会阻塞服务启动，会降级为惰性加载。
func (w *SettingCacheWarmer) WarmUp(ctx context.Context) error {
	start := time.Now()
	slog.Info("Starting Setting cache warmup")

	settings, err := w.settingQueryRepo.FindAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to load settings for warmup: %w", err)
	}

	if len(settings) == 0 {
		slog.Warn("No settings found for cache warmup")
		return nil
	}

	if err := w.cacheService.SetAll(ctx, settings); err != nil {
		return fmt.Errorf("failed to warmup cache: %w", err)
	}

	slog.Info("Setting cache warmup completed",
		"count", len(settings),
		"duration", time.Since(start),
	)

	return nil
}

// WarmUpWithTimeout 带超时的预热。
//
// 这是 WarmUp 的便捷包装，使用默认超时时间。
func (w *SettingCacheWarmer) WarmUpWithTimeout(ctx context.Context) error {
	warmupCtx, cancel := context.WithTimeout(ctx, warmupTimeout)
	defer cancel()

	return w.WarmUp(warmupCtx)
}
