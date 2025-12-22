package query

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
)

// listMockUserQueryRepo 用于列表测试的查询仓储 Mock
type listMockUserQueryRepo struct {
	mockUserQueryRepo

	users          []*user.User
	total          int64
	listErr        error
	countErr       error
	searchUsers    []*user.User
	searchTotal    int64
	searchErr      error
	searchCountErr error
}

func (m *listMockUserQueryRepo) List(ctx context.Context, offset, limit int) ([]*user.User, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	// 模拟分页
	start := offset
	end := offset + limit
	if start >= len(m.users) {
		return []*user.User{}, nil
	}
	if end > len(m.users) {
		end = len(m.users)
	}
	return m.users[start:end], nil
}

func (m *listMockUserQueryRepo) Count(ctx context.Context) (int64, error) {
	if m.countErr != nil {
		return 0, m.countErr
	}
	return m.total, nil
}

func (m *listMockUserQueryRepo) Search(ctx context.Context, keyword string, offset, limit int) ([]*user.User, error) {
	if m.searchErr != nil {
		return nil, m.searchErr
	}
	return m.searchUsers, nil
}

func (m *listMockUserQueryRepo) CountBySearch(ctx context.Context, keyword string) (int64, error) {
	if m.searchCountErr != nil {
		return 0, m.searchCountErr
	}
	return m.searchTotal, nil
}

type listUsersTestSetup struct {
	handler   *ListUsersHandler
	queryRepo *listMockUserQueryRepo
}

func setupListUsersTest() *listUsersTestSetup {
	queryRepo := &listMockUserQueryRepo{mockUserQueryRepo: *newMockUserQueryRepo()}
	handler := NewListUsersHandler(queryRepo)

	return &listUsersTestSetup{
		handler:   handler,
		queryRepo: queryRepo,
	}
}

func TestListUsersHandler_Handle(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	t.Run("成功获取用户列表", func(t *testing.T) {
		setup := setupListUsersTest()
		setup.queryRepo.users = []*user.User{
			{ID: 1, Username: "user1", Email: "user1@example.com", Status: "active", CreatedAt: now, UpdatedAt: now},
			{ID: 2, Username: "user2", Email: "user2@example.com", Status: "active", CreatedAt: now, UpdatedAt: now},
			{ID: 3, Username: "user3", Email: "user3@example.com", Status: "inactive", CreatedAt: now, UpdatedAt: now},
		}
		setup.queryRepo.total = 3

		query := ListUsersQuery{
			Offset: 0,
			Limit:  10,
		}

		result, err := setup.handler.Handle(ctx, query)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result.Users, 3)
		assert.Equal(t, int64(3), result.Total)
		assert.Equal(t, "user1", result.Users[0].Username)
	})

	t.Run("分页 - 第二页", func(t *testing.T) {
		setup := setupListUsersTest()
		setup.queryRepo.users = []*user.User{
			{ID: 1, Username: "user1", Email: "user1@example.com", Status: "active", CreatedAt: now, UpdatedAt: now},
			{ID: 2, Username: "user2", Email: "user2@example.com", Status: "active", CreatedAt: now, UpdatedAt: now},
			{ID: 3, Username: "user3", Email: "user3@example.com", Status: "active", CreatedAt: now, UpdatedAt: now},
			{ID: 4, Username: "user4", Email: "user4@example.com", Status: "active", CreatedAt: now, UpdatedAt: now},
			{ID: 5, Username: "user5", Email: "user5@example.com", Status: "active", CreatedAt: now, UpdatedAt: now},
		}
		setup.queryRepo.total = 5

		query := ListUsersQuery{
			Offset: 2,
			Limit:  2,
		}

		result, err := setup.handler.Handle(ctx, query)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result.Users, 2)
		assert.Equal(t, int64(5), result.Total)
		assert.Equal(t, "user3", result.Users[0].Username)
		assert.Equal(t, "user4", result.Users[1].Username)
	})

	t.Run("空列表", func(t *testing.T) {
		setup := setupListUsersTest()
		setup.queryRepo.users = []*user.User{}
		setup.queryRepo.total = 0

		query := ListUsersQuery{
			Offset: 0,
			Limit:  10,
		}

		result, err := setup.handler.Handle(ctx, query)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Empty(t, result.Users)
		assert.Equal(t, int64(0), result.Total)
	})

	t.Run("搜索用户", func(t *testing.T) {
		setup := setupListUsersTest()
		setup.queryRepo.searchUsers = []*user.User{
			{ID: 1, Username: "admin", Email: "admin@example.com", FullName: "Admin User", Status: "active", CreatedAt: now, UpdatedAt: now},
		}
		setup.queryRepo.searchTotal = 1

		query := ListUsersQuery{
			Search: "admin",
			Offset: 0,
			Limit:  10,
		}

		result, err := setup.handler.Handle(ctx, query)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result.Users, 1)
		assert.Equal(t, int64(1), result.Total)
		assert.Equal(t, "admin", result.Users[0].Username)
	})

	t.Run("搜索无结果", func(t *testing.T) {
		setup := setupListUsersTest()
		setup.queryRepo.searchUsers = []*user.User{}
		setup.queryRepo.searchTotal = 0

		query := ListUsersQuery{
			Search: "nonexistent",
			Offset: 0,
			Limit:  10,
		}

		result, err := setup.handler.Handle(ctx, query)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Empty(t, result.Users)
		assert.Equal(t, int64(0), result.Total)
	})

	t.Run("列表查询失败", func(t *testing.T) {
		setup := setupListUsersTest()
		setup.queryRepo.listErr = errors.New("database error")

		query := ListUsersQuery{
			Offset: 0,
			Limit:  10,
		}

		result, err := setup.handler.Handle(ctx, query)

		require.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("计数查询失败", func(t *testing.T) {
		setup := setupListUsersTest()
		setup.queryRepo.users = []*user.User{
			{ID: 1, Username: "user1", CreatedAt: now, UpdatedAt: now},
		}
		setup.queryRepo.countErr = errors.New("count error")

		query := ListUsersQuery{
			Offset: 0,
			Limit:  10,
		}

		result, err := setup.handler.Handle(ctx, query)

		require.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("搜索查询失败", func(t *testing.T) {
		setup := setupListUsersTest()
		setup.queryRepo.searchErr = errors.New("search error")

		query := ListUsersQuery{
			Search: "test",
			Offset: 0,
			Limit:  10,
		}

		result, err := setup.handler.Handle(ctx, query)

		require.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("搜索计数失败", func(t *testing.T) {
		setup := setupListUsersTest()
		setup.queryRepo.searchUsers = []*user.User{
			{ID: 1, Username: "user1", CreatedAt: now, UpdatedAt: now},
		}
		setup.queryRepo.searchCountErr = errors.New("search count error")

		query := ListUsersQuery{
			Search: "test",
			Offset: 0,
			Limit:  10,
		}

		result, err := setup.handler.Handle(ctx, query)

		require.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestNewListUsersHandler(t *testing.T) {
	t.Run("创建 ListUsersHandler", func(t *testing.T) {
		queryRepo := newMockUserQueryRepo()
		handler := NewListUsersHandler(queryRepo)

		assert.NotNil(t, handler)
	})
}
