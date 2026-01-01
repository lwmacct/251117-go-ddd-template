package product

import (
	"context"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/product"
)

// UpdateHandler 更新产品处理器
type UpdateHandler struct {
	commandRepo product.CommandRepository
	queryRepo   product.QueryRepository
}

// NewUpdateHandler 创建 UpdateHandler 实例
func NewUpdateHandler(
	commandRepo product.CommandRepository,
	queryRepo product.QueryRepository,
) *UpdateHandler {
	return &UpdateHandler{
		commandRepo: commandRepo,
		queryRepo:   queryRepo,
	}
}

// Handle 处理更新产品命令
func (h *UpdateHandler) Handle(ctx context.Context, cmd UpdateProductCommand) (*ProductDTO, error) {
	// 获取现有产品
	p, err := h.queryRepo.GetByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}

	// 如果要更新代码，检查是否与其他产品冲突
	if cmd.Code != nil && *cmd.Code != p.Code {
		exists, err := h.queryRepo.ExistsByCode(ctx, *cmd.Code)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, product.ErrProductCodeExists
		}
		p.Code = *cmd.Code
	}

	// 如果要更新名称，检查是否与其他产品冲突
	if cmd.Name != nil && *cmd.Name != p.Name {
		exists, err := h.queryRepo.ExistsByName(ctx, *cmd.Name)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, product.ErrProductNameExists
		}
		p.Name = *cmd.Name
	}

	// 更新其他字段
	if cmd.Type != nil {
		p.Type = product.Type(*cmd.Type)
	}
	if cmd.Description != nil {
		p.Description = *cmd.Description
	}
	if cmd.Price != nil {
		p.Price = *cmd.Price
	}
	if cmd.Status != nil {
		p.Status = product.Status(*cmd.Status)
	}
	if cmd.LayoutRef != nil {
		p.LayoutRef = *cmd.LayoutRef
	}
	if cmd.MaxSeats != nil {
		p.MaxSeats = *cmd.MaxSeats
	}
	if cmd.TrialDays != nil {
		p.TrialDays = *cmd.TrialDays
	}

	// 持久化
	if err := h.commandRepo.Update(ctx, p); err != nil {
		return nil, err
	}

	return ToProductDTO(p), nil
}
