// Package query 定义用户查询处理器
package query

import (
	"context"

	appUserDTO "github.com/lwmacct/251117-go-ddd-template/internal/application/user"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
)

// ListUsersHandler 获取用户列表查询处理器
type ListUsersHandler struct {
	userQueryRepo user.QueryRepository
}

// NewListUsersHandler 创建获取用户列表查询处理器
func NewListUsersHandler(userQueryRepo user.QueryRepository) *ListUsersHandler {
	return &ListUsersHandler{
		userQueryRepo: userQueryRepo,
	}
}

// Handle 处理获取用户列表查询
func (h *ListUsersHandler) Handle(ctx context.Context, query ListUsersQuery) (*appUserDTO.UserListResponse, error) {
	// 获取用户列表
	users, err := h.userQueryRepo.List(ctx, query.Offset, query.Limit)
	if err != nil {
		return nil, err
	}

	// 获取总数
	total, err := h.userQueryRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	// 转换为 DTO
	userResponses := make([]*appUserDTO.UserResponse, 0, len(users))
	for _, u := range users {
		userResponses = append(userResponses, &appUserDTO.UserResponse{
			ID:        u.ID,
			Username:  u.Username,
			Email:     u.Email,
			FullName:  u.FullName,
			Avatar:    u.Avatar,
			Bio:       u.Bio,
			Status:    u.Status,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		})
	}

	return &appUserDTO.UserListResponse{
		Users: userResponses,
		Total: total,
	}, nil
}
