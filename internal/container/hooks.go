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

// RegisterEventHandlers sets up event subscriptions for cache invalidation and audit logging.
//
// Subscriptions:
//   - user.role_assigned, user.deleted, role.permissions_changed → cache invalidation
//   - * (all events) → audit logging
func RegisterEventHandlers(
	eventBus event.EventBus,
	repos *RepositoriesModule,
	services *ServicesModule,
) {
	// Cache invalidation handler
	cacheHandler := eventhandler.NewCacheInvalidationHandler(
		services.PermissionCache,
		repos.User.Query,
	)

	// Audit log handler
	auditHandler := eventhandler.NewAuditLogHandler(repos.AuditLog.Command)

	// Subscribe to cache invalidation events
	eventBus.Subscribe("user.role_assigned", cacheHandler)
	eventBus.Subscribe("user.deleted", cacheHandler)
	eventBus.Subscribe("role.permissions_changed", cacheHandler)

	// Subscribe to all events for audit logging
	eventBus.Subscribe("*", auditHandler)

	slog.Info("Event handlers initialized",
		"handlers", []string{"CacheInvalidationHandler", "AuditLogHandler"},
		"cache_subscriptions", []string{"user.role_assigned", "user.deleted", "role.permissions_changed"},
		"audit_subscriptions", []string{"*"},
	)
}

// WarmupCaches pre-populates caches with frequently accessed data.
//
// Caches warmed:
//   - Setting cache
//   - SettingCategory cache
//
// Note: Uses raw repositories (not cache decorators) to avoid circular dependencies.
// Warmup failures are logged but don't block application startup.
func WarmupCaches(lc fx.Lifecycle, db *gorm.DB, cache *CacheServicesModule) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			warmupSettingCache(ctx, db, cache)
			warmupSettingCategoryCache(ctx, db, cache)
			return nil // Never fail startup due to warmup
		},
	})
}

// warmupSettingCache pre-populates the Setting cache.
func warmupSettingCache(ctx context.Context, db *gorm.DB, cache *CacheServicesModule) {
	// Use raw repository to avoid circular dependency with cache decorator
	rawRepos := persistence.NewSettingRepositories(db)
	warmer := warmup.NewSettingCacheWarmer(rawRepos.Query, cache.Setting)

	if err := warmer.WarmUpWithTimeout(ctx); err != nil {
		slog.Warn("Setting cache warmup failed, will use lazy loading", "err", err)
	} else {
		slog.Info("Setting cache warmup completed")
	}
}

// warmupSettingCategoryCache pre-populates the SettingCategory cache.
func warmupSettingCategoryCache(ctx context.Context, db *gorm.DB, cache *CacheServicesModule) {
	// Use raw repository to avoid circular dependency with cache decorator
	rawRepos := persistence.NewSettingRepositories(db)
	warmer := warmup.NewSettingCategoryCacheWarmer(rawRepos.CategoryQuery, cache.SettingCategory)

	if err := warmer.WarmUpWithTimeout(ctx); err != nil {
		slog.Warn("SettingCategory cache warmup failed, will use lazy loading", "err", err)
	} else {
		slog.Info("SettingCategory cache warmup completed")
	}
}
