package organization

import (
	"context"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/organization"
)

// TeamDeleteHandler 删除团队命令处理器
type TeamDeleteHandler struct {
	teamCommand     organization.TeamCommandRepository
	teamQuery       organization.TeamQueryRepository
	teamMemberQuery organization.TeamMemberQueryRepository
}

// NewTeamDeleteHandler 创建删除团队命令处理器
func NewTeamDeleteHandler(
	teamCommand organization.TeamCommandRepository,
	teamQuery organization.TeamQueryRepository,
	teamMemberQuery organization.TeamMemberQueryRepository,
) *TeamDeleteHandler {
	return &TeamDeleteHandler{
		teamCommand:     teamCommand,
		teamQuery:       teamQuery,
		teamMemberQuery: teamMemberQuery,
	}
}

// Handle 处理删除团队命令
func (h *TeamDeleteHandler) Handle(ctx context.Context, cmd DeleteTeamCommand) error {
	// 1. 获取团队
	team, err := h.teamQuery.GetByID(ctx, cmd.TeamID)
	if err != nil {
		return err
	}

	// 2. 验证团队属于指定组织
	if !team.BelongsTo(cmd.OrgID) {
		return organization.ErrTeamNotInOrg
	}

	// 3. 检查是否还有成员
	memberCount, err := h.teamMemberQuery.CountByTeam(ctx, cmd.TeamID)
	if err != nil {
		return fmt.Errorf("failed to count team members: %w", err)
	}
	if memberCount > 0 {
		return organization.ErrTeamHasMembers
	}

	// 4. 删除团队
	if err := h.teamCommand.Delete(ctx, cmd.TeamID); err != nil {
		return fmt.Errorf("failed to delete team: %w", err)
	}

	return nil
}
