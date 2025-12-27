package manualtest

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/251117-go-ddd-template/internal/application/user"
	"github.com/lwmacct/251117-go-ddd-template/internal/manualtest/helper"
)

// TestAdminUsersFlow 用户管理完整流程测试。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestAdminUsersFlow ./internal/manualtest/
func TestAdminUsersFlow(t *testing.T) {
	c := helper.LoginAsAdmin(t)

	// 测试 1: 获取用户列表
	t.Log("\n测试 1: 获取用户列表")
	users, meta, err := helper.GetList[user.UserDTO](c, "/api/admin/users", map[string]string{
		"page":  "1",
		"limit": "10",
	})
	require.NoError(t, err, "获取用户列表失败")
	t.Logf("  用户数量: %d", len(users))
	if meta != nil {
		t.Logf("  总数: %d", meta.Total)
	}

	// 测试 2: 创建用户（使用工厂函数，返回清理控制）
	t.Log("\n测试 2: 创建用户")
	testUser, markDeleted := helper.CreateTestUserWithCleanupControl(t, c, "testuser")
	t.Logf("  创建成功! 用户 ID: %d", testUser.ID)

	// 验证创建的用户数据
	assert.NotEmpty(t, testUser.Username, "用户名不应为空")
	assert.NotEmpty(t, testUser.Email, "邮箱不应为空")

	// 测试 3: 获取用户详情
	t.Log("\n测试 3: 获取用户详情")
	userDetail, err := helper.Get[user.UserDTO](c, fmt.Sprintf("/api/admin/users/%d", testUser.ID), nil)
	require.NoError(t, err, "获取用户详情失败")
	t.Logf("  用户名: %s, 邮箱: %s", userDetail.Username, userDetail.Email)

	// 验证用户详情
	assert.Equal(t, testUser.ID, userDetail.ID, "用户 ID 不匹配")

	// 测试 4: 更新用户
	t.Log("\n测试 4: 更新用户")
	newFullName := "测试用户（已更新）"
	updateReq := user.UpdateDTO{
		FullName: &newFullName,
	}
	updatedUser, err := helper.Put[user.UserDTO](c, fmt.Sprintf("/api/admin/users/%d", testUser.ID), updateReq)
	require.NoError(t, err, "更新用户失败")
	t.Logf("  更新成功! 全名: %s", updatedUser.FullName)

	// 验证更新后的字段
	assert.Equal(t, newFullName, updatedUser.FullName, "全名未更新")

	// 测试 5: 删除用户
	t.Log("\n测试 5: 删除用户")
	err = c.Delete(fmt.Sprintf("/api/admin/users/%d", testUser.ID))
	require.NoError(t, err, "删除用户失败")
	t.Log("  删除成功!")

	// 标记已删除，避免 t.Cleanup 重复删除
	markDeleted()

	t.Log("\n用户管理流程测试完成!")
}

// TestListUsers 测试获取用户列表。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestListUsers ./internal/manualtest/
func TestListUsers(t *testing.T) {
	c := helper.LoginAsAdmin(t)

	t.Log("获取用户列表...")
	users, meta, err := helper.GetList[user.UserDTO](c, "/api/admin/users", map[string]string{
		"page":  "1",
		"limit": "10",
	})
	require.NoError(t, err, "获取用户列表失败")

	t.Logf("用户数量: %d", len(users))
	if meta != nil {
		t.Logf("总数: %d, 总页数: %d", meta.Total, meta.TotalPages)
	}

	for _, u := range users {
		t.Logf("  - [%d] %s <%s> 状态: %s", u.ID, u.Username, u.Email, u.Status)
	}
}

