package setting

import (
	"time"
)

// Setting 系统配置
type Setting struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Key       string    `gorm:"uniqueIndex;size:100;not null" json:"key"`
	Value     string    `gorm:"type:text" json:"value"`
	Category  string    `gorm:"size:50;index;not null" json:"category"`
	ValueType string    `gorm:"size:20;default:'string'" json:"value_type"` // string, number, boolean, json
	Label     string    `gorm:"size:200" json:"label"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 指定表名
func (Setting) TableName() string {
	return "settings"
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
