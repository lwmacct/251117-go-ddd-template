//nolint:forcetypeassert // Mock 返回值类型在测试中总是已知的
package role

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/event"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/role"
)

// MockRoleCommandRepository 角色写仓储 Mock
type MockRoleCommandRepository struct {
	mock.Mock
}

func (m *MockRoleCommandRepository) Create(ctx context.Context, role *role.Role) error {
	args := m.Called(ctx, role)
	// 模拟数据库生成 ID
	if args.Error(0) == nil && role.ID == 0 {
		role.ID = 1
	}
	return args.Error(0)
}

func (m *MockRoleCommandRepository) Update(ctx context.Context, role *role.Role) error {
	args := m.Called(ctx, role)
	return args.Error(0)
}

func (m *MockRoleCommandRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRoleCommandRepository) SetPermissions(ctx context.Context, roleID uint, permissionIDs []uint) error {
	args := m.Called(ctx, roleID, permissionIDs)
	return args.Error(0)
}

func (m *MockRoleCommandRepository) AddPermissions(ctx context.Context, roleID uint, permissionIDs []uint) error {
	args := m.Called(ctx, roleID, permissionIDs)
	return args.Error(0)
}

func (m *MockRoleCommandRepository) RemovePermissions(ctx context.Context, roleID uint, permissionIDs []uint) error {
	args := m.Called(ctx, roleID, permissionIDs)
	return args.Error(0)
}

// MockRoleQueryRepository 角色读仓储 Mock
type MockRoleQueryRepository struct {
	mock.Mock
}

func (m *MockRoleQueryRepository) FindByID(ctx context.Context, id uint) (*role.Role, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*role.Role), args.Error(1)
}

func (m *MockRoleQueryRepository) FindByName(ctx context.Context, name string) (*role.Role, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*role.Role), args.Error(1)
}

func (m *MockRoleQueryRepository) FindByIDWithPermissions(ctx context.Context, id uint) (*role.Role, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*role.Role), args.Error(1)
}

func (m *MockRoleQueryRepository) List(ctx context.Context, page, limit int) ([]role.Role, int64, error) {
	args := m.Called(ctx, page, limit)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]role.Role), args.Get(1).(int64), args.Error(2)
}

func (m *MockRoleQueryRepository) GetPermissions(ctx context.Context, roleID uint) ([]role.Permission, error) {
	args := m.Called(ctx, roleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]role.Permission), args.Error(1)
}

func (m *MockRoleQueryRepository) Exists(ctx context.Context, id uint) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockRoleQueryRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	args := m.Called(ctx, name)
	return args.Bool(0), args.Error(1)
}

// MockPermissionQueryRepository 权限读仓储 Mock
type MockPermissionQueryRepository struct {
	mock.Mock
}

func (m *MockPermissionQueryRepository) FindByID(ctx context.Context, id uint) (*role.Permission, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*role.Permission), args.Error(1)
}

func (m *MockPermissionQueryRepository) FindByCode(ctx context.Context, code string) (*role.Permission, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*role.Permission), args.Error(1)
}

func (m *MockPermissionQueryRepository) FindByIDs(ctx context.Context, ids []uint) ([]role.Permission, error) {
	args := m.Called(ctx, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]role.Permission), args.Error(1)
}

func (m *MockPermissionQueryRepository) List(ctx context.Context, page, limit int) ([]role.Permission, int64, error) {
	args := m.Called(ctx, page, limit)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]role.Permission), args.Get(1).(int64), args.Error(2)
}

func (m *MockPermissionQueryRepository) ListByResource(ctx context.Context, resource string) ([]role.Permission, error) {
	args := m.Called(ctx, resource)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]role.Permission), args.Error(1)
}

func (m *MockPermissionQueryRepository) Exists(ctx context.Context, id uint) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockPermissionQueryRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	args := m.Called(ctx, code)
	return args.Bool(0), args.Error(1)
}

// MockEventBus 事件总线 Mock
type MockEventBus struct {
	mock.Mock
}

func (m *MockEventBus) Publish(ctx context.Context, events ...event.Event) error {
	args := m.Called(ctx, events)
	return args.Error(0)
}

func (m *MockEventBus) Subscribe(eventName string, handler event.EventHandler) {
	m.Called(eventName, handler)
}

func (m *MockEventBus) Unsubscribe(eventName string, handler event.EventHandler) {
	m.Called(eventName, handler)
}

func (m *MockEventBus) Close() error {
	args := m.Called()
	return args.Error(0)
}
