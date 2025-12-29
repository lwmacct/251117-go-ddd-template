package setting

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Setting 配置定义实体。
// 存储配置项的 Schema 和默认值，支持分类、类型标注和 UI 元数据。
//
// DefaultValue 字段（JSONB）直接存储原生 JSON 值：
//   - 字符串: "My Site"
//   - 数值: 30
//   - 布尔值: true
//   - JSON 对象/数组: {"key": "value"} 或 [1, 2, 3]
//
// Scope 字段决定配置的作用域：
//   - "system": 系统设置，全局唯一，管理员直接修改 DefaultValue
//   - "user": 用户设置，DefaultValue 作为初始值，用户可在 user_settings 表覆盖
//
// 权限字段用于 RBAC 控制：
//   - ViewPermission: 查看权限，如 "user:settings:read"
//   - EditPermission: 编辑权限，如 "admin:settings:update"
//
// InputType 决定前端控件类型和后端自动校验规则（email/url/password 等）。
// Validation 存储自定义 JSON Logic 规则，用于业务级增强校验。
// UIConfig 存储前端展示配置：hint（提示）、options（下拉选项）、depends_on（依赖关系）。
type Setting struct {
	ID           uint   // 唯一标识
	Key          string // 配置键，唯一约束
	DefaultValue any    // 默认值（JSONB 原生值）
	Scope        string // 作用域：system（全局唯一）| user（可覆盖）
	CategoryID   uint   // 外键关联 SettingCategory.ID
	Group        string // 分类内子分组：basic, locale, appearance 等
	ValueType    string // 值类型：string, number, boolean, json（用于类型校验）
	Label        string // 显示标签
	Order        int    // 排序权重（小的在前）

	// UI 配置
	InputType  string // 控件类型：text, email, url, password, select 等（决定自动校验规则）
	Validation string // 自定义校验规则（JSON Logic 格式）
	UIConfig   string // 前端展示配置：hint、options、depends_on（JSONB 字符串）

	// 权限控制（复用 RBAC）
	ViewPermission string // 查看权限
	EditPermission string // 编辑权限

	CreatedAt time.Time
	UpdatedAt time.Time
}

// =============================================================================
// 验证方法
// =============================================================================

// Validate 验证实体完整性。
//
// 检查：
//   - Key 非空且格式有效
//   - CategoryID 非零（由数据库外键保证引用完整性）
//   - ValueType 有效
//   - InputType 有效
//   - Scope 有效
//   - ViewPermission 和 EditPermission 非空
//   - DefaultValue 与 ValueType 匹配
//   - DefaultValue 通过 InputType 格式校验
func (s *Setting) Validate() error {
	if s.Key == "" {
		return ErrInvalidValue
	}
	if _, err := NewSettingKey(s.Key); err != nil {
		return err
	}
	if s.CategoryID == 0 {
		return ErrCategoryNotFound
	}
	if !s.IsValidValueType() {
		return ErrInvalidValueType
	}
	if !s.IsValidInputType() {
		return ErrInvalidInputType
	}
	if !s.IsValidScope() {
		return ErrInvalidScope
	}
	if s.ViewPermission == "" || s.EditPermission == "" {
		return ErrInvalidPermission
	}
	if err := s.ValidateValue(s.DefaultValue); err != nil {
		return err
	}
	if err := s.ValidateByInputType(s.DefaultValue); err != nil {
		return err
	}
	return nil
}

// ValidateValue 验证给定值是否符合配置定义的类型。
func (s *Setting) ValidateValue(value any) error {
	if value == nil {
		return nil // nil 值总是允许的
	}

	switch s.ValueType {
	case ValueTypeString:
		if _, ok := value.(string); !ok {
			return ErrInvalidValueType
		}
	case ValueTypeNumber:
		switch value.(type) {
		case int, int8, int16, int32, int64,
			uint, uint8, uint16, uint32, uint64,
			float32, float64, json.Number:
			// 有效的数值类型
		default:
			return ErrInvalidValueType
		}
	case ValueTypeBoolean:
		if _, ok := value.(bool); !ok {
			return ErrInvalidValueType
		}
	case ValueTypeJSON:
		// JSON 类型接受 map 或 slice
		switch value.(type) {
		case map[string]any, []any:
			// 有效的 JSON 复合类型
		default:
			return ErrInvalidValueType
		}
	default:
		return ErrInvalidValueType
	}
	return nil
}

