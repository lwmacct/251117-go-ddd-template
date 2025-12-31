package persistence

import (
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/organization"
	"gorm.io/gorm"
)

// OrgMemberRepositories 聚合组织成员读写仓储
type OrgMemberRepositories struct {
	Command organization.MemberCommandRepository
	Query   organization.MemberQueryRepository
}

// NewOrgMemberRepositories 创建组织成员仓储聚合实例
func NewOrgMemberRepositories(db *gorm.DB) OrgMemberRepositories {
	return OrgMemberRepositories{
		Command: NewOrgMemberCommandRepository(db),
		Query:   NewOrgMemberQueryRepository(db),
	}
}
