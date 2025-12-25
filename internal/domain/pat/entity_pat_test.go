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

func TestPersonalAccessToken_IsIPAllowed(t *testing.T) {
	tests := []struct {
		name        string
		ipWhitelist StringList
		ip          string
		want        bool
	}{
		{
			name:        "空白名单 - 允许所有",
			ipWhitelist: nil,
			ip:          "192.168.1.1",
			want:        true,
		},
		{
			name:        "空数组白名单 - 允许所有",
			ipWhitelist: StringList{},
			ip:          "192.168.1.1",
			want:        true,
		},
		{
			name:        "IP 在白名单中",
			ipWhitelist: StringList{"192.168.1.1", "10.0.0.1"},
			ip:          "192.168.1.1",
			want:        true,
		},
		{
			name:        "IP 不在白名单中",
			ipWhitelist: StringList{"192.168.1.1", "10.0.0.1"},
			ip:          "172.16.0.1",
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pat := &PersonalAccessToken{IPWhitelist: tt.ipWhitelist}
			assert.Equal(t, tt.want, pat.IsIPAllowed(tt.ip))
		})
	}
}

func TestPersonalAccessToken_HasPermission(t *testing.T) {
	pat := &PersonalAccessToken{
		Permissions: PermissionList{"read", "write", "admin"},
	}

	tests := []struct {
		name  string
		scope string
		want  bool
	}{
		{"存在的权限", "read", true},
		{"存在的权限 2", "admin", true},
		{"不存在的权限", "delete", false},
		{"空权限", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, pat.HasPermission(tt.scope))
		})
	}
}

func TestPersonalAccessToken_HasAnyPermission(t *testing.T) {
	pat := &PersonalAccessToken{
		Permissions: PermissionList{"read", "write"},
	}

	tests := []struct {
		name   string
		scopes []string
		want   bool
	}{
		{"有一个匹配", []string{"read", "delete"}, true},
		{"全部匹配", []string{"read", "write"}, true},
		{"无匹配", []string{"delete", "admin"}, false},
		{"空数组", []string{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, pat.HasAnyPermission(tt.scopes...))
		})
	}
}

func TestPersonalAccessToken_HasAllPermissions(t *testing.T) {
	pat := &PersonalAccessToken{
		Permissions: PermissionList{"read", "write", "admin"},
	}

	tests := []struct {
		name   string
		scopes []string
		want   bool
	}{
		{"全部存在", []string{"read", "write"}, true},
		{"部分存在", []string{"read", "delete"}, false},
		{"全部不存在", []string{"delete", "super"}, false},
		{"空数组", []string{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, pat.HasAllPermissions(tt.scopes...))
		})
	}
}

func TestPersonalAccessToken_DisableEnable(t *testing.T) {
	t.Run("禁用 Token", func(t *testing.T) {
		pat := newTestPAT("active", nil)
		assert.False(t, pat.IsDisabled())

		pat.Disable()
		assert.True(t, pat.IsDisabled())
		assert.Equal(t, StatusDisabled, pat.Status)
		assert.False(t, pat.IsActive())
	})

	t.Run("启用 Token", func(t *testing.T) {
		pat := newTestPAT("disabled", nil)
		assert.True(t, pat.IsDisabled())

		pat.Enable()
		assert.False(t, pat.IsDisabled())
		assert.Equal(t, StatusActive, pat.Status)
		assert.True(t, pat.IsActive())
	})

	t.Run("标记过期", func(t *testing.T) {
		pat := newTestPAT("active", nil)
		pat.MarkExpired()
		assert.Equal(t, StatusExpired, pat.Status)
		assert.False(t, pat.IsActive())
	})
}

func TestPersonalAccessToken_GetPermissionCount(t *testing.T) {
	tests := []struct {
		name        string
		permissions PermissionList
		want        int
	}{
		{"有权限", PermissionList{"read", "write"}, 2},
		{"空权限", PermissionList{}, 0},
		{"nil 权限", nil, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pat := &PersonalAccessToken{Permissions: tt.permissions}
			assert.Equal(t, tt.want, pat.GetPermissionCount())
		})
	}
}

func TestPersonalAccessToken_CanBeUsed(t *testing.T) {
	futureTime := time.Now().Add(24 * time.Hour)

	tests := []struct {
		name        string
		status      string
		expiresAt   *time.Time
		ipWhitelist StringList
		ip          string
		want        bool
	}{
		{
			name:        "活跃且 IP 允许",
			status:      "active",
			expiresAt:   &futureTime,
			ipWhitelist: nil,
			ip:          "192.168.1.1",
			want:        true,
		},
		{
			name:        "活跃但 IP 不允许",
			status:      "active",
			expiresAt:   &futureTime,
			ipWhitelist: StringList{"10.0.0.1"},
			ip:          "192.168.1.1",
			want:        false,
		},
		{
			name:        "禁用状态",
			status:      "disabled",
			expiresAt:   &futureTime,
			ipWhitelist: nil,
			ip:          "192.168.1.1",
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pat := &PersonalAccessToken{
				Status:      tt.status,
				ExpiresAt:   tt.expiresAt,
				IPWhitelist: tt.ipWhitelist,
			}
			assert.Equal(t, tt.want, pat.CanBeUsed(tt.ip))
		})
	}
}
