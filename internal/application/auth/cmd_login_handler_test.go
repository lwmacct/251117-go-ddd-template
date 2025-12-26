package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	domainAuth "github.com/lwmacct/251117-go-ddd-template/internal/domain/auth"
	domainTwoFA "github.com/lwmacct/251117-go-ddd-template/internal/domain/twofa"
	domainUser "github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
	authInfra "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/auth"
)

func TestLoginHandler_Handle_Success_Without2FA(t *testing.T) {
	// Arrange
	mockUserQryRepo := new(MockUserQueryRepository)
	mockCaptchaRepo := new(MockCaptchaCommandRepository)
	mockTwofaQryRepo := new(MockTwoFAQueryRepository)
	mockAuthService := new(MockAuthService)
	loginSession := authInfra.NewLoginSessionService()

	expiresAt := time.Now().Add(24 * time.Hour)

	user := &domainUser.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashed_password",
		Status:   "active",
	}

	mockCaptchaRepo.On("Verify", mock.Anything, "captcha_id", "captcha_code").Return(true, nil)
	mockUserQryRepo.On("GetByUsernameWithRoles", mock.Anything, "testuser").Return(user, nil)
	mockAuthService.On("VerifyPassword", mock.Anything, "hashed_password", "password123").Return(nil)
	mockTwofaQryRepo.On("FindByUserID", mock.Anything, uint(1)).Return(nil, nil) // 2FA 未启用
	mockAuthService.On("GenerateAccessToken", mock.Anything, uint(1), "testuser").Return("access_token", expiresAt, nil)
	mockAuthService.On("GenerateRefreshToken", mock.Anything, uint(1)).Return("refresh_token", expiresAt.Add(7*24*time.Hour), nil)

	handler := NewLoginHandler(mockUserQryRepo, mockCaptchaRepo, mockTwofaQryRepo, mockAuthService, loginSession, nil)

	// Act
	result, err := handler.Handle(context.Background(), LoginCommand{
		Account:   "testuser",
		Password:  "password123",
		CaptchaID: "captcha_id",
		Captcha:   "captcha_code",
		ClientIP:  "127.0.0.1",
		UserAgent: "TestAgent/1.0",
	})

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "access_token", result.AccessToken)
	assert.Equal(t, "refresh_token", result.RefreshToken)
	assert.Equal(t, "Bearer", result.TokenType)
	assert.Equal(t, uint(1), result.UserID)
	assert.Equal(t, "testuser", result.Username)
	assert.False(t, result.Requires2FA)

	mockCaptchaRepo.AssertExpectations(t)
	mockUserQryRepo.AssertExpectations(t)
	mockTwofaQryRepo.AssertExpectations(t)
	mockAuthService.AssertExpectations(t)
}

func TestLoginHandler_Handle_Success_With2FA(t *testing.T) {
	// Arrange
	mockUserQryRepo := new(MockUserQueryRepository)
	mockCaptchaRepo := new(MockCaptchaCommandRepository)
	mockTwofaQryRepo := new(MockTwoFAQueryRepository)
	mockAuthService := new(MockAuthService)
	loginSession := authInfra.NewLoginSessionService()

	user := &domainUser.User{
		ID:       1,
		Username: "admin",
		Email:    "admin@example.com",
		Password: "hashed_password",
		Status:   "active",
	}

	twofa := &domainTwoFA.TwoFA{
		ID:      1,
		UserID:  1,
		Enabled: true,
		Secret:  "ABCDEFGHIJ",
	}

	mockCaptchaRepo.On("Verify", mock.Anything, "captcha_id", "captcha_code").Return(true, nil)
	mockUserQryRepo.On("GetByUsernameWithRoles", mock.Anything, "admin").Return(user, nil)
	mockAuthService.On("VerifyPassword", mock.Anything, "hashed_password", "password123").Return(nil)
	mockTwofaQryRepo.On("FindByUserID", mock.Anything, uint(1)).Return(twofa, nil) // 2FA 已启用

	handler := NewLoginHandler(mockUserQryRepo, mockCaptchaRepo, mockTwofaQryRepo, mockAuthService, loginSession, nil)

	// Act
	result, err := handler.Handle(context.Background(), LoginCommand{
		Account:   "admin",
		Password:  "password123",
		CaptchaID: "captcha_id",
		Captcha:   "captcha_code",
		ClientIP:  "127.0.0.1",
		UserAgent: "TestAgent/1.0",
	})

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Requires2FA)
	assert.NotEmpty(t, result.SessionToken) // 应生成 session token
	assert.Equal(t, uint(1), result.UserID)
	assert.Equal(t, "admin", result.Username)
	assert.Empty(t, result.AccessToken)  // 不应生成访问令牌
	assert.Empty(t, result.RefreshToken) // 不应生成刷新令牌

	mockCaptchaRepo.AssertExpectations(t)
	mockUserQryRepo.AssertExpectations(t)
	mockTwofaQryRepo.AssertExpectations(t)
	mockAuthService.AssertExpectations(t)
}

