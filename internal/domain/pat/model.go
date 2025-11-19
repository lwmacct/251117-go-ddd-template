package pat

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

// PersonalAccessToken represents a personal access token for API authentication
type PersonalAccessToken struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	UserID      uint           `gorm:"not null;index" json:"user_id"`                // 所属用户
	Name        string         `gorm:"size:100;not null" json:"name"`                // Token 名称
	Token       string         `gorm:"size:255;uniqueIndex;not null" json:"-"`       // Token 哈希值（不返回）
	TokenPrefix string         `gorm:"size:20;not null;index" json:"token_prefix"`   // Token 前缀（明文）

	Permissions PermissionList `gorm:"type:jsonb;not null" json:"permissions"`       // 权限列表（JSON）
	Scopes      string         `gorm:"type:text" json:"scopes,omitempty"`            // 权限范围描述

	ExpiresAt   *time.Time     `json:"expires_at,omitempty"`                         // 过期时间（nil=永久）
	LastUsedAt  *time.Time     `json:"last_used_at,omitempty"`                       // 最后使用时间
	Status      string         `gorm:"size:20;not null;default:'active';index" json:"status"` // active, revoked, expired

	IPWhitelist StringList     `gorm:"type:jsonb" json:"ip_whitelist,omitempty"`     // IP 白名单（可选）
	Description string         `gorm:"type:text" json:"description,omitempty"`       // 描述
}

// TableName specifies the table name for PersonalAccessToken model
func (PersonalAccessToken) TableName() string {
	return "personal_access_tokens"
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

// PermissionList is a custom type for JSON array of permissions
type PermissionList []string

// Scan implements sql.Scanner interface
func (p *PermissionList) Scan(value interface{}) error {
	if value == nil {
		*p = []string{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan PermissionList")
	}

	return json.Unmarshal(bytes, p)
}

// Value implements driver.Valuer interface
func (p PermissionList) Value() (driver.Value, error) {
	if p == nil {
		return json.Marshal([]string{})
	}
	return json.Marshal(p)
}

// StringList is a custom type for JSON array of strings
type StringList []string

// Scan implements sql.Scanner interface
func (s *StringList) Scan(value interface{}) error {
	if value == nil {
		*s = []string{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan StringList")
	}

	return json.Unmarshal(bytes, s)
}

// Value implements driver.Valuer interface
func (s StringList) Value() (driver.Value, error) {
	if s == nil {
		return json.Marshal([]string{})
	}
	return json.Marshal(s)
}

// CreateTokenRequest represents the request to create a new PAT
type CreateTokenRequest struct {
	UserID      uint     `json:"-"`                                     // 从 Context 获取
	Name        string   `json:"name" binding:"required,min=1,max=100"` // Token 名称
	Permissions []string `json:"permissions" binding:"required,min=1"`  // 选择的权限
	ExpiresIn   *int     `json:"expires_in"`                            // 过期天数：7, 30, 90, null(永久)
	IPWhitelist []string `json:"ip_whitelist"`                          // IP 白名单
	Description string   `json:"description"`                           // 描述
}

// TokenResponse represents the response after creating a new PAT
type TokenResponse struct {
	Token       string    `json:"token"`         // 明文 Token（只显示一次！）
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	TokenPrefix string    `json:"token_prefix"`
	Permissions []string  `json:"permissions"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// TokenListItem represents a single PAT in the list (without full token)
type TokenListItem struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	TokenPrefix string     `json:"token_prefix"`  // 用于识别
	Permissions []string   `json:"permissions"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	LastUsedAt  *time.Time `json:"last_used_at,omitempty"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
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
