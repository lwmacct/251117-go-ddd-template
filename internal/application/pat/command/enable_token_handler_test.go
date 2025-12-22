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

// enableMockPATCommandRepo 用于启用测试的命令仓储 Mock
type enableMockPATCommandRepo struct {
	mockPATCommandRepo

	enableError error
	enabledID   uint
}

func (m *enableMockPATCommandRepo) Enable(ctx context.Context, id uint) error {
	if m.enableError != nil {
		return m.enableError
	}
	m.enabledID = id
	return nil
}

type enableTokenTestSetup struct {
	handler     *EnableTokenHandler
	commandRepo *enableMockPATCommandRepo
	queryRepo   *mockPATQueryRepo
}

func setupEnableTokenTest() *enableTokenTestSetup {
	commandRepo := &enableMockPATCommandRepo{mockPATCommandRepo: *newMockPATCommandRepo()}
	queryRepo := newMockPATQueryRepo()
	handler := NewEnableTokenHandler(commandRepo, queryRepo)

	return &enableTokenTestSetup{
		handler:     handler,
		commandRepo: commandRepo,
		queryRepo:   queryRepo,
	}
}

func TestEnableTokenHandler_Handle(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	t.Run("成功启用Token", func(t *testing.T) {
		setup := setupEnableTokenTest()
		setup.queryRepo.tokens[1] = &pat.PersonalAccessToken{
			ID:        1,
			UserID:    100,
			Name:      "Disabled Token",
			Status:    pat.StatusDisabled,
			ExpiresAt: nil, // 不过期
			CreatedAt: now,
		}

		cmd := EnableTokenCommand{
			UserID:  100,
			TokenID: 1,
		}

		err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		assert.Equal(t, uint(1), setup.commandRepo.enabledID)
	})

	t.Run("Token不存在", func(t *testing.T) {
		setup := setupEnableTokenTest()

		cmd := EnableTokenCommand{
			UserID:  100,
			TokenID: 999,
		}

		err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "token not found")
	})

	t.Run("Token不属于该用户", func(t *testing.T) {
		setup := setupEnableTokenTest()
		setup.queryRepo.tokens[1] = &pat.PersonalAccessToken{
			ID:        1,
			UserID:    200,
			Name:      "Other User Token",
			Status:    pat.StatusDisabled,
			CreatedAt: now,
		}

		cmd := EnableTokenCommand{
			UserID:  100,
			TokenID: 1,
		}

		err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "token does not belong to this user")
	})

	t.Run("不能启用已过期的Token", func(t *testing.T) {
		setup := setupEnableTokenTest()
		expiredTime := now.Add(-24 * time.Hour) // 已过期
		setup.queryRepo.tokens[1] = &pat.PersonalAccessToken{
			ID:        1,
			UserID:    100,
			Name:      "Expired Token",
			Status:    pat.StatusDisabled,
			ExpiresAt: &expiredTime,
			CreatedAt: now,
		}

		cmd := EnableTokenCommand{
			UserID:  100,
			TokenID: 1,
		}

		err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "token is expired and cannot be enabled")
	})

	t.Run("启用操作失败", func(t *testing.T) {
		setup := setupEnableTokenTest()
		setup.queryRepo.tokens[1] = &pat.PersonalAccessToken{
			ID:        1,
			UserID:    100,
			Name:      "Disabled Token",
			Status:    pat.StatusDisabled,
			CreatedAt: now,
		}
		setup.commandRepo.enableError = errors.New("database error")

		cmd := EnableTokenCommand{
			UserID:  100,
			TokenID: 1,
		}

		err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to enable token")
	})

	t.Run("启用未过期Token且有过期时间", func(t *testing.T) {
		setup := setupEnableTokenTest()
		futureTime := now.Add(24 * time.Hour) // 未过期
		setup.queryRepo.tokens[1] = &pat.PersonalAccessToken{
			ID:        1,
			UserID:    100,
			Name:      "Valid Token",
			Status:    pat.StatusDisabled,
			ExpiresAt: &futureTime,
			CreatedAt: now,
		}

		cmd := EnableTokenCommand{
			UserID:  100,
			TokenID: 1,
		}

		err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		assert.Equal(t, uint(1), setup.commandRepo.enabledID)
	})
}

func TestNewEnableTokenHandler(t *testing.T) {
	commandRepo := newMockPATCommandRepo()
	queryRepo := newMockPATQueryRepo()

	handler := NewEnableTokenHandler(commandRepo, queryRepo)

	assert.NotNil(t, handler)
}
