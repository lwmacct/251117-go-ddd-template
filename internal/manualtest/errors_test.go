package manualtest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/251117-go-ddd-template/internal/manualtest/helper"
)

// TestNotFoundEndpoints 测试不存在资源的 404 响应（Table-Driven）。
//
// 统一验证各模块的 NotFound 错误处理
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestNotFoundEndpoints ./internal/manualtest/
func TestNotFoundEndpoints(t *testing.T) {
	helper.SkipIfNotManual(t)

	c := helper.NewClient()
	_, err := c.Login("admin", "admin123")
	require.NoError(t, err, "登录失败")

	cases := []struct {
		name       string
		method     string
		endpoint   string
		wantStatus int // 期望的 HTTP 状态码（通常是 404）
	}{
		{
			name:       "不存在的配置",
			method:     "GET",
			endpoint:   "/api/system/settings/non_existent_setting_key_12345",
			wantStatus: 404,
		},
		{
			name:       "不存在的用户",
			method:     "GET",
			endpoint:   "/api/system/users/99999999",
			wantStatus: 404,
		},
		{
			name:       "不存在的角色",
			method:     "GET",
			endpoint:   "/api/system/roles/99999999",
			wantStatus: 404,
		},
		{
			name:       "不存在的菜单",
			method:     "GET",
			endpoint:   "/api/system/menus/99999999",
			wantStatus: 404,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := c.R().Execute(tc.method, tc.endpoint)
			require.NoError(t, err, "请求失败")
			assert.Equal(t, tc.wantStatus, resp.StatusCode(), "状态码不匹配")
			t.Logf("正确返回 %d", resp.StatusCode())
		})
	}
}

// TestUnauthorizedEndpoints 测试未认证访问受保护资源（Table-Driven）。
//
// 验证各模块的认证保护
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestUnauthorizedEndpoints ./internal/manualtest/
func TestUnauthorizedEndpoints(t *testing.T) {
	helper.SkipIfNotManual(t)

	// 不登录，直接访问
	c := helper.NewClient()

	cases := []struct {
		name       string
		method     string
		endpoint   string
		wantStatus int // 期望 401 Unauthorized
	}{
		{
			name:       "未认证访问用户列表",
			method:     "GET",
			endpoint:   "/api/system/users",
			wantStatus: 401,
		},
		{
			name:       "未认证访问角色列表",
			method:     "GET",
			endpoint:   "/api/system/roles",
			wantStatus: 401,
		},
		{
			name:       "未认证访问个人资料",
			method:     "GET",
			endpoint:   "/api/user/profile",
			wantStatus: 401,
		},
		{
			name:       "未认证访问配置列表",
			method:     "GET",
			endpoint:   "/api/system/settings",
			wantStatus: 401,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := c.R().Execute(tc.method, tc.endpoint)
			require.NoError(t, err, "请求失败")
			assert.Equal(t, tc.wantStatus, resp.StatusCode(), "状态码不匹配")
			t.Logf("正确返回 %d Unauthorized", resp.StatusCode())
		})
	}
}
