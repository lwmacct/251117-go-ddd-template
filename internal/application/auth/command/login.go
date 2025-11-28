// Package command 定义认证模块的命令。
//
// 本包处理用户认证相关的所有命令：
//   - LoginCommand: 用户登录（支持用户名/邮箱 + 验证码）
//   - Login2FACommand: 双因素认证验证
//   - RegisterCommand: 用户注册
//   - RefreshTokenCommand: 刷新访问令牌
//
// 认证流程：
//  1. 客户端提交 LoginCommand
//  2. 若用户启用 2FA，返回临时 SessionToken
//  3. 客户端提交 Login2FACommand 完成二次验证
//  4. 返回 AccessToken 和 RefreshToken
package command

// LoginCommand 登录命令
type LoginCommand struct {
	Account   string // 用户名或邮箱
	Password  string
	CaptchaID string
	Captcha   string
}

// Login2FACommand 二次认证命令
type Login2FACommand struct {
	SessionToken  string
	TwoFactorCode string
}
