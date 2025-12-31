package persistence

import (
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/org"
	"gorm.io/gorm"
)

// OrganizationRepositories 聚合组织读写仓储
type OrganizationRepositories struct {
	Command org.CommandRepository
	Query   org.QueryRepository
}

// NewOrganizationRepositories 创建组织仓储聚合实例
func NewOrganizationRepositories(db *gorm.DB) OrganizationRepositories {
	return OrganizationRepositories{
		Command: NewOrganizationCommandRepository(db),
		Query:   NewOrganizationQueryRepository(db),
	}
}
