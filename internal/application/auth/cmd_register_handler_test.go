package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	domainUser "github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
)

func TestRegisterHandler_Handle_Success(t *testing.T) {
	tests := []struct {
		name string
		cmd  RegisterCommand
	}{
		{
			name: "成功注册新用户",
			cmd: RegisterCommand{
				Username: "newuser",
				Email:    "newuser@example.com",
				Password: "Password123!",
				FullName: "New User",
			},
		},
		{
			name: "注册用户不填全名",
			cmd: RegisterCommand{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "SecurePass456",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockUserCmdRepo := new(MockUserCommandRepository)
			mockUserQryRepo := new(MockUserQueryRepository)
			mockAuthService := new(MockAuthService)

			expiresAt := time.Now().Add(24 * time.Hour)

			mockAuthService.On("ValidatePasswordPolicy", mock.Anything, tt.cmd.Password).Return(nil)
			mockUserQryRepo.On("ExistsByUsername", mock.Anything, tt.cmd.Username).Return(false, nil)
			mockUserQryRepo.On("ExistsByEmail", mock.Anything, tt.cmd.Email).Return(false, nil)
			mockAuthService.On("GeneratePasswordHash", mock.Anything, tt.cmd.Password).Return("hashed_password", nil)
			mockUserCmdRepo.On("Create", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)
			mockAuthService.On("GenerateAccessToken", mock.Anything, uint(1), tt.cmd.Username).Return("access_token", expiresAt, nil)
			mockAuthService.On("GenerateRefreshToken", mock.Anything, uint(1)).Return("refresh_token", expiresAt.Add(7*24*time.Hour), nil)

			handler := NewRegisterHandler(mockUserCmdRepo, mockUserQryRepo, mockAuthService)

			// Act
			result, err := handler.Handle(context.Background(), tt.cmd)

			// Assert
			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, uint(1), result.UserID)
			assert.Equal(t, tt.cmd.Username, result.Username)
			assert.Equal(t, tt.cmd.Email, result.Email)
			assert.Equal(t, "access_token", result.AccessToken)
			assert.Equal(t, "refresh_token", result.RefreshToken)
			assert.Equal(t, "Bearer", result.TokenType)

			mockUserCmdRepo.AssertExpectations(t)
			mockUserQryRepo.AssertExpectations(t)
			mockAuthService.AssertExpectations(t)
		})
	}
}

