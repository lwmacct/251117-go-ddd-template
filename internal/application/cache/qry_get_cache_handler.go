package cache

import (
	"context"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/cache"
)

// GetCacheHandler 获取缓存查询处理器
type GetCacheHandler struct {
	cacheQueryRepo cache.QueryRepository
}

// NewGetCacheHandler 创建获取缓存查询处理器
func NewGetCacheHandler(cacheQueryRepo cache.QueryRepository) *GetCacheHandler {
	return &GetCacheHandler{
		cacheQueryRepo: cacheQueryRepo,
	}
}

// Handle 处理获取缓存查询
func (h *GetCacheHandler) Handle(ctx context.Context, query GetCacheQuery) (*GetCacheResultDTO, error) {
	var value any
	if err := h.cacheQueryRepo.Get(ctx, query.Key, &value); err != nil {
		return nil, err
	}

	return &GetCacheResultDTO{
		Key:   query.Key,
		Value: value,
	}, nil
}
