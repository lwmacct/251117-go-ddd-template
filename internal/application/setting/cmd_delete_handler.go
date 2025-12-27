package setting

import (
	"context"
	"errors"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

// DeleteHandler 删除配置命令处理器
type DeleteHandler struct {
	commandRepo setting.CommandRepository
	queryRepo   setting.QueryRepository
}

// NewDeleteHandler 创建 DeleteHandler 实例
func NewDeleteHandler(
	commandRepo setting.CommandRepository,
	queryRepo setting.QueryRepository,
) *DeleteHandler {
	return &DeleteHandler{
		commandRepo: commandRepo,
		queryRepo:   queryRepo,
	}
}

// Handle 处理删除配置命令
func (h *DeleteHandler) Handle(ctx context.Context, cmd DeleteCommand) error {
	// 1. 查询配置定义
	def, err := h.queryRepo.FindByKey(ctx, cmd.Key)
	if err != nil {
		return fmt.Errorf("failed to find setting: %w", err)
	}
	if def == nil {
		return errors.New("setting not found")
	}

	// 2. 删除配置定义
	if err := h.commandRepo.Delete(ctx, cmd.Key); err != nil {
		return fmt.Errorf("failed to delete setting: %w", err)
	}

	return nil
}
