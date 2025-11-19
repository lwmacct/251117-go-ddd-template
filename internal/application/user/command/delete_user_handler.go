// Package command 定义用户命令处理器
package command

import (
	"context"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
)

// DeleteUserHandler 删除用户命令处理器
type DeleteUserHandler struct {
	userCommandRepo user.CommandRepository
	userQueryRepo   user.QueryRepository
}

// NewDeleteUserHandler 创建删除用户命令处理器
func NewDeleteUserHandler(
	userCommandRepo user.CommandRepository,
	userQueryRepo user.QueryRepository,
) *DeleteUserHandler {
	return &DeleteUserHandler{
		userCommandRepo: userCommandRepo,
		userQueryRepo:   userQueryRepo,
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

	return nil
}
