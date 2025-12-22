package setting

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	domainSetting "github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

func TestToSettingResponse(t *testing.T) {
	t.Run("转换正常设置", func(t *testing.T) {
		now := time.Now()
		s := &domainSetting.Setting{
			ID:        1,
			Key:       "site.name",
			Value:     "My Application",
			Category:  "site",
			ValueType: "string",
			Label:     "站点名称",
			CreatedAt: now,
			UpdatedAt: now,
		}

		result := ToSettingResponse(s)

		assert.NotNil(t, result)
		assert.Equal(t, uint(1), result.ID)
		assert.Equal(t, "site.name", result.Key)
		assert.Equal(t, "My Application", result.Value)
		assert.Equal(t, "site", result.Category)
		assert.Equal(t, "string", result.ValueType)
		assert.Equal(t, "站点名称", result.Label)
		assert.Equal(t, now, result.CreatedAt)
		assert.Equal(t, now, result.UpdatedAt)
	})

	t.Run("转换nil返回nil", func(t *testing.T) {
		result := ToSettingResponse(nil)
		assert.Nil(t, result)
	})

	t.Run("转换布尔类型设置", func(t *testing.T) {
		s := &domainSetting.Setting{
			ID:        2,
			Key:       "feature.enabled",
			Value:     "true",
			Category:  "feature",
			ValueType: "boolean",
			Label:     "功能开关",
		}

		result := ToSettingResponse(s)

		assert.NotNil(t, result)
		assert.Equal(t, "feature.enabled", result.Key)
		assert.Equal(t, "true", result.Value)
		assert.Equal(t, "boolean", result.ValueType)
	})

	t.Run("转换数字类型设置", func(t *testing.T) {
		s := &domainSetting.Setting{
			ID:        3,
			Key:       "limit.max_users",
			Value:     "1000",
			Category:  "limit",
			ValueType: "number",
			Label:     "最大用户数",
		}

		result := ToSettingResponse(s)

		assert.NotNil(t, result)
		assert.Equal(t, "limit.max_users", result.Key)
		assert.Equal(t, "1000", result.Value)
		assert.Equal(t, "number", result.ValueType)
	})

	t.Run("转换JSON类型设置", func(t *testing.T) {
		s := &domainSetting.Setting{
			ID:        4,
			Key:       "smtp.config",
			Value:     `{"host":"smtp.example.com","port":587}`,
			Category:  "email",
			ValueType: "json",
			Label:     "SMTP配置",
		}

		result := ToSettingResponse(s)

		assert.NotNil(t, result)
		assert.Equal(t, "smtp.config", result.Key)
		assert.Contains(t, result.Value, "smtp.example.com")
		assert.Equal(t, "json", result.ValueType)
	})

	t.Run("转换空字段设置", func(t *testing.T) {
		s := &domainSetting.Setting{
			ID:  5,
			Key: "empty.setting",
		}

		result := ToSettingResponse(s)

		assert.NotNil(t, result)
		assert.Equal(t, uint(5), result.ID)
		assert.Equal(t, "empty.setting", result.Key)
		assert.Empty(t, result.Value)
		assert.Empty(t, result.Category)
		assert.Empty(t, result.ValueType)
		assert.Empty(t, result.Label)
	})
}
