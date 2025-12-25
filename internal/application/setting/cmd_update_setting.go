package setting

// UpdateSettingCommand 更新设置命令
type UpdateSettingCommand struct {
	Key       string
	Value     string
	ValueType string
	Label     string
}
