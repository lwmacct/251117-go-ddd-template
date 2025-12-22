package command

import (
	"context"
	"time"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/captcha"
)

// GenerateCaptchaHandler 生成验证码处理器
type GenerateCaptchaHandler struct {
	captchaCommandRepo captcha.CommandRepository
	captchaService     captcha.Service
}

// NewGenerateCaptchaHandler 创建 GenerateCaptchaHandler 实例
func NewGenerateCaptchaHandler(
	captchaCommandRepo captcha.CommandRepository,
	captchaService captcha.Service,
) *GenerateCaptchaHandler {
	return &GenerateCaptchaHandler{
		captchaCommandRepo: captchaCommandRepo,
		captchaService:     captchaService,
	}
}

// Handle 处理生成验证码命令
func (h *GenerateCaptchaHandler) Handle(ctx context.Context, cmd GenerateCaptchaCommand) (*GenerateCaptchaResult, error) {
	// 根据模式生成验证码
	captchaID, imageBase64, codeValue, err := h.generateCaptcha(cmd)
	if err != nil {
		return nil, err
	}

	// 存储验证码到 Repository
	expiration := time.Duration(h.captchaService.GetDefaultExpiration()) * time.Second
	if err := h.captchaCommandRepo.Create(ctx, captchaID, codeValue, expiration); err != nil {
		return nil, err
	}

	// 构建返回结果
	result := &GenerateCaptchaResult{
		ID:       captchaID,
		Image:    imageBase64,
		ExpireAt: time.Now().Add(expiration).Unix(),
	}

	// 开发模式下返回验证码值
	if cmd.DevMode {
		result.Code = codeValue
	}

	return result, nil
}

// generateCaptcha 根据命令生成验证码
func (h *GenerateCaptchaHandler) generateCaptcha(cmd GenerateCaptchaCommand) (string, string, string, error) {
	if cmd.DevMode && cmd.CustomCode != "" {
		// 开发模式：使用自定义验证码
		captchaID := h.captchaService.GenerateDevCaptchaID()
		imageBase64, err := h.captchaService.GenerateCustomCodeImage(cmd.CustomCode)
		if err != nil {
			return "", "", "", err
		}
		return captchaID, imageBase64, cmd.CustomCode, nil
	}

	// 普通模式：生成随机验证码
	return h.captchaService.GenerateRandomCode()
}
