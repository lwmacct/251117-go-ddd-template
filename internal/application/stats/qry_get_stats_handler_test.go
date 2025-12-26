//nolint:gosec // 测试中的整数转换是可控的
package stats

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/stats"
)

func TestGetStatsHandler_Handle_Success(t *testing.T) {
	// Arrange
	mockQueryRepo := new(MockStatsQueryRepository)

	now := time.Now()
	expectedStats := &stats.SystemStats{
		TotalUsers:       100,
		ActiveUsers:      80,
		InactiveUsers:    15,
		BannedUsers:      5,
		TotalRoles:       10,
		TotalPermissions: 50,
		TotalMenus:       20,
		RecentAuditLogs: []stats.AuditLogSummary{
			{
				ID:        1,
				UserID:    1,
				Username:  "admin",
				Action:    "login",
				Resource:  "auth",
				Status:    "success",
				CreatedAt: now,
			},
			{
				ID:        2,
				UserID:    2,
				Username:  "user1",
				Action:    "create",
				Resource:  "user",
				Status:    "success",
				CreatedAt: now.Add(-time.Hour),
			},
		},
	}

	mockQueryRepo.On("GetSystemStats", 5).Return(expectedStats, nil)

	handler := NewGetStatsHandler(mockQueryRepo)

	// Act
	result, err := handler.Handle(context.Background(), GetStatsQuery{
		RecentLogsLimit: 5,
	})

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(100), result.TotalUsers)
	assert.Equal(t, int64(80), result.ActiveUsers)
	assert.Equal(t, int64(15), result.InactiveUsers)
	assert.Equal(t, int64(5), result.BannedUsers)
	assert.Equal(t, int64(10), result.TotalRoles)
	assert.Equal(t, int64(50), result.TotalPermissions)
	assert.Equal(t, int64(20), result.TotalMenus)
	assert.Len(t, result.RecentAuditLogs, 2)
	assert.Equal(t, uint(1), result.RecentAuditLogs[0].ID)
	assert.Equal(t, "admin", result.RecentAuditLogs[0].Username)
	assert.Equal(t, "login", result.RecentAuditLogs[0].Action)

	mockQueryRepo.AssertExpectations(t)
}

func TestGetStatsHandler_Handle_DefaultLimit(t *testing.T) {
	// Arrange
	mockQueryRepo := new(MockStatsQueryRepository)

	expectedStats := &stats.SystemStats{
		TotalUsers:       50,
		ActiveUsers:      40,
		InactiveUsers:    8,
		BannedUsers:      2,
		TotalRoles:       5,
		TotalPermissions: 25,
		TotalMenus:       10,
		RecentAuditLogs:  []stats.AuditLogSummary{},
	}

	// 验证默认 limit 为 5
	mockQueryRepo.On("GetSystemStats", 5).Return(expectedStats, nil)

	handler := NewGetStatsHandler(mockQueryRepo)

	// Act - 未指定 limit 或 limit <= 0
	result, err := handler.Handle(context.Background(), GetStatsQuery{
		RecentLogsLimit: 0, // 应该使用默认值 5
	})

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(50), result.TotalUsers)

	mockQueryRepo.AssertExpectations(t)
}

func TestGetStatsHandler_Handle_CustomLimit(t *testing.T) {
	tests := []struct {
		name         string
		requestLimit int
		expectedCall int
	}{
		{
			name:         "使用默认 limit（输入为 0）",
			requestLimit: 0,
			expectedCall: 5,
		},
		{
			name:         "使用默认 limit（输入为负数）",
			requestLimit: -10,
			expectedCall: 5,
		},
		{
			name:         "自定义 limit 为 10",
			requestLimit: 10,
			expectedCall: 10,
		},
		{
			name:         "自定义 limit 为 20",
			requestLimit: 20,
			expectedCall: 20,
		},
		{
			name:         "自定义 limit 为 1",
			requestLimit: 1,
			expectedCall: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockQueryRepo := new(MockStatsQueryRepository)

			expectedStats := &stats.SystemStats{
				TotalUsers:       10,
				ActiveUsers:      8,
				InactiveUsers:    2,
				BannedUsers:      0,
				TotalRoles:       3,
				TotalPermissions: 15,
				TotalMenus:       5,
				RecentAuditLogs:  []stats.AuditLogSummary{},
			}

			mockQueryRepo.On("GetSystemStats", tt.expectedCall).Return(expectedStats, nil)

			handler := NewGetStatsHandler(mockQueryRepo)

			// Act
			result, err := handler.Handle(context.Background(), GetStatsQuery{
				RecentLogsLimit: tt.requestLimit,
			})

			// Assert
			require.NoError(t, err)
			assert.NotNil(t, result)
			mockQueryRepo.AssertExpectations(t)
		})
	}
}

func TestGetStatsHandler_Handle_EmptyAuditLogs(t *testing.T) {
	// Arrange
	mockQueryRepo := new(MockStatsQueryRepository)

	expectedStats := &stats.SystemStats{
		TotalUsers:       30,
		ActiveUsers:      25,
		InactiveUsers:    5,
		BannedUsers:      0,
		TotalRoles:       4,
		TotalPermissions: 20,
		TotalMenus:       8,
		RecentAuditLogs:  []stats.AuditLogSummary{}, // 空审计日志
	}

	mockQueryRepo.On("GetSystemStats", 5).Return(expectedStats, nil)

	handler := NewGetStatsHandler(mockQueryRepo)

	// Act
	result, err := handler.Handle(context.Background(), GetStatsQuery{
		RecentLogsLimit: 5,
	})

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result.RecentAuditLogs)

	mockQueryRepo.AssertExpectations(t)
}

