package organization

import (
	"context"
	"fmt"
	"time"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/organization"
)

// TeamCreateHandler 创建团队命令处理器
type TeamCreateHandler struct {
	teamCommand       organization.TeamCommandRepository
	teamQuery         organization.TeamQueryRepository
	orgQuery          organization.QueryRepository
	teamMemberCommand organization.TeamMemberCommandRepository
}

// NewTeamCreateHandler 创建团队命令处理器
func NewTeamCreateHandler(
	teamCommand organization.TeamCommandRepository,
	teamQuery organization.TeamQueryRepository,
	orgQuery organization.QueryRepository,
	teamMemberCommand organization.TeamMemberCommandRepository,
) *TeamCreateHandler {
	return &TeamCreateHandler{
		teamCommand:       teamCommand,
		teamQuery:         teamQuery,
		orgQuery:          orgQuery,
		teamMemberCommand: teamMemberCommand,
	}
}

// Handle 处理创建团队命令
func (h *TeamCreateHandler) Handle(ctx context.Context, cmd CreateTeamCommand) (*TeamDTO, error) {
	// 1. 检查组织是否存在
	exists, err := h.orgQuery.Exists(ctx, cmd.OrgID)
	if err != nil {
		return nil, fmt.Errorf("failed to check org existence: %w", err)
	}
	if !exists {
		return nil, organization.ErrOrgNotFound
	}

	// 2. 检查团队名称是否已存在（组织内唯一）
	exists, err = h.teamQuery.ExistsByOrgAndName(ctx, cmd.OrgID, cmd.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to check team name existence: %w", err)
	}
	if exists {
		return nil, organization.ErrTeamNameAlreadyExists
	}

	// 3. 创建团队实体
	team := &organization.Team{
		OrganizationID: cmd.OrgID,
		Name:           cmd.Name,
		DisplayName:    cmd.DisplayName,
		Description:    cmd.Description,
	}

	// 4. 保存团队
	if err := h.teamCommand.Create(ctx, team); err != nil {
		return nil, fmt.Errorf("failed to create team: %w", err)
	}

	// 5. 添加团队负责人（如果提供）
	if cmd.LeadUserID > 0 {
		member := &organization.TeamMember{
			TeamID:   team.ID,
			UserID:   cmd.LeadUserID,
			Role:     organization.TeamMemberRoleLead,
			JoinedAt: time.Now(),
		}
		if err := h.teamMemberCommand.Add(ctx, member); err != nil {
			return nil, fmt.Errorf("failed to add team lead: %w", err)
		}
	}

	return ToTeamDTO(team), nil
}
