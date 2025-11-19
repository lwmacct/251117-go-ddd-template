package pat

import "context"

// Repository defines the interface for PersonalAccessToken persistence
type Repository interface {
	// Create creates a new personal access token
	Create(ctx context.Context, pat *PersonalAccessToken) error

	// FindByToken retrieves a token by its hash value
	// This is used for authentication - token should be SHA-256 hashed before lookup
	FindByToken(ctx context.Context, tokenHash string) (*PersonalAccessToken, error)

	// FindByID retrieves a token by its ID
	FindByID(ctx context.Context, id uint) (*PersonalAccessToken, error)

	// FindByPrefix retrieves a token by its prefix (for user identification)
	FindByPrefix(ctx context.Context, prefix string) (*PersonalAccessToken, error)

	// ListByUser retrieves all tokens for a specific user
	ListByUser(ctx context.Context, userID uint) ([]*PersonalAccessToken, error)

	// Update updates an existing token (mainly for LastUsedAt)
	Update(ctx context.Context, pat *PersonalAccessToken) error

	// Delete hard deletes a token by ID
	Delete(ctx context.Context, id uint) error

	// Revoke revokes a token by setting its status to "revoked"
	Revoke(ctx context.Context, id uint) error

	// RevokeByUserID revokes all tokens for a specific user
	RevokeByUserID(ctx context.Context, userID uint) error

	// CleanupExpired deletes expired tokens (soft delete)
	CleanupExpired(ctx context.Context) error
}
