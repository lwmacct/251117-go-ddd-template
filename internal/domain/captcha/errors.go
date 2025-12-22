package captcha

import "errors"

var (
	// ErrCaptchaNotFound 验证码不存在
	ErrCaptchaNotFound = errors.New("captcha not found")

	// ErrCaptchaExpired 验证码已过期
	ErrCaptchaExpired = errors.New("captcha has expired")

	// ErrInvalidCaptcha 验证码错误
	ErrInvalidCaptcha = errors.New("invalid captcha")

	// ErrCaptchaAlreadyUsed 验证码已被使用
	ErrCaptchaAlreadyUsed = errors.New("captcha already used")

	// ErrTooManyAttempts 尝试次数过多
	ErrTooManyAttempts = errors.New("too many captcha attempts")

	// ErrCaptchaGenerationFailed 验证码生成失败
	ErrCaptchaGenerationFailed = errors.New("captcha generation failed")
)
