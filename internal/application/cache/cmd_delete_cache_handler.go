package cache

import (
	"context"

	domainCache "github.com/lwmacct/251117-go-ddd-template/internal/domain/cache"
)

// DeleteCacheHandler 删除缓存命令处理器
type DeleteCacheHandler struct {
	cacheCommandRepo domainCache.CommandRepository
}

// NewDeleteCacheHandler 创建删除缓存命令处理器
func NewDeleteCacheHandler(cacheCommandRepo domainCache.CommandRepository) *DeleteCacheHandler {
	return &DeleteCacheHandler{
		cacheCommandRepo: cacheCommandRepo,
	}
}

// Handle 处理删除缓存命令
func (h *DeleteCacheHandler) Handle(ctx context.Context, cmd DeleteCacheCommand) error {
	return h.cacheCommandRepo.Delete(ctx, cmd.Key)
}
