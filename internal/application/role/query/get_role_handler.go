// Package query 定义角色查询处理器
package query

import (
	"context"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/application/role"
	domainRole "github.com/lwmacct/251117-go-ddd-template/internal/domain/role"
)

// GetRoleHandler 获取角色查询处理器
type GetRoleHandler struct {
	roleQueryRepo domainRole.QueryRepository
}

// NewGetRoleHandler 创建获取角色查询处理器
func NewGetRoleHandler(roleQueryRepo domainRole.QueryRepository) *GetRoleHandler {
	return &GetRoleHandler{
		roleQueryRepo: roleQueryRepo,
	}
}

// Handle 处理获取角色查询
func (h *GetRoleHandler) Handle(ctx context.Context, query GetRoleQuery) (*role.RoleResponse, error) {
	// 查询角色（包含权限）
	roleEntity, err := h.roleQueryRepo.FindByIDWithPermissions(ctx, query.RoleID)
	if err != nil {
		return nil, fmt.Errorf("failed to find role: %w", err)
	}
	if roleEntity == nil {
		return nil, fmt.Errorf("role not found with id: %d", query.RoleID)
	}

	// 转换为 DTO
	return role.ToRoleResponse(roleEntity), nil
}
