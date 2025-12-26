package cache

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

// ============================================================
// MockCacheCommandRepository
// ============================================================

// MockCacheCommandRepository 缓存写操作仓储 Mock
type MockCacheCommandRepository struct {
	mock.Mock
}

// Set 模拟设置缓存值
func (m *MockCacheCommandRepository) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

// Delete 模拟删除缓存值
func (m *MockCacheCommandRepository) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

// SetNX 模拟仅当键不存在时设置值
func (m *MockCacheCommandRepository) SetNX(ctx context.Context, key string, value any, expiration time.Duration) (bool, error) {
	args := m.Called(ctx, key, value, expiration)
	return args.Bool(0), args.Error(1)
}

// ============================================================
// MockCacheQueryRepository
// ============================================================

// MockCacheQueryRepository 缓存读操作仓储 Mock
type MockCacheQueryRepository struct {
	mock.Mock
}

// Get 模拟获取缓存值
func (m *MockCacheQueryRepository) Get(ctx context.Context, key string, dest any) error {
	args := m.Called(ctx, key, dest)
	return args.Error(0)
}

// Exists 模拟检查键是否存在
func (m *MockCacheQueryRepository) Exists(ctx context.Context, key string) (bool, error) {
	args := m.Called(ctx, key)
	return args.Bool(0), args.Error(1)
}
