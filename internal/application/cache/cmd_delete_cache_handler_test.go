package cache

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDeleteCacheHandler_Handle_Success(t *testing.T) {
	// Arrange
	mockCommandRepo := new(MockCacheCommandRepository)
	mockCommandRepo.On("Delete", mock.Anything, "test_key").Return(nil)

	handler := NewDeleteCacheHandler(mockCommandRepo)

	// Act
	err := handler.Handle(context.Background(), DeleteCacheCommand{
		Key: "test_key",
	})

	// Assert
	require.NoError(t, err)
	mockCommandRepo.AssertExpectations(t)
}

func TestDeleteCacheHandler_Handle_Error(t *testing.T) {
	tests := []struct {
		name       string
		cmd        DeleteCacheCommand
		setupMocks func(*MockCacheCommandRepository)
		wantErr    string
	}{
		{
			name: "仓储删除失败",
			cmd: DeleteCacheCommand{
				Key: "nonexistent_key",
			},
			setupMocks: func(repo *MockCacheCommandRepository) {
				repo.On("Delete", mock.Anything, "nonexistent_key").
					Return(errors.New("redis: nil"))
			},
			wantErr: "redis: nil",
		},
		{
			name: "空键名",
			cmd: DeleteCacheCommand{
				Key: "",
			},
			setupMocks: func(repo *MockCacheCommandRepository) {
				repo.On("Delete", mock.Anything, "").
					Return(errors.New("empty key"))
			},
			wantErr: "empty key",
		},
		{
			name: "网络错误",
			cmd: DeleteCacheCommand{
				Key: "test_key",
			},
			setupMocks: func(repo *MockCacheCommandRepository) {
				repo.On("Delete", mock.Anything, "test_key").
					Return(errors.New("connection timeout"))
			},
			wantErr: "connection timeout",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockCommandRepo := new(MockCacheCommandRepository)
			tt.setupMocks(mockCommandRepo)

			handler := NewDeleteCacheHandler(mockCommandRepo)

			// Act
			err := handler.Handle(context.Background(), tt.cmd)

			// Assert
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
			mockCommandRepo.AssertExpectations(t)
		})
	}
}

func TestDeleteCacheHandler_Handle_ContextCancellation(t *testing.T) {
	// Arrange
	mockCommandRepo := new(MockCacheCommandRepository)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // 立即取消 context

	mockCommandRepo.On("Delete", mock.Anything, "test_key").
		Return(context.Canceled)

	handler := NewDeleteCacheHandler(mockCommandRepo)

	// Act
	err := handler.Handle(ctx, DeleteCacheCommand{
		Key: "test_key",
	})

	// Assert
	require.Error(t, err)
	require.ErrorIs(t, err, context.Canceled)
	mockCommandRepo.AssertExpectations(t)
}
