package cache

import (
	"context"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/cache"
)

// GetHandler 获取缓存查询处理器
type GetHandler struct {
	cacheQueryRepo cache.QueryRepository
}

// NewGetHandler 创建获取缓存查询处理器
func NewGetHandler(cacheQueryRepo cache.QueryRepository) *GetHandler {
	return &GetHandler{
		cacheQueryRepo: cacheQueryRepo,
	}
}

// Handle 处理获取缓存查询
func (h *GetHandler) Handle(ctx context.Context, query GetQuery) (*GetResultDTO, error) {
	var value any
	if err := h.cacheQueryRepo.Get(ctx, query.Key, &value); err != nil {
		return nil, err
	}

	return &GetResultDTO{
		Key:   query.Key,
		Value: value,
	}, nil
}
