package product

import (
	"context"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/product"
)

// ListHandler 产品列表处理器
type ListHandler struct {
	queryRepo product.QueryRepository
}

// NewListHandler 创建 ListHandler 实例
func NewListHandler(queryRepo product.QueryRepository) *ListHandler {
	return &ListHandler{
		queryRepo: queryRepo,
	}
}

// ListResult 列表查询结果
type ListResult struct {
	Items []*ProductDTO
	Total int64
}

// Handle 处理产品列表查询
func (h *ListHandler) Handle(ctx context.Context, query ListProductsQuery) (*ListResult, error) {
	// 获取总数
	total, err := h.queryRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	// 获取列表
	products, err := h.queryRepo.List(ctx, query.Offset, query.Limit)
	if err != nil {
		return nil, err
	}

	return &ListResult{
		Items: ToProductDTOs(products),
		Total: total,
	}, nil
}
