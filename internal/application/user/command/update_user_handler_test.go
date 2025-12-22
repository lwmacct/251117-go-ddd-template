package command

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
)

// updateMockUserQueryRepo 用于更新测试的查询仓储 Mock
type updateMockUserQueryRepo struct {
	mockUserQueryRepo

	user       *user.User
	getUserErr error
}

func (m *updateMockUserQueryRepo) GetByID(ctx context.Context, id uint) (*user.User, error) {
	if m.getUserErr != nil {
		return nil, m.getUserErr
	}
	if m.user != nil && m.user.ID == id {
		return m.user, nil
	}
	return nil, user.ErrUserNotFound
}

// updateMockUserCommandRepo 用于更新测试的命令仓储 Mock
type updateMockUserCommandRepo struct {
	mockUserCommandRepo

	updateError   error
	updatedUser   *user.User
	updateCounter int
}

func (m *updateMockUserCommandRepo) Update(ctx context.Context, u *user.User) error {
	if m.updateError != nil {
		return m.updateError
	}
	m.updateCounter++
	m.updatedUser = u
	return nil
}

type updateUserTestSetup struct {
	handler     *UpdateUserHandler
	commandRepo *updateMockUserCommandRepo
	queryRepo   *updateMockUserQueryRepo
}

func setupUpdateUserTest() *updateUserTestSetup {
	commandRepo := &updateMockUserCommandRepo{mockUserCommandRepo: *newMockUserCommandRepo()}
	queryRepo := &updateMockUserQueryRepo{mockUserQueryRepo: *newMockUserQueryRepo()}

	handler := NewUpdateUserHandler(commandRepo, queryRepo)

	return &updateUserTestSetup{
		handler:     handler,
		commandRepo: commandRepo,
		queryRepo:   queryRepo,
	}
}

func TestUpdateUserHandler_Handle(t *testing.T) {
	ctx := context.Background()

	t.Run("成功更新用户 - 更新 FullName", func(t *testing.T) {
		setup := setupUpdateUserTest()
		setup.queryRepo.user = &user.User{
			ID:       1,
			Username: "testuser",
			FullName: "Old Name",
			Status:   "active",
		}

		newName := "New Name"
		cmd := UpdateUserCommand{
			UserID:   1,
			FullName: &newName,
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, uint(1), result.UserID)
		assert.Equal(t, "New Name", setup.commandRepo.updatedUser.FullName)
	})

	t.Run("成功更新用户 - 更新多个字段", func(t *testing.T) {
		setup := setupUpdateUserTest()
		setup.queryRepo.user = &user.User{
			ID:       1,
			Username: "testuser",
			FullName: "Old Name",
			Avatar:   "old_avatar.jpg",
			Bio:      "Old bio",
			Status:   "active",
		}

		newName := "New Name"
		newAvatar := "new_avatar.jpg"
		newBio := "New bio"
		cmd := UpdateUserCommand{
			UserID:   1,
			FullName: &newName,
			Avatar:   &newAvatar,
			Bio:      &newBio,
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "New Name", setup.commandRepo.updatedUser.FullName)
		assert.Equal(t, "new_avatar.jpg", setup.commandRepo.updatedUser.Avatar)
		assert.Equal(t, "New bio", setup.commandRepo.updatedUser.Bio)
	})

	t.Run("成功更新用户状态 - active", func(t *testing.T) {
		setup := setupUpdateUserTest()
		setup.queryRepo.user = &user.User{
			ID:     1,
			Status: "inactive",
		}

		status := "active"
		cmd := UpdateUserCommand{
			UserID: 1,
			Status: &status,
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "active", setup.commandRepo.updatedUser.Status)
	})

	t.Run("成功更新用户状态 - inactive", func(t *testing.T) {
		setup := setupUpdateUserTest()
		setup.queryRepo.user = &user.User{
			ID:     1,
			Status: "active",
		}

		status := "inactive"
		cmd := UpdateUserCommand{
			UserID: 1,
			Status: &status,
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "inactive", setup.commandRepo.updatedUser.Status)
	})

	t.Run("成功更新用户状态 - banned", func(t *testing.T) {
		setup := setupUpdateUserTest()
		setup.queryRepo.user = &user.User{
			ID:     1,
			Status: "active",
		}

		status := "banned"
		cmd := UpdateUserCommand{
			UserID: 1,
			Status: &status,
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "banned", setup.commandRepo.updatedUser.Status)
	})

	t.Run("用户不存在", func(t *testing.T) {
		setup := setupUpdateUserTest()
		setup.queryRepo.getUserErr = user.ErrUserNotFound

		newName := "New Name"
		cmd := UpdateUserCommand{
			UserID:   999,
			FullName: &newName,
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.ErrorIs(t, err, user.ErrUserNotFound)
		assert.Nil(t, result)
	})

	t.Run("无效的状态值", func(t *testing.T) {
		setup := setupUpdateUserTest()
		setup.queryRepo.user = &user.User{
			ID:     1,
			Status: "active",
		}

		invalidStatus := "invalid_status"
		cmd := UpdateUserCommand{
			UserID: 1,
			Status: &invalidStatus,
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.ErrorIs(t, err, user.ErrInvalidUserStatus)
		assert.Nil(t, result)
	})

	t.Run("更新失败", func(t *testing.T) {
		setup := setupUpdateUserTest()
		setup.queryRepo.user = &user.User{
			ID:     1,
			Status: "active",
		}
		setup.commandRepo.updateError = errors.New("database error")

		newName := "New Name"
		cmd := UpdateUserCommand{
			UserID:   1,
			FullName: &newName,
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to update user")
		assert.Nil(t, result)
	})

	t.Run("空更新 - 不修改任何字段", func(t *testing.T) {
		setup := setupUpdateUserTest()
		setup.queryRepo.user = &user.User{
			ID:       1,
			FullName: "Original Name",
			Status:   "active",
		}

		cmd := UpdateUserCommand{
			UserID: 1,
			// 所有字段都是 nil
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		require.NotNil(t, result)
		// 应该仍然调用了 Update（即使没有修改）
		assert.Equal(t, 1, setup.commandRepo.updateCounter)
		assert.Equal(t, "Original Name", setup.commandRepo.updatedUser.FullName)
	})
}

func TestNewUpdateUserHandler(t *testing.T) {
	t.Run("创建 UpdateUserHandler", func(t *testing.T) {
		commandRepo := newMockUserCommandRepo()
		queryRepo := newMockUserQueryRepo()

		handler := NewUpdateUserHandler(commandRepo, queryRepo)

		assert.NotNil(t, handler)
	})
}
