package pat

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/pat"
)

func TestToTokenResponse(t *testing.T) {
	t.Run("转换正常令牌", func(t *testing.T) {
		now := time.Now()
		expiresAt := now.Add(30 * 24 * time.Hour)
		lastUsedAt := now.Add(-1 * time.Hour)
		token := &pat.PersonalAccessToken{
			ID:          1,
			UserID:      100,
			Name:        "My API Token",
			TokenPrefix: "abc123",
			Permissions: []string{"read", "write"},
			IPWhitelist: []string{"192.168.1.0/24"},
			Status:      "active",
			ExpiresAt:   &expiresAt,
			LastUsedAt:  &lastUsedAt,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		result := ToTokenResponse(token)

		assert.NotNil(t, result)
		assert.Equal(t, uint(1), result.ID)
		assert.Equal(t, uint(100), result.UserID)
		assert.Equal(t, "My API Token", result.Name)
		assert.Equal(t, "abc123", result.TokenPrefix)
		assert.Equal(t, []string{"read", "write"}, result.Permissions)
		assert.Equal(t, []string{"192.168.1.0/24"}, result.IPWhitelist)
		assert.Equal(t, "active", result.Status)
		assert.NotNil(t, result.ExpiresAt)
		assert.NotNil(t, result.LastUsedAt)
		assert.Equal(t, now, result.CreatedAt)
		assert.Equal(t, now, result.UpdatedAt)
	})

	t.Run("转换nil返回nil", func(t *testing.T) {
		result := ToTokenResponse(nil)
		assert.Nil(t, result)
	})

	t.Run("转换无过期时间令牌", func(t *testing.T) {
		token := &pat.PersonalAccessToken{
			ID:          2,
			UserID:      101,
			Name:        "永久令牌",
			TokenPrefix: "def456",
			Status:      "active",
			ExpiresAt:   nil,
			LastUsedAt:  nil,
		}

		result := ToTokenResponse(token)

		assert.NotNil(t, result)
		assert.Nil(t, result.ExpiresAt)
		assert.Nil(t, result.LastUsedAt)
	})

	t.Run("转换禁用令牌", func(t *testing.T) {
		token := &pat.PersonalAccessToken{
			ID:     3,
			UserID: 102,
			Name:   "已禁用令牌",
			Status: "disabled",
		}

		result := ToTokenResponse(token)

		assert.NotNil(t, result)
		assert.Equal(t, "disabled", result.Status)
	})
}

func TestToCreateTokenResponse(t *testing.T) {
	t.Run("转换创建令牌响应", func(t *testing.T) {
		now := time.Now()
		token := &pat.PersonalAccessToken{
			ID:          1,
			UserID:      100,
			Name:        "New Token",
			TokenPrefix: "xyz789",
			Status:      "active",
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		plainToken := "pat_xyz789_secrettoken123456" //nolint:gosec // test credential

		result := ToCreateTokenResponse(token, plainToken)

		assert.NotNil(t, result)
		assert.Equal(t, "pat_xyz789_secrettoken123456", result.PlainToken)
		assert.NotNil(t, result.Token)
		assert.Equal(t, uint(1), result.Token.ID)
		assert.Equal(t, "New Token", result.Token.Name)
		assert.Equal(t, "xyz789", result.Token.TokenPrefix)
	})

	t.Run("转换nil返回nil", func(t *testing.T) {
		result := ToCreateTokenResponse(nil, "some_token")
		assert.Nil(t, result)
	})
}

func TestToTokenListResponse(t *testing.T) {
	t.Run("转换令牌列表", func(t *testing.T) {
		now := time.Now()
		expiresAt := now.Add(30 * 24 * time.Hour)
		items := []*pat.TokenListItem{
			{
				ID:          1,
				Name:        "Token 1",
				TokenPrefix: "tok1",
				Permissions: []string{"read"},
				Status:      "active",
				ExpiresAt:   &expiresAt,
				LastUsedAt:  nil,
				CreatedAt:   now,
			},
			{
				ID:          2,
				Name:        "Token 2",
				TokenPrefix: "tok2",
				Permissions: []string{"read", "write"},
				Status:      "disabled",
				ExpiresAt:   nil,
				LastUsedAt:  nil,
				CreatedAt:   now,
			},
		}

		result := ToTokenListResponse(items)

		assert.NotNil(t, result)
		assert.Equal(t, int64(2), result.Total)
		assert.Len(t, result.Tokens, 2)
		assert.Equal(t, uint(1), result.Tokens[0].ID)
		assert.Equal(t, "Token 1", result.Tokens[0].Name)
		assert.Equal(t, "tok1", result.Tokens[0].TokenPrefix)
		assert.Equal(t, uint(2), result.Tokens[1].ID)
		assert.Equal(t, "Token 2", result.Tokens[1].Name)
	})

	t.Run("转换空列表", func(t *testing.T) {
		items := []*pat.TokenListItem{}

		result := ToTokenListResponse(items)

		assert.NotNil(t, result)
		assert.Equal(t, int64(0), result.Total)
		assert.Empty(t, result.Tokens)
	})

	t.Run("转换nil列表", func(t *testing.T) {
		result := ToTokenListResponse(nil)

		assert.NotNil(t, result)
		assert.Equal(t, int64(0), result.Total)
		assert.Empty(t, result.Tokens)
	})
}

func TestToTokenInfoResponse(t *testing.T) {
	t.Run("转换令牌信息", func(t *testing.T) {
		now := time.Now()
		token := &pat.PersonalAccessToken{
			ID:          1,
			UserID:      100,
			Name:        "Info Token",
			TokenPrefix: "info123",
			Status:      "active",
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		result := ToTokenInfoResponse(token)

		assert.NotNil(t, result)
		assert.Equal(t, uint(1), result.ID)
		assert.Equal(t, "Info Token", result.Name)
		assert.Equal(t, "info123", result.TokenPrefix)
	})

	t.Run("转换nil返回nil", func(t *testing.T) {
		result := ToTokenInfoResponse(nil)
		assert.Nil(t, result)
	})
}
