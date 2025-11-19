package captcha

import (
	"fmt"
	"image/color"
	"time"

	"github.com/mojocn/base64Captcha"
)

const (
	// DefaultExpiration 默认过期时间（5分钟）
	DefaultExpiration = 5 * time.Minute
	// DefaultLength 验证码长度
	DefaultLength = 4
	// DefaultWidth 验证码图片宽度
	DefaultWidth = 140
	// DefaultHeight 验证码图片高度
	DefaultHeight = 50
)

// 验证码字符集（a-z, 0-9）
const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

// Service 验证码服务
// 负责验证码图片的生成，使用 base64Captcha 库
type Service struct {
	driver *base64Captcha.DriverString
}

// NewService 创建验证码服务
func NewService() *Service {
	driver := base64Captcha.NewDriverString(
		DefaultHeight,                      // 高度
		DefaultWidth,                       // 宽度
		2,                                  // 干扰噪点数量（适中）
		base64Captcha.OptionShowSlimeLine,  // 显示干扰线选项
		DefaultLength,                      // 验证码长度
		charset,                            // 字符集
		&color.RGBA{245, 245, 245, 255},    // 背景颜色（浅灰）
		nil,                                // 使用默认字体
		[]string{"wqy-microhei.ttc"},       // 字体列表
	).ConvertFonts()

	return &Service{
		driver: driver,
	}
}

// GenerateRandomCode 生成随机验证码
// 返回 (captchaID, imageBase64, code, error)
func (s *Service) GenerateRandomCode() (string, string, string, error) {
	captchaInstance := base64Captcha.NewCaptcha(s.driver, base64Captcha.DefaultMemStore)
	captchaID, b64s, _, err := captchaInstance.Generate()
	if err != nil {
		return "", "", "", err
	}

	// 从 base64Captcha 的 store 获取验证码值
	code := base64Captcha.DefaultMemStore.Get(captchaID, false) // false 表示不删除

	return captchaID, b64s, code, nil
}

// GenerateCustomCodeImage 生成指定文本的验证码图片
// 用于开发模式
func (s *Service) GenerateCustomCodeImage(text string) (string, error) {
	item, err := s.driver.DrawCaptcha(text)
	if err != nil {
		return "", err
	}
	return item.EncodeB64string(), nil
}

// GenerateDevCaptchaID 生成开发模式验证码ID
func GenerateDevCaptchaID() string {
	return fmt.Sprintf("dev-%d", time.Now().Unix())
}
