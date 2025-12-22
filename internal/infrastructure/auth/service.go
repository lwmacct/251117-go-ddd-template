package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	userdto "github.com/lwmacct/251117-go-ddd-template/internal/application/user"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/captcha"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/twofa"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

// Service 认证服务
type Service struct {
	userCommandRepo    user.CommandRepository
	userQueryRepo      user.QueryRepository
	twofaCommandRepo   twofa.CommandRepository
	twofaQueryRepo     twofa.QueryRepository
	captchaCommandRepo captcha.CommandRepository
	jwtManager         *JWTManager
	sessionService     *LoginSessionService
}

// NewService 创建认证服务
func NewService(
	userCommandRepo user.CommandRepository,
	userQueryRepo user.QueryRepository,
	twofaCommandRepo twofa.CommandRepository,
	twofaQueryRepo twofa.QueryRepository,
	captchaCommandRepo captcha.CommandRepository,
	jwtManager *JWTManager,
	sessionService *LoginSessionService,
) *Service {
	return &Service{
		userCommandRepo:    userCommandRepo,
		userQueryRepo:      userQueryRepo,
		twofaCommandRepo:   twofaCommandRepo,
		twofaQueryRepo:     twofaQueryRepo,
		captchaCommandRepo: captchaCommandRepo,
		jwtManager:         jwtManager,
		sessionService:     sessionService,
	}
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"max=100"`
}

// LoginRequest 登录请求
// 支持两种登录流程：
// 1. 第一次登录：login + password + captcha_id + captcha (必须)
// 2. 第二次登录（2FA）：session_token + two_factor_code (必须)
type LoginRequest struct {
	Login         string `json:"login,omitempty"`           // 用户名或邮箱（第一次登录必须）
	Password      string `json:"password,omitempty"`        // 密码（第一次登录必须）
	CaptchaID     string `json:"captcha_id,omitempty"`      // 验证码ID（第一次登录必须）
	Captcha       string `json:"captcha,omitempty"`         // 验证码（第一次登录必须）
	TwoFactorCode string `json:"two_factor_code,omitempty"` // 2FA验证码（第二次登录必须）
	SessionToken  string `json:"session_token,omitempty"`   // 临时会话token（第二次登录必须）
}

// AuthResponse 认证响应
type AuthResponse struct {
	AccessToken  string                         `json:"access_token,omitempty"`  // 第一次登录时不返回（需要2FA验证）
	RefreshToken string                         `json:"refresh_token,omitempty"` // 第一次登录时不返回（需要2FA验证）
	TokenType    string                         `json:"token_type,omitempty"`
	ExpiresIn    int                            `json:"expires_in,omitempty"` // 秒
	User         *userdto.UserWithRolesResponse `json:"user,omitempty"`
	SessionToken string                         `json:"session_token,omitempty"` // 临时会话token（第一次登录后返回，用于第二次2FA验证）
	Requires2FA  bool                           `json:"requires_2fa,omitempty"`  // 是否需要2FA验证
}

// TwoFactorRequiredError 表示需要2FA验证的错误
type TwoFactorRequiredError struct{}

func (e *TwoFactorRequiredError) Error() string {
	return "two factor authentication required"
}

// Register 注册新用户
func (s *Service) Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error) {
	// 检查用户名是否已存在
	if _, err := s.userQueryRepo.GetByUsername(ctx, req.Username); err == nil {
		return nil, errors.New("username already exists")
	}

	// 检查邮箱是否已存在
	if _, err := s.userQueryRepo.GetByEmail(ctx, req.Email); err == nil {
		return nil, errors.New("email already exists")
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 创建用户
	newUser := &user.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		FullName: req.FullName,
		Status:   "active",
	}

	if err = s.userCommandRepo.Create(ctx, newUser); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Reload user with roles
	userWithRoles, err := s.userQueryRepo.GetByIDWithRoles(ctx, newUser.ID)
	if err != nil {
		// Fallback to empty roles if failed to load
		userWithRoles = newUser
	}

	// 生成 token（新架构：不传递 roles/permissions，权限从缓存查询）
	accessToken, refreshToken, err := s.jwtManager.GenerateTokenPair(userWithRoles.ID, userWithRoles.Username, userWithRoles.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(s.jwtManager.accessTokenDuration.Seconds()),
		User:         userdto.ToUserWithRolesResponse(userWithRoles),
	}, nil
}

// Login 用户登录
// 🔒 安全策略：使用临时session token防止2FA暴力破解
// 流程1：第一次登录（账号密码+图形验证码）
//   - 验证图形验证码
//   - 验证账号密码
//   - 如果启用2FA，生成临时session token返回，不返回访问令牌
//
// 流程2：第二次登录（session token + 2FA验证码）
//   - 验证session token，获取用户ID（防止2FA暴力破解）
//   - 验证2FA验证码
//   - 生成访问令牌返回
func (s *Service) Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error) {
	// 根据是否有 session token 分流处理
	if req.SessionToken != "" {
		return s.loginWith2FA(ctx, req)
	}
	return s.loginWithCredentials(ctx, req)
}

// RefreshToken 刷新访问令牌
func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (*AuthResponse, error) {
	// 验证刷新令牌
	claims, err := s.jwtManager.ValidateToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// 获取用户信息（包含角色和权限）
	u, err := s.userQueryRepo.GetByIDWithRoles(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// 检查用户状态
	if u.Status != "active" {
		return nil, fmt.Errorf("user account is %s", u.Status)
	}

	// 生成新的 token 对（新架构：不传递 roles/permissions，权限从缓存查询）
	accessToken, newRefreshToken, err := s.jwtManager.GenerateTokenPair(u.ID, u.Username, u.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(s.jwtManager.accessTokenDuration.Seconds()),
		User:         userdto.ToUserWithRolesResponse(u),
	}, nil
}

// ============================================================================
// 以下为私有辅助方法
// ============================================================================

// loginWith2FA 流程2：第二次登录（2FA验证）
func (s *Service) loginWith2FA(ctx context.Context, req *LoginRequest) (*AuthResponse, error) {
	// 参数验证
	if req.TwoFactorCode == "" {
		return nil, errors.New("two factor code is required")
	}

	// 1. 验证session token（防止2FA暴力破解）
	sessionData, err := s.sessionService.VerifySessionToken(ctx, req.SessionToken)
	if err != nil {
		return nil, errors.New("session expired or invalid, please login again")
	}

	// 2. 根据session中的用户ID查找用户
	u, err := s.userQueryRepo.GetByIDWithRoles(ctx, sessionData.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// 3. 检查用户状态
	if err := s.checkUserStatus(u); err != nil {
		return nil, err
	}

	// 4. 验证2FA验证码
	if err := s.verify2FACode(ctx, u.ID, req.TwoFactorCode); err != nil {
		return nil, err
	}

	// 5. 2FA验证成功，生成访问令牌
	return s.generateTokens(u)
}

// loginWithCredentials 流程1：第一次登录（账号密码验证）
func (s *Service) loginWithCredentials(ctx context.Context, req *LoginRequest) (*AuthResponse, error) {
	// 参数验证
	if err := s.validateCredentialRequest(req); err != nil {
		return nil, err
	}

	// 1. 验证图形验证码
	if err := s.verifyCaptcha(ctx, req.CaptchaID, req.Captcha); err != nil {
		return nil, err
	}

	// 2. 查找用户（支持用户名或邮箱）
	u, err := s.findUserByLogin(ctx, req.Login)
	if err != nil {
		return nil, err
	}

	// 3. 验证密码
	if err := s.verifyPassword(u.Password, req.Password); err != nil {
		return nil, err
	}

	// 4. 检查用户状态
	if err := s.checkUserStatus(u); err != nil {
		return nil, err
	}

	// 5. 检查是否启用2FA，返回相应结果
	return s.handlePostCredentialLogin(ctx, u, req.Login)
}

// generateTokens 生成访问令牌
func (s *Service) generateTokens(u *user.User) (*AuthResponse, error) {
	// 新架构：不传递 roles/permissions，权限从缓存查询
	accessToken, refreshToken, err := s.jwtManager.GenerateTokenPair(u.ID, u.Username, u.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(s.jwtManager.accessTokenDuration.Seconds()),
		User:         userdto.ToUserWithRolesResponse(u),
	}, nil
}

// validateCredentialRequest 验证第一次登录请求参数
func (s *Service) validateCredentialRequest(req *LoginRequest) error {
	if req.Login == "" {
		return errors.New("login is required")
	}
	if req.Password == "" {
		return errors.New("password is required")
	}
	if req.CaptchaID == "" || req.Captcha == "" {
		return errors.New("captcha is required")
	}
	return nil
}

// checkUserStatus 检查用户状态
func (s *Service) checkUserStatus(u *user.User) error {
	if u.Status != "active" {
		return fmt.Errorf("user account is %s", u.Status)
	}
	return nil
}

// verifyCaptcha 验证图形验证码
func (s *Service) verifyCaptcha(ctx context.Context, captchaID, captcha string) error {
	valid, err := s.captchaCommandRepo.Verify(ctx, captchaID, captcha)
	if err != nil || !valid {
		return errors.New("invalid or expired captcha")
	}
	return nil
}

// findUserByLogin 根据用户名或邮箱查找用户
func (s *Service) findUserByLogin(ctx context.Context, login string) (*user.User, error) {
	// 尝试用户名
	u, err := s.userQueryRepo.GetByUsernameWithRoles(ctx, login)
	if err == nil {
		return u, nil
	}

	// 尝试邮箱
	tempUser, emailErr := s.userQueryRepo.GetByEmail(ctx, login)
	if emailErr != nil {
		return nil, errors.New("invalid credentials")
	}

	// Reload with roles (如果失败，返回不带角色的用户)
	u, _ = s.userQueryRepo.GetByIDWithRoles(ctx, tempUser.ID)
	if u == nil {
		return tempUser, nil
	}
	return u, nil
}

// verifyPassword 验证密码
func (s *Service) verifyPassword(hashedPassword, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return errors.New("invalid credentials")
	}
	return nil
}

// verify2FACode 验证2FA验证码
func (s *Service) verify2FACode(ctx context.Context, userID uint, code string) error {
	tfa, err := s.twofaQueryRepo.FindByUserID(ctx, userID)
	if err != nil || tfa == nil || !tfa.Enabled {
		return errors.New("2FA not enabled for this account")
	}

	// 使用2FA服务验证（支持TOTP和恢复码）
	twofaSvc := &twofaService{
		twofaCommandRepo: s.twofaCommandRepo,
		twofaQueryRepo:   s.twofaQueryRepo,
	}
	valid, err := twofaSvc.Verify(ctx, userID, code)
	if err != nil || !valid {
		return errors.New("invalid two factor code")
	}
	return nil
}

// handlePostCredentialLogin 处理密码验证后的2FA检查
func (s *Service) handlePostCredentialLogin(ctx context.Context, u *user.User, login string) (*AuthResponse, error) {
	// 检查是否启用2FA
	tfa, err := s.twofaQueryRepo.FindByUserID(ctx, u.ID)
	if err != nil || tfa == nil || !tfa.Enabled {
		// 未启用2FA，直接生成访问令牌
		return s.generateTokens(u)
	}

	// 用户启用了 2FA，生成临时session token
	sessionToken, err := s.sessionService.GenerateSessionToken(ctx, u.ID, login)
	if err != nil {
		return nil, errors.New("failed to generate session token")
	}

	// 返回session token，前端需要再次提交2FA验证码
	return &AuthResponse{
		SessionToken: sessionToken,
		Requires2FA:  true,
	}, &TwoFactorRequiredError{}
}

// twofaService 临时 2FA 服务包装（避免循环依赖）
type twofaService struct {
	twofaCommandRepo twofa.CommandRepository
	twofaQueryRepo   twofa.QueryRepository
}

func (s *twofaService) Verify(ctx context.Context, userID uint, code string) (bool, error) {
	tfa, err := s.twofaQueryRepo.FindByUserID(ctx, userID)
	if err != nil {
		return false, err
	}
	if tfa == nil || !tfa.Enabled {
		return false, errors.New("2FA not enabled")
	}

	// 首先尝试 TOTP 验证
	if totp.Validate(code, tfa.Secret) {
		return true, nil
	}

	// TOTP 验证失败，尝试恢复码
	code = strings.TrimSpace(code)
	for i, recoveryCode := range tfa.RecoveryCodes {
		if recoveryCode == code {
			// 移除已使用的恢复码
			tfa.RecoveryCodes = append(tfa.RecoveryCodes[:i], tfa.RecoveryCodes[i+1:]...)
			if err := s.twofaCommandRepo.CreateOrUpdate(ctx, tfa); err != nil {
				return false, fmt.Errorf("failed to update recovery codes: %w", err)
			}
			return true, nil
		}
	}

	return false, nil
}
