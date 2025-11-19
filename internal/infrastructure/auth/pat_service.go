package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/pat"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
)

// PATService handles Personal Access Token business logic
type PATService struct {
	patRepo      pat.Repository
	userRepo     user.Repository
	tokenGen     *TokenGenerator
}

// NewPATService creates a new PAT service
func NewPATService(patRepo pat.Repository, userRepo user.Repository, tokenGen *TokenGenerator) *PATService {
	return &PATService{
		patRepo:  patRepo,
		userRepo: userRepo,
		tokenGen: tokenGen,
	}
}

// CreateToken creates a new Personal Access Token
func (s *PATService) CreateToken(ctx context.Context, req *pat.CreateTokenRequest) (*pat.TokenResponse, error) {
	// Validate user exists and get their permissions
	u, err := s.userRepo.GetByIDWithRoles(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Get all user permissions
	userPermissions := u.GetPermissionCodes()
	if len(userPermissions) == 0 {
		return nil, errors.New("user has no permissions")
	}

	// Validate requested permissions are subset of user's permissions
	if err := s.validatePermissions(req.Permissions, userPermissions); err != nil {
		return nil, err
	}

	// Generate token
	plainToken, tokenHash, prefix, err := s.tokenGen.GeneratePAT()
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Calculate expiry
	var expiresAt *time.Time
	if req.ExpiresIn != nil {
		expiry := time.Now().Add(time.Duration(*req.ExpiresIn) * 24 * time.Hour)
		expiresAt = &expiry
	}

	// Create PAT record
	token := &pat.PersonalAccessToken{
		UserID:      req.UserID,
		Name:        req.Name,
		Token:       tokenHash,
		TokenPrefix: prefix,
		Permissions: req.Permissions,
		ExpiresAt:   expiresAt,
		Status:      "active",
		IPWhitelist: req.IPWhitelist,
		Description: req.Description,
	}

	if err := s.patRepo.Create(ctx, token); err != nil {
		return nil, fmt.Errorf("failed to create token: %w", err)
	}

	// Return response with plain token (only shown once!)
	return &pat.TokenResponse{
		Token:       plainToken, // Full plain token: pat_xxxxx_yyy...
		ID:          token.ID,
		Name:        token.Name,
		TokenPrefix: token.TokenPrefix,
		Permissions: token.Permissions,
		ExpiresAt:   token.ExpiresAt,
		CreatedAt:   token.CreatedAt,
	}, nil
}

// ListTokens lists all tokens for a user (without sensitive data)
func (s *PATService) ListTokens(ctx context.Context, userID uint) ([]*pat.TokenListItem, error) {
	tokens, err := s.patRepo.ListByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list tokens: %w", err)
	}

	items := make([]*pat.TokenListItem, len(tokens))
	for i, t := range tokens {
		items[i] = t.ToListItem()
	}

	return items, nil
}

// RevokeToken revokes a token by ID
func (s *PATService) RevokeToken(ctx context.Context, userID, tokenID uint) error {
	// Get token to verify ownership
	token, err := s.patRepo.FindByID(ctx, tokenID)
	if err != nil {
		return fmt.Errorf("token not found: %w", err)
	}

	// Verify ownership
	if token.UserID != userID {
		return errors.New("unauthorized: token does not belong to user")
	}

	// Revoke token
	if err := s.patRepo.Revoke(ctx, tokenID); err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}

	return nil
}

// ValidateToken validates a PAT and returns associated permissions
// This is used in the authentication middleware
func (s *PATService) ValidateToken(ctx context.Context, plainToken string) (*pat.PersonalAccessToken, error) {
	// Validate format
	if !s.tokenGen.ValidateTokenFormat(plainToken) {
		return nil, errors.New("invalid token format")
	}

	// Hash the token
	tokenHash := s.tokenGen.HashToken(plainToken)

	// Find token in database
	token, err := s.patRepo.FindByToken(ctx, tokenHash)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	// Check if token is active
	if !token.IsActive() {
		return nil, errors.New("token is inactive or expired")
	}

	// Update last used time (asynchronously to avoid blocking)
	go func() {
		now := time.Now()
		token.LastUsedAt = &now
		_ = s.patRepo.Update(context.Background(), token)
	}()

	return token, nil
}

// ValidateTokenWithIP validates a PAT and checks IP whitelist
func (s *PATService) ValidateTokenWithIP(ctx context.Context, plainToken, clientIP string) (*pat.PersonalAccessToken, error) {
	token, err := s.ValidateToken(ctx, plainToken)
	if err != nil {
		return nil, err
	}

	// Check IP whitelist if configured
	if len(token.IPWhitelist) > 0 {
		allowed := false
		for _, ip := range token.IPWhitelist {
			if ip == clientIP {
				allowed = true
				break
			}
		}
		if !allowed {
			return nil, fmt.Errorf("access denied: IP %s not in whitelist", clientIP)
		}
	}

	return token, nil
}

// RevokeAllUserTokens revokes all tokens for a user (e.g., on password change)
func (s *PATService) RevokeAllUserTokens(ctx context.Context, userID uint) error {
	return s.patRepo.RevokeByUserID(ctx, userID)
}

// CleanupExpiredTokens cleans up expired tokens (should be run periodically)
func (s *PATService) CleanupExpiredTokens(ctx context.Context) error {
	return s.patRepo.CleanupExpired(ctx)
}

// validatePermissions checks if requested permissions are a subset of user permissions
func (s *PATService) validatePermissions(requested, userPerms []string) error {
	userPermSet := make(map[string]bool)
	for _, perm := range userPerms {
		userPermSet[perm] = true
	}

	for _, perm := range requested {
		if !userPermSet[perm] {
			return fmt.Errorf("permission '%s' not granted to user", perm)
		}
	}

	return nil
}
