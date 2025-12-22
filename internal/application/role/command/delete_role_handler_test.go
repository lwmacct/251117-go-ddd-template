package command

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/role"
)

// deleteMockRoleQueryRepo 用于删除测试的查询仓储 Mock
type deleteMockRoleQueryRepo struct {
	mockRoleQueryRepo

	findByIDErr error
}

func (m *deleteMockRoleQueryRepo) FindByID(ctx context.Context, id uint) (*role.Role, error) {
	if m.findByIDErr != nil {
		return nil, m.findByIDErr
	}
	if r, ok := m.roles[id]; ok {
		return r, nil
	}
	return nil, nil //nolint:nilnil // test mock returns nil for not found
}

// deleteMockRoleCommandRepo 用于删除测试的命令仓储 Mock
type deleteMockRoleCommandRepo struct {
	mockRoleCommandRepo

	deleteError   error
	deletedRoleID uint
}

func (m *deleteMockRoleCommandRepo) Delete(ctx context.Context, id uint) error {
	if m.deleteError != nil {
		return m.deleteError
	}
	m.deletedRoleID = id
	return nil
}

type deleteRoleTestSetup struct {
	handler     *DeleteRoleHandler
	commandRepo *deleteMockRoleCommandRepo
	queryRepo   *deleteMockRoleQueryRepo
}

func setupDeleteRoleTest() *deleteRoleTestSetup {
	commandRepo := &deleteMockRoleCommandRepo{mockRoleCommandRepo: *newMockRoleCommandRepo()}
	queryRepo := &deleteMockRoleQueryRepo{mockRoleQueryRepo: *newMockRoleQueryRepo()}
	handler := NewDeleteRoleHandler(commandRepo, queryRepo)

	return &deleteRoleTestSetup{
		handler:     handler,
		commandRepo: commandRepo,
		queryRepo:   queryRepo,
	}
}

func TestDeleteRoleHandler_Handle(t *testing.T) {
	ctx := context.Background()

	t.Run("成功删除角色", func(t *testing.T) {
		setup := setupDeleteRoleTest()
		setup.queryRepo.roles[1] = &role.Role{
			ID:       1,
			Name:     "editor",
			IsSystem: false,
		}

		cmd := DeleteRoleCommand{
			RoleID: 1,
		}

		err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		assert.Equal(t, uint(1), setup.commandRepo.deletedRoleID)
	})

	t.Run("角色不存在", func(t *testing.T) {
		setup := setupDeleteRoleTest()

		cmd := DeleteRoleCommand{
			RoleID: 999,
		}

		err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "role not found")
	})

	t.Run("查询角色失败", func(t *testing.T) {
		setup := setupDeleteRoleTest()
		setup.queryRepo.findByIDErr = errors.New("database error")

		cmd := DeleteRoleCommand{
			RoleID: 1,
		}

		err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to find role")
	})

	t.Run("系统角色不可删除", func(t *testing.T) {
		setup := setupDeleteRoleTest()
		setup.queryRepo.roles[1] = &role.Role{
			ID:       1,
			Name:     "admin",
			IsSystem: true,
		}

		cmd := DeleteRoleCommand{
			RoleID: 1,
		}

		err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot delete system role")
	})

	t.Run("删除失败", func(t *testing.T) {
		setup := setupDeleteRoleTest()
		setup.queryRepo.roles[1] = &role.Role{
			ID:       1,
			Name:     "editor",
			IsSystem: false,
		}
		setup.commandRepo.deleteError = errors.New("database error")

		cmd := DeleteRoleCommand{
			RoleID: 1,
		}

		err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to delete role")
	})
}

func TestNewDeleteRoleHandler(t *testing.T) {
	commandRepo := newMockRoleCommandRepo()
	queryRepo := newMockRoleQueryRepo()

	handler := NewDeleteRoleHandler(commandRepo, queryRepo)

	assert.NotNil(t, handler)
}
