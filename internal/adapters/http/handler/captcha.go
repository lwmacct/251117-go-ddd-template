package handler

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/captcha"
	captchaService "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/captcha"
)

// CaptchaHandler 验证码处理器
type CaptchaHandler struct {
	captchaRepo    captcha.CommandRepository
	captchaService *captchaService.Service
	devSecret      string // 开发模式密钥
}

// NewCaptchaHandler 创建验证码处理器
func NewCaptchaHandler(
	captchaRepo captcha.CommandRepository,
	captchaService *captchaService.Service,
	devSecret string,
) *CaptchaHandler {
	return &CaptchaHandler{
		captchaRepo:    captchaRepo,
		captchaService: captchaService,
		devSecret:      devSecret,
	}
}

// GenerateCaptchaResponse 生成验证码响应
type GenerateCaptchaResponse struct {
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
// @Success      200 {object} response.Response{data=GenerateCaptchaResponse} "验证码生成成功"
// @Failure      500 {object} response.ErrorResponse "生成失败"
// @Router       /api/auth/captcha [get]
func (h *CaptchaHandler) GetCaptcha(c *gin.Context) {
	ctx := c.Request.Context()

	// 开发模式参数
	code := c.Query("code")
	secret := c.Query("secret")

	var captchaID, imageBase64, codeValue string
	var err error

	// 检查是否为开发模式
	isDevMode := code != "" && secret != "" && secret == h.devSecret

	if isDevMode {
		// 开发模式：使用自定义验证码
		codeValue = code
		captchaID = captchaService.GenerateDevCaptchaID()
		imageBase64, err = h.captchaService.GenerateCustomCodeImage(codeValue)
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to generate captcha image"})
			return
		}
	} else {
		// 普通模式：生成随机验证码
		captchaID, imageBase64, codeValue, err = h.captchaService.GenerateRandomCode()
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to generate captcha"})
			return
		}
	}

	// 存储验证码到 Repository
	expiration := captchaService.DefaultExpiration
	if err := h.captchaRepo.Create(ctx, captchaID, codeValue, expiration); err != nil {
		c.JSON(500, gin.H{"error": "failed to save captcha"})
		return
	}

	// 构建返回结果
	result := &GenerateCaptchaResponse{
		ID:       captchaID,
		Image:    imageBase64,
		ExpireAt: time.Now().Add(expiration).Unix(),
	}

	// 开发模式下返回验证码值
	if isDevMode {
		result.Code = codeValue
	}

	c.JSON(200, gin.H{
		"message": "success",
		"data":    result,
	})
}
