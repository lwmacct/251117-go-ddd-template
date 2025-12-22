package testutil

import (
	"context"
	"errors"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/role"
)

// MockRoleCommandRepo 是角色命令仓储的 Mock 实现。
type MockRoleCommandRepo struct {
	Roles       []*role.Role
	CreateError error
	UpdateError error
	DeleteError error
	IDCounter   uint
}

// NewMockRoleCommandRepo 创建 MockRoleCommandRepo 实例
func NewMockRoleCommandRepo() *MockRoleCommandRepo {
	return &MockRoleCommandRepo{
		Roles:     make([]*role.Role, 0),
		IDCounter: 0,
	}
}

func (m *MockRoleCommandRepo) Create(ctx context.Context, r *role.Role) error {
	if m.CreateError != nil {
		return m.CreateError
	}
	m.IDCounter++
	r.ID = m.IDCounter
	m.Roles = append(m.Roles, r)
	return nil
}

func (m *MockRoleCommandRepo) Update(ctx context.Context, r *role.Role) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}
	return nil
}

func (m *MockRoleCommandRepo) Delete(ctx context.Context, id uint) error {
	if m.DeleteError != nil {
		return m.DeleteError
	}
	return nil
}

func (m *MockRoleCommandRepo) SetPermissions(ctx context.Context, roleID uint, permissionIDs []uint) error {
	return nil
}

func (m *MockRoleCommandRepo) AddPermissions(ctx context.Context, roleID uint, permissionIDs []uint) error {
	return nil
}

func (m *MockRoleCommandRepo) RemovePermissions(ctx context.Context, roleID uint, permissionIDs []uint) error {
	return nil
}

// MockRoleQueryRepo 是角色查询仓储的 Mock 实现。
type MockRoleQueryRepo struct {
	Roles         map[uint]*role.Role
	ExistingNames map[string]bool
	CheckError    error
	FindError     error
}

// NewMockRoleQueryRepo 创建 MockRoleQueryRepo 实例
func NewMockRoleQueryRepo() *MockRoleQueryRepo {
	return &MockRoleQueryRepo{
		Roles:         make(map[uint]*role.Role),
		ExistingNames: make(map[string]bool),
	}
}

func (m *MockRoleQueryRepo) ExistsByName(ctx context.Context, name string) (bool, error) {
	if m.CheckError != nil {
		return false, m.CheckError
	}
	return m.ExistingNames[name], nil
}

func (m *MockRoleQueryRepo) FindByID(ctx context.Context, id uint) (*role.Role, error) {
	if m.FindError != nil {
		return nil, m.FindError
	}
	if r, ok := m.Roles[id]; ok {
		return r, nil
	}
	return nil, errors.New("role not found")
}

func (m *MockRoleQueryRepo) FindByName(ctx context.Context, name string) (*role.Role, error) {
	return nil, nil //nolint:nilnil // test mock returns nil for not found
}

func (m *MockRoleQueryRepo) FindByIDWithPermissions(ctx context.Context, id uint) (*role.Role, error) {
	return m.FindByID(ctx, id)
}

func (m *MockRoleQueryRepo) List(ctx context.Context, page, limit int) ([]role.Role, int64, error) {
	return nil, 0, nil
}

func (m *MockRoleQueryRepo) GetPermissions(ctx context.Context, roleID uint) ([]role.Permission, error) {
	return nil, nil
}

func (m *MockRoleQueryRepo) Exists(ctx context.Context, id uint) (bool, error) {
	_, ok := m.Roles[id]
	return ok, nil
}
