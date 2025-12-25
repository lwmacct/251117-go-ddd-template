package user

import (
	"context"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/event"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/event/events"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
)

// AssignRolesHandler 负责分配用户角色
type AssignRolesHandler struct {
	userCommandRepo user.CommandRepository
	userQueryRepo   user.QueryRepository
	eventBus        event.EventBus
}

// NewAssignRolesHandler 创建新的分配角色处理器
func NewAssignRolesHandler(
	userCommandRepo user.CommandRepository,
	userQueryRepo user.QueryRepository,
	eventBus event.EventBus,
) *AssignRolesHandler {
	return &AssignRolesHandler{
		userCommandRepo: userCommandRepo,
		userQueryRepo:   userQueryRepo,
		eventBus:        eventBus,
	}
}

// Handle 处理分配角色命令
func (h *AssignRolesHandler) Handle(ctx context.Context, cmd AssignRolesCommand) error {
	if _, err := h.userQueryRepo.GetByID(ctx, cmd.UserID); err != nil {
		return err
	}

	if err := h.userCommandRepo.AssignRoles(ctx, cmd.UserID, cmd.RoleIDs); err != nil {
		return err
	}

	// 发布用户角色分配事件，触发缓存失效
	evt := events.NewUserRoleAssignedEvent(cmd.UserID, cmd.RoleIDs)
	if h.eventBus != nil {
		_ = h.eventBus.Publish(ctx, evt) // 缓存失效失败不阻塞业务
	}

	return nil
}
