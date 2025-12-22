package role

import (
	domainRole "github.com/lwmacct/251117-go-ddd-template/internal/domain/role"
)

// ToRoleResponse 将领域实体转换为 DTO
func ToRoleResponse(role *domainRole.Role) *RoleResponse {
	if role == nil {
		return nil
	}

	response := &RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		DisplayName: role.DisplayName,
		Description: role.Description,
		IsSystem:    role.IsSystem,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}

	// 转换权限列表
	if len(role.Permissions) > 0 {
		response.Permissions = make([]*PermissionResponse, 0, len(role.Permissions))
		for i := range role.Permissions {
			response.Permissions = append(response.Permissions, ToPermissionResponse(&role.Permissions[i]))
		}
	}

	return response
}

// ToPermissionResponse 将权限实体转换为 DTO
func ToPermissionResponse(permission *domainRole.Permission) *PermissionResponse {
	if permission == nil {
		return nil
	}

	return &PermissionResponse{
		ID:          permission.ID,
		Code:        permission.Code,
		Name:        permission.Code, // 使用 Code 作为 Name
		Description: permission.Description,
		Resource:    permission.Resource,
		Action:      permission.Action,
		CreatedAt:   permission.CreatedAt,
		UpdatedAt:   permission.UpdatedAt,
	}
}
