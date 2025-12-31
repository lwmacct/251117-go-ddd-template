package org

import (
	"context"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/org"
)

// TeamDeleteHandler 删除团队命令处理器
type TeamDeleteHandler struct {
	teamCommand     org.TeamCommandRepository
	teamQuery       org.TeamQueryRepository
	teamMemberQuery org.TeamMemberQueryRepository
}

// NewTeamDeleteHandler 创建删除团队命令处理器
func NewTeamDeleteHandler(
	teamCommand org.TeamCommandRepository,
	teamQuery org.TeamQueryRepository,
	teamMemberQuery org.TeamMemberQueryRepository,
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
		return org.ErrTeamNotInOrg
	}

	// 3. 检查是否还有成员
	memberCount, err := h.teamMemberQuery.CountByTeam(ctx, cmd.TeamID)
	if err != nil {
		return fmt.Errorf("failed to count team members: %w", err)
	}
	if memberCount > 0 {
		return org.ErrTeamHasMembers
	}

	// 4. 删除团队
	if err := h.teamCommand.Delete(ctx, cmd.TeamID); err != nil {
		return fmt.Errorf("failed to delete team: %w", err)
	}

	return nil
}
