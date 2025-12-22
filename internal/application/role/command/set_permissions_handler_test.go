package command

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/role"
)

// mockPermissionQueryRepo 权限查询仓储 Mock
type mockPermissionQueryRepo struct {
	existingPermissions map[uint]bool
	existsError         error
}

func newMockPermissionQueryRepo() *mockPermissionQueryRepo {
	return &mockPermissionQueryRepo{
		existingPermissions: make(map[uint]bool),
	}
}

func (m *mockPermissionQueryRepo) FindByID(ctx context.Context, id uint) (*role.Permission, error) {
	return nil, nil //nolint:nilnil // test mock returns nil for not found
}

func (m *mockPermissionQueryRepo) FindByCode(ctx context.Context, code string) (*role.Permission, error) {
	return nil, nil //nolint:nilnil // test mock returns nil for not found
}

func (m *mockPermissionQueryRepo) FindByIDs(ctx context.Context, ids []uint) ([]role.Permission, error) {
	return nil, nil
}

func (m *mockPermissionQueryRepo) List(ctx context.Context, page, limit int) ([]role.Permission, int64, error) {
	return nil, 0, nil
}

func (m *mockPermissionQueryRepo) ListByResource(ctx context.Context, resource string) ([]role.Permission, error) {
	return nil, nil
}

func (m *mockPermissionQueryRepo) Exists(ctx context.Context, id uint) (bool, error) {
	if m.existsError != nil {
		return false, m.existsError
	}
	return m.existingPermissions[id], nil
}

func (m *mockPermissionQueryRepo) ExistsByCode(ctx context.Context, code string) (bool, error) {
	return false, nil
}

// setPermissionsCommandRepo 设置权限专用 Mock
type setPermissionsCommandRepo struct {
	setPermissionsErr error
	setPermCalled     bool
	setPermRoleID     uint
	setPermIDs        []uint
}

func newSetPermissionsCommandRepo() *setPermissionsCommandRepo {
	return &setPermissionsCommandRepo{}
}

func (m *setPermissionsCommandRepo) Create(ctx context.Context, r *role.Role) error {
	return nil
}

func (m *setPermissionsCommandRepo) Update(ctx context.Context, r *role.Role) error {
	return nil
}

func (m *setPermissionsCommandRepo) Delete(ctx context.Context, id uint) error {
	return nil
}

func (m *setPermissionsCommandRepo) SetPermissions(ctx context.Context, roleID uint, permissionIDs []uint) error {
	m.setPermCalled = true
	m.setPermRoleID = roleID
	m.setPermIDs = permissionIDs
	return m.setPermissionsErr
}

func (m *setPermissionsCommandRepo) AddPermissions(ctx context.Context, roleID uint, permissionIDs []uint) error {
	return nil
}

func (m *setPermissionsCommandRepo) RemovePermissions(ctx context.Context, roleID uint, permissionIDs []uint) error {
	return nil
}

// setPermissionsQueryRepo 设置权限专用查询 Mock
type setPermissionsQueryRepo struct {
	existingRoles map[uint]bool
	existsError   error
}

func newSetPermissionsQueryRepo() *setPermissionsQueryRepo {
	return &setPermissionsQueryRepo{
		existingRoles: make(map[uint]bool),
	}
}

func (m *setPermissionsQueryRepo) ExistsByName(ctx context.Context, name string) (bool, error) {
	return false, nil
}

func (m *setPermissionsQueryRepo) FindByID(ctx context.Context, id uint) (*role.Role, error) {
	return nil, nil //nolint:nilnil // test mock returns nil for not found
}

func (m *setPermissionsQueryRepo) FindByName(ctx context.Context, name string) (*role.Role, error) {
	return nil, nil //nolint:nilnil // test mock returns nil for not found
}

func (m *setPermissionsQueryRepo) FindByIDWithPermissions(ctx context.Context, id uint) (*role.Role, error) {
	return nil, nil //nolint:nilnil // test mock returns nil for not found
}

func (m *setPermissionsQueryRepo) List(ctx context.Context, page, limit int) ([]role.Role, int64, error) {
	return nil, 0, nil
}

func (m *setPermissionsQueryRepo) GetPermissions(ctx context.Context, roleID uint) ([]role.Permission, error) {
	return nil, nil
}

