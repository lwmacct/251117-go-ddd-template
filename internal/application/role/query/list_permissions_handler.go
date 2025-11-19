// Package query 定义角色查询处理器
package query

import (
	"context"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/application/role"
	domainRole "github.com/lwmacct/251117-go-ddd-template/internal/domain/role"
)

// ListPermissionsHandler 列出权限查询处理器
type ListPermissionsHandler struct {
	permissionQueryRepo domainRole.PermissionQueryRepository
}

// NewListPermissionsHandler 创建列出权限查询处理器
func NewListPermissionsHandler(permissionQueryRepo domainRole.PermissionQueryRepository) *ListPermissionsHandler {
	return &ListPermissionsHandler{
		permissionQueryRepo: permissionQueryRepo,
	}
}

// Handle 处理列出权限查询
func (h *ListPermissionsHandler) Handle(ctx context.Context, query ListPermissionsQuery) (*role.ListPermissionsResponse, error) {
	// 查询权限列表
	permissions, total, err := h.permissionQueryRepo.List(ctx, query.Page, query.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list permissions: %w", err)
	}

	// 转换为 DTO
	permissionResponses := make([]*role.PermissionResponse, 0, len(permissions))
	for i := range permissions {
		permissionResponses = append(permissionResponses, role.ToPermissionResponse(&permissions[i]))
	}

	return &role.ListPermissionsResponse{
		Permissions: permissionResponses,
		Total:       total,
		Page:        query.Page,
		Limit:       query.Limit,
	}, nil
}
