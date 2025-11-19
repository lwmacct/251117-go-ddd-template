package captcha

import "time"

// CaptchaData 验证码数据实体
// 用于存储验证码的值和过期时间（内存存储，无 GORM 标签）
type CaptchaData struct {
	Code      string    // 验证码值（小写，不区分大小写）
	ExpireAt  time.Time // 过期时间
	CreatedAt time.Time // 创建时间
}

// IsExpired 检查验证码是否过期
func (c *CaptchaData) IsExpired() bool {
	return time.Now().After(c.ExpireAt)
}
