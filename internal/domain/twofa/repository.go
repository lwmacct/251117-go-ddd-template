package twofa

import "context"

// Repository 2FA 仓储接口
type Repository interface {
	// CreateOrUpdate 创建或更新 2FA 配置
	CreateOrUpdate(ctx context.Context, twoFA *TwoFA) error

	// FindByUserID 根据用户ID查找 2FA 配置
	FindByUserID(ctx context.Context, userID uint) (*TwoFA, error)

	// Delete 删除 2FA 配置
	Delete(ctx context.Context, userID uint) error

	// IsEnabled 检查用户是否启用了 2FA
	IsEnabled(ctx context.Context, userID uint) (bool, error)
}