func TestLoginHandler_Handle_LoginByEmail(t *testing.T) {
	// Arrange
	mockUserQryRepo := new(MockUserQueryRepository)
	mockCaptchaRepo := new(MockCaptchaCommandRepository)
	mockTwofaQryRepo := new(MockTwoFAQueryRepository)
	mockAuthService := new(MockAuthService)
	loginSession := authInfra.NewLoginSessionService()

	expiresAt := time.Now().Add(24 * time.Hour)

	user := &domainUser.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashed_password",
		Status:   "active",
	}

	mockCaptchaRepo.On("Verify", mock.Anything, "captcha_id", "code").Return(true, nil)
	// 用户名查找失败
	mockUserQryRepo.On("GetByUsernameWithRoles", mock.Anything, "test@example.com").Return(nil, errors.New("not found"))
	// 通过邮箱查找成功
	mockUserQryRepo.On("GetByEmailWithRoles", mock.Anything, "test@example.com").Return(user, nil)
	mockAuthService.On("VerifyPassword", mock.Anything, "hashed_password", "password").Return(nil)
	mockTwofaQryRepo.On("FindByUserID", mock.Anything, uint(1)).Return(nil, nil)
	mockAuthService.On("GenerateAccessToken", mock.Anything, uint(1), "testuser").Return("token", expiresAt, nil)
	mockAuthService.On("GenerateRefreshToken", mock.Anything, uint(1)).Return("refresh", expiresAt, nil)

	handler := NewLoginHandler(mockUserQryRepo, mockCaptchaRepo, mockTwofaQryRepo, mockAuthService, loginSession, nil)

	// Act
	result, err := handler.Handle(context.Background(), LoginCommand{
		Account:   "test@example.com",
		Password:  "password",
		CaptchaID: "captcha_id",
		Captcha:   "code",
	})

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "testuser", result.Username)
}

