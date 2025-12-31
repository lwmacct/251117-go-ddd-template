package user_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/251117-go-ddd-template/internal/application/setting"
	"github.com/lwmacct/251117-go-ddd-template/internal/manualtest"
)

// extractSettingsFromSchema 从层级结构的 Schema 中提取扁平化的配置列表
func extractSettingsFromSchema(schema []setting.SettingsCategoryDTO) []setting.SettingsItemDTO {
	var result []setting.SettingsItemDTO
	for _, cat := range schema {
		for _, grp := range cat.Groups {
			result = append(result, grp.Settings...)
		}
	}
	return result
}

// findBooleanSettingFromSchema 从 Schema 中查找第一个 boolean 类型配置（无复杂验证规则）
func findBooleanSettingFromSchema(schema []setting.SettingsCategoryDTO) *setting.SettingsItemDTO {
	for _, cat := range schema {
		for _, grp := range cat.Groups {
			for i := range grp.Settings {
				if grp.Settings[i].ValueType == "boolean" {
					return &grp.Settings[i]
				}
			}
		}
	}
	return nil
}

// findTwoBooleanSettingsFromSchema 从 Schema 中查找两个 boolean 类型配置
func findTwoBooleanSettingsFromSchema(schema []setting.SettingsCategoryDTO) (*setting.SettingsItemDTO, *setting.SettingsItemDTO) {
	var settings []*setting.SettingsItemDTO
	for _, cat := range schema {
		for _, grp := range cat.Groups {
			for i := range grp.Settings {
				if grp.Settings[i].ValueType == "boolean" {
					settings = append(settings, &grp.Settings[i])
					if len(settings) >= 2 {
						return settings[0], settings[1]
					}
				}
			}
		}
	}
	if len(settings) >= 2 {
		return settings[0], settings[1]
	}
	if len(settings) == 1 {
		return settings[0], nil
	}
	return nil, nil
}

// TestUserSettingsFlow 用户配置完整流程测试。
//
// 测试流程：获取配置列表 → 设置配置 → 验证 IsCustomized → 重置配置
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestUserSettingsFlow ./internal/manualtest/user/
func TestUserSettingsFlow(t *testing.T) {
	c := manualtest.LoginAsAdmin(t)

	// 测试 1: 获取用户配置 Schema（层级结构）
	t.Log("\n测试 1: 获取用户配置 Schema")
	schema, err := manualtest.Get[[]setting.SettingsCategoryDTO](c, "/api/user/settings", nil)
	require.NoError(t, err, "获取用户配置 Schema 失败")
	t.Logf("  分类数: %d", len(*schema))

	// 从 Schema 中提取第一个 boolean 类型配置（无复杂验证规则）
	testSetting := findBooleanSettingFromSchema(*schema)
	if testSetting == nil {
		t.Log("  ⚠ 没有 boolean 类型的用户配置，跳过后续测试")
		return
	}

	testKey := testSetting.Key
	originalValue := testSetting.Value
	t.Logf("  选取测试配置: %s", testKey)
	t.Logf("  当前值: %v (IsCustomized: %v)", testSetting.Value, testSetting.IsCustomized)

	// 测试 2: 获取单个用户配置
	t.Log("\n测试 2: 获取单个用户配置")
	detail, err := manualtest.Get[setting.UserSettingDTO](c, "/api/user/settings/"+testKey, nil)
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
	updated, err := manualtest.Put[setting.UserSettingDTO](c, "/api/user/settings/"+testKey, setReq)
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
	resetDetail, err := manualtest.Get[setting.UserSettingDTO](c, "/api/user/settings/"+testKey, nil)
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
		_, err = manualtest.Put[setting.UserSettingDTO](c, "/api/user/settings/"+testKey, restoreReq)
		if err != nil {
			t.Logf("  恢复原始值失败: %v", err)
		} else {
			t.Log("  恢复成功")
		}
	}

	t.Log("\n用户配置流程测试完成!")
}

