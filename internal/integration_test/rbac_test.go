package integration_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	roleCommand "github.com/lwmacct/251117-go-ddd-template/internal/application/role/command"
	roleQuery "github.com/lwmacct/251117-go-ddd-template/internal/application/role/query"
	userCommand "github.com/lwmacct/251117-go-ddd-template/internal/application/user/command"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/role"
)

func TestRBACFlow_RoleLifecycle(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	ctx := context.Background()

	t.Run("完整角色生命周期: 创建 → 查询 → 更新 → 删除", func(t *testing.T) {
		// 步骤 1: 创建角色
		createResult, err := env.CreateRoleHandler.Handle(ctx, roleCommand.CreateRoleCommand{
			Name:        "test_role",
			DisplayName: "测试角色",
			Description: "用于测试的角色",
		})

		require.NoError(t, err)
		require.NotNil(t, createResult)
		assert.NotZero(t, createResult.RoleID)
		assert.Equal(t, "test_role", createResult.Name)
		assert.Equal(t, "测试角色", createResult.DisplayName)

		roleID := createResult.RoleID

		// 步骤 2: 查询角色
		foundRole, err := env.GetRoleHandler.Handle(ctx, roleQuery.GetRoleQuery{
			RoleID: roleID,
		})

		require.NoError(t, err)
		require.NotNil(t, foundRole)
		assert.Equal(t, "test_role", foundRole.Name)
		assert.Equal(t, "测试角色", foundRole.DisplayName)
		assert.Equal(t, "用于测试的角色", foundRole.Description)

		// 步骤 3: 更新角色
		newDisplayName := "更新后的测试角色"
		newDescription := "更新后的描述"
		updateResult, err := env.UpdateRoleHandler.Handle(ctx, roleCommand.UpdateRoleCommand{
			RoleID:      roleID,
			DisplayName: &newDisplayName,
			Description: &newDescription,
		})

		require.NoError(t, err)
		require.NotNil(t, updateResult)
		assert.Equal(t, "更新后的测试角色", updateResult.DisplayName)
		assert.Equal(t, "更新后的描述", updateResult.Description)

		// 步骤 4: 验证更新生效
		updatedRole, err := env.GetRoleHandler.Handle(ctx, roleQuery.GetRoleQuery{
			RoleID: roleID,
		})

		require.NoError(t, err)
		assert.Equal(t, "更新后的测试角色", updatedRole.DisplayName)

		// 步骤 5: 删除角色
		err = env.DeleteRoleHandler.Handle(ctx, roleCommand.DeleteRoleCommand{
			RoleID: roleID,
		})

		require.NoError(t, err)

		// 步骤 6: 验证角色已删除
		_, err = env.GetRoleHandler.Handle(ctx, roleQuery.GetRoleQuery{
			RoleID: roleID,
		})

		assert.Error(t, err)
	})
}