// IsValidValueType 报告 ValueType 是否有效。
func (s *Setting) IsValidValueType() bool {
	switch s.ValueType {
	case ValueTypeString, ValueTypeNumber, ValueTypeBoolean, ValueTypeJSON:
		return true
	default:
		return false
	}
}

// IsValidScope 报告 Scope 是否有效。
func (s *Setting) IsValidScope() bool {
	switch s.Scope {
	case ScopeSystem, ScopeUser:
		return true
	default:
		return false
	}
}

// =============================================================================
// Scope 方法
// =============================================================================

// IsSystemScope 报告是否为系统级配置。
//
// 系统级配置全局唯一，管理员直接修改 DefaultValue。
func (s *Setting) IsSystemScope() bool {
	return s.Scope == ScopeSystem
}

// IsUserScope 报告是否为用户级配置。
//
// 用户级配置允许用户在 user_settings 表中覆盖。
func (s *Setting) IsUserScope() bool {
	return s.Scope == ScopeUser
}

// =============================================================================
// 权限方法
// =============================================================================

// CanView 报告给定权限列表是否可以查看此配置。
//
// 权限格式遵循 RBAC 三段式：domain:resource:action
// 支持通配符匹配：
//   - "*:settings:read" 匹配所有 domain
//   - "admin:*:*" 匹配 admin 的所有操作
func (s *Setting) CanView(permissions []string) bool {
	return matchPermission(permissions, s.ViewPermission)
}

// CanEdit 报告给定权限列表是否可以编辑此配置。
//
// 对于 system scope 配置，需要检查 EditPermission。
// 对于 user scope 配置，用户可以编辑自己的覆盖值（此检查在 Service 层处理）。
func (s *Setting) CanEdit(permissions []string) bool {
	return matchPermission(permissions, s.EditPermission)
}

