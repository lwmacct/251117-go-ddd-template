package captcha

import "context"

// QueryRepository 定义只读操作
type QueryRepository interface {
	// GetStats 获取验证码存储统计信息
	GetStats(ctx context.Context) map[string]interface{}
}
