package persistence

import (
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/organization"
	"gorm.io/gorm"
)

// TeamMemberRepositories 聚合团队成员读写仓储
type TeamMemberRepositories struct {
	Command organization.TeamMemberCommandRepository
	Query   organization.TeamMemberQueryRepository
}

// NewTeamMemberRepositories 创建团队成员仓储聚合实例
func NewTeamMemberRepositories(db *gorm.DB) TeamMemberRepositories {
	return TeamMemberRepositories{
		Command: NewTeamMemberCommandRepository(db),
		Query:   NewTeamMemberQueryRepository(db),
	}
}
