package query

import (
	"context"
	"fmt"

	menuApp "github.com/lwmacct/251117-go-ddd-template/internal/application/menu"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/menu"
)

// ListMenusHandler 获取菜单列表查询处理器
type ListMenusHandler struct {
	menuQueryRepo menu.QueryRepository
}

// NewListMenusHandler 创建 ListMenusHandler 实例
func NewListMenusHandler(menuQueryRepo menu.QueryRepository) *ListMenusHandler {
	return &ListMenusHandler{
		menuQueryRepo: menuQueryRepo,
	}
}

// Handle 处理获取菜单列表查询
func (h *ListMenusHandler) Handle(ctx context.Context, query ListMenusQuery) ([]*menuApp.MenuResponse, error) {
	menus, err := h.menuQueryRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch menus: %w", err)
	}

	// 转换为 DTO
	menuResponses := make([]*menuApp.MenuResponse, 0, len(menus))
	for _, m := range menus {
		menuResponses = append(menuResponses, menuApp.ToMenuResponse(m))
	}

	return menuResponses, nil
}
