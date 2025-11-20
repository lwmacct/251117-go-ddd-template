// Package command 定义 PAT 命令处理器
package command

// EnableTokenCommand 启用 Token 命令
type EnableTokenCommand struct {
	UserID  uint
	TokenID uint
}
