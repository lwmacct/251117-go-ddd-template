package setting

import "errors"

var (
	// ErrDefinitionNotFound 配置定义不存在
	ErrDefinitionNotFound = errors.New("setting definition not found")

	// ErrDefinitionKeyExists 配置定义键已存在
	ErrDefinitionKeyExists = errors.New("setting definition key already exists")

	// ErrUserSettingNotFound 用户配置不存在
	ErrUserSettingNotFound = errors.New("user setting not found")

	// ErrInvalidValueType 无效的值类型
	ErrInvalidValueType = errors.New("invalid value type")

	// ErrInvalidValue 无效的配置值
	ErrInvalidValue = errors.New("invalid setting value")

	// ErrCategoryNotFound 配置分类不存在
	ErrCategoryNotFound = errors.New("category not found")

	// ErrValidationFailed 验证失败
	ErrValidationFailed = errors.New("validation failed")

	// ErrInvalidValidationRule 无效的验证规则
	ErrInvalidValidationRule = errors.New("invalid validation rule")
)
