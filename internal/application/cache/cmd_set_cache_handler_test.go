package cache

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSetCacheHandler_Handle_Success(t *testing.T) {
	tests := []struct {
		name  string
		cmd   SetCacheCommand
		value any
	}{
		{
			name: "设置字符串值",
			cmd: SetCacheCommand{
				Key:   "user:name",
				Value: "John Doe",
				TTL:   5 * time.Minute,
			},
			value: "John Doe",
		},
		{
			name: "设置整数值",
			cmd: SetCacheCommand{
				Key:   "user:age",
				Value: 30,
				TTL:   10 * time.Minute,
			},
			value: 30,
		},
		{
			name: "设置结构体值",
			cmd: SetCacheCommand{
				Key: "user:profile",
				Value: map[string]any{
					"id":   1,
					"name": "John",
				},
				TTL: time.Hour,
			},
			value: map[string]any{
				"id":   1,
				"name": "John",
			},
		},
		{
			name: "设置永久缓存",
			cmd: SetCacheCommand{
				Key:   "config:app",
				Value: "production",
				TTL:   0, // 0 表示永不过期
			},
			value: "production",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockCommandRepo := new(MockCacheCommandRepository)
			mockCommandRepo.On("Set", mock.Anything, tt.cmd.Key, tt.value, tt.cmd.TTL).
				Return(nil)

			handler := NewSetCacheHandler(mockCommandRepo)

			// Act
			err := handler.Handle(context.Background(), tt.cmd)

			// Assert
			require.NoError(t, err)
			mockCommandRepo.AssertExpectations(t)
		})
	}
}

func TestSetCacheHandler_Handle_Error(t *testing.T) {
	tests := []struct {
		name       string
		cmd        SetCacheCommand
		setupMocks func(*MockCacheCommandRepository)
		wantErr    string
	}{
		{
			name: "仓储设置失败",
			cmd: SetCacheCommand{
				Key:   "test_key",
				Value: "test_value",
				TTL:   time.Minute,
			},
			setupMocks: func(repo *MockCacheCommandRepository) {
				repo.On("Set", mock.Anything, "test_key", "test_value", time.Minute).
					Return(errors.New("redis: connection refused"))
			},
			wantErr: "redis: connection refused",
		},
		{
			name: "空键名",
			cmd: SetCacheCommand{
				Key:   "",
				Value: "value",
				TTL:   time.Minute,
			},
			setupMocks: func(repo *MockCacheCommandRepository) {
				repo.On("Set", mock.Anything, "", "value", time.Minute).
					Return(errors.New("empty key not allowed"))
			},
			wantErr: "empty key not allowed",
		},
		{
			name: "内存不足",
			cmd: SetCacheCommand{
				Key:   "large_key",
				Value: make([]byte, 1<<20), // 1MB
				TTL:   time.Hour,
			},
			setupMocks: func(repo *MockCacheCommandRepository) {
				repo.On("Set", mock.Anything, "large_key", mock.Anything, time.Hour).
					Return(errors.New("OOM command not allowed"))
			},
			wantErr: "OOM command not allowed",
		},
		{
			name: "超时",
			cmd: SetCacheCommand{
				Key:   "timeout_key",
				Value: "value",
				TTL:   time.Second,
			},
			setupMocks: func(repo *MockCacheCommandRepository) {
				repo.On("Set", mock.Anything, "timeout_key", "value", time.Second).
					Return(errors.New("i/o timeout"))
			},
			wantErr: "i/o timeout",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockCommandRepo := new(MockCacheCommandRepository)
			tt.setupMocks(mockCommandRepo)

			handler := NewSetCacheHandler(mockCommandRepo)

			// Act
			err := handler.Handle(context.Background(), tt.cmd)

			// Assert
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
			mockCommandRepo.AssertExpectations(t)
		})
	}
}

func TestSetCacheHandler_Handle_ContextCancellation(t *testing.T) {
	// Arrange
	mockCommandRepo := new(MockCacheCommandRepository)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // 立即取消 context

	mockCommandRepo.On("Set", mock.Anything, "test_key", "test_value", time.Minute).
		Return(context.Canceled)

	handler := NewSetCacheHandler(mockCommandRepo)

	// Act
	err := handler.Handle(ctx, SetCacheCommand{
		Key:   "test_key",
		Value: "test_value",
		TTL:   time.Minute,
	})

	// Assert
	require.Error(t, err)
	require.ErrorIs(t, err, context.Canceled)
	mockCommandRepo.AssertExpectations(t)
}
