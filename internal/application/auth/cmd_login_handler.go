package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/auth"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/captcha"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/twofa"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
	authInfra "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/auth"
)

// LoginHandler 登录命令处理器
type LoginHandler struct {
	userQueryRepo      user.QueryRepository
	captchaCommandRepo captcha.CommandRepository
	twofaQueryRepo     twofa.QueryRepository
	authService        auth.Service
	loginSession       *authInfra.LoginSessionService
}

// NewLoginHandler 创建登录命令处理器
func NewLoginHandler(
	userQueryRepo user.QueryRepository,
	captchaCommandRepo captcha.CommandRepository,
	twofaQueryRepo twofa.QueryRepository,
	authService auth.Service,
	loginSession *authInfra.LoginSessionService,
) *LoginHandler {
	return &LoginHandler{
		userQueryRepo:      userQueryRepo,
		captchaCommandRepo: captchaCommandRepo,
		twofaQueryRepo:     twofaQueryRepo,
		authService:        authService,
		loginSession:       loginSession,
	}
}

// Handle 处理登录命令
func (h *LoginHandler) Handle(ctx context.Context, cmd LoginCommand) (*LoginResultDTO, error) {
	// 1. 验证图形验证码
	valid, err := h.captchaCommandRepo.Verify(ctx, cmd.CaptchaID, cmd.Captcha)
	if err != nil {
		return nil, fmt.Errorf("failed to verify captcha: %w", err)
	}
	if !valid {
		return nil, auth.ErrInvalidCaptcha
	}

	// 2. 查找用户（支持用户名或邮箱登录）
	var u *user.User
	if u, err = h.userQueryRepo.GetByUsernameWithRoles(ctx, cmd.Account); err != nil {
		// 尝试通过邮箱查找
		if u, err = h.userQueryRepo.GetByEmailWithRoles(ctx, cmd.Account); err != nil {
			return nil, auth.ErrInvalidCredentials
		}
	}

	// 3. 检查用户状态
	if !u.CanLogin() {
		if u.IsBanned() {
			return nil, auth.ErrUserBanned
		}
		if u.IsInactive() {
			return nil, auth.ErrUserInactive
		}
	}

	// 4. 验证密码
	if err = h.authService.VerifyPassword(ctx, u.Password, cmd.Password); err != nil {
		return nil, auth.ErrInvalidCredentials
	}

	// 5. 检查是否启用 2FA
	tfa, err := h.twofaQueryRepo.FindByUserID(ctx, u.ID)
	if err == nil && tfa != nil && tfa.Enabled {
		// 需要 2FA 验证，生成临时 session token
		sessionToken, sessionErr := h.loginSession.GenerateSessionToken(ctx, u.ID, cmd.Account)
		if sessionErr != nil {
			return nil, fmt.Errorf("failed to generate session token: %w", sessionErr)
		}

		return &LoginResultDTO{
			Requires2FA:  true,
			SessionToken: sessionToken,
			UserID:       u.ID,
			Username:     u.Username,
		}, nil
	}

	// 6. 生成访问令牌（新架构：不传递 roles，权限从缓存查询）
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
