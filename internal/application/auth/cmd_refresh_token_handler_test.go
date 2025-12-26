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
	domainUser "github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
)

func TestRefreshTokenHandler_Handle_Success(t *testing.T) {
	// Arrange
	mockUserQryRepo := new(MockUserQueryRepository)
	mockAuthService := new(MockAuthService)

	expiresAt := time.Now().Add(24 * time.Hour)
	refreshExpiresAt := time.Now().Add(7 * 24 * time.Hour)

	user := &domainUser.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Status:   "active",
	}

	mockAuthService.On("ValidateRefreshToken", mock.Anything, "valid_refresh_token").Return(uint(1), nil)
	mockUserQryRepo.On("GetByIDWithRoles", mock.Anything, uint(1)).Return(user, nil)
	mockAuthService.On("GenerateAccessToken", mock.Anything, uint(1), "testuser").Return("new_access_token", expiresAt, nil)
	mockAuthService.On("GenerateRefreshToken", mock.Anything, uint(1)).Return("new_refresh_token", refreshExpiresAt, nil)

	handler := NewRefreshTokenHandler(mockUserQryRepo, mockAuthService)

	// Act
	result, err := handler.Handle(context.Background(), RefreshTokenCommand{
		RefreshToken: "valid_refresh_token",
	})

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "new_access_token", result.AccessToken)
	assert.Equal(t, "new_refresh_token", result.RefreshToken)
	assert.Equal(t, "Bearer", result.TokenType)

	mockUserQryRepo.AssertExpectations(t)
	mockAuthService.AssertExpectations(t)
}

func TestRefreshTokenHandler_Handle_Error(t *testing.T) {
	tests := []struct {
		name       string
		cmd        RefreshTokenCommand
		setupMocks func(*MockUserQueryRepository, *MockAuthService)
		wantErr    string
	}{
		{
			name: "无效的刷新令牌",
			cmd:  RefreshTokenCommand{RefreshToken: "invalid_token"},
			setupMocks: func(qryRepo *MockUserQueryRepository, authService *MockAuthService) {
				authService.On("ValidateRefreshToken", mock.Anything, "invalid_token").Return(uint(0), errors.New("invalid token"))
			},
			wantErr: "invalid token",
		},
		{
			name: "用户不存在",
			cmd:  RefreshTokenCommand{RefreshToken: "valid_token"},
			setupMocks: func(qryRepo *MockUserQueryRepository, authService *MockAuthService) {
				authService.On("ValidateRefreshToken", mock.Anything, "valid_token").Return(uint(999), nil)
				qryRepo.On("GetByIDWithRoles", mock.Anything, uint(999)).Return(nil, errors.New("user not found"))
			},
			wantErr: domainAuth.ErrUserNotFound.Error(),
		},
		{
			name: "用户已被禁用",
			cmd:  RefreshTokenCommand{RefreshToken: "valid_token"},
			setupMocks: func(qryRepo *MockUserQueryRepository, authService *MockAuthService) {
				authService.On("ValidateRefreshToken", mock.Anything, "valid_token").Return(uint(1), nil)
				qryRepo.On("GetByIDWithRoles", mock.Anything, uint(1)).Return(&domainUser.User{
					ID:       1,
					Username: "banned",
					Status:   "banned",
				}, nil)
			},
			wantErr: domainAuth.ErrUserBanned.Error(),
		},
		{
			name: "用户未激活",
			cmd:  RefreshTokenCommand{RefreshToken: "valid_token"},
			setupMocks: func(qryRepo *MockUserQueryRepository, authService *MockAuthService) {
				authService.On("ValidateRefreshToken", mock.Anything, "valid_token").Return(uint(1), nil)
				qryRepo.On("GetByIDWithRoles", mock.Anything, uint(1)).Return(&domainUser.User{
					ID:       1,
					Username: "inactive",
					Status:   "inactive",
				}, nil)
			},
			wantErr: domainAuth.ErrUserInactive.Error(),
		},
		{
			name: "生成访问令牌失败",
			cmd:  RefreshTokenCommand{RefreshToken: "valid_token"},
			setupMocks: func(qryRepo *MockUserQueryRepository, authService *MockAuthService) {
				authService.On("ValidateRefreshToken", mock.Anything, "valid_token").Return(uint(1), nil)
				qryRepo.On("GetByIDWithRoles", mock.Anything, uint(1)).Return(&domainUser.User{
					ID:       1,
					Username: "testuser",
					Status:   "active",
				}, nil)
				authService.On("GenerateAccessToken", mock.Anything, uint(1), "testuser").Return("", time.Time{}, errors.New("token error"))
			},
			wantErr: "failed to generate access token",
		},
		{
			name: "生成刷新令牌失败",
			cmd:  RefreshTokenCommand{RefreshToken: "valid_token"},
			setupMocks: func(qryRepo *MockUserQueryRepository, authService *MockAuthService) {
				authService.On("ValidateRefreshToken", mock.Anything, "valid_token").Return(uint(1), nil)
				qryRepo.On("GetByIDWithRoles", mock.Anything, uint(1)).Return(&domainUser.User{
					ID:       1,
					Username: "testuser",
					Status:   "active",
				}, nil)
				authService.On("GenerateAccessToken", mock.Anything, uint(1), "testuser").Return("access_token", time.Now().Add(time.Hour), nil)
				authService.On("GenerateRefreshToken", mock.Anything, uint(1)).Return("", time.Time{}, errors.New("refresh error"))
			},
			wantErr: "failed to generate refresh token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockUserQryRepo := new(MockUserQueryRepository)
			mockAuthService := new(MockAuthService)
			tt.setupMocks(mockUserQryRepo, mockAuthService)

			handler := NewRefreshTokenHandler(mockUserQryRepo, mockAuthService)

			// Act
			result, err := handler.Handle(context.Background(), tt.cmd)

			// Assert
			assert.Nil(t, result)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}
