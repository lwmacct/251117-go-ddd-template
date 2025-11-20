package pat

import "context"

// CommandRepository 定义 PAT 写操作接口
type CommandRepository interface {
	// Create 创建新的个人访问令牌
	Create(ctx context.Context, pat *PersonalAccessToken) error

	// Update 更新令牌（主要用于 LastUsedAt）
	Update(ctx context.Context, pat *PersonalAccessToken) error

	// Delete 硬删除令牌
	Delete(ctx context.Context, id uint) error

	// Disable 禁用令牌（设置状态为 disabled）
	Disable(ctx context.Context, id uint) error

	// Enable 启用令牌（设置状态为 active）
	Enable(ctx context.Context, id uint) error

	// DeleteByUserID 删除指定用户的所有令牌
	DeleteByUserID(ctx context.Context, userID uint) error

	// CleanupExpired 清理过期令牌
	CleanupExpired(ctx context.Context) error
}
