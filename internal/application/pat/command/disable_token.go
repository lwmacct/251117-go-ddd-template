// Package command 定义 PAT 命令处理器
package command

// DisableTokenCommand 禁用 Token 命令
type DisableTokenCommand struct {
	UserID  uint
	TokenID uint
}