// TestAssignRoles 测试分配用户角色。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestAssignRoles ./internal/manualtest/
func TestAssignRoles(t *testing.T) {
	c := helper.LoginAsAdmin(t)

	// 创建测试用户
	t.Log("\n步骤 1: 创建测试用户")
	testUser := helper.CreateTestUser(t, c, "roletest")
	t.Logf("  创建成功! 用户 ID: %d", testUser.ID)

	// 分配角色（使用 user 角色 ID=2）
	t.Log("\n步骤 2: 分配角色")
	assignReq := user.AssignRolesDTO{
		RoleIDs: []uint{2}, // user 角色
	}
	t.Logf("  分配角色 IDs: %v", assignReq.RoleIDs)

	assignResp, err := helper.Put[user.UserWithRolesDTO](c, fmt.Sprintf("/api/admin/users/%d/roles", testUser.ID), assignReq)
	require.NoError(t, err, "分配角色失败")

	t.Logf("  分配成功! 用户现有角色数: %d", len(assignResp.Roles))
	for _, r := range assignResp.Roles {
		t.Logf("    - [%d] %s (%s)", r.ID, r.DisplayName, r.Name)
	}

	// 验证角色已分配
	require.NotEmpty(t, assignResp.Roles, "角色分配失败，用户没有角色")

	// 使用 assert.Contains 验证是否包含指定的角色 ID
	roleIDs := helper.ExtractIDs(assignResp.Roles, func(r user.RoleDTO) uint { return r.ID })
	assert.Contains(t, roleIDs, uint(2), "未找到预期的角色 ID=2")

	t.Log("\n角色分配测试完成!")
}

// TestBatchCreateUsers 测试批量创建用户。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestBatchCreateUsers ./internal/manualtest/
func TestBatchCreateUsers(t *testing.T) {
	c := helper.LoginAsAdmin(t)

	timestamp := time.Now().Unix()
	username1 := fmt.Sprintf("batch1_%d", timestamp)
	username2 := fmt.Sprintf("batch2_%d", timestamp)

	// 确保测试结束时清理资源
	t.Cleanup(func() {
		// 获取用户列表并删除测试用户
		users, _, _ := helper.GetList[user.UserWithRolesDTO](c, "/api/admin/users", nil)
		for _, u := range users {
			if u.Username == username1 || u.Username == username2 {
				_ = c.Delete(fmt.Sprintf("/api/admin/users/%d", u.ID))
			}
		}
	})

	// 步骤 1: 批量创建用户（2个成功 + 1个重复失败）
	t.Log("\n步骤 1: 批量创建用户")
	t.Logf("  用户1: %s", username1)
	t.Logf("  用户2: %s", username2)
	t.Logf("  用户3: %s (重复，应失败)", username1)

	batchReq := user.BatchCreateDTO{
		Users: []user.BatchItemDTO{
			{
				Username: username1,
				Email:    username1 + "@example.com",
				Password: "test123456",
				FullName: "批量用户1",
			},
			{
				Username: username2,
				Email:    username2 + "@example.com",
				Password: "test123456",
				FullName: "批量用户2",
			},
			{
				Username: username1, // 重复用户名
				Email:    "dup_" + username1 + "@example.com",
				Password: "test123456",
				FullName: "重复用户",
			},
		},
	}

	result, err := helper.Post[user.BatchCreateResultDTO](c, "/api/admin/users/batch", batchReq)
	require.NoError(t, err, "批量创建请求失败")

	t.Logf("\n批量创建结果:")
	t.Logf("  总数: %d", result.Total)
	t.Logf("  成功: %d", result.Success)
	t.Logf("  失败: %d", result.Failed)

	// 步骤 2: 验证结果
	t.Log("\n步骤 2: 验证结果")
	assert.Equal(t, 3, result.Total, "总数应为 3")
	assert.Equal(t, 2, result.Success, "成功数应为 2")
	assert.Equal(t, 1, result.Failed, "失败数应为 1")

	// 验证错误详情
	if len(result.Errors) > 0 {
		t.Log("  错误详情:")
		for _, e := range result.Errors {
			t.Logf("    - [%d] %s: %s", e.Index, e.Username, e.Error)
		}
	}

	// 步骤 3: 验证用户已创建
	t.Log("\n步骤 3: 验证用户已创建")
	users, _, _ := helper.GetList[user.UserWithRolesDTO](c, "/api/admin/users", nil)

	// 使用 assert.Contains 验证用户名
	usernames := helper.ExtractStrings(users, func(u user.UserWithRolesDTO) string { return u.Username })
	assert.Contains(t, usernames, username1, "用户1未创建")
	assert.Contains(t, usernames, username2, "用户2未创建")

	t.Log("\n批量创建用户测试完成!")
}
