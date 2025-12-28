package setting

import (
	"context"
	"errors"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

// UpdateCategoryHandler 更新配置分类命令处理器
type UpdateCategoryHandler struct {
	commandRepo setting.SettingCategoryCommandRepository
	queryRepo   setting.SettingCategoryQueryRepository
}

// NewUpdateCategoryHandler 创建 UpdateCategoryHandler 实例
func NewUpdateCategoryHandler(
	commandRepo setting.SettingCategoryCommandRepository,
	queryRepo setting.SettingCategoryQueryRepository,
) *UpdateCategoryHandler {
	return &UpdateCategoryHandler{
		commandRepo: commandRepo,
		queryRepo:   queryRepo,
	}
}

// Handle 处理更新配置分类命令
func (h *UpdateCategoryHandler) Handle(ctx context.Context, cmd UpdateCategoryCommand) (*CategoryDTO, error) {
	// 1. 查询现有分类
	category, err := h.queryRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find category: %w", err)
	}
	if category == nil {
		return nil, errors.New("category not found")
	}

	// 2. 更新字段
	if cmd.Label != "" {
		category.UpdateLabel(cmd.Label)
	}
	if cmd.Icon != "" {
		category.UpdateIcon(cmd.Icon)
	}
	category.UpdateOrder(cmd.Order)

	// 3. 验证实体
	if validateErr := category.Validate(); validateErr != nil {
		return nil, fmt.Errorf("invalid category: %w", validateErr)
	}

	// 4. 保存更新
	if updateErr := h.commandRepo.Update(ctx, category); updateErr != nil {
		return nil, fmt.Errorf("failed to update category: %w", updateErr)
	}

	// 5. 重新查询返回最新数据
	updated, err := h.queryRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated category: %w", err)
	}

	return ToCategoryDTO(updated), nil
}
