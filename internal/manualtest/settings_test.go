package manualtest

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/251117-go-ddd-template/internal/application/setting"
	"github.com/lwmacct/251117-go-ddd-template/internal/manualtest/helper"
)

// 测试配置前缀，用于隔离测试数据
const settingTestPrefix = "test_setting_"

// TestSettingsFlow 系统配置完整流程测试。
//
// 测试 CRUD 完整流程：创建 → 获取 → 更新 → 批量更新 → 删除
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestSettingsFlow ./internal/manualtest/
func TestSettingsFlow(t *testing.T) {
	c := helper.LoginAsAdmin(t)

	// 测试 1: 创建配置
	t.Log("\n测试 1: 创建配置")
	settingKey := fmt.Sprintf("%s%d", settingTestPrefix, time.Now().Unix())
	createReq := map[string]any{
		"key":           settingKey,
		"default_value": "测试值",
		"category_id":   1, // 使用 general 分类（ID=1）
		"group":         "basic",
		"value_type":    "string",
		"label":         "测试配置",
		"ui_config":     `{"input_type":"text"}`,
		"order":         100,
	}
	t.Logf("  创建配置: %s", settingKey)

	created, err := helper.Post[setting.SettingDTO](c, "/api/admin/settings", createReq)
	require.NoError(t, err, "创建配置失败")
	require.NotZero(t, created.ID, "创建的配置 ID 为 0")
	assert.Equal(t, settingKey, created.Key, "Key 不匹配")
	t.Logf("  创建成功! ID: %d, Key: %s", created.ID, created.Key)

	// 确保清理
	t.Cleanup(func() {
		if deleteErr := c.Delete("/api/admin/settings/" + settingKey); deleteErr != nil {
			t.Logf("清理配置失败: %v", deleteErr)
		}
	})

	// 测试 2: 获取单个配置
	t.Log("\n测试 2: 获取单个配置")
	detail, err := helper.Get[setting.SettingDTO](c, "/api/admin/settings/"+settingKey, nil)
	require.NoError(t, err, "获取配置失败")
	assert.Equal(t, settingKey, detail.Key, "Key 不匹配")
	assert.Equal(t, "测试配置", detail.Label, "Label 不匹配")
	t.Logf("  Key: %s", detail.Key)
	t.Logf("  Label: %s", detail.Label)
	t.Logf("  CategoryID: %d", detail.CategoryID)
	t.Logf("  DefaultValue: %v", detail.DefaultValue)

	// 测试 3: 更新配置
	t.Log("\n测试 3: 更新配置")
	updateReq := map[string]any{
		"default_value": "更新后的值",
		"label":         "更新后的标签",
		"order":         200,
	}
	updated, err := helper.Put[setting.SettingDTO](c, "/api/admin/settings/"+settingKey, updateReq)
	require.NoError(t, err, "更新配置失败")
	assert.Equal(t, "更新后的标签", updated.Label, "Label 更新失败")
	t.Logf("  更新成功!")
	t.Logf("  新 Label: %s", updated.Label)
	t.Logf("  新 DefaultValue: %v", updated.DefaultValue)

	// 测试 4: 获取配置列表
	t.Log("\n测试 4: 获取配置列表")
	settings, err := helper.Get[[]setting.SettingDTO](c, "/api/admin/settings", nil)
	require.NoError(t, err, "获取配置列表失败")
	t.Logf("  配置总数: %d", len(*settings))
	// 验证刚创建的配置在列表中
	found := false
	for _, s := range *settings {
		if s.Key == settingKey {
			found = true
			break
		}
	}
	assert.True(t, found, "创建的配置不在列表中")
	t.Log("  ✓ 创建的配置存在于列表中")

	// 测试 5: 获取配置 Schema
	t.Log("\n测试 5: 获取配置 Schema")
	schema, err := helper.Get[[]setting.SchemaCategoryDTO](c, "/api/admin/settings/schema", nil)
	require.NoError(t, err, "获取 Schema 失败")
	t.Logf("  Schema 分类数: %d", len(*schema))
	for _, cat := range *schema {
		t.Logf("    - %s (%s): %d 分组", cat.Label, cat.Category, len(cat.Groups))
	}

	t.Log("\n系统配置流程测试完成!")
}

// TestGetSettingsWithFilters 测试配置列表查询（Table-Driven）。
//
// 覆盖场景：获取所有配置、按 general/security 类别筛选
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestGetSettingsWithFilters ./internal/manualtest/
func TestGetSettingsWithFilters(t *testing.T) {
	c := helper.LoginAsAdmin(t)

	cases := []struct {
		name     string
		query    map[string]string
		validate func(t *testing.T, settings []setting.SettingDTO)
	}{
		{
			name:  "获取所有配置",
			query: nil,
			validate: func(t *testing.T, settings []setting.SettingDTO) {
				t.Helper()
				assert.NotEmpty(t, settings, "配置列表为空")
				t.Logf("配置总数: %d", len(settings))
				for _, s := range settings {
					t.Logf("  [%d] %s (%s): %v", s.ID, s.Key, s.ValueType, s.DefaultValue)
				}
			},
		},
		{
			name:  "按 category_id=1 筛选（general）",
			query: map[string]string{"category_id": "1"},
			validate: func(t *testing.T, settings []setting.SettingDTO) {
				t.Helper()
				for _, s := range settings {
					assert.Equal(t, uint(1), s.CategoryID, "配置 %s 的 CategoryID 不是 1", s.Key)
				}
				t.Logf("category_id=1 配置数: %d", len(settings))
			},
		},
		{
			name:  "按 category_id=2 筛选（security）",
			query: map[string]string{"category_id": "2"},
			validate: func(t *testing.T, settings []setting.SettingDTO) {
				t.Helper()
				for _, s := range settings {
					assert.Equal(t, uint(2), s.CategoryID, "配置 %s 的 CategoryID 不是 2", s.Key)
				}
				t.Logf("category_id=2 配置数: %d", len(settings))
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			settings, err := helper.Get[[]setting.SettingDTO](c, "/api/admin/settings", tc.query)
			require.NoError(t, err, "获取配置失败")
			tc.validate(t, *settings)
		})
	}
}

