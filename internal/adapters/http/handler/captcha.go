package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/response"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/captcha"
)

// CaptchaHandler 验证码处理器
type CaptchaHandler struct {
	generateCaptchaHandler *captcha.GenerateCaptchaHandler
	devSecret              string
}

// NewCaptchaHandler 创建验证码处理器
func NewCaptchaHandler(
	generateCaptchaHandler *captcha.GenerateCaptchaHandler,
	devSecret string,
) *CaptchaHandler {
	return &CaptchaHandler{
		generateCaptchaHandler: generateCaptchaHandler,
		devSecret:              devSecret,
	}
}

// GenerateCaptchaDTO 生成验证码响应
type GenerateCaptchaDTO struct {
	ID       string `json:"id"`             // 验证码ID
	Image    string `json:"image"`          // Base64编码的图片
	ExpireAt int64  `json:"expire_at"`      // 过期时间戳
	Code     string `json:"code,omitempty"` // 验证码值（开发模式下返回）
}

// GetCaptcha 获取验证码
//
// @Summary      获取图形验证码
// @Description  生成图形验证码用于登录。支持开发模式（通过code和secret参数指定验证码值）
// @Tags         认证 (Authentication)
// @Accept       json
// @Produce      json
// @Param        code query string false "开发模式：指定验证码值" example:"9999"
// @Param        secret query string false "开发模式：密钥" example:"dev-secret"
// @Success      200 {object} response.DataResponse[captcha.GenerateCaptchaResultDTO] "验证码生成成功"
// @Failure      500 {object} response.ErrorResponse "生成失败"
// @Router       /api/auth/captcha [get]
func (h *CaptchaHandler) GetCaptcha(c *gin.Context) {
	ctx := c.Request.Context()

	// 解析开发模式参数
	code := c.Query("code")
	secret := c.Query("secret")
	isDevMode := code != "" && secret != "" && secret == h.devSecret

	// 构建命令
	cmd := captcha.GenerateCaptchaCommand{
		DevMode:    isDevMode,
		CustomCode: code,
	}

	// 调用 Application Handler
	result, err := h.generateCaptchaHandler.Handle(ctx, cmd)
	if err != nil {
		response.InternalError(c, "failed to generate captcha", err.Error())
		return
	}

	response.OK(c, "success", result)
}
