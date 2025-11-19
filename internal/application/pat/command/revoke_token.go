// Package command 定义 PAT 命令处理器
package command

// RevokeTokenCommand 撤销 Token 命令
type RevokeTokenCommand struct {
	UserID  uint
	TokenID uint
}