// TestGetSettingsSchema 测试获取配置 Schema。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestGetSettingsSchema ./internal/manualtest/
func TestGetSettingsSchema(t *testing.T) {
	c := helper.LoginAsAdmin(t)

	t.Log("\n获取配置 Schema...")
	schema, err := helper.Get[[]setting.SchemaCategoryDTO](c, "/api/admin/settings/schema", nil)
	require.NoError(t, err, "获取 Schema 失败")

	t.Logf("Schema 层级结构:")
	for _, cat := range *schema {
		t.Logf("  📁 %s (%s) [icon: %s]", cat.Label, cat.Category, cat.Icon)
		for _, group := range cat.Groups {
			t.Logf("    📂 %s (%s)", group.Label, group.Group)
			for _, setting := range group.Settings {
				t.Logf("      - %s (%s): %v", setting.Label, setting.Key, setting.Value)
			}
		}
	}
}

// TestBatchUpdateSettings 测试批量更新配置。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestBatchUpdateSettings ./internal/manualtest/
func TestBatchUpdateSettings(t *testing.T) {
	c := helper.LoginAsAdmin(t)

	// 先创建两个测试配置
	timestamp := time.Now().Unix()
	key1 := fmt.Sprintf("%sbatch1_%d", settingTestPrefix, timestamp)
	key2 := fmt.Sprintf("%sbatch2_%d", settingTestPrefix, timestamp)

	t.Log("\n准备: 创建两个测试配置...")
	for _, key := range []string{key1, key2} {
		createReq := map[string]any{
			"key":           key,
			"default_value": "初始值",
			"category_id":   1, // 使用 general 分类（ID=1）
			"group":         "batch",
			"value_type":    "string",
			"label":         "批量测试",
		}
		_, createErr := helper.Post[setting.SettingDTO](c, "/api/admin/settings", createReq)
		require.NoError(t, createErr, "创建配置 %s 失败", key)
		t.Logf("  创建配置: %s", key)
	}

	// 确保清理
	t.Cleanup(func() {
		for _, key := range []string{key1, key2} {
			if deleteErr := c.Delete("/api/admin/settings/" + key); deleteErr != nil {
				t.Logf("清理配置 %s 失败: %v", key, deleteErr)
			}
		}
	})

	// 批量更新
	t.Log("\n测试: 批量更新配置...")
	batchReq := map[string]any{
		"settings": []map[string]any{
			{"key": key1, "value": "批量更新值1"},
			{"key": key2, "value": "批量更新值2"},
		},
	}

	resp, err := c.R().
		SetBody(batchReq).
		Post("/api/admin/settings/batch")
	require.NoError(t, err, "批量更新请求失败")
	require.False(t, resp.IsError(), "批量更新失败: 状态码 %d", resp.StatusCode())
	t.Log("  批量更新成功!")

	// 验证更新结果
	t.Log("\n验证更新结果...")
	for i, key := range []string{key1, key2} {
		detail, getErr := helper.Get[setting.SettingDTO](c, "/api/admin/settings/"+key, nil)
		require.NoError(t, getErr, "获取配置 %s 失败", key)
		expected := fmt.Sprintf("批量更新值%d", i+1)
		assert.Equal(t, expected, detail.DefaultValue, "配置 %s 值不匹配", key)
		t.Logf("  ✓ %s = %v", key, detail.DefaultValue)
	}
}

// TestDeleteSetting 测试删除配置。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestDeleteSetting ./internal/manualtest/
func TestDeleteSetting(t *testing.T) {
	c := helper.LoginAsAdmin(t)

	// 先创建一个测试配置
	settingKey := fmt.Sprintf("%sdelete_%d", settingTestPrefix, time.Now().Unix())
	t.Logf("\n准备: 创建测试配置 %s...", settingKey)
	createReq := map[string]any{
		"key":           settingKey,
		"default_value": "待删除",
		"category_id":   1, // 使用 general 分类（ID=1）
		"group":         "delete",
		"value_type":    "string",
		"label":         "删除测试",
	}
	_, err := helper.Post[setting.SettingDTO](c, "/api/admin/settings", createReq)
	require.NoError(t, err, "创建配置失败")
	t.Log("  创建成功")

	// 删除配置
	t.Log("\n测试: 删除配置...")
	err = c.Delete("/api/admin/settings/" + settingKey)
	require.NoError(t, err, "删除配置失败")
	t.Log("  删除成功!")

	// 验证删除
	t.Log("\n验证: 确认配置已删除...")
	_, err = helper.Get[setting.SettingDTO](c, "/api/admin/settings/"+settingKey, nil)
	require.Error(t, err, "配置应该已被删除，但仍能获取")
	t.Log("  ✓ 配置已成功删除")
}
