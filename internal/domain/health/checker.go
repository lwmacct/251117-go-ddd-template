// Package health 定义系统健康检查的领域模型
package health

import "context"

// Status 表示组件的健康状态
type Status string

const (
	StatusHealthy   Status = "healthy"
	StatusUnhealthy Status = "unhealthy"
	StatusDegraded  Status = "degraded"
)

// CheckResult 单个组件的检查结果
type CheckResult struct {
	Status Status         `json:"status"`
	Error  string         `json:"error,omitempty"`
	Stats  map[string]any `json:"stats,omitempty"`
}

// HealthReport 系统健康报告
type HealthReport struct {
	Status Status                 `json:"status"`
	Checks map[string]CheckResult `json:"checks"`
}

// Checker 定义健康检查服务接口
type Checker interface {
	// Check 执行健康检查并返回报告
	Check(ctx context.Context) *HealthReport
}
