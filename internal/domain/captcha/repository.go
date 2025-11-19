package captcha

import (
	"context"
	"time"
)

// Repository 验证码仓储接口
type Repository interface {
	// Create 创建验证码并存储
	Create(ctx context.Context, captchaID string, code string, expiration time.Duration) error

	// Verify 验证验证码（不区分大小写，一次性使用）
	// 验证后会自动删除验证码
	Verify(ctx context.Context, captchaID string, code string) (bool, error)

	// Delete 删除验证码
	Delete(ctx context.Context, captchaID string) error

	// GetStats 获取统计信息
	GetStats(ctx context.Context) map[string]interface{}
}
