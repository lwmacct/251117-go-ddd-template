package command

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// batchTestSetup 批量创建测试设置
type batchTestSetup struct {
	handler     *BatchCreateUsersHandler
	commandRepo *mockUserCommandRepo
	queryRepo   *mockUserQueryRepo
	authService *mockAuthService
}

func setupBatchTest() *batchTestSetup {
	commandRepo := newMockUserCommandRepo()
	queryRepo := newMockUserQueryRepo()
	authService := newMockAuthService()

	handler := NewBatchCreateUsersHandler(commandRepo, queryRepo, authService)

	return &batchTestSetup{
		handler:     handler,
		commandRepo: commandRepo,
		queryRepo:   queryRepo,
		authService: authService,
	}
}

func TestBatchCreateUsersHandler_Handle(t *testing.T) { //nolint:maintidx // 测试函数的复杂度可接受
	ctx := context.Background()

	t.Run("成功创建所有用户", func(t *testing.T) {
		setup := setupBatchTest()

		cmd := BatchCreateUsersCommand{
			Users: []BatchUserItem{
				{Username: "user1", Email: "user1@example.com", Password: "Password123!", FullName: "User One"},
				{Username: "user2", Email: "user2@example.com", Password: "Password123!", FullName: "User Two"},
			},
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		assert.Equal(t, 2, result.Total)
		assert.Equal(t, 2, result.Success)
		assert.Equal(t, 0, result.Failed)
		assert.Empty(t, result.Errors)

		// 验证用户已创建
		assert.Len(t, setup.commandRepo.users, 2)
	})

	t.Run("部分用户创建失败 - 用户名已存在", func(t *testing.T) {
		setup := setupBatchTest()
		setup.queryRepo.existingUsernames["existing_user"] = true

		cmd := BatchCreateUsersCommand{
			Users: []BatchUserItem{
				{Username: "existing_user", Email: "new1@example.com", Password: "Password123!", FullName: "New User"},
				{Username: "new_user", Email: "new2@example.com", Password: "Password123!", FullName: "New User"},
			},
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		assert.Equal(t, 2, result.Total)
		assert.Equal(t, 1, result.Success)
		assert.Equal(t, 1, result.Failed)
		assert.Len(t, result.Errors, 1)
		assert.Equal(t, 0, result.Errors[0].Index)
		assert.Equal(t, "existing_user", result.Errors[0].Username)
		assert.Contains(t, result.Errors[0].Error, "用户名已存在")
	})

	t.Run("部分用户创建失败 - 邮箱已存在", func(t *testing.T) {
		setup := setupBatchTest()
		setup.queryRepo.existingEmails["existing@example.com"] = true

		cmd := BatchCreateUsersCommand{
			Users: []BatchUserItem{
				{Username: "user1", Email: "existing@example.com", Password: "Password123!", FullName: "User One"},
				{Username: "user2", Email: "new@example.com", Password: "Password123!", FullName: "User Two"},
			},
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		assert.Equal(t, 2, result.Total)
		assert.Equal(t, 1, result.Success)
		assert.Equal(t, 1, result.Failed)
		assert.Len(t, result.Errors, 1)
		assert.Contains(t, result.Errors[0].Error, "邮箱已存在")
	})

	t.Run("用户名验证失败 - 太短", func(t *testing.T) {
		setup := setupBatchTest()

		cmd := BatchCreateUsersCommand{
			Users: []BatchUserItem{
				{Username: "ab", Email: "user@example.com", Password: "Password123!", FullName: "User"},
			},
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		assert.Equal(t, 1, result.Total)
		assert.Equal(t, 0, result.Success)
		assert.Equal(t, 1, result.Failed)
		assert.Contains(t, result.Errors[0].Error, "用户名长度必须在 3-50 个字符之间")
	})

	t.Run("邮箱格式无效", func(t *testing.T) {
		setup := setupBatchTest()

		cmd := BatchCreateUsersCommand{
			Users: []BatchUserItem{
				{Username: "user1", Email: "invalid-email", Password: "Password123!", FullName: "User"},
			},
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		assert.Equal(t, 1, result.Total)
		assert.Equal(t, 0, result.Success)
		assert.Equal(t, 1, result.Failed)
		assert.Contains(t, result.Errors[0].Error, "邮箱格式无效")
	})

	t.Run("密码策略验证失败", func(t *testing.T) {
		setup := setupBatchTest()
		setup.authService.weakPassword = true

		cmd := BatchCreateUsersCommand{
			Users: []BatchUserItem{
				{Username: "user1", Email: "user@example.com", Password: "weak", FullName: "User"},
			},
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		assert.Equal(t, 1, result.Total)
		assert.Equal(t, 0, result.Success)
		assert.Equal(t, 1, result.Failed)
		assert.Contains(t, result.Errors[0].Error, "密码不符合策略")
	})

	t.Run("空用户列表", func(t *testing.T) {
		setup := setupBatchTest()

		cmd := BatchCreateUsersCommand{
			Users: []BatchUserItem{},
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		assert.Equal(t, 0, result.Total)
		assert.Equal(t, 0, result.Success)
		assert.Equal(t, 0, result.Failed)
		assert.Empty(t, result.Errors)
	})

	t.Run("带角色分配的用户创建", func(t *testing.T) {
		setup := setupBatchTest()

		cmd := BatchCreateUsersCommand{
			Users: []BatchUserItem{
				{Username: "user1", Email: "user1@example.com", Password: "Password123!", FullName: "User One", RoleIDs: []uint{1, 2}},
			},
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		assert.Equal(t, 1, result.Total)
		assert.Equal(t, 1, result.Success)
		assert.Equal(t, 0, result.Failed)

		// 验证角色分配
		assignedRoles, exists := setup.commandRepo.assignedRoles[uint(1)]
		assert.True(t, exists)
		assert.Equal(t, []uint{1, 2}, assignedRoles)
	})

	t.Run("数据库创建失败", func(t *testing.T) {
		setup := setupBatchTest()
		setup.commandRepo.createError = errors.New("database error")

		cmd := BatchCreateUsersCommand{
			Users: []BatchUserItem{
				{Username: "user1", Email: "user1@example.com", Password: "Password123!", FullName: "User One"},
			},
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		assert.Equal(t, 1, result.Total)
		assert.Equal(t, 0, result.Success)
		assert.Equal(t, 1, result.Failed)
		assert.Contains(t, result.Errors[0].Error, "创建用户失败")
	})

	t.Run("无效状态值", func(t *testing.T) {
		setup := setupBatchTest()

		cmd := BatchCreateUsersCommand{
			Users: []BatchUserItem{
				{Username: "user1", Email: "user1@example.com", Password: "Password123!", FullName: "User One", Status: "invalid"},
			},
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		assert.Equal(t, 1, result.Total)
		assert.Equal(t, 0, result.Success)
		assert.Equal(t, 1, result.Failed)
		assert.Contains(t, result.Errors[0].Error, "无效的状态值")
	})

	t.Run("混合结果场景", func(t *testing.T) {
		setup := setupBatchTest()
		setup.queryRepo.existingUsernames["user2"] = true

		cmd := BatchCreateUsersCommand{
			Users: []BatchUserItem{
				{Username: "user1", Email: "user1@example.com", Password: "Password123!", FullName: "User One"},
				{Username: "user2", Email: "user2@example.com", Password: "Password123!", FullName: "User Two"},   // 用户名已存在
				{Username: "user3", Email: "user3@example.com", Password: "Password123!", FullName: "User Three"}, // 成功
			},
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		assert.Equal(t, 3, result.Total)
		assert.Equal(t, 2, result.Success)
		assert.Equal(t, 1, result.Failed)
		assert.Len(t, result.Errors, 1)
		assert.Equal(t, 1, result.Errors[0].Index)
		assert.Equal(t, "user2", result.Errors[0].Username)
	})

	t.Run("默认状态为 active", func(t *testing.T) {
		setup := setupBatchTest()

		cmd := BatchCreateUsersCommand{
			Users: []BatchUserItem{
				{Username: "user1", Email: "user1@example.com", Password: "Password123!", FullName: "User One"},
			},
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		assert.Equal(t, 1, result.Success)

		// 验证状态
		require.Len(t, setup.commandRepo.users, 1)
		assert.Equal(t, "active", setup.commandRepo.users[0].Status)
	})

	t.Run("显式指定 inactive 状态", func(t *testing.T) {
		setup := setupBatchTest()

		cmd := BatchCreateUsersCommand{
			Users: []BatchUserItem{
				{Username: "user1", Email: "user1@example.com", Password: "Password123!", FullName: "User One", Status: "inactive"},
			},
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		assert.Equal(t, 1, result.Success)

		// 验证状态
		require.Len(t, setup.commandRepo.users, 1)
		assert.Equal(t, "inactive", setup.commandRepo.users[0].Status)
	})
}

func TestBatchCreateUsersHandler_ValidateUserItem(t *testing.T) {
	handler := &BatchCreateUsersHandler{}

	tests := []struct {
		name      string
		item      BatchUserItem
		expectErr bool
		errMsg    string
	}{
		{
			name:      "有效用户数据",
			item:      BatchUserItem{Username: "validuser", Email: "valid@example.com", Password: "Password123!"},
			expectErr: false,
		},
		{
			name:      "空用户名",
			item:      BatchUserItem{Username: "", Email: "valid@example.com", Password: "Password123!"},
			expectErr: true,
			errMsg:    "用户名不能为空",
		},
		{
			name:      "用户名太短",
			item:      BatchUserItem{Username: "ab", Email: "valid@example.com", Password: "Password123!"},
			expectErr: true,
			errMsg:    "用户名长度必须在 3-50 个字符之间",
		},
		{
			name:      "空邮箱",
			item:      BatchUserItem{Username: "validuser", Email: "", Password: "Password123!"},
			expectErr: true,
			errMsg:    "邮箱不能为空",
		},
		{
			name:      "无效邮箱格式",
			item:      BatchUserItem{Username: "validuser", Email: "invalid-email", Password: "Password123!"},
			expectErr: true,
			errMsg:    "邮箱格式无效",
		},
		{
			name:      "空密码",
			item:      BatchUserItem{Username: "validuser", Email: "valid@example.com", Password: ""},
			expectErr: true,
			errMsg:    "密码不能为空",
		},
		{
			name:      "只有空格的用户名",
			item:      BatchUserItem{Username: "   ", Email: "valid@example.com", Password: "Password123!"},
			expectErr: true,
			errMsg:    "用户名不能为空",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.validateUserItem(tt.item)
			if tt.expectErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNewBatchCreateUsersHandler(t *testing.T) {
	t.Run("创建 BatchCreateUsersHandler", func(t *testing.T) {
		commandRepo := newMockUserCommandRepo()
		queryRepo := newMockUserQueryRepo()
		authService := newMockAuthService()

		handler := NewBatchCreateUsersHandler(commandRepo, queryRepo, authService)

		assert.NotNil(t, handler)
	})
}

func TestBatchCreateUsersHandler_PasswordPolicyIntegration(t *testing.T) {
	ctx := context.Background()

	t.Run("密码策略验证在其他验证之后", func(t *testing.T) {
		setup := setupBatchTest()
		setup.authService.weakPassword = true

		// 即使密码弱，如果基本验证失败，应该先返回基本验证错误
		cmd := BatchCreateUsersCommand{
			Users: []BatchUserItem{
				{Username: "ab", Email: "user@example.com", Password: "weak", FullName: "User"}, // 用户名太短
			},
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		assert.Equal(t, 1, result.Failed)
		// 应该是用户名验证错误，而不是密码错误
		assert.Contains(t, result.Errors[0].Error, "用户名长度")
	})
}

func TestBatchCreateUsersHandler_RoleAssignmentFailure(t *testing.T) {
	ctx := context.Background()

	t.Run("角色分配失败不回滚用户创建", func(t *testing.T) {
		setup := setupBatchTest()
		setup.commandRepo.assignRolesError = errors.New("role assignment failed")

		cmd := BatchCreateUsersCommand{
			Users: []BatchUserItem{
				{Username: "user1", Email: "user1@example.com", Password: "Password123!", FullName: "User One", RoleIDs: []uint{1, 2}},
			},
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		// 用户创建成功但角色分配失败，应该记录为失败
		assert.Equal(t, 1, result.Total)
		assert.Equal(t, 0, result.Success)
		assert.Equal(t, 1, result.Failed)
		assert.Contains(t, result.Errors[0].Error, "用户已创建但角色分配失败")

		// 用户实际上已经创建了
		assert.Len(t, setup.commandRepo.users, 1)
	})
}

func BenchmarkBatchCreateUsersHandler_Handle(b *testing.B) {
	ctx := context.Background()
	setup := setupBatchTest()

	users := make([]BatchUserItem, 10)
	for i := range 10 {
		users[i] = BatchUserItem{
			Username: "user",
			Email:    "user@example.com",
			Password: "Password123!",
			FullName: "Test User",
		}
	}

	cmd := BatchCreateUsersCommand{Users: users}

	for b.Loop() {
		_, _ = setup.handler.Handle(ctx, cmd)
	}
}

func TestBatchCreateUsersHandler_LargeInput(t *testing.T) {
	ctx := context.Background()
	setup := setupBatchTest()

	// 测试 100 个用户（最大允许数量）
	users := make([]BatchUserItem, 100)
	for i := range 100 {
		users[i] = BatchUserItem{
			Username: "user",
			Email:    "user@example.com",
			Password: "Password123!",
			FullName: "Test User",
		}
	}

	cmd := BatchCreateUsersCommand{Users: users}
	result, err := setup.handler.Handle(ctx, cmd)

	require.NoError(t, err)
	assert.Equal(t, 100, result.Total)
	// 由于 mock 的 ExistsByUsername 不会检查重复，所有用户都应该成功
	// 但实际上由于 mock 的实现，所有用户都会成功
	assert.Equal(t, 100, result.Success)
}

func TestBatchCreateUsersResult_Structure(t *testing.T) {
	t.Run("BatchCreateUsersResult 正确初始化", func(t *testing.T) {
		result := &BatchCreateUsersResult{
			Total:   10,
			Success: 8,
			Failed:  2,
			Errors: []BatchCreateUserError{
				{Index: 0, Username: "user1", Email: "user1@example.com", Error: "用户名已存在"},
				{Index: 5, Username: "user6", Email: "user6@example.com", Error: "邮箱已存在"},
			},
		}

		assert.Equal(t, 10, result.Total)
		assert.Equal(t, 8, result.Success)
		assert.Equal(t, 2, result.Failed)
		assert.Len(t, result.Errors, 2)
		assert.Equal(t, 0, result.Errors[0].Index)
		assert.Equal(t, 5, result.Errors[1].Index)
	})
}

func TestBatchCreateUsersHandler_HashPasswordFailure(t *testing.T) {
	ctx := context.Background()
	setup := setupBatchTest()
	setup.authService.hashError = errors.New("hash failed")

	cmd := BatchCreateUsersCommand{
		Users: []BatchUserItem{
			{Username: "user1", Email: "user1@example.com", Password: "Password123!", FullName: "User One"},
		},
	}

	result, err := setup.handler.Handle(ctx, cmd)

	require.NoError(t, err)
	assert.Equal(t, 1, result.Failed)
	assert.Contains(t, result.Errors[0].Error, "密码加密失败")
}

func TestBatchCreateUsersHandler_CheckUsernameQueryError(t *testing.T) {
	ctx := context.Background()
	setup := setupBatchTest()
	setup.queryRepo.checkError = errors.New("database error")

	cmd := BatchCreateUsersCommand{
		Users: []BatchUserItem{
			{Username: "user1", Email: "user1@example.com", Password: "Password123!", FullName: "User One"},
		},
	}

	result, err := setup.handler.Handle(ctx, cmd)

	require.NoError(t, err)
	assert.Equal(t, 1, result.Failed)
	assert.Contains(t, result.Errors[0].Error, "检查用户名失败")
}

// 测试密码策略正常工作（不设置 weakPassword 标志）
func TestBatchCreateUsersHandler_PasswordPolicyNormal(t *testing.T) {
	ctx := context.Background()

	t.Run("正常密码通过验证", func(t *testing.T) {
		setup := setupBatchTest()
		// 不设置 weakPassword，正常密码应该通过

		cmd := BatchCreateUsersCommand{
			Users: []BatchUserItem{
				{Username: "user1", Email: "user1@example.com", Password: "Password123!", FullName: "User One"},
			},
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		assert.Equal(t, 1, result.Success)
		assert.Equal(t, 0, result.Failed)
	})

	t.Run("短密码被拒绝", func(t *testing.T) {
		setup := setupBatchTest()

		cmd := BatchCreateUsersCommand{
			Users: []BatchUserItem{
				{Username: "user1", Email: "user1@example.com", Password: "12345", FullName: "User One"}, // 太短
			},
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		assert.Equal(t, 0, result.Success)
		assert.Equal(t, 1, result.Failed)
		// mockAuthService 会检查长度 < 6，返回 ErrWeakPassword
		// 批量操作不直接返回错误，而是记录失败计数
	})
}
