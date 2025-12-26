package twofa

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/twofa"
)

func TestSetupHandler_Handle_Success(t *testing.T) {
	tests := []struct {
		name string
		cmd  SetupCommand
		want *SetupResultDTO
	}{
		{
			name: "成功设置 2FA",
			cmd:  SetupCommand{UserID: 1},
			want: &SetupResultDTO{
				Secret:    "JBSWY3DPEHPK3PXP",
				QRCodeURL: "otpauth://totp/Test:user1?secret=JBSWY3DPEHPK3PXP&issuer=Test",
				QRCodeImg: "data:image/png;base64,iVBORw0KGgoAAAANSUhEUg...",
			},
		},
		{
			name: "不同用户设置 2FA",
			cmd:  SetupCommand{UserID: 100},
			want: &SetupResultDTO{
				Secret:    "HXDMVJECJJWSRB3H",
				QRCodeURL: "otpauth://totp/Test:user100?secret=HXDMVJECJJWSRB3H&issuer=Test",
				QRCodeImg: "data:image/png;base64,iVBORw0KGgoAAAANSUhEUg...",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockService := new(MockTwoFAService)
			mockService.On("Setup", mock.Anything, tt.cmd.UserID).Return(&twofa.SetupResult{
				Secret:    tt.want.Secret,
				QRCodeURL: tt.want.QRCodeURL,
				QRCodeImg: tt.want.QRCodeImg,
			}, nil)

			handler := NewSetupHandler(mockService)

			// Act
			result, err := handler.Handle(context.Background(), tt.cmd)

			// Assert
			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tt.want.Secret, result.Secret)
			assert.Equal(t, tt.want.QRCodeURL, result.QRCodeURL)
			assert.Equal(t, tt.want.QRCodeImg, result.QRCodeImg)
			mockService.AssertExpectations(t)
		})
	}
}

func TestSetupHandler_Handle_Error(t *testing.T) {
	tests := []struct {
		name    string
		cmd     SetupCommand
		err     error
		wantErr string
	}{
		{
			name:    "用户不存在",
			cmd:     SetupCommand{UserID: 999},
			err:     errors.New("user not found"),
			wantErr: "user not found",
		},
		{
			name:    "服务内部错误",
			cmd:     SetupCommand{UserID: 1},
			err:     errors.New("failed to generate TOTP key"),
			wantErr: "failed to generate TOTP key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockService := new(MockTwoFAService)
			mockService.On("Setup", mock.Anything, tt.cmd.UserID).Return(nil, tt.err)

			handler := NewSetupHandler(mockService)

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
