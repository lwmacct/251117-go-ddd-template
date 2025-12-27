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
	roles, meta, err := helper.GetList[role.RoleDTO](c, "/api/admin/roles", map[string]string{
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
	roleDetail, err := helper.Get[role.RoleDTO](c, fmt.Sprintf("/api/admin/roles/%d", testRole.RoleID), nil)
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
	updatedRole, err := helper.Put[role.RoleDTO](c, fmt.Sprintf("/api/admin/roles/%d", testRole.RoleID), updateReq)
	require.NoError(t, err, "更新角色失败")
	t.Logf("  更新成功! 显示名: %s", updatedRole.DisplayName)

	// 验证更新后的字段
	assert.Equal(t, newDisplayName, updatedRole.DisplayName, "显示名未更新")
	assert.Equal(t, newDescription, updatedRole.Description, "描述未更新")

	// 测试 5: 获取权限列表
	t.Log("\n测试 5: 获取权限列表")
	permissions, permMeta, err := helper.GetList[role.PermissionDTO](c, "/api/admin/permissions", map[string]string{
		"page":  "1",
		"limit": "50",
	})
	require.NoError(t, err, "获取权限列表失败")
	t.Logf("  权限数量: %d", len(permissions))
	if permMeta != nil {
		t.Logf("  总数: %d", permMeta.Total)
	}

	// 显示前 5 个权限
	for i, p := range permissions {
		if i >= 5 {
			t.Logf("    ... 还有 %d 个权限", len(permissions)-5)
			break
		}
		t.Logf("    - [%d] %s: %s", p.ID, p.Code, p.Description)
	}

	// 测试 6: 设置角色权限
	t.Log("\n测试 6: 设置角色权限")
	if len(permissions) < 3 {
		t.Log("  跳过：权限数量不足")
	} else {
		testSetRolePermissions(t, c, testRole.RoleID, permissions[:3])
	}

	// 测试 7: 删除角色
	t.Log("\n测试 7: 删除角色")
	err = c.Delete(fmt.Sprintf("/api/admin/roles/%d", testRole.RoleID))
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
	roles, meta, err := helper.GetList[role.RoleDTO](c, "/api/admin/roles", map[string]string{
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
		t.Logf("  - [%d] %s (%s)%s", r.ID, r.DisplayName, r.Name, systemFlag)
	}
}

// TestListPermissions 测试获取权限列表。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestListPermissions ./internal/manualtest/
func TestListPermissions(t *testing.T) {
	c := helper.LoginAsAdmin(t)

	t.Log("获取权限列表...")
	permissions, meta, err := helper.GetList[role.PermissionDTO](c, "/api/admin/permissions", map[string]string{
		"page":  "1",
		"limit": "50",
	})
	require.NoError(t, err, "获取权限列表失败")

	t.Logf("权限数量: %d", len(permissions))
	if meta != nil {
		t.Logf("总数: %d, 总页数: %d", meta.Total, meta.TotalPages)
	}

	// 按 domain 分组显示
	domains := make(map[string][]role.PermissionDTO)
	for _, p := range permissions {
		domains[p.Resource] = append(domains[p.Resource], p)
	}

	for domain, perms := range domains {
		t.Logf("\n  [%s] %d 个权限:", domain, len(perms))
		for _, p := range perms {
			t.Logf("    - %s: %s", p.Code, p.Description)
		}
	}
}

// testSetRolePermissions 设置角色权限并验证（辅助函数，降低嵌套复杂度）。
func testSetRolePermissions(t *testing.T, c *helper.Client, roleID uint, permissions []role.PermissionDTO) {
	t.Helper()

	permIDs := helper.ExtractIDs(permissions, func(p role.PermissionDTO) uint { return p.ID })

	setPermReq := role.SetPermissionsDTO{
		PermissionIDs: permIDs,
	}
	t.Logf("  设置权限 IDs: %v", permIDs)

	resp, err := c.R().
		SetBody(setPermReq).
		Put(fmt.Sprintf("/api/admin/roles/%d/permissions", roleID))
	require.NoError(t, err, "设置权限请求失败")
	require.False(t, resp.IsError(), "设置权限失败，状态码: %d", resp.StatusCode())
	t.Log("  权限设置成功!")

	// 验证权限已设置
	roleWithPerms, err := helper.Get[role.RoleDTO](c, fmt.Sprintf("/api/admin/roles/%d", roleID), nil)
	require.NoError(t, err, "获取角色详情失败")
	t.Logf("  验证：角色现有 %d 个权限", len(roleWithPerms.Permissions))

	// 验证权限 ID 是否匹配
	assert.Len(t, roleWithPerms.Permissions, len(permIDs), "权限数量不匹配")

	// 使用 assert.ElementsMatch 验证权限 ID 集合
	actualIDs := helper.ExtractIDs(roleWithPerms.Permissions, func(p *role.PermissionDTO) uint { return p.ID })
	assert.ElementsMatch(t, permIDs, actualIDs, "权限 ID 集合不匹配")
}