func TestRegisterHandler_Handle_Error(t *testing.T) {
	tests := []struct {
		name       string
		cmd        RegisterCommand
		setupMocks func(*MockUserCommandRepository, *MockUserQueryRepository, *MockAuthService)
		wantErr    string
	}{
		{
			name: "密码策略验证失败",
			cmd:  RegisterCommand{Username: "user", Email: "user@example.com", Password: "123"},
			setupMocks: func(cmdRepo *MockUserCommandRepository, qryRepo *MockUserQueryRepository, authService *MockAuthService) {
				authService.On("ValidatePasswordPolicy", mock.Anything, "123").Return(errors.New("password too short"))
			},
			wantErr: "password too short",
		},
		{
			name: "用户名已存在",
			cmd:  RegisterCommand{Username: "existing", Email: "new@example.com", Password: "ValidPass123"},
			setupMocks: func(cmdRepo *MockUserCommandRepository, qryRepo *MockUserQueryRepository, authService *MockAuthService) {
				authService.On("ValidatePasswordPolicy", mock.Anything, "ValidPass123").Return(nil)
				qryRepo.On("ExistsByUsername", mock.Anything, "existing").Return(true, nil)
			},
			wantErr: domainUser.ErrUsernameAlreadyExists.Error(),
		},
		{
			name: "邮箱已存在",
			cmd:  RegisterCommand{Username: "newuser", Email: "existing@example.com", Password: "ValidPass123"},
			setupMocks: func(cmdRepo *MockUserCommandRepository, qryRepo *MockUserQueryRepository, authService *MockAuthService) {
				authService.On("ValidatePasswordPolicy", mock.Anything, "ValidPass123").Return(nil)
				qryRepo.On("ExistsByUsername", mock.Anything, "newuser").Return(false, nil)
				qryRepo.On("ExistsByEmail", mock.Anything, "existing@example.com").Return(true, nil)
			},
			wantErr: domainUser.ErrEmailAlreadyExists.Error(),
		},
		{
			name: "检查用户名时数据库错误",
			cmd:  RegisterCommand{Username: "user", Email: "user@example.com", Password: "ValidPass123"},
			setupMocks: func(cmdRepo *MockUserCommandRepository, qryRepo *MockUserQueryRepository, authService *MockAuthService) {
				authService.On("ValidatePasswordPolicy", mock.Anything, "ValidPass123").Return(nil)
				qryRepo.On("ExistsByUsername", mock.Anything, "user").Return(false, errors.New("database error"))
			},
			wantErr: "failed to check username existence",
		},
		{
			name: "检查邮箱时数据库错误",
			cmd:  RegisterCommand{Username: "user", Email: "user@example.com", Password: "ValidPass123"},
			setupMocks: func(cmdRepo *MockUserCommandRepository, qryRepo *MockUserQueryRepository, authService *MockAuthService) {
				authService.On("ValidatePasswordPolicy", mock.Anything, "ValidPass123").Return(nil)
				qryRepo.On("ExistsByUsername", mock.Anything, "user").Return(false, nil)
				qryRepo.On("ExistsByEmail", mock.Anything, "user@example.com").Return(false, errors.New("database error"))
			},
			wantErr: "failed to check email existence",
		},
		{
			name: "密码哈希生成失败",
			cmd:  RegisterCommand{Username: "user", Email: "user@example.com", Password: "ValidPass123"},
			setupMocks: func(cmdRepo *MockUserCommandRepository, qryRepo *MockUserQueryRepository, authService *MockAuthService) {
				authService.On("ValidatePasswordPolicy", mock.Anything, "ValidPass123").Return(nil)
				qryRepo.On("ExistsByUsername", mock.Anything, "user").Return(false, nil)
				qryRepo.On("ExistsByEmail", mock.Anything, "user@example.com").Return(false, nil)
				authService.On("GeneratePasswordHash", mock.Anything, "ValidPass123").Return("", errors.New("hash error"))
			},
			wantErr: "failed to hash password",
		},
		{
			name: "创建用户失败",
			cmd:  RegisterCommand{Username: "user", Email: "user@example.com", Password: "ValidPass123"},
			setupMocks: func(cmdRepo *MockUserCommandRepository, qryRepo *MockUserQueryRepository, authService *MockAuthService) {
				authService.On("ValidatePasswordPolicy", mock.Anything, "ValidPass123").Return(nil)
				qryRepo.On("ExistsByUsername", mock.Anything, "user").Return(false, nil)
				qryRepo.On("ExistsByEmail", mock.Anything, "user@example.com").Return(false, nil)
				authService.On("GeneratePasswordHash", mock.Anything, "ValidPass123").Return("hashed", nil)
				cmdRepo.On("Create", mock.Anything, mock.AnythingOfType("*user.User")).Return(errors.New("create error"))
			},
			wantErr: "failed to create user",
		},
		{
			name: "生成访问令牌失败",
			cmd:  RegisterCommand{Username: "user", Email: "user@example.com", Password: "ValidPass123"},
			setupMocks: func(cmdRepo *MockUserCommandRepository, qryRepo *MockUserQueryRepository, authService *MockAuthService) {
				authService.On("ValidatePasswordPolicy", mock.Anything, "ValidPass123").Return(nil)
				qryRepo.On("ExistsByUsername", mock.Anything, "user").Return(false, nil)
				qryRepo.On("ExistsByEmail", mock.Anything, "user@example.com").Return(false, nil)
				authService.On("GeneratePasswordHash", mock.Anything, "ValidPass123").Return("hashed", nil)
				cmdRepo.On("Create", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)
				authService.On("GenerateAccessToken", mock.Anything, uint(1), "user").Return("", time.Time{}, errors.New("token error"))
			},
			wantErr: "failed to generate access token",
		},
		{
			name: "生成刷新令牌失败",
			cmd:  RegisterCommand{Username: "user", Email: "user@example.com", Password: "ValidPass123"},
			setupMocks: func(cmdRepo *MockUserCommandRepository, qryRepo *MockUserQueryRepository, authService *MockAuthService) {
				authService.On("ValidatePasswordPolicy", mock.Anything, "ValidPass123").Return(nil)
				qryRepo.On("ExistsByUsername", mock.Anything, "user").Return(false, nil)
				qryRepo.On("ExistsByEmail", mock.Anything, "user@example.com").Return(false, nil)
				authService.On("GeneratePasswordHash", mock.Anything, "ValidPass123").Return("hashed", nil)
				cmdRepo.On("Create", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)
				authService.On("GenerateAccessToken", mock.Anything, uint(1), "user").Return("access_token", time.Now().Add(time.Hour), nil)
				authService.On("GenerateRefreshToken", mock.Anything, uint(1)).Return("", time.Time{}, errors.New("refresh token error"))
			},
			wantErr: "failed to generate refresh token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockUserCmdRepo := new(MockUserCommandRepository)
			mockUserQryRepo := new(MockUserQueryRepository)
			mockAuthService := new(MockAuthService)
			tt.setupMocks(mockUserCmdRepo, mockUserQryRepo, mockAuthService)

			handler := NewRegisterHandler(mockUserCmdRepo, mockUserQryRepo, mockAuthService)

			// Act
			result, err := handler.Handle(context.Background(), tt.cmd)

			// Assert
			assert.Nil(t, result)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}
