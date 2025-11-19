package setting

import (
	"time"
)

// Setting 系统配置实体
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

// Category 常量
const (
	CategoryGeneral      = "general"
	CategorySecurity     = "security"
	CategoryNotification = "notification"
	CategoryBackup       = "backup"
)

// ValueType 常量
const (
	ValueTypeString  = "string"
	ValueTypeNumber  = "number"
	ValueTypeBoolean = "boolean"
	ValueTypeJSON    = "json"
)
