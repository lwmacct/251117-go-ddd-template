package manualtest

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/lwmacct/251117-go-ddd-template/internal/application/user"
	"github.com/lwmacct/251117-go-ddd-template/internal/manualtest/helper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 测试用户名前缀，用于隔离测试数据
const userTestPrefix = "test_user_"

// TestUsersAPIFlow 通用用户接口完整流程测试。
//
// 测试 CRUD 完整流程：创建 → 获取 → 更新 → 删除
// 注意：/api/users/* 与 /api/admin/users/* 的区别是权限要求不同
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestUsersAPIFlow ./internal/manualtest/
func TestUsersAPIFlow(t *testing.T) {
	helper.SkipIfNotManual(t)

	c := helper.NewClient()

	t.Log("准备工作: 登录账户")
	_, err := c.Login("admin", "admin123")
	require.NoError(t, err, "登录失败")
	t.Log("  登录成功")

	// 测试 1: 创建用户
	t.Log("\n测试 1: 创建用户")
	timestamp := time.Now().Unix()
	username := fmt.Sprintf("%s%d", userTestPrefix, timestamp)
	email := fmt.Sprintf("%s%d@test.com", userTestPrefix, timestamp)

	createReq := user.CreateDTO{
		Username: username,
		Email:    email,
		Password: "test123456",
		FullName: "测试用户",
	}
	t.Logf("  创建用户: %s (%s)", username, email)

	created, err := helper.Post[user.UserWithRolesDTO](c, "/api/users", createReq)
	require.NoError(t, err, "创建用户失败")
	require.NotZero(t, created.ID, "创建的用户 ID 为 0")
	assert.Equal(t, username, created.Username, "Username 不匹配")
	t.Logf("  创建成功! ID: %d", created.ID)
	t.Logf("  Username: %s", created.Username)
	t.Logf("  Email: %s", created.Email)
	t.Logf("  Status: %s", created.Status)

	userID := created.ID

	// 确保清理
	t.Cleanup(func() {
		if deleteErr := c.Delete(fmt.Sprintf("/api/users/%d", userID)); deleteErr != nil {
			t.Logf("清理用户失败: %v", deleteErr)
		}
	})

	// 测试 2: 获取用户详情
	t.Log("\n测试 2: 获取用户详情")
	detail, err := helper.Get[user.UserWithRolesDTO](c, fmt.Sprintf("/api/users/%d", userID), nil)
	require.NoError(t, err, "获取用户详情失败")
	assert.Equal(t, userID, detail.ID, "ID 不匹配")
	assert.Equal(t, username, detail.Username, "Username 不匹配")
	t.Logf("  ID: %d", detail.ID)
	t.Logf("  Username: %s", detail.Username)
	t.Logf("  Email: %s", detail.Email)
	t.Logf("  FullName: %s", detail.FullName)
	t.Logf("  Status: %s", detail.Status)
	t.Logf("  Roles: %d 个", len(detail.Roles))

	// 测试 3: 更新用户
	t.Log("\n测试 3: 更新用户")
	newFullName := "更新后的名称"
	newBio := "这是更新后的简介"
	updateReq := user.UpdateDTO{
		FullName: &newFullName,
		Bio:      &newBio,
	}

	resp, err := c.R().
		SetBody(updateReq).
		Put(fmt.Sprintf("/api/users/%d", userID))
	require.NoError(t, err, "更新用户失败")
	require.False(t, resp.IsError(), "更新用户失败: 状态码 %d", resp.StatusCode())
	t.Log("  更新成功!")

	// 验证更新结果
	updated, err := helper.Get[user.UserWithRolesDTO](c, fmt.Sprintf("/api/users/%d", userID), nil)
	require.NoError(t, err, "获取更新后用户失败")
	assert.Equal(t, newFullName, updated.FullName, "FullName 更新失败")
	assert.Equal(t, newBio, updated.Bio, "Bio 更新失败")
	t.Logf("  新 FullName: %s", updated.FullName)
	t.Logf("  新 Bio: %s", updated.Bio)

	// 测试 4: 获取用户列表
	t.Log("\n测试 4: 获取用户列表")
	users, meta, err := helper.GetList[user.UserWithRolesDTO](c, "/api/users", nil)
	require.NoError(t, err, "获取用户列表失败")
	t.Logf("  用户数: %d", len(users))
	if meta != nil {
		t.Logf("  总数: %d", meta.Total)
	}

	// 验证创建的用户在列表中
	found := false
	for _, u := range users {
		if u.ID == userID {
			found = true
			break
		}
	}
	assert.True(t, found, "创建的用户不在列表中")
	if found {
		t.Log("  ✓ 创建的用户存在于列表中")
	}

	t.Log("\n通用用户接口流程测试完成!")
}

