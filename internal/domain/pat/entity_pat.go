package pat

import "time"

// PersonalAccessToken represents a personal access token for API authentication
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
	Status     string     `json:"status"`                 // active, revoked, expired

	IPWhitelist StringList `json:"ip_whitelist,omitempty"` // IP 白名单（可选）
	Description string     `json:"description,omitempty"`  // 描述
}

// IsExpired checks if the token has expired
func (p *PersonalAccessToken) IsExpired() bool {
	if p.ExpiresAt == nil {
		return false
	}
	return p.ExpiresAt.Before(time.Now())
}

// IsActive checks if the token is active and not expired
func (p *PersonalAccessToken) IsActive() bool {
	return p.Status == "active" && !p.IsExpired()
}

// ToListItem converts PAT to TokenListItem
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

// TokenListItem represents a single PAT in the list (without full token)
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
