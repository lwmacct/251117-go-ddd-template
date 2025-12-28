package container

import (
	"context"
	"log/slog"

	"go.uber.org/fx"
	"gorm.io/gorm"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/cache"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/event"
	infra_auth "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/auth"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/cache/warmup"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/eventhandler"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/persistence"
)

// eventHandlersParams 聚合事件处理器所需的依赖。
type eventHandlersParams struct {
	fx.In

	EventBus        event.EventBus
	PermissionCache *infra_auth.PermissionCacheService
	UserRepos       persistence.UserRepositories
	AuditLogRepos   persistence.AuditLogRepositories
}

// RegisterEventHandlers 设置缓存失效和审计日志的事件订阅。
//
// 订阅事件：
//   - user.role_assigned, user.deleted, role.permissions_changed → 缓存失效
//   - *（所有事件）→ 审计日志
func RegisterEventHandlers(p eventHandlersParams) {
	// 缓存失效处理器
	cacheHandler := eventhandler.NewCacheInvalidationHandler(
		p.PermissionCache,
		p.UserRepos.Query,
	)

	// 审计日志处理器
	auditHandler := eventhandler.NewAuditLogHandler(p.AuditLogRepos.Command)

	// 订阅缓存失效事件
	p.EventBus.Subscribe("user.role_assigned", cacheHandler)
	p.EventBus.Subscribe("user.deleted", cacheHandler)
	p.EventBus.Subscribe("role.permissions_changed", cacheHandler)

	// 订阅所有事件用于审计日志
	p.EventBus.Subscribe("*", auditHandler)

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
func WarmupCaches(lc fx.Lifecycle, db *gorm.DB, settingCache cache.SettingCacheService) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			warmupSettingCache(ctx, db, settingCache)
			return nil // 预热失败不影响启动
		},
	})
}

// warmupSettingCache 预热 Setting 缓存。
func warmupSettingCache(ctx context.Context, db *gorm.DB, settingCache cache.SettingCacheService) {
	// 使用原始仓储以避免与缓存装饰器的循环依赖
	rawRepos := persistence.NewSettingRepositories(db)
	warmer := warmup.NewSettingCacheWarmer(rawRepos.Query, settingCache)

	if err := warmer.WarmUpWithTimeout(ctx); err != nil {
		slog.Warn("Setting cache warmup failed, will use lazy loading", "err", err)
	} else {
		slog.Info("Setting cache warmup completed")
	}
}
