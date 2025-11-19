package command

import (
	"context"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
)

// AssignRolesHandler 负责分配用户角色
type AssignRolesHandler struct {
	userCommandRepo user.CommandRepository
	userQueryRepo   user.QueryRepository
}

// NewAssignRolesHandler 创建新的分配角色处理器
func NewAssignRolesHandler(
	userCommandRepo user.CommandRepository,
	userQueryRepo user.QueryRepository,
) *AssignRolesHandler {
	return &AssignRolesHandler{
		userCommandRepo: userCommandRepo,
		userQueryRepo:   userQueryRepo,
	}
}

// Handle 处理分配角色命令
func (h *AssignRolesHandler) Handle(ctx context.Context, cmd AssignRolesCommand) error {
	if _, err := h.userQueryRepo.GetByID(ctx, cmd.UserID); err != nil {
		return err
	}

	return h.userCommandRepo.AssignRoles(ctx, cmd.UserID, cmd.RoleIDs)
}
