package setting

import (
	"encoding/json"
	"strconv"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

// uiConfigRaw 内部结构用于解析 JSONB
type uiConfigRaw struct {
	InputType  string              `json:"input_type"`
	Hint       string              `json:"hint"`
	Options    []SelectOptionDTO   `json:"options"`
	Validation any                 `json:"validation"`
	DependsOn  *DependsOnConfigDTO `json:"depends_on"`
}

// parseUIConfig 解析 UIConfig JSON 字符串
func parseUIConfig(jsonStr string) UIConfigDTO {
	if jsonStr == "" || jsonStr == "{}" {
		return UIConfigDTO{InputType: "text"}
	}

	var raw uiConfigRaw
	if err := json.Unmarshal([]byte(jsonStr), &raw); err != nil {
		return UIConfigDTO{InputType: "text"}
	}

	// 设置默认 input_type
	if raw.InputType == "" {
		raw.InputType = "text"
	}

	return UIConfigDTO(raw)
}

// ToSettingDTO 将领域实体转换为 DTO
func ToSettingDTO(s *setting.Setting) *SettingDTO {
	if s == nil {
		return nil
	}

	return &SettingDTO{
		ID:        s.ID,
		Key:       s.Key,
		Value:     s.Value,
		Category:  s.Category,
		Group:     s.Group,
		ValueType: s.ValueType,
		Label:     s.Label,
		UIConfig:  parseUIConfig(s.UIConfig),
		Order:     s.Order,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

// ToSettingDTOs 将领域实体列表转换为 DTO 列表
func ToSettingDTOs(settings []*setting.Setting) []SettingDTO {
	if len(settings) == 0 {
		return []SettingDTO{}
	}

	dtos := make([]SettingDTO, 0, len(settings))
	for _, s := range settings {
		if dto := ToSettingDTO(s); dto != nil {
			dtos = append(dtos, *dto)
		}
	}

	return dtos
}

// ToSchemaSettingDTO 将领域实体转换为 Schema DTO（精简版）
func ToSchemaSettingDTO(s *setting.Setting) *SchemaSettingDTO {
	if s == nil {
		return nil
	}

	return &SchemaSettingDTO{
		Key:       s.Key,
		Value:     s.Value,
		ValueType: s.ValueType,
		Label:     s.Label,
		UIConfig:  parseUIConfig(s.UIConfig),
		Order:     s.Order,
	}
}

// ToSchemaSettingDTOs 将领域实体列表转换为 Schema DTO 列表
func ToSchemaSettingDTOs(settings []*setting.Setting) []SchemaSettingDTO {
	if len(settings) == 0 {
		return []SchemaSettingDTO{}
	}

	dtos := make([]SchemaSettingDTO, 0, len(settings))
	for _, s := range settings {
		if dto := ToSchemaSettingDTO(s); dto != nil {
			dtos = append(dtos, *dto)
		}
	}

	return dtos
}

// parseValueForValidation 根据值类型解析字符串值为对应类型。
// 用于验证时将字符串值转换为正确的类型进行比较。
func parseValueForValidation(value string, valueType string) any {
	switch valueType {
	case "boolean":
		return value == "true"
	case "number":
		if f, err := strconv.ParseFloat(value, 64); err == nil {
			return f
		}
		return 0.0
	case "json":
		var v any
		if err := json.Unmarshal([]byte(value), &v); err == nil {
			return v
		}
		return nil
	default:
		return value
	}
}
