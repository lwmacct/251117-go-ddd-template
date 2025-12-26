//nolint:forcetypeassert // Mock 返回值类型在测试中总是已知的
package setting

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

// MockSettingCommandRepository 设置写仓储 Mock
type MockSettingCommandRepository struct {
	mock.Mock
}

func (m *MockSettingCommandRepository) Create(ctx context.Context, s *setting.Setting) error {
	args := m.Called(ctx, s)
	if args.Error(0) == nil && s.ID == 0 {
		s.ID = 1
	}
	return args.Error(0)
}

func (m *MockSettingCommandRepository) Update(ctx context.Context, s *setting.Setting) error {
	args := m.Called(ctx, s)
	return args.Error(0)
}

func (m *MockSettingCommandRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSettingCommandRepository) BatchUpsert(ctx context.Context, settings []*setting.Setting) error {
	args := m.Called(ctx, settings)
	return args.Error(0)
}

// MockSettingQueryRepository 设置读仓储 Mock
type MockSettingQueryRepository struct {
	mock.Mock
}

func (m *MockSettingQueryRepository) FindByID(ctx context.Context, id uint) (*setting.Setting, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*setting.Setting), args.Error(1)
}

func (m *MockSettingQueryRepository) FindByKey(ctx context.Context, key string) (*setting.Setting, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*setting.Setting), args.Error(1)
}

func (m *MockSettingQueryRepository) FindByKeys(ctx context.Context, keys []string) ([]*setting.Setting, error) {
	args := m.Called(ctx, keys)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*setting.Setting), args.Error(1)
}

func (m *MockSettingQueryRepository) FindByCategory(ctx context.Context, category string) ([]*setting.Setting, error) {
	args := m.Called(ctx, category)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*setting.Setting), args.Error(1)
}

func (m *MockSettingQueryRepository) FindAll(ctx context.Context) ([]*setting.Setting, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*setting.Setting), args.Error(1)
}
