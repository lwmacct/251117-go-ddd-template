package product

import (
	"context"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/product"
)

// GetHandler 获取产品处理器
type GetHandler struct {
	queryRepo product.QueryRepository
}

// NewGetHandler 创建 GetHandler 实例
func NewGetHandler(queryRepo product.QueryRepository) *GetHandler {
	return &GetHandler{
		queryRepo: queryRepo,
	}
}

// Handle 处理获取产品查询
func (h *GetHandler) Handle(ctx context.Context, query GetProductQuery) (*ProductDTO, error) {
	p, err := h.queryRepo.GetByID(ctx, query.ID)
	if err != nil {
		return nil, err
	}
	return ToProductDTO(p), nil
}
