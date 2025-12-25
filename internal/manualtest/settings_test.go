package manualtest

import (
	"fmt"
	"testing"
	"time"

	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/handler"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/setting"
	"github.com/lwmacct/251117-go-ddd-template/internal/manualtest/helper"
)

// 使用 time 包生成唯一标识符
var _ = time.Now

// TestSettingsFlow 系统设置完整流程测试。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestSettingsFlow ./internal/manualtest/
func TestSettingsFlow(t *testing.T) {
	helper.SkipIfNotManual(t)

	c := helper.NewClient()

	t.Log("准备工作: 登录管理员账户")
	_, err := c.Login("admin", "admin123")
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}
	t.Log("  登录成功")

	// 测试 1: 获取设置列表
	t.Log("\n测试 1: 获取设置列表")
	settings, err := helper.Get[[]setting.SettingDTO](c, "/api/admin/settings", nil)
	if err != nil {
		t.Fatalf("获取设置列表失败: %v", err)
	}
	t.Logf("  设置数量: %d", len(*settings))

	// 测试 2: 创建设置
	t.Log("\n测试 2: 创建设置")
	settingKey := fmt.Sprintf("test_setting_%d", time.Now().Unix())
	createReq := handler.CreateSettingRequest{
		Key:       settingKey,
		Value:     "test_value",
		Category:  "test",
		ValueType: "string",
		Label:     "测试设置",
	}
	t.Logf("  创建设置: %s", settingKey)

	created, err := helper.Post[setting.SettingDTO](c, "/api/admin/settings", createReq)
	if err != nil {
		t.Fatalf("创建设置失败: %v", err)
	}
	t.Logf("  创建成功! 设置 ID: %d", created.ID)

	// 测试 3: 获取单个设置
	t.Log("\n测试 3: 获取单个设置")
	detail, err := helper.Get[setting.SettingDTO](c, "/api/admin/settings/"+settingKey, nil)
	if err != nil {
		t.Fatalf("获取设置详情失败: %v", err)
	}
	t.Logf("  Key: %s, Value: %s", detail.Key, detail.Value)
	t.Logf("  Category: %s, Label: %s", detail.Category, detail.Label)

	// 测试 4: 更新设置
	t.Log("\n测试 4: 更新设置")
	updateReq := handler.UpdateSettingRequest{
		Value: "updated_value",
		Label: "更新后的测试设置",
	}
	updated, err := helper.Put[setting.SettingDTO](c, "/api/admin/settings/"+settingKey, updateReq)
	if err != nil {
		t.Fatalf("更新设置失败: %v", err)
	}
	t.Logf("  更新成功! 新值: %s", updated.Value)

	// 测试 5: 删除设置
	t.Log("\n测试 5: 删除设置")
	err = c.Delete("/api/admin/settings/" + settingKey)
	if err != nil {
		t.Fatalf("删除设置失败: %v", err)
	}
	t.Log("  删除成功!")

	t.Log("\n系统设置流程测试完成!")
}

// TestListSettings 测试获取设置列表。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestListSettings ./internal/manualtest/
func TestListSettings(t *testing.T) {
	helper.SkipIfNotManual(t)

	c := helper.NewClient()

	t.Log("获取设置列表...")
	_, err := c.Login("admin", "admin123")
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}

	settings, err := helper.Get[[]setting.SettingDTO](c, "/api/admin/settings", nil)
	if err != nil {
		t.Fatalf("获取设置列表失败: %v", err)
	}

	t.Logf("设置数量: %d", len(*settings))

	// 按类别分组显示
	categories := make(map[string][]setting.SettingDTO)
	for _, s := range *settings {
		categories[s.Category] = append(categories[s.Category], s)
	}

	for cat, items := range categories {
		t.Logf("\n[%s] %d 个设置:", cat, len(items))
		for _, s := range items {
			t.Logf("  - %s = %s (%s)", s.Key, s.Value, s.Label)
		}
	}
}

// TestSettingsByCategory 测试按类别获取设置。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestSettingsByCategory ./internal/manualtest/
func TestSettingsByCategory(t *testing.T) {
	helper.SkipIfNotManual(t)

	c := helper.NewClient()

	t.Log("按类别获取设置...")
	_, err := c.Login("admin", "admin123")
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}

	// 测试按 general 类别筛选
	settings, err := helper.Get[[]setting.SettingDTO](c, "/api/admin/settings", map[string]string{
		"category": "general",
	})
	if err != nil {
		t.Logf("  筛选失败: %v（可能该类别不存在）", err)
	} else {
		t.Logf("general 类别设置数量: %d", len(*settings))
		for _, s := range *settings {
			t.Logf("  - %s = %s", s.Key, s.Value)
		}
	}
}
