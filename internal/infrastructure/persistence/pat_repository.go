package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/pat"
	"gorm.io/gorm"
)

// patRepository implements pat.CommandRepository and pat.QueryRepository using GORM
type patRepository struct {
	db *gorm.DB
}

// NewPATRepository creates a new PAT repository instance (deprecated, use NewPATCommandRepository/NewPATQueryRepository)
func NewPATRepository(db *gorm.DB) pat.Repository {
	return &patRepository{db: db}
}

// NewPATCommandRepository creates a new PAT command repository instance
func NewPATCommandRepository(db *gorm.DB) pat.CommandRepository {
	return &patRepository{db: db}
}

// NewPATQueryRepository creates a new PAT query repository instance
func NewPATQueryRepository(db *gorm.DB) pat.QueryRepository {
	return &patRepository{db: db}
}

// Create creates a new personal access token
func (r *patRepository) Create(ctx context.Context, token *pat.PersonalAccessToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

// FindByToken retrieves a token by its hash value
func (r *patRepository) FindByToken(ctx context.Context, tokenHash string) (*pat.PersonalAccessToken, error) {
	var token pat.PersonalAccessToken
	err := r.db.WithContext(ctx).
		Where("token = ?", tokenHash).
		First(&token).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("token not found")
		}
		return nil, err
	}

	return &token, nil
}

// FindByID retrieves a token by its ID
func (r *patRepository) FindByID(ctx context.Context, id uint) (*pat.PersonalAccessToken, error) {
	var token pat.PersonalAccessToken
	err := r.db.WithContext(ctx).
		First(&token, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("token not found")
		}
		return nil, err
	}

	return &token, nil
}

// FindByPrefix retrieves a token by its prefix
func (r *patRepository) FindByPrefix(ctx context.Context, prefix string) (*pat.PersonalAccessToken, error) {
	var token pat.PersonalAccessToken
	err := r.db.WithContext(ctx).
		Where("token_prefix = ?", prefix).
		First(&token).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("token not found")
		}
		return nil, err
	}

	return &token, nil
}

// ListByUser retrieves all tokens for a specific user
func (r *patRepository) ListByUser(ctx context.Context, userID uint) ([]*pat.PersonalAccessToken, error) {
	var tokens []*pat.PersonalAccessToken
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&tokens).Error

	if err != nil {
		return nil, err
	}

	return tokens, nil
}

// Update updates an existing token
func (r *patRepository) Update(ctx context.Context, token *pat.PersonalAccessToken) error {
	return r.db.WithContext(ctx).Save(token).Error
}

// Delete hard deletes a token by ID
func (r *patRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Unscoped(). // Hard delete
		Delete(&pat.PersonalAccessToken{}, id).Error
}

// Revoke revokes a token by setting its status to "revoked"
func (r *patRepository) Revoke(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Model(&pat.PersonalAccessToken{}).
		Where("id = ?", id).
		Update("status", "revoked").Error
}

// RevokeByUserID revokes all tokens for a specific user
func (r *patRepository) RevokeByUserID(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).
		Model(&pat.PersonalAccessToken{}).
		Where("user_id = ? AND status = ?", userID, "active").
		Update("status", "revoked").Error
}

// CleanupExpired soft deletes expired tokens
func (r *patRepository) CleanupExpired(ctx context.Context) error {
	now := time.Now()

	// Find all expired tokens
	return r.db.WithContext(ctx).
		Model(&pat.PersonalAccessToken{}).
		Where("expires_at IS NOT NULL AND expires_at < ? AND status != ?", now, "expired").
		Update("status", "expired").Error
}
