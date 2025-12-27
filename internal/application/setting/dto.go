package setting

import "time"

// UIConfigDTO UI 配置（从 JSONB 解析）
type UIConfigDTO struct {
	InputType  string              `json:"input_type"`           // 控件类型: text, number, switch, select...
	Hint       string              `json:"hint,omitempty"`       // 输入提示
	Options    []SelectOptionDTO   `json:"options,omitempty"`    // 下拉选项
	Validation any                 `json:"validation,omitempty"` // JSON Logic 验证规则
	DependsOn  *DependsOnConfigDTO `json:"depends_on,omitempty"` // 依赖关系
}

// SelectOptionDTO 下拉选项
type SelectOptionDTO struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

// DependsOnConfigDTO 依赖关系
type DependsOnConfigDTO struct {
	Key      string `json:"key"`
	Value    any    `json:"value,omitempty"`
	Operator string `json:"operator,omitempty"` // eq, ne, gt, lt（默认 eq）
}

// SettingDTO 设置响应 DTO（CRUD API 使用）
type SettingDTO struct {
	ID        uint        `json:"id"`
	Key       string      `json:"key"`
	Value     string      `json:"value"`
	Category  string      `json:"category"`
	Group     string      `json:"group"`
	ValueType string      `json:"value_type"`
	Label     string      `json:"label"`
	UIConfig  UIConfigDTO `json:"ui_config"`
	Order     int         `json:"order"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// SchemaSettingDTO Schema API 专用 DTO（精简版）
type SchemaSettingDTO struct {
	Key       string      `json:"key"`
	Value     string      `json:"value"`
	ValueType string      `json:"value_type"`
	Label     string      `json:"label"`
	UIConfig  UIConfigDTO `json:"ui_config"`
	Order     int         `json:"order"`
}

// SchemaGroupDTO Schema API 分组
type SchemaGroupDTO struct {
	Group    string             `json:"group"`
	Label    string             `json:"label"`
	Settings []SchemaSettingDTO `json:"settings"`
}

// SchemaCategoryDTO Schema API 分类
type SchemaCategoryDTO struct {
	Category string           `json:"category"`
	Label    string           `json:"label"`
	Icon     string           `json:"icon"`
	Groups   []SchemaGroupDTO `json:"groups"`
}

// CreateResultDTO 创建设置结果 DTO
type CreateResultDTO struct {
	ID uint `json:"id"`
}

// SettingGroupDTO 按分组聚合的设置列表（CRUD API 使用）
type SettingGroupDTO struct {
	Group    string       `json:"group"`
	Label    string       `json:"label"`
	Settings []SettingDTO `json:"settings"`
}

// CategorySettingsDTO 按 Category 聚合的响应（CRUD API 使用）
type CategorySettingsDTO struct {
	Category string            `json:"category"`
	Label    string            `json:"label"`
	Icon     string            `json:"icon"`
	Groups   []SettingGroupDTO `json:"groups"`
}
