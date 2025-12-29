package role

import (
	"context"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/role"
)

// GetHandler 获取角色查询处理器
type GetHandler struct {
	roleQueryRepo role.QueryRepository
}

// NewGetHandler 创建获取角色查询处理器
func NewGetHandler(roleQueryRepo role.QueryRepository) *GetHandler {
	return &GetHandler{
		roleQueryRepo: roleQueryRepo,
	}
}

// Handle 处理获取角色查询
func (h *GetHandler) Handle(ctx context.Context, query GetQuery) (*RoleDTO, error) {
	// 查询角色（包含权限）
	roleEntity, err := h.roleQueryRepo.FindByIDWithPermissions(ctx, query.RoleID)
	if err != nil {
		return nil, fmt.Errorf("failed to find role: %w", err)
	}
	if roleEntity == nil {
		return nil, fmt.Errorf("role not found with id: %d", query.RoleID)
	}

	// 转换为 DTO
	return ToRoleDTO(roleEntity), nil
}
