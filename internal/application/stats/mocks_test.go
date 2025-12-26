//nolint:forcetypeassert // Mock 返回值类型在测试中总是已知的
package stats

import (
	"github.com/stretchr/testify/mock"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/stats"
)

// ============================================================
// MockStatsQueryRepository
// ============================================================

// MockStatsQueryRepository 统计查询仓储 Mock
type MockStatsQueryRepository struct {
	mock.Mock
}

// GetSystemStats 模拟获取系统统计信息
func (m *MockStatsQueryRepository) GetSystemStats(recentLogsLimit int) (*stats.SystemStats, error) {
	args := m.Called(recentLogsLimit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*stats.SystemStats), args.Error(1)
}

// GetUserCountByStatus 模拟按状态统计用户数量
func (m *MockStatsQueryRepository) GetUserCountByStatus(status string) (int64, error) {
	args := m.Called(status)
	return args.Get(0).(int64), args.Error(1)
}

// GetTotalUsers 模拟获取用户总数
func (m *MockStatsQueryRepository) GetTotalUsers() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

// GetTotalRoles 模拟获取角色总数
func (m *MockStatsQueryRepository) GetTotalRoles() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

// GetTotalPermissions 模拟获取权限总数
func (m *MockStatsQueryRepository) GetTotalPermissions() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

// GetTotalMenus 模拟获取菜单总数
func (m *MockStatsQueryRepository) GetTotalMenus() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

// GetRecentAuditLogs 模拟获取最近的审计日志
func (m *MockStatsQueryRepository) GetRecentAuditLogs(limit int) ([]stats.AuditLogSummary, error) {
	args := m.Called(limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]stats.AuditLogSummary), args.Error(1)
}
