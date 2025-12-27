package validation_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/validation"
)

func TestJSONLogicValidator_Validate_SimpleRules(t *testing.T) {
	validator := validation.NewJSONLogicValidator()
	ctx := context.Background()

	tests := []struct {
		name      string
		rule      string
		value     any
		wantValid bool
	}{
		{
			name:      "简单最小值规则 - 通过",
			rule:      `{"min":6}`,
			value:     10.0,
			wantValid: true,
		},
		{
			name:      "简单最小值规则 - 失败",
			rule:      `{"min":6}`,
			value:     3.0,
			wantValid: false,
		},
		{
			name:      "简单最大值规则 - 通过",
			rule:      `{"max":100}`,
			value:     50.0,
			wantValid: true,
		},
		{
			name:      "简单最大值规则 - 失败",
			rule:      `{"max":100}`,
			value:     150.0,
			wantValid: false,
		},
		{
			name:      "简单范围规则 - 通过",
			rule:      `{"min":6,"max":32}`,
			value:     8.0,
			wantValid: true,
		},
		{
			name:      "简单范围规则 - 失败（太小）",
			rule:      `{"min":6,"max":32}`,
			value:     3.0,
			wantValid: false,
		},
		{
			name:      "简单范围规则 - 失败（太大）",
			rule:      `{"min":6,"max":32}`,
			value:     100.0,
			wantValid: false,
		},
		{
			name:      "必填规则 - 通过",
			rule:      `{"required":true}`,
			value:     "hello",
			wantValid: true,
		},
		{
			name:      "必填规则 - 失败（空字符串）",
			rule:      `{"required":true}`,
			value:     "",
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vctx := &setting.ValidationContext{
				Key:   "test.key",
				Value: tt.value,
				Rule:  tt.rule,
			}
			result, err := validator.Validate(ctx, vctx)
			require.NoError(t, err)
			assert.Equal(t, tt.wantValid, result.Valid, "expected valid=%v, got %v, message=%s", tt.wantValid, result.Valid, result.Message)
		})
	}
}

func TestJSONLogicValidator_Validate_JSONLogicRules(t *testing.T) {
	validator := validation.NewJSONLogicValidator()
	ctx := context.Background()

	tests := []struct {
		name      string
		rule      string
		value     any
		wantValid bool
	}{
		{
			name:      "JSON Logic 大于等于 - 通过",
			rule:      `{">=": [{"var": "value"}, 6]}`,
			value:     10.0,
			wantValid: true,
		},
		{
			name:      "JSON Logic 大于等于 - 失败",
			rule:      `{">=": [{"var": "value"}, 6]}`,
			value:     3.0,
			wantValid: false,
		},
		{
			name:      "JSON Logic AND 组合 - 通过",
			rule:      `{"and": [{">=": [{"var": "value"}, 6]}, {"<=": [{"var": "value"}, 32]}]}`,
			value:     8.0,
			wantValid: true,
		},
		{
			name:      "JSON Logic AND 组合 - 失败",
			rule:      `{"and": [{">=": [{"var": "value"}, 6]}, {"<=": [{"var": "value"}, 32]}]}`,
			value:     100.0,
			wantValid: false,
		},
		{
			name:      "JSON Logic 真值检查 - 通过",
			rule:      `{"!!": {"var": "value"}}`,
			value:     "hello",
			wantValid: true,
		},
		{
			name:      "JSON Logic 真值检查 - 失败",
			rule:      `{"!!": {"var": "value"}}`,
			value:     "",
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vctx := &setting.ValidationContext{
				Key:   "test.key",
				Value: tt.value,
				Rule:  tt.rule,
			}
			result, err := validator.Validate(ctx, vctx)
			require.NoError(t, err)
			assert.Equal(t, tt.wantValid, result.Valid, "expected valid=%v, got %v, message=%s", tt.wantValid, result.Valid, result.Message)
		})
	}
}

func TestJSONLogicValidator_Validate_CrossFieldValidation(t *testing.T) {
	validator := validation.NewJSONLogicValidator()
	ctx := context.Background()

	// 备份频率必须小于 (保留天数 * 24)
	// 注意：JSON Logic 的 var 使用点符号访问嵌套对象
	rule := `{"<": [{"var": "value"}, {"*": [{"var": "settings.backup_retention_days"}, 24]}]}`

	tests := []struct {
		name          string
		value         any
		retentionDays any
		wantValid     bool
	}{
		{
			name:          "备份频率 24h，保留 30 天 - 通过",
			value:         24.0,
			retentionDays: 30.0,
			wantValid:     true,
		},
		{
			name:          "备份频率 168h，保留 7 天 - 失败（168 >= 7*24=168）",
			value:         168.0,
			retentionDays: 7.0,
			wantValid:     false,
		},
		{
			name:          "备份频率 1h，保留 1 天 - 通过（1 < 24）",
			value:         1.0,
			retentionDays: 1.0,
			wantValid:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vctx := &setting.ValidationContext{
				Key:   "backup.backup_frequency",
				Value: tt.value,
				Rule:  rule,
				AllSettings: map[string]any{
					// JSON Logic 需要嵌套结构：settings.backup_retention_days
					"backup_retention_days": tt.retentionDays,
				},
			}
			result, err := validator.Validate(ctx, vctx)
			require.NoError(t, err)
			assert.Equal(t, tt.wantValid, result.Valid, "expected valid=%v, got %v, message=%s", tt.wantValid, result.Valid, result.Message)
		})
	}
}

func TestJSONLogicValidator_ValidateBatch(t *testing.T) {
	validator := validation.NewJSONLogicValidator()
	ctx := context.Background()

	items := []*setting.ValidationContext{
		{
			Key:   "password_length",
			Value: 4.0, // 小于 6，应该失败
			Rule:  `{"min":6,"max":32}`,
		},
		{
			Key:   "session_timeout",
			Value: 30.0, // 在范围内，应该通过
			Rule:  `{"min":5,"max":1440}`,
		},
		{
			Key:   "max_attempts",
			Value: 15.0, // 大于 10，应该失败
			Rule:  `{"min":3,"max":10}`,
		},
	}

	errors, err := validator.ValidateBatch(ctx, items)
	require.NoError(t, err)

	// 应该有 2 个错误
	assert.Len(t, errors, 2)
	assert.Contains(t, errors, "password_length")
	assert.Contains(t, errors, "max_attempts")
	assert.NotContains(t, errors, "session_timeout")
}

func TestJSONLogicValidator_Validate_InvalidRule(t *testing.T) {
	validator := validation.NewJSONLogicValidator()
	ctx := context.Background()

	vctx := &setting.ValidationContext{
		Key:   "test.key",
		Value: 10.0,
		Rule:  `{invalid json`,
	}

	_, err := validator.Validate(ctx, vctx)
	assert.Error(t, err)
}

func TestJSONLogicValidator_Validate_EmptyRule(t *testing.T) {
	validator := validation.NewJSONLogicValidator()
	ctx := context.Background()

	vctx := &setting.ValidationContext{
		Key:   "test.key",
		Value: 10.0,
		Rule:  "",
	}

	result, err := validator.Validate(ctx, vctx)
	require.NoError(t, err)
	assert.True(t, result.Valid, "empty rule should pass validation")
}