func TestRBACFlow_DuplicateRoleName(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	ctx := context.Background()

	t.Run("重复角色名应该失败", func(t *testing.T) {
		// 创建第一个角色
		_, err := env.CreateRoleHandler.Handle(ctx, roleCommand.CreateRoleCommand{
			Name:        "duplicate_role",
			DisplayName: "重复角色",
			Description: "第一个",
		})
		require.NoError(t, err)

		// 尝试创建同名角色
		_, err = env.CreateRoleHandler.Handle(ctx, roleCommand.CreateRoleCommand{
			Name:        "duplicate_role",
			DisplayName: "重复角色2",
			Description: "第二个",
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestRBACFlow_ListRoles(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	ctx := context.Background()

	t.Run("分页列表查询角色", func(t *testing.T) {
		// 创建多个角色
		for i := 1; i <= 5; i++ {
			_, err := env.CreateRoleHandler.Handle(ctx, roleCommand.CreateRoleCommand{
				Name:        "list_role_" + string(rune('0'+i)),
				DisplayName: "列表角色" + string(rune('0'+i)),
				Description: "用于列表测试",
			})
			require.NoError(t, err)
		}

		// 测试分页
		result, err := env.ListRolesHandler.Handle(ctx, roleQuery.ListRolesQuery{
			Page:  1,
			Limit: 3,
		})

		require.NoError(t, err)
		assert.GreaterOrEqual(t, result.Total, int64(5))
		assert.Len(t, result.Roles, 3)

		// 测试第二页
		result2, err := env.ListRolesHandler.Handle(ctx, roleQuery.ListRolesQuery{
			Page:  2,
			Limit: 3,
		})

		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(result2.Roles), 2)
	})
}

func TestRBACFlow_UserRoleAssignment(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	ctx := context.Background()

	t.Run("用户角色分配流程", func(t *testing.T) {
		// 步骤 1: 创建用户
		testUser := env.CreateTestUser(ctx, t, "rbac_user", "rbac@example.com", "Password123!")

		// 步骤 2: 创建角色
		adminResult, err := env.CreateRoleHandler.Handle(ctx, roleCommand.CreateRoleCommand{
			Name:        "admin_role",
			DisplayName: "管理员角色",
			Description: "拥有管理权限",
		})
		require.NoError(t, err)

		editorResult, err := env.CreateRoleHandler.Handle(ctx, roleCommand.CreateRoleCommand{
			Name:        "editor_role",
			DisplayName: "编辑者角色",
			Description: "拥有编辑权限",
		})
		require.NoError(t, err)

		// 步骤 3: 分配角色给用户
		err = env.AssignRolesHandler.Handle(ctx, userCommand.AssignRolesCommand{
			UserID:  testUser.ID,
			RoleIDs: []uint{adminResult.RoleID, editorResult.RoleID},
		})

		require.NoError(t, err)

		// 步骤 4: 验证用户角色
		userRoleIDs, err := env.UserQueryRepo.GetRoles(ctx, testUser.ID)
		require.NoError(t, err)
		assert.Len(t, userRoleIDs, 2)
		assert.Contains(t, userRoleIDs, adminResult.RoleID)
		assert.Contains(t, userRoleIDs, editorResult.RoleID)
	})

	t.Run("更新用户角色（替换）", func(t *testing.T) {
		// 创建用户
		testUser := env.CreateTestUser(ctx, t, "rbac_user2", "rbac2@example.com", "Password123!")

		// 创建角色
		role1, err := env.CreateRoleHandler.Handle(ctx, roleCommand.CreateRoleCommand{
			Name:        "role_a",
			DisplayName: "角色A",
		})
		require.NoError(t, err)

		role2, err := env.CreateRoleHandler.Handle(ctx, roleCommand.CreateRoleCommand{
			Name:        "role_b",
			DisplayName: "角色B",
		})
		require.NoError(t, err)

		role3, err := env.CreateRoleHandler.Handle(ctx, roleCommand.CreateRoleCommand{
			Name:        "role_c",
			DisplayName: "角色C",
		})
		require.NoError(t, err)

		// 先分配角色A和B
		err = env.AssignRolesHandler.Handle(ctx, userCommand.AssignRolesCommand{
			UserID:  testUser.ID,
			RoleIDs: []uint{role1.RoleID, role2.RoleID},
		})
		require.NoError(t, err)

		// 重新分配为角色B和C（替换）
		err = env.AssignRolesHandler.Handle(ctx, userCommand.AssignRolesCommand{
			UserID:  testUser.ID,
			RoleIDs: []uint{role2.RoleID, role3.RoleID},
		})
		require.NoError(t, err)

		// 验证角色已更新
		userRoleIDs, err := env.UserQueryRepo.GetRoles(ctx, testUser.ID)
		require.NoError(t, err)
		assert.Len(t, userRoleIDs, 2)
		assert.NotContains(t, userRoleIDs, role1.RoleID) // 角色A 应被移除
		assert.Contains(t, userRoleIDs, role2.RoleID)    // 角色B 应保留
		assert.Contains(t, userRoleIDs, role3.RoleID)    // 角色C 应新增
	})

	t.Run("清空用户角色", func(t *testing.T) {
		// 创建用户
		testUser := env.CreateTestUser(ctx, t, "rbac_user3", "rbac3@example.com", "Password123!")

		// 创建并分配角色
		role1, err := env.CreateRoleHandler.Handle(ctx, roleCommand.CreateRoleCommand{
			Name:        "temp_role",
			DisplayName: "临时角色",
		})
		require.NoError(t, err)

		err = env.AssignRolesHandler.Handle(ctx, userCommand.AssignRolesCommand{
			UserID:  testUser.ID,
			RoleIDs: []uint{role1.RoleID},
		})
		require.NoError(t, err)

		// 清空角色
		err = env.AssignRolesHandler.Handle(ctx, userCommand.AssignRolesCommand{
			UserID:  testUser.ID,
			RoleIDs: []uint{}, // 空数组
		})
		require.NoError(t, err)

		// 验证角色已清空
		userRoleIDs, err := env.UserQueryRepo.GetRoles(ctx, testUser.ID)
		require.NoError(t, err)
		assert.Empty(t, userRoleIDs)
	})
}

func TestRBACFlow_SystemRoleProtection(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	ctx := context.Background()

	t.Run("系统角色不能删除", func(t *testing.T) {
		// 直接通过 Repository 创建系统角色
		systemRole := &role.Role{
			Name:        "super_admin",
			DisplayName: "超级管理员",
			Description: "系统内置角色",
			IsSystem:    true,
		}

		err := env.RoleCommandRepo.Create(ctx, systemRole)
		require.NoError(t, err)

		// 尝试删除系统角色
		err = env.DeleteRoleHandler.Handle(ctx, roleCommand.DeleteRoleCommand{
			RoleID: systemRole.ID,
		})

		// 应该失败
		require.Error(t, err)
		assert.Contains(t, err.Error(), "system role")
	})
}