// matchPermission 检查权限列表中是否有匹配的权限。
//
// 支持通配符匹配规则：
//   - 精确匹配
//   - "*" 匹配任意单段
//   - "*:*:*" 匹配所有
func matchPermission(permissions []string, required string) bool {
	if required == "" {
		return true // 无需权限
	}

	requiredParts := strings.Split(required, ":")
	if len(requiredParts) != 3 {
		return false // 无效的权限格式
	}

	for _, perm := range permissions {
		permParts := strings.Split(perm, ":")
		if len(permParts) != 3 {
			continue
		}

		match := true
		for i := range 3 {
			if permParts[i] != "*" && permParts[i] != requiredParts[i] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

// =============================================================================
// 查询方法
// =============================================================================

// BelongsToCategoryID 报告是否属于指定分类 ID。
func (s *Setting) BelongsToCategoryID(categoryID uint) bool {
	return s.CategoryID == categoryID
}

// HasValidationRule 报告是否配置了自定义验证规则。
//
// 检查 Validation 字段是否非空。
func (s *Setting) HasValidationRule() bool {
	return s.Validation != ""
}

// IsRequired 报告是否为必填配置。
//
// 通过检查 Validation 中的 required 字段判断。
func (s *Setting) IsRequired() bool {
	if s.Validation == "" {
		return false
	}
	// 简单检查是否包含 required: true
	return strings.Contains(s.Validation, `"required":true`) ||
		strings.Contains(s.Validation, `"required": true`)
}

// GetKeyCategory 从 Key 提取 category 部分。
//
// 例如 "general.site_name" 返回 "general"。
func (s *Setting) GetKeyCategory() string {
	key, err := NewSettingKey(s.Key)
	if err != nil {
		return ""
	}
	return key.Category()
}

// GetKeyName 从 Key 提取 name 部分。
//
// 例如 "general.site_name" 返回 "site_name"。
func (s *Setting) GetKeyName() string {
	key, err := NewSettingKey(s.Key)
	if err != nil {
		return ""
	}
	return key.Name()
}

// =============================================================================
// 值处理方法
// =============================================================================

// CoerceValue 尝试将任意值转换为正确类型。
//
// 支持的转换：
//   - string -> number（解析为 float64）
//   - string -> boolean（"true"/"false"）
//   - number -> string（格式化）
func (s *Setting) CoerceValue(raw any) (any, error) {
	if raw == nil {
		return nil, nil //nolint:nilnil // nil input is valid and returns nil output
	}

	// 如果类型已匹配，直接返回
	if err := s.ValidateValue(raw); err == nil {
		return raw, nil
	}

	// 尝试类型转换
	switch s.ValueType {
	case ValueTypeString:
		return coerceToString(raw)
	case ValueTypeNumber:
		return coerceToNumber(raw)
	case ValueTypeBoolean:
		return coerceToBool(raw)
	case ValueTypeJSON:
		return coerceToJSON(raw)
	default:
		return nil, ErrInvalidValueType
	}
}

// GetDefaultValue 返回默认值。
func (s *Setting) GetDefaultValue() any {
	return s.DefaultValue
}

// =============================================================================
// 状态变更方法
// =============================================================================

// UpdateDefault 更新默认值。
//
// 验证新值是否与 ValueType 匹配，不匹配则返回错误。
func (s *Setting) UpdateDefault(value any) error {
	if err := s.ValidateValue(value); err != nil {
		return err
	}
	s.DefaultValue = value
	return nil
}

// UpdateLabel 更新显示标签。
func (s *Setting) UpdateLabel(label string) {
	s.Label = label
}

// UpdateOrder 更新排序权重。
func (s *Setting) UpdateOrder(order int) {
	s.Order = order
}

// =============================================================================
// 常量定义
// =============================================================================

// Scope 常量。
// 决定配置的作用域和访问方式。
const (
	ScopeSystem = "system" // 系统设置，全局唯一，管理员直接修改 DefaultValue
	ScopeUser   = "user"   // 用户设置，DefaultValue 作为初始值，用户可覆盖
)

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

// =============================================================================
// 辅助函数
// =============================================================================

func coerceToString(v any) (string, error) {
	switch val := v.(type) {
	case string:
		return val, nil
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", val), nil
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", val), nil
	case float32, float64:
		return fmt.Sprintf("%v", val), nil
	case bool:
		return strconv.FormatBool(val), nil
	default:
		return "", ErrInvalidValue
	}
}

func coerceToNumber(v any) (float64, error) {
	switch val := v.(type) {
	case float64:
		return val, nil
	case float32:
		return float64(val), nil
	case int:
		return float64(val), nil
	case int64:
		return float64(val), nil
	case string:
		return strconv.ParseFloat(val, 64)
	case json.Number:
		return val.Float64()
	default:
		return 0, ErrInvalidValue
	}
}

func coerceToBool(v any) (bool, error) {
	switch val := v.(type) {
	case bool:
		return val, nil
	case string:
		return strconv.ParseBool(val)
	case int, int64:
		return val != 0, nil
	case float64:
		return val != 0, nil
	default:
		return false, ErrInvalidValue
	}
}

func coerceToJSON(v any) (any, error) {
	switch val := v.(type) {
	case map[string]any:
		return val, nil
	case []any:
		return val, nil
	case string:
		// 尝试解析 JSON 字符串
		var result any
		if err := json.Unmarshal([]byte(val), &result); err != nil {
			return nil, ErrInvalidValue
		}
		return result, nil
	default:
		return nil, ErrInvalidValue
	}
}
