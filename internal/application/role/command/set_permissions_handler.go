package command

import (
	"context"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/role"
)

// SetPermissionsHandler 设置权限命令处理器
type SetPermissionsHandler struct {
	roleCommandRepo     role.CommandRepository
	roleQueryRepo       role.QueryRepository
	permissionQueryRepo role.PermissionQueryRepository
}

// NewSetPermissionsHandler 创建设置权限命令处理器
func NewSetPermissionsHandler(
	roleCommandRepo role.CommandRepository,
	roleQueryRepo role.QueryRepository,
	permissionQueryRepo role.PermissionQueryRepository,
) *SetPermissionsHandler {
	return &SetPermissionsHandler{
		roleCommandRepo:     roleCommandRepo,
		roleQueryRepo:       roleQueryRepo,
		permissionQueryRepo: permissionQueryRepo,
	}
}

// Handle 处理设置权限命令
func (h *SetPermissionsHandler) Handle(ctx context.Context, cmd SetPermissionsCommand) error {
	// 1. 验证角色是否存在
	exists, err := h.roleQueryRepo.Exists(ctx, cmd.RoleID)
	if err != nil {
		return fmt.Errorf("failed to check role existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("role not found with id: %d", cmd.RoleID)
	}

	// 2. 验证所有权限ID是否有效
	for _, permID := range cmd.PermissionIDs {
		permExists, err := h.permissionQueryRepo.Exists(ctx, permID)
		if err != nil {
			return fmt.Errorf("failed to check permission existence: %w", err)
		}
		if !permExists {
			return fmt.Errorf("permission not found with id: %d", permID)
		}
	}

	// 3. 设置权限
	if err := h.roleCommandRepo.SetPermissions(ctx, cmd.RoleID, cmd.PermissionIDs); err != nil {
		return fmt.Errorf("failed to set permissions: %w", err)
	}

	return nil
}
