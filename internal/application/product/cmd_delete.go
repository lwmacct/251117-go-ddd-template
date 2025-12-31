package product

import (
	"context"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/product"
)

// DeleteHandler 删除产品处理器
type DeleteHandler struct {
	commandRepo product.CommandRepository
}

// NewDeleteHandler 创建 DeleteHandler 实例
func NewDeleteHandler(commandRepo product.CommandRepository) *DeleteHandler {
	return &DeleteHandler{
		commandRepo: commandRepo,
	}
}

// Handle 处理删除产品命令
func (h *DeleteHandler) Handle(ctx context.Context, cmd DeleteProductCommand) error {
	return h.commandRepo.Delete(ctx, cmd.ID)
}
