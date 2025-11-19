// Package persistence 提供 PAT 命令仓储的 GORM 实现
package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/pat"
	"gorm.io/gorm"
)

// patCommandRepository PAT 命令仓储的 GORM 实现
type patCommandRepository struct {
	db *gorm.DB
}

// NewPATCommandRepository 创建 PAT 命令仓储实例
func NewPATCommandRepository(db *gorm.DB) pat.CommandRepository {
	return &patCommandRepository{db: db}
}

// Create 创建新的个人访问令牌
func (r *patCommandRepository) Create(ctx context.Context, token *pat.PersonalAccessToken) error {
	if err := r.db.WithContext(ctx).Create(token).Error; err != nil {
		return fmt.Errorf("failed to create PAT: %w", err)
	}
	return nil
}

// Update 更新令牌（主要用于 LastUsedAt）
func (r *patCommandRepository) Update(ctx context.Context, token *pat.PersonalAccessToken) error {
	if err := r.db.WithContext(ctx).Save(token).Error; err != nil {
		return fmt.Errorf("failed to update PAT: %w", err)
	}
	return nil
}

// Delete 硬删除令牌
func (r *patCommandRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).
		Unscoped(). // Hard delete
		Delete(&pat.PersonalAccessToken{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete PAT: %w", err)
	}
	return nil
}

// Revoke 撤销令牌（设置状态为 revoked）
func (r *patCommandRepository) Revoke(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).
		Model(&pat.PersonalAccessToken{}).
		Where("id = ?", id).
		Update("status", "revoked").Error; err != nil {
		return fmt.Errorf("failed to revoke PAT: %w", err)
	}
	return nil
}

// RevokeByUserID 撤销指定用户的所有令牌
func (r *patCommandRepository) RevokeByUserID(ctx context.Context, userID uint) error {
	if err := r.db.WithContext(ctx).
		Model(&pat.PersonalAccessToken{}).
		Where("user_id = ? AND status = ?", userID, "active").
		Update("status", "revoked").Error; err != nil {
		return fmt.Errorf("failed to revoke PATs by user ID: %w", err)
	}
	return nil
}

// CleanupExpired 清理过期令牌
func (r *patCommandRepository) CleanupExpired(ctx context.Context) error {
	now := time.Now()

	if err := r.db.WithContext(ctx).
		Model(&pat.PersonalAccessToken{}).
		Where("expires_at IS NOT NULL AND expires_at < ? AND status != ?", now, "expired").
		Update("status", "expired").Error; err != nil {
		return fmt.Errorf("failed to cleanup expired PATs: %w", err)
	}
	return nil
}
