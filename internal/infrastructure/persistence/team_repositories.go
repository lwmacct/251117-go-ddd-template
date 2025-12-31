package persistence

import (
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/organization"
	"gorm.io/gorm"
)

// TeamRepositories 聚合团队读写仓储
type TeamRepositories struct {
	Command organization.TeamCommandRepository
	Query   organization.TeamQueryRepository
}

// NewTeamRepositories 创建团队仓储聚合实例
func NewTeamRepositories(db *gorm.DB) TeamRepositories {
	return TeamRepositories{
		Command: NewTeamCommandRepository(db),
		Query:   NewTeamQueryRepository(db),
	}
}
