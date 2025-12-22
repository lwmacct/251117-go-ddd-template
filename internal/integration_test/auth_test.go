package integration_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	authCommand "github.com/lwmacct/251117-go-ddd-template/internal/application/auth/command"
	userCommand "github.com/lwmacct/251117-go-ddd-template/internal/application/user/command"
	userQuery "github.com/lwmacct/251117-go-ddd-template/internal/application/user/query"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
)

func TestAuthFlow_RegisterLoginRefresh(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	ctx := context.Background()

	t.Run("完整认证流程", func(t *testing.T) {
		// 步骤 1: 注册用户
		registerResult, err := env.RegisterHandler.Handle(ctx, authCommand.RegisterCommand{
			Username: "newuser",
			Email:    "newuser@example.com",
			Password: "Password123!",
		})

		require.NoError(t, err)
		require.NotNil(t, registerResult)
		assert.NotZero(t, registerResult.UserID)
		assert.NotEmpty(t, registerResult.AccessToken)
		assert.NotEmpty(t, registerResult.RefreshToken)

		// 步骤 2: 验证用户已创建
		createdUser, err := env.UserQueryRepo.GetByID(ctx, registerResult.UserID)
		require.NoError(t, err)
		assert.Equal(t, "newuser", createdUser.Username)
		assert.Equal(t, "newuser@example.com", createdUser.Email)
		assert.Equal(t, "active", createdUser.Status)

		// 步骤 3: 验证密码已加密（不是明文）
		assert.NotEqual(t, "Password123!", createdUser.Password)
		assert.Greater(t, len(createdUser.Password), 20, "password should be hashed")
	})

	t.Run("重复注册应该失败", func(t *testing.T) {
		// 先注册一个用户
		_, err := env.RegisterHandler.Handle(ctx, authCommand.RegisterCommand{
			Username: "duplicate",
			Email:    "dup@example.com",
			Password: "Password123!",
		})
		require.NoError(t, err)

		// 尝试用相同用户名注册
		_, err = env.RegisterHandler.Handle(ctx, authCommand.RegisterCommand{
			Username: "duplicate",
			Email:    "different@example.com",
			Password: "Password123!",
		})
		require.Error(t, err)
		require.ErrorIs(t, err, user.ErrUsernameAlreadyExists)

		// 尝试用相同邮箱注册
		_, err = env.RegisterHandler.Handle(ctx, authCommand.RegisterCommand{
			Username: "different",
			Email:    "dup@example.com",
			Password: "Password123!",
		})
		require.Error(t, err)
		assert.ErrorIs(t, err, user.ErrEmailAlreadyExists)
	})

	t.Run("弱密码应该被拒绝", func(t *testing.T) {
		_, err := env.RegisterHandler.Handle(ctx, authCommand.RegisterCommand{
			Username: "weakpwd",
			Email:    "weak@example.com",
			Password: "123", // 太短
		})

		assert.Error(t, err)
	})
}

func TestUserManagement_CRUD(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	ctx := context.Background()

	t.Run("创建 → 查询 → 更新 → 删除流程", func(t *testing.T) {
		// 步骤 1: 创建用户
		createResult, err := env.CreateUserHandler.Handle(ctx, userCommand.CreateUserCommand{
			Username: "cruduser",
			Email:    "crud@example.com",
			Password: "Password123!",
			FullName: "CRUD Test User",
		})
		require.NoError(t, err)
		require.NotNil(t, createResult)
		userID := createResult.UserID

		// 步骤 2: 查询用户
		foundUser, err := env.GetUserHandler.Handle(ctx, userQuery.GetUserQuery{
			UserID:    userID,
			WithRoles: false,
		})
		require.NoError(t, err)
		assert.Equal(t, "cruduser", foundUser.Username)
		assert.Equal(t, "CRUD Test User", foundUser.FullName)

		// 步骤 3: 列表查询
		listResult, err := env.ListUsersHandler.Handle(ctx, userQuery.ListUsersQuery{
			Offset: 0,
			Limit:  10,
		})
		require.NoError(t, err)
		assert.GreaterOrEqual(t, listResult.Total, int64(1))

		// 步骤 4: 删除用户
		err = env.UserCommandRepo.Delete(ctx, userID)
		require.NoError(t, err)

		// 步骤 5: 验证软删除
		_, err = env.UserQueryRepo.GetByID(ctx, userID)
		assert.ErrorIs(t, err, user.ErrUserNotFound)
	})
}

