package manualtest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/251117-go-ddd-template/internal/application/setting"
	"github.com/lwmacct/251117-go-ddd-template/internal/manualtest/helper"
)

// TestUserSettingsFlow 用户配置完整流程测试。
//
// 测试流程：获取配置列表 → 设置配置 → 验证 IsCustomized → 重置配置
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestUserSettingsFlow ./internal/manualtest/
func TestUserSettingsFlow(t *testing.T) {
	c := helper.LoginAsAdmin(t)

	// 测试 1: 获取用户配置列表
	t.Log("\n测试 1: 获取用户配置列表")
	settings, err := helper.Get[[]setting.UserSettingDTO](c, "/api/user/settings", nil)
	require.NoError(t, err, "获取用户配置列表失败")
	t.Logf("  用户配置数: %d", len(*settings))

	if len(*settings) == 0 {
		t.Log("  ⚠ 没有用户配置，跳过后续测试")
		t.Log("  提示: 需要先通过种子数据创建系统配置")
		return
	}

	// 选取第一个配置进行测试
	testSetting := (*settings)[0]
	testKey := testSetting.Key
	originalValue := testSetting.Value
	t.Logf("  选取测试配置: %s", testKey)
	t.Logf("  当前值: %v (IsCustomized: %v)", testSetting.Value, testSetting.IsCustomized)

	// 测试 2: 获取单个用户配置
	t.Log("\n测试 2: 获取单个用户配置")
	detail, err := helper.Get[setting.UserSettingDTO](c, "/api/user/settings/"+testKey, nil)
	require.NoError(t, err, "获取用户配置失败")
	t.Logf("  Key: %s", detail.Key)
	t.Logf("  Value: %v", detail.Value)
	t.Logf("  DefaultValue: %v", detail.DefaultValue)
	t.Logf("  IsCustomized: %v", detail.IsCustomized)
	t.Logf("  Label: %s", detail.Label)

	// 测试 3: 设置用户配置
	t.Log("\n测试 3: 设置用户配置")
	var newValue any
	// 根据值类型设置合适的新值
	switch detail.ValueType {
	case "string":
		newValue = "测试自定义值"
	case "number", "integer":
		newValue = 999
	case "boolean":
		// 取反
		if v, ok := detail.Value.(bool); ok {
			newValue = !v
		} else {
			newValue = true
		}
	default:
		newValue = "测试值"
	}

	setReq := map[string]any{
		"value": newValue,
	}
	updated, err := helper.Put[setting.UserSettingDTO](c, "/api/user/settings/"+testKey, setReq)
	require.NoError(t, err, "设置用户配置失败")
	t.Logf("  设置成功!")
	t.Logf("  新 Value: %v", updated.Value)
	t.Logf("  IsCustomized: %v", updated.IsCustomized)

	// 验证 IsCustomized 应该为 true
	assert.True(t, updated.IsCustomized, "设置后 IsCustomized 应该为 true")
	if updated.IsCustomized {
		t.Log("  ✓ IsCustomized 正确设置为 true")
	}

	// 测试 4: 重置用户配置
	t.Log("\n测试 4: 重置用户配置（恢复默认值）")
	resp, err := c.R().Delete("/api/user/settings/" + testKey)
	require.NoError(t, err, "重置用户配置失败")
	require.False(t, resp.IsError(), "重置用户配置失败: 状态码 %d", resp.StatusCode())
	t.Log("  重置成功!")

	// 测试 5: 验证重置结果
	t.Log("\n测试 5: 验证重置结果")
	resetDetail, err := helper.Get[setting.UserSettingDTO](c, "/api/user/settings/"+testKey, nil)
	require.NoError(t, err, "获取重置后配置失败")
	t.Logf("  Value: %v", resetDetail.Value)
	t.Logf("  DefaultValue: %v", resetDetail.DefaultValue)
	t.Logf("  IsCustomized: %v", resetDetail.IsCustomized)

	// 验证 IsCustomized 应该为 false
	assert.False(t, resetDetail.IsCustomized, "重置后 IsCustomized 应该为 false")
	if !resetDetail.IsCustomized {
		t.Log("  ✓ IsCustomized 正确恢复为 false")
	}

	// 如果原来有自定义值，恢复它
	if testSetting.IsCustomized {
		t.Log("\n清理: 恢复原始自定义值...")
		restoreReq := map[string]any{"value": originalValue}
		_, err = helper.Put[setting.UserSettingDTO](c, "/api/user/settings/"+testKey, restoreReq)
		if err != nil {
			t.Logf("  恢复原始值失败: %v", err)
		} else {
			t.Log("  恢复成功")
		}
	}

	t.Log("\n用户配置流程测试完成!")
}

