package task

import (
	"context"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/task"
)

// GetHandler 获取任务详情处理器。
type GetHandler struct {
	queryRepo task.QueryRepository
}

// NewGetHandler 创建 GetHandler 实例。
func NewGetHandler(queryRepo task.QueryRepository) *GetHandler {
	return &GetHandler{
		queryRepo: queryRepo,
	}
}

// Handle 处理获取任务详情查询。
func (h *GetHandler) Handle(ctx context.Context, query GetTaskQuery) (*TaskDTO, error) {
	t, err := h.queryRepo.GetByIDAndTeam(ctx, query.ID, query.OrgID, query.TeamID)
	if err != nil {
		return nil, err
	}
	return ToTaskDTO(t), nil
}
