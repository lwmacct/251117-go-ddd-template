// Package command 定义认证命令处理器
package command

import (
	"context"
	"fmt"

	domainAuth "github.com/lwmacct/251117-go-ddd-template/internal/domain/auth"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
)

// RegisterResult 注册结果
type RegisterResult struct {
	UserID       uint
	Username     string
	Email        string
	AccessToken  string
	RefreshToken string
	TokenType    string
	ExpiresIn    int
}

// RegisterHandler 注册命令处理器
type RegisterHandler struct {
	userCommandRepo user.CommandRepository
	userQueryRepo   user.QueryRepository
	authService     domainAuth.Service
}

// NewRegisterHandler 创建注册命令处理器
func NewRegisterHandler(
	userCommandRepo user.CommandRepository,
	userQueryRepo user.QueryRepository,
	authService domainAuth.Service,
) *RegisterHandler {
	return &RegisterHandler{
		userCommandRepo: userCommandRepo,
		userQueryRepo:   userQueryRepo,
		authService:     authService,
	}
}

// Handle 处理注册命令
func (h *RegisterHandler) Handle(ctx context.Context, cmd RegisterCommand) (*RegisterResult, error) {
	// 1. 验证密码策略
	if err := h.authService.ValidatePasswordPolicy(ctx, cmd.Password); err != nil {
		return nil, err
	}

	// 2. 检查用户名是否已存在
	exists, err := h.userQueryRepo.ExistsByUsername(ctx, cmd.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to check username existence: %w", err)
	}
	if exists {
		return nil, user.ErrUsernameAlreadyExists
	}

	// 3. 检查邮箱是否已存在
	exists, err = h.userQueryRepo.ExistsByEmail(ctx, cmd.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	}
	if exists {
		return nil, user.ErrEmailAlreadyExists
	}

	// 4. 生成密码哈希
	hashedPassword, err := h.authService.GeneratePasswordHash(ctx, cmd.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 5. 创建用户
	newUser := &user.User{
		Username: cmd.Username,
		Email:    cmd.Email,
		Password: hashedPassword,
		FullName: cmd.FullName,
		Status:   "active",
	}

	if err := h.userCommandRepo.Create(ctx, newUser); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// 6. 生成访问令牌
	accessToken, expiresAt, err := h.authService.GenerateAccessToken(ctx, newUser.ID, newUser.Username, []string{})
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, _, err := h.authService.GenerateRefreshToken(ctx, newUser.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &RegisterResult{
		UserID:       newUser.ID,
		Username:     newUser.Username,
		Email:        newUser.Email,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(expiresAt.Sub(expiresAt).Seconds()),
	}, nil
}
