// Package command 定义菜单命令处理器
package command

import (
	"context"
	"fmt"

	menuApp "github.com/lwmacct/251117-go-ddd-template/internal/application/menu"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/menu"
)

// UpdateMenuHandler 更新菜单命令处理器
type UpdateMenuHandler struct {
	menuCommandRepo menu.CommandRepository
	menuQueryRepo   menu.QueryRepository
}

// NewUpdateMenuHandler 创建 UpdateMenuHandler 实例
func NewUpdateMenuHandler(
	menuCommandRepo menu.CommandRepository,
	menuQueryRepo menu.QueryRepository,
) *UpdateMenuHandler {
	return &UpdateMenuHandler{
		menuCommandRepo: menuCommandRepo,
		menuQueryRepo:   menuQueryRepo,
	}
}

// Handle 处理更新菜单命令
func (h *UpdateMenuHandler) Handle(ctx context.Context, cmd UpdateMenuCommand) (*menuApp.MenuResponse, error) {
	// 1. 查询菜单
	menuEntity, err := h.menuQueryRepo.FindByID(ctx, cmd.MenuID)
	if err != nil {
		return nil, fmt.Errorf("menu not found: %w", err)
	}

	// 2. 如果更新父菜单，验证父菜单是否存在且不能设置自己为父菜单
	if cmd.ParentID != nil {
		if *cmd.ParentID == menuEntity.ID {
			return nil, fmt.Errorf("cannot set menu as its own parent")
		}
		if *cmd.ParentID != 0 {
			_, err := h.menuQueryRepo.FindByID(ctx, *cmd.ParentID)
			if err != nil {
				return nil, fmt.Errorf("parent menu not found: %w", err)
			}
		}
	}

	// 3. 更新字段
	if cmd.Title != nil {
		menuEntity.Title = *cmd.Title
	}
	if cmd.Path != nil {
		menuEntity.Path = *cmd.Path
	}
	if cmd.Icon != nil {
		menuEntity.Icon = *cmd.Icon
	}
	if cmd.ParentID != nil {
		menuEntity.ParentID = cmd.ParentID
	}
	if cmd.Order != nil {
		menuEntity.Order = *cmd.Order
	}
	if cmd.Visible != nil {
		menuEntity.Visible = *cmd.Visible
	}

	// 4. 保存更新
	if err := h.menuCommandRepo.Update(ctx, menuEntity); err != nil {
		return nil, fmt.Errorf("failed to update menu: %w", err)
	}

	return menuApp.ToMenuResponse(menuEntity), nil
}
