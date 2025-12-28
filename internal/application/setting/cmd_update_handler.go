package setting

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

// UpdateHandler 更新配置命令处理器
type UpdateHandler struct {
	commandRepo setting.CommandRepository
	queryRepo   setting.QueryRepository
	validator   setting.Validator
	schemaCache SchemaCacheService
}

// NewUpdateHandler 创建 UpdateHandler 实例
func NewUpdateHandler(
	commandRepo setting.CommandRepository,
	queryRepo setting.QueryRepository,
	validator setting.Validator,
	schemaCache SchemaCacheService,
) *UpdateHandler {
	return &UpdateHandler{
		commandRepo: commandRepo,
		queryRepo:   queryRepo,
		validator:   validator,
		schemaCache: schemaCache,
	}
}

// Handle 处理更新配置命令
func (h *UpdateHandler) Handle(ctx context.Context, cmd UpdateCommand) (*SettingDTO, error) {
	// 1. 查询配置定义
	def, err := h.queryRepo.FindByKey(ctx, cmd.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to find setting: %w", err)
	}
	if def == nil {
		return nil, errors.New("setting not found")
	}

	// 2. 执行验证（如果有验证器和验证规则）
	validationRule := extractValidationFromUIConfig(def.UIConfig)
	if h.validator != nil && validationRule != "" {
		// 获取所有设置用于跨字段验证
		allSettings, _ := h.getAllSettingsMap(ctx)

		vctx := &setting.ValidationContext{
			Key:         cmd.Key,
			Value:       cmd.DefaultValue, // 直接使用 any 类型，无需转换
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
	def.DefaultValue = cmd.DefaultValue
	if cmd.Label != "" {
		def.Label = cmd.Label
	}
	if cmd.UIConfig != "" {
		def.UIConfig = cmd.UIConfig
	}
	if cmd.Order != 0 {
		def.Order = cmd.Order
	}

	// 4. 保存更新
	if err := h.commandRepo.Update(ctx, def); err != nil {
		return nil, fmt.Errorf("failed to update setting: %w", err)
	}

	// 5. 失效 Schema 缓存
	if h.schemaCache != nil {
		if err := h.schemaCache.DeleteAdminSchemaAll(ctx); err != nil {
			slog.Warn("admin schema cache invalidation failed", "key", cmd.Key, "err", err)
		}
	}

	return ToSettingDTO(def), nil
}

// getAllSettingsMap 获取所有配置的 key -> value 映射
func (h *UpdateHandler) getAllSettingsMap(ctx context.Context) (map[string]any, error) {
	defs, err := h.queryRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	result := make(map[string]any, len(defs))
	for _, d := range defs {
		result[d.Key] = d.DefaultValue // 直接使用 any 类型
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
