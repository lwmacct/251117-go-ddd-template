package setting

import (
	"context"
	"errors"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

// DeleteHandler 删除设置命令处理器
type DeleteHandler struct {
	settingCommandRepo setting.CommandRepository
	settingQueryRepo   setting.QueryRepository
}

// NewDeleteHandler 创建 DeleteHandler 实例
func NewDeleteHandler(
	settingCommandRepo setting.CommandRepository,
	settingQueryRepo setting.QueryRepository,
) *DeleteHandler {
	return &DeleteHandler{
		settingCommandRepo: settingCommandRepo,
		settingQueryRepo:   settingQueryRepo,
	}
}

// Handle 处理删除设置命令
func (h *DeleteHandler) Handle(ctx context.Context, cmd DeleteCommand) error {
	// 1. 查询设置
	settingEntity, err := h.settingQueryRepo.FindByKey(ctx, cmd.Key)
	if err != nil {
		return fmt.Errorf("failed to find setting: %w", err)
	}
	if settingEntity == nil {
		return errors.New("setting not found")
	}

	// 2. 删除设置
	if err := h.settingCommandRepo.Delete(ctx, settingEntity.ID); err != nil {
		return fmt.Errorf("failed to delete setting: %w", err)
	}

	return nil
}
