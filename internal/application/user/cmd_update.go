package user

import (
	"context"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
)

// UpdateHandler 更新用户命令处理器
type UpdateHandler struct {
	userCommandRepo user.CommandRepository
	userQueryRepo   user.QueryRepository
}

// NewUpdateHandler 创建更新用户命令处理器
func NewUpdateHandler(
	userCommandRepo user.CommandRepository,
	userQueryRepo user.QueryRepository,
) *UpdateHandler {
	return &UpdateHandler{
		userCommandRepo: userCommandRepo,
		userQueryRepo:   userQueryRepo,
	}
}

// Handle 处理更新用户命令
func (h *UpdateHandler) Handle(ctx context.Context, cmd UpdateCommand) (*UpdateResultDTO, error) {
	// 1. 获取用户
	u, err := h.userQueryRepo.GetByID(ctx, cmd.UserID)
	if err != nil {
		return nil, user.ErrUserNotFound
	}

	// 2. 系统用户保护检查
	if cmd.Username != nil && *cmd.Username != u.Username {
		// 系统用户不可修改用户名
		if !u.CanModifyUsername() {
			return nil, user.ErrCannotModifySystemUsername
		}
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
		// root 用户状态不可修改
		if !u.CanModifyStatus() {
			return nil, user.ErrCannotModifyRootStatus
		}
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

	return &UpdateResultDTO{
		UserID: u.ID,
	}, nil
}
