// Package command 定义认证模块的命令
package command

// RegisterCommand 注册命令
type RegisterCommand struct {
	Username string
	Email    string
	Password string
	FullName string
}
