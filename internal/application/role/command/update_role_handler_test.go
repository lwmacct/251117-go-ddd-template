package command

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/role"
)

// updateMockRoleQueryRepo 用于更新测试的查询仓储 Mock
type updateMockRoleQueryRepo struct {
	mockRoleQueryRepo

	findByIDErr error
}

func (m *updateMockRoleQueryRepo) FindByID(ctx context.Context, id uint) (*role.Role, error) {
	if m.findByIDErr != nil {
		return nil, m.findByIDErr
	}
	if r, ok := m.roles[id]; ok {
		return r, nil
	}
	return nil, nil //nolint:nilnil // test mock returns nil for not found
}

// updateMockRoleCommandRepo 用于更新测试的命令仓储 Mock
type updateMockRoleCommandRepo struct {
	mockRoleCommandRepo

	updateError error
	updatedRole *role.Role
}

func (m *updateMockRoleCommandRepo) Update(ctx context.Context, r *role.Role) error {
	if m.updateError != nil {
		return m.updateError
	}
	m.updatedRole = r
	return nil
}

type updateRoleTestSetup struct {
	handler     *UpdateRoleHandler
	commandRepo *updateMockRoleCommandRepo
	queryRepo   *updateMockRoleQueryRepo
}

func setupUpdateRoleTest() *updateRoleTestSetup {
	commandRepo := &updateMockRoleCommandRepo{mockRoleCommandRepo: *newMockRoleCommandRepo()}
	queryRepo := &updateMockRoleQueryRepo{mockRoleQueryRepo: *newMockRoleQueryRepo()}
	handler := NewUpdateRoleHandler(commandRepo, queryRepo)

	return &updateRoleTestSetup{
		handler:     handler,
		commandRepo: commandRepo,
		queryRepo:   queryRepo,
	}
}

func TestUpdateRoleHandler_Handle(t *testing.T) {
	ctx := context.Background()

	t.Run("成功更新角色", func(t *testing.T) {
		setup := setupUpdateRoleTest()
		setup.queryRepo.roles[1] = &role.Role{
			ID:          1,
			Name:        "editor",
			DisplayName: "Editor",
			Description: "Old description",
			IsSystem:    false,
		}

		newDisplayName := "Senior Editor"
		newDescription := "New description"
		cmd := UpdateRoleCommand{
			RoleID:      1,
			DisplayName: &newDisplayName,
			Description: &newDescription,
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, uint(1), result.RoleID)
		assert.Equal(t, "Senior Editor", result.DisplayName)
		assert.Equal(t, "New description", result.Description)
	})

	t.Run("角色不存在", func(t *testing.T) {
		setup := setupUpdateRoleTest()

		newDisplayName := "New Name"
		cmd := UpdateRoleCommand{
			RoleID:      999,
			DisplayName: &newDisplayName,
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "role not found")
		assert.Nil(t, result)
	})

	t.Run("查询角色失败", func(t *testing.T) {
		setup := setupUpdateRoleTest()
		setup.queryRepo.findByIDErr = errors.New("database error")

		newDisplayName := "New Name"
		cmd := UpdateRoleCommand{
			RoleID:      1,
			DisplayName: &newDisplayName,
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to find role")
		assert.Nil(t, result)
	})

	t.Run("系统角色不可修改", func(t *testing.T) {
		setup := setupUpdateRoleTest()
		setup.queryRepo.roles[1] = &role.Role{
			ID:       1,
			Name:     "admin",
			IsSystem: true,
		}

		newDisplayName := "New Admin"
		cmd := UpdateRoleCommand{
			RoleID:      1,
			DisplayName: &newDisplayName,
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot modify system role")
		assert.Nil(t, result)
	})

	t.Run("更新失败", func(t *testing.T) {
		setup := setupUpdateRoleTest()
		setup.queryRepo.roles[1] = &role.Role{
			ID:       1,
			Name:     "editor",
			IsSystem: false,
		}
		setup.commandRepo.updateError = errors.New("database error")

		newDisplayName := "New Name"
		cmd := UpdateRoleCommand{
			RoleID:      1,
			DisplayName: &newDisplayName,
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to update role")
		assert.Nil(t, result)
	})

	t.Run("只更新部分字段", func(t *testing.T) {
		setup := setupUpdateRoleTest()
		setup.queryRepo.roles[1] = &role.Role{
			ID:          1,
			Name:        "editor",
			DisplayName: "Original",
			Description: "Original desc",
			IsSystem:    false,
		}

		newDisplayName := "Updated"
		cmd := UpdateRoleCommand{
			RoleID:      1,
			DisplayName: &newDisplayName,
			// Description 为 nil，不更新
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "Updated", result.DisplayName)
		assert.Equal(t, "Original desc", result.Description)
	})
}

func TestNewUpdateRoleHandler(t *testing.T) {
	commandRepo := newMockRoleCommandRepo()
	queryRepo := newMockRoleQueryRepo()

	handler := NewUpdateRoleHandler(commandRepo, queryRepo)

	assert.NotNil(t, handler)
}
