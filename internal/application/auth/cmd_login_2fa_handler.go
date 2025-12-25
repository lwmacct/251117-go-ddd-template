package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/auth"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
	authInfra "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/auth"
	twofaInfra "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/twofa"
)

// Login2FAHandler 二次认证登录命令处理器
type Login2FAHandler struct {
	userQueryRepo user.QueryRepository
	authService   auth.Service
	loginSession  *authInfra.LoginSessionService
	twofaService  *twofaInfra.Service
}

// NewLogin2FAHandler 创建二次认证登录命令处理器
func NewLogin2FAHandler(
	userQueryRepo user.QueryRepository,
	authService auth.Service,
	loginSession *authInfra.LoginSessionService,
	twofaService *twofaInfra.Service,
) *Login2FAHandler {
	return &Login2FAHandler{
		userQueryRepo: userQueryRepo,
		authService:   authService,
		loginSession:  loginSession,
		twofaService:  twofaService,
	}
}

// Handle 处理二次认证登录命令
func (h *Login2FAHandler) Handle(ctx context.Context, cmd Login2FACommand) (*LoginResultDTO, error) {
	// 1. 验证 session token（防止 2FA 暴力破解）
	sessionData, err := h.loginSession.VerifySessionToken(ctx, cmd.SessionToken)
	if err != nil {
		return nil, errors.New("session expired or invalid, please login again")
	}

	// 2. 验证 2FA 验证码
	valid, err := h.twofaService.Verify(ctx, sessionData.UserID, cmd.TwoFactorCode)
	if err != nil {
		return nil, fmt.Errorf("2FA verification failed: %w", err)
	}
	if !valid {
		return nil, errors.New("invalid two factor code")
	}

	// 3. 获取用户信息
	u, err := h.userQueryRepo.GetByIDWithRoles(ctx, sessionData.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// 4. 检查用户状态
	if !u.CanLogin() {
		if u.IsBanned() {
			return nil, auth.ErrUserBanned
		}
		if u.IsInactive() {
			return nil, auth.ErrUserInactive
		}
	}

	// 5. 生成访问令牌
	accessToken, expiresAt, err := h.authService.GenerateAccessToken(ctx, u.ID, u.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, _, err := h.authService.GenerateRefreshToken(ctx, u.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	expiresIn := int(time.Until(expiresAt).Seconds())

	return &LoginResultDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
		UserID:       u.ID,
		Username:     u.Username,
		Requires2FA:  false,
	}, nil
}