func (m *setPermissionsQueryRepo) Exists(ctx context.Context, id uint) (bool, error) {
	if m.existsError != nil {
		return false, m.existsError
	}
	return m.existingRoles[id], nil
}

type setPermissionsTestSetup struct {
	handler         *SetPermissionsHandler
	commandRepo     *setPermissionsCommandRepo
	queryRepo       *setPermissionsQueryRepo
	permissionQuery *mockPermissionQueryRepo
}

func setupSetPermissionsTest() *setPermissionsTestSetup {
	commandRepo := newSetPermissionsCommandRepo()
	queryRepo := newSetPermissionsQueryRepo()
	permissionQuery := newMockPermissionQueryRepo()
	handler := NewSetPermissionsHandler(commandRepo, queryRepo, permissionQuery)

	return &setPermissionsTestSetup{
		handler:         handler,
		commandRepo:     commandRepo,
		queryRepo:       queryRepo,
		permissionQuery: permissionQuery,
	}
}

func TestSetPermissionsHandler_Handle(t *testing.T) {
	ctx := context.Background()

	t.Run("成功设置权限", func(t *testing.T) {
		setup := setupSetPermissionsTest()
		setup.queryRepo.existingRoles[1] = true
		setup.permissionQuery.existingPermissions[10] = true
		setup.permissionQuery.existingPermissions[20] = true

		cmd := SetPermissionsCommand{
			RoleID:        1,
			PermissionIDs: []uint{10, 20},
		}

		err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		assert.True(t, setup.commandRepo.setPermCalled)
		assert.Equal(t, uint(1), setup.commandRepo.setPermRoleID)
		assert.Equal(t, []uint{10, 20}, setup.commandRepo.setPermIDs)
	})

	t.Run("角色不存在", func(t *testing.T) {
		setup := setupSetPermissionsTest()
		// 角色 999 不存在

		cmd := SetPermissionsCommand{
			RoleID:        999,
			PermissionIDs: []uint{10},
		}

		err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "role not found")
	})

	t.Run("检查角色存在性失败", func(t *testing.T) {
		setup := setupSetPermissionsTest()
		setup.queryRepo.existsError = errors.New("database error")

		cmd := SetPermissionsCommand{
			RoleID:        1,
			PermissionIDs: []uint{10},
		}

		err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to check role existence")
	})

	t.Run("权限不存在", func(t *testing.T) {
		setup := setupSetPermissionsTest()
		setup.queryRepo.existingRoles[1] = true
		setup.permissionQuery.existingPermissions[10] = true
		// 权限 999 不存在

		cmd := SetPermissionsCommand{
			RoleID:        1,
			PermissionIDs: []uint{10, 999},
		}

		err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "permission not found")
	})

	t.Run("检查权限存在性失败", func(t *testing.T) {
		setup := setupSetPermissionsTest()
		setup.queryRepo.existingRoles[1] = true
		setup.permissionQuery.existsError = errors.New("database error")

		cmd := SetPermissionsCommand{
			RoleID:        1,
			PermissionIDs: []uint{10},
		}

		err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to check permission existence")
	})

	t.Run("设置权限失败", func(t *testing.T) {
		setup := setupSetPermissionsTest()
		setup.queryRepo.existingRoles[1] = true
		setup.permissionQuery.existingPermissions[10] = true
		setup.commandRepo.setPermissionsErr = errors.New("database error")

		cmd := SetPermissionsCommand{
			RoleID:        1,
			PermissionIDs: []uint{10},
		}

		err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to set permissions")
	})

	t.Run("设置空权限列表（清空权限）", func(t *testing.T) {
		setup := setupSetPermissionsTest()
		setup.queryRepo.existingRoles[1] = true

		cmd := SetPermissionsCommand{
			RoleID:        1,
			PermissionIDs: []uint{},
		}

		err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		assert.True(t, setup.commandRepo.setPermCalled)
		assert.Empty(t, setup.commandRepo.setPermIDs)
	})
}

func TestNewSetPermissionsHandler(t *testing.T) {
	commandRepo := newSetPermissionsCommandRepo()
	queryRepo := newSetPermissionsQueryRepo()
	permissionQuery := newMockPermissionQueryRepo()

	handler := NewSetPermissionsHandler(commandRepo, queryRepo, permissionQuery)

	assert.NotNil(t, handler)
}
