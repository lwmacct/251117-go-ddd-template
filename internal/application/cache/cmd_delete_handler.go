package cache

import (
	"context"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/cache"
)

// DeleteHandler 删除缓存命令处理器
type DeleteHandler struct {
	cacheCommandRepo cache.CommandRepository
}

// NewDeleteHandler 创建删除缓存命令处理器
func NewDeleteHandler(cacheCommandRepo cache.CommandRepository) *DeleteHandler {
	return &DeleteHandler{
		cacheCommandRepo: cacheCommandRepo,
	}
}

// Handle 处理删除缓存命令
func (h *DeleteHandler) Handle(ctx context.Context, cmd DeleteCommand) error {
	return h.cacheCommandRepo.Delete(ctx, cmd.Key)
}
