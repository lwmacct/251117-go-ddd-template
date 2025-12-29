package menu

import (
	"context"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/menu"
)

// GetHandler 获取菜单查询处理器
type GetHandler struct {
	menuQueryRepo menu.QueryRepository
}

// NewGetHandler 创建 GetHandler 实例
func NewGetHandler(menuQueryRepo menu.QueryRepository) *GetHandler {
	return &GetHandler{
		menuQueryRepo: menuQueryRepo,
	}
}

// Handle 处理获取菜单查询
func (h *GetHandler) Handle(ctx context.Context, query GetQuery) (*MenuDTO, error) {
	menuEntity, err := h.menuQueryRepo.FindByID(ctx, query.MenuID)
	if err != nil {
		return nil, fmt.Errorf("menu not found: %w", err)
	}

	return ToMenuDTO(menuEntity), nil
}
