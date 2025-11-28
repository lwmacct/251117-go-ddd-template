// Package auth 提供认证领域服务的基础设施实现。
//
// 本包实现 domain/auth.Service 接口，提供：
//   - 密码管理：BCrypt 哈希生成与验证
//   - JWT 令牌：访问令牌和刷新令牌的生成与验证
//   - PAT 令牌：个人访问令牌的生成与哈希
//   - 密码策略：可配置的密码强度验证
//
// 组件：
//   - JWTManager: JWT 令牌的生成、签名和验证
//   - TokenGenerator: PAT 等安全令牌的生成
//   - authServiceImpl: domain/auth.Service 的完整实现
//
// 安全特性：
//   - 使用 BCrypt 进行密码哈希（默认 cost=10）
//   - JWT 使用 HS256 签名算法
//   - Token 过期时间可配置
package auth

import (
	"context"
	"fmt"
	"time"

	domainAuth "github.com/lwmacct/251117-go-ddd-template/internal/domain/auth"
	"golang.org/x/crypto/bcrypt"
)

// authServiceImpl 认证服务实现
type authServiceImpl struct {
	jwtManager     *JWTManager
	tokenGenerator *TokenGenerator
	passwordPolicy *domainAuth.PasswordPolicy
}

// NewAuthService 创建认证服务实例
func NewAuthService(
	jwtManager *JWTManager,
	tokenGenerator *TokenGenerator,
	passwordPolicy *domainAuth.PasswordPolicy,
) domainAuth.Service {
	if passwordPolicy == nil {
		passwordPolicy = domainAuth.DefaultPasswordPolicy()
	}
	return &authServiceImpl{
		jwtManager:     jwtManager,
		tokenGenerator: tokenGenerator,
		passwordPolicy: passwordPolicy,
	}
}

// VerifyPassword 验证密码是否正确
func (s *authServiceImpl) VerifyPassword(ctx context.Context, hashedPassword, plainPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword)); err != nil {
		return domainAuth.ErrPasswordMismatch
	}
	return nil
}

// GeneratePasswordHash 生成密码哈希
func (s *authServiceImpl) GeneratePasswordHash(ctx context.Context, password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedBytes), nil
}

// ValidatePasswordPolicy 验证密码是否符合策略
func (s *authServiceImpl) ValidatePasswordPolicy(ctx context.Context, password string) error {
	if !s.passwordPolicy.Validate(password) {
		return domainAuth.ErrWeakPassword
	}
	return nil
}

// GenerateAccessToken 生成访问令牌
// 新架构：Token 只包含 user_id/username，权限信息从缓存实时查询
func (s *authServiceImpl) GenerateAccessToken(ctx context.Context, userID uint, username string) (string, time.Time, error) {
	token, err := s.jwtManager.GenerateAccessToken(userID, username, "")
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to generate access token: %w", err)
	}

	expiresAt := time.Now().Add(s.jwtManager.accessTokenDuration)
	return token, expiresAt, nil
}

// GenerateRefreshToken 生成刷新令牌
func (s *authServiceImpl) GenerateRefreshToken(ctx context.Context, userID uint) (string, time.Time, error) {
	token, err := s.jwtManager.GenerateRefreshToken(userID)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	expiresAt := time.Now().Add(s.jwtManager.refreshTokenDuration)
	return token, expiresAt, nil
}

// ValidateAccessToken 验证访问令牌
func (s *authServiceImpl) ValidateAccessToken(ctx context.Context, token string) (*domainAuth.TokenClaims, error) {
	claims, err := s.jwtManager.ValidateToken(token)
	if err != nil {
		return nil, domainAuth.ErrInvalidToken
	}

	if time.Unix(claims.ExpiresAt.Unix(), 0).Before(time.Now()) {
		return nil, domainAuth.ErrTokenExpired
	}

	return &domainAuth.TokenClaims{
		UserID:   claims.UserID,
		Username: claims.Username,
		Roles:    claims.Roles,
		Exp:      claims.ExpiresAt.Unix(),
	}, nil
}

// ValidateRefreshToken 验证刷新令牌
func (s *authServiceImpl) ValidateRefreshToken(ctx context.Context, token string) (uint, error) {
	claims, err := s.jwtManager.ValidateToken(token)
	if err != nil {
		return 0, domainAuth.ErrInvalidToken
	}

	if time.Unix(claims.ExpiresAt.Unix(), 0).Before(time.Now()) {
		return 0, domainAuth.ErrTokenExpired
	}

	return claims.UserID, nil
}

// GeneratePATToken 生成个人访问令牌
func (s *authServiceImpl) GeneratePATToken(ctx context.Context) (string, error) {
	plainToken, _, _, err := s.tokenGenerator.GeneratePAT()
	if err != nil {
		return "", fmt.Errorf("failed to generate PAT: %w", err)
	}
	return plainToken, nil
}

// HashPATToken 哈希个人访问令牌（用于存储）
func (s *authServiceImpl) HashPATToken(ctx context.Context, token string) string {
	return s.tokenGenerator.HashToken(token)
}