// TestGetUserSettings 测试获取用户配置列表。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestGetUserSettings ./internal/manualtest/
func TestGetUserSettings(t *testing.T) {
	c := helper.LoginAsAdmin(t)

	t.Log("\n获取用户配置列表...")
	settings, err := helper.Get[[]setting.UserSettingDTO](c, "/api/user/settings", nil)
	require.NoError(t, err, "获取用户配置列表失败")

	t.Logf("用户配置数: %d", len(*settings))
	for _, s := range *settings {
		customIcon := " "
		if s.IsCustomized {
			customIcon = "✓"
		}
		t.Logf("  [%s] %s (%s): %v", customIcon, s.Key, s.ValueType, s.Value)
	}
}

// TestGetUserSetting 测试获取单个用户配置。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestGetUserSetting ./internal/manualtest/
func TestGetUserSetting(t *testing.T) {
	c := helper.LoginAsAdmin(t)

	// 先获取配置列表，取第一个 key
	t.Log("\n获取配置列表...")
	settings, err := helper.Get[[]setting.UserSettingDTO](c, "/api/user/settings", nil)
	require.NoError(t, err, "获取用户配置列表失败")
	if len(*settings) == 0 {
		t.Skip("没有用户配置可供测试")
	}

	testKey := (*settings)[0].Key
	t.Logf("  选取测试配置: %s", testKey)

	// 获取单个配置
	t.Log("\n获取单个用户配置...")
	detail, err := helper.Get[setting.UserSettingDTO](c, "/api/user/settings/"+testKey, nil)
	require.NoError(t, err, "获取用户配置失败")

	// 验证字段完整性
	t.Logf("配置详情:")
	t.Logf("  Key: %s", detail.Key)
	t.Logf("  Value: %v", detail.Value)
	t.Logf("  DefaultValue: %v", detail.DefaultValue)
	t.Logf("  ValueType: %s", detail.ValueType)
	t.Logf("  CategoryID: %d", detail.CategoryID)
	t.Logf("  Group: %s", detail.Group)
	t.Logf("  Label: %s", detail.Label)
	t.Logf("  IsCustomized: %v", detail.IsCustomized)
	t.Logf("  Order: %d", detail.Order)

	assert.Equal(t, testKey, detail.Key, "Key 不匹配")
	assert.NotEmpty(t, detail.ValueType, "ValueType 不应为空")
}

// TestGetUserSettingsByCategory 测试按类别筛选用户配置。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestGetUserSettingsByCategory ./internal/manualtest/
func TestGetUserSettingsByCategory(t *testing.T) {
	c := helper.LoginAsAdmin(t)

	t.Log("\n按类别筛选用户配置 (category_id=1)...")
	settings, err := helper.Get[[]setting.UserSettingDTO](c, "/api/user/settings", map[string]string{
		"category_id": "1",
	})
	require.NoError(t, err, "获取用户配置失败")

	t.Logf("category_id=1 配置数: %d", len(*settings))
	for _, s := range *settings {
		assert.Equal(t, uint(1), s.CategoryID, "配置 %s 的 CategoryID 不是 1", s.Key)
		t.Logf("  %s: %v (自定义: %v)", s.Key, s.Value, s.IsCustomized)
	}
}

// TestGetUserSettingsSchema 测试获取用户配置 Schema。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestGetUserSettingsSchema ./internal/manualtest/
func TestGetUserSettingsSchema(t *testing.T) {
	c := helper.LoginAsAdmin(t)

	t.Log("\n获取用户配置 Schema...")
	schema, err := helper.Get[[]setting.SchemaCategoryDTO](c, "/api/user/settings/schema", nil)
	require.NoError(t, err, "获取 Schema 失败")

	t.Logf("Schema 层级结构:")
	for _, cat := range *schema {
		t.Logf("  📁 %s (%s) [icon: %s]", cat.Label, cat.Category, cat.Icon)
		for _, group := range cat.Groups {
			t.Logf("    📂 %s", group.Name)
			for _, s := range group.Settings {
				customIcon := " "
				if s.IsCustomized {
					customIcon = "✓"
				}
				t.Logf("      [%s] %s (%s): %v", customIcon, s.Label, s.Key, s.Value)
			}
		}
	}
}

