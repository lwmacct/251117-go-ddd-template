package query

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/role"
)

type listRolesTestSetup struct {
	handler   *ListRolesHandler
	queryRepo *mockRoleQueryRepo
}

func setupListRolesTest() *listRolesTestSetup {
	queryRepo := newMockRoleQueryRepo()
	handler := NewListRolesHandler(queryRepo)

	return &listRolesTestSetup{
		handler:   handler,
		queryRepo: queryRepo,
	}
}

func TestListRolesHandler_Handle(t *testing.T) {
	ctx := context.Background()

	t.Run("成功获取角色列表", func(t *testing.T) {
		setup := setupListRolesTest()
		setup.queryRepo.listRoles = []role.Role{
			{ID: 1, Name: "admin", DisplayName: "Administrator", IsSystem: true},
			{ID: 2, Name: "editor", DisplayName: "Editor", IsSystem: false},
			{ID: 3, Name: "viewer", DisplayName: "Viewer", IsSystem: false},
		}
		setup.queryRepo.listTotal = 3

		query := ListRolesQuery{
			Page:  1,
			Limit: 10,
		}

		result, err := setup.handler.Handle(ctx, query)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result.Roles, 3)
		assert.Equal(t, int64(3), result.Total)
		assert.Equal(t, 1, result.Page)
		assert.Equal(t, 10, result.Limit)
	})

	t.Run("分页", func(t *testing.T) {
		setup := setupListRolesTest()
		setup.queryRepo.listRoles = []role.Role{
			{ID: 3, Name: "viewer", DisplayName: "Viewer"},
		}
		setup.queryRepo.listTotal = 5

		query := ListRolesQuery{
			Page:  2,
			Limit: 2,
		}

		result, err := setup.handler.Handle(ctx, query)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result.Roles, 1)
		assert.Equal(t, int64(5), result.Total)
		assert.Equal(t, 2, result.Page)
		assert.Equal(t, 2, result.Limit)
	})

	t.Run("空列表", func(t *testing.T) {
		setup := setupListRolesTest()
		setup.queryRepo.listRoles = []role.Role{}
		setup.queryRepo.listTotal = 0

		query := ListRolesQuery{
			Page:  1,
			Limit: 10,
		}

		result, err := setup.handler.Handle(ctx, query)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Empty(t, result.Roles)
		assert.Equal(t, int64(0), result.Total)
	})

	t.Run("查询失败", func(t *testing.T) {
		setup := setupListRolesTest()
		setup.queryRepo.listErr = errors.New("database error")

		query := ListRolesQuery{
			Page:  1,
			Limit: 10,
		}

		result, err := setup.handler.Handle(ctx, query)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to list roles")
		assert.Nil(t, result)
	})

	t.Run("角色带权限", func(t *testing.T) {
		setup := setupListRolesTest()
		setup.queryRepo.listRoles = []role.Role{
			{
				ID:          1,
				Name:        "editor",
				DisplayName: "Editor",
				Permissions: []role.Permission{
					{ID: 1, Code: "article:read"},
					{ID: 2, Code: "article:write"},
				},
			},
		}
		setup.queryRepo.listTotal = 1

		query := ListRolesQuery{
			Page:  1,
			Limit: 10,
		}

		result, err := setup.handler.Handle(ctx, query)

		require.NoError(t, err)
		require.NotNil(t, result)
		require.Len(t, result.Roles, 1)
		assert.Len(t, result.Roles[0].Permissions, 2)
	})
}

func TestNewListRolesHandler(t *testing.T) {
	queryRepo := newMockRoleQueryRepo()
	handler := NewListRolesHandler(queryRepo)

	assert.NotNil(t, handler)
}
