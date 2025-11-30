// Package pat 提供个人访问令牌领域模型单元测试。
package pat

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// newTestPAT 创建测试用 PAT。
func newTestPAT(status string, expiresAt *time.Time) *PersonalAccessToken {
	return &PersonalAccessToken{
		ID:          1,
		UserID:      100,
		Name:        "Test Token",
		Token:       "hashed_token_value",
		TokenPrefix: "ghp_xxxx",
		Permissions: PermissionList{"read", "write"},
		Status:      status,
		ExpiresAt:   expiresAt,
		CreatedAt:   time.Now(),
	}
}

func TestPersonalAccessToken_IsExpired(t *testing.T) {
	now := time.Now()
	pastTime := now.Add(-24 * time.Hour)
	futureTime := now.Add(24 * time.Hour)

	tests := []struct {
		name      string
		expiresAt *time.Time
		want      bool
	}{
		{
			name:      "无过期时间 - 永久有效",
			expiresAt: nil,
			want:      false,
		},
		{
			name:      "已过期",
			expiresAt: &pastTime,
			want:      true,
		},
		{
			name:      "未过期",
			expiresAt: &futureTime,
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pat := newTestPAT("active", tt.expiresAt)
			got := pat.IsExpired()
			assert.Equal(t, tt.want, got, "PersonalAccessToken.IsExpired()")
		})
	}
}

func TestPersonalAccessToken_IsActive(t *testing.T) {
	now := time.Now()
	pastTime := now.Add(-24 * time.Hour)
	futureTime := now.Add(24 * time.Hour)

	tests := []struct {
		name      string
		status    string
		expiresAt *time.Time
		want      bool
	}{
		{
			name:      "活跃且未过期",
			status:    "active",
			expiresAt: &futureTime,
			want:      true,
		},
		{
			name:      "活跃且永久有效",
			status:    "active",
			expiresAt: nil,
			want:      true,
		},
		{
			name:      "活跃但已过期",
			status:    "active",
			expiresAt: &pastTime,
			want:      false,
		},
		{
			name:      "已禁用",
			status:    "disabled",
			expiresAt: &futureTime,
			want:      false,
		},
		{
			name:      "已过期状态",
			status:    "expired",
			expiresAt: nil,
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pat := newTestPAT(tt.status, tt.expiresAt)
			got := pat.IsActive()
			assert.Equal(t, tt.want, got, "PersonalAccessToken.IsActive()")
		})
	}
}

func TestPersonalAccessToken_ToListItem(t *testing.T) {
	now := time.Now()
	expiresAt := now.Add(30 * 24 * time.Hour)
	lastUsedAt := now.Add(-1 * time.Hour)

	pat := &PersonalAccessToken{
		ID:          42,
		Name:        "CI Token",
		TokenPrefix: "ghp_1234",
		Permissions: PermissionList{"repo:read", "repo:write"},
		ExpiresAt:   &expiresAt,
		LastUsedAt:  &lastUsedAt,
		Status:      "active",
		CreatedAt:   now,
	}

	item := pat.ToListItem()

	assert.Equal(t, pat.ID, item.ID, "ID 应该相等")
	assert.Equal(t, pat.Name, item.Name, "Name 应该相等")
	assert.Equal(t, pat.TokenPrefix, item.TokenPrefix, "TokenPrefix 应该相等")
	assert.Equal(t, []string(pat.Permissions), item.Permissions, "Permissions 应该相等")
	assert.Equal(t, pat.ExpiresAt, item.ExpiresAt, "ExpiresAt 应该相等")
	assert.Equal(t, pat.LastUsedAt, item.LastUsedAt, "LastUsedAt 应该相等")
	assert.Equal(t, pat.Status, item.Status, "Status 应该相等")
	assert.Equal(t, pat.CreatedAt, item.CreatedAt, "CreatedAt 应该相等")
}

func TestPersonalAccessToken_ToListItem_NilFields(t *testing.T) {
	pat := &PersonalAccessToken{
		ID:          1,
		Name:        "Simple Token",
		TokenPrefix: "ghp_xxxx",
		Permissions: nil,
		ExpiresAt:   nil,
		LastUsedAt:  nil,
		Status:      "active",
	}

	item := pat.ToListItem()

	assert.Nil(t, item.ExpiresAt, "ExpiresAt 应该为 nil")
	assert.Nil(t, item.LastUsedAt, "LastUsedAt 应该为 nil")
	assert.Nil(t, item.Permissions, "Permissions 应该为 nil")
}

func TestPersonalAccessToken_EdgeCases(t *testing.T) {
	t.Run("边界时间 - 刚刚过期", func(t *testing.T) {
		// 1 纳秒前过期
		justExpired := time.Now().Add(-1 * time.Nanosecond)
		pat := newTestPAT("active", &justExpired)

		assert.True(t, pat.IsExpired(), "刚刚过期的 Token 应该返回 true")
		assert.False(t, pat.IsActive(), "刚刚过期的 Token 不应该是活跃的")
	})

	t.Run("零值时间", func(t *testing.T) {
		zeroTime := time.Time{}
		pat := newTestPAT("active", &zeroTime)

		// 零值时间是过去的，应该过期
		assert.True(t, pat.IsExpired(), "零值时间应该被视为过期")
	})
}
