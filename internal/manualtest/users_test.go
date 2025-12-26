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

	// 验证创建的用户数据
	if createResp.Username != createReq.Username {
		t.Errorf("用户名不匹配: 期望 %s, 实际 %s", createReq.Username, createResp.Username)
	}
	if createResp.Email != createReq.Email {
		t.Errorf("邮箱不匹配: 期望 %s, 实际 %s", createReq.Email, createResp.Email)
	}
	if createResp.FullName != createReq.FullName {
		t.Errorf("全名不匹配: 期望 %s, 实际 %s", createReq.FullName, createResp.FullName)
	}

	// 测试 3: 获取用户详情
	t.Log("\n测试 3: 获取用户详情")
	userDetail, err := helper.Get[user.UserDTO](c, fmt.Sprintf("/api/admin/users/%d", testUserID), nil)
	if err != nil {
		t.Fatalf("获取用户详情失败: %v", err)
	}
	t.Logf("  用户名: %s, 邮箱: %s", userDetail.Username, userDetail.Email)

	// 验证用户详情
	if userDetail.ID != testUserID {
		t.Errorf("用户 ID 不匹配: 期望 %d, 实际 %d", testUserID, userDetail.ID)
	}
	if userDetail.Username != testUsername {
		t.Errorf("用户名不匹配: 期望 %s, 实际 %s", testUsername, userDetail.Username)
	}

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

	// 验证更新后的字段
	if updatedUser.FullName != newFullName {
		t.Errorf("全名未更新: 期望 %s, 实际 %s", newFullName, updatedUser.FullName)
	}

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

	// 验证是否包含指定的角色 ID
	foundRole := false
	for _, r := range assignResp.Roles {
		if r.ID == 2 {
			foundRole = true
			break
		}
	}
	if !foundRole {
		t.Error("未找到预期的角色 ID=2")
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

// TestBatchCreateUsers 测试批量创建用户。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestBatchCreateUsers ./internal/manualtest/
func TestBatchCreateUsers(t *testing.T) {
	helper.SkipIfNotManual(t)

	c := helper.NewClient()

	t.Log("准备工作: 登录管理员账户")
	_, err := c.Login("admin", "admin123")
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}
	t.Log("  登录成功")

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

	batchReq := user.BatchCreateUserDTO{
		Users: []user.BatchUserItemDTO{
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

	result, err := helper.Post[user.BatchCreateUserResultDTO](c, "/api/admin/users/batch", batchReq)
	if err != nil {
		t.Fatalf("批量创建请求失败: %v", err)
	}

	t.Logf("\n批量创建结果:")
	t.Logf("  总数: %d", result.Total)
	t.Logf("  成功: %d", result.Success)
	t.Logf("  失败: %d", result.Failed)

	// 步骤 2: 验证结果
	t.Log("\n步骤 2: 验证结果")
	if result.Total != 3 {
		t.Errorf("总数应为 3，实际 %d", result.Total)
	}
	if result.Success != 2 {
		t.Errorf("成功数应为 2，实际 %d", result.Success)
	}
	if result.Failed != 1 {
		t.Errorf("失败数应为 1，实际 %d", result.Failed)
	}

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
	found1, found2 := false, false
	for _, u := range users {
		if u.Username == username1 {
			found1 = true
			t.Logf("  找到用户: %s (ID: %d)", u.Username, u.ID)
		}
		if u.Username == username2 {
			found2 = true
			t.Logf("  找到用户: %s (ID: %d)", u.Username, u.ID)
		}
	}
	if !found1 || !found2 {
		t.Error("部分用户未创建成功")
	}

	t.Log("\n批量创建用户测试完成!")
}
