// Package pat 定义个人访问令牌应用层的 DTO
package pat

import "time"

// CreateTokenRequest 创建令牌请求 DTO
type CreateTokenRequest struct {
	Name        string   `json:"name" binding:"required,min=3,max=100"`
	Permissions []string `json:"permissions,omitempty"` // 可选，限制令牌权限范围
	ExpiresAt   *int64   `json:"expires_at,omitempty"`  // 可选，Unix 时间戳（秒），不设置则永不过期
}

// TokenResponse 令牌响应 DTO（创建时返回完整令牌）
type TokenResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Token       string    `json:"token,omitempty"`       // 仅在创建时返回，其他时候为空
	Permissions []string  `json:"permissions,omitempty"` // 令牌权限范围
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	LastUsedAt  *time.Time `json:"last_used_at,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// TokenListResponse 令牌列表响应 DTO
type TokenListResponse struct {
	Tokens []*TokenResponse `json:"tokens"`
	Total  int64            `json:"total"`
}
