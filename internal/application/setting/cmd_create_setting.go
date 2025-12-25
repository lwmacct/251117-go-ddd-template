package setting

// CreateSettingCommand 创建设置命令
type CreateSettingCommand struct {
	Key       string
	Value     string
	Category  string
	ValueType string
	Label     string
}
