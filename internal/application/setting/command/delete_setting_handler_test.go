package command

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

// deleteMockSettingCommandRepo 用于删除测试的命令仓储 Mock
type deleteMockSettingCommandRepo struct {
	mockSettingCommandRepo

	deleteError error
	deletedID   uint
}

func (m *deleteMockSettingCommandRepo) Delete(ctx context.Context, id uint) error {
	if m.deleteError != nil {
		return m.deleteError
	}
	m.deletedID = id
	return nil
}

type deleteSettingTestSetup struct {
	handler     *DeleteSettingHandler
	commandRepo *deleteMockSettingCommandRepo
	queryRepo   *mockSettingQueryRepo
}

func setupDeleteSettingTest() *deleteSettingTestSetup {
	commandRepo := &deleteMockSettingCommandRepo{mockSettingCommandRepo: *newMockSettingCommandRepo()}
	queryRepo := newMockSettingQueryRepo()
	handler := NewDeleteSettingHandler(commandRepo, queryRepo)

	return &deleteSettingTestSetup{
		handler:     handler,
		commandRepo: commandRepo,
		queryRepo:   queryRepo,
	}
}

func TestDeleteSettingHandler_Handle(t *testing.T) {
	ctx := context.Background()

	t.Run("成功删除设置", func(t *testing.T) {
		setup := setupDeleteSettingTest()
		setup.queryRepo.settings["site.name"] = &setting.Setting{
			ID:    1,
			Key:   "site.name",
			Value: "My Site",
		}

		cmd := DeleteSettingCommand{
			Key: "site.name",
		}

		err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		assert.Equal(t, uint(1), setup.commandRepo.deletedID)
	})

	t.Run("设置不存在", func(t *testing.T) {
		setup := setupDeleteSettingTest()

		cmd := DeleteSettingCommand{
			Key: "nonexistent.key",
		}

		err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "setting not found")
	})

	t.Run("查询设置失败", func(t *testing.T) {
		setup := setupDeleteSettingTest()
		setup.queryRepo.findError = errors.New("database error")

		cmd := DeleteSettingCommand{
			Key: "any.key",
		}

		err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to find setting")
	})

	t.Run("删除操作失败", func(t *testing.T) {
		setup := setupDeleteSettingTest()
		setup.queryRepo.settings["test.key"] = &setting.Setting{
			ID:  1,
			Key: "test.key",
		}
		setup.commandRepo.deleteError = errors.New("delete failed")

		cmd := DeleteSettingCommand{
			Key: "test.key",
		}

		err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to delete setting")
	})
}

func TestNewDeleteSettingHandler(t *testing.T) {
	commandRepo := newMockSettingCommandRepo()
	queryRepo := newMockSettingQueryRepo()

	handler := NewDeleteSettingHandler(commandRepo, queryRepo)

	assert.NotNil(t, handler)
}
