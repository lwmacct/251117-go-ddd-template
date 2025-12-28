package container

import (
	"context"
	"log/slog"

	"go.uber.org/fx"
	"gorm.io/gorm"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/event"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/cache/warmup"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/eventhandler"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/persistence"
)

// RegisterEventHandlers 设置缓存失效和审计日志的事件订阅。
//
// 订阅事件：
//   - user.role_assigned, user.deleted, role.permissions_changed → 缓存失效
//   - *（所有事件）→ 审计日志
func RegisterEventHandlers(
	eventBus event.EventBus,
	repos *RepositoriesModule,
	services *ServicesModule,
) {
	// 缓存失效处理器
	cacheHandler := eventhandler.NewCacheInvalidationHandler(
		services.PermissionCache,
		repos.User.Query,
	)

	// 审计日志处理器
	auditHandler := eventhandler.NewAuditLogHandler(repos.AuditLog.Command)

	// 订阅缓存失效事件
	eventBus.Subscribe("user.role_assigned", cacheHandler)
	eventBus.Subscribe("user.deleted", cacheHandler)
	eventBus.Subscribe("role.permissions_changed", cacheHandler)

	// 订阅所有事件用于审计日志
	eventBus.Subscribe("*", auditHandler)

	slog.Info("Event handlers initialized",
		"handlers", []string{"CacheInvalidationHandler", "AuditLogHandler"},
		"cache_subscriptions", []string{"user.role_assigned", "user.deleted", "role.permissions_changed"},
		"audit_subscriptions", []string{"*"},
	)
}

// WarmupCaches 预热缓存，加载常用数据。
//
// 预热的缓存：
//   - Setting 缓存
//
// 注意：使用原始仓储（非缓存装饰器）以避免循环依赖。
// 预热失败会记录日志但不阻塞应用启动。
func WarmupCaches(lc fx.Lifecycle, db *gorm.DB, cache *CacheServicesModule) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			warmupSettingCache(ctx, db, cache)
			return nil // 预热失败不影响启动
		},
	})
}

// warmupSettingCache 预热 Setting 缓存。
func warmupSettingCache(ctx context.Context, db *gorm.DB, cache *CacheServicesModule) {
	// 使用原始仓储以避免与缓存装饰器的循环依赖
	rawRepos := persistence.NewSettingRepositories(db)
	warmer := warmup.NewSettingCacheWarmer(rawRepos.Query, cache.Setting)

	if err := warmer.WarmUpWithTimeout(ctx); err != nil {
		slog.Warn("Setting cache warmup failed, will use lazy loading", "err", err)
	} else {
		slog.Info("Setting cache warmup completed")
	}
}
