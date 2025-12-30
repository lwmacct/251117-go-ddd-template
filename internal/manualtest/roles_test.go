package manualtest

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/251117-go-ddd-template/internal/application/role"
	"github.com/lwmacct/251117-go-ddd-template/internal/manualtest/helper"
)

// TestRolesFlow 角色管理完整流程测试。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestRolesFlow ./internal/manualtest/
func TestRolesFlow(t *testing.T) {
	c := helper.LoginAsAdmin(t)

	// 测试 1: 获取角色列表
	t.Log("\n测试 1: 获取角色列表")
	roles, meta, err := helper.GetList[role.RoleDTO](c, "/api/system/roles", map[string]string{
		"page":  "1",
		"limit": "10",
	})
	require.NoError(t, err, "获取角色列表失败")
	t.Logf("  角色数量: %d", len(roles))
	if meta != nil {
		t.Logf("  总数: %d", meta.Total)
	}
	for _, r := range roles {
		t.Logf("    - [%d] %s (%s)", r.ID, r.DisplayName, r.Name)
	}

	// 测试 2: 创建角色（使用工厂函数）
	t.Log("\n测试 2: 创建角色")
	testRole, markDeleted := helper.CreateTestRoleWithCleanupControl(t, c, "testrole")
	t.Logf("  创建成功! 角色 ID: %d", testRole.RoleID)

	// 验证角色 ID 有效
	require.NotZero(t, testRole.RoleID, "创建角色失败: 返回的角色 ID 为 0")

	// 测试 3: 获取角色详情
	t.Log("\n测试 3: 获取角色详情")
	roleDetail, err := helper.Get[role.RoleDTO](c, fmt.Sprintf("/api/system/roles/%d", testRole.RoleID), nil)
	require.NoError(t, err, "获取角色详情失败")
	t.Logf("  角色名: %s, 显示名: %s", roleDetail.Name, roleDetail.DisplayName)
	t.Logf("  描述: %s", roleDetail.Description)
	t.Logf("  权限数量: %d", len(roleDetail.Permissions))

	// 验证角色详情
	assert.Equal(t, testRole.RoleID, roleDetail.ID, "角色 ID 不匹配")

	// 测试 4: 更新角色
	t.Log("\n测试 4: 更新角色")
	newDisplayName := "测试角色（已更新）"
	newDescription := "更新后的描述"
	updateReq := role.UpdateDTO{
		DisplayName: &newDisplayName,
		Description: &newDescription,
	}
	updatedRole, err := helper.Put[role.RoleDTO](c, fmt.Sprintf("/api/system/roles/%d", testRole.RoleID), updateReq)
	require.NoError(t, err, "更新角色失败")
	t.Logf("  更新成功! 显示名: %s", updatedRole.DisplayName)

	// 验证更新后的字段
	assert.Equal(t, newDisplayName, updatedRole.DisplayName, "显示名未更新")
	assert.Equal(t, newDescription, updatedRole.Description, "描述未更新")

	// 测试 5: 设置角色权限（新 RBAC 模型）
	t.Log("\n测试 5: 设置角色权限")
	testSetRolePermissions(t, c, testRole.RoleID)

	// 测试 6: 删除角色
	t.Log("\n测试 6: 删除角色")
	err = c.Delete(fmt.Sprintf("/api/system/roles/%d", testRole.RoleID))
	require.NoError(t, err, "删除角色失败")
	t.Log("  删除成功!")

	// 标记已删除，避免 t.Cleanup 重复删除
	markDeleted()

	t.Log("\n角色管理流程测试完成!")
}

// TestListRoles 测试获取角色列表。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestListRoles ./internal/manualtest/
func TestListRoles(t *testing.T) {
	c := helper.LoginAsAdmin(t)

	t.Log("获取角色列表...")
	roles, meta, err := helper.GetList[role.RoleDTO](c, "/api/system/roles", map[string]string{
		"page":  "1",
		"limit": "10",
	})
	require.NoError(t, err, "获取角色列表失败")

	t.Logf("角色数量: %d", len(roles))
	if meta != nil {
		t.Logf("总数: %d, 总页数: %d", meta.Total, meta.TotalPages)
	}

	for _, r := range roles {
		systemFlag := ""
		if r.IsSystem {
			systemFlag = " [系统]"
		}
		t.Logf("  - [%d] %s (%s)%s - 权限数: %d", r.ID, r.DisplayName, r.Name, systemFlag, len(r.Permissions))
	}
}

