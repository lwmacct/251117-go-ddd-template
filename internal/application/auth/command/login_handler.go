// Package command 定义认证命令处理器
package command

import (
	"context"
	"fmt"

	domainAuth "github.com/lwmacct/251117-go-ddd-template/internal/domain/auth"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/captcha"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/twofa"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
)

// LoginResult 登录结果
type LoginResult struct {
	AccessToken  string
	RefreshToken string
	TokenType    string
	ExpiresIn    int
	UserID       uint
	Username     string
	Requires2FA  bool
	SessionToken string
}

// LoginHandler 登录命令处理器
type LoginHandler struct {
	userQueryRepo    user.QueryRepository
	captchaQueryRepo captcha.Repository
	twofaQueryRepo   twofa.QueryRepository
	authService      domainAuth.Service
}

// NewLoginHandler 创建登录命令处理器
func NewLoginHandler(
	userQueryRepo user.QueryRepository,
	captchaQueryRepo captcha.Repository,
	twofaQueryRepo twofa.QueryRepository,
	authService domainAuth.Service,
) *LoginHandler {
	return &LoginHandler{
		userQueryRepo:    userQueryRepo,
		captchaQueryRepo: captchaQueryRepo,
		twofaQueryRepo:   twofaQueryRepo,
		authService:      authService,
	}
}

// Handle 处理登录命令
func (h *LoginHandler) Handle(ctx context.Context, cmd LoginCommand) (*LoginResult, error) {
	// 1. 验证图形验证码
	valid, err := h.captchaQueryRepo.Verify(ctx, cmd.CaptchaID, cmd.Captcha)
	if err != nil {
		return nil, fmt.Errorf("failed to verify captcha: %w", err)
	}
	if !valid {
		return nil, domainAuth.ErrInvalidCaptcha
	}

	// 2. 查找用户（支持用户名或邮箱登录）
	var u *user.User
	if u, err = h.userQueryRepo.GetByUsernameWithRoles(ctx, cmd.Login); err != nil {
		// 尝试通过邮箱查找
		if u, err = h.userQueryRepo.GetByEmailWithRoles(ctx, cmd.Login); err != nil {
			return nil, domainAuth.ErrInvalidCredentials
		}
	}

	// 3. 检查用户状态
	if !u.CanLogin() {
		if u.IsBanned() {
			return nil, domainAuth.ErrUserBanned
		}
		if u.IsInactive() {
			return nil, domainAuth.ErrUserInactive
		}
	}

	// 4. 验证密码
	if err := h.authService.VerifyPassword(ctx, u.Password, cmd.Password); err != nil {
		return nil, domainAuth.ErrInvalidCredentials
	}

	// 5. 检查是否启用 2FA
	tfa, err := h.twofaQueryRepo.FindByUserID(ctx, u.ID)
	if err == nil && tfa != nil && tfa.Enabled {
		// 需要 2FA 验证，生成临时 session token
		sessionToken, _, err := h.authService.GenerateRefreshToken(ctx, u.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to generate session token: %w", err)
		}

		return &LoginResult{
			Requires2FA:  true,
			SessionToken: sessionToken,
			UserID:       u.ID,
			Username:     u.Username,
		}, nil
	}

	// 6. 生成访问令牌
	accessToken, expiresAt, err := h.authService.GenerateAccessToken(ctx, u.ID, u.Username, u.GetRoleNames())
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, _, err := h.authService.GenerateRefreshToken(ctx, u.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	expiresIn := int(expiresAt.Sub(expiresAt).Seconds())

	return &LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
		UserID:       u.ID,
		Username:     u.Username,
		Requires2FA:  false,
	}, nil
}
