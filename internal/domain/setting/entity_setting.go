package setting

import "time"

// Setting 配置定义实体。
// 存储配置项的 Schema 和默认值，支持分类、类型标注和 UI 元数据。
//
// DefaultValue 字段（JSONB）直接存储原生 JSON 值：
//   - 字符串: "My Site"
//   - 数值: 30
//   - 布尔值: true
//   - JSON 对象/数组: {"key": "value"} 或 [1, 2, 3]
//
// UIConfig 字段（JSONB）统一存储前端渲染所需的 UI 配置：
//   - input_type: 控件类型（text, switch, select 等）
//   - hint: 输入提示文字
//   - options: 下拉/单选/多选的选项列表
//   - validation: JSON Logic 验证规则
//   - depends_on: 依赖关系配置
type Setting struct {
	ID           uint   // 唯一标识
	Key          string // 配置键，唯一约束
	DefaultValue any    // 默认值（JSONB 原生值）
	Category     string // 分类：general, security, notification, backup
	Group        string // 分类内子分组：basic, locale, appearance 等
	ValueType    string // 值类型：string, number, boolean, json（用于 UI 提示）
	Label        string // 显示标签
	UIConfig     string // UI 配置（JSONB 字符串）
	Order        int    // 排序权重（小的在前）

	CreatedAt time.Time
	UpdatedAt time.Time
}

// IsValidValueType 检查 ValueType 是否有效。
func (s *Setting) IsValidValueType() bool {
	switch s.ValueType {
	case ValueTypeString, ValueTypeNumber, ValueTypeBoolean, ValueTypeJSON:
		return true
	default:
		return false
	}
}

// IsValidCategory 检查 Category 是否有效。
func (s *Setting) IsValidCategory() bool {
	switch s.Category {
	case CategoryGeneral, CategorySecurity, CategoryNotification, CategoryBackup:
		return true
	default:
		return false
	}
}

// 配置分类常量。
// 用于对配置项进行逻辑分组，便于管理界面展示和权限控制。
const (
	CategoryGeneral      = "general"      // 通用配置（站点名称、Logo 等）
	CategorySecurity     = "security"     // 安全配置（密码策略、登录限制等）
	CategoryNotification = "notification" // 通知配置（邮件、短信等）
	CategoryBackup       = "backup"       // 备份配置（备份周期、保留策略等）
)

// 值类型常量。
// 指示配置值的数据类型，前端可据此渲染不同的输入控件。
const (
	ValueTypeString  = "string"  // 字符串类型，使用文本输入框
	ValueTypeNumber  = "number"  // 数值类型，使用数字输入框
	ValueTypeBoolean = "boolean" // 布尔类型，使用开关控件
	ValueTypeJSON    = "json"    // JSON 类型，使用 JSON 编辑器
)

// 控件类型常量。
// 用于前端动态渲染对应的输入控件。
const (
	InputTypeText     = "text"     // 单行文本框
	InputTypeTextarea = "textarea" // 多行文本框
	InputTypeNumber   = "number"   // 数字输入框
	InputTypeSwitch   = "switch"   // 开关
	InputTypeSelect   = "select"   // 下拉选择
	InputTypeRadio    = "radio"    // 单选按钮组
	InputTypeCheckbox = "checkbox" // 多选复选框
	InputTypePassword = "password" // 密码输入框
	InputTypeEmail    = "email"    // 邮箱输入框
	InputTypeURL      = "url"      // URL 输入框
	InputTypeJSON     = "json"     // JSON 编辑器
	InputTypeColor    = "color"    // 颜色选择器（预留）
)
