package role

import (
	"context"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/event"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/event/events"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/role"
)

// SetPermissionsHandler 设置权限命令处理器
type SetPermissionsHandler struct {
	roleCommandRepo     role.CommandRepository
	roleQueryRepo       role.QueryRepository
	permissionQueryRepo role.PermissionQueryRepository
	eventBus            event.EventBus
}

// NewSetPermissionsHandler 创建设置权限命令处理器
func NewSetPermissionsHandler(
	roleCommandRepo role.CommandRepository,
	roleQueryRepo role.QueryRepository,
	permissionQueryRepo role.PermissionQueryRepository,
	eventBus event.EventBus,
) *SetPermissionsHandler {
	return &SetPermissionsHandler{
		roleCommandRepo:     roleCommandRepo,
		roleQueryRepo:       roleQueryRepo,
		permissionQueryRepo: permissionQueryRepo,
		eventBus:            eventBus,
	}
}

// Handle 处理设置权限命令
func (h *SetPermissionsHandler) Handle(ctx context.Context, cmd SetPermissionsCommand) error {
	// 1. 验证角色是否存在
	exists, err := h.roleQueryRepo.Exists(ctx, cmd.RoleID)
	if err != nil {
		return fmt.Errorf("failed to check role existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("role not found with id: %d", cmd.RoleID)
	}

	// 2. 验证所有权限ID是否有效
	for _, permID := range cmd.PermissionIDs {
		permExists, err := h.permissionQueryRepo.Exists(ctx, permID)
		if err != nil {
			return fmt.Errorf("failed to check permission existence: %w", err)
		}
		if !permExists {
			return fmt.Errorf("permission not found with id: %d", permID)
		}
	}

	// 3. 设置权限
	if err := h.roleCommandRepo.SetPermissions(ctx, cmd.RoleID, cmd.PermissionIDs); err != nil {
		return fmt.Errorf("failed to set permissions: %w", err)
	}

	// 4. 发布角色权限变更事件，触发缓存失效
	evt := events.NewRolePermissionsChangedEvent(cmd.RoleID, cmd.PermissionIDs)
	if h.eventBus != nil {
		_ = h.eventBus.Publish(ctx, evt) // 缓存失效失败不阻塞业务
	}

	return nil
}
