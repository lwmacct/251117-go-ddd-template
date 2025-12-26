package twofa

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetStatusHandler_Handle_Success(t *testing.T) {
	tests := []struct {
		name  string
		query GetStatusQuery
		want  *GetStatusResultDTO
	}{
		{
			name:  "2FA 已启用，有恢复码",
			query: GetStatusQuery{UserID: 1},
			want: &GetStatusResultDTO{
				Enabled:            true,
				RecoveryCodesCount: 8,
			},
		},
		{
			name:  "2FA 已启用，部分恢复码已使用",
			query: GetStatusQuery{UserID: 2},
			want: &GetStatusResultDTO{
				Enabled:            true,
				RecoveryCodesCount: 3,
			},
		},
		{
			name:  "2FA 未启用",
			query: GetStatusQuery{UserID: 3},
			want: &GetStatusResultDTO{
				Enabled:            false,
				RecoveryCodesCount: 0,
			},
		},
		{
			name:  "用户从未设置过 2FA",
			query: GetStatusQuery{UserID: 100},
			want: &GetStatusResultDTO{
				Enabled:            false,
				RecoveryCodesCount: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockService := new(MockTwoFAService)
			mockService.On("GetStatus", mock.Anything, tt.query.UserID).
				Return(tt.want.Enabled, tt.want.RecoveryCodesCount, nil)

			handler := NewGetStatusHandler(mockService)

			// Act
			result, err := handler.Handle(context.Background(), tt.query)

			// Assert
			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tt.want.Enabled, result.Enabled)
			assert.Equal(t, tt.want.RecoveryCodesCount, result.RecoveryCodesCount)
			mockService.AssertExpectations(t)
		})
	}
}

func TestGetStatusHandler_Handle_Error(t *testing.T) {
	tests := []struct {
		name    string
		query   GetStatusQuery
		err     error
		wantErr string
	}{
		{
			name:    "用户不存在",
			query:   GetStatusQuery{UserID: 999},
			err:     errors.New("user not found"),
			wantErr: "user not found",
		},
		{
			name:    "数据库查询失败",
			query:   GetStatusQuery{UserID: 1},
			err:     errors.New("failed to get 2FA config"),
			wantErr: "failed to get 2FA config",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockService := new(MockTwoFAService)
			mockService.On("GetStatus", mock.Anything, tt.query.UserID).
				Return(false, 0, tt.err)

			handler := NewGetStatusHandler(mockService)

			// Act
			result, err := handler.Handle(context.Background(), tt.query)

			// Assert
			assert.Nil(t, result)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
			mockService.AssertExpectations(t)
		})
	}
}
