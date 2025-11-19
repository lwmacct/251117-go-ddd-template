package handler

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/captcha"
	captchaService "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/captcha"
)

// CaptchaHandler 验证码处理器
type CaptchaHandler struct {
	captchaRepo    captcha.Repository
	captchaService *captchaService.Service
	devSecret      string // 开发模式密钥
}

// NewCaptchaHandler 创建验证码处理器
func NewCaptchaHandler(
	captchaRepo captcha.Repository,
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
// GET /api/auth/captcha
// 开发模式：?code=9999&secret=your-secret
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