// TestSetUserSetting 测试设置单个用户配置。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestSetUserSetting ./internal/manualtest/
func TestSetUserSetting(t *testing.T) {
	c := helper.LoginAsAdmin(t)

	// 获取一个 string 类型配置
	t.Log("\n获取可用配置...")
	settings, err := helper.Get[[]setting.UserSettingDTO](c, "/api/user/settings", nil)
	require.NoError(t, err, "获取用户配置列表失败")

	var testSetting *setting.UserSettingDTO
	for _, s := range *settings {
		if s.ValueType == "string" {
			testSetting = &s
			break
		}
	}
	if testSetting == nil {
		t.Skip("没有 string 类型的配置可供测试")
	}

	testKey := testSetting.Key
	origValue := testSetting.Value
	origCustom := testSetting.IsCustomized
	t.Logf("  选取配置: %s (当前值: %v, IsCustomized: %v)", testKey, origValue, origCustom)

	// 注册清理函数
	t.Cleanup(func() {
		if origCustom {
			restoreReq := map[string]any{"value": origValue}
			_, _ = helper.Put[setting.UserSettingDTO](c, "/api/user/settings/"+testKey, restoreReq)
		} else {
			_, _ = c.R().Delete("/api/user/settings/" + testKey)
		}
	})

	// 设置新值
	t.Log("\n设置用户配置...")
	newValue := "测试设置值_" + testKey
	setReq := map[string]any{"value": newValue}
	updated, err := helper.Put[setting.UserSettingDTO](c, "/api/user/settings/"+testKey, setReq)
	require.NoError(t, err, "设置用户配置失败")

	t.Logf("  新 Value: %v", updated.Value)
	t.Logf("  IsCustomized: %v", updated.IsCustomized)

	// 验证
	assert.True(t, updated.IsCustomized, "设置后 IsCustomized 应该为 true")
	if updated.IsCustomized {
		t.Log("  ✓ IsCustomized 正确设置为 true")
	}
}

// TestResetUserSetting 测试重置用户配置（恢复默认值）。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestResetUserSetting ./internal/manualtest/
func TestResetUserSetting(t *testing.T) {
	c := helper.LoginAsAdmin(t)

	// 获取一个 string 类型配置
	t.Log("\n获取可用配置...")
	settings, err := helper.Get[[]setting.UserSettingDTO](c, "/api/user/settings", nil)
	require.NoError(t, err, "获取用户配置列表失败")

	var testSetting *setting.UserSettingDTO
	for _, s := range *settings {
		if s.ValueType == "string" {
			testSetting = &s
			break
		}
	}
	if testSetting == nil {
		t.Skip("没有 string 类型的配置可供测试")
	}

	testKey := testSetting.Key
	defaultValue := testSetting.DefaultValue
	t.Logf("  选取配置: %s (DefaultValue: %v)", testKey, defaultValue)

	// 先设置一个自定义值
	t.Log("\n先设置自定义值...")
	setReq := map[string]any{"value": "临时测试值"}
	_, err = helper.Put[setting.UserSettingDTO](c, "/api/user/settings/"+testKey, setReq)
	require.NoError(t, err, "设置用户配置失败")
	t.Log("  设置成功")

	// 重置配置
	t.Log("\n重置用户配置...")
	resp, err := c.R().Delete("/api/user/settings/" + testKey)
	require.NoError(t, err, "重置用户配置失败")
	require.False(t, resp.IsError(), "重置用户配置失败: 状态码 %d", resp.StatusCode())
	t.Log("  重置成功")

	// 验证重置结果
	t.Log("\n验证重置结果...")
	resetDetail, err := helper.Get[setting.UserSettingDTO](c, "/api/user/settings/"+testKey, nil)
	require.NoError(t, err, "获取重置后配置失败")

	t.Logf("  Value: %v", resetDetail.Value)
	t.Logf("  DefaultValue: %v", resetDetail.DefaultValue)
	t.Logf("  IsCustomized: %v", resetDetail.IsCustomized)

	assert.False(t, resetDetail.IsCustomized, "重置后 IsCustomized 应该为 false")
	if !resetDetail.IsCustomized {
		t.Log("  ✓ IsCustomized 正确恢复为 false")
	}
}

