package setting

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

// UpdateHandler 更新设置命令处理器
type UpdateHandler struct {
	settingCommandRepo setting.CommandRepository
	settingQueryRepo   setting.QueryRepository
	validator          setting.Validator
}

// NewUpdateHandler 创建 UpdateHandler 实例
func NewUpdateHandler(
	settingCommandRepo setting.CommandRepository,
	settingQueryRepo setting.QueryRepository,
	validator setting.Validator,
) *UpdateHandler {
	return &UpdateHandler{
		settingCommandRepo: settingCommandRepo,
		settingQueryRepo:   settingQueryRepo,
		validator:          validator,
	}
}

// Handle 处理更新设置命令
func (h *UpdateHandler) Handle(ctx context.Context, cmd UpdateCommand) (*SettingDTO, error) {
	// 1. 查询设置
	settingEntity, err := h.settingQueryRepo.FindByKey(ctx, cmd.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to find setting: %w", err)
	}
	if settingEntity == nil {
		return nil, errors.New("setting not found")
	}

	// 2. 执行验证（如果有验证器和验证规则）
	validationRule := extractValidationFromUIConfig(settingEntity.UIConfig)
	if h.validator != nil && validationRule != "" {
		// 获取所有设置用于跨字段验证
		allSettings, _ := h.getAllSettingsMap(ctx)

		// 转换值类型
		value := parseValueForValidation(cmd.Value, settingEntity.ValueType)

		vctx := &setting.ValidationContext{
			Key:         cmd.Key,
			Value:       value,
			Rule:        validationRule,
			AllSettings: allSettings,
		}

		result, validateErr := h.validator.Validate(ctx, vctx)
		if validateErr != nil {
			return nil, fmt.Errorf("validation error: %w", validateErr)
		}
		if !result.Valid {
			return nil, fmt.Errorf("%w: %s", setting.ErrValidationFailed, result.Message)
		}
	}

	// 3. 更新字段
	settingEntity.Value = cmd.Value
	if cmd.ValueType != "" {
		settingEntity.ValueType = cmd.ValueType
	}
	if cmd.Label != "" {
		settingEntity.Label = cmd.Label
	}

	// 4. 保存更新
	if err := h.settingCommandRepo.Update(ctx, settingEntity); err != nil {
		return nil, fmt.Errorf("failed to update setting: %w", err)
	}

	return ToSettingDTO(settingEntity), nil
}

// getAllSettingsMap 获取所有设置的 key -> value 映射
func (h *UpdateHandler) getAllSettingsMap(ctx context.Context) (map[string]any, error) {
	settings, err := h.settingQueryRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	result := make(map[string]any, len(settings))
	for _, s := range settings {
		result[s.Key] = parseValueForValidation(s.Value, s.ValueType)
	}
	return result, nil
}

// extractValidationFromUIConfig 从 UIConfig JSON 中提取验证规则
func extractValidationFromUIConfig(uiConfig string) string {
	if uiConfig == "" || uiConfig == "{}" {
		return ""
	}

	var raw struct {
		Validation any `json:"validation"`
	}
	if err := json.Unmarshal([]byte(uiConfig), &raw); err != nil {
		return ""
	}

	if raw.Validation == nil {
		return ""
	}

	// 将验证规则转回 JSON 字符串
	data, err := json.Marshal(raw.Validation)
	if err != nil {
		return ""
	}
	return string(data)
}
