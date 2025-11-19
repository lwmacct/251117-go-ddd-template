// Package auth 提供认证服务
package auth

import (
	"context"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
	"golang.org/x/crypto/bcrypt"
)

// Service 认证服务
type Service struct {
	userRepo   user.Repository
	jwtManager *JWTManager
}

// NewService 创建认证服务
func NewService(userRepo user.Repository, jwtManager *JWTManager) *Service {
	return &Service{
		userRepo:   userRepo,
		jwtManager: jwtManager,
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
type LoginRequest struct {
	Login    string `json:"login" binding:"required"` // 用户名或邮箱
	Password string `json:"password" binding:"required"`
}

// AuthResponse 认证响应
type AuthResponse struct {
	AccessToken  string             `json:"access_token"`
	RefreshToken string             `json:"refresh_token"`
	TokenType    string             `json:"token_type"`
	ExpiresIn    int                `json:"expires_in"` // 秒
	User         *user.UserResponse `json:"user"`
}

// Register 注册新用户
func (s *Service) Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error) {
	// 检查用户名是否已存在
	if _, err := s.userRepo.GetByUsername(ctx, req.Username); err == nil {
		return nil, fmt.Errorf("username already exists")
	}

	// 检查邮箱是否已存在
	if _, err := s.userRepo.GetByEmail(ctx, req.Email); err == nil {
		return nil, fmt.Errorf("email already exists")
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

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Reload user with roles
	userWithRoles, err := s.userRepo.GetByIDWithRoles(ctx, newUser.ID)
	if err != nil {
		// Fallback to empty roles if failed to load
		userWithRoles = newUser
	}

	// 生成 token
	roles := userWithRoles.GetRoleNames()
	permissions := userWithRoles.GetPermissionCodes()
	accessToken, refreshToken, err := s.jwtManager.GenerateTokenPair(userWithRoles.ID, userWithRoles.Username, userWithRoles.Email, roles, permissions)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(s.jwtManager.accessTokenDuration.Seconds()),
		User:         userWithRoles.ToResponse(),
	}, nil
}

// Login 用户登录
func (s *Service) Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error) {
	// 尝试通过用户名或邮箱查找用户（包含角色和权限）
	var u *user.User
	var err error

	// 先尝试用户名
	u, err = s.userRepo.GetByUsernameWithRoles(ctx, req.Login)
	if err != nil {
		// 再尝试邮箱
		tempUser, err2 := s.userRepo.GetByEmail(ctx, req.Login)
		if err2 != nil {
			return nil, fmt.Errorf("invalid credentials")
		}
		// Reload with roles
		u, err = s.userRepo.GetByIDWithRoles(ctx, tempUser.ID)
		if err != nil {
			u = tempUser
		}
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// 检查用户状态
	if u.Status != "active" {
		return nil, fmt.Errorf("user account is %s", u.Status)
	}

	// 生成 token
	roles := u.GetRoleNames()
	permissions := u.GetPermissionCodes()
	accessToken, refreshToken, err := s.jwtManager.GenerateTokenPair(u.ID, u.Username, u.Email, roles, permissions)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(s.jwtManager.accessTokenDuration.Seconds()),
		User:         u.ToResponse(),
	}, nil
}

// RefreshToken 刷新访问令牌
func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (*AuthResponse, error) {
	// 验证刷新令牌
	claims, err := s.jwtManager.ValidateToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// 获取用户信息（包含角色和权限）
	u, err := s.userRepo.GetByIDWithRoles(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// 检查用户状态
	if u.Status != "active" {
		return nil, fmt.Errorf("user account is %s", u.Status)
	}

	// 生成新的 token 对
	roles := u.GetRoleNames()
	permissions := u.GetPermissionCodes()
	accessToken, newRefreshToken, err := s.jwtManager.GenerateTokenPair(u.ID, u.Username, u.Email, roles, permissions)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(s.jwtManager.accessTokenDuration.Seconds()),
		User:         u.ToResponse(),
	}, nil
}
