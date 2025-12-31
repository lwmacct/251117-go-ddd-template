package organization

import (
	"context"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/organization"
)

// TeamListHandler 团队列表查询处理器
type TeamListHandler struct {
	teamQuery organization.TeamQueryRepository
}

// NewTeamListHandler 创建团队列表查询处理器
func NewTeamListHandler(teamQuery organization.TeamQueryRepository) *TeamListHandler {
	return &TeamListHandler{teamQuery: teamQuery}
}

// TeamListResult 团队列表查询结果
type TeamListResult struct {
	Items []*TeamDTO
	Total int64
}

// Handle 处理团队列表查询
func (h *TeamListHandler) Handle(ctx context.Context, query ListTeamsQuery) (*TeamListResult, error) {
	teams, err := h.teamQuery.ListByOrg(ctx, query.OrgID, query.Offset, query.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list teams: %w", err)
	}

	total, err := h.teamQuery.CountByOrg(ctx, query.OrgID)
	if err != nil {
		return nil, fmt.Errorf("failed to count teams: %w", err)
	}

	return &TeamListResult{
		Items: ToTeamDTOs(teams),
		Total: total,
	}, nil
}
