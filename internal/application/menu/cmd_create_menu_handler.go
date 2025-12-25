package menu

import (
	"context"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/menu"
)

// CreateMenuHandler 创建菜单命令处理器
type CreateMenuHandler struct {
	menuCommandRepo menu.CommandRepository
	menuQueryRepo   menu.QueryRepository
}

// NewCreateMenuHandler 创建 CreateMenuHandler 实例
func NewCreateMenuHandler(
	menuCommandRepo menu.CommandRepository,
	menuQueryRepo menu.QueryRepository,
) *CreateMenuHandler {
	return &CreateMenuHandler{
		menuCommandRepo: menuCommandRepo,
		menuQueryRepo:   menuQueryRepo,
	}
}

// Handle 处理创建菜单命令
func (h *CreateMenuHandler) Handle(ctx context.Context, cmd CreateMenuCommand) (*CreateMenuResultDTO, error) {
	// 1. 如果指定了父菜单，验证父菜单是否存在
	if cmd.ParentID != nil {
		_, err := h.menuQueryRepo.FindByID(ctx, *cmd.ParentID)
		if err != nil {
			return nil, fmt.Errorf("parent menu not found: %w", err)
		}
	}

	// 2. 创建菜单实体
	menuEntity := &menu.Menu{
		Title:    cmd.Title,
		Path:     cmd.Path,
		Icon:     cmd.Icon,
		ParentID: cmd.ParentID,
		Order:    cmd.Order,
		Visible:  cmd.Visible,
	}

	// 3. 保存菜单
	if err := h.menuCommandRepo.Create(ctx, menuEntity); err != nil {
		return nil, fmt.Errorf("failed to create menu: %w", err)
	}

	return &CreateMenuResultDTO{
		ID: menuEntity.ID,
	}, nil
}
