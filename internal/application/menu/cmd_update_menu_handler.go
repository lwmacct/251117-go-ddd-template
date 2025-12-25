package menu

import (
	"context"
	"errors"
	"fmt"

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
func (h *UpdateMenuHandler) Handle(ctx context.Context, cmd UpdateMenuCommand) (*MenuDTO, error) {
	// 1. 查询菜单
	menuEntity, err := h.menuQueryRepo.FindByID(ctx, cmd.MenuID)
	if err != nil {
		return nil, fmt.Errorf("menu not found: %w", err)
	}

	// 2. 如果更新父菜单，验证父菜单是否存在且不能设置自己为父菜单
	if err := h.validateParentID(ctx, cmd.ParentID, menuEntity.ID); err != nil {
		return nil, err
	}

	// 3. 更新字段
	h.applyUpdates(menuEntity, &cmd)

	// 4. 保存更新
	if err := h.menuCommandRepo.Update(ctx, menuEntity); err != nil {
		return nil, fmt.Errorf("failed to update menu: %w", err)
	}

	return ToMenuDTO(menuEntity), nil
}

// validateParentID 验证父菜单 ID
func (h *UpdateMenuHandler) validateParentID(ctx context.Context, parentID *uint, selfID uint) error {
	if parentID == nil {
		return nil
	}
	if *parentID == selfID {
		return errors.New("cannot set menu as its own parent")
	}
	if *parentID == 0 {
		return nil
	}
	if _, err := h.menuQueryRepo.FindByID(ctx, *parentID); err != nil {
		return fmt.Errorf("parent menu not found: %w", err)
	}
	return nil
}

// applyUpdates 应用更新字段
func (h *UpdateMenuHandler) applyUpdates(menuEntity *menu.Menu, cmd *UpdateMenuCommand) {
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
}
