package auth

// LoginCommand 登录命令
type LoginCommand struct {
	Account   string // 用户名或邮箱
	Password  string
	CaptchaID string
	Captcha   string
	ClientIP  string // 客户端 IP（用于审计日志）
	UserAgent string // 用户代理（用于审计日志）
}

// Login2FACommand 二次认证命令
type Login2FACommand struct {
	SessionToken  string
	TwoFactorCode string
	ClientIP      string // 客户端 IP（用于审计日志）
	UserAgent     string // 用户代理（用于审计日志）
}
