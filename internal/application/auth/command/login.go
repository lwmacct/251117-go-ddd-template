// Package command 定义认证模块的命令
package command

// LoginCommand 登录命令
type LoginCommand struct {
	Login    string // 用户名或邮箱
	Password string
	CaptchaID string
	Captcha   string
}

// Login2FACommand 二次认证命令
type Login2FACommand struct {
	SessionToken  string
	TwoFactorCode string
}
