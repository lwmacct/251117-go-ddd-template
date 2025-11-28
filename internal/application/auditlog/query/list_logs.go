// Package query 定义审计日志的查询操作。
//
// 本包提供审计日志的查询和过滤功能：
//   - ListLogsQuery: 分页查询，支持多维度过滤
//   - GetLogQuery: 获取单条日志详情
//
// 过滤选项：
//   - UserID: 按用户筛选
//   - Action: 按操作类型筛选（如 create, update, delete）
//   - Resource: 按资源类型筛选（如 user, role）
//   - Status: 按操作状态筛选（success, failed）
//   - StartDate/EndDate: 时间范围筛选
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