func TestGetStatsHandler_Handle_Error(t *testing.T) {
	tests := []struct {
		name       string
		query      GetStatsQuery
		setupMocks func(*MockStatsQueryRepository)
		wantErr    string
	}{
		{
			name: "仓储查询失败",
			query: GetStatsQuery{
				RecentLogsLimit: 5,
			},
			setupMocks: func(repo *MockStatsQueryRepository) {
				repo.On("GetSystemStats", 5).
					Return(nil, errors.New("database connection failed"))
			},
			wantErr: "database connection failed",
		},
		{
			name: "数据库超时",
			query: GetStatsQuery{
				RecentLogsLimit: 10,
			},
			setupMocks: func(repo *MockStatsQueryRepository) {
				repo.On("GetSystemStats", 10).
					Return(nil, errors.New("query timeout"))
			},
			wantErr: "query timeout",
		},
		{
			name: "权限不足",
			query: GetStatsQuery{
				RecentLogsLimit: 5,
			},
			setupMocks: func(repo *MockStatsQueryRepository) {
				repo.On("GetSystemStats", 5).
					Return(nil, errors.New("permission denied"))
			},
			wantErr: "permission denied",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockQueryRepo := new(MockStatsQueryRepository)
			tt.setupMocks(mockQueryRepo)

			handler := NewGetStatsHandler(mockQueryRepo)

			// Act
			result, err := handler.Handle(context.Background(), tt.query)

			// Assert
			require.Error(t, err)
			assert.Nil(t, result)
			assert.Contains(t, err.Error(), tt.wantErr)
			mockQueryRepo.AssertExpectations(t)
		})
	}
}

func TestGetStatsHandler_Handle_ZeroValues(t *testing.T) {
	// Arrange - 测试全零值场景（新系统或数据被清空）
	mockQueryRepo := new(MockStatsQueryRepository)

	expectedStats := &stats.SystemStats{
		TotalUsers:       0,
		ActiveUsers:      0,
		InactiveUsers:    0,
		BannedUsers:      0,
		TotalRoles:       0,
		TotalPermissions: 0,
		TotalMenus:       0,
		RecentAuditLogs:  []stats.AuditLogSummary{},
	}

	mockQueryRepo.On("GetSystemStats", 5).Return(expectedStats, nil)

	handler := NewGetStatsHandler(mockQueryRepo)

	// Act
	result, err := handler.Handle(context.Background(), GetStatsQuery{
		RecentLogsLimit: 5,
	})

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(0), result.TotalUsers)
	assert.Equal(t, int64(0), result.ActiveUsers)
	assert.Equal(t, int64(0), result.TotalRoles)
	assert.Empty(t, result.RecentAuditLogs)

	mockQueryRepo.AssertExpectations(t)
}

func TestGetStatsHandler_Handle_LargeDataset(t *testing.T) {
	// Arrange - 测试大数据量场景
	mockQueryRepo := new(MockStatsQueryRepository)

	now := time.Now()
	auditLogs := make([]stats.AuditLogSummary, 100)
	for i := range 100 {
		auditLogs[i] = stats.AuditLogSummary{
			ID:        uint(i + 1),
			UserID:    uint((i % 10) + 1),
			Username:  "user" + string(rune(i)),
			Action:    "action",
			Resource:  "resource",
			Status:    "success",
			CreatedAt: now.Add(-time.Duration(i) * time.Minute),
		}
	}

	expectedStats := &stats.SystemStats{
		TotalUsers:       10000,
		ActiveUsers:      8500,
		InactiveUsers:    1200,
		BannedUsers:      300,
		TotalRoles:       100,
		TotalPermissions: 500,
		TotalMenus:       200,
		RecentAuditLogs:  auditLogs,
	}

	mockQueryRepo.On("GetSystemStats", 100).Return(expectedStats, nil)

	handler := NewGetStatsHandler(mockQueryRepo)

	// Act
	result, err := handler.Handle(context.Background(), GetStatsQuery{
		RecentLogsLimit: 100,
	})

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(10000), result.TotalUsers)
	assert.Len(t, result.RecentAuditLogs, 100)

	mockQueryRepo.AssertExpectations(t)
}

func TestGetStatsHandler_ToStatsDTO_NilAuditLogs(t *testing.T) {
	// Arrange
	mockQueryRepo := new(MockStatsQueryRepository)

	expectedStats := &stats.SystemStats{
		TotalUsers:       20,
		ActiveUsers:      15,
		InactiveUsers:    5,
		BannedUsers:      0,
		TotalRoles:       3,
		TotalPermissions: 10,
		TotalMenus:       5,
		RecentAuditLogs:  nil, // nil 而非空切片
	}

	mockQueryRepo.On("GetSystemStats", 5).Return(expectedStats, nil)

	handler := NewGetStatsHandler(mockQueryRepo)

	// Act
	result, err := handler.Handle(context.Background(), GetStatsQuery{
		RecentLogsLimit: 5,
	})

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Nil(t, result.RecentAuditLogs) // 应保持 nil

	mockQueryRepo.AssertExpectations(t)
}