func TestLoginHandler_Handle_Error(t *testing.T) {
	tests := []struct {
		name       string
		cmd        LoginCommand
		setupMocks func(*MockUserQueryRepository, *MockCaptchaCommandRepository, *MockTwoFAQueryRepository, *MockAuthService)
		wantErr    error
	}{
		{
			name: "验证码验证失败",
			cmd:  LoginCommand{Account: "user", Password: "pass", CaptchaID: "id", Captcha: "wrong"},
			setupMocks: func(userQry *MockUserQueryRepository, captcha *MockCaptchaCommandRepository, twofa *MockTwoFAQueryRepository, auth *MockAuthService) {
				captcha.On("Verify", mock.Anything, "id", "wrong").Return(false, nil)
			},
			wantErr: domainAuth.ErrInvalidCaptcha,
		},
		{
			name: "验证码服务错误",
			cmd:  LoginCommand{Account: "user", Password: "pass", CaptchaID: "id", Captcha: "code"},
			setupMocks: func(userQry *MockUserQueryRepository, captcha *MockCaptchaCommandRepository, twofa *MockTwoFAQueryRepository, auth *MockAuthService) {
				captcha.On("Verify", mock.Anything, "id", "code").Return(false, errors.New("captcha service error"))
			},
			wantErr: errors.New("failed to verify captcha"),
		},
		{
			name: "用户不存在",
			cmd:  LoginCommand{Account: "unknown", Password: "pass", CaptchaID: "id", Captcha: "code"},
			setupMocks: func(userQry *MockUserQueryRepository, captcha *MockCaptchaCommandRepository, twofa *MockTwoFAQueryRepository, auth *MockAuthService) {
				captcha.On("Verify", mock.Anything, "id", "code").Return(true, nil)
				userQry.On("GetByUsernameWithRoles", mock.Anything, "unknown").Return(nil, errors.New("not found"))
				userQry.On("GetByEmailWithRoles", mock.Anything, "unknown").Return(nil, errors.New("not found"))
			},
			wantErr: domainAuth.ErrInvalidCredentials,
		},
		{
			name: "用户已被禁用",
			cmd:  LoginCommand{Account: "banned", Password: "pass", CaptchaID: "id", Captcha: "code"},
			setupMocks: func(userQry *MockUserQueryRepository, captcha *MockCaptchaCommandRepository, twofa *MockTwoFAQueryRepository, auth *MockAuthService) {
				captcha.On("Verify", mock.Anything, "id", "code").Return(true, nil)
				userQry.On("GetByUsernameWithRoles", mock.Anything, "banned").Return(&domainUser.User{
					ID:       1,
					Username: "banned",
					Status:   "banned",
					Password: "hash",
				}, nil)
			},
			wantErr: domainAuth.ErrUserBanned,
		},
		{
			name: "用户未激活",
			cmd:  LoginCommand{Account: "inactive", Password: "pass", CaptchaID: "id", Captcha: "code"},
			setupMocks: func(userQry *MockUserQueryRepository, captcha *MockCaptchaCommandRepository, twofa *MockTwoFAQueryRepository, auth *MockAuthService) {
				captcha.On("Verify", mock.Anything, "id", "code").Return(true, nil)
				userQry.On("GetByUsernameWithRoles", mock.Anything, "inactive").Return(&domainUser.User{
					ID:       1,
					Username: "inactive",
					Status:   "inactive",
					Password: "hash",
				}, nil)
			},
			wantErr: domainAuth.ErrUserInactive,
		},
		{
			name: "密码错误",
			cmd:  LoginCommand{Account: "user", Password: "wrongpass", CaptchaID: "id", Captcha: "code"},
			setupMocks: func(userQry *MockUserQueryRepository, captcha *MockCaptchaCommandRepository, twofa *MockTwoFAQueryRepository, auth *MockAuthService) {
				captcha.On("Verify", mock.Anything, "id", "code").Return(true, nil)
				userQry.On("GetByUsernameWithRoles", mock.Anything, "user").Return(&domainUser.User{
					ID:       1,
					Username: "user",
					Status:   "active",
					Password: "hashed",
				}, nil)
				auth.On("VerifyPassword", mock.Anything, "hashed", "wrongpass").Return(errors.New("password mismatch"))
			},
			wantErr: domainAuth.ErrInvalidCredentials,
		},
		{
			name: "生成访问令牌失败",
			cmd:  LoginCommand{Account: "user", Password: "pass", CaptchaID: "id", Captcha: "code"},
			setupMocks: func(userQry *MockUserQueryRepository, captcha *MockCaptchaCommandRepository, twofa *MockTwoFAQueryRepository, auth *MockAuthService) {
				captcha.On("Verify", mock.Anything, "id", "code").Return(true, nil)
				userQry.On("GetByUsernameWithRoles", mock.Anything, "user").Return(&domainUser.User{
					ID:       1,
					Username: "user",
					Status:   "active",
					Password: "hashed",
				}, nil)
				auth.On("VerifyPassword", mock.Anything, "hashed", "pass").Return(nil)
				twofa.On("FindByUserID", mock.Anything, uint(1)).Return(nil, nil)
				auth.On("GenerateAccessToken", mock.Anything, uint(1), "user").Return("", time.Time{}, errors.New("token error"))
			},
			wantErr: errors.New("failed to generate access token"),
		},
		{
			name: "生成刷新令牌失败",
			cmd:  LoginCommand{Account: "user", Password: "pass", CaptchaID: "id", Captcha: "code"},
			setupMocks: func(userQry *MockUserQueryRepository, captcha *MockCaptchaCommandRepository, twofa *MockTwoFAQueryRepository, auth *MockAuthService) {
				captcha.On("Verify", mock.Anything, "id", "code").Return(true, nil)
				userQry.On("GetByUsernameWithRoles", mock.Anything, "user").Return(&domainUser.User{
					ID:       1,
					Username: "user",
					Status:   "active",
					Password: "hashed",
				}, nil)
				auth.On("VerifyPassword", mock.Anything, "hashed", "pass").Return(nil)
				twofa.On("FindByUserID", mock.Anything, uint(1)).Return(nil, nil)
				auth.On("GenerateAccessToken", mock.Anything, uint(1), "user").Return("token", time.Now().Add(time.Hour), nil)
				auth.On("GenerateRefreshToken", mock.Anything, uint(1)).Return("", time.Time{}, errors.New("refresh error"))
			},
			wantErr: errors.New("failed to generate refresh token"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockUserQryRepo := new(MockUserQueryRepository)
			mockCaptchaRepo := new(MockCaptchaCommandRepository)
			mockTwofaQryRepo := new(MockTwoFAQueryRepository)
			mockAuthService := new(MockAuthService)
			loginSession := authInfra.NewLoginSessionService()

			tt.setupMocks(mockUserQryRepo, mockCaptchaRepo, mockTwofaQryRepo, mockAuthService)

			handler := NewLoginHandler(mockUserQryRepo, mockCaptchaRepo, mockTwofaQryRepo, mockAuthService, loginSession, nil)

			// Act
			result, err := handler.Handle(context.Background(), tt.cmd)

			// Assert
			assert.Nil(t, result)
			require.Error(t, err)
			if errors.Is(tt.wantErr, domainAuth.ErrInvalidCaptcha) ||
				errors.Is(tt.wantErr, domainAuth.ErrInvalidCredentials) ||
				errors.Is(tt.wantErr, domainAuth.ErrUserBanned) ||
				errors.Is(tt.wantErr, domainAuth.ErrUserInactive) {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.Contains(t, err.Error(), tt.wantErr.Error())
			}
		})
	}
}

func TestLoginHandler_Handle_NilAuditLogHandler(t *testing.T) {
	// 测试 auditLogHandler 为 nil 时不会 panic
	mockUserQryRepo := new(MockUserQueryRepository)
	mockCaptchaRepo := new(MockCaptchaCommandRepository)
	mockTwofaQryRepo := new(MockTwoFAQueryRepository)
	mockAuthService := new(MockAuthService)
	loginSession := authInfra.NewLoginSessionService()

	// 模拟验证码无效（会触发 logLoginEvent）
	mockCaptchaRepo.On("Verify", mock.Anything, "id", "wrong").Return(false, nil)

	handler := NewLoginHandler(mockUserQryRepo, mockCaptchaRepo, mockTwofaQryRepo, mockAuthService, loginSession, nil)

	// Act - 不应 panic
	_, err := handler.Handle(context.Background(), LoginCommand{
		Account:   "user",
		Password:  "pass",
		CaptchaID: "id",
		Captcha:   "wrong",
	})

	// Assert
	require.Error(t, err)
	assert.ErrorIs(t, err, domainAuth.ErrInvalidCaptcha)
}
