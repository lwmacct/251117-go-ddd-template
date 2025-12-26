package twofa

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDisableHandler_Handle_Success(t *testing.T) {
	tests := []struct {
		name string
		cmd  DisableCommand
	}{
		{
			name: "成功禁用 2FA",
			cmd:  DisableCommand{UserID: 1},
		},
		{
			name: "禁用另一用户的 2FA",
			cmd:  DisableCommand{UserID: 100},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockService := new(MockTwoFAService)
			mockService.On("Disable", mock.Anything, tt.cmd.UserID).Return(nil)

			handler := NewDisableHandler(mockService)

			// Act
			err := handler.Handle(context.Background(), tt.cmd)

			// Assert
			require.NoError(t, err)
			mockService.AssertExpectations(t)
		})
	}
}

func TestDisableHandler_Handle_Error(t *testing.T) {
	tests := []struct {
		name    string
		cmd     DisableCommand
		err     error
		wantErr string
	}{
		{
			name:    "用户不存在",
			cmd:     DisableCommand{UserID: 999},
			err:     errors.New("user not found"),
			wantErr: "user not found",
		},
		{
			name:    "2FA 未启用",
			cmd:     DisableCommand{UserID: 1},
			err:     errors.New("2FA not enabled"),
			wantErr: "2FA not enabled",
		},
		{
			name:    "数据库错误",
			cmd:     DisableCommand{UserID: 2},
			err:     errors.New("database connection failed"),
			wantErr: "database connection failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockService := new(MockTwoFAService)
			mockService.On("Disable", mock.Anything, tt.cmd.UserID).Return(tt.err)

			handler := NewDisableHandler(mockService)

			// Act
			err := handler.Handle(context.Background(), tt.cmd)

			// Assert
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
			mockService.AssertExpectations(t)
		})
	}
}
