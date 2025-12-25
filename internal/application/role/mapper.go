package role

import (
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/role"
)

// ToRoleDTO 将领域实体转换为 DTO
func ToRoleDTO(role *role.Role) *RoleDTO {
	if role == nil {
		return nil
	}

	response := &RoleDTO{
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
		response.Permissions = make([]*PermissionDTO, 0, len(role.Permissions))
		for i := range role.Permissions {
			response.Permissions = append(response.Permissions, ToPermissionDTO(&role.Permissions[i]))
		}
	}

	return response
}

// ToPermissionDTO 将权限实体转换为 DTO
func ToPermissionDTO(permission *role.Permission) *PermissionDTO {
	if permission == nil {
		return nil
	}

	return &PermissionDTO{
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
