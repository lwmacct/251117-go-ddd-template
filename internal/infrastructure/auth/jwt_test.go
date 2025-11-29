// Package auth 提供 JWT 管理器的单元测试。
package auth

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewJWTManager(t *testing.T) {
	t.Run("创建 JWT 管理器", func(t *testing.T) {
		manager := NewJWTManager("secret", time.Hour, 24*time.Hour)

		require.NotNil(t, manager, "NewJWTManager() 不应返回 nil")
		assert.Equal(t, "secret", manager.secretKey, "secretKey 应该匹配")
		assert.Equal(t, time.Hour, manager.accessTokenDuration, "accessTokenDuration 应该匹配")
		assert.Equal(t, 24*time.Hour, manager.refreshTokenDuration, "refreshTokenDuration 应该匹配")
	})
}

func TestJWTManager_GenerateAccessToken(t *testing.T) {
	manager := NewJWTManager("test-secret", time.Hour, 24*time.Hour)

	t.Run("成功生成访问令牌", func(t *testing.T) {
		token, err := manager.GenerateAccessToken(1, "testuser", "test@example.com")

		require.NoError(t, err, "GenerateAccessToken() 应该成功")
		assert.NotEmpty(t, token, "GenerateAccessToken() 不应返回空令牌")

		// JWT 令牌应该有三部分（header.payload.signature）
		parts := strings.Split(token, ".")
		assert.Len(t, parts, 3, "JWT 令牌应该有 3 部分")
	})

	t.Run("包含正确的用户信息", func(t *testing.T) {
		token, _ := manager.GenerateAccessToken(123, "john", "john@example.com")
		claims, err := manager.ValidateToken(token)

		require.NoError(t, err, "ValidateToken() 应该成功")
		assert.Equal(t, uint(123), claims.UserID, "UserID 应该匹配")
		assert.Equal(t, "john", claims.Username, "Username 应该匹配")
		assert.Equal(t, "john@example.com", claims.Email, "Email 应该匹配")
	})

	t.Run("生成的令牌有正确的过期时间", func(t *testing.T) {
		manager := NewJWTManager("secret", time.Minute*30, time.Hour*24)
		token, _ := manager.GenerateAccessToken(1, "user", "")
		claims, _ := manager.ValidateToken(token)

		expectedExpiry := time.Now().Add(time.Minute * 30)
		actualExpiry := claims.ExpiresAt.Time

		// 允许 5 秒误差
		diff := actualExpiry.Sub(expectedExpiry)
		assert.True(t, diff >= -5*time.Second && diff <= 5*time.Second,
			"过期时间偏差应该在 5 秒以内, 实际偏差: %v", diff)
	})

	t.Run("空用户名和邮箱", func(t *testing.T) {
		token, err := manager.GenerateAccessToken(1, "", "")

		assert.NoError(t, err, "空用户名应该能生成令牌")
		assert.NotEmpty(t, token, "空用户名应该能生成非空令牌")
	})
}

func TestJWTManager_GenerateRefreshToken(t *testing.T) {
	manager := NewJWTManager("test-secret", time.Hour, 24*time.Hour)

	t.Run("成功生成刷新令牌", func(t *testing.T) {
		token, err := manager.GenerateRefreshToken(456)

		require.NoError(t, err, "GenerateRefreshToken() 应该成功")
		assert.NotEmpty(t, token, "GenerateRefreshToken() 不应返回空令牌")
	})

	t.Run("刷新令牌只包含 UserID", func(t *testing.T) {
		token, _ := manager.GenerateRefreshToken(789)
		claims, err := manager.ValidateToken(token)

		require.NoError(t, err, "ValidateToken() 应该成功")
		assert.Equal(t, uint(789), claims.UserID, "UserID 应该匹配")
		// Refresh token 不包含 username 和 email
		assert.Empty(t, claims.Username, "Refresh token 不应该包含 Username")
	})

	t.Run("刷新令牌过期时间更长", func(t *testing.T) {
		accessToken, _ := manager.GenerateAccessToken(1, "user", "")
		refreshToken, _ := manager.GenerateRefreshToken(1)

		accessClaims, _ := manager.ValidateToken(accessToken)
		refreshClaims, _ := manager.ValidateToken(refreshToken)

		assert.True(t, refreshClaims.ExpiresAt.After(accessClaims.ExpiresAt.Time),
			"刷新令牌过期时间应该比访问令牌长")
	})
}

func TestJWTManager_ValidateToken(t *testing.T) {
	manager := NewJWTManager("test-secret-key", time.Hour, 24*time.Hour)

	t.Run("验证有效令牌", func(t *testing.T) {
		token, _ := manager.GenerateAccessToken(100, "testuser", "test@example.com")

		claims, err := manager.ValidateToken(token)

		require.NoError(t, err, "ValidateToken() 应该成功")
		assert.Equal(t, uint(100), claims.UserID, "UserID 应该匹配")
	})

	t.Run("无效令牌格式", func(t *testing.T) {
		invalidTokens := []string{
			"",
			"invalid",
			"not.a.jwt",
			"three.parts.but.invalid",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.invalid.signature",
		}

		for _, token := range invalidTokens {
			_, err := manager.ValidateToken(token)
			assert.Error(t, err, "ValidateToken(%q) 应该返回错误", token)
		}
	})

	t.Run("错误签名密钥的令牌", func(t *testing.T) {
		otherManager := NewJWTManager("other-secret-key", time.Hour, 24*time.Hour)
		token, _ := otherManager.GenerateAccessToken(1, "user", "")

		_, err := manager.ValidateToken(token)

		assert.Error(t, err, "使用不同密钥签名的令牌应该验证失败")
	})

	t.Run("过期令牌", func(t *testing.T) {
		// 创建一个立即过期的管理器
		quickManager := NewJWTManager("secret", -time.Hour, 24*time.Hour)
		token, _ := quickManager.GenerateAccessToken(1, "user", "")

		_, err := manager.ValidateToken(token)

		assert.Error(t, err, "过期令牌应该验证失败")
	})
}