// TestUsersAPICreate 测试创建用户。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestUsersAPICreate ./internal/manualtest/
func TestUsersAPICreate(t *testing.T) {
	helper.SkipIfNotManual(t)

	c := helper.NewClient()

	t.Log("登录账户...")
	_, err := c.Login("admin", "admin123")
	require.NoError(t, err, "登录失败")

	timestamp := time.Now().Unix()
	username := fmt.Sprintf("%screate_%d", userTestPrefix, timestamp)
	email := fmt.Sprintf("%screate_%d@test.com", userTestPrefix, timestamp)

	t.Logf("\n创建用户: %s", username)
	createReq := user.CreateDTO{
		Username: username,
		Email:    email,
		Password: "test123456",
		FullName: "创建测试用户",
	}

	created, err := helper.Post[user.UserWithRolesDTO](c, "/api/users", createReq)
	require.NoError(t, err, "创建用户失败")

	t.Logf("创建成功!")
	t.Logf("  ID: %d", created.ID)
	t.Logf("  Username: %s", created.Username)
	t.Logf("  Email: %s", created.Email)
	t.Logf("  Status: %s", created.Status)

	// 清理
	t.Cleanup(func() {
		if err := c.Delete(fmt.Sprintf("/api/users/%d", created.ID)); err != nil {
			t.Logf("清理用户失败: %v", err)
		}
	})
}

// TestUsersAPIList 测试获取用户列表。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestUsersAPIList ./internal/manualtest/
func TestUsersAPIList(t *testing.T) {
	helper.SkipIfNotManual(t)

	c := helper.NewClient()

	t.Log("登录账户...")
	_, err := c.Login("admin", "admin123")
	require.NoError(t, err, "登录失败")

	t.Log("\n获取用户列表...")
	users, meta, err := helper.GetList[user.UserWithRolesDTO](c, "/api/users", nil)
	require.NoError(t, err, "获取用户列表失败")

	t.Logf("用户数: %d", len(users))
	if meta != nil {
		t.Logf("总数: %d, 页码: %d, 每页: %d", meta.Total, meta.Page, meta.PerPage)
	}

	for _, u := range users {
		roleNames := make([]string, len(u.Roles))
		for i, r := range u.Roles {
			roleNames[i] = r.Name
		}
		rolesStr := "无"
		if len(roleNames) > 0 {
			rolesStr = strings.Join(roleNames, ", ")
		}
		t.Logf("  [%d] %s <%s> (状态: %s, 角色: %s)", u.ID, u.Username, u.Email, u.Status, rolesStr)
	}
}

// TestUsersAPIListWithPagination 测试用户列表分页。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestUsersAPIListWithPagination ./internal/manualtest/
func TestUsersAPIListWithPagination(t *testing.T) {
	helper.SkipIfNotManual(t)

	c := helper.NewClient()

	t.Log("登录账户...")
	_, err := c.Login("admin", "admin123")
	require.NoError(t, err, "登录失败")

	t.Log("\n测试分页: page=1, limit=2...")
	users, meta, err := helper.GetList[user.UserWithRolesDTO](c, "/api/users", map[string]string{
		"page":  "1",
		"limit": "2",
	})
	require.NoError(t, err, "获取用户列表失败")

	t.Logf("返回用户数: %d", len(users))
	if meta != nil {
		t.Logf("分页信息: 总数=%d, 页码=%d, 每页=%d", meta.Total, meta.Page, meta.PerPage)
	}

	assert.LessOrEqual(t, len(users), 2, "分页限制失败")

	for _, u := range users {
		t.Logf("  [%d] %s", u.ID, u.Username)
	}
}

// TestUsersAPIGetByID 测试获取用户详情。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestUsersAPIGetByID ./internal/manualtest/
func TestUsersAPIGetByID(t *testing.T) {
	helper.SkipIfNotManual(t)

	c := helper.NewClient()

	t.Log("登录账户...")
	_, err := c.Login("admin", "admin123")
	require.NoError(t, err, "登录失败")

	// 先获取列表，取第一个用户测试
	t.Log("\n获取用户列表以获取测试 ID...")
	users, _, err := helper.GetList[user.UserWithRolesDTO](c, "/api/users", map[string]string{
		"limit": "1",
	})
	require.NoError(t, err, "获取用户列表失败")
	if len(users) == 0 {
		t.Skip("没有用户可供测试")
	}

	testUserID := users[0].ID
	t.Logf("\n获取用户详情: ID=%d", testUserID)

	detail, err := helper.Get[user.UserWithRolesDTO](c, fmt.Sprintf("/api/users/%d", testUserID), nil)
	require.NoError(t, err, "获取用户详情失败")

	t.Logf("用户详情:")
	t.Logf("  ID: %d", detail.ID)
	t.Logf("  Username: %s", detail.Username)
	t.Logf("  Email: %s", detail.Email)
	t.Logf("  FullName: %s", detail.FullName)
	t.Logf("  Bio: %s", detail.Bio)
	t.Logf("  Avatar: %s", detail.Avatar)
	t.Logf("  Status: %s", detail.Status)
	t.Logf("  Roles: %d 个", len(detail.Roles))
	for _, r := range detail.Roles {
		t.Logf("    - %s (%s)", r.DisplayName, r.Name)
	}
	t.Logf("  CreatedAt: %s", detail.CreatedAt.Format("2006-01-02 15:04:05"))
}

