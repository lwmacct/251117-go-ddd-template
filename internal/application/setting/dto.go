package setting

import "time"

// SettingDTO 设置响应 DTO
type SettingDTO struct {
	ID        uint      `json:"id"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	Category  string    `json:"category"`
	ValueType string    `json:"value_type"`
	Label     string    `json:"label"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateSettingResultDTO 创建设置结果 DTO
type CreateSettingResultDTO struct {
	ID uint `json:"id"`
}
