package setting

import (
	"context"
	"errors"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

// CreateSettingHandler 创建设置命令处理器
type CreateSettingHandler struct {
	settingCommandRepo setting.CommandRepository
	settingQueryRepo   setting.QueryRepository
}

// NewCreateSettingHandler 创建 CreateSettingHandler 实例
func NewCreateSettingHandler(
	settingCommandRepo setting.CommandRepository,
	settingQueryRepo setting.QueryRepository,
) *CreateSettingHandler {
	return &CreateSettingHandler{
		settingCommandRepo: settingCommandRepo,
		settingQueryRepo:   settingQueryRepo,
	}
}

// Handle 处理创建设置命令
func (h *CreateSettingHandler) Handle(ctx context.Context, cmd CreateSettingCommand) (*CreateSettingResultDTO, error) {
	// 1. 验证 Key 是否已存在
	existing, err := h.settingQueryRepo.FindByKey(ctx, cmd.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing setting: %w", err)
	}
	if existing != nil {
		return nil, errors.New("setting key already exists")
	}

	// 2. 设置默认值类型
	valueType := cmd.ValueType
	if valueType == "" {
		valueType = setting.ValueTypeString
	}

	// 3. 创建设置实体
	settingEntity := &setting.Setting{
		Key:       cmd.Key,
		Value:     cmd.Value,
		Category:  cmd.Category,
		ValueType: valueType,
		Label:     cmd.Label,
	}

	// 4. 保存设置
	if err := h.settingCommandRepo.Create(ctx, settingEntity); err != nil {
		return nil, fmt.Errorf("failed to create setting: %w", err)
	}

	return &CreateSettingResultDTO{
		ID: settingEntity.ID,
	}, nil
}