// TestBatchSetUserSettings 测试批量设置用户配置。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestBatchSetUserSettings ./internal/manualtest/
func TestBatchSetUserSettings(t *testing.T) {
	c := helper.LoginAsAdmin(t)

	// 获取可用的配置
	t.Log("\n获取可用配置...")
	settings, err := helper.Get[[]setting.UserSettingDTO](c, "/api/user/settings", nil)
	require.NoError(t, err, "获取用户配置列表失败")

	if len(*settings) < 2 {
		t.Skip("需要至少 2 个配置才能测试批量设置")
	}

	// 选取没有复杂验证规则的配置（string 类型）
	var setting1, setting2 setting.UserSettingDTO
	foundCount := 0
	for _, s := range *settings {
		if s.ValueType == "string" {
			if foundCount == 0 {
				setting1 = s
				foundCount++
			} else if foundCount == 1 {
				setting2 = s
				foundCount++
				break
			}
		}
	}

	if foundCount < 2 {
		t.Skip("需要至少 2 个 string 类型的 general 配置才能测试批量设置")
	}

	t.Logf("  选取配置: %s, %s", setting1.Key, setting2.Key)

	// 保存原始值用于恢复
	origValue1 := setting1.Value
	origValue2 := setting2.Value
	origCustom1 := setting1.IsCustomized
	origCustom2 := setting2.IsCustomized

	// 注册清理函数（确保即使测试失败也会执行）
	t.Cleanup(func() {
		if origCustom1 {
			restoreReq := map[string]any{"value": origValue1}
			_, _ = helper.Put[setting.UserSettingDTO](c, "/api/user/settings/"+setting1.Key, restoreReq)
		} else {
			_, _ = c.R().Delete("/api/user/settings/" + setting1.Key)
		}
		if origCustom2 {
			restoreReq := map[string]any{"value": origValue2}
			_, _ = helper.Put[setting.UserSettingDTO](c, "/api/user/settings/"+setting2.Key, restoreReq)
		} else {
			_, _ = c.R().Delete("/api/user/settings/" + setting2.Key)
		}
	})

	// 批量设置
	t.Log("\n测试: 批量设置用户配置...")
	batchReq := map[string]any{
		"settings": []map[string]any{
			{"key": setting1.Key, "value": "批量测试值1"},
			{"key": setting2.Key, "value": "批量测试值2"},
		},
	}

	resp, err := c.R().
		SetBody(batchReq).
		Post("/api/user/settings/batch")
	require.NoError(t, err, "批量设置失败")
	require.False(t, resp.IsError(), "批量设置失败: 状态码 %d, 响应: %s", resp.StatusCode(), resp.String())
	t.Log("  批量设置成功!")

	// 验证设置结果
	t.Log("\n验证设置结果...")
	for _, key := range []string{setting1.Key, setting2.Key} {
		detail, getErr := helper.Get[setting.UserSettingDTO](c, "/api/user/settings/"+key, nil)
		require.NoError(t, getErr, "获取配置 %s 失败", key)
		assert.True(t, detail.IsCustomized, "配置 %s 的 IsCustomized 应该为 true", key)
		if detail.IsCustomized {
			t.Logf("  ✓ %s = %v (IsCustomized: true)", key, detail.Value)
		}
	}
}

// TestUserSettingNotFound 测试获取不存在的用户配置。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestUserSettingNotFound ./internal/manualtest/
func TestUserSettingNotFound(t *testing.T) {
	c := helper.LoginAsAdmin(t)

	t.Log("\n获取不存在的用户配置...")
	nonExistentKey := "non_existent_user_setting_key_12345"
	_, err := helper.Get[setting.UserSettingDTO](c, "/api/user/settings/"+nonExistentKey, nil)
	require.Error(t, err, "期望获取不存在的配置返回错误，但成功了")
	t.Logf("  ✓ 正确返回错误: %v", err)
}
