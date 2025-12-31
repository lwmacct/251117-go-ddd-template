package manualtest

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/lwmacct/251117-go-ddd-template/internal/application/role"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/setting"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/user"
)

// CreateTestUser 创建测试用户并自动注册清理。
// 测试结束时会自动删除创建的用户。
func CreateTestUser(t *testing.T, c *Client, prefix string) *user.UserWithRolesDTO {
	t.Helper()

	username := fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
	req := user.CreateDTO{
		Username: username,
		Email:    username + "@test.local",
		Password: "test123456",
		FullName: "测试用户",
	}

	resp, err := Post[user.UserWithRolesDTO](c, "/api/system/users", req)
	require.NoError(t, err, "创建测试用户失败: %s", username)

	userID := resp.ID
	t.Cleanup(func() {
		if userID > 0 {
			_ = c.Delete(fmt.Sprintf("/api/system/users/%d", userID))
		}
	})

	return resp
}

// CreateTestRole 创建测试角色并自动注册清理。
// 测试结束时会自动删除创建的角色。
func CreateTestRole(t *testing.T, c *Client, prefix string) *role.CreateResultDTO {
	t.Helper()

	roleName := fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
	req := role.CreateDTO{
		Name:        roleName,
		DisplayName: "测试角色",
		Description: "自动创建的测试角色",
	}

	resp, err := Post[role.CreateResultDTO](c, "/api/system/roles", req)
	require.NoError(t, err, "创建测试角色失败: %s", roleName)

	roleID := resp.RoleID
	t.Cleanup(func() {
		if roleID > 0 {
			_ = c.Delete(fmt.Sprintf("/api/system/roles/%d", roleID))
		}
	})

	return resp
}

// CreateTestUserWithCleanupControl 创建测试用户，返回清理控制函数。
// 当测试本身需要删除用户时使用，删除成功后调用返回的 markDeleted 函数。
func CreateTestUserWithCleanupControl(t *testing.T, c *Client, prefix string) (*user.UserWithRolesDTO, func()) {
	t.Helper()

	username := fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
	req := user.CreateDTO{
		Username: username,
		Email:    username + "@test.local",
		Password: "test123456",
		FullName: "测试用户",
	}

	result, err := Post[user.UserWithRolesDTO](c, "/api/system/users", req)
	require.NoError(t, err, "创建测试用户失败: %s", username)

	userID := result.ID
	deleted := false

	t.Cleanup(func() {
		if !deleted && userID > 0 {
			_ = c.Delete(fmt.Sprintf("/api/system/users/%d", userID))
		}
	})

	return result, func() { deleted = true }
}

// CreateTestRoleWithCleanupControl 创建测试角色，返回清理控制函数。
// 当测试本身需要删除角色时使用，删除成功后调用返回的 markDeleted 函数。
func CreateTestRoleWithCleanupControl(t *testing.T, c *Client, prefix string) (*role.CreateResultDTO, func()) {
	t.Helper()

	roleName := fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
	req := role.CreateDTO{
		Name:        roleName,
		DisplayName: "测试角色",
		Description: "自动创建的测试角色",
	}

	result, err := Post[role.CreateResultDTO](c, "/api/system/roles", req)
	require.NoError(t, err, "创建测试角色失败: %s", roleName)

	roleID := result.RoleID
	deleted := false

	t.Cleanup(func() {
		if !deleted && roleID > 0 {
			_ = c.Delete(fmt.Sprintf("/api/system/roles/%d", roleID))
		}
	})

	return result, func() { deleted = true }
}

// CreateTestSetting 创建测试配置并自动注册清理。
// 测试结束时会自动删除创建的配置。
func CreateTestSetting(t *testing.T, c *Client, prefix string) *setting.SettingDTO {
	t.Helper()

	key := fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
	createReq := map[string]any{
		"key":           key,
		"default_value": "测试值",
		"category_id":   1, // general 分类
		"group":         "test",
		"value_type":    "string",
		"label":         "测试配置",
	}

	result, err := Post[setting.SettingDTO](c, "/api/system/settings", createReq)
	require.NoError(t, err, "创建测试配置失败: %s", key)

	t.Cleanup(func() {
		_ = c.Delete("/api/system/settings/" + key)
	})

	return result
}

// CreateTestSettingWithCleanupControl 创建测试配置，返回清理控制函数。
// 当测试本身需要删除配置时使用，删除成功后调用返回的 markDeleted 函数。
func CreateTestSettingWithCleanupControl(t *testing.T, c *Client, prefix string) (*setting.SettingDTO, func()) {
	t.Helper()

	key := fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
	createReq := map[string]any{
		"key":           key,
		"default_value": "测试值",
		"category_id":   1, // general 分类
		"group":         "test",
		"value_type":    "string",
		"label":         "测试配置",
	}

	result, err := Post[setting.SettingDTO](c, "/api/system/settings", createReq)
	require.NoError(t, err, "创建测试配置失败: %s", key)

	deleted := false
	t.Cleanup(func() {
		if !deleted {
			_ = c.Delete("/api/system/settings/" + key)
		}
	})

	return result, func() { deleted = true }
}
