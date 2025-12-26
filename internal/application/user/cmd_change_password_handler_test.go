package user

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
)

func TestChangePasswordHandler_Handle_Success(t *testing.T) {
	tests := []struct {
		name string
		cmd  ChangePasswordCommand
	}{
		{
			name: "成功修改密码",
			cmd: ChangePasswordCommand{
				UserID:      1,
				OldPassword: "oldpass123",
				NewPassword: "newpass456",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockCmdRepo := new(MockUserCommandRepository)
			mockQryRepo := new(MockUserQueryRepository)
			mockAuthService := new(MockAuthService)

			mockQryRepo.On("GetByID", mock.Anything, tt.cmd.UserID).Return(&user.User{
				ID:       tt.cmd.UserID,
				Password: "hashed_old_password",
			}, nil)
			mockAuthService.On("VerifyPassword", mock.Anything, "hashed_old_password", tt.cmd.OldPassword).Return(nil)
			mockAuthService.On("ValidatePasswordPolicy", mock.Anything, tt.cmd.NewPassword).Return(nil)
			mockAuthService.On("GeneratePasswordHash", mock.Anything, tt.cmd.NewPassword).Return("hashed_new_password", nil)
			mockCmdRepo.On("UpdatePassword", mock.Anything, tt.cmd.UserID, "hashed_new_password").Return(nil)

			handler := NewChangePasswordHandler(mockCmdRepo, mockQryRepo, mockAuthService)

			// Act
			err := handler.Handle(context.Background(), tt.cmd)

			// Assert
			require.NoError(t, err)
			mockCmdRepo.AssertExpectations(t)
			mockQryRepo.AssertExpectations(t)
			mockAuthService.AssertExpectations(t)
		})
	}
}

func TestChangePasswordHandler_Handle_Error(t *testing.T) {
	tests := []struct {
		name       string
		cmd        ChangePasswordCommand
		setupMocks func(*MockUserCommandRepository, *MockUserQueryRepository, *MockAuthService)
		wantErr    string
	}{
		{
			name: "用户不存在",
			cmd:  ChangePasswordCommand{UserID: 999, OldPassword: "old", NewPassword: "new"},
			setupMocks: func(cmdRepo *MockUserCommandRepository, qryRepo *MockUserQueryRepository, authService *MockAuthService) {
				qryRepo.On("GetByID", mock.Anything, uint(999)).Return(nil, errors.New("user not found"))
			},
			wantErr: "user not found",
		},
		{
			name: "旧密码错误",
			cmd:  ChangePasswordCommand{UserID: 1, OldPassword: "wrongpass", NewPassword: "newpass"},
			setupMocks: func(cmdRepo *MockUserCommandRepository, qryRepo *MockUserQueryRepository, authService *MockAuthService) {
				qryRepo.On("GetByID", mock.Anything, uint(1)).Return(&user.User{
					ID:       1,
					Password: "hashed_password",
				}, nil)
				authService.On("VerifyPassword", mock.Anything, "hashed_password", "wrongpass").Return(errors.New("password mismatch"))
			},
			wantErr: "invalid password",
		},
		{
			name: "新密码不符合策略",
			cmd:  ChangePasswordCommand{UserID: 1, OldPassword: "oldpass", NewPassword: "123"},
			setupMocks: func(cmdRepo *MockUserCommandRepository, qryRepo *MockUserQueryRepository, authService *MockAuthService) {
				qryRepo.On("GetByID", mock.Anything, uint(1)).Return(&user.User{
					ID:       1,
					Password: "hashed_password",
				}, nil)
				authService.On("VerifyPassword", mock.Anything, "hashed_password", "oldpass").Return(nil)
				authService.On("ValidatePasswordPolicy", mock.Anything, "123").Return(errors.New("password too short"))
			},
			wantErr: "password too short",
		},
		{
			name: "密码哈希失败",
			cmd:  ChangePasswordCommand{UserID: 1, OldPassword: "oldpass", NewPassword: "newpass"},
			setupMocks: func(cmdRepo *MockUserCommandRepository, qryRepo *MockUserQueryRepository, authService *MockAuthService) {
				qryRepo.On("GetByID", mock.Anything, uint(1)).Return(&user.User{
					ID:       1,
					Password: "hashed_password",
				}, nil)
				authService.On("VerifyPassword", mock.Anything, "hashed_password", "oldpass").Return(nil)
				authService.On("ValidatePasswordPolicy", mock.Anything, "newpass").Return(nil)
				authService.On("GeneratePasswordHash", mock.Anything, "newpass").Return("", errors.New("hash failed"))
			},
			wantErr: "failed to hash password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockCmdRepo := new(MockUserCommandRepository)
			mockQryRepo := new(MockUserQueryRepository)
			mockAuthService := new(MockAuthService)
			tt.setupMocks(mockCmdRepo, mockQryRepo, mockAuthService)

			handler := NewChangePasswordHandler(mockCmdRepo, mockQryRepo, mockAuthService)

			// Act
			err := handler.Handle(context.Background(), tt.cmd)

			// Assert
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}
