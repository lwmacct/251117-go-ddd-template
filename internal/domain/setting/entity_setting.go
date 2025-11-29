package setting

import (
	"time"
)

// Setting 系统配置实体。
// 采用 Key-Value 模式存储配置项，支持分类和类型标注。
type Setting struct {
	ID        uint      `json:"id"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	Category  string    `json:"category"`
	ValueType string    `json:"value_type"`
	Label     string    `json:"label"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
