// Package twofa 定义双因素认证 (Two-Factor Authentication) 领域模型。
//
// 本包实现基于 TOTP (时间同步一次性密码) 的双因素认证，兼容：
//   - Google Authenticator
//   - Microsoft Authenticator
//   - 其他标准 TOTP 应用
//
// 核心功能：
//   - TwoFA 实体：存储用户的 2FA 配置和状态
//   - RecoveryCodes：一次性恢复码，用于设备丢失时的账户恢复
//   - TOTP 密钥管理：Secret 字段存储 Base32 编码的密钥
package twofa

import "time"

// TwoFA 用户双因素认证配置实体
type TwoFA struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	// 用户关联
	UserID uint `json:"user_id"`

	// 2FA 状态
	Enabled bool `json:"enabled"` // 是否启用 2FA

	// TOTP 密钥（加密存储）
	Secret string `json:"-"` // TOTP 密钥（Base32 编码）

	// 恢复码（加密存储，JSON 数组）
	RecoveryCodes RecoveryCodes `json:"-"` // 恢复码列表

	// 设置信息
	SetupCompletedAt *time.Time `json:"setup_completed_at,omitempty"` // 完成设置的时间
	LastUsedAt       *time.Time `json:"last_used_at,omitempty"`       // 最后使用时间
}

// HasRecoveryCodes 检查是否有可用的恢复码
func (t *TwoFA) HasRecoveryCodes() bool {
	return len(t.RecoveryCodes) > 0
}

// IsEnabled 检查 2FA 是否已启用
func (t *TwoFA) IsEnabled() bool {
	return t.Enabled
}

// IsSetupComplete 检查 2FA 设置是否已完成
func (t *TwoFA) IsSetupComplete() bool {
	return t.SetupCompletedAt != nil
}

// Enable 启用 2FA
func (t *TwoFA) Enable() {
	t.Enabled = true
	now := time.Now()
	t.SetupCompletedAt = &now
}

// Disable 禁用 2FA
func (t *TwoFA) Disable() {
	t.Enabled = false
}

// MarkUsed 标记最后使用时间
func (t *TwoFA) MarkUsed() {
	now := time.Now()
	t.LastUsedAt = &now
}

// UseRecoveryCode 使用恢复码（从列表中移除已使用的码）
// 返回 true 表示恢复码有效并已使用，false 表示无效
func (t *TwoFA) UseRecoveryCode(code string) bool {
	for i, rc := range t.RecoveryCodes {
		if rc == code {
			// 移除已使用的恢复码
			t.RecoveryCodes = append(t.RecoveryCodes[:i], t.RecoveryCodes[i+1:]...)
			t.MarkUsed()
			return true
		}
	}
	return false
}

// GetRecoveryCodesCount 获取剩余恢复码数量
func (t *TwoFA) GetRecoveryCodesCount() int {
	return len(t.RecoveryCodes)
}

// SetRecoveryCodes 设置恢复码（覆盖现有的）
func (t *TwoFA) SetRecoveryCodes(codes []string) {
	t.RecoveryCodes = codes
}

// HasSecret 检查是否已配置 TOTP 密钥
func (t *TwoFA) HasSecret() bool {
	return t.Secret != ""
}

// ClearSecret 清除 TOTP 密钥（禁用时使用）
func (t *TwoFA) ClearSecret() {
	t.Secret = ""
}

// Reset 重置 2FA 配置到初始状态
func (t *TwoFA) Reset() {
	t.Enabled = false
	t.Secret = ""
	t.RecoveryCodes = nil
	t.SetupCompletedAt = nil
	t.LastUsedAt = nil
}
