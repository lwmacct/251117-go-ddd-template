package manualtest

import (
	"fmt"
	"testing"
	"time"

	"github.com/lwmacct/251117-go-ddd-template/internal/application/role"
	"github.com/lwmacct/251117-go-ddd-template/internal/manualtest/helper"
)

// TestRolesFlow 角色管理完整流程测试。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestRolesFlow ./internal/manualtest/
func TestRolesFlow(t *testing.T) {
	helper.SkipIfNotManual(t)

	c := helper.NewClient()

	t.Log("准备工作: 登录管理员账户")
	_, err := c.Login("admin", "admin123")
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}
	t.Log("  登录成功")

	var testRoleID uint

	// 测试 1: 获取角色列表
	t.Log("\n测试 1: 获取角色列表")
	roles, meta, err := helper.GetList[role.RoleDTO](c, "/api/admin/roles", map[string]string{
		"page":  "1",
		"limit": "10",
	})
	if err != nil {
		t.Fatalf("获取角色列表失败: %v", err)
	}
	t.Logf("  角色数量: %d", len(roles))
	if meta != nil {
		t.Logf("  总数: %d", meta.Total)
	}
	for _, r := range roles {
		t.Logf("    - [%d] %s (%s)", r.ID, r.DisplayName, r.Name)
	}

	// 测试 2: 创建角色
	t.Log("\n测试 2: 创建角色")
	testRoleName := fmt.Sprintf("testrole_%d", time.Now().Unix())
	createReq := role.CreateRoleDTO{
		Name:        testRoleName,
		DisplayName: "测试角色",
		Description: "这是一个测试角色",
	}
	t.Logf("  创建角色: %s", createReq.Name)

	createResp, err := helper.Post[role.CreateRoleResultDTO](c, "/api/admin/roles", createReq)
	if err != nil {
		t.Fatalf("创建角色失败: %v", err)
	}
	testRoleID = createResp.RoleID
	t.Logf("  创建成功! 角色 ID: %d", testRoleID)

	// 测试 3: 获取角色详情
	t.Log("\n测试 3: 获取角色详情")
	roleDetail, err := helper.Get[role.RoleDTO](c, fmt.Sprintf("/api/admin/roles/%d", testRoleID), nil)
	if err != nil {
		t.Fatalf("获取角色详情失败: %v", err)
	}
	t.Logf("  角色名: %s, 显示名: %s", roleDetail.Name, roleDetail.DisplayName)
	t.Logf("  描述: %s", roleDetail.Description)
	t.Logf("  权限数量: %d", len(roleDetail.Permissions))

	// 测试 4: 更新角色
	t.Log("\n测试 4: 更新角色")
	newDisplayName := "测试角色（已更新）"
	newDescription := "更新后的描述"
	updateReq := role.UpdateRoleDTO{
		DisplayName: &newDisplayName,
		Description: &newDescription,
	}
	updatedRole, err := helper.Put[role.RoleDTO](c, fmt.Sprintf("/api/admin/roles/%d", testRoleID), updateReq)
	if err != nil {
		t.Fatalf("更新角色失败: %v", err)
	}
	t.Logf("  更新成功! 显示名: %s", updatedRole.DisplayName)

	// 测试 5: 获取权限列表
	t.Log("\n测试 5: 获取权限列表")
	permissions, permMeta, err := helper.GetList[role.PermissionDTO](c, "/api/admin/permissions", map[string]string{
		"page":  "1",
		"limit": "50",
	})
	if err != nil {
		t.Fatalf("获取权限列表失败: %v", err)
	}
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
		testSetRolePermissions(t, c, testRoleID, permissions[:3])
	}

	// 测试 7: 删除角色
	t.Log("\n测试 7: 删除角色")
	err = c.Delete(fmt.Sprintf("/api/admin/roles/%d", testRoleID))
	if err != nil {
		t.Fatalf("删除角色失败: %v", err)
	}
	t.Log("  删除成功!")

	t.Log("\n角色管理流程测试完成!")
}

// TestListRoles 测试获取角色列表。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestListRoles ./internal/manualtest/
func TestListRoles(t *testing.T) {
	helper.SkipIfNotManual(t)

	c := helper.NewClient()

	_, err := c.Login("admin", "admin123")
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}

	t.Log("获取角色列表...")
	roles, meta, err := helper.GetList[role.RoleDTO](c, "/api/admin/roles", map[string]string{
		"page":  "1",
		"limit": "10",
	})
	if err != nil {
		t.Fatalf("获取角色列表失败: %v", err)
	}

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
	helper.SkipIfNotManual(t)

	c := helper.NewClient()

	_, err := c.Login("admin", "admin123")
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}

	t.Log("获取权限列表...")
	permissions, meta, err := helper.GetList[role.PermissionDTO](c, "/api/admin/permissions", map[string]string{
		"page":  "1",
		"limit": "50",
	})
	if err != nil {
		t.Fatalf("获取权限列表失败: %v", err)
	}

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

	permIDs := make([]uint, len(permissions))
	for i, p := range permissions {
		permIDs[i] = p.ID
	}

	setPermReq := role.SetPermissionsDTO{
		PermissionIDs: permIDs,
	}
	t.Logf("  设置权限 IDs: %v", permIDs)

	resp, err := c.R().
		SetBody(setPermReq).
		Put(fmt.Sprintf("/api/admin/roles/%d/permissions", roleID))
	if err != nil {
		t.Fatalf("设置权限请求失败: %v", err)
	}
	if resp.IsError() {
		t.Fatalf("设置权限失败，状态码: %d", resp.StatusCode())
	}
	t.Log("  权限设置成功!")

	// 验证权限已设置
	roleWithPerms, err := helper.Get[role.RoleDTO](c, fmt.Sprintf("/api/admin/roles/%d", roleID), nil)
	if err != nil {
		t.Fatalf("获取角色详情失败: %v", err)
	}
	t.Logf("  验证：角色现有 %d 个权限", len(roleWithPerms.Permissions))
}
