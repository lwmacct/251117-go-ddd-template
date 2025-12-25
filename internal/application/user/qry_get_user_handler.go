package user

import (
	"context"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
)

// GetUserHandler 获取用户查询处理器
type GetUserHandler struct {
	userQueryRepo user.QueryRepository
}

// NewGetUserHandler 创建获取用户查询处理器
func NewGetUserHandler(userQueryRepo user.QueryRepository) *GetUserHandler {
	return &GetUserHandler{
		userQueryRepo: userQueryRepo,
	}
}

// Handle 处理获取用户查询
func (h *GetUserHandler) Handle(ctx context.Context, query GetUserQuery) (*UserWithRolesDTO, error) {
	var u *user.User
	var err error

	if query.WithRoles {
		u, err = h.userQueryRepo.GetByIDWithRoles(ctx, query.UserID)
	} else {
		u, err = h.userQueryRepo.GetByID(ctx, query.UserID)
	}

	if err != nil {
		return nil, err
	}

	// 转换为 DTO
	response := &UserWithRolesDTO{
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

	if query.WithRoles {
		roles := make([]RoleDTO, 0, len(u.Roles))
		for _, r := range u.Roles {
			roles = append(roles, RoleDTO{
				ID:          r.ID,
				Name:        r.Name,
				DisplayName: r.DisplayName,
				Description: r.Description,
			})
		}
		response.Roles = roles
	}

	return response, nil
}
