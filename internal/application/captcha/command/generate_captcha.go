// Package command 定义验证码生成命令
package command

// GenerateCaptchaCommand 生成验证码命令
type GenerateCaptchaCommand struct {
	// DevMode 是否为开发模式
	DevMode bool
	// CustomCode 自定义验证码（开发模式）
	CustomCode string
}

// GenerateCaptchaResult 验证码生成结果
type GenerateCaptchaResult struct {
	// ID 验证码ID
	ID string `json:"id"`
	// Image Base64编码的图片
	Image string `json:"image"`
	// ExpireAt 过期时间戳（秒）
	ExpireAt int64 `json:"expire_at"`
	// Code 验证码值（仅开发模式返回）
	Code string `json:"code,omitempty"`
}
