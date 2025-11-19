// Package command 定义设置命令处理器
package command

// UpdateSettingCommand 更新设置命令
type UpdateSettingCommand struct {
	Key       string
	Value     string
	ValueType string
	Label     string
}
