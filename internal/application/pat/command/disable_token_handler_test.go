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

// disableMockPATCommandRepo 用于禁用测试的命令仓储 Mock
type disableMockPATCommandRepo struct {
	mockPATCommandRepo

	disableError error
	disabledID   uint
}

func (m *disableMockPATCommandRepo) Disable(ctx context.Context, id uint) error {
	if m.disableError != nil {
		return m.disableError
	}
	m.disabledID = id
	return nil
}

type disableTokenTestSetup struct {
	handler     *DisableTokenHandler
	commandRepo *disableMockPATCommandRepo
	queryRepo   *mockPATQueryRepo
}

func setupDisableTokenTest() *disableTokenTestSetup {
	commandRepo := &disableMockPATCommandRepo{mockPATCommandRepo: *newMockPATCommandRepo()}
	queryRepo := newMockPATQueryRepo()
	handler := NewDisableTokenHandler(commandRepo, queryRepo)

	return &disableTokenTestSetup{
		handler:     handler,
		commandRepo: commandRepo,
		queryRepo:   queryRepo,
	}
}

func TestDisableTokenHandler_Handle(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	t.Run("成功禁用Token", func(t *testing.T) {
		setup := setupDisableTokenTest()
		setup.queryRepo.tokens[1] = &pat.PersonalAccessToken{
			ID:        1,
			UserID:    100,
			Name:      "Test Token",
			Status:    pat.StatusActive,
			CreatedAt: now,
		}

		cmd := DisableTokenCommand{
			UserID:  100,
			TokenID: 1,
		}

		err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		assert.Equal(t, uint(1), setup.commandRepo.disabledID)
	})

	t.Run("Token不存在", func(t *testing.T) {
		setup := setupDisableTokenTest()

		cmd := DisableTokenCommand{
			UserID:  100,
			TokenID: 999,
		}

		err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "token not found")
	})

	t.Run("Token不属于该用户", func(t *testing.T) {
		setup := setupDisableTokenTest()
		setup.queryRepo.tokens[1] = &pat.PersonalAccessToken{
			ID:        1,
			UserID:    200,
			Name:      "Other User Token",
			Status:    pat.StatusActive,
			CreatedAt: now,
		}

		cmd := DisableTokenCommand{
			UserID:  100,
			TokenID: 1,
		}

		err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "token does not belong to this user")
	})

	t.Run("禁用操作失败", func(t *testing.T) {
		setup := setupDisableTokenTest()
		setup.queryRepo.tokens[1] = &pat.PersonalAccessToken{
			ID:        1,
			UserID:    100,
			Name:      "Test Token",
			Status:    pat.StatusActive,
			CreatedAt: now,
		}
		setup.commandRepo.disableError = errors.New("database error")

		cmd := DisableTokenCommand{
			UserID:  100,
			TokenID: 1,
		}

		err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to disable token")
	})

	t.Run("禁用已禁用的Token", func(t *testing.T) {
		setup := setupDisableTokenTest()
		setup.queryRepo.tokens[1] = &pat.PersonalAccessToken{
			ID:        1,
			UserID:    100,
			Name:      "Disabled Token",
			Status:    pat.StatusDisabled,
			CreatedAt: now,
		}

		cmd := DisableTokenCommand{
			UserID:  100,
			TokenID: 1,
		}

		// 再次禁用已禁用的Token应该成功（幂等操作）
		err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
	})
}

func TestNewDisableTokenHandler(t *testing.T) {
	commandRepo := newMockPATCommandRepo()
	queryRepo := newMockPATQueryRepo()

	handler := NewDisableTokenHandler(commandRepo, queryRepo)

	assert.NotNil(t, handler)
}
