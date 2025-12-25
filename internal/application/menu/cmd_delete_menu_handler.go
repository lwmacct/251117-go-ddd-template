package menu

import (
	"context"
	"errors"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/menu"
)

// DeleteMenuHandler 删除菜单命令处理器
type DeleteMenuHandler struct {
	menuCommandRepo menu.CommandRepository
	menuQueryRepo   menu.QueryRepository
}

// NewDeleteMenuHandler 创建 DeleteMenuHandler 实例
func NewDeleteMenuHandler(
	menuCommandRepo menu.CommandRepository,
	menuQueryRepo menu.QueryRepository,
) *DeleteMenuHandler {
	return &DeleteMenuHandler{
		menuCommandRepo: menuCommandRepo,
		menuQueryRepo:   menuQueryRepo,
	}
}

// Handle 处理删除菜单命令
func (h *DeleteMenuHandler) Handle(ctx context.Context, cmd DeleteMenuCommand) error {
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
