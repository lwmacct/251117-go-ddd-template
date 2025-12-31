package org

import (
	"context"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/org"
)

// DeleteHandler 删除组织命令处理器
type DeleteHandler struct {
	orgCommandRepo org.CommandRepository
	orgQueryRepo   org.QueryRepository
	memberQuery    org.MemberQueryRepository
	teamQuery      org.TeamQueryRepository
}

// NewDeleteHandler 创建删除组织命令处理器
func NewDeleteHandler(
	orgCommandRepo org.CommandRepository,
	orgQueryRepo org.QueryRepository,
	memberQuery org.MemberQueryRepository,
	teamQuery org.TeamQueryRepository,
) *DeleteHandler {
	return &DeleteHandler{
		orgCommandRepo: orgCommandRepo,
		orgQueryRepo:   orgQueryRepo,
		memberQuery:    memberQuery,
		teamQuery:      teamQuery,
	}
}

// Handle 处理删除组织命令
func (h *DeleteHandler) Handle(ctx context.Context, cmd DeleteOrgCommand) error {
	// 1. 检查组织是否存在
	exists, err := h.orgQueryRepo.Exists(ctx, cmd.OrgID)
	if err != nil {
		return fmt.Errorf("failed to check org existence: %w", err)
	}
	if !exists {
		return org.ErrOrgNotFound
	}

	// 2. 检查是否还有成员
	memberCount, err := h.memberQuery.CountByOrg(ctx, cmd.OrgID)
	if err != nil {
		return fmt.Errorf("failed to count members: %w", err)
	}
	if memberCount > 0 {
		return org.ErrOrgHasMembers
	}

	// 3. 检查是否还有团队
	teamCount, err := h.teamQuery.CountByOrg(ctx, cmd.OrgID)
	if err != nil {
		return fmt.Errorf("failed to count teams: %w", err)
	}
	if teamCount > 0 {
		return org.ErrOrgHasTeams
	}

	// 4. 删除组织
	if err := h.orgCommandRepo.Delete(ctx, cmd.OrgID); err != nil {
		return fmt.Errorf("failed to delete organization: %w", err)
	}

	return nil
}
