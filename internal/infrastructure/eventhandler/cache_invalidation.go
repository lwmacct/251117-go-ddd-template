package eventhandler

import (
	"context"
	"log/slog"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/event"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/event/events"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/auth"
)

// CacheInvalidationHandler 缓存失效处理器
// 处理角色权限变更和用户角色分配事件，自动失效相关缓存
type CacheInvalidationHandler struct {
	permissionCache *auth.PermissionCacheService
	userQueryRepo   user.QueryRepository
	logger          *slog.Logger
}

// NewCacheInvalidationHandler 创建缓存失效处理器
func NewCacheInvalidationHandler(
	permissionCache *auth.PermissionCacheService,
	userQueryRepo user.QueryRepository,
) *CacheInvalidationHandler {
	return &CacheInvalidationHandler{
		permissionCache: permissionCache,
		userQueryRepo:   userQueryRepo,
		logger:          slog.Default(),
	}
}

// Handle 处理事件
func (h *CacheInvalidationHandler) Handle(ctx context.Context, e event.Event) error {
	switch evt := e.(type) {
	case *events.UserRoleAssignedEvent:
		return h.handleUserRoleAssigned(ctx, evt)
	case *events.RolePermissionsChangedEvent:
		return h.handleRolePermissionsChanged(ctx, evt)
	case *events.UserDeletedEvent:
		return h.handleUserDeleted(ctx, evt)
	default:
		// 忽略不处理的事件
		return nil
	}
}

// handleUserRoleAssigned 处理用户角色分配事件
// 失效单个用户的权限缓存
func (h *CacheInvalidationHandler) handleUserRoleAssigned(ctx context.Context, evt *events.UserRoleAssignedEvent) error {
	h.logger.Info("invalidating permission cache for user",
		"event", evt.EventName(),
		"user_id", evt.UserID,
	)

	if err := h.permissionCache.InvalidateUser(ctx, evt.UserID); err != nil {
		h.logger.Error("failed to invalidate user permission cache",
			"user_id", evt.UserID,
			"error", err,
		)
		// 缓存失效失败不应该阻塞业务流程，只记录错误
		return nil
	}

	return nil
}

// handleRolePermissionsChanged 处理角色权限变更事件
// 失效拥有该角色的所有用户的权限缓存
func (h *CacheInvalidationHandler) handleRolePermissionsChanged(ctx context.Context, evt *events.RolePermissionsChangedEvent) error {
	h.logger.Info("invalidating permission cache for role",
		"event", evt.EventName(),
		"role_id", evt.RoleID,
	)

	// 获取拥有该角色的所有用户
	userIDs, err := h.userQueryRepo.GetUserIDsByRole(ctx, evt.RoleID)
	if err != nil {
		h.logger.Error("failed to get users by role",
			"role_id", evt.RoleID,
			"error", err,
		)
		return nil
	}

	// 逐个失效用户缓存
	for _, userID := range userIDs {
		if err := h.permissionCache.InvalidateUser(ctx, userID); err != nil {
			h.logger.Error("failed to invalidate user permission cache",
				"user_id", userID,
				"error", err,
			)
			// 继续处理其他用户
		}
	}

	h.logger.Info("invalidated permission cache for users",
		"role_id", evt.RoleID,
		"user_count", len(userIDs),
	)

	return nil
}

// handleUserDeleted 处理用户删除事件
// 清理用户相关缓存
func (h *CacheInvalidationHandler) handleUserDeleted(ctx context.Context, evt *events.UserDeletedEvent) error {
	h.logger.Info("cleaning up cache for deleted user",
		"event", evt.EventName(),
		"user_id", evt.UserID,
	)

	if err := h.permissionCache.InvalidateUser(ctx, evt.UserID); err != nil {
		h.logger.Error("failed to invalidate deleted user cache",
			"user_id", evt.UserID,
			"error", err,
		)
	}

	return nil
}

// Ensure interface is implemented
var _ event.EventHandler = (*CacheInvalidationHandler)(nil)
