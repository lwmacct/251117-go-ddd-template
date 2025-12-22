package persistence

import (
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/menu"
	"gorm.io/gorm"
)

// MenuRepositories 聚合菜单读写仓储
type MenuRepositories struct {
	Command menu.CommandRepository
	Query   menu.QueryRepository
}

// NewMenuRepositories 创建菜单仓储聚合实例
func NewMenuRepositories(db *gorm.DB) MenuRepositories {
	return MenuRepositories{
		Command: NewMenuCommandRepository(db),
		Query:   NewMenuQueryRepository(db),
	}
}
