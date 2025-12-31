package organization

import (
	"context"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/organization"
)

// MemberUpdateRoleHandler 更新成员角色命令处理器
type MemberUpdateRoleHandler struct {
	memberCommand organization.MemberCommandRepository
	memberQuery   organization.MemberQueryRepository
}

// NewMemberUpdateRoleHandler 创建更新成员角色命令处理器
func NewMemberUpdateRoleHandler(
	memberCommand organization.MemberCommandRepository,
	memberQuery organization.MemberQueryRepository,
) *MemberUpdateRoleHandler {
	return &MemberUpdateRoleHandler{
		memberCommand: memberCommand,
		memberQuery:   memberQuery,
	}
}

// Handle 处理更新成员角色命令
func (h *MemberUpdateRoleHandler) Handle(ctx context.Context, cmd UpdateMemberRoleCommand) error {
	// 1. 验证角色
	newRole := organization.MemberRole(cmd.Role)
	if !organization.IsValidMemberRole(newRole) {
		return organization.ErrInvalidMemberRole
	}

	// 2. 获取当前成员信息
	member, err := h.memberQuery.GetByOrgAndUser(ctx, cmd.OrgID, cmd.UserID)
	if err != nil {
		return err
	}

	// 3. 如果当前是 owner 且要降级，检查是否是最后一个
	if member.IsOwner() && newRole != organization.MemberRoleOwner {
		ownerCount, err := h.memberQuery.CountOwners(ctx, cmd.OrgID)
		if err != nil {
			return fmt.Errorf("failed to count owners: %w", err)
		}
		if ownerCount <= 1 {
			return organization.ErrCannotDemoteLastOwner
		}
	}

	// 4. 更新角色
	if err := h.memberCommand.UpdateRole(ctx, cmd.OrgID, cmd.UserID, newRole); err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}

	return nil
}
