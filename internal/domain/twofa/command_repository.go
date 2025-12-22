package twofa

import "context"

// CommandRepository 定义 2FA 写操作接口
type CommandRepository interface {
	// CreateOrUpdate 创建或更新 2FA 配置
	CreateOrUpdate(ctx context.Context, twoFA *TwoFA) error

	// Delete 删除 2FA 配置
	Delete(ctx context.Context, userID uint) error
}
