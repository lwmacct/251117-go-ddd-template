package query

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/pat"
)

type listTokensTestSetup struct {
	handler   *ListTokensHandler
	queryRepo *mockPATQueryRepo
}

func setupListTokensTest() *listTokensTestSetup {
	queryRepo := newMockPATQueryRepo()
	handler := NewListTokensHandler(queryRepo)

	return &listTokensTestSetup{
		handler:   handler,
		queryRepo: queryRepo,
	}
}

func TestListTokensHandler_Handle(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	t.Run("成功获取Token列表", func(t *testing.T) {
		setup := setupListTokensTest()
		setup.queryRepo.userTokens = []*pat.PersonalAccessToken{
			{
				ID:          1,
				UserID:      100,
				Name:        "Token 1",
				TokenPrefix: "pat_abc",
				Status:      pat.StatusActive,
				Permissions: []string{"read:user"},
				CreatedAt:   now,
			},
			{
				ID:          2,
				UserID:      100,
				Name:        "Token 2",
				TokenPrefix: "pat_def",
				Status:      pat.StatusDisabled,
				Permissions: []string{"write:user"},
				CreatedAt:   now,
			},
		}

		query := ListTokensQuery{
			UserID: 100,
		}

		result, err := setup.handler.Handle(ctx, query)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result, 2)
		assert.Equal(t, "Token 1", result[0].Name)
		assert.Equal(t, "Token 2", result[1].Name)
	})

	t.Run("空Token列表", func(t *testing.T) {
		setup := setupListTokensTest()
		setup.queryRepo.userTokens = []*pat.PersonalAccessToken{}

		query := ListTokensQuery{
			UserID: 100,
		}

		result, err := setup.handler.Handle(ctx, query)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Empty(t, result)
	})

	t.Run("查询失败", func(t *testing.T) {
		setup := setupListTokensTest()
		setup.queryRepo.listErr = errors.New("database error")

		query := ListTokensQuery{
			UserID: 100,
		}

		result, err := setup.handler.Handle(ctx, query)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to fetch tokens")
		assert.Nil(t, result)
	})

	t.Run("返回Token包含所有字段", func(t *testing.T) {
		setup := setupListTokensTest()
		expiresAt := now.Add(30 * 24 * time.Hour)
		lastUsedAt := now.Add(-1 * time.Hour)
		setup.queryRepo.userTokens = []*pat.PersonalAccessToken{
			{
				ID:          1,
				UserID:      100,
				Name:        "Full Token",
				TokenPrefix: "pat_full",
				Status:      pat.StatusActive,
				Permissions: []string{"read:user", "write:user"},
				ExpiresAt:   &expiresAt,
				LastUsedAt:  &lastUsedAt,
				IPWhitelist: []string{"192.168.1.1"},
				Description: "Test token with all fields",
				CreatedAt:   now,
			},
		}

		query := ListTokensQuery{
			UserID: 100,
		}

		result, err := setup.handler.Handle(ctx, query)

		require.NoError(t, err)
		require.Len(t, result, 1)
		assert.Equal(t, "Full Token", result[0].Name)
		assert.Equal(t, "pat_full", result[0].TokenPrefix)
		assert.NotNil(t, result[0].ExpiresAt)
		assert.NotNil(t, result[0].LastUsedAt)
	})

	t.Run("混合状态Token列表", func(t *testing.T) {
		setup := setupListTokensTest()
		expiredTime := now.Add(-24 * time.Hour)
		setup.queryRepo.userTokens = []*pat.PersonalAccessToken{
			{ID: 1, UserID: 100, Name: "Active", Status: pat.StatusActive, CreatedAt: now},
			{ID: 2, UserID: 100, Name: "Disabled", Status: pat.StatusDisabled, CreatedAt: now},
			{ID: 3, UserID: 100, Name: "Expired", Status: pat.StatusExpired, ExpiresAt: &expiredTime, CreatedAt: now},
		}

		query := ListTokensQuery{
			UserID: 100,
		}

		result, err := setup.handler.Handle(ctx, query)

		require.NoError(t, err)
		assert.Len(t, result, 3)
	})
}

func TestNewListTokensHandler(t *testing.T) {
	queryRepo := newMockPATQueryRepo()

	handler := NewListTokensHandler(queryRepo)

	assert.NotNil(t, handler)
}
