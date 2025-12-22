package setting

import "errors"

var (
	// ErrSettingNotFound 设置不存在
	ErrSettingNotFound = errors.New("setting not found")

	// ErrSettingKeyExists 设置键已存在
	ErrSettingKeyExists = errors.New("setting key already exists")

	// ErrInvalidValueType 无效的值类型
	ErrInvalidValueType = errors.New("invalid value type")

	// ErrInvalidValue 无效的设置值
	ErrInvalidValue = errors.New("invalid setting value")

	// ErrCategoryNotFound 设置分类不存在
	ErrCategoryNotFound = errors.New("category not found")

	// ErrValueTypeMismatch 值类型不匹配
	ErrValueTypeMismatch = errors.New("value type mismatch")

	// ErrInvalidJSONValue 无效的 JSON 值
	ErrInvalidJSONValue = errors.New("invalid JSON value")

	// ErrInvalidBoolValue 无效的布尔值
	ErrInvalidBoolValue = errors.New("invalid boolean value")

	// ErrInvalidNumberValue 无效的数值
	ErrInvalidNumberValue = errors.New("invalid number value")
)
