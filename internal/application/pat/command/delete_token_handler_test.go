package command

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/pat"
)

// mockPATCommandRepo 是 PAT 命令仓储的 Mock 实现
type mockPATCommandRepo struct {
	deleteError error
	deletedID   uint
}

func newMockPATCommandRepo() *mockPATCommandRepo {
	return &mockPATCommandRepo{}
}

func (m *mockPATCommandRepo) Create(ctx context.Context, p *pat.PersonalAccessToken) error {
	return nil
}

func (m *mockPATCommandRepo) Update(ctx context.Context, p *pat.PersonalAccessToken) error {
	return nil
}

func (m *mockPATCommandRepo) Delete(ctx context.Context, id uint) error {
	if m.deleteError != nil {
		return m.deleteError
	}
	m.deletedID = id
	return nil
}

func (m *mockPATCommandRepo) Disable(ctx context.Context, id uint) error {
	return nil
}

func (m *mockPATCommandRepo) Enable(ctx context.Context, id uint) error {
	return nil
}

func (m *mockPATCommandRepo) DeleteByUserID(ctx context.Context, userID uint) error {
	return nil
}

func (m *mockPATCommandRepo) CleanupExpired(ctx context.Context) error {
	return nil
}

// mockPATQueryRepo 是 PAT 查询仓储的 Mock 实现
type mockPATQueryRepo struct {
	tokens      map[uint]*pat.PersonalAccessToken
	findByIDErr error
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
	return nil, nil
}

type deleteTokenTestSetup struct {
	handler     *DeleteTokenHandler
	commandRepo *mockPATCommandRepo
	queryRepo   *mockPATQueryRepo
}

func setupDeleteTokenTest() *deleteTokenTestSetup {
	commandRepo := newMockPATCommandRepo()
	queryRepo := newMockPATQueryRepo()
	handler := NewDeleteTokenHandler(commandRepo, queryRepo)

	return &deleteTokenTestSetup{
		handler:     handler,
		commandRepo: commandRepo,
		queryRepo:   queryRepo,
	}
}

func TestDeleteTokenHandler_Handle(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	t.Run("成功删除Token", func(t *testing.T) {
		setup := setupDeleteTokenTest()
		setup.queryRepo.tokens[1] = &pat.PersonalAccessToken{
			ID:        1,
			UserID:    100,
			Name:      "Test Token",
			Status:    pat.StatusActive,
			CreatedAt: now,
		}

		cmd := DeleteTokenCommand{
			UserID:  100,
			TokenID: 1,
		}

		err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		assert.Equal(t, uint(1), setup.commandRepo.deletedID)
	})

	t.Run("Token不存在", func(t *testing.T) {
		setup := setupDeleteTokenTest()

		cmd := DeleteTokenCommand{
			UserID:  100,
			TokenID: 999,
		}

		err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "token not found")
	})

	t.Run("Token不属于该用户", func(t *testing.T) {
		setup := setupDeleteTokenTest()
		setup.queryRepo.tokens[1] = &pat.PersonalAccessToken{
			ID:        1,
			UserID:    200, // 属于另一个用户
			Name:      "Other User Token",
			Status:    pat.StatusActive,
			CreatedAt: now,
		}

		cmd := DeleteTokenCommand{
			UserID:  100, // 不同的用户
			TokenID: 1,
		}

		err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "token does not belong to this user")
	})

	t.Run("删除操作失败", func(t *testing.T) {
		setup := setupDeleteTokenTest()
		setup.queryRepo.tokens[1] = &pat.PersonalAccessToken{
			ID:        1,
			UserID:    100,
			Name:      "Test Token",
			Status:    pat.StatusActive,
			CreatedAt: now,
		}
		setup.commandRepo.deleteError = errors.New("database error")

		cmd := DeleteTokenCommand{
			UserID:  100,
			TokenID: 1,
		}

		err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to delete token")
	})

	t.Run("查询Token失败", func(t *testing.T) {
		setup := setupDeleteTokenTest()
		setup.queryRepo.findByIDErr = errors.New("database error")

		cmd := DeleteTokenCommand{
			UserID:  100,
			TokenID: 1,
		}

		err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "token not found")
	})
}

func TestNewDeleteTokenHandler(t *testing.T) {
	commandRepo := newMockPATCommandRepo()
	queryRepo := newMockPATQueryRepo()

	handler := NewDeleteTokenHandler(commandRepo, queryRepo)

	assert.NotNil(t, handler)
}
