// Package setting 定义设置模块的 DTO
package setting

import "time"

// SettingResponse 设置响应 DTO
type SettingResponse struct {
	ID        uint      `json:"id"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	Category  string    `json:"category"`
	ValueType string    `json:"value_type"`
	Label     string    `json:"label"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
