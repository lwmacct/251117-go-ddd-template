package setting

import "time"

// ==================== Category DTO ====================

// CategoryDTO 配置分类响应 DTO
type CategoryDTO struct {
	ID        uint      `json:"id"`
	Key       string    `json:"key"`
	Label     string    `json:"label"`
	Icon      string    `json:"icon"`
	Order     int       `json:"order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CategoryMetaDTO 分类元信息 DTO（不含 settings，用于懒加载场景）
type CategoryMetaDTO struct {
	Category string `json:"category"` // key
	Label    string `json:"label"`
	Icon     string `json:"icon"`
	Order    int    `json:"order"`
}

// CreateCategoryResultDTO 创建分类结果 DTO
type CreateCategoryResultDTO struct {
	ID uint `json:"id"`
}

// ==================== UI 配置 DTO ====================

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

// ==================== Setting DTO ====================

// SettingDTO 配置响应 DTO
type SettingDTO struct {
	ID             uint        `json:"id"`
	Key            string      `json:"key"`
	DefaultValue   any         `json:"default_value"` // JSONB 原生值
	Scope          string      `json:"scope"`         // system | user
	CategoryID     uint        `json:"category_id"`
	Group          string      `json:"group"`
	ValueType      string      `json:"value_type"`
	Label          string      `json:"label"`
	UIConfig       UIConfigDTO `json:"ui_config"`
	Order          int         `json:"order"`
	ViewPermission string      `json:"view_permission"` // 查看权限
	EditPermission string      `json:"edit_permission"` // 编辑权限
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}

// CreateResultDTO 创建配置结果 DTO
type CreateResultDTO struct {
	ID uint `json:"id"`
}

// ==================== Schema DTO ====================

// SchemaSettingDTO Schema API 专用 DTO（Admin/User 统一）
// Admin 场景返回全部字段，User 场景通过 omitempty 省略权限字段
type SchemaSettingDTO struct {
	Key            string      `json:"key"`
	Value          any         `json:"value"`           // 实际生效值
	DefaultValue   any         `json:"default_value"`   // 系统默认值
	IsCustomized   bool        `json:"is_customized"`   // 是否用户自定义
	Scope          string      `json:"scope,omitempty"` // 仅 Admin: system | user
	ValueType      string      `json:"value_type"`
	Label          string      `json:"label"`
	UIConfig       UIConfigDTO `json:"ui_config"`
	Order          int         `json:"order"`
	ViewPermission string      `json:"view_permission,omitempty"` // 仅 Admin
	EditPermission string      `json:"edit_permission,omitempty"` // 仅 Admin
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

// ==================== 分组聚合 DTO ====================

// SettingGroupDTO 按分组聚合的配置列表
type SettingGroupDTO struct {
	Group    string       `json:"group"`
	Label    string       `json:"label"`
	Settings []SettingDTO `json:"settings"`
}

// CategorySettingsDTO 按 Category 聚合的配置响应
type CategorySettingsDTO struct {
	Category string            `json:"category"`
	Label    string            `json:"label"`
	Icon     string            `json:"icon"`
	Groups   []SettingGroupDTO `json:"groups"`
}

// ==================== UserSetting DTO ====================

// UserSettingDTO 用户配置响应 DTO（合并视图）
type UserSettingDTO struct {
	Key          string      `json:"key"`
	Value        any         `json:"value"`         // 实际生效值（用户值或默认值）
	DefaultValue any         `json:"default_value"` // 系统默认值
	IsCustomized bool        `json:"is_customized"` // 是否为用户自定义
	CategoryID   uint        `json:"category_id"`
	Group        string      `json:"group"`
	ValueType    string      `json:"value_type"`
	Label        string      `json:"label"`
	UIConfig     UIConfigDTO `json:"ui_config"`
	Order        int         `json:"order"`
}

// UserSettingGroupDTO 按分组聚合的用户配置列表
type UserSettingGroupDTO struct {
	Group    string           `json:"group"`
	Label    string           `json:"label"`
	Settings []UserSettingDTO `json:"settings"`
}

// UserCategorySettingsDTO 按 Category 聚合的用户配置响应
type UserCategorySettingsDTO struct {
	Category string                `json:"category"`
	Label    string                `json:"label"`
	Icon     string                `json:"icon"`
	Groups   []UserSettingGroupDTO `json:"groups"`
}
