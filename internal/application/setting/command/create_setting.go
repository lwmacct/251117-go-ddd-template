// Package command 定义设置命令处理器
package command

// CreateSettingCommand 创建设置命令
type CreateSettingCommand struct {
	Key       string
	Value     string
	Category  string
	ValueType string
	Label     string
}

// CreateSettingResult 创建设置结果
type CreateSettingResult struct {
	ID uint
}
