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

	// 用于清理的变量
	var settingKey string

	// 确保测试结束时清理资源
	t.Cleanup(func() {
		if settingKey != "" {
			_ = c.Delete("/api/admin/settings/" + settingKey)
		}
	})

	// 测试 1: 获取设置列表
	t.Log("\n测试 1: 获取设置列表")
	settings, err := helper.Get[[]setting.SettingDTO](c, "/api/admin/settings", nil)
	if err != nil {
		t.Fatalf("获取设置列表失败: %v", err)
	}
	initialCount := len(*settings)
	t.Logf("  设置数量: %d", initialCount)

	// 测试 2: 创建设置
	t.Log("\n测试 2: 创建设置")
	settingKey = fmt.Sprintf("test_setting_%d", time.Now().Unix())
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

	// 验证创建结果（Create API 只返回 ID，完整字段在 Get 时验证）
	if created.ID == 0 {
		t.Fatal("创建设置后 ID 不应为 0")
	}
	t.Logf("  创建成功! 设置 ID: %d", created.ID)

	// 测试 3: 获取单个设置并验证字段
	t.Log("\n测试 3: 获取单个设置并验证字段")
	detail, err := helper.Get[setting.SettingDTO](c, "/api/admin/settings/"+settingKey, nil)
	if err != nil {
		t.Fatalf("获取设置详情失败: %v", err)
	}
	if detail.ID != created.ID {
		t.Errorf("详情 ID 不匹配: 期望 %d, 实际 %d", created.ID, detail.ID)
	}
	if detail.Key != settingKey {
		t.Errorf("Key 不匹配: 期望 %q, 实际 %q", settingKey, detail.Key)
	}
	if detail.Value != "test_value" {
		t.Errorf("Value 不匹配: 期望 %q, 实际 %q", "test_value", detail.Value)
	}
	if detail.Category != "test" {
		t.Errorf("Category 不匹配: 期望 %q, 实际 %q", "test", detail.Category)
	}
	if detail.ValueType != "string" {
		t.Errorf("ValueType 不匹配: 期望 %q, 实际 %q", "string", detail.ValueType)
	}
	if detail.Label != "测试设置" {
		t.Errorf("Label 不匹配: 期望 %q, 实际 %q", "测试设置", detail.Label)
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
	if updated.Value != "updated_value" {
		t.Errorf("更新后 Value 不匹配: 期望 %q, 实际 %q", "updated_value", updated.Value)
	}
	if updated.Label != "更新后的测试设置" {
		t.Errorf("更新后 Label 不匹配: 期望 %q, 实际 %q", "更新后的测试设置", updated.Label)
	}
	// 验证未更新的字段保持不变
	if updated.Key != settingKey {
		t.Errorf("未更新的 Key 不应改变: 期望 %q, 实际 %q", settingKey, updated.Key)
	}
	if updated.Category != "test" {
		t.Errorf("未更新的 Category 不应改变: 期望 %q, 实际 %q", "test", updated.Category)
	}
	if updated.ValueType != "string" {
		t.Errorf("未更新的 ValueType 不应改变: 期望 %q, 实际 %q", "string", updated.ValueType)
	}
	t.Logf("  更新成功! 新值: %s, 新标签: %s", updated.Value, updated.Label)

	// 测试 5: 删除设置
	t.Log("\n测试 5: 删除设置")
	err = c.Delete("/api/admin/settings/" + settingKey)
	if err != nil {
		t.Fatalf("删除设置失败: %v", err)
	}
	settingKey = "" // 已删除，清理时不需要再删
	t.Log("  删除成功!")

	// 测试 6: 验证删除后获取应失败
	t.Log("\n测试 6: 验证删除后获取应返回 404")
	_, err = helper.Get[setting.SettingDTO](c, "/api/admin/settings/"+createReq.Key, nil)
	if err == nil {
		t.Error("删除后获取应返回错误")
	} else {
		t.Logf("  正确返回错误: %v", err)
	}

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

// TestBatchUpdateSettings 测试批量更新设置。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestBatchUpdateSettings ./internal/manualtest/
func TestBatchUpdateSettings(t *testing.T) {
	helper.SkipIfNotManual(t)

	c := helper.NewClient()

	t.Log("准备工作: 登录管理员账户")
	_, err := c.Login("admin", "admin123")
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}
	t.Log("  登录成功")

	// 用于清理的变量
	timestamp := time.Now().Unix()
	key1 := fmt.Sprintf("batch_test_%d_1", timestamp)
	key2 := fmt.Sprintf("batch_test_%d_2", timestamp)

	// 确保测试结束时清理资源
	t.Cleanup(func() {
		_ = c.Delete("/api/admin/settings/" + key1)
		_ = c.Delete("/api/admin/settings/" + key2)
	})

	// 步骤 1: 创建两个测试设置
	t.Log("\n步骤 1: 创建测试设置")
	for _, key := range []string{key1, key2} {
		createReq := handler.CreateSettingRequest{
			Key:       key,
			Value:     "original_value",
			Category:  "test",
			ValueType: "string",
			Label:     "批量更新测试",
		}
		_, createErr := helper.Post[setting.SettingDTO](c, "/api/admin/settings", createReq)
		if createErr != nil {
			t.Fatalf("创建设置 %s 失败: %v", key, createErr)
		}
		t.Logf("  创建设置: %s", key)
	}

	// 步骤 2: 批量更新设置
	t.Log("\n步骤 2: 批量更新设置")
	batchReq := map[string]any{
		"settings": []map[string]string{
			{"key": key1, "value": "updated_value_1"},
			{"key": key2, "value": "updated_value_2"},
		},
	}

	resp, err := c.R().
		SetBody(batchReq).
		Post("/api/admin/settings/batch")
	if err != nil {
		t.Fatalf("批量更新请求失败: %v", err)
	}
	if resp.IsError() {
		t.Fatalf("批量更新失败，状态码: %d, 响应: %s", resp.StatusCode(), resp.String())
	}
	t.Log("  批量更新成功!")

	// 步骤 3: 验证更新结果
	t.Log("\n步骤 3: 验证更新结果")
	detail1, err := helper.Get[setting.SettingDTO](c, "/api/admin/settings/"+key1, nil)
	if err != nil {
		t.Fatalf("获取设置 %s 失败: %v", key1, err)
	}
	if detail1.Value != "updated_value_1" {
		t.Errorf("设置 %s 值未更新，期望 %q，实际 %q", key1, "updated_value_1", detail1.Value)
	} else {
		t.Logf("  %s = %s ✓", key1, detail1.Value)
	}

	detail2, err := helper.Get[setting.SettingDTO](c, "/api/admin/settings/"+key2, nil)
	if err != nil {
		t.Fatalf("获取设置 %s 失败: %v", key2, err)
	}
	if detail2.Value != "updated_value_2" {
		t.Errorf("设置 %s 值未更新，期望 %q，实际 %q", key2, "updated_value_2", detail2.Value)
	} else {
		t.Logf("  %s = %s ✓", key2, detail2.Value)
	}

	t.Log("\n批量更新设置测试完成!")
}

