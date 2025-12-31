package persistence

import (
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/audit"
	"gorm.io/gorm"
)

// AuditRepositories 聚合审计日志读写仓储，便于同时注入 Command/Query
type AuditRepositories struct {
	Command audit.CommandRepository
	Query   audit.QueryRepository
}

// NewAuditRepositories 初始化审计日志仓储聚合
func NewAuditRepositories(db *gorm.DB) AuditRepositories {
	return AuditRepositories{
		Command: NewAuditCommandRepository(db),
		Query:   NewAuditQueryRepository(db),
	}
}
