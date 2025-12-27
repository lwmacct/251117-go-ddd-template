package setting

import (
	"context"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

// BatchUpdateHandler 批量更新设置命令处理器
type BatchUpdateHandler struct {
	settingCommandRepo setting.CommandRepository
	settingQueryRepo   setting.QueryRepository
	validator          setting.Validator
}

// NewBatchUpdateHandler 创建 BatchUpdateHandler 实例
func NewBatchUpdateHandler(
	settingCommandRepo setting.CommandRepository,
	settingQueryRepo setting.QueryRepository,
	validator setting.Validator,
) *BatchUpdateHandler {
	return &BatchUpdateHandler{
		settingCommandRepo: settingCommandRepo,
		settingQueryRepo:   settingQueryRepo,
		validator:          validator,
	}
}

// Handle 处理批量更新设置命令
func (h *BatchUpdateHandler) Handle(ctx context.Context, cmd BatchUpdateCommand) error {
	if len(cmd.Settings) == 0 {
		return nil
	}

	// 1. 提取所有 keys 并构建 key -> value 映射
	keys := make([]string, len(cmd.Settings))
	keyValueMap := make(map[string]string, len(cmd.Settings))
	for i, item := range cmd.Settings {
		keys[i] = item.Key
		keyValueMap[item.Key] = item.Value
	}

	// 2. 批量查询所有现有配置（1 次查询，解决 N+1 问题）
	existingSettings, err := h.settingQueryRepo.FindByKeys(ctx, keys)
	if err != nil {
		return fmt.Errorf("failed to find settings: %w", err)
	}

	// 3. 构建 key -> setting 映射，用于验证存在性
	existingMap := make(map[string]*setting.Setting, len(existingSettings))
	for _, s := range existingSettings {
		existingMap[s.Key] = s
	}

	// 4. 验证所有 key 存在并更新值
	settings := make([]*setting.Setting, 0, len(cmd.Settings))
	for _, key := range keys {
		existing, ok := existingMap[key]
		if !ok {
			return fmt.Errorf("setting key %s does not exist", key)
		}
		existing.Value = keyValueMap[key]
		settings = append(settings, existing)
	}

	// 5. 执行验证（如果有验证器）
	if err := h.validateSettings(ctx, settings); err != nil {
		return err
	}

	// 6. 批量更新（1 次写入）
	if err := h.settingCommandRepo.BatchUpsert(ctx, settings); err != nil {
		return fmt.Errorf("failed to batch update settings: %w", err)
	}

	return nil
}

// validateSettings 执行设置验证
func (h *BatchUpdateHandler) validateSettings(ctx context.Context, settings []*setting.Setting) error {
	if h.validator == nil {
		return nil
	}

	// 获取所有设置用于跨字段验证
	allSettings, _ := h.getAllSettingsMap(ctx, settings)

	// 构建验证上下文
	validationItems := make([]*setting.ValidationContext, 0)
	for _, s := range settings {
		rule := extractValidationFromUIConfig(s.UIConfig)
		if rule == "" {
			continue
		}
		value := parseValueForValidation(s.Value, s.ValueType)
		validationItems = append(validationItems, &setting.ValidationContext{
			Key:         s.Key,
			Value:       value,
			Rule:        rule,
			AllSettings: allSettings,
		})
	}

	if len(validationItems) == 0 {
		return nil
	}

	errors, validateErr := h.validator.ValidateBatch(ctx, validationItems)
	if validateErr != nil {
		return fmt.Errorf("validation error: %w", validateErr)
	}
	if len(errors) > 0 {
		// 返回第一个错误（或可以返回所有错误）
		for key, msg := range errors {
			return fmt.Errorf("%w: %s - %s", setting.ErrValidationFailed, key, msg)
		}
	}

	return nil
}

// getAllSettingsMap 获取所有设置的 key -> value 映射，合并待更新的值
func (h *BatchUpdateHandler) getAllSettingsMap(ctx context.Context, pendingUpdates []*setting.Setting) (map[string]any, error) {
	settings, err := h.settingQueryRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	result := make(map[string]any, len(settings))
	for _, s := range settings {
		result[s.Key] = parseValueForValidation(s.Value, s.ValueType)
	}

	// 合并待更新的值（使用新值进行验证）
	for _, s := range pendingUpdates {
		result[s.Key] = parseValueForValidation(s.Value, s.ValueType)
	}

	return result, nil
}
