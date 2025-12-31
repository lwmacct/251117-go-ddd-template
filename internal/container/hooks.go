package container

import (
	"log/slog"

	"go.uber.org/fx"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/event"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/auth"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/eventhandler"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/persistence"
)

// eventHandlersParams 聚合事件处理器所需的依赖。
type eventHandlersParams struct {
	fx.In

	EventBus        event.EventBus
	PermissionCache *auth.PermissionCacheService
	UserRepos       persistence.UserRepositories
	AuditRepos      persistence.AuditRepositories
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
	)

	// 审计日志处理器
	auditHandler := eventhandler.NewAuditEventHandler(p.AuditRepos.Command)

	// 订阅缓存失效事件
	p.EventBus.Subscribe("user.role_assigned", cacheHandler)
	p.EventBus.Subscribe("user.deleted", cacheHandler)
	p.EventBus.Subscribe("role.permissions_changed", cacheHandler)

	// 订阅所有事件用于审计日志
	p.EventBus.Subscribe("*", auditHandler)

	slog.Info("Event handlers initialized",
		"handlers", []string{"CacheInvalidationHandler", "AuditEventHandler"},
		"cache_subscriptions", []string{"user.role_assigned", "user.deleted", "role.permissions_changed"},
		"audit_subscriptions", []string{"*"},
	)
}
