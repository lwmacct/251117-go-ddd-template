package setting

import (
	"context"
	"errors"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

// CreateHandler 创建设置命令处理器
type CreateHandler struct {
	settingCommandRepo setting.CommandRepository
	settingQueryRepo   setting.QueryRepository
}

// NewCreateHandler 创建 CreateHandler 实例
func NewCreateHandler(
	settingCommandRepo setting.CommandRepository,
	settingQueryRepo setting.QueryRepository,
) *CreateHandler {
	return &CreateHandler{
		settingCommandRepo: settingCommandRepo,
		settingQueryRepo:   settingQueryRepo,
	}
}

// Handle 处理创建设置命令
func (h *CreateHandler) Handle(ctx context.Context, cmd CreateCommand) (*CreateResultDTO, error) {
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

	return &CreateResultDTO{
		ID: settingEntity.ID,
	}, nil
}
