// Package auth 定义认证模块的 DTO。
//
// 本包提供认证相关的数据传输对象：
//   - 请求 DTO：用于 HTTP 层接收客户端请求
//   - 响应 DTO：用于统一格式化响应数据
//
// 设计原则：
//   - DTO 包含验证标签（binding）和文档标签（example）
//   - 与 Command/Query 分离，支持独立演进
package auth

// LoginDTO 登录请求
type LoginDTO struct {
	Account   string `json:"account" binding:"required" example:"admin"`         // 手机号/用户名/邮箱
	Password  string `json:"password" binding:"required" example:"admin123"`     // 密码
	CaptchaID string `json:"captcha_id" binding:"required" example:"dev-123456"` // 验证码ID
	Captcha   string `json:"captcha" binding:"required" example:"9999"`          // 验证码
}

// Login2FADTO 二次认证请求
type Login2FADTO struct {
	SessionToken  string `json:"session_token" binding:"required" example:"eyJhbGciOiJIUzI1NiIs..."` // 登录时返回的临时会话令牌
	TwoFactorCode string `json:"two_factor_code" binding:"required,len=6" example:"123456"`          // 6位TOTP验证码
}

// RegisterDTO 注册请求
type RegisterDTO struct {
	Username string `json:"username" binding:"required,min=3,max=50" example:"john_doe"`
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"password123"`
	FullName string `json:"full_name" binding:"max=100" example:"John Doe"`
}

// RefreshTokenDTO 刷新令牌请求
type RefreshTokenDTO struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// TokenResponse 令牌响应
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	TokenResponse

	User *UserBriefResponse `json:"user,omitempty"`
}

// UserBriefResponse 用户简要信息响应
type UserBriefResponse struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
}

// TwoFARequiredResponse 需要二次认证响应
type TwoFARequiredResponse struct {
	Requires2FA  bool   `json:"requires_2fa"`
	SessionToken string `json:"session_token"`
}

// RegisterResponse 注册响应
type RegisterResponse struct {
	TokenResponse

	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
