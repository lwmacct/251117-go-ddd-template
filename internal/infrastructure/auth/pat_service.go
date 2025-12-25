package auth

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/pat"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
)

// PATService handles Personal Access Token business logic
type PATService struct {
	patCommandRepo pat.CommandRepository
	patQueryRepo   pat.QueryRepository
	userQueryRepo  user.QueryRepository
	tokenGen       *TokenGenerator
}

// NewPATService creates a new PAT service
func NewPATService(
	patCommandRepo pat.CommandRepository,
	patQueryRepo pat.QueryRepository,
	userQueryRepo user.QueryRepository,
	tokenGen *TokenGenerator,
) *PATService {
	return &PATService{
		patCommandRepo: patCommandRepo,
		patQueryRepo:   patQueryRepo,
		userQueryRepo:  userQueryRepo,
		tokenGen:       tokenGen,
	}
}

// CreateTokenRequest 创建令牌请求（Infrastructure 层内部使用）
type CreateTokenRequest struct {
	Name        string
	Permissions []string
	ExpiresAt   *string
	ExpiresIn   *int
	IPWhitelist []string
	Description string
}

// CreateTokenResult 创建令牌结果（返回领域实体和明文 token）
type CreateTokenResult struct {
	Token      *pat.PersonalAccessToken
	PlainToken string
}

// CreateToken creates a new Personal Access Token
// 返回领域实体，由 Application 层负责转换为 DTO
func (s *PATService) CreateToken(ctx context.Context, req *CreateTokenRequest, userID uint) (*CreateTokenResult, error) {
	// Validate user exists and get their permissions
	u, err := s.userQueryRepo.GetByIDWithRoles(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Get all user permissions
	userPermissions := u.GetPermissionCodes()
	if len(userPermissions) == 0 {
		return nil, errors.New("user has no permissions")
	}

	requestedPerms := req.Permissions
	if len(requestedPerms) == 0 {
		requestedPerms = userPermissions
	}

	// Validate requested permissions are subset of user's permissions
	if err = s.validatePermissions(requestedPerms, userPermissions); err != nil {
		return nil, err
	}

	// Generate token
	plainToken, tokenHash, prefix, err := s.tokenGen.GeneratePAT()
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Calculate expiry
	var expiresAt *time.Time
	if req.ExpiresAt != nil && *req.ExpiresAt != "" {
		parsed, parseErr := time.Parse(time.RFC3339, *req.ExpiresAt)
		if parseErr != nil {
			return nil, fmt.Errorf("invalid expires_at format: %w", parseErr)
		}
		tmp := parsed.UTC()
		expiresAt = &tmp
	} else if req.ExpiresIn != nil && *req.ExpiresIn > 0 {
		tmp := time.Now().Add(time.Duration(*req.ExpiresIn) * 24 * time.Hour).UTC()
		expiresAt = &tmp
	}

	// Create PAT record
	token := &pat.PersonalAccessToken{
		UserID:      userID,
		Name:        req.Name,
		Token:       tokenHash,
		TokenPrefix: prefix,
		Permissions: requestedPerms,
		ExpiresAt:   expiresAt,
		Status:      "active",
		IPWhitelist: req.IPWhitelist,
		Description: req.Description,
	}

	if err := s.patCommandRepo.Create(ctx, token); err != nil {
		return nil, fmt.Errorf("failed to create token: %w", err)
	}

	// Return response with plain token (only shown once!)
	return &CreateTokenResult{Token: token, PlainToken: plainToken}, nil
}

// ListTokens lists all tokens for a user (without sensitive data)
// 返回领域实体列表，由 Application 层负责转换为 DTO
func (s *PATService) ListTokens(ctx context.Context, userID uint) ([]*pat.PersonalAccessToken, error) {
	tokens, err := s.patQueryRepo.ListByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list tokens: %w", err)
	}

	return tokens, nil
}

// DisableToken pauses a token (GitHub-style)
func (s *PATService) DisableToken(ctx context.Context, userID, tokenID uint) error {
	// Get token to verify ownership
	token, err := s.patQueryRepo.FindByID(ctx, tokenID)
	if err != nil {
		return fmt.Errorf("token not found: %w", err)
	}

	// Verify ownership
	if token.UserID != userID {
		return errors.New("unauthorized: token does not belong to user")
	}

	if token.Status == "disabled" {
		return nil
	}

	// Disable token
	if err := s.patCommandRepo.Disable(ctx, tokenID); err != nil {
		return fmt.Errorf("failed to disable token: %w", err)
	}

	return nil
}

// EnableToken enables a disabled token
func (s *PATService) EnableToken(ctx context.Context, userID, tokenID uint) error {
	token, err := s.patQueryRepo.FindByID(ctx, tokenID)
	if err != nil {
		return fmt.Errorf("token not found: %w", err)
	}

	if token.UserID != userID {
		return errors.New("unauthorized: token does not belong to user")
	}

	if token.IsExpired() {
		return errors.New("token is expired and cannot be enabled")
	}

	if token.Status == "active" {
		return nil
	}

	if err := s.patCommandRepo.Enable(ctx, tokenID); err != nil {
		return fmt.Errorf("failed to enable token: %w", err)
	}

	return nil
}

// DeleteToken deletes a token permanently
func (s *PATService) DeleteToken(ctx context.Context, userID, tokenID uint) error {
	token, err := s.patQueryRepo.FindByID(ctx, tokenID)
	if err != nil {
		return fmt.Errorf("token not found: %w", err)
	}
	if token.UserID != userID {
		return errors.New("unauthorized: token does not belong to user")
	}

	if err := s.patCommandRepo.Delete(ctx, tokenID); err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
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
	token, err := s.patQueryRepo.FindByToken(ctx, tokenHash)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	// Check if token is active
	if !token.IsActive() {
		return nil, errors.New("token is inactive or expired")
	}

	// Update last used time (asynchronously to avoid blocking)
	// 使用 WithoutCancel 保留 trace 信息，但不受请求取消影响
	go func(updateCtx context.Context) {
		now := time.Now()
		token.LastUsedAt = &now
		_ = s.patCommandRepo.Update(updateCtx, token)
	}(context.WithoutCancel(ctx))

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
		if !slices.Contains(token.IPWhitelist, clientIP) {
			return nil, fmt.Errorf("access denied: IP %s not in whitelist", clientIP)
		}
	}

	return token, nil
}

// DeleteAllUserTokens deletes all tokens for a user (e.g., on password change)
func (s *PATService) DeleteAllUserTokens(ctx context.Context, userID uint) error {
	return s.patCommandRepo.DeleteByUserID(ctx, userID)
}

// CleanupExpiredTokens cleans up expired tokens (should be run periodically)
func (s *PATService) CleanupExpiredTokens(ctx context.Context) error {
	return s.patCommandRepo.CleanupExpired(ctx)
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
