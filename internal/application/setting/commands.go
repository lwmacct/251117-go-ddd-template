package setting

// CreateCommand 创建设置命令
type CreateCommand struct {
	Key       string
	Value     string
	Category  string
	ValueType string
	Label     string
}

// UpdateCommand 更新设置命令
type UpdateCommand struct {
	Key       string
	Value     string
	ValueType string
	Label     string
}

// DeleteCommand 删除设置命令
type DeleteCommand struct {
	Key string
}

// SettingItemCommand 设置项
type SettingItemCommand struct {
	Key   string
	Value string
}

// BatchUpdateCommand 批量更新设置命令
type BatchUpdateCommand struct {
	Settings []SettingItemCommand
}
