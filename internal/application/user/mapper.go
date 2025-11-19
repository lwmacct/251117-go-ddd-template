// Package user 提供领域模型到应用层 DTO 的映射函数
package user

import (
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
)

// ToUserResponse 将领域模型 User 转换为应用层 UserResponse DTO
func ToUserResponse(u *user.User) *UserResponse {
	if u == nil {
		return nil
	}

	return &UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		FullName:  u.FullName,
		Avatar:    u.Avatar,
		Bio:       u.Bio,
		Status:    u.Status,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// ToUserWithRolesResponse 将领域模型 User 转换为应用层 UserWithRolesResponse DTO（包含角色信息）
func ToUserWithRolesResponse(u *user.User) *UserWithRolesResponse {
	if u == nil {
		return nil
	}

	// 转换角色信息
	var roles []RoleDTO
	for _, role := range u.Roles {
		roles = append(roles, RoleDTO{
			ID:          role.ID,
			Name:        role.Name,
			DisplayName: role.DisplayName,
			Description: role.Description,
		})
	}

	return &UserWithRolesResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		FullName:  u.FullName,
		Avatar:    u.Avatar,
		Bio:       u.Bio,
		Status:    u.Status,
		Roles:     roles,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
