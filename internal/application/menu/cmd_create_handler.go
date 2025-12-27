package menu

import (
	"context"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/menu"
)

// CreateHandler 创建菜单命令处理器
type CreateHandler struct {
	menuCommandRepo menu.CommandRepository
	menuQueryRepo   menu.QueryRepository
}

// NewCreateHandler 创建 CreateHandler 实例
func NewCreateHandler(
	menuCommandRepo menu.CommandRepository,
	menuQueryRepo menu.QueryRepository,
) *CreateHandler {
	return &CreateHandler{
		menuCommandRepo: menuCommandRepo,
		menuQueryRepo:   menuQueryRepo,
	}
}

// Handle 处理创建菜单命令
func (h *CreateHandler) Handle(ctx context.Context, cmd CreateCommand) (*CreateResultDTO, error) {
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

	return &CreateResultDTO{
		ID: menuEntity.ID,
	}, nil
}
