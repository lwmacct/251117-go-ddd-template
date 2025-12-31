package product

import (
	"context"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/product"
)

// CreateHandler 创建产品处理器
type CreateHandler struct {
	commandRepo product.CommandRepository
	queryRepo   product.QueryRepository
}

// NewCreateHandler 创建 CreateHandler 实例
func NewCreateHandler(
	commandRepo product.CommandRepository,
	queryRepo product.QueryRepository,
) *CreateHandler {
	return &CreateHandler{
		commandRepo: commandRepo,
		queryRepo:   queryRepo,
	}
}

// Handle 处理创建产品命令
func (h *CreateHandler) Handle(ctx context.Context, cmd CreateProductCommand) (*ProductDTO, error) {
	// 检查产品名称是否已存在
	exists, err := h.queryRepo.ExistsByName(ctx, cmd.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, product.ErrProductNameExists
	}

	// 设置默认状态
	status := product.Status(cmd.Status)
	if status == "" {
		status = product.StatusActive
	}

	// 创建产品实体
	p := &product.Product{
		Name:        cmd.Name,
		Description: cmd.Description,
		Price:       cmd.Price,
		Status:      status,
	}

	// 持久化
	if err := h.commandRepo.Create(ctx, p); err != nil {
		return nil, err
	}

	return ToProductDTO(p), nil
}
