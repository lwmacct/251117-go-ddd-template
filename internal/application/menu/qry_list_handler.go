package menu

import (
	"context"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/menu"
)

// ListHandler 获取菜单列表查询处理器
type ListHandler struct {
	menuQueryRepo menu.QueryRepository
}

// NewListHandler 创建 ListHandler 实例
func NewListHandler(menuQueryRepo menu.QueryRepository) *ListHandler {
	return &ListHandler{
		menuQueryRepo: menuQueryRepo,
	}
}

// Handle 处理获取菜单列表查询
func (h *ListHandler) Handle(ctx context.Context, query ListQuery) ([]*MenuDTO, error) {
	menus, err := h.menuQueryRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch menus: %w", err)
	}

	// 转换为 DTO
	menuResponses := make([]*MenuDTO, 0, len(menus))
	for _, m := range menus {
		menuResponses = append(menuResponses, ToMenuDTO(m))
	}

	return menuResponses, nil
}
