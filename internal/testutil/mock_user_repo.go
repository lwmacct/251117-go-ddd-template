package testutil

import (
	"context"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
)

// MockUserCommandRepo 是用户命令仓储的 Mock 实现。
type MockUserCommandRepo struct {
	Users            []*user.User
	CreateError      error
	UpdateError      error
	DeleteError      error
	AssignRolesError error
	IDCounter        uint
	AssignedRoles    map[uint][]uint // userID -> roleIDs
}

// NewMockUserCommandRepo 创建 MockUserCommandRepo 实例
func NewMockUserCommandRepo() *MockUserCommandRepo {
	return &MockUserCommandRepo{
		Users:         make([]*user.User, 0),
		IDCounter:     0,
		AssignedRoles: make(map[uint][]uint),
	}
}

func (m *MockUserCommandRepo) Create(ctx context.Context, u *user.User) error {
	if m.CreateError != nil {
		return m.CreateError
	}
	m.IDCounter++
	u.ID = m.IDCounter
	m.Users = append(m.Users, u)
	return nil
}

func (m *MockUserCommandRepo) Update(ctx context.Context, u *user.User) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}
	return nil
}

func (m *MockUserCommandRepo) Delete(ctx context.Context, id uint) error {
	if m.DeleteError != nil {
		return m.DeleteError
	}
	return nil
}

func (m *MockUserCommandRepo) AssignRoles(ctx context.Context, userID uint, roleIDs []uint) error {
	if m.AssignRolesError != nil {
		return m.AssignRolesError
	}
	m.AssignedRoles[userID] = roleIDs
	return nil
}

func (m *MockUserCommandRepo) RemoveRoles(ctx context.Context, userID uint, roleIDs []uint) error {
	return nil
}

func (m *MockUserCommandRepo) UpdatePassword(ctx context.Context, userID uint, hashedPassword string) error {
	return nil
}

func (m *MockUserCommandRepo) UpdateStatus(ctx context.Context, userID uint, status string) error {
	return nil
}

// MockUserQueryRepo 是用户查询仓储的 Mock 实现。
type MockUserQueryRepo struct {
	Users             map[uint]*user.User
	ExistingUsernames map[string]bool
	ExistingEmails    map[string]bool
	CheckError        error
	FindError         error
}

// NewMockUserQueryRepo 创建 MockUserQueryRepo 实例
func NewMockUserQueryRepo() *MockUserQueryRepo {
	return &MockUserQueryRepo{
		Users:             make(map[uint]*user.User),
		ExistingUsernames: make(map[string]bool),
		ExistingEmails:    make(map[string]bool),
	}
}

func (m *MockUserQueryRepo) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	if m.CheckError != nil {
		return false, m.CheckError
	}
	return m.ExistingUsernames[username], nil
}

func (m *MockUserQueryRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	if m.CheckError != nil {
		return false, m.CheckError
	}
	return m.ExistingEmails[email], nil
}

func (m *MockUserQueryRepo) GetByID(ctx context.Context, id uint) (*user.User, error) {
	if m.FindError != nil {
		return nil, m.FindError
	}
	if u, ok := m.Users[id]; ok {
		return u, nil
	}
	return nil, nil //nolint:nilnil // test mock returns nil for not found
}

func (m *MockUserQueryRepo) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	return nil, nil //nolint:nilnil // test mock returns nil for not found
}

func (m *MockUserQueryRepo) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	return nil, nil //nolint:nilnil // test mock returns nil for not found
}

func (m *MockUserQueryRepo) GetByIDWithRoles(ctx context.Context, id uint) (*user.User, error) {
	return m.GetByID(ctx, id)
}

func (m *MockUserQueryRepo) GetByUsernameWithRoles(ctx context.Context, username string) (*user.User, error) {
	return nil, nil //nolint:nilnil // test mock returns nil for not found
}

func (m *MockUserQueryRepo) GetByEmailWithRoles(ctx context.Context, email string) (*user.User, error) {
	return nil, nil //nolint:nilnil // test mock returns nil for not found
}

func (m *MockUserQueryRepo) List(ctx context.Context, offset, limit int) ([]*user.User, error) {
	return nil, nil
}

func (m *MockUserQueryRepo) Count(ctx context.Context) (int64, error) {
	return 0, nil
}

func (m *MockUserQueryRepo) GetRoles(ctx context.Context, userID uint) ([]uint, error) {
	return nil, nil
}

func (m *MockUserQueryRepo) Search(ctx context.Context, keyword string, offset, limit int) ([]*user.User, error) {
	return nil, nil
}

func (m *MockUserQueryRepo) CountBySearch(ctx context.Context, keyword string) (int64, error) {
	return 0, nil
}

func (m *MockUserQueryRepo) Exists(ctx context.Context, id uint) (bool, error) {
	_, ok := m.Users[id]
	return ok, nil
}

func (m *MockUserQueryRepo) GetUserIDsByRole(ctx context.Context, roleID uint) ([]uint, error) {
	return nil, nil
}
