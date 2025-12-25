package bootstrap

import (
	"log/slog"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/event"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/eventhandler"
)

// initEventHandlers 初始化事件处理器并订阅事件
// 依赖：InfrastructureModule.EventBus, RepositoriesModule, ServicesModule
func initEventHandlers(eventBus event.EventBus, repos *RepositoriesModule, services *ServicesModule) {
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

	// 订阅审计日志事件（使用通配符订阅所有事件）
	eventBus.Subscribe("*", auditHandler)

	slog.Info("Event handlers initialized",
		"handlers", []string{"CacheInvalidationHandler", "AuditLogHandler"},
		"cache_subscriptions", []string{"user.role_assigned", "user.deleted", "role.permissions_changed"},
		"audit_subscriptions", []string{"*"},
	)
}