// TestSettingsSchema 测试获取设置 Schema（层级结构）。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestSettingsSchema ./internal/manualtest/
func TestSettingsSchema(t *testing.T) {
	helper.SkipIfNotManual(t)

	c := helper.NewClient()

	t.Log("准备工作: 登录管理员账户")
	_, err := c.Login("admin", "admin123")
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}
	t.Log("  登录成功")

	// 测试: 获取设置 Schema
	t.Log("\n测试: 获取设置 Schema")
	schema, err := helper.Get[[]setting.SchemaCategoryDTO](c, "/api/admin/settings/schema", nil)
	if err != nil {
		t.Fatalf("获取设置 Schema 失败: %v", err)
	}

	t.Logf("  返回 %d 个分类", len(*schema))

	// 验证期望的分类存在
	expectedCategories := map[string]bool{
		"general":      false,
		"security":     false,
		"notification": false,
		"backup":       false,
	}

	for _, cat := range *schema {
		t.Logf("\n[%s] %s (图标: %s)", cat.Category, cat.Label, cat.Icon)

		if _, exists := expectedCategories[cat.Category]; exists {
			expectedCategories[cat.Category] = true
		}

		// 验证有 Groups
		if len(cat.Groups) == 0 {
			t.Errorf("  分类 %s 没有分组", cat.Category)
		}

		for _, group := range cat.Groups {
			t.Logf("  └─ [%s] %s (%d 个设置)", group.Group, group.Label, len(group.Settings))

			// 验证 Settings 有 UI 元数据
			for _, s := range group.Settings {
				if s.UIConfig.InputType == "" {
					t.Errorf("    设置 %s 缺少 ui_config.input_type", s.Key)
				}
				t.Logf("      - %s (%s) %s", s.Key, s.UIConfig.InputType, s.UIConfig.Hint)
			}
		}
	}

	// 检查所有期望的分类是否都存在
	for cat, found := range expectedCategories {
		if !found {
			t.Errorf("期望的分类 %s 未找到", cat)
		}
	}

	t.Log("\nSchema 测试完成!")
}

