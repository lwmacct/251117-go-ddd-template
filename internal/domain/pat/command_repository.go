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

	// Revoke 撤销令牌（设置状态为 revoked）
	Revoke(ctx context.Context, id uint) error

	// RevokeByUserID 撤销指定用户的所有令牌
	RevokeByUserID(ctx context.Context, userID uint) error

	// CleanupExpired 清理过期令牌
	CleanupExpired(ctx context.Context) error
}
