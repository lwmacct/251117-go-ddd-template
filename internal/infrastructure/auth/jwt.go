package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims JWT 自定义声明
type Claims struct {
	jwt.RegisteredClaims

	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	// Roles 和 Permissions 字段已废弃，保留仅用于向后兼容旧 token
	// 新 token 不再包含这些字段，权限信息改为从缓存/数据库实时查询
	Roles       []string `json:"roles,omitempty"`       // Deprecated: 仅用于向后兼容
	Permissions []string `json:"permissions,omitempty"` // Deprecated: 仅用于向后兼容
}

// JWTManager JWT 管理器
type JWTManager struct {
	secretKey            string
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

// NewJWTManager 创建 JWT 管理器
func NewJWTManager(secretKey string, accessTokenDuration, refreshTokenDuration time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:            secretKey,
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
	}
}

// GenerateAccessToken 生成访问令牌
// 新架构：Token 只包含 user_id/username/email，权限信息从缓存/数据库实时查询
func (m *JWTManager) GenerateAccessToken(userID uint, username, email string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		Email:    email,
		// Roles 和 Permissions 不再包含在 token 中
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.accessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secretKey))
}

// GenerateRefreshToken 生成刷新令牌
// Refresh Token 同样不包含权限信息，刷新时从数据库查询最新权限
func (m *JWTManager) GenerateRefreshToken(userID uint) (string, error) {
	claims := Claims{
		UserID: userID,
		// Refresh Token 只需要 user_id，username/email 在刷新时重新获取
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.refreshTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secretKey))
}

// ValidateToken 验证令牌
func (m *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (any, error) {
			// 验证签名方法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(m.secretKey), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// GenerateTokenPair 生成访问令牌和刷新令牌对
//
//nolint:nonamedreturns // named returns for self-documenting API
func (m *JWTManager) GenerateTokenPair(userID uint, username, email string) (accessToken, refreshToken string, err error) {
	accessToken, err = m.GenerateAccessToken(userID, username, email)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err = m.GenerateRefreshToken(userID)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}
