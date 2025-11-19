// Package query 定义审计日志查询处理器
package query

import "time"

// ListLogsQuery 获取审计日志列表查询
type ListLogsQuery struct {
	Page      int
	Limit     int
	UserID    *uint
	Action    string
	Resource  string
	Status    string
	StartDate *time.Time
	EndDate   *time.Time
}
