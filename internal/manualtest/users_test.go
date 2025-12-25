package manualtest

import (
	"fmt"
	"testing"
	"time"

	"github.com/lwmacct/251117-go-ddd-template/internal/application/user"
	"github.com/lwmacct/251117-go-ddd-template/internal/manualtest/helper"
)

// TestAdminUsersFlow 用户管理完整流程测试。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestAdminUsersFlow ./internal/manualtest/
func TestAdminUsersFlow(t *testing.T) {
	helper.SkipIfNotManual(t)

	c := helper.NewClient()

	t.Log("准备工作: 登录管理员账户")
	_, err := c.Login("admin", "admin123")
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}
	t.Log("  登录成功")

	var testUserID uint

	// 测试 1: 获取用户列表
	t.Log("\n测试 1: 获取用户列表")
	users, meta, err := helper.GetList[user.UserDTO](c, "/api/admin/users", map[string]string{
		"page":  "1",
		"limit": "10",
	})
	if err != nil {
		t.Fatalf("获取用户列表失败: %v", err)
	}
	t.Logf("  用户数量: %d", len(users))
	if meta != nil {
		t.Logf("  总数: %d", meta.Total)
	}

	// 测试 2: 创建用户
	t.Log("\n测试 2: 创建用户")
	testUsername := fmt.Sprintf("testuser_%d", time.Now().Unix())
	createReq := user.CreateUserDTO{
		Username: testUsername,
		Email:    testUsername + "@example.com",
		Password: "password123",
		FullName: "测试用户",
	}
	t.Logf("  创建用户: %s", createReq.Username)

	createResp, err := helper.Post[user.UserWithRolesDTO](c, "/api/admin/users", createReq)
	if err != nil {
		t.Fatalf("创建用户失败: %v", err)
	}
	testUserID = createResp.ID
	t.Logf("  创建成功! 用户 ID: %d", createResp.ID)

	// 测试 3: 获取用户详情
	t.Log("\n测试 3: 获取用户详情")
	userDetail, err := helper.Get[user.UserDTO](c, fmt.Sprintf("/api/admin/users/%d", testUserID), nil)
	if err != nil {
		t.Fatalf("获取用户详情失败: %v", err)
	}
	t.Logf("  用户名: %s, 邮箱: %s", userDetail.Username, userDetail.Email)

	// 测试 4: 更新用户
	t.Log("\n测试 4: 更新用户")
	newFullName := "测试用户（已更新）"
	updateReq := user.UpdateUserDTO{
		FullName: &newFullName,
	}
	updatedUser, err := helper.Put[user.UserDTO](c, fmt.Sprintf("/api/admin/users/%d", testUserID), updateReq)
	if err != nil {
		t.Fatalf("更新用户失败: %v", err)
	}
	t.Logf("  更新成功! 全名: %s", updatedUser.FullName)

	// 测试 5: 删除用户
	t.Log("\n测试 5: 删除用户")
	err = c.Delete(fmt.Sprintf("/api/admin/users/%d", testUserID))
	if err != nil {
		t.Fatalf("删除用户失败: %v", err)
	}
	t.Log("  删除成功!")

	t.Log("\n用户管理流程测试完成!")
}

// TestListUsers 测试获取用户列表。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestListUsers ./internal/manualtest/
func TestListUsers(t *testing.T) {
	helper.SkipIfNotManual(t)

	c := helper.NewClient()

	_, err := c.Login("admin", "admin123")
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}

	t.Log("获取用户列表...")
	users, meta, err := helper.GetList[user.UserDTO](c, "/api/admin/users", map[string]string{
		"page":  "1",
		"limit": "10",
	})
	if err != nil {
		t.Fatalf("获取用户列表失败: %v", err)
	}

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
	helper.SkipIfNotManual(t)

	c := helper.NewClient()

	t.Log("准备工作: 登录管理员账户")
	_, err := c.Login("admin", "admin123")
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}
	t.Log("  登录成功")

	// 创建测试用户
	t.Log("\n步骤 1: 创建测试用户")
	testUsername := fmt.Sprintf("roletest_%d", time.Now().Unix())
	createReq := user.CreateUserDTO{
		Username: testUsername,
		Email:    testUsername + "@example.com",
		Password: "password123",
		FullName: "角色测试用户",
	}

	createResp, err := helper.Post[user.UserWithRolesDTO](c, "/api/admin/users", createReq)
	if err != nil {
		t.Fatalf("创建用户失败: %v", err)
	}
	testUserID := createResp.ID
	t.Logf("  创建成功! 用户 ID: %d", testUserID)

	// 分配角色（使用 user 角色 ID=2）
	t.Log("\n步骤 2: 分配角色")
	assignReq := user.AssignRolesDTO{
		RoleIDs: []uint{2}, // user 角色
	}
	t.Logf("  分配角色 IDs: %v", assignReq.RoleIDs)

	assignResp, err := helper.Put[user.UserWithRolesDTO](c, fmt.Sprintf("/api/admin/users/%d/roles", testUserID), assignReq)
	if err != nil {
		t.Fatalf("分配角色失败: %v", err)
	}

	t.Logf("  分配成功! 用户现有角色数: %d", len(assignResp.Roles))
	for _, r := range assignResp.Roles {
		t.Logf("    - [%d] %s (%s)", r.ID, r.DisplayName, r.Name)
	}

	// 验证角色已分配
	if len(assignResp.Roles) == 0 {
		t.Fatal("角色分配失败，用户没有角色")
	}

	// 清理
	t.Log("\n步骤 3: 清理测试用户")
	err = c.Delete(fmt.Sprintf("/api/admin/users/%d", testUserID))
	if err != nil {
		t.Logf("警告：无法删除测试用户: %v", err)
	} else {
		t.Log("  测试用户已删除")
	}

	t.Log("\n角色分配测试完成!")
}
