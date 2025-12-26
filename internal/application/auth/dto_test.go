package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================
// ToLoginResponse 测试
// ============================================================

func TestLoginResultDTO_ToLoginResponse(t *testing.T) {
	tests := []struct {
		name     string
		input    *LoginResultDTO
		expected *LoginResponseDTO
	}{
		{
			name: "完整登录结果转换",
			input: &LoginResultDTO{
				AccessToken:  "access_token_123",
				RefreshToken: "refresh_token_456",
				TokenType:    "Bearer",
				ExpiresIn:    3600,
				UserID:       1,
				Username:     "testuser",
				Requires2FA:  false,
				SessionToken: "",
			},
			expected: &LoginResponseDTO{
				AccessToken:  "access_token_123",
				RefreshToken: "refresh_token_456",
				TokenType:    "Bearer",
				ExpiresIn:    3600,
				User: UserBriefDTO{
					UserID:   1,
					Username: "testuser",
				},
				Requires2FA:  false,
				SessionToken: "",
			},
		},
		{
			name: "需要 2FA 的登录结果转换",
			input: &LoginResultDTO{
				AccessToken:  "",
				RefreshToken: "",
				TokenType:    "",
				ExpiresIn:    0,
				UserID:       2,
				Username:     "admin",
				Requires2FA:  true,
				SessionToken: "session_token_789",
			},
			expected: &LoginResponseDTO{
				AccessToken:  "",
				RefreshToken: "",
				TokenType:    "",
				ExpiresIn:    0,
				User: UserBriefDTO{
					UserID:   2,
					Username: "admin",
				},
				Requires2FA:  true,
				SessionToken: "session_token_789",
			},
		},
		{
			name: "空值处理",
			input: &LoginResultDTO{
				UserID:   0,
				Username: "",
			},
			expected: &LoginResponseDTO{
				User: UserBriefDTO{
					UserID:   0,
					Username: "",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.ToLoginResponse()

			assert.Equal(t, tt.expected.AccessToken, result.AccessToken)
			assert.Equal(t, tt.expected.RefreshToken, result.RefreshToken)
			assert.Equal(t, tt.expected.TokenType, result.TokenType)
			assert.Equal(t, tt.expected.ExpiresIn, result.ExpiresIn)
			assert.Equal(t, tt.expected.User.UserID, result.User.UserID)
			assert.Equal(t, tt.expected.User.Username, result.User.Username)
			assert.Equal(t, tt.expected.Requires2FA, result.Requires2FA)
			assert.Equal(t, tt.expected.SessionToken, result.SessionToken)
		})
	}
}

// ============================================================
// ToTwoFARequiredDTO 测试 (from mapper.go)
// ============================================================

func TestToTwoFARequiredDTO(t *testing.T) {
	tests := []struct {
		name         string
		sessionToken string
	}{
		{
			name:         "正常 session token",
			sessionToken: "valid_session_token_123",
		},
		{
			name:         "较长的 session token",
			sessionToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ",
		},
		{
			name:         "空 session token",
			sessionToken: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToTwoFARequiredDTO(tt.sessionToken)

			assert.NotNil(t, result)
			assert.True(t, result.Requires2FA, "Requires2FA 应始终为 true")
			assert.Equal(t, tt.sessionToken, result.SessionToken)
		})
	}
}

// ============================================================
// DTO 结构体字段验证测试
// ============================================================

func TestLoginDTO_Structure(t *testing.T) {
	dto := LoginDTO{
		Account:   "testuser",
		Password:  "password123",
		CaptchaID: "captcha-123",
		Captcha:   "1234",
	}

	assert.Equal(t, "testuser", dto.Account)
	assert.Equal(t, "password123", dto.Password)
	assert.Equal(t, "captcha-123", dto.CaptchaID)
	assert.Equal(t, "1234", dto.Captcha)
}

func TestLogin2FADTO_Structure(t *testing.T) {
	dto := Login2FADTO{
		SessionToken:  "session_123",
		TwoFactorCode: "123456",
	}

	assert.Equal(t, "session_123", dto.SessionToken)
	assert.Equal(t, "123456", dto.TwoFactorCode)
}

func TestRegisterDTO_Structure(t *testing.T) {
	dto := RegisterDTO{
		Username: "newuser",
		Email:    "newuser@example.com",
		Password: "securepassword",
		FullName: "New User",
	}

	assert.Equal(t, "newuser", dto.Username)
	assert.Equal(t, "newuser@example.com", dto.Email)
	assert.Equal(t, "securepassword", dto.Password)
	assert.Equal(t, "New User", dto.FullName)
}

func TestRefreshTokenDTO_Structure(t *testing.T) {
	dto := RefreshTokenDTO{
		RefreshToken: "refresh_token_xyz",
	}

	assert.Equal(t, "refresh_token_xyz", dto.RefreshToken)
}

func TestTokenDTO_Structure(t *testing.T) {
	dto := TokenDTO{
		AccessToken:  "access",
		RefreshToken: "refresh",
		TokenType:    "Bearer",
		ExpiresIn:    7200,
	}

	assert.Equal(t, "access", dto.AccessToken)
	assert.Equal(t, "refresh", dto.RefreshToken)
	assert.Equal(t, "Bearer", dto.TokenType)
	assert.Equal(t, 7200, dto.ExpiresIn)
}

func TestRefreshTokenResultDTO_Structure(t *testing.T) {
	dto := RefreshTokenResultDTO{
		AccessToken:  "new_access",
		RefreshToken: "new_refresh",
		TokenType:    "Bearer",
		ExpiresIn:    3600,
	}

	assert.Equal(t, "new_access", dto.AccessToken)
	assert.Equal(t, "new_refresh", dto.RefreshToken)
	assert.Equal(t, "Bearer", dto.TokenType)
	assert.Equal(t, 3600, dto.ExpiresIn)
}

func TestRegisterResultDTO_Structure(t *testing.T) {
	dto := RegisterResultDTO{
		UserID:       10,
		Username:     "registered",
		Email:        "registered@example.com",
		AccessToken:  "access_token",
		RefreshToken: "refresh_token",
		TokenType:    "Bearer",
		ExpiresIn:    3600,
	}

	assert.Equal(t, uint(10), dto.UserID)
	assert.Equal(t, "registered", dto.Username)
	assert.Equal(t, "registered@example.com", dto.Email)
	assert.Equal(t, "access_token", dto.AccessToken)
	assert.Equal(t, "refresh_token", dto.RefreshToken)
	assert.Equal(t, "Bearer", dto.TokenType)
	assert.Equal(t, 3600, dto.ExpiresIn)
}

func TestUserBriefDTO_Structure(t *testing.T) {
	dto := UserBriefDTO{
		UserID:   5,
		Username: "brief_user",
	}

	assert.Equal(t, uint(5), dto.UserID)
	assert.Equal(t, "brief_user", dto.Username)
}

func TestTwoFARequiredDTO_Structure(t *testing.T) {
	dto := TwoFARequiredDTO{
		Requires2FA:  true,
		SessionToken: "session_abc",
	}

	assert.True(t, dto.Requires2FA)
	assert.Equal(t, "session_abc", dto.SessionToken)
}

// ============================================================
// 重新导出的错误变量测试
// ============================================================

func TestReExportedErrors(t *testing.T) {
	// 验证重新导出的错误变量不为 nil
	require.Error(t, ErrInvalidToken)
	require.Error(t, ErrTokenExpired)
}
