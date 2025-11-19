// Package persistence 提供 PAT 查询仓储的 GORM 实现
package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/pat"
	"gorm.io/gorm"
)

// patQueryRepository PAT 查询仓储的 GORM 实现
type patQueryRepository struct {
	db *gorm.DB
}

// NewPATQueryRepository 创建 PAT 查询仓储实例
func NewPATQueryRepository(db *gorm.DB) pat.QueryRepository {
	return &patQueryRepository{db: db}
}

// FindByToken 通过令牌哈希查找（用于认证）
func (r *patQueryRepository) FindByToken(ctx context.Context, tokenHash string) (*pat.PersonalAccessToken, error) {
	var token pat.PersonalAccessToken
	err := r.db.WithContext(ctx).
		Where("token = ?", tokenHash).
		First(&token).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("token not found")
		}
		return nil, fmt.Errorf("failed to find PAT by token: %w", err)
	}

	return &token, nil
}

// FindByID 通过 ID 查找令牌
func (r *patQueryRepository) FindByID(ctx context.Context, id uint) (*pat.PersonalAccessToken, error) {
	var token pat.PersonalAccessToken
	err := r.db.WithContext(ctx).
		First(&token, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("token not found")
		}
		return nil, fmt.Errorf("failed to find PAT by ID: %w", err)
	}

	return &token, nil
}

// FindByPrefix 通过前缀查找令牌
func (r *patQueryRepository) FindByPrefix(ctx context.Context, prefix string) (*pat.PersonalAccessToken, error) {
	var token pat.PersonalAccessToken
	err := r.db.WithContext(ctx).
		Where("token_prefix = ?", prefix).
		First(&token).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("token not found")
		}
		return nil, fmt.Errorf("failed to find PAT by prefix: %w", err)
	}

	return &token, nil
}

// ListByUser 获取指定用户的所有令牌
func (r *patQueryRepository) ListByUser(ctx context.Context, userID uint) ([]*pat.PersonalAccessToken, error) {
	var tokens []*pat.PersonalAccessToken
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&tokens).Error

	if err != nil {
		return nil, fmt.Errorf("failed to list PATs by user: %w", err)
	}

	return tokens, nil
}
