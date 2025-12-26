package twofa

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestVerifyEnableHandler_Handle_Success(t *testing.T) {
	tests := []struct {
		name string
		cmd  VerifyEnableCommand
		want *VerifyEnableResultDTO
	}{
		{
			name: "成功验证并启用 2FA",
			cmd:  VerifyEnableCommand{UserID: 1, Code: "123456"},
			want: &VerifyEnableResultDTO{
				RecoveryCodes: []string{
					"AAAA-BBBB-CCCC",
					"DDDD-EEEE-FFFF",
					"GGGG-HHHH-IIII",
					"JJJJ-KKKK-LLLL",
				},
			},
		},
		{
			name: "使用 6 位验证码启用",
			cmd:  VerifyEnableCommand{UserID: 2, Code: "654321"},
			want: &VerifyEnableResultDTO{
				RecoveryCodes: []string{
					"CODE-0001-XXXX",
					"CODE-0002-YYYY",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockService := new(MockTwoFAService)
			mockService.On("VerifyAndEnable", mock.Anything, tt.cmd.UserID, tt.cmd.Code).
				Return(tt.want.RecoveryCodes, nil)

			handler := NewVerifyEnableHandler(mockService)

			// Act
			result, err := handler.Handle(context.Background(), tt.cmd)

			// Assert
			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tt.want.RecoveryCodes, result.RecoveryCodes)
			assert.Len(t, result.RecoveryCodes, len(tt.want.RecoveryCodes))
			mockService.AssertExpectations(t)
		})
	}
}

func TestVerifyEnableHandler_Handle_Error(t *testing.T) {
	tests := []struct {
		name    string
		cmd     VerifyEnableCommand
		err     error
		wantErr string
	}{
		{
			name:    "验证码错误",
			cmd:     VerifyEnableCommand{UserID: 1, Code: "000000"},
			err:     errors.New("invalid verification code"),
			wantErr: "invalid verification code",
		},
		{
			name:    "未设置 2FA 密钥",
			cmd:     VerifyEnableCommand{UserID: 2, Code: "123456"},
			err:     errors.New("please setup 2FA first"),
			wantErr: "please setup 2FA first",
		},
		{
			name:    "用户不存在",
			cmd:     VerifyEnableCommand{UserID: 999, Code: "123456"},
			err:     errors.New("failed to get 2FA config"),
			wantErr: "failed to get 2FA config",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockService := new(MockTwoFAService)
			mockService.On("VerifyAndEnable", mock.Anything, tt.cmd.UserID, tt.cmd.Code).
				Return(nil, tt.err)

			handler := NewVerifyEnableHandler(mockService)

			// Act
			result, err := handler.Handle(context.Background(), tt.cmd)

			// Assert
			assert.Nil(t, result)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
			mockService.AssertExpectations(t)
		})
	}
}