// TestGetUserSettings 测试获取用户配置 Schema（层级结构）。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestGetUserSettings ./internal/manualtest/user/
func TestGetUserSettings(t *testing.T) {
	c := manualtest.LoginAsAdmin(t)

	t.Log("\n获取用户配置 Schema...")
	schema, err := manualtest.Get[[]setting.SettingsCategoryDTO](c, "/api/user/settings", nil)
	require.NoError(t, err, "获取用户配置 Schema 失败")

	t.Logf("Schema 层级结构 (分类数: %d):", len(*schema))
	for _, cat := range *schema {
		t.Logf("  📁 %s (%s)", cat.Label, cat.Category)
		for _, grp := range cat.Groups {
			t.Logf("    📂 %s", grp.Name)
			for _, s := range grp.Settings {
				customIcon := " "
				if s.IsCustomized {
					customIcon = "✓"
				}
				t.Logf("      [%s] %s (%s): %v", customIcon, s.Key, s.ValueType, s.Value)
			}
		}
	}
}

// TestGetUserSetting 测试获取单个用户配置。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestGetUserSetting ./internal/manualtest/user/
func TestGetUserSetting(t *testing.T) {
	c := manualtest.LoginAsAdmin(t)

	// 先获取 Schema，取第一个配置 key
	t.Log("\n获取配置 Schema...")
	schema, err := manualtest.Get[[]setting.SettingsCategoryDTO](c, "/api/user/settings", nil)
	require.NoError(t, err, "获取用户配置 Schema 失败")

	settings := extractSettingsFromSchema(*schema)
	if len(settings) == 0 {
		t.Skip("没有用户配置可供测试")
	}

	testKey := settings[0].Key
	t.Logf("  选取测试配置: %s", testKey)

	// 获取单个配置
	t.Log("\n获取单个用户配置...")
	detail, err := manualtest.Get[setting.UserSettingDTO](c, "/api/user/settings/"+testKey, nil)
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
//	MANUAL=1 go test -v -run TestGetUserSettingsByCategory ./internal/manualtest/user/
func TestGetUserSettingsByCategory(t *testing.T) {
	c := manualtest.LoginAsAdmin(t)

	// 先获取全量 Schema 确定可用的分类
	fullSchema, err := manualtest.Get[[]setting.SettingsCategoryDTO](c, "/api/user/settings", nil)
	require.NoError(t, err, "获取全量 Schema 失败")
	if len(*fullSchema) == 0 {
		t.Skip("没有配置分类可供测试")
	}

	// 选取第一个分类的 Key 进行测试
	testCategory := (*fullSchema)[0].Category
	t.Logf("\n按类别筛选用户配置 (category=%s)...", testCategory)

	schema, err := manualtest.Get[[]setting.SettingsCategoryDTO](c, "/api/user/settings", map[string]string{
		"category": testCategory,
	})
	require.NoError(t, err, "获取用户配置失败")

	// 验证只返回了指定分类
	require.Len(t, *schema, 1, "按分类筛选应只返回 1 个分类")
	assert.Equal(t, testCategory, (*schema)[0].Category, "返回的分类 Key 不匹配")

	settings := extractSettingsFromSchema(*schema)
	t.Logf("category=%s 配置数: %d", testCategory, len(settings))
	for _, s := range settings {
		t.Logf("  %s: %v (自定义: %v)", s.Key, s.Value, s.IsCustomized)
	}
}

// TestSetUserSetting 测试设置单个用户配置。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestSetUserSetting ./internal/manualtest/user/
func TestSetUserSetting(t *testing.T) {
	c := manualtest.LoginAsAdmin(t)

	// 获取 Schema 并找一个 boolean 类型配置（无复杂验证规则）
	t.Log("\n获取可用配置...")
	schema, err := manualtest.Get[[]setting.SettingsCategoryDTO](c, "/api/user/settings", nil)
	require.NoError(t, err, "获取用户配置 Schema 失败")

	testSetting := findBooleanSettingFromSchema(*schema)
	if testSetting == nil {
		t.Skip("没有 boolean 类型的配置可供测试")
	}

	testKey := testSetting.Key
	origCustom := testSetting.IsCustomized
	// 获取当前 boolean 值并取反
	origBool, _ := testSetting.Value.(bool)
	t.Logf("  选取配置: %s (当前值: %v, IsCustomized: %v)", testKey, origBool, origCustom)

	// 注册清理函数
	t.Cleanup(func() {
		// 重置配置（删除用户自定义值）
		_, _ = c.R().Delete("/api/user/settings/" + testKey)
	})

	// 设置新值（取反）
	t.Log("\n设置用户配置...")
	newValue := !origBool
	setReq := map[string]any{"value": newValue}
	updated, err := manualtest.Put[setting.UserSettingDTO](c, "/api/user/settings/"+testKey, setReq)
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
//	MANUAL=1 go test -v -run TestResetUserSetting ./internal/manualtest/user/
func TestResetUserSetting(t *testing.T) {
	c := manualtest.LoginAsAdmin(t)

	// 获取 Schema 并找一个 boolean 类型配置（无复杂验证规则）
	t.Log("\n获取可用配置...")
	schema, err := manualtest.Get[[]setting.SettingsCategoryDTO](c, "/api/user/settings", nil)
	require.NoError(t, err, "获取用户配置 Schema 失败")

	testSetting := findBooleanSettingFromSchema(*schema)
	if testSetting == nil {
		t.Skip("没有 boolean 类型的配置可供测试")
	}

	testKey := testSetting.Key
	defaultValue := testSetting.DefaultValue
	t.Logf("  选取配置: %s (DefaultValue: %v)", testKey, defaultValue)

	// 先设置一个自定义值（取反）
	t.Log("\n先设置自定义值...")
	origBool, _ := testSetting.Value.(bool)
	setReq := map[string]any{"value": !origBool}
	_, err = manualtest.Put[setting.UserSettingDTO](c, "/api/user/settings/"+testKey, setReq)
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
	resetDetail, err := manualtest.Get[setting.UserSettingDTO](c, "/api/user/settings/"+testKey, nil)
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
//	MANUAL=1 go test -v -run TestBatchSetUserSettings ./internal/manualtest/user/
func TestBatchSetUserSettings(t *testing.T) {
	c := manualtest.LoginAsAdmin(t)

	// 获取 Schema 并找两个 boolean 类型配置（无复杂验证规则）
	t.Log("\n获取可用配置...")
	schema, err := manualtest.Get[[]setting.SettingsCategoryDTO](c, "/api/user/settings", nil)
	require.NoError(t, err, "获取用户配置 Schema 失败")

	setting1, setting2 := findTwoBooleanSettingsFromSchema(*schema)
	if setting1 == nil || setting2 == nil {
		t.Skip("需要至少 2 个 boolean 类型的配置才能测试批量设置")
	}

	t.Logf("  选取配置: %s, %s", setting1.Key, setting2.Key)

	// 获取当前 boolean 值
	origBool1, _ := setting1.Value.(bool)
	origBool2, _ := setting2.Value.(bool)

	// 注册清理函数（确保即使测试失败也会执行）
	t.Cleanup(func() {
		// 重置配置
		_, _ = c.R().Delete("/api/user/settings/" + setting1.Key)
		_, _ = c.R().Delete("/api/user/settings/" + setting2.Key)
	})

	// 批量设置（取反）
	t.Log("\n测试: 批量设置用户配置...")
	batchReq := map[string]any{
		"settings": []map[string]any{
			{"key": setting1.Key, "value": !origBool1},
			{"key": setting2.Key, "value": !origBool2},
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
		detail, getErr := manualtest.Get[setting.UserSettingDTO](c, "/api/user/settings/"+key, nil)
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
//	MANUAL=1 go test -v -run TestUserSettingNotFound ./internal/manualtest/user/
func TestUserSettingNotFound(t *testing.T) {
	c := manualtest.LoginAsAdmin(t)

	t.Log("\n获取不存在的用户配置...")
	nonExistentKey := "non_existent_user_setting_key_12345"
	_, err := manualtest.Get[setting.UserSettingDTO](c, "/api/user/settings/"+nonExistentKey, nil)
	require.Error(t, err, "期望获取不存在的配置返回错误，但成功了")
	t.Logf("  ✓ 正确返回错误: %v", err)
}
