package setting

import (
	"context"
	"errors"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

// UpdateSettingHandler 更新设置命令处理器
type UpdateSettingHandler struct {
	settingCommandRepo setting.CommandRepository
	settingQueryRepo   setting.QueryRepository
}

// NewUpdateSettingHandler 创建 UpdateSettingHandler 实例
func NewUpdateSettingHandler(
	settingCommandRepo setting.CommandRepository,
	settingQueryRepo setting.QueryRepository,
) *UpdateSettingHandler {
	return &UpdateSettingHandler{
		settingCommandRepo: settingCommandRepo,
		settingQueryRepo:   settingQueryRepo,
	}
}

// Handle 处理更新设置命令
func (h *UpdateSettingHandler) Handle(ctx context.Context, cmd UpdateSettingCommand) (*SettingDTO, error) {
	// 1. 查询设置
	settingEntity, err := h.settingQueryRepo.FindByKey(ctx, cmd.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to find setting: %w", err)
	}
	if settingEntity == nil {
		return nil, errors.New("setting not found")
	}

	// 2. 更新字段
	settingEntity.Value = cmd.Value
	if cmd.ValueType != "" {
		settingEntity.ValueType = cmd.ValueType
	}
	if cmd.Label != "" {
		settingEntity.Label = cmd.Label
	}

	// 3. 保存更新
	if err := h.settingCommandRepo.Update(ctx, settingEntity); err != nil {
		return nil, fmt.Errorf("failed to update setting: %w", err)
	}

	return ToSettingDTO(settingEntity), nil
}
