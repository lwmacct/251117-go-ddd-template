// Package command 定义 PAT 命令处理器
package command

// CreateTokenCommand 创建 Token 命令
type CreateTokenCommand struct {
	UserID      uint
	Name        string
	Permissions []string
	ExpiresAt   *string
}

// CreateTokenResult 创建 Token 结果
type CreateTokenResult struct {
	ID          uint
	Token       string
	Name        string
	Permissions []string
}
