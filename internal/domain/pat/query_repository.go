package pat

import "context"

// QueryRepository 定义 PAT 读操作接口
type QueryRepository interface {
	// FindByToken 通过令牌哈希查找（用于认证）
	FindByToken(ctx context.Context, tokenHash string) (*PersonalAccessToken, error)

	// FindByID 通过 ID 查找令牌
	FindByID(ctx context.Context, id uint) (*PersonalAccessToken, error)

	// FindByPrefix 通过前缀查找令牌
	FindByPrefix(ctx context.Context, prefix string) (*PersonalAccessToken, error)

	// ListByUser 获取指定用户的所有令牌
	ListByUser(ctx context.Context, userID uint) ([]*PersonalAccessToken, error)
}
