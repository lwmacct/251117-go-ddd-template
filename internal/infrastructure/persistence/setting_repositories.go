package persistence

import (
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
	"gorm.io/gorm"
)

// SettingRepositories 聚合系统配置读写仓储
type SettingRepositories struct {
	Command setting.CommandRepository
	Query   setting.QueryRepository
}

// NewSettingRepositories 创建配置仓储聚合实例
func NewSettingRepositories(db *gorm.DB) SettingRepositories {
	return SettingRepositories{
		Command: NewSettingCommandRepository(db),
		Query:   NewSettingQueryRepository(db),
	}
}
