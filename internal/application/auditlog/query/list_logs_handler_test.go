package query

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/auditlog"
)

// mockAuditLogQueryRepo 是审计日志查询仓储的 Mock 实现
type mockAuditLogQueryRepo struct {
	logs    []auditlog.AuditLog
	total   int64
	listErr error
}

func newMockAuditLogQueryRepo() *mockAuditLogQueryRepo {
	return &mockAuditLogQueryRepo{}
}

func (m *mockAuditLogQueryRepo) FindByID(ctx context.Context, id uint) (*auditlog.AuditLog, error) {
	return nil, nil //nolint:nilnil // test mock returns nil for not found
}

func (m *mockAuditLogQueryRepo) List(ctx context.Context, filter auditlog.FilterOptions) ([]auditlog.AuditLog, int64, error) {
	if m.listErr != nil {
		return nil, 0, m.listErr
	}
	return m.logs, m.total, nil
}

func (m *mockAuditLogQueryRepo) ListByUser(ctx context.Context, userID uint, page, limit int) ([]auditlog.AuditLog, int64, error) {
	return nil, 0, nil
}

func (m *mockAuditLogQueryRepo) ListByResource(ctx context.Context, resource string, page, limit int) ([]auditlog.AuditLog, int64, error) {
	return nil, 0, nil
}

func (m *mockAuditLogQueryRepo) ListByAction(ctx context.Context, action string, page, limit int) ([]auditlog.AuditLog, int64, error) {
	return nil, 0, nil
}

func (m *mockAuditLogQueryRepo) Count(ctx context.Context, filter auditlog.FilterOptions) (int64, error) {
	return m.total, nil
}

func (m *mockAuditLogQueryRepo) Search(ctx context.Context, keyword string, page, limit int) ([]auditlog.AuditLog, int64, error) {
	return nil, 0, nil
}

type listLogsTestSetup struct {
	handler   *ListLogsHandler
	queryRepo *mockAuditLogQueryRepo
}

func setupListLogsTest() *listLogsTestSetup {
	queryRepo := newMockAuditLogQueryRepo()
	handler := NewListLogsHandler(queryRepo)

	return &listLogsTestSetup{
		handler:   handler,
		queryRepo: queryRepo,
	}
}

func TestListLogsHandler_Handle(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	t.Run("成功获取审计日志列表", func(t *testing.T) {
		setup := setupListLogsTest()
		setup.queryRepo.logs = []auditlog.AuditLog{
			{ID: 1, UserID: 1, Username: "admin", Action: "login", Resource: "auth", Status: "success", CreatedAt: now},
			{ID: 2, UserID: 1, Username: "admin", Action: "update", Resource: "user", Status: "success", CreatedAt: now},
			{ID: 3, UserID: 2, Username: "user1", Action: "login", Resource: "auth", Status: "failed", CreatedAt: now},
		}
		setup.queryRepo.total = 3

		query := ListLogsQuery{
			Page:  1,
			Limit: 10,
		}

		result, err := setup.handler.Handle(ctx, query)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result.Logs, 3)
		assert.Equal(t, int64(3), result.Total)
		assert.Equal(t, 1, result.Page)
		assert.Equal(t, 10, result.Limit)
	})

	t.Run("分页", func(t *testing.T) {
		setup := setupListLogsTest()
		setup.queryRepo.logs = []auditlog.AuditLog{
			{ID: 3, UserID: 1, Action: "update", Resource: "role", Status: "success", CreatedAt: now},
		}
		setup.queryRepo.total = 10

		query := ListLogsQuery{
			Page:  2,
			Limit: 5,
		}

		result, err := setup.handler.Handle(ctx, query)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result.Logs, 1)
		assert.Equal(t, int64(10), result.Total)
		assert.Equal(t, 2, result.Page)
	})

	t.Run("空列表", func(t *testing.T) {
		setup := setupListLogsTest()
		setup.queryRepo.logs = []auditlog.AuditLog{}
		setup.queryRepo.total = 0

		query := ListLogsQuery{
			Page:  1,
			Limit: 10,
		}

		result, err := setup.handler.Handle(ctx, query)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Empty(t, result.Logs)
		assert.Equal(t, int64(0), result.Total)
	})

	t.Run("带筛选条件", func(t *testing.T) {
		setup := setupListLogsTest()
		userID := uint(1)
		setup.queryRepo.logs = []auditlog.AuditLog{
			{ID: 1, UserID: 1, Username: "admin", Action: "login", Resource: "auth", Status: "success", CreatedAt: now},
		}
		setup.queryRepo.total = 1

		query := ListLogsQuery{
			Page:     1,
			Limit:    10,
			UserID:   &userID,
			Action:   "login",
			Resource: "auth",
			Status:   "success",
		}

		result, err := setup.handler.Handle(ctx, query)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result.Logs, 1)
	})

	t.Run("查询失败", func(t *testing.T) {
		setup := setupListLogsTest()
		setup.queryRepo.listErr = errors.New("database error")

		query := ListLogsQuery{
			Page:  1,
			Limit: 10,
		}

		result, err := setup.handler.Handle(ctx, query)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to list audit logs")
		assert.Nil(t, result)
	})

	t.Run("带日期范围筛选", func(t *testing.T) {
		setup := setupListLogsTest()
		startDate := now.Add(-24 * time.Hour)
		endDate := now
		setup.queryRepo.logs = []auditlog.AuditLog{
			{ID: 1, UserID: 1, Action: "login", Resource: "auth", Status: "success", CreatedAt: now},
		}
		setup.queryRepo.total = 1

		query := ListLogsQuery{
			Page:      1,
			Limit:     10,
			StartDate: &startDate,
			EndDate:   &endDate,
		}

		result, err := setup.handler.Handle(ctx, query)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result.Logs, 1)
	})
}

func TestNewListLogsHandler(t *testing.T) {
	queryRepo := newMockAuditLogQueryRepo()
	handler := NewListLogsHandler(queryRepo)

	assert.NotNil(t, handler)
}
