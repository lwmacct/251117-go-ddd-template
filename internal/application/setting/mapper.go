package setting

import (
	"encoding/json"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

// ==================== Category Mappers ====================

// ToCategoryDTO 将 SettingCategory 实体转换为 CategoryDTO
func ToCategoryDTO(c *setting.SettingCategory) *CategoryDTO {
	if c == nil {
		return nil
	}

	return &CategoryDTO{
		ID:        c.ID,
		Key:       c.Key,
		Label:     c.Label,
		Icon:      c.Icon,
		Order:     c.Order,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

// ToCategoryDTOs 将 SettingCategory 实体列表转换为 CategoryDTO 列表
func ToCategoryDTOs(categories []*setting.SettingCategory) []CategoryDTO {
	if len(categories) == 0 {
		return []CategoryDTO{}
	}

	dtos := make([]CategoryDTO, 0, len(categories))
	for _, c := range categories {
		if dto := ToCategoryDTO(c); dto != nil {
			dtos = append(dtos, *dto)
		}
	}

	return dtos
}

// ==================== Setting Mappers ====================

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

// ToSettingDTO 将 Setting 实体转换为 SettingDTO
func ToSettingDTO(s *setting.Setting) *SettingDTO {
	if s == nil {
		return nil
	}

	return &SettingDTO{
		ID:             s.ID,
		Key:            s.Key,
		DefaultValue:   s.DefaultValue,
		Scope:          s.Scope,
		CategoryID:     s.CategoryID,
		Group:          s.Group,
		ValueType:      s.ValueType,
		Label:          s.Label,
		UIConfig:       parseUIConfig(s.UIConfig),
		Order:          s.Order,
		ViewPermission: s.ViewPermission,
		EditPermission: s.EditPermission,
		CreatedAt:      s.CreatedAt,
		UpdatedAt:      s.UpdatedAt,
	}
}

// ToSettingDTOs 将 Setting 实体列表转换为 SettingDTO 列表
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

// ToSchemaSettingDTO 将 Setting 转换为 SchemaSettingDTO（Admin 场景，包含全部字段）
func ToSchemaSettingDTO(s *setting.Setting) *SchemaSettingDTO {
	if s == nil {
		return nil
	}

	return &SchemaSettingDTO{
		Key:            s.Key,
		Value:          s.DefaultValue,
		DefaultValue:   s.DefaultValue,
		IsCustomized:   false,
		Scope:          s.Scope, // Admin 返回
		ValueType:      s.ValueType,
		Label:          s.Label,
		UIConfig:       parseUIConfig(s.UIConfig),
		Order:          s.Order,
		ViewPermission: s.ViewPermission, // Admin 返回
		EditPermission: s.EditPermission, // Admin 返回
	}
}

// ==================== UserSetting Mappers ====================

// ToUserSettingDTO 将 Setting 定义和可选的 UserSetting 合并为 UserSettingDTO
func ToUserSettingDTO(s *setting.Setting, us *setting.UserSetting) *UserSettingDTO {
	if s == nil {
		return nil
	}

	dto := &UserSettingDTO{
		Key:          s.Key,
		Value:        s.DefaultValue, // 默认使用系统默认值
		DefaultValue: s.DefaultValue,
		IsCustomized: false,
		CategoryID:   s.CategoryID,
		Group:        s.Group,
		ValueType:    s.ValueType,
		Label:        s.Label,
		UIConfig:     parseUIConfig(s.UIConfig),
		Order:        s.Order,
	}

	// 如果有用户自定义值，使用用户值
	if us != nil {
		dto.Value = us.Value
		dto.IsCustomized = true
	}

	return dto
}

// ToUserSchemaSettingDTO 将 Setting 定义和可选的 UserSetting 合并为 SchemaSettingDTO
// User 场景省略权限字段（Scope, ViewPermission, EditPermission 为空，omitempty 不输出）
func ToUserSchemaSettingDTO(s *setting.Setting, us *setting.UserSetting) *SchemaSettingDTO {
	if s == nil {
		return nil
	}

	dto := &SchemaSettingDTO{
		Key:          s.Key,
		Value:        s.DefaultValue, // 默认使用系统默认值
		DefaultValue: s.DefaultValue,
		IsCustomized: false,
		// Scope, ViewPermission, EditPermission 留空（omitempty 不输出）
		ValueType: s.ValueType,
		Label:     s.Label,
		UIConfig:  parseUIConfig(s.UIConfig),
		Order:     s.Order,
	}

	// 如果有用户自定义值，使用用户值
	if us != nil {
		dto.Value = us.Value
		dto.IsCustomized = true
	}

	return dto
}

// extractValidationRule 从 UIConfig 中提取验证规则
func extractValidationRule(uiConfig string) string {
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

	data, err := json.Marshal(raw.Validation)
	if err != nil {
		return ""
	}
	return string(data)
}
