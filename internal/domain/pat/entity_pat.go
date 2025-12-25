package pat

import (
	"slices"
	"time"
)

// PersonalAccessToken 个人访问令牌实体，用于 API 认证
type PersonalAccessToken struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`

	UserID      uint   `json:"user_id"`      // 所属用户
	Name        string `json:"name"`         // Token 名称
	Token       string `json:"-"`            // Token 哈希值（不返回）
	TokenPrefix string `json:"token_prefix"` // Token 前缀（明文）

	Permissions PermissionList `json:"permissions"`      // 权限列表（JSON）
	Scopes      string         `json:"scopes,omitempty"` // 权限范围描述

	ExpiresAt  *time.Time `json:"expires_at,omitempty"`   // 过期时间（nil=永久）
	LastUsedAt *time.Time `json:"last_used_at,omitempty"` // 最后使用时间
	Status     string     `json:"status"`                 // active, disabled, expired

	IPWhitelist StringList `json:"ip_whitelist,omitempty"` // IP 白名单（可选）
	Description string     `json:"description,omitempty"`  // 描述
}

// IsExpired 检查 Token 是否已过期
func (p *PersonalAccessToken) IsExpired() bool {
	if p.ExpiresAt == nil {
		return false
	}
	return p.ExpiresAt.Before(time.Now())
}

// IsActive 检查 Token 是否处于活跃状态（未过期且未禁用）
func (p *PersonalAccessToken) IsActive() bool {
	return p.Status == "active" && !p.IsExpired()
}

// Token 状态常量
const (
	StatusActive   = "active"
	StatusDisabled = "disabled"
	StatusExpired  = "expired"
)

// IsIPAllowed 检查给定 IP 是否允许使用此 Token。
// 如果 IP 白名单为空，则允许所有 IP。
func (p *PersonalAccessToken) IsIPAllowed(ip string) bool {
	if len(p.IPWhitelist) == 0 {
		return true
	}
	return slices.Contains(p.IPWhitelist, ip)
}

// HasPermission 检查 Token 是否有指定权限
func (p *PersonalAccessToken) HasPermission(scope string) bool {
	return slices.Contains(p.Permissions, scope)
}

// HasAnyPermission 检查 Token 是否有任一指定权限
func (p *PersonalAccessToken) HasAnyPermission(scopes ...string) bool {
	return slices.ContainsFunc(scopes, p.HasPermission)
}

// HasAllPermissions 检查 Token 是否有所有指定权限
func (p *PersonalAccessToken) HasAllPermissions(scopes ...string) bool {
	for _, scope := range scopes {
		if !p.HasPermission(scope) {
			return false
		}
	}
	return true
}

// Disable 禁用 Token
func (p *PersonalAccessToken) Disable() {
	p.Status = StatusDisabled
}

// Enable 启用 Token
func (p *PersonalAccessToken) Enable() {
	p.Status = StatusActive
}

// MarkExpired 标记 Token 为已过期
func (p *PersonalAccessToken) MarkExpired() {
	p.Status = StatusExpired
}

// IsDisabled 检查 Token 是否被禁用
func (p *PersonalAccessToken) IsDisabled() bool {
	return p.Status == StatusDisabled
}

// GetPermissionCount 返回权限数量
func (p *PersonalAccessToken) GetPermissionCount() int {
	return len(p.Permissions)
}

// CanBeUsed 检查 Token 是否可以被使用（活跃且 IP 允许）
func (p *PersonalAccessToken) CanBeUsed(ip string) bool {
	return p.IsActive() && p.IsIPAllowed(ip)
}

// ToListItem 将 PAT 转换为列表项（不含完整 Token）
func (p *PersonalAccessToken) ToListItem() *TokenListItem {
	return &TokenListItem{
		ID:          p.ID,
		Name:        p.Name,
		TokenPrefix: p.TokenPrefix,
		Permissions: p.Permissions,
		ExpiresAt:   p.ExpiresAt,
		LastUsedAt:  p.LastUsedAt,
		Status:      p.Status,
		CreatedAt:   p.CreatedAt,
	}
}

// TokenListItem PAT 列表项（不含完整 Token）
type TokenListItem struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	TokenPrefix string     `json:"token_prefix"` // 用于识别
	Permissions []string   `json:"permissions"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	LastUsedAt  *time.Time `json:"last_used_at,omitempty"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
}
