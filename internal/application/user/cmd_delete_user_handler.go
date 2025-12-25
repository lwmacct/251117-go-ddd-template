package user

import (
	"context"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/event"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/event/events"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
)

// DeleteUserHandler 删除用户命令处理器
type DeleteUserHandler struct {
	userCommandRepo user.CommandRepository
	userQueryRepo   user.QueryRepository
	eventBus        event.EventBus
}

// NewDeleteUserHandler 创建删除用户命令处理器
func NewDeleteUserHandler(
	userCommandRepo user.CommandRepository,
	userQueryRepo user.QueryRepository,
	eventBus event.EventBus,
) *DeleteUserHandler {
	return &DeleteUserHandler{
		userCommandRepo: userCommandRepo,
		userQueryRepo:   userQueryRepo,
		eventBus:        eventBus,
	}
}

// Handle 处理删除用户命令
func (h *DeleteUserHandler) Handle(ctx context.Context, cmd DeleteUserCommand) error {
	// 1. 检查用户是否存在
	exists, err := h.userQueryRepo.Exists(ctx, cmd.UserID)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if !exists {
		return user.ErrUserNotFound
	}

	// 2. 执行删除（软删除）
	if err := h.userCommandRepo.Delete(ctx, cmd.UserID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	// 3. 发布用户删除事件，触发缓存清理
	evt := events.NewUserDeletedEvent(cmd.UserID)
	if h.eventBus != nil {
		_ = h.eventBus.Publish(ctx, evt) // 缓存清理失败不阻塞业务
	}

	return nil
}
