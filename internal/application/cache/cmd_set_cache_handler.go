package cache

import (
	"context"

	domainCache "github.com/lwmacct/251117-go-ddd-template/internal/domain/cache"
)

// SetCacheHandler 设置缓存命令处理器
type SetCacheHandler struct {
	cacheCommandRepo domainCache.CommandRepository
}

// NewSetCacheHandler 创建设置缓存命令处理器
func NewSetCacheHandler(cacheCommandRepo domainCache.CommandRepository) *SetCacheHandler {
	return &SetCacheHandler{
		cacheCommandRepo: cacheCommandRepo,
	}
}

// Handle 处理设置缓存命令
func (h *SetCacheHandler) Handle(ctx context.Context, cmd SetCacheCommand) error {
	return h.cacheCommandRepo.Set(ctx, cmd.Key, cmd.Value, cmd.TTL)
}
