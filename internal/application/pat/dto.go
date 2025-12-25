package pat

import (
	"time"
)

// CreateTokenDTO 创建令牌请求 DTO
type CreateTokenDTO struct {
	Name        string   `json:"name" binding:"required,min=3,max=100"`
	Permissions []string `json:"permissions,omitempty"`  // 可选，限制令牌权限范围（为空则默认用户全部权限）
	ExpiresAt   *string  `json:"expires_at,omitempty"`   // 可选，过期时间（RFC3339 或 yyyy-MM-ddTHH:mm）
	ExpiresIn   *int     `json:"expires_in,omitempty"`   // 可选，以天为单位的有效期（兜底，前端未使用时可忽略）
	IPWhitelist []string `json:"ip_whitelist,omitempty"` // 可选，IP 白名单
	Description string   `json:"description,omitempty"`  // 可选，备注
}

// TokenDTO 令牌响应 DTO（不含明文 token）
type TokenDTO struct {
	ID          uint       `json:"id"`
	UserID      uint       `json:"user_id"`
	Name        string     `json:"name"`
	TokenPrefix string     `json:"token_prefix"`
	Permissions []string   `json:"permissions"`
	IPWhitelist []string   `json:"ip_whitelist,omitempty"`
	Status      string     `json:"status"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	LastUsedAt  *time.Time `json:"last_used_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// CreateTokenResultDTO 令牌创建响应（包含一次性明文 token）
type CreateTokenResultDTO struct {
	Token      *TokenDTO `json:"token"`
	PlainToken string    `json:"plain_token"`
}

// TokenListDTO 令牌列表响应 DTO
type TokenListDTO struct {
	Tokens []*TokenDTO `json:"tokens"`
	Total  int64       `json:"total"`
}

// TokenInfoDTO Token 信息响应（与 TokenDTO 结构相同，用于语义表达）
type TokenInfoDTO struct {
	ID          uint       `json:"id"`
	UserID      uint       `json:"user_id"`
	Name        string     `json:"name"`
	TokenPrefix string     `json:"token_prefix"`
	Permissions []string   `json:"permissions"`
	IPWhitelist []string   `json:"ip_whitelist,omitempty"`
	Status      string     `json:"status"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	LastUsedAt  *time.Time `json:"last_used_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
