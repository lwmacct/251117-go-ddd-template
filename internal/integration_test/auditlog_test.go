// Package integration_test 提供审计日志流程集成测试
package integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	auditlogQuery "github.com/lwmacct/251117-go-ddd-template/internal/application/auditlog/query"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/auditlog"
)

func TestAuditLogFlow_CreateAndQuery(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	ctx := context.Background()

	t.Run("完整审计日志流程: 创建 → 查询 → 过滤", func(t *testing.T) {
		// 步骤 1: 创建测试用户
		testUser := env.CreateTestUser(ctx, t, "auditlog_user", "auditlog@example.com", "Password123!")

		// 步骤 2: 记录审计日志
		log1 := &auditlog.AuditLog{
			UserID:     testUser.ID,
			Username:   testUser.Username,
			Action:     auditlog.ActionLogin,
			Resource:   "auth",
			ResourceID: "",
			IPAddress:  "192.168.1.1",
			UserAgent:  "Mozilla/5.0 Test",
			Details:    "用户登录成功",
			Status:     auditlog.StatusSuccess,
		}

		err := env.AuditLogCommandRepo.Create(ctx, log1)
		require.NoError(t, err)
		assert.NotZero(t, log1.ID)

		// 步骤 3: 记录另一条审计日志
		log2 := &auditlog.AuditLog{
			UserID:     testUser.ID,
			Username:   testUser.Username,
			Action:     auditlog.ActionUpdate,
			Resource:   "user",
			ResourceID: "profile",
			IPAddress:  "192.168.1.1",
			UserAgent:  "Mozilla/5.0 Test",
			Details:    "用户更新个人资料",
			Status:     auditlog.StatusSuccess,
		}

		err = env.AuditLogCommandRepo.Create(ctx, log2)
		require.NoError(t, err)

		// 步骤 4: 查询审计日志列表
		result, err := env.ListLogsHandler.Handle(ctx, auditlogQuery.ListLogsQuery{
			Page:  1,
			Limit: 10,
		})

		require.NoError(t, err)
		assert.GreaterOrEqual(t, result.Total, int64(2))
		assert.GreaterOrEqual(t, len(result.Logs), 2)

		// 步骤 5: 按用户过滤
		userID := testUser.ID
		userResult, err := env.ListLogsHandler.Handle(ctx, auditlogQuery.ListLogsQuery{
			Page:   1,
			Limit:  10,
			UserID: &userID,
		})

		require.NoError(t, err)
		assert.GreaterOrEqual(t, userResult.Total, int64(2))
		for _, log := range userResult.Logs {
			assert.Equal(t, testUser.ID, log.UserID)
		}
	})
}

func TestAuditLogFlow_FilterByAction(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	ctx := context.Background()

	t.Run("按操作类型过滤审计日志", func(t *testing.T) {
		testUser := env.CreateTestUser(ctx, t, "action_filter_user", "action@example.com", "Password123!")

		// 创建不同操作类型的日志
		actions := []string{auditlog.ActionCreate, auditlog.ActionUpdate, auditlog.ActionDelete, auditlog.ActionLogin}
		for _, action := range actions {
			log := &auditlog.AuditLog{
				UserID:   testUser.ID,
				Username: testUser.Username,
				Action:   action,
				Resource: "test_resource",
				Status:   auditlog.StatusSuccess,
			}
			err := env.AuditLogCommandRepo.Create(ctx, log)
			require.NoError(t, err)
		}

		// 按 login 操作过滤
		result, err := env.ListLogsHandler.Handle(ctx, auditlogQuery.ListLogsQuery{
			Page:   1,
			Limit:  10,
			Action: auditlog.ActionLogin,
		})

		require.NoError(t, err)
		assert.GreaterOrEqual(t, result.Total, int64(1))
		for _, log := range result.Logs {
			assert.Equal(t, auditlog.ActionLogin, log.Action)
		}
	})
}

func TestAuditLogFlow_FilterByResource(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	ctx := context.Background()

	t.Run("按资源类型过滤审计日志", func(t *testing.T) {
		testUser := env.CreateTestUser(ctx, t, "resource_filter_user", "resource@example.com", "Password123!")

		// 创建不同资源类型的日志
		resources := []string{"user", "role", "menu", "setting"}
		for _, resource := range resources {
			log := &auditlog.AuditLog{
				UserID:   testUser.ID,
				Username: testUser.Username,
				Action:   auditlog.ActionUpdate,
				Resource: resource,
				Status:   auditlog.StatusSuccess,
			}
			err := env.AuditLogCommandRepo.Create(ctx, log)
			require.NoError(t, err)
		}

		// 按 role 资源过滤
		result, err := env.ListLogsHandler.Handle(ctx, auditlogQuery.ListLogsQuery{
			Page:     1,
			Limit:    10,
			Resource: "role",
		})

		require.NoError(t, err)
		assert.GreaterOrEqual(t, result.Total, int64(1))
		for _, log := range result.Logs {
			assert.Equal(t, "role", log.Resource)
		}
	})
}

