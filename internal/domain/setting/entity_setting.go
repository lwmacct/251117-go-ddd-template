package setting

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

// Setting 系统配置实体。
// 采用 Key-Value 模式存储配置项，支持分类、类型标注和 UI 元数据。
//
// UIConfig 字段（JSONB）统一存储前端渲染所需的 UI 配置：
//   - input_type: 控件类型（text, switch, select 等）
//   - hint: 输入提示文字
//   - options: 下拉/单选/多选的选项列表
//   - validation: JSON Logic 验证规则
//   - depends_on: 依赖关系配置
type Setting struct {
	ID        uint   `json:"id"`
	Key       string `json:"key"`
	Value     string `json:"value"`
	Category  string `json:"category"`
	Group     string `json:"group"` // Category 内子分组
	ValueType string `json:"value_type"`
	Label     string `json:"label"`
	UIConfig  string `json:"ui_config"` // JSONB: 统一的 UI 配置
	Order     int    `json:"order"`     // 排序权重（小的在前）

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// IsValidValueType 检查 ValueType 是否有效
func (s *Setting) IsValidValueType() bool {
	switch s.ValueType {
	case ValueTypeString, ValueTypeNumber, ValueTypeBoolean, ValueTypeJSON:
		return true
	default:
		return false
	}
}

// IsValidCategory 检查 Category 是否有效
func (s *Setting) IsValidCategory() bool {
	switch s.Category {
	case CategoryGeneral, CategorySecurity, CategoryNotification, CategoryBackup:
		return true
	default:
		return false
	}
}

// ParseBool 将 Value 解析为布尔值
func (s *Setting) ParseBool() (bool, error) {
	if s.ValueType != ValueTypeBoolean {
		return false, ErrValueTypeMismatch
	}
	lower := strings.ToLower(s.Value)
	switch lower {
	case "true", "1", "yes", "on":
		return true, nil
	case "false", "0", "no", "off":
		return false, nil
	default:
		return false, ErrInvalidBoolValue
	}
}

// ParseInt 将 Value 解析为整数
func (s *Setting) ParseInt() (int, error) {
	if s.ValueType != ValueTypeNumber {
		return 0, ErrValueTypeMismatch
	}
	val, err := strconv.Atoi(s.Value)
	if err != nil {
		return 0, ErrInvalidNumberValue
	}
	return val, nil
}

// ParseFloat 将 Value 解析为浮点数
func (s *Setting) ParseFloat() (float64, error) {
	if s.ValueType != ValueTypeNumber {
		return 0, ErrValueTypeMismatch
	}
	val, err := strconv.ParseFloat(s.Value, 64)
	if err != nil {
		return 0, ErrInvalidNumberValue
	}
	return val, nil
}

// ParseJSON 将 Value 解析为 JSON 对象
func (s *Setting) ParseJSON(v any) error {
	if s.ValueType != ValueTypeJSON {
		return ErrValueTypeMismatch
	}
	if err := json.Unmarshal([]byte(s.Value), v); err != nil {
		return ErrInvalidJSONValue
	}
	return nil
}

// SetBool 设置布尔值
func (s *Setting) SetBool(val bool) {
	s.ValueType = ValueTypeBoolean
	if val {
		s.Value = "true"
	} else {
		s.Value = "false"
	}
}

// SetInt 设置整数值
func (s *Setting) SetInt(val int) {
	s.ValueType = ValueTypeNumber
	s.Value = strconv.Itoa(val)
}

// SetJSON 设置 JSON 值
func (s *Setting) SetJSON(v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	s.ValueType = ValueTypeJSON
	s.Value = string(data)
	return nil
}

// IsEmpty 检查值是否为空
func (s *Setting) IsEmpty() bool {
	return s.Value == ""
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