// TestUsersAPIUpdate 测试更新用户。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestUsersAPIUpdate ./internal/manualtest/
func TestUsersAPIUpdate(t *testing.T) {
	helper.SkipIfNotManual(t)

	c := helper.NewClient()

	t.Log("准备工作: 登录账户")
	_, err := c.Login("admin", "admin123")
	require.NoError(t, err, "登录失败")
	t.Log("  登录成功")

	// 先创建一个测试用户
	timestamp := time.Now().Unix()
	username := fmt.Sprintf("%supdate_%d", userTestPrefix, timestamp)
	email := fmt.Sprintf("%supdate_%d@test.com", userTestPrefix, timestamp)

	t.Logf("\n准备: 创建测试用户 %s...", username)
	createReq := user.CreateDTO{
		Username: username,
		Email:    email,
		Password: "test123456",
		FullName: "待更新用户",
	}
	created, err := helper.Post[user.UserWithRolesDTO](c, "/api/users", createReq)
	require.NoError(t, err, "创建用户失败")
	t.Logf("  创建成功, ID: %d", created.ID)

	// 确保清理
	t.Cleanup(func() {
		if deleteErr := c.Delete(fmt.Sprintf("/api/users/%d", created.ID)); deleteErr != nil {
			t.Logf("清理用户失败: %v", deleteErr)
		}
	})

	// 更新用户
	t.Log("\n测试: 更新用户...")
	newFullName := "已更新的用户名"
	newBio := "这是更新后的个人简介"
	updateReq := user.UpdateDTO{
		FullName: &newFullName,
		Bio:      &newBio,
	}

	resp, err := c.R().
		SetBody(updateReq).
		Put(fmt.Sprintf("/api/users/%d", created.ID))
	require.NoError(t, err, "更新用户失败")
	require.False(t, resp.IsError(), "更新用户失败: 状态码 %d", resp.StatusCode())
	t.Log("  更新成功!")

	// 验证更新
	t.Log("\n验证更新结果...")
	updated, err := helper.Get[user.UserWithRolesDTO](c, fmt.Sprintf("/api/users/%d", created.ID), nil)
	require.NoError(t, err, "获取更新后用户失败")
	assert.Equal(t, newFullName, updated.FullName, "FullName 更新失败")
	assert.Equal(t, newBio, updated.Bio, "Bio 更新失败")
	if updated.FullName == newFullName {
		t.Logf("  ✓ FullName: %s", updated.FullName)
	}
	if updated.Bio == newBio {
		t.Logf("  ✓ Bio: %s", updated.Bio)
	}
}

// TestUsersAPIDelete 测试删除用户。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestUsersAPIDelete ./internal/manualtest/
func TestUsersAPIDelete(t *testing.T) {
	helper.SkipIfNotManual(t)

	c := helper.NewClient()

	t.Log("准备工作: 登录账户")
	_, err := c.Login("admin", "admin123")
	require.NoError(t, err, "登录失败")
	t.Log("  登录成功")

	// 先创建一个测试用户
	timestamp := time.Now().Unix()
	username := fmt.Sprintf("%sdelete_%d", userTestPrefix, timestamp)
	email := fmt.Sprintf("%sdelete_%d@test.com", userTestPrefix, timestamp)

	t.Logf("\n准备: 创建测试用户 %s...", username)
	createReq := user.CreateDTO{
		Username: username,
		Email:    email,
		Password: "test123456",
		FullName: "待删除用户",
	}
	created, err := helper.Post[user.UserWithRolesDTO](c, "/api/users", createReq)
	require.NoError(t, err, "创建用户失败")
	t.Logf("  创建成功, ID: %d", created.ID)

	// 删除用户
	t.Log("\n测试: 删除用户...")
	err = c.Delete(fmt.Sprintf("/api/users/%d", created.ID))
	require.NoError(t, err, "删除用户失败")
	t.Log("  删除成功!")

	// 验证删除
	t.Log("\n验证: 确认用户已删除...")
	_, err = helper.Get[user.UserWithRolesDTO](c, fmt.Sprintf("/api/users/%d", created.ID), nil)
	require.Error(t, err, "用户应该已被删除，但仍能获取")
	t.Log("  ✓ 用户已成功删除")
}

// TestUsersAPINotFound 测试获取不存在的用户。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestUsersAPINotFound ./internal/manualtest/
func TestUsersAPINotFound(t *testing.T) {
	helper.SkipIfNotManual(t)

	c := helper.NewClient()

	t.Log("登录账户...")
	_, err := c.Login("admin", "admin123")
	require.NoError(t, err, "登录失败")

	t.Log("\n获取不存在的用户...")
	nonExistentID := 999999
	_, err = helper.Get[user.UserWithRolesDTO](c, fmt.Sprintf("/api/users/%d", nonExistentID), nil)
	require.Error(t, err, "期望获取不存在的用户返回错误，但成功了")
	t.Logf("  ✓ 正确返回错误: %v", err)
}
