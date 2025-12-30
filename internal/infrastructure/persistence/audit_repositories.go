package persistence

import (
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/audit"
	"gorm.io/gorm"
)

// AuditLogRepositories 聚合审计日志读写仓储，便于同时注入 Command/Query
type AuditLogRepositories struct {
	Command audit.CommandRepository
	Query   audit.QueryRepository
}

// NewAuditLogRepositories 初始化审计日志仓储聚合
func NewAuditLogRepositories(db *gorm.DB) AuditLogRepositories {
	return AuditLogRepositories{
		Command: NewAuditLogCommandRepository(db),
		Query:   NewAuditLogQueryRepository(db),
	}
}