func TestJWTManager_GenerateTokenPair(t *testing.T) {
	manager := NewJWTManager("test-secret", time.Hour, 24*time.Hour)

	t.Run("成功生成令牌对", func(t *testing.T) {
		accessToken, refreshToken, err := manager.GenerateTokenPair(1, "user", "user@example.com")

		require.NoError(t, err, "GenerateTokenPair() 应该成功")
		assert.NotEmpty(t, accessToken, "accessToken 不应为空")
		assert.NotEmpty(t, refreshToken, "refreshToken 不应为空")
		assert.NotEqual(t, accessToken, refreshToken, "访问令牌和刷新令牌应该不同")
	})

	t.Run("令牌对可以独立验证", func(t *testing.T) {
		accessToken, refreshToken, _ := manager.GenerateTokenPair(123, "testuser", "test@example.com")

		accessClaims, err := manager.ValidateToken(accessToken)
		require.NoError(t, err, "访问令牌验证应该成功")
		assert.Equal(t, uint(123), accessClaims.UserID, "访问令牌 UserID 应该匹配")

		refreshClaims, err := manager.ValidateToken(refreshToken)
		require.NoError(t, err, "刷新令牌验证应该成功")
		assert.Equal(t, uint(123), refreshClaims.UserID, "刷新令牌 UserID 应该匹配")
	})
}

func TestJWTManager_EdgeCases(t *testing.T) {
	t.Run("空密钥", func(t *testing.T) {
		manager := NewJWTManager("", time.Hour, 24*time.Hour)
		token, err := manager.GenerateAccessToken(1, "user", "")

		assert.NoError(t, err, "空密钥也应该能生成令牌")

		// 但是空密钥的令牌安全性较低，验证应该成功
		_, err = manager.ValidateToken(token)
		assert.NoError(t, err, "空密钥令牌验证应该成功")
	})

	t.Run("非常短的过期时间", func(t *testing.T) {
		manager := NewJWTManager("secret", time.Nanosecond, time.Nanosecond)
		token, _ := manager.GenerateAccessToken(1, "user", "")

		// 等待令牌过期
		time.Sleep(time.Millisecond)

		_, err := manager.ValidateToken(token)
		assert.Error(t, err, "过期令牌应该验证失败")
	})

	t.Run("UserID 为 0", func(t *testing.T) {
		manager := NewJWTManager("secret", time.Hour, 24*time.Hour)
		token, err := manager.GenerateAccessToken(0, "user", "")

		assert.NoError(t, err, "UserID 为 0 也应该能生成令牌")

		claims, _ := manager.ValidateToken(token)
		assert.Equal(t, uint(0), claims.UserID, "UserID 应该是 0")
	})

	t.Run("特殊字符的用户名", func(t *testing.T) {
		manager := NewJWTManager("secret", time.Hour, 24*time.Hour)
		specialUsername := "user@domain.com<script>alert('xss')</script>"
		token, err := manager.GenerateAccessToken(1, specialUsername, "")

		assert.NoError(t, err, "特殊字符用户名应该能生成令牌")

		claims, _ := manager.ValidateToken(token)
		assert.Equal(t, specialUsername, claims.Username, "Username 应该匹配")
	})

	t.Run("Unicode 用户名", func(t *testing.T) {
		manager := NewJWTManager("secret", time.Hour, 24*time.Hour)
		unicodeUsername := "用户名🚀测试"
		token, err := manager.GenerateAccessToken(1, unicodeUsername, "")

		assert.NoError(t, err, "Unicode 用户名应该能生成令牌")

		claims, _ := manager.ValidateToken(token)
		assert.Equal(t, unicodeUsername, claims.Username, "Username 应该匹配")
	})
}

func BenchmarkJWTManager_GenerateAccessToken(b *testing.B) {
	manager := NewJWTManager("benchmark-secret-key", time.Hour, 24*time.Hour)

	b.ResetTimer()
	for range b.N {
		_, _ = manager.GenerateAccessToken(1, "benchuser", "bench@example.com")
	}
}

func BenchmarkJWTManager_ValidateToken(b *testing.B) {
	manager := NewJWTManager("benchmark-secret-key", time.Hour, 24*time.Hour)
	token, _ := manager.GenerateAccessToken(1, "benchuser", "bench@example.com")

	b.ResetTimer()
	for range b.N {
		_, _ = manager.ValidateToken(token)
	}
}

func BenchmarkJWTManager_GenerateTokenPair(b *testing.B) {
	manager := NewJWTManager("benchmark-secret-key", time.Hour, 24*time.Hour)

	b.ResetTimer()
	for range b.N {
		_, _, _ = manager.GenerateTokenPair(1, "benchuser", "bench@example.com")
	}
}
