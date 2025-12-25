package user

import (
	"context"

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
func (h *ListUsersHandler) Handle(ctx context.Context, query ListUsersQuery) (*UserListDTO, error) {
	// 根据是否有搜索关键词选择不同的查询方法
	users, total, err := h.fetchUsers(ctx, query)
	if err != nil {
		return nil, err
	}

	// 转换为 DTO
	userResponses := make([]*UserDTO, 0, len(users))
	for _, u := range users {
		userResponses = append(userResponses, &UserDTO{
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

	return &UserListDTO{
		Users: userResponses,
		Total: total,
	}, nil
}

// fetchUsers 根据查询条件获取用户列表和总数
func (h *ListUsersHandler) fetchUsers(ctx context.Context, query ListUsersQuery) ([]*user.User, int64, error) {
	if query.Search != "" {
		return h.searchUsers(ctx, query.Search, query.Offset, query.Limit)
	}
	return h.listAllUsers(ctx, query.Offset, query.Limit)
}

// searchUsers 搜索用户
func (h *ListUsersHandler) searchUsers(ctx context.Context, keyword string, offset, limit int) ([]*user.User, int64, error) {
	users, err := h.userQueryRepo.Search(ctx, keyword, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	total, err := h.userQueryRepo.CountBySearch(ctx, keyword)
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

// listAllUsers 获取所有用户列表
func (h *ListUsersHandler) listAllUsers(ctx context.Context, offset, limit int) ([]*user.User, int64, error) {
	users, err := h.userQueryRepo.List(ctx, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	total, err := h.userQueryRepo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}
