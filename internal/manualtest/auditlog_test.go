package manualtest

import (
	"fmt"
	"testing"

	"github.com/lwmacct/251117-go-ddd-template/internal/application/auditlog"
	"github.com/lwmacct/251117-go-ddd-template/internal/manualtest/helper"
)

// TestListAuditLogs 测试获取审计日志列表。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestListAuditLogs ./internal/manualtest/
func TestListAuditLogs(t *testing.T) {
	helper.SkipIfNotManual(t)

	c := helper.NewClient()

	t.Log("准备工作: 登录管理员账户")
	_, err := c.Login("admin", "admin123")
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}
	t.Log("  登录成功")

	t.Log("\n获取审计日志列表...")
	logs, meta, err := helper.GetList[auditlog.AuditLogDTO](c, "/api/admin/auditlogs", map[string]string{
		"page":  "1",
		"limit": "10",
	})
	if err != nil {
		t.Fatalf("获取审计日志列表失败: %v", err)
	}

	t.Logf("日志数量: %d", len(logs))
	if meta != nil {
		t.Logf("总数: %d, 总页数: %d", meta.Total, meta.TotalPages)
	}

	// 显示前 5 条日志
	displayCount := min(len(logs), 5)

	for i := range displayCount {
		log := logs[i]
		t.Logf("  - [%d] %s %s (%s) - %s", log.ID, log.Action, log.Resource, log.Status, log.CreatedAt.Format("2006-01-02 15:04:05"))
	}

	if len(logs) > 5 {
		t.Logf("  ... 还有 %d 条日志", len(logs)-5)
	}
}

// TestGetAuditLogDetail 测试获取审计日志详情。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestGetAuditLogDetail ./internal/manualtest/
func TestGetAuditLogDetail(t *testing.T) {
	helper.SkipIfNotManual(t)

	c := helper.NewClient()

	t.Log("准备工作: 登录管理员账户")
	_, err := c.Login("admin", "admin123")
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}
	t.Log("  登录成功")

	// 先获取日志列表，取第一条的 ID
	t.Log("\n步骤 1: 获取日志列表")
	logs, _, err := helper.GetList[auditlog.AuditLogDTO](c, "/api/admin/auditlogs", map[string]string{
		"page":  "1",
		"limit": "1",
	})
	if err != nil {
		t.Fatalf("获取审计日志列表失败: %v", err)
	}

	if len(logs) == 0 {
		t.Skip("没有审计日志可供测试")
		return
	}

	logID := logs[0].ID
	t.Logf("  获取日志 ID: %d", logID)

	// 获取详情
	t.Log("\n步骤 2: 获取日志详情")
	detail, err := helper.Get[auditlog.AuditLogDTO](c, fmt.Sprintf("/api/admin/auditlogs/%d", logID), nil)
	if err != nil {
		t.Fatalf("获取审计日志详情失败: %v", err)
	}

	t.Logf("日志详情:")
	t.Logf("  ID: %d", detail.ID)
	t.Logf("  用户 ID: %d", detail.UserID)
	t.Logf("  操作: %s", detail.Action)
	t.Logf("  资源: %s", detail.Resource)
	t.Logf("  状态: %s", detail.Status)
	t.Logf("  IP 地址: %s", detail.IPAddress)
	t.Logf("  创建时间: %s", detail.CreatedAt.Format("2006-01-02 15:04:05"))
	if detail.Details != "" {
		// 截断过长的详情
		details := detail.Details
		if len(details) > 100 {
			details = details[:100] + "..."
		}
		t.Logf("  详情: %s", details)
	}
}

// TestAuditLogFilters 测试审计日志筛选功能。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestAuditLogFilters ./internal/manualtest/
func TestAuditLogFilters(t *testing.T) {
	helper.SkipIfNotManual(t)

	c := helper.NewClient()

	t.Log("准备工作: 登录管理员账户")
	_, err := c.Login("admin", "admin123")
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}
	t.Log("  登录成功")

	// 测试按操作类型筛选
	t.Log("\n测试 1: 按操作类型筛选 (login)")
	loginLogs, meta, err := helper.GetList[auditlog.AuditLogDTO](c, "/api/admin/auditlogs", map[string]string{
		"action": "login",
		"limit":  "5",
	})
	if err != nil {
		t.Logf("  筛选失败: %v（可能不支持此筛选）", err)
	} else {
		t.Logf("  找到 %d 条登录日志", len(loginLogs))
		if meta != nil {
			t.Logf("  总数: %d", meta.Total)
		}
	}

	// 测试按用户 ID 筛选
	t.Log("\n测试 2: 按用户 ID 筛选 (user_id=1)")
	userLogs, meta, err := helper.GetList[auditlog.AuditLogDTO](c, "/api/admin/auditlogs", map[string]string{
		"user_id": "1",
		"limit":   "5",
	})
	if err != nil {
		t.Logf("  筛选失败: %v（可能不支持此筛选）", err)
	} else {
		t.Logf("  找到 %d 条用户 1 的日志", len(userLogs))
		if meta != nil {
			t.Logf("  总数: %d", meta.Total)
		}
	}

	t.Log("\n审计日志筛选测试完成!")
}
