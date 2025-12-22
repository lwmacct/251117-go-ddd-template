package testutil

import (
	"context"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

// MockSettingCommandRepo 是设置命令仓储的 Mock 实现。
type MockSettingCommandRepo struct {
	CreateError error
	UpdateError error
	DeleteError error
}

// NewMockSettingCommandRepo 创建 MockSettingCommandRepo 实例
func NewMockSettingCommandRepo() *MockSettingCommandRepo {
	return &MockSettingCommandRepo{}
}

func (m *MockSettingCommandRepo) Create(ctx context.Context, s *setting.Setting) error {
	if m.CreateError != nil {
		return m.CreateError
	}
	return nil
}

func (m *MockSettingCommandRepo) Update(ctx context.Context, s *setting.Setting) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}
	return nil
}

func (m *MockSettingCommandRepo) Delete(ctx context.Context, id uint) error {
	if m.DeleteError != nil {
		return m.DeleteError
	}
	return nil
}

func (m *MockSettingCommandRepo) BatchUpsert(ctx context.Context, settings []*setting.Setting) error {
	return nil
}

// MockSettingQueryRepo 是设置查询仓储的 Mock 实现。
type MockSettingQueryRepo struct {
	Settings  map[string]*setting.Setting
	FindError error
}

// NewMockSettingQueryRepo 创建 MockSettingQueryRepo 实例
func NewMockSettingQueryRepo() *MockSettingQueryRepo {
	return &MockSettingQueryRepo{
		Settings: make(map[string]*setting.Setting),
	}
}

func (m *MockSettingQueryRepo) FindByID(ctx context.Context, id uint) (*setting.Setting, error) {
	return nil, nil //nolint:nilnil // test mock returns nil for not found
}

func (m *MockSettingQueryRepo) FindByKey(ctx context.Context, key string) (*setting.Setting, error) {
	if m.FindError != nil {
		return nil, m.FindError
	}
	if s, ok := m.Settings[key]; ok {
		return s, nil
	}
	return nil, nil //nolint:nilnil // test mock returns nil for not found
}

func (m *MockSettingQueryRepo) FindByCategory(ctx context.Context, category string) ([]*setting.Setting, error) {
	return nil, nil
}

func (m *MockSettingQueryRepo) FindAll(ctx context.Context) ([]*setting.Setting, error) {
	return nil, nil
}
