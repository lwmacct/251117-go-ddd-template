package captcha

// CaptchaDTO 验证码响应 DTO
type CaptchaDTO struct {
	ID       string `json:"id"`             // 验证码ID
	Image    string `json:"image"`          // Base64编码的图片
	ExpireAt int64  `json:"expire_at"`      // 过期时间戳（秒）
	Code     string `json:"code,omitempty"` // 验证码值（仅开发模式返回）
}

// GenerateCaptchaResultDTO 验证码生成结果
type GenerateCaptchaResultDTO struct {
	// ID 验证码ID
	ID string `json:"id"`
	// Image Base64编码的图片
	Image string `json:"image"`
	// ExpireAt 过期时间戳（秒）
	ExpireAt int64 `json:"expire_at"`
	// Code 验证码值（仅开发模式返回）
	Code string `json:"code,omitempty"`
}
