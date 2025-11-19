// Package query 定义设置查询处理器
package query

// ListSettingsQuery 获取设置列表查询
type ListSettingsQuery struct {
	Category string // 可选: 按类别过滤
}