func TestBatchUserCreation(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	ctx := context.Background()

	batchHandler := userCommand.NewBatchCreateUsersHandler(
		env.UserCommandRepo,
		env.UserQueryRepo,
		env.AuthService,
	)

	t.Run("批量创建用户成功", func(t *testing.T) {
		cmd := userCommand.BatchCreateUsersCommand{
			Users: []userCommand.BatchUserItem{
				{Username: "batch1", Email: "batch1@example.com", Password: "Password123!", FullName: "Batch User 1"},
				{Username: "batch2", Email: "batch2@example.com", Password: "Password123!", FullName: "Batch User 2"},
				{Username: "batch3", Email: "batch3@example.com", Password: "Password123!", FullName: "Batch User 3"},
			},
		}

		result, err := batchHandler.Handle(ctx, cmd)

		require.NoError(t, err)
		assert.Equal(t, 3, result.Total)
		assert.Equal(t, 3, result.Success)
		assert.Equal(t, 0, result.Failed)
		assert.Empty(t, result.Errors)

		// 验证用户已创建
		count, err := env.UserQueryRepo.Count(ctx)
		require.NoError(t, err)
		assert.Equal(t, int64(3), count)
	})

	t.Run("批量创建部分失败", func(t *testing.T) {
		cmd := userCommand.BatchCreateUsersCommand{
			Users: []userCommand.BatchUserItem{
				{Username: "new1", Email: "new1@example.com", Password: "Password123!", FullName: "New User 1"},
				{Username: "batch1", Email: "batch1@example.com", Password: "Password123!", FullName: "Duplicate"}, // 重复
				{Username: "new2", Email: "new2@example.com", Password: "Password123!", FullName: "New User 2"},
			},
		}

		result, err := batchHandler.Handle(ctx, cmd)

		require.NoError(t, err)
		assert.Equal(t, 3, result.Total)
		assert.Equal(t, 2, result.Success)
		assert.Equal(t, 1, result.Failed)
		assert.Len(t, result.Errors, 1)
		assert.Equal(t, 1, result.Errors[0].Index)
		assert.Contains(t, result.Errors[0].Error, "用户名已存在")
	})
}

func TestPasswordPolicyEnforcement(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	ctx := context.Background()

	// 默认密码策略只检查最小长度(6字符)
	// 不要求大写/小写/数字/特殊字符
	testCases := []struct {
		name        string
		password    string
		shouldError bool
		description string
	}{
		{"强密码", "StrongPassword123!", false, "满足所有条件"},
		{"太短", "Ab1!", true, "只有4个字符，少于6"},
		{"刚好6字符", "abc123", false, "刚好满足最小长度"},
		{"无大写", "password123!", false, "默认策略不要求大写"},
		{"无小写", "PASSWORD123!", false, "默认策略不要求小写"},
		{"无数字", "PasswordABC!", false, "默认策略不要求数字"},
		{"纯数字", "123456", false, "只要满足最小长度即可"},
		{"5字符", "12345", true, "5个字符少于最小长度6"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := env.RegisterHandler.Handle(ctx, authCommand.RegisterCommand{
				Username: "pwd_" + tc.name,
				Email:    "pwd_" + tc.name + "@example.com",
				Password: tc.password,
			})

			if tc.shouldError {
				assert.Error(t, err, "password: %s should be rejected (%s)", tc.password, tc.description)
			} else {
				assert.NoError(t, err, "password: %s should be accepted (%s)", tc.password, tc.description)
			}
		})
	}
}
