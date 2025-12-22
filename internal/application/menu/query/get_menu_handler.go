package query

import (
	"context"
	"fmt"

	menuApp "github.com/lwmacct/251117-go-ddd-template/internal/application/menu"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/menu"
)

// GetMenuHandler 获取菜单查询处理器
type GetMenuHandler struct {
	menuQueryRepo menu.QueryRepository
}

// NewGetMenuHandler 创建 GetMenuHandler 实例
func NewGetMenuHandler(menuQueryRepo menu.QueryRepository) *GetMenuHandler {
	return &GetMenuHandler{
		menuQueryRepo: menuQueryRepo,
	}
}

// Handle 处理获取菜单查询
func (h *GetMenuHandler) Handle(ctx context.Context, query GetMenuQuery) (*menuApp.MenuResponse, error) {
	menuEntity, err := h.menuQueryRepo.FindByID(ctx, query.MenuID)
	if err != nil {
		return nil, fmt.Errorf("menu not found: %w", err)
	}

	return menuApp.ToMenuResponse(menuEntity), nil
}
