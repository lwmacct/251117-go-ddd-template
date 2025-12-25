package setting

// SettingItemCommand 设置项
type SettingItemCommand struct {
	Key   string
	Value string
}

// BatchUpdateSettingsCommand 批量更新设置命令
type BatchUpdateSettingsCommand struct {
	Settings []SettingItemCommand
}
