package user

import (
	"context"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
)

// UpdateUserHandler 更新用户命令处理器
type UpdateUserHandler struct {
	userCommandRepo user.CommandRepository
	userQueryRepo   user.QueryRepository
}

// NewUpdateUserHandler 创建更新用户命令处理器
func NewUpdateUserHandler(
	userCommandRepo user.CommandRepository,
	userQueryRepo user.QueryRepository,
) *UpdateUserHandler {
	return &UpdateUserHandler{
		userCommandRepo: userCommandRepo,
		userQueryRepo:   userQueryRepo,
	}
}

// Handle 处理更新用户命令
func (h *UpdateUserHandler) Handle(ctx context.Context, cmd UpdateUserCommand) (*UpdateUserResultDTO, error) {
	// 1. 获取用户
	u, err := h.userQueryRepo.GetByID(ctx, cmd.UserID)
	if err != nil {
		return nil, user.ErrUserNotFound
	}

	// 2. 更新用户属性
	if cmd.Username != nil && *cmd.Username != u.Username {
		// 检查用户名是否已存在
		exists, err := h.userQueryRepo.ExistsByUsername(ctx, *cmd.Username)
		if err != nil {
			return nil, fmt.Errorf("failed to check username existence: %w", err)
		}
		if exists {
			return nil, user.ErrUsernameAlreadyExists
		}
		u.Username = *cmd.Username
	}
	if cmd.Email != nil && *cmd.Email != u.Email {
		// 检查邮箱是否已存在
		exists, err := h.userQueryRepo.ExistsByEmail(ctx, *cmd.Email)
		if err != nil {
			return nil, fmt.Errorf("failed to check email existence: %w", err)
		}
		if exists {
			return nil, user.ErrEmailAlreadyExists
		}
		u.Email = *cmd.Email
	}
	if cmd.FullName != nil {
		u.FullName = *cmd.FullName
	}
	if cmd.Avatar != nil {
		u.Avatar = *cmd.Avatar
	}
	if cmd.Bio != nil {
		u.Bio = *cmd.Bio
	}
	if cmd.Status != nil {
		// 使用领域模型方法
		switch *cmd.Status {
		case "active":
			u.Activate()
		case "inactive":
			u.Deactivate()
		case "banned":
			u.Ban()
		default:
			return nil, user.ErrInvalidUserStatus
		}
	}

	// 3. 保存更新
	if err := h.userCommandRepo.Update(ctx, u); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &UpdateUserResultDTO{
		UserID: u.ID,
	}, nil
}
