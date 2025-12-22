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

// mockPATQueryRepo 是 PAT 查询仓储的 Mock 实现
type mockPATQueryRepo struct {
	tokens      map[uint]*pat.PersonalAccessToken
	userTokens  []*pat.PersonalAccessToken
	findByIDErr error
	listErr     error
}

func newMockPATQueryRepo() *mockPATQueryRepo {
	return &mockPATQueryRepo{
		tokens: make(map[uint]*pat.PersonalAccessToken),
	}
}

func (m *mockPATQueryRepo) FindByToken(ctx context.Context, tokenHash string) (*pat.PersonalAccessToken, error) {
	return nil, nil //nolint:nilnil // test mock returns nil for not found
}

func (m *mockPATQueryRepo) FindByID(ctx context.Context, id uint) (*pat.PersonalAccessToken, error) {
	if m.findByIDErr != nil {
		return nil, m.findByIDErr
	}
	if t, ok := m.tokens[id]; ok {
		return t, nil
	}
	return nil, nil //nolint:nilnil // test mock returns nil for not found
}

func (m *mockPATQueryRepo) FindByPrefix(ctx context.Context, prefix string) (*pat.PersonalAccessToken, error) {
	return nil, nil //nolint:nilnil // test mock returns nil for not found
}

func (m *mockPATQueryRepo) ListByUser(ctx context.Context, userID uint) ([]*pat.PersonalAccessToken, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	return m.userTokens, nil
}

type getTokenTestSetup struct {
	handler   *GetTokenHandler
	queryRepo *mockPATQueryRepo
}

func setupGetTokenTest() *getTokenTestSetup {
	queryRepo := newMockPATQueryRepo()
	handler := NewGetTokenHandler(queryRepo)

	return &getTokenTestSetup{
		handler:   handler,
		queryRepo: queryRepo,
	}
}

func TestGetTokenHandler_Handle(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	t.Run("成功获取Token", func(t *testing.T) {
		setup := setupGetTokenTest()
		setup.queryRepo.tokens[1] = &pat.PersonalAccessToken{
			ID:          1,
			UserID:      100,
			Name:        "My API Token",
			TokenPrefix: "pat_abc",
			Status:      pat.StatusActive,
			Permissions: []string{"read:user", "write:user"},
			CreatedAt:   now,
		}

		query := GetTokenQuery{
			UserID:  100,
			TokenID: 1,
		}

		result, err := setup.handler.Handle(ctx, query)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, uint(1), result.ID)
		assert.Equal(t, "My API Token", result.Name)
		assert.Equal(t, "pat_abc", result.TokenPrefix)
		assert.Equal(t, pat.StatusActive, result.Status)
	})

	t.Run("Token不存在", func(t *testing.T) {
		setup := setupGetTokenTest()

		query := GetTokenQuery{
			UserID:  100,
			TokenID: 999,
		}

		result, err := setup.handler.Handle(ctx, query)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "token not found")
		assert.Nil(t, result)
	})

	t.Run("Token不属于该用户", func(t *testing.T) {
		setup := setupGetTokenTest()
		setup.queryRepo.tokens[1] = &pat.PersonalAccessToken{
			ID:        1,
			UserID:    200, // 属于另一个用户
			Name:      "Other Token",
			Status:    pat.StatusActive,
			CreatedAt: now,
		}

		query := GetTokenQuery{
			UserID:  100, // 不同的用户
			TokenID: 1,
		}

		result, err := setup.handler.Handle(ctx, query)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "token does not belong to this user")
		assert.Nil(t, result)
	})

	t.Run("查询失败", func(t *testing.T) {
		setup := setupGetTokenTest()
		setup.queryRepo.findByIDErr = errors.New("database error")

		query := GetTokenQuery{
			UserID:  100,
			TokenID: 1,
		}

		result, err := setup.handler.Handle(ctx, query)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "token not found")
		assert.Nil(t, result)
	})

	t.Run("获取带过期时间的Token", func(t *testing.T) {
		setup := setupGetTokenTest()
		expiresAt := now.Add(30 * 24 * time.Hour)
		setup.queryRepo.tokens[1] = &pat.PersonalAccessToken{
			ID:          1,
			UserID:      100,
			Name:        "Expiring Token",
			TokenPrefix: "pat_xyz",
			Status:      pat.StatusActive,
			ExpiresAt:   &expiresAt,
			CreatedAt:   now,
		}

		query := GetTokenQuery{
			UserID:  100,
			TokenID: 1,
		}

		result, err := setup.handler.Handle(ctx, query)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotNil(t, result.ExpiresAt)
	})

	t.Run("获取禁用状态的Token", func(t *testing.T) {
		setup := setupGetTokenTest()
		setup.queryRepo.tokens[1] = &pat.PersonalAccessToken{
			ID:        1,
			UserID:    100,
			Name:      "Disabled Token",
			Status:    pat.StatusDisabled,
			CreatedAt: now,
		}

		query := GetTokenQuery{
			UserID:  100,
			TokenID: 1,
		}

		result, err := setup.handler.Handle(ctx, query)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, pat.StatusDisabled, result.Status)
	})
}

func TestNewGetTokenHandler(t *testing.T) {
	queryRepo := newMockPATQueryRepo()

	handler := NewGetTokenHandler(queryRepo)

	assert.NotNil(t, handler)
}
