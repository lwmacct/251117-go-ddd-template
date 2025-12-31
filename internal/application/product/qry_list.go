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
	products, err := h.queryProducts(ctx, query.Status, query.Offset, query.Limit)
	if err != nil {
		return nil, err
	}

	total, err := h.countProducts(ctx, query.Status)
	if err != nil {
		return nil, err
	}

	return &ListResult{
		Items: ToProductDTOs(products),
		Total: total,
	}, nil
}

// queryProducts 查询产品列表（按状态筛选或全部）
func (h *ListHandler) queryProducts(ctx context.Context, status string, offset, limit int) ([]*product.Product, error) {
	if status != "" {
		return h.queryRepo.ListByStatus(ctx, status, offset, limit)
	}
	return h.queryRepo.List(ctx, offset, limit)
}

// countProducts 统计产品数量（按状态筛选或全部）
func (h *ListHandler) countProducts(ctx context.Context, status string) (int64, error) {
	if status != "" {
		return h.queryRepo.CountByStatus(ctx, status)
	}
	return h.queryRepo.Count(ctx)
}
