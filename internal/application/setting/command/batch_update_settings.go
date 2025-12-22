// Package command 定义设置命令处理器
package command

// SettingItem 设置项
type SettingItem struct {
	Key   string
	Value string
}

// BatchUpdateSettingsCommand 批量更新设置命令
type BatchUpdateSettingsCommand struct {
	Settings []SettingItem
}
