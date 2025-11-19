package persistence

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/twofa"
)

// twofaRepository 2FA 仓储实现
type twofaRepository struct {
	db *gorm.DB
}

// NewTwoFARepository 创建 2FA 仓储（deprecated，使用 NewTwoFACommandRepository/NewTwoFAQueryRepository）
func NewTwoFARepository(db *gorm.DB) twofa.Repository {
	return &twofaRepository{db: db}
}

// NewTwoFACommandRepository 创建 2FA 写操作仓储
func NewTwoFACommandRepository(db *gorm.DB) twofa.CommandRepository {
	return &twofaRepository{db: db}
}

// NewTwoFAQueryRepository 创建 2FA 读操作仓储
func NewTwoFAQueryRepository(db *gorm.DB) twofa.QueryRepository {
	return &twofaRepository{db: db}
}

// CreateOrUpdate 创建或更新 2FA 配置
func (r *twofaRepository) CreateOrUpdate(ctx context.Context, tfa *twofa.TwoFA) error {
	// 使用 Upsert：如果存在则更新，不存在则创建
	result := r.db.WithContext(ctx).
		Where("user_id = ?", tfa.UserID).
		Assign(tfa).
		FirstOrCreate(tfa)

	if result.Error != nil {
		return fmt.Errorf("failed to create or update 2FA: %w", result.Error)
	}

	return nil
}

// FindByUserID 根据用户ID查找 2FA 配置
func (r *twofaRepository) FindByUserID(ctx context.Context, userID uint) (*twofa.TwoFA, error) {
	var tfa twofa.TwoFA
	result := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		First(&tfa)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find 2FA by user ID: %w", result.Error)
	}

	return &tfa, nil
}

// Delete 删除 2FA 配置
func (r *twofaRepository) Delete(ctx context.Context, userID uint) error {
	result := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&twofa.TwoFA{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete 2FA: %w", result.Error)
	}

	return nil
}

// IsEnabled 检查用户是否启用了 2FA
func (r *twofaRepository) IsEnabled(ctx context.Context, userID uint) (bool, error) {
	var count int64
	result := r.db.WithContext(ctx).
		Model(&twofa.TwoFA{}).
		Where("user_id = ? AND enabled = ?", userID, true).
		Count(&count)

	if result.Error != nil {
		return false, fmt.Errorf("failed to check 2FA status: %w", result.Error)
	}

	return count > 0, nil
}
