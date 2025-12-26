package manualtest

import (
	"testing"

	"github.com/lwmacct/251117-go-ddd-template/internal/application/cache"
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

	// 测试 1: 设置缓存（POST /api/cache，body 包含 key/value/ttl）
	t.Log("测试 1: 设置缓存")
	setResult, err := helper.Post[cache.SetCacheResultDTO](c, "/api/cache", cache.SetCacheDTO{
		Key:   cacheKey,
		Value: cacheValue,
		TTL:   60,
	})
	if err != nil {
		t.Skipf("设置缓存失败 (可能功能未开启): %v", err)
		return
	}
	if setResult.Key != cacheKey {
		t.Errorf("返回的 Key 不匹配: 期望 %q, 实际 %q", cacheKey, setResult.Key)
	}
	t.Logf("  设置成功! Key=%s, TTL=%d", setResult.Key, setResult.TTL)

	// 测试 2: 获取缓存（GET /api/cache/:key）
	t.Log("\n测试 2: 获取缓存")
	getResult, err := helper.Get[cache.GetCacheResultDTO](c, "/api/cache/"+cacheKey, nil)
	if err != nil {
		t.Errorf("获取缓存失败: %v", err)
	} else {
		t.Logf("  获取成功! Key=%s, Value=%v", getResult.Key, getResult.Value)
		// 验证获取的值是否与设置的值一致
		if getResult.Value != cacheValue {
			t.Errorf("缓存值不匹配: 期望 %q, 实际 %v", cacheValue, getResult.Value)
		}
	}

	// 测试 3: 删除缓存（DELETE /api/cache/:key）
	t.Log("\n测试 3: 删除缓存")
	if delErr := c.Delete("/api/cache/" + cacheKey); delErr != nil {
		t.Errorf("删除缓存失败: %v", delErr)
	} else {
		t.Logf("  删除成功!")
	}

	// 测试 4: 验证删除后获取缓存应该失败
	t.Log("\n测试 4: 验证删除后缓存不存在")
	_, err = helper.Get[cache.GetCacheResultDTO](c, "/api/cache/"+cacheKey, nil)
	if err == nil {
		t.Errorf("删除后仍能获取缓存，期望失败但成功了")
	} else {
		t.Logf("  验证通过: 删除后无法获取缓存")
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