// TestSettingsUIMetadata 测试设置的 UI 元数据字段。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestSettingsUIMetadata ./internal/manualtest/
func TestSettingsUIMetadata(t *testing.T) {
	helper.SkipIfNotManual(t)

	c := helper.NewClient()

	t.Log("准备工作: 登录管理员账户")
	_, err := c.Login("admin", "admin123")
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}
	t.Log("  登录成功")

	// 获取 Schema
	schema, err := helper.Get[[]setting.SchemaCategoryDTO](c, "/api/admin/settings/schema", nil)
	if err != nil {
		t.Fatalf("获取设置 Schema 失败: %v", err)
	}

	// 构建 key -> setting 映射
	settingsMap := make(map[string]setting.SchemaSettingDTO)
	for _, cat := range *schema {
		for _, group := range cat.Groups {
			for _, s := range group.Settings {
				settingsMap[s.Key] = s
			}
		}
	}

	// 测试 1: 验证 select 类型有 options
	t.Log("\n测试 1: 验证 select 类型设置有 options")
	timezone, ok := settingsMap["general.timezone"]
	if !ok {
		t.Fatal("找不到 general.timezone 设置")
	}
	if timezone.UIConfig.InputType != "select" {
		t.Errorf("general.timezone 应该是 select 类型，实际是 %s", timezone.UIConfig.InputType)
	}
	if timezone.UIConfig.Options == nil {
		t.Error("general.timezone 应该有 ui_config.options")
	} else {
		t.Logf("  general.timezone options: %v", timezone.UIConfig.Options)
	}

	// 测试 2: 验证依赖关系配置
	t.Log("\n测试 2: 验证依赖关系配置")
	emailNotif, ok := settingsMap["notification.enable_email"]
	if !ok {
		t.Fatal("找不到 notification.enable_email 设置")
	}
	if emailNotif.UIConfig.DependsOn == nil {
		t.Error("notification.enable_email 应该有 ui_config.depends_on 配置")
	} else {
		t.Logf("  notification.enable_email depends_on: %v", emailNotif.UIConfig.DependsOn)
	}

	// 测试 3: 验证 number 类型有 validation
	t.Log("\n测试 3: 验证 number 类型设置有 validation")
	pwdMinLen, ok := settingsMap["security.password_min_length"]
	if !ok {
		t.Fatal("找不到 security.password_min_length 设置")
	}
	if pwdMinLen.UIConfig.InputType != "number" {
		t.Errorf("security.password_min_length 应该是 number 类型，实际是 %s", pwdMinLen.UIConfig.InputType)
	}
	if pwdMinLen.UIConfig.Validation == nil {
		t.Error("security.password_min_length 应该有 ui_config.validation 规则")
	} else {
		t.Logf("  security.password_min_length validation: %v", pwdMinLen.UIConfig.Validation)
	}

	// 测试 4: 验证 switch 类型
	t.Log("\n测试 4: 验证 switch 类型设置")
	enableNotif, ok := settingsMap["notification.enable_notifications"]
	if !ok {
		t.Fatal("找不到 notification.enable_notifications 设置")
	}
	if enableNotif.UIConfig.InputType != "switch" {
		t.Errorf("notification.enable_notifications 应该是 switch 类型，实际是 %s", enableNotif.UIConfig.InputType)
	}
	t.Logf("  notification.enable_notifications hint: %s", enableNotif.UIConfig.Hint)

	// 测试 5: 验证排序字段
	t.Log("\n测试 5: 验证排序字段")
	for _, cat := range *schema {
		for _, group := range cat.Groups {
			prevOrder := -1
			for _, s := range group.Settings {
				if s.Order < prevOrder {
					t.Errorf("设置 %s 排序异常: Order=%d 小于前一个 %d", s.Key, s.Order, prevOrder)
				}
				prevOrder = s.Order
			}
		}
	}
	t.Log("  排序验证通过")

	t.Log("\nUI 元数据测试完成!")
}
