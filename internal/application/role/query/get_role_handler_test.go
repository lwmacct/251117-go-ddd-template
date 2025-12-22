package query

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/role"
)

// mockRoleQueryRepo 是角色查询仓储的 Mock 实现
type mockRoleQueryRepo struct {
	roles             map[uint]*role.Role
	findByIDErr       error
	findWithPermsErr  error
	listRoles         []role.Role
	listTotal         int64
	listErr           error
	permissions       []role.Permission
	getPermissionsErr error
}

func newMockRoleQueryRepo() *mockRoleQueryRepo {
	return &mockRoleQueryRepo{
		roles: make(map[uint]*role.Role),
	}
}

func (m *mockRoleQueryRepo) FindByID(ctx context.Context, id uint) (*role.Role, error) {
	if m.findByIDErr != nil {
		return nil, m.findByIDErr
	}
	if r, ok := m.roles[id]; ok {
		return r, nil
	}
	return nil, nil //nolint:nilnil // test mock returns nil for not found
}

func (m *mockRoleQueryRepo) FindByIDWithPermissions(ctx context.Context, id uint) (*role.Role, error) {
	if m.findWithPermsErr != nil {
		return nil, m.findWithPermsErr
	}
	if r, ok := m.roles[id]; ok {
		return r, nil
	}
	return nil, nil //nolint:nilnil // test mock returns nil for not found
}

func (m *mockRoleQueryRepo) FindByName(ctx context.Context, name string) (*role.Role, error) {
	return nil, nil //nolint:nilnil // test mock returns nil for not found
}

func (m *mockRoleQueryRepo) List(ctx context.Context, page, limit int) ([]role.Role, int64, error) {
	if m.listErr != nil {
		return nil, 0, m.listErr
	}
	return m.listRoles, m.listTotal, nil
}

func (m *mockRoleQueryRepo) GetPermissions(ctx context.Context, roleID uint) ([]role.Permission, error) {
	if m.getPermissionsErr != nil {
		return nil, m.getPermissionsErr
	}
	return m.permissions, nil
}

func (m *mockRoleQueryRepo) ExistsByName(ctx context.Context, name string) (bool, error) {
	return false, nil
}

func (m *mockRoleQueryRepo) Exists(ctx context.Context, id uint) (bool, error) {
	_, ok := m.roles[id]
	return ok, nil
}

type getRoleTestSetup struct {
	handler   *GetRoleHandler
	queryRepo *mockRoleQueryRepo
}

func setupGetRoleTest() *getRoleTestSetup {
	queryRepo := newMockRoleQueryRepo()
	handler := NewGetRoleHandler(queryRepo)

	return &getRoleTestSetup{
		handler:   handler,
		queryRepo: queryRepo,
	}
}

func TestGetRoleHandler_Handle(t *testing.T) {
	ctx := context.Background()

	t.Run("成功获取角色", func(t *testing.T) {
		setup := setupGetRoleTest()
		setup.queryRepo.roles[1] = &role.Role{
			ID:          1,
			Name:        "editor",
			DisplayName: "Editor",
			Description: "Can edit content",
			IsSystem:    false,
			Permissions: []role.Permission{
				{ID: 1, Code: "article:read", Description: "Read Articles"},
				{ID: 2, Code: "article:write", Description: "Write Articles"},
			},
		}

		query := GetRoleQuery{
			RoleID: 1,
		}

		result, err := setup.handler.Handle(ctx, query)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, uint(1), result.ID)
		assert.Equal(t, "editor", result.Name)
		assert.Equal(t, "Editor", result.DisplayName)
		assert.Len(t, result.Permissions, 2)
	})

	t.Run("角色不存在", func(t *testing.T) {
		setup := setupGetRoleTest()

		query := GetRoleQuery{
			RoleID: 999,
		}

		result, err := setup.handler.Handle(ctx, query)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "role not found")
		assert.Nil(t, result)
	})

	t.Run("查询失败", func(t *testing.T) {
		setup := setupGetRoleTest()
		setup.queryRepo.findWithPermsErr = errors.New("database error")

		query := GetRoleQuery{
			RoleID: 1,
		}

		result, err := setup.handler.Handle(ctx, query)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to find role")
		assert.Nil(t, result)
	})

	t.Run("系统角色", func(t *testing.T) {
		setup := setupGetRoleTest()
		setup.queryRepo.roles[1] = &role.Role{
			ID:          1,
			Name:        "admin",
			DisplayName: "Administrator",
			IsSystem:    true,
		}

		query := GetRoleQuery{
			RoleID: 1,
		}

		result, err := setup.handler.Handle(ctx, query)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.True(t, result.IsSystem)
	})

	t.Run("角色无权限", func(t *testing.T) {
		setup := setupGetRoleTest()
		setup.queryRepo.roles[1] = &role.Role{
			ID:          1,
			Name:        "guest",
			DisplayName: "Guest",
			Permissions: []role.Permission{},
		}

		query := GetRoleQuery{
			RoleID: 1,
		}

		result, err := setup.handler.Handle(ctx, query)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Empty(t, result.Permissions)
	})
}

func TestNewGetRoleHandler(t *testing.T) {
	queryRepo := newMockRoleQueryRepo()
	handler := NewGetRoleHandler(queryRepo)

	assert.NotNil(t, handler)
}
