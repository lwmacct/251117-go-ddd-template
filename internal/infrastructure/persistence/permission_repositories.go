package persistence

import (
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/role"
	"gorm.io/gorm"
)

// PermissionRepositories 聚合权限读写仓储
type PermissionRepositories struct {
	Command role.PermissionCommandRepository
	Query   role.PermissionQueryRepository
}

// NewPermissionRepositories 创建权限仓储聚合实例
func NewPermissionRepositories(db *gorm.DB) PermissionRepositories {
	return PermissionRepositories{
		Command: NewPermissionCommandRepository(db),
		Query:   NewPermissionQueryRepository(db),
	}
}
