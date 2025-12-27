package manualtest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
	require.NoError(t, err, "健康检查请求失败")
	require.False(t, resp.IsError(), "健康检查失败: 状态码 %d", resp.StatusCode())

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
	c := helper.LoginAsAdmin(t)

	t.Log("\n获取系统统计...")
	statsResult, err := helper.Get[stats.StatsDTO](c, "/api/admin/overview/stats", nil)
	require.NoError(t, err, "获取系统统计失败")

	// 验证返回的数据
	assert.GreaterOrEqual(t, statsResult.TotalUsers, int64(0), "总用户数不应为负数")
	assert.GreaterOrEqual(t, statsResult.ActiveUsers, int64(0), "活跃用户数不应为负数")
	assert.GreaterOrEqual(t, statsResult.TotalRoles, int64(0), "总角色数不应为负数")

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
	setResult, err := helper.Post[cache.SetResultDTO](c, "/api/cache", cache.SetDTO{
		Key:   cacheKey,
		Value: cacheValue,
		TTL:   60,
	})
	if err != nil {
		t.Skipf("设置缓存失败 (可能功能未开启): %v", err)
		return
	}
	assert.Equal(t, cacheKey, setResult.Key, "返回的 Key 不匹配")
	t.Logf("  设置成功! Key=%s, TTL=%d", setResult.Key, setResult.TTL)

	// 测试 2: 获取缓存（GET /api/cache/:key）
	t.Log("\n测试 2: 获取缓存")
	getResult, err := helper.Get[cache.GetResultDTO](c, "/api/cache/"+cacheKey, nil)
	require.NoError(t, err, "获取缓存失败")
	t.Logf("  获取成功! Key=%s, Value=%v", getResult.Key, getResult.Value)
	// 验证获取的值是否与设置的值一致
	assert.Equal(t, cacheValue, getResult.Value, "缓存值不匹配")

	// 测试 3: 删除缓存（DELETE /api/cache/:key）
	t.Log("\n测试 3: 删除缓存")
	delErr := c.Delete("/api/cache/" + cacheKey)
	require.NoError(t, delErr, "删除缓存失败")
	t.Logf("  删除成功!")

	// 测试 4: 验证删除后获取缓存应该失败
	t.Log("\n测试 4: 验证删除后缓存不存在")
	_, err = helper.Get[cache.GetResultDTO](c, "/api/cache/"+cacheKey, nil)
	require.Error(t, err, "删除后仍能获取缓存，期望失败但成功了")
	t.Logf("  验证通过: 删除后无法获取缓存")

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
	require.NoError(t, err, "Swagger 请求失败")

	if resp.IsError() {
		t.Logf("Swagger 文档不可用: 状态码 %d", resp.StatusCode())
	} else {
		t.Logf("Swagger 文档可用!")
		t.Logf("  状态码: %d", resp.StatusCode())
		t.Logf("  内容长度: %d bytes", len(resp.Body()))
	}
}
