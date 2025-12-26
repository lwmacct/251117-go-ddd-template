//nolint:forcetypeassert // 测试中的类型断言是可控的
package cache

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetCacheHandler_Handle_Success(t *testing.T) {
	tests := []struct {
		name          string
		query         GetCacheQuery
		mockValue     any
		expectedValue any
	}{
		{
			name: "获取字符串值",
			query: GetCacheQuery{
				Key: "user:name",
			},
			mockValue:     "John Doe",
			expectedValue: "John Doe",
		},
		{
			name: "获取整数值",
			query: GetCacheQuery{
				Key: "user:age",
			},
			mockValue:     30,
			expectedValue: 30,
		},
		{
			name: "获取布尔值",
			query: GetCacheQuery{
				Key: "user:active",
			},
			mockValue:     true,
			expectedValue: true,
		},
		{
			name: "获取结构体值",
			query: GetCacheQuery{
				Key: "user:profile",
			},
			mockValue: map[string]any{
				"id":   1,
				"name": "John",
			},
			expectedValue: map[string]any{
				"id":   1,
				"name": "John",
			},
		},
		{
			name: "获取数组值",
			query: GetCacheQuery{
				Key: "user:tags",
			},
			mockValue:     []string{"admin", "editor"},
			expectedValue: []string{"admin", "editor"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockQueryRepo := new(MockCacheQueryRepository)
			mockQueryRepo.On("Get", mock.Anything, tt.query.Key, mock.Anything).
				Run(func(args mock.Arguments) {
					// 模拟设置 dest 参数的值
					dest := args.Get(2).(*any)
					*dest = tt.mockValue
				}).
				Return(nil)

			handler := NewGetCacheHandler(mockQueryRepo)

			// Act
			result, err := handler.Handle(context.Background(), tt.query)

			// Assert
			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tt.query.Key, result.Key)
			assert.Equal(t, tt.expectedValue, result.Value)
			mockQueryRepo.AssertExpectations(t)
		})
	}
}

func TestGetCacheHandler_Handle_Error(t *testing.T) {
	tests := []struct {
		name       string
		query      GetCacheQuery
		setupMocks func(*MockCacheQueryRepository)
		wantErr    string
	}{
		{
			name: "键不存在",
			query: GetCacheQuery{
				Key: "nonexistent_key",
			},
			setupMocks: func(repo *MockCacheQueryRepository) {
				repo.On("Get", mock.Anything, "nonexistent_key", mock.Anything).
					Return(errors.New("redis: nil"))
			},
			wantErr: "redis: nil",
		},
		{
			name: "空键名",
			query: GetCacheQuery{
				Key: "",
			},
			setupMocks: func(repo *MockCacheQueryRepository) {
				repo.On("Get", mock.Anything, "", mock.Anything).
					Return(errors.New("empty key"))
			},
			wantErr: "empty key",
		},
		{
			name: "网络错误",
			query: GetCacheQuery{
				Key: "test_key",
			},
			setupMocks: func(repo *MockCacheQueryRepository) {
				repo.On("Get", mock.Anything, "test_key", mock.Anything).
					Return(errors.New("connection timeout"))
			},
			wantErr: "connection timeout",
		},
		{
			name: "反序列化失败",
			query: GetCacheQuery{
				Key: "corrupt_key",
			},
			setupMocks: func(repo *MockCacheQueryRepository) {
				repo.On("Get", mock.Anything, "corrupt_key", mock.Anything).
					Return(errors.New("invalid character"))
			},
			wantErr: "invalid character",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockQueryRepo := new(MockCacheQueryRepository)
			tt.setupMocks(mockQueryRepo)

			handler := NewGetCacheHandler(mockQueryRepo)

			// Act
			result, err := handler.Handle(context.Background(), tt.query)

			// Assert
			require.Error(t, err)
			assert.Nil(t, result)
			assert.Contains(t, err.Error(), tt.wantErr)
			mockQueryRepo.AssertExpectations(t)
		})
	}
}

func TestGetCacheHandler_Handle_ContextCancellation(t *testing.T) {
	// Arrange
	mockQueryRepo := new(MockCacheQueryRepository)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // 立即取消 context

	mockQueryRepo.On("Get", mock.Anything, "test_key", mock.Anything).
		Return(context.Canceled)

	handler := NewGetCacheHandler(mockQueryRepo)

	// Act
	result, err := handler.Handle(ctx, GetCacheQuery{
		Key: "test_key",
	})

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	require.ErrorIs(t, err, context.Canceled)
	mockQueryRepo.AssertExpectations(t)
}

func TestGetCacheHandler_Handle_NilValue(t *testing.T) {
	// Arrange
	mockQueryRepo := new(MockCacheQueryRepository)
	mockQueryRepo.On("Get", mock.Anything, "nil_key", mock.Anything).
		Run(func(args mock.Arguments) {
			dest := args.Get(2).(*any)
			*dest = nil
		}).
		Return(nil)

	handler := NewGetCacheHandler(mockQueryRepo)

	// Act
	result, err := handler.Handle(context.Background(), GetCacheQuery{
		Key: "nil_key",
	})

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "nil_key", result.Key)
	assert.Nil(t, result.Value)
	mockQueryRepo.AssertExpectations(t)
}