func TestAuditLogFlow_FilterByStatus(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	ctx := context.Background()

	t.Run("按状态过滤审计日志", func(t *testing.T) {
		testUser := env.CreateTestUser(ctx, t, "status_filter_user", "status@example.com", "Password123!")

		// 使用唯一资源名以隔离测试
		uniqueResource := "status_test_resource"

		// 创建成功和失败的日志
		successLog := &auditlog.AuditLog{
			UserID:   testUser.ID,
			Username: testUser.Username,
			Action:   auditlog.ActionLogin,
			Resource: uniqueResource,
			Status:   auditlog.StatusSuccess,
			Details:  "登录成功",
		}
		err := env.AuditLogCommandRepo.Create(ctx, successLog)
		require.NoError(t, err)

		failedLog := &auditlog.AuditLog{
			UserID:   testUser.ID,
			Username: testUser.Username,
			Action:   auditlog.ActionLogin,
			Resource: uniqueResource,
			Status:   auditlog.StatusFailed,
			Details:  "密码错误",
		}
		err = env.AuditLogCommandRepo.Create(ctx, failedLog)
		require.NoError(t, err)

		// 按失败状态过滤（同时按资源过滤以隔离）
		result, err := env.ListLogsHandler.Handle(ctx, auditlogQuery.ListLogsQuery{
			Page:     1,
			Limit:    10,
			Status:   auditlog.StatusFailed,
			Resource: uniqueResource,
		})

		require.NoError(t, err)
		// 验证返回的所有日志都是 failed 状态
		assert.GreaterOrEqual(t, result.Total, int64(1), "should have at least one failed log")
		for _, log := range result.Logs {
			assert.Equal(t, auditlog.StatusFailed, log.Status, "all returned logs should have failed status")
		}
	})
}

func TestAuditLogFlow_FilterByDateRange(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	ctx := context.Background()

	t.Run("按日期范围过滤审计日志", func(t *testing.T) {
		testUser := env.CreateTestUser(ctx, t, "date_filter_user", "date@example.com", "Password123!")

		// 创建日志
		log := &auditlog.AuditLog{
			UserID:   testUser.ID,
			Username: testUser.Username,
			Action:   auditlog.ActionCreate,
			Resource: "test",
			Status:   auditlog.StatusSuccess,
		}
		err := env.AuditLogCommandRepo.Create(ctx, log)
		require.NoError(t, err)

		// 按今天的日期范围过滤
		now := time.Now()
		startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		endOfDay := startOfDay.Add(24 * time.Hour)

		result, err := env.ListLogsHandler.Handle(ctx, auditlogQuery.ListLogsQuery{
			Page:      1,
			Limit:     10,
			StartDate: &startOfDay,
			EndDate:   &endOfDay,
		})

		require.NoError(t, err)
		assert.GreaterOrEqual(t, result.Total, int64(1))
	})
}

func TestAuditLogFlow_BatchCreate(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	ctx := context.Background()

	t.Run("批量创建审计日志", func(t *testing.T) {
		testUser := env.CreateTestUser(ctx, t, "batch_log_user", "batch@example.com", "Password123!")

		// 批量创建日志
		logs := make([]*auditlog.AuditLog, 5)
		for i := range logs {
			logs[i] = &auditlog.AuditLog{
				UserID:   testUser.ID,
				Username: testUser.Username,
				Action:   auditlog.ActionUpdate,
				Resource: "batch_test",
				Status:   auditlog.StatusSuccess,
				Details:  "批量测试日志",
			}
		}

		err := env.AuditLogCommandRepo.BatchCreate(ctx, logs)
		require.NoError(t, err)

		// 验证所有日志都已创建
		for _, log := range logs {
			assert.NotZero(t, log.ID)
		}

		// 查询验证
		result, err := env.ListLogsHandler.Handle(ctx, auditlogQuery.ListLogsQuery{
			Page:     1,
			Limit:    10,
			Resource: "batch_test",
		})

		require.NoError(t, err)
		assert.GreaterOrEqual(t, result.Total, int64(5))
	})
}

func TestAuditLogFlow_Pagination(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	ctx := context.Background()

	t.Run("审计日志分页", func(t *testing.T) {
		testUser := env.CreateTestUser(ctx, t, "pagination_user", "pagination@example.com", "Password123!")

		// 创建 15 条日志
		for range 15 {
			log := &auditlog.AuditLog{
				UserID:   testUser.ID,
				Username: testUser.Username,
				Action:   auditlog.ActionUpdate,
				Resource: "pagination_test",
				Status:   auditlog.StatusSuccess,
			}
			err := env.AuditLogCommandRepo.Create(ctx, log)
			require.NoError(t, err)
		}

		// 第一页
		page1, err := env.ListLogsHandler.Handle(ctx, auditlogQuery.ListLogsQuery{
			Page:     1,
			Limit:    5,
			Resource: "pagination_test",
		})
		require.NoError(t, err)
		assert.Equal(t, int64(15), page1.Total)
		assert.Len(t, page1.Logs, 5)

		// 第二页
		page2, err := env.ListLogsHandler.Handle(ctx, auditlogQuery.ListLogsQuery{
			Page:     2,
			Limit:    5,
			Resource: "pagination_test",
		})
		require.NoError(t, err)
		assert.Len(t, page2.Logs, 5)

		// 第三页
		page3, err := env.ListLogsHandler.Handle(ctx, auditlogQuery.ListLogsQuery{
			Page:     3,
			Limit:    5,
			Resource: "pagination_test",
		})
		require.NoError(t, err)
		assert.Len(t, page3.Logs, 5)
	})
}
