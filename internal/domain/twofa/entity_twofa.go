package twofa

import (
	"time"

	"gorm.io/gorm"
)

// TwoFA 用户双因素认证配置实体
type TwoFA struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// 用户关联
	UserID uint `gorm:"uniqueIndex;not null" json:"user_id"`

	// 2FA 状态
	Enabled bool `gorm:"default:false;not null" json:"enabled"` // 是否启用 2FA

	// TOTP 密钥（加密存储）
	Secret string `gorm:"size:255;not null" json:"-"` // TOTP 密钥（Base32 编码）

	// 恢复码（加密存储，JSON 数组）
	RecoveryCodes RecoveryCodes `gorm:"type:text" json:"-"` // 恢复码列表

	// 设置信息
	SetupCompletedAt *time.Time `json:"setup_completed_at,omitempty"` // 完成设置的时间
	LastUsedAt       *time.Time `json:"last_used_at,omitempty"`       // 最后使用时间
}

// TableName 指定表名
func (TwoFA) TableName() string {
	return "user_2fas"
}

// HasRecoveryCodes 检查是否有可用的恢复码
func (t *TwoFA) HasRecoveryCodes() bool {
	return len(t.RecoveryCodes) > 0
}

// BeforeCreate GORM Hook: 创建前初始化
func (t *TwoFA) BeforeCreate(tx *gorm.DB) error {
	if len(t.RecoveryCodes) == 0 {
		t.RecoveryCodes = []string{}
	}
	return nil
}