// testSetRolePermissions 设置角色权限并验证（辅助函数）。
// 新 RBAC 模型：使用 Operation Pattern + Resource Pattern
func testSetRolePermissions(t *testing.T, c *helper.Client, roleID uint) {
	t.Helper()

	// 新 RBAC 模型：使用 Operation Pattern + Resource Pattern
	permissions := []role.PermissionInputDTO{
		{OperationPattern: "sys:users.read", ResourcePattern: "*"},
		{OperationPattern: "sys:users.list", ResourcePattern: "*"},
		{OperationPattern: "user:profile.*", ResourcePattern: "user/self"},
	}

	setPermReq := role.SetPermissionsDTO{
		Permissions: permissions,
	}
	t.Logf("  设置权限模式: %v", permissions)

	resp, err := c.R().
		SetBody(setPermReq).
		Put(fmt.Sprintf("/api/system/roles/%d/permissions", roleID))
	require.NoError(t, err, "设置权限请求失败")
	require.False(t, resp.IsError(), "设置权限失败，状态码: %d", resp.StatusCode())
	t.Log("  权限设置成功!")

	// 验证权限已设置
	roleWithPerms, err := helper.Get[role.RoleDTO](c, fmt.Sprintf("/api/system/roles/%d", roleID), nil)
	require.NoError(t, err, "获取角色详情失败")
	t.Logf("  验证：角色现有 %d 个权限", len(roleWithPerms.Permissions))

	// 验证权限数量
	assert.Len(t, roleWithPerms.Permissions, len(permissions), "权限数量不匹配")

	// 显示权限详情
	for _, p := range roleWithPerms.Permissions {
		t.Logf("    - %s | %s", p.OperationPattern, p.ResourcePattern)
	}
}

// TestSystemRoleProtection 测试系统角色保护机制。
//
// 系统角色（admin、user）不可删除。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestSystemRoleProtection ./internal/manualtest/
func TestSystemRoleProtection(t *testing.T) {
	c := helper.LoginAsAdmin(t)

	// 获取系统角色
	t.Log("\n步骤 1: 获取系统角色")
	roles, _, err := helper.GetList[role.RoleDTO](c, "/api/system/roles", nil)
	require.NoError(t, err, "获取角色列表失败")

	var adminRole, userRole *role.RoleDTO
	for i := range roles {
		if roles[i].Name == "admin" {
			adminRole = &roles[i]
		}
		if roles[i].Name == "user" {
			userRole = &roles[i]
		}
	}
	require.NotNil(t, adminRole, "未找到 admin 角色")
	require.NotNil(t, userRole, "未找到 user 角色")

	t.Logf("  admin 角色: ID=%d, IsSystem=%v", adminRole.ID, adminRole.IsSystem)
	t.Logf("  user 角色: ID=%d, IsSystem=%v", userRole.ID, userRole.IsSystem)

	// 验证系统角色标记
	assert.True(t, adminRole.IsSystem, "admin 应为系统角色")
	assert.True(t, userRole.IsSystem, "user 应为系统角色")

	// 测试 2: 尝试删除 admin 角色（应失败）
	t.Log("\n步骤 2: 尝试删除 admin 角色（应失败）")
	err = c.Delete(fmt.Sprintf("/api/system/roles/%d", adminRole.ID))
	require.Error(t, err, "删除 admin 角色应该失败")
	t.Logf("  预期失败: %v", err)

	// 测试 3: 尝试删除 user 角色（应失败）
	t.Log("\n步骤 3: 尝试删除 user 角色（应失败）")
	err = c.Delete(fmt.Sprintf("/api/system/roles/%d", userRole.ID))
	require.Error(t, err, "删除 user 角色应该失败")
	t.Logf("  预期失败: %v", err)

	t.Log("\n系统角色保护测试完成!")
}
