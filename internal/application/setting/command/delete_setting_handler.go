package command

import (
	"context"
	"errors"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

// DeleteSettingHandler 删除设置命令处理器
type DeleteSettingHandler struct {
	settingCommandRepo setting.CommandRepository
	settingQueryRepo   setting.QueryRepository
}

// NewDeleteSettingHandler 创建 DeleteSettingHandler 实例
func NewDeleteSettingHandler(
	settingCommandRepo setting.CommandRepository,
	settingQueryRepo setting.QueryRepository,
) *DeleteSettingHandler {
	return &DeleteSettingHandler{
		settingCommandRepo: settingCommandRepo,
		settingQueryRepo:   settingQueryRepo,
	}
}

// Handle 处理删除设置命令
func (h *DeleteSettingHandler) Handle(ctx context.Context, cmd DeleteSettingCommand) error {
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
