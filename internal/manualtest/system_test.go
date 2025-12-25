package manualtest

import (
	"testing"

	"github.com/lwmacct/251117-go-ddd-template/internal/application/stats"
	"github.com/lwmacct/251117-go-ddd-template/internal/manualtest/helper"
)

// TestHealthCheck 测试健康检查端点。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestHealthCheck ./internal/manualtest/
func TestHealthCheck(t *testing.T) {
	helper.SkipIfNotManual(t)

	c := helper.NewClient()

	t.Log("检查健康状态...")
	resp, err := c.R().Get("/health")
	if err != nil {
		t.Fatalf("健康检查请求失败: %v", err)
	}

	if resp.IsError() {
		t.Fatalf("健康检查失败: 状态码 %d", resp.StatusCode())
	}

	t.Logf("健康检查通过!")
	t.Logf("  状态码: %d", resp.StatusCode())
	t.Logf("  响应: %s", string(resp.Body()))
}

// TestSystemStats 测试系统统计端点。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestSystemStats ./internal/manualtest/
func TestSystemStats(t *testing.T) {
	helper.SkipIfNotManual(t)

	c := helper.NewClient()

	t.Log("准备工作: 登录管理员账户")
	_, err := c.Login("admin", "admin123")
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}
	t.Log("  登录成功")

	t.Log("\n获取系统统计...")
	statsResult, err := helper.Get[stats.StatsDTO](c, "/api/admin/overview/stats", nil)
	if err != nil {
		t.Fatalf("获取系统统计失败: %v", err)
	}

	// 验证返回的数据
	if statsResult.TotalUsers < 0 {
		t.Error("总用户数不应为负数")
	}
	if statsResult.ActiveUsers < 0 {
		t.Error("活跃用户数不应为负数")
	}
	if statsResult.TotalRoles < 0 {
		t.Error("总角色数不应为负数")
	}

	t.Logf("系统统计:")
	t.Logf("  总用户数: %d", statsResult.TotalUsers)
	t.Logf("  活跃用户数: %d", statsResult.ActiveUsers)
	t.Logf("  非活跃用户数: %d", statsResult.InactiveUsers)
	t.Logf("  封禁用户数: %d", statsResult.BannedUsers)
	t.Logf("  总角色数: %d", statsResult.TotalRoles)
	t.Logf("  总权限数: %d", statsResult.TotalPermissions)
	t.Logf("  总菜单数: %d", statsResult.TotalMenus)
	if len(statsResult.RecentAuditLogs) > 0 {
		t.Logf("  最近审计日志: %d 条", len(statsResult.RecentAuditLogs))
	}
}

// TestCacheOperations 测试缓存操作端点。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestCacheOperations ./internal/manualtest/
func TestCacheOperations(t *testing.T) {
	helper.SkipIfNotManual(t)

	c := helper.NewClient()

	cacheKey := "test_cache_key"
	cacheValue := "test_cache_value"

	// 测试 1: 设置缓存
	t.Log("测试 1: 设置缓存")
	resp, err := c.R().
		SetBody(map[string]string{"value": cacheValue}).
		Post("/api/cache/" + cacheKey)
	if err != nil {
		t.Fatalf("设置缓存请求失败: %v", err)
	}
	if resp.IsError() {
		t.Logf("  设置缓存失败 (可能需要认证): 状态码 %d", resp.StatusCode())
	} else {
		t.Logf("  设置成功!")
	}

	// 测试 2: 获取缓存
	t.Log("\n测试 2: 获取缓存")
	resp, err = c.R().Get("/api/cache/" + cacheKey)
	if err != nil {
		t.Fatalf("获取缓存请求失败: %v", err)
	}
	if resp.IsError() {
		t.Logf("  获取缓存失败 (可能需要认证): 状态码 %d", resp.StatusCode())
	} else {
		t.Logf("  获取成功! 响应: %s", string(resp.Body()))
	}

	// 测试 3: 删除缓存
	t.Log("\n测试 3: 删除缓存")
	resp, err = c.R().Delete("/api/cache/" + cacheKey)
	if err != nil {
		t.Fatalf("删除缓存请求失败: %v", err)
	}
	if resp.IsError() {
		t.Logf("  删除缓存失败 (可能需要认证): 状态码 %d", resp.StatusCode())
	} else {
		t.Logf("  删除成功!")
	}

	t.Log("\n缓存操作测试完成!")
}

// TestSwaggerDocs 测试 Swagger 文档端点。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestSwaggerDocs ./internal/manualtest/
func TestSwaggerDocs(t *testing.T) {
	helper.SkipIfNotManual(t)

	c := helper.NewClient()

	t.Log("检查 Swagger 文档...")
	resp, err := c.R().Get("/swagger/index.html")
	if err != nil {
		t.Fatalf("Swagger 请求失败: %v", err)
	}

	if resp.IsError() {
		t.Logf("Swagger 文档不可用: 状态码 %d", resp.StatusCode())
	} else {
		t.Logf("Swagger 文档可用!")
		t.Logf("  状态码: %d", resp.StatusCode())
		t.Logf("  内容长度: %d bytes", len(resp.Body()))
	}
}
