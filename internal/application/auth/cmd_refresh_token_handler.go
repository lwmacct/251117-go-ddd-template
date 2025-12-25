package auth

import (
	"context"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/auth"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
)

// RefreshTokenHandler 刷新令牌命令处理器
type RefreshTokenHandler struct {
	userQueryRepo user.QueryRepository
	authService   auth.Service
}

// NewRefreshTokenHandler 创建刷新令牌命令处理器
func NewRefreshTokenHandler(
	userQueryRepo user.QueryRepository,
	authService auth.Service,
) *RefreshTokenHandler {
	return &RefreshTokenHandler{
		userQueryRepo: userQueryRepo,
		authService:   authService,
	}
}

// Handle 处理刷新令牌命令
func (h *RefreshTokenHandler) Handle(ctx context.Context, cmd RefreshTokenCommand) (*RefreshTokenResultDTO, error) {
	// 1. 验证 refresh token
	userID, err := h.authService.ValidateRefreshToken(ctx, cmd.RefreshToken)
	if err != nil {
		return nil, err
	}

	// 2. 获取用户信息
	u, err := h.userQueryRepo.GetByIDWithRoles(ctx, userID)
	if err != nil {
		return nil, auth.ErrUserNotFound
	}

	// 3. 检查用户状态
	if !u.CanLogin() {
		if u.IsBanned() {
			return nil, auth.ErrUserBanned
		}
		return nil, auth.ErrUserInactive
	}

	// 4. 生成新的访问令牌（新架构：不传递 roles，权限从缓存查询）
	accessToken, expiresAt, err := h.authService.GenerateAccessToken(ctx, u.ID, u.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// 5. 生成新的刷新令牌
	newRefreshToken, _, err := h.authService.GenerateRefreshToken(ctx, u.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &RefreshTokenResultDTO{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(expiresAt.Sub(expiresAt).Seconds()),
	}, nil
}
