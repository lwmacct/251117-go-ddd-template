package organization

import (
	"context"
	"fmt"
	"time"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/organization"
)

// MemberAddHandler 添加组织成员命令处理器
type MemberAddHandler struct {
	memberCommand organization.MemberCommandRepository
	memberQuery   organization.MemberQueryRepository
	orgQuery      organization.QueryRepository
}

// NewMemberAddHandler 创建添加成员命令处理器
func NewMemberAddHandler(
	memberCommand organization.MemberCommandRepository,
	memberQuery organization.MemberQueryRepository,
	orgQuery organization.QueryRepository,
) *MemberAddHandler {
	return &MemberAddHandler{
		memberCommand: memberCommand,
		memberQuery:   memberQuery,
		orgQuery:      orgQuery,
	}
}

// Handle 处理添加成员命令
func (h *MemberAddHandler) Handle(ctx context.Context, cmd AddMemberCommand) (*MemberDTO, error) {
	// 1. 检查组织是否存在
	exists, err := h.orgQuery.Exists(ctx, cmd.OrgID)
	if err != nil {
		return nil, fmt.Errorf("failed to check org existence: %w", err)
	}
	if !exists {
		return nil, organization.ErrOrgNotFound
	}

	// 2. 检查用户是否已是成员
	isMember, err := h.memberQuery.IsMember(ctx, cmd.OrgID, cmd.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to check membership: %w", err)
	}
	if isMember {
		return nil, organization.ErrMemberAlreadyExists
	}

	// 3. 验证角色
	role := organization.MemberRole(cmd.Role)
	if !organization.IsValidMemberRole(role) {
		return nil, organization.ErrInvalidMemberRole
	}

	// 4. 创建成员
	member := &organization.Member{
		OrganizationID: cmd.OrgID,
		UserID:         cmd.UserID,
		Role:           role,
		JoinedAt:       time.Now(),
	}

	if err := h.memberCommand.Add(ctx, member); err != nil {
		return nil, fmt.Errorf("failed to add member: %w", err)
	}

	return ToMemberDTO(member), nil
}
