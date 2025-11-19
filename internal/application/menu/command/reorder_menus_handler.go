// Package command 定义菜单命令处理器
package command

import (
	"context"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/menu"
)

// ReorderMenusHandler 批量更新菜单排序命令处理器
type ReorderMenusHandler struct {
	menuCommandRepo menu.CommandRepository
	menuQueryRepo   menu.QueryRepository
}

// NewReorderMenusHandler 创建 ReorderMenusHandler 实例
func NewReorderMenusHandler(
	menuCommandRepo menu.CommandRepository,
	menuQueryRepo menu.QueryRepository,
) *ReorderMenusHandler {
	return &ReorderMenusHandler{
		menuCommandRepo: menuCommandRepo,
		menuQueryRepo:   menuQueryRepo,
	}
}

// Handle 处理批量更新菜单排序命令
func (h *ReorderMenusHandler) Handle(ctx context.Context, cmd ReorderMenusCommand) error {
	// 1. 验证所有菜单是否存在
	for _, item := range cmd.Menus {
		_, err := h.menuQueryRepo.FindByID(ctx, item.ID)
		if err != nil {
			return fmt.Errorf("menu %d not found: %w", item.ID, err)
		}

		// 如果指定了父菜单，验证父菜单是否存在
		if item.ParentID != nil && *item.ParentID != 0 {
			_, err := h.menuQueryRepo.FindByID(ctx, *item.ParentID)
			if err != nil {
				return fmt.Errorf("parent menu %d not found: %w", *item.ParentID, err)
			}
		}
	}

	// 2. 转换为 Repository 需要的格式
	reorderData := make([]struct {
		ID       uint
		Order    int
		ParentID *uint
	}, len(cmd.Menus))

	for i, item := range cmd.Menus {
		reorderData[i].ID = item.ID
		reorderData[i].Order = item.Order
		reorderData[i].ParentID = item.ParentID
	}

	// 3. 批量更新排序
	if err := h.menuCommandRepo.UpdateOrder(ctx, reorderData); err != nil {
		return fmt.Errorf("failed to update menu order: %w", err)
	}

	return nil
}
