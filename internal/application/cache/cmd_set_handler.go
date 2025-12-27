package cache

import (
	"context"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/cache"
)

// SetHandler 设置缓存命令处理器
type SetHandler struct {
	cacheCommandRepo cache.CommandRepository
}

// NewSetHandler 创建设置缓存命令处理器
func NewSetHandler(cacheCommandRepo cache.CommandRepository) *SetHandler {
	return &SetHandler{
		cacheCommandRepo: cacheCommandRepo,
	}
}

// Handle 处理设置缓存命令
func (h *SetHandler) Handle(ctx context.Context, cmd SetCommand) error {
	return h.cacheCommandRepo.Set(ctx, cmd.Key, cmd.Value, cmd.TTL)
}
