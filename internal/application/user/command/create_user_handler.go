// Package command 定义用户命令处理器
package command

import (
	"context"
	"fmt"

	domainAuth "github.com/lwmacct/251117-go-ddd-template/internal/domain/auth"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
)

// CreateUserResult 创建用户结果
type CreateUserResult struct {
	UserID   uint
	Username string
	Email    string
}

// CreateUserHandler 创建用户命令处理器
type CreateUserHandler struct {
	userCommandRepo user.CommandRepository
	userQueryRepo   user.QueryRepository
	authService     domainAuth.Service
}

// NewCreateUserHandler 创建用户命令处理器
func NewCreateUserHandler(
	userCommandRepo user.CommandRepository,
	userQueryRepo user.QueryRepository,
	authService domainAuth.Service,
) *CreateUserHandler {
	return &CreateUserHandler{
		userCommandRepo: userCommandRepo,
		userQueryRepo:   userQueryRepo,
		authService:     authService,
	}
}

// Handle 处理创建用户命令
func (h *CreateUserHandler) Handle(ctx context.Context, cmd CreateUserCommand) (*CreateUserResult, error) {
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

	// 5. 创建用户实体
	newUser := &user.User{
		Username: cmd.Username,
		Email:    cmd.Email,
		Password: hashedPassword,
		FullName: cmd.FullName,
		Status:   "active",
	}

	// 6. 保存用户
	if err := h.userCommandRepo.Create(ctx, newUser); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// 7. 分配角色（如果提供）
	if len(cmd.RoleIDs) > 0 {
		if err := h.userCommandRepo.AssignRoles(ctx, newUser.ID, cmd.RoleIDs); err != nil {
			return nil, fmt.Errorf("failed to assign roles: %w", err)
		}
	}

	return &CreateUserResult{
		UserID:   newUser.ID,
		Username: newUser.Username,
		Email:    newUser.Email,
	}, nil
}
