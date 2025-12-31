package organization

import (
	"context"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/organization"
)

// TeamGetHandler 获取团队查询处理器
type TeamGetHandler struct {
	teamQuery organization.TeamQueryRepository
}

// NewTeamGetHandler 创建获取团队查询处理器
func NewTeamGetHandler(teamQuery organization.TeamQueryRepository) *TeamGetHandler {
	return &TeamGetHandler{teamQuery: teamQuery}
}

// Handle 处理获取团队查询
func (h *TeamGetHandler) Handle(ctx context.Context, query GetTeamQuery) (*TeamDTO, error) {
	team, err := h.teamQuery.GetByID(ctx, query.TeamID)
	if err != nil {
		return nil, err
	}

	// 验证团队属于指定组织
	if !team.BelongsTo(query.OrgID) {
		return nil, organization.ErrTeamNotInOrg
	}

	return ToTeamDTO(team), nil
}
