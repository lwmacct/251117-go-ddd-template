package query

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
	permissions     []role.Permission
	permission      *role.Permission
	total           int64
	listErr         error
	findByIDErr     error
	findByCodeErr   error
	findByIDsErr    error
	listByResErr    error
	existsErr       error
	existsByCodeErr error
}

func newMockPermissionQueryRepo() *mockPermissionQueryRepo {
	return &mockPermissionQueryRepo{}
}

func (m *mockPermissionQueryRepo) FindByID(ctx context.Context, id uint) (*role.Permission, error) {
	if m.findByIDErr != nil {
		return nil, m.findByIDErr
	}
	return m.permission, nil
}

func (m *mockPermissionQueryRepo) FindByCode(ctx context.Context, code string) (*role.Permission, error) {
	if m.findByCodeErr != nil {
		return nil, m.findByCodeErr
	}
	return m.permission, nil
}

func (m *mockPermissionQueryRepo) FindByIDs(ctx context.Context, ids []uint) ([]role.Permission, error) {
	if m.findByIDsErr != nil {
		return nil, m.findByIDsErr
	}
	return m.permissions, nil
}

func (m *mockPermissionQueryRepo) List(ctx context.Context, page, limit int) ([]role.Permission, int64, error) {
	if m.listErr != nil {
		return nil, 0, m.listErr
	}
	return m.permissions, m.total, nil
}

func (m *mockPermissionQueryRepo) ListByResource(ctx context.Context, resource string) ([]role.Permission, error) {
	if m.listByResErr != nil {
		return nil, m.listByResErr
	}
	return m.permissions, nil
}

func (m *mockPermissionQueryRepo) Exists(ctx context.Context, id uint) (bool, error) {
	if m.existsErr != nil {
		return false, m.existsErr
	}
	return m.permission != nil, nil
}

func (m *mockPermissionQueryRepo) ExistsByCode(ctx context.Context, code string) (bool, error) {
	if m.existsByCodeErr != nil {
		return false, m.existsByCodeErr
	}
	return m.permission != nil, nil
}

type listPermissionsTestSetup struct {
	handler   *ListPermissionsHandler
	queryRepo *mockPermissionQueryRepo
}

func setupListPermissionsTest() *listPermissionsTestSetup {
	queryRepo := newMockPermissionQueryRepo()
	handler := NewListPermissionsHandler(queryRepo)

	return &listPermissionsTestSetup{
		handler:   handler,
		queryRepo: queryRepo,
	}
}

func TestListPermissionsHandler_Handle(t *testing.T) {
	ctx := context.Background()

	t.Run("成功获取权限列表", func(t *testing.T) {
		setup := setupListPermissionsTest()
		setup.queryRepo.permissions = []role.Permission{
			{ID: 1, Code: "user:read", Description: "读取用户", Resource: "user", Action: "read"},
			{ID: 2, Code: "user:write", Description: "写入用户", Resource: "user", Action: "write"},
			{ID: 3, Code: "role:read", Description: "读取角色", Resource: "role", Action: "read"},
		}
		setup.queryRepo.total = 3

		query := ListPermissionsQuery{
			Page:  1,
			Limit: 10,
		}

		result, err := setup.handler.Handle(ctx, query)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result.Permissions, 3)
		assert.Equal(t, int64(3), result.Total)
		assert.Equal(t, 1, result.Page)
		assert.Equal(t, 10, result.Limit)
	})

	t.Run("分页查询", func(t *testing.T) {
		setup := setupListPermissionsTest()
		setup.queryRepo.permissions = []role.Permission{
			{ID: 6, Code: "menu:read", Description: "读取菜单"},
		}
		setup.queryRepo.total = 10

		query := ListPermissionsQuery{
			Page:  2,
			Limit: 5,
		}

		result, err := setup.handler.Handle(ctx, query)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result.Permissions, 1)
		assert.Equal(t, int64(10), result.Total)
		assert.Equal(t, 2, result.Page)
		assert.Equal(t, 5, result.Limit)
	})

	t.Run("空列表", func(t *testing.T) {
		setup := setupListPermissionsTest()
		setup.queryRepo.permissions = []role.Permission{}
		setup.queryRepo.total = 0

		query := ListPermissionsQuery{
			Page:  1,
			Limit: 10,
		}

		result, err := setup.handler.Handle(ctx, query)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Empty(t, result.Permissions)
		assert.Equal(t, int64(0), result.Total)
	})

	t.Run("查询失败", func(t *testing.T) {
		setup := setupListPermissionsTest()
		setup.queryRepo.listErr = errors.New("database error")

		query := ListPermissionsQuery{
			Page:  1,
			Limit: 10,
		}

		result, err := setup.handler.Handle(ctx, query)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to list permissions")
		assert.Nil(t, result)
	})

	t.Run("权限包含完整信息", func(t *testing.T) {
		setup := setupListPermissionsTest()
		setup.queryRepo.permissions = []role.Permission{
			{
				ID:          1,
				Code:        "user:create",
				Description: "创建用户",
				Resource:    "user",
				Action:      "create",
			},
		}
		setup.queryRepo.total = 1

		query := ListPermissionsQuery{
			Page:  1,
			Limit: 10,
		}

		result, err := setup.handler.Handle(ctx, query)

		require.NoError(t, err)
		require.Len(t, result.Permissions, 1)
		assert.Equal(t, "user:create", result.Permissions[0].Code)
		assert.Equal(t, "创建用户", result.Permissions[0].Description)
		assert.Equal(t, "user", result.Permissions[0].Resource)
		assert.Equal(t, "create", result.Permissions[0].Action)
	})
}

func TestNewListPermissionsHandler(t *testing.T) {
	queryRepo := newMockPermissionQueryRepo()
	handler := NewListPermissionsHandler(queryRepo)

	assert.NotNil(t, handler)
}
