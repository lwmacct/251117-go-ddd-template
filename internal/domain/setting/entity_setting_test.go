package setting

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetting_IsValidValueType(t *testing.T) {
	tests := []struct {
		name      string
		valueType string
		want      bool
	}{
		{"string type", ValueTypeString, true},
		{"number type", ValueTypeNumber, true},
		{"boolean type", ValueTypeBoolean, true},
		{"json type", ValueTypeJSON, true},
		{"invalid type", "invalid", false},
		{"empty type", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Setting{ValueType: tt.valueType}
			assert.Equal(t, tt.want, s.IsValidValueType())
		})
	}
}

func TestSetting_IsValidCategory(t *testing.T) {
	tests := []struct {
		name     string
		category string
		want     bool
	}{
		{"general", CategoryGeneral, true},
		{"security", CategorySecurity, true},
		{"notification", CategoryNotification, true},
		{"backup", CategoryBackup, true},
		{"invalid", "invalid", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Setting{Category: tt.category}
			assert.Equal(t, tt.want, s.IsValidCategory())
		})
	}
}

func TestSetting_ParseBool(t *testing.T) {
	tests := []struct {
		name      string
		valueType string
		value     string
		want      bool
		wantErr   error
	}{
		{"true", ValueTypeBoolean, "true", true, nil},
		{"false", ValueTypeBoolean, "false", false, nil},
		{"1", ValueTypeBoolean, "1", true, nil},
		{"0", ValueTypeBoolean, "0", false, nil},
		{"yes", ValueTypeBoolean, "yes", true, nil},
		{"no", ValueTypeBoolean, "no", false, nil},
		{"on", ValueTypeBoolean, "on", true, nil},
		{"off", ValueTypeBoolean, "off", false, nil},
		{"TRUE uppercase", ValueTypeBoolean, "TRUE", true, nil},
		{"invalid value", ValueTypeBoolean, "invalid", false, ErrInvalidBoolValue},
		{"wrong type", ValueTypeString, "true", false, ErrValueTypeMismatch},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Setting{ValueType: tt.valueType, Value: tt.value}
			got, err := s.ParseBool()
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestSetting_ParseInt(t *testing.T) {
	tests := []struct {
		name      string
		valueType string
		value     string
		want      int
		wantErr   error
	}{
		{"positive", ValueTypeNumber, "42", 42, nil},
		{"zero", ValueTypeNumber, "0", 0, nil},
		{"negative", ValueTypeNumber, "-10", -10, nil},
		{"invalid number", ValueTypeNumber, "abc", 0, ErrInvalidNumberValue},
		{"float value", ValueTypeNumber, "3.14", 0, ErrInvalidNumberValue},
		{"wrong type", ValueTypeString, "42", 0, ErrValueTypeMismatch},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Setting{ValueType: tt.valueType, Value: tt.value}
			got, err := s.ParseInt()
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestSetting_ParseFloat(t *testing.T) {
	tests := []struct {
		name      string
		valueType string
		value     string
		want      float64
		wantErr   error
	}{
		{"integer", ValueTypeNumber, "42", 42.0, nil},
		{"float", ValueTypeNumber, "3.14", 3.14, nil},
		{"negative", ValueTypeNumber, "-2.5", -2.5, nil},
		{"invalid", ValueTypeNumber, "abc", 0, ErrInvalidNumberValue},
		{"wrong type", ValueTypeString, "3.14", 0, ErrValueTypeMismatch},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Setting{ValueType: tt.valueType, Value: tt.value}
			got, err := s.ParseFloat()
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				assert.InDelta(t, tt.want, got, 0.001)
			}
		})
	}
}

func TestSetting_ParseJSON(t *testing.T) {
	t.Run("valid json object", func(t *testing.T) {
		s := &Setting{ValueType: ValueTypeJSON, Value: `{"name":"test","count":10}`}
		var result map[string]any
		err := s.ParseJSON(&result)
		require.NoError(t, err)
		assert.Equal(t, "test", result["name"])
		assert.InEpsilon(t, float64(10), result["count"], 0.001)
	})

	t.Run("valid json array", func(t *testing.T) {
		s := &Setting{ValueType: ValueTypeJSON, Value: `[1,2,3]`}
		var result []int
		err := s.ParseJSON(&result)
		require.NoError(t, err)
		assert.Equal(t, []int{1, 2, 3}, result)
	})

	t.Run("invalid json", func(t *testing.T) {
		s := &Setting{ValueType: ValueTypeJSON, Value: `{invalid}`}
		var result map[string]any
		err := s.ParseJSON(&result)
		assert.ErrorIs(t, err, ErrInvalidJSONValue)
	})

	t.Run("wrong type", func(t *testing.T) {
		s := &Setting{ValueType: ValueTypeString, Value: `{"key":"value"}`}
		var result map[string]any
		err := s.ParseJSON(&result)
		assert.ErrorIs(t, err, ErrValueTypeMismatch)
	})
}

func TestSetting_SetBool(t *testing.T) {
	t.Run("set true", func(t *testing.T) {
		s := &Setting{}
		s.SetBool(true)
		assert.Equal(t, ValueTypeBoolean, s.ValueType)
		assert.Equal(t, "true", s.Value)
	})

	t.Run("set false", func(t *testing.T) {
		s := &Setting{}
		s.SetBool(false)
		assert.Equal(t, ValueTypeBoolean, s.ValueType)
		assert.Equal(t, "false", s.Value)
	})
}

func TestSetting_SetInt(t *testing.T) {
	tests := []struct {
		name  string
		value int
		want  string
	}{
		{"positive", 42, "42"},
		{"zero", 0, "0"},
		{"negative", -10, "-10"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Setting{}
			s.SetInt(tt.value)
			assert.Equal(t, ValueTypeNumber, s.ValueType)
			assert.Equal(t, tt.want, s.Value)
		})
	}
}

func TestSetting_SetJSON(t *testing.T) {
	t.Run("set object", func(t *testing.T) {
		s := &Setting{}
		err := s.SetJSON(map[string]string{"key": "value"})
		require.NoError(t, err)
		assert.Equal(t, ValueTypeJSON, s.ValueType) //nolint:testifylint // ValueTypeJSON 是类型常量，不是 JSON 内容
		assert.JSONEq(t, `{"key":"value"}`, s.Value)
	})

	t.Run("set array", func(t *testing.T) {
		s := &Setting{}
		err := s.SetJSON([]int{1, 2, 3})
		require.NoError(t, err)
		assert.Equal(t, ValueTypeJSON, s.ValueType) //nolint:testifylint // ValueTypeJSON 是类型常量，不是 JSON 内容
		assert.Equal(t, "[1,2,3]", s.Value)
	})
}

func TestSetting_IsEmpty(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"empty string", "", true},
		{"non-empty string", "value", false},
		{"whitespace only", " ", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Setting{Value: tt.value}
			assert.Equal(t, tt.want, s.IsEmpty())
		})
	}
}
