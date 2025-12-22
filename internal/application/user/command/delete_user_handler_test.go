package command

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
)

// deleteMockUserQueryRepo 用于删除测试的查询仓储 Mock
type deleteMockUserQueryRepo struct {
	mockUserQueryRepo

	existsResult bool
	existsError  error
}

func (m *deleteMockUserQueryRepo) Exists(ctx context.Context, id uint) (bool, error) {
	if m.existsError != nil {
		return false, m.existsError
	}
	return m.existsResult, nil
}

// deleteMockUserCommandRepo 用于删除测试的命令仓储 Mock
type deleteMockUserCommandRepo struct {
	mockUserCommandRepo

	deleteError   error
	deletedUserID uint
}

func (m *deleteMockUserCommandRepo) Delete(ctx context.Context, id uint) error {
	if m.deleteError != nil {
		return m.deleteError
	}
	m.deletedUserID = id
	return nil
}

type deleteUserTestSetup struct {
	handler     *DeleteUserHandler
	commandRepo *deleteMockUserCommandRepo
	queryRepo   *deleteMockUserQueryRepo
}

func setupDeleteUserTest() *deleteUserTestSetup {
	commandRepo := &deleteMockUserCommandRepo{mockUserCommandRepo: *newMockUserCommandRepo()}
	queryRepo := &deleteMockUserQueryRepo{mockUserQueryRepo: *newMockUserQueryRepo()}

	handler := NewDeleteUserHandler(commandRepo, queryRepo)

	return &deleteUserTestSetup{
		handler:     handler,
		commandRepo: commandRepo,
		queryRepo:   queryRepo,
	}
}

func TestDeleteUserHandler_Handle(t *testing.T) {
	ctx := context.Background()

	t.Run("成功删除用户", func(t *testing.T) {
		setup := setupDeleteUserTest()
		setup.queryRepo.existsResult = true

		cmd := DeleteUserCommand{
			UserID: 1,
		}

		err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		assert.Equal(t, uint(1), setup.commandRepo.deletedUserID)
	})

	t.Run("用户不存在", func(t *testing.T) {
		setup := setupDeleteUserTest()
		setup.queryRepo.existsResult = false

		cmd := DeleteUserCommand{
			UserID: 999,
		}

		err := setup.handler.Handle(ctx, cmd)

		assert.ErrorIs(t, err, user.ErrUserNotFound)
	})

	t.Run("检查用户存在性失败", func(t *testing.T) {
		setup := setupDeleteUserTest()
		setup.queryRepo.existsError = errors.New("database error")

		cmd := DeleteUserCommand{
			UserID: 1,
		}

		err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to check user existence")
	})

	t.Run("删除操作失败", func(t *testing.T) {
		setup := setupDeleteUserTest()
		setup.queryRepo.existsResult = true
		setup.commandRepo.deleteError = errors.New("database error")

		cmd := DeleteUserCommand{
			UserID: 1,
		}

		err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to delete user")
	})
}

func TestNewDeleteUserHandler(t *testing.T) {
	t.Run("创建 DeleteUserHandler", func(t *testing.T) {
		commandRepo := newMockUserCommandRepo()
		queryRepo := newMockUserQueryRepo()

		handler := NewDeleteUserHandler(commandRepo, queryRepo)

		assert.NotNil(t, handler)
	})
}
