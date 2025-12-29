package menu

import (
	"context"
	"errors"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/menu"
)

// DeleteHandler 删除菜单命令处理器
type DeleteHandler struct {
	menuCommandRepo menu.CommandRepository
	menuQueryRepo   menu.QueryRepository
}

// NewDeleteHandler 创建 DeleteHandler 实例
func NewDeleteHandler(
	menuCommandRepo menu.CommandRepository,
	menuQueryRepo menu.QueryRepository,
) *DeleteHandler {
	return &DeleteHandler{
		menuCommandRepo: menuCommandRepo,
		menuQueryRepo:   menuQueryRepo,
	}
}

// Handle 处理删除菜单命令
func (h *DeleteHandler) Handle(ctx context.Context, cmd DeleteCommand) error {
	// 1. 验证菜单是否存在
	_, err := h.menuQueryRepo.FindByID(ctx, cmd.MenuID)
	if err != nil {
		return fmt.Errorf("menu not found: %w", err)
	}

	// 2. 检查是否有子菜单
	parentID := &cmd.MenuID
	children, err := h.menuQueryRepo.FindByParentID(ctx, parentID)
	if err != nil {
		return fmt.Errorf("failed to check children: %w", err)
	}

	if len(children) > 0 {
		return errors.New("cannot delete menu with children")
	}

	// 3. 删除菜单
	if err := h.menuCommandRepo.Delete(ctx, cmd.MenuID); err != nil {
		return fmt.Errorf("failed to delete menu: %w", err)
	}

	return nil
}
