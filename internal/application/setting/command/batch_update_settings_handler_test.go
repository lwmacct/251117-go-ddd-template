package command

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

// batchMockSettingCommandRepo 用于批量更新测试的命令仓储 Mock
type batchMockSettingCommandRepo struct {
	mockSettingCommandRepo

	batchUpsertError    error
	batchUpsertSettings []*setting.Setting
}

func (m *batchMockSettingCommandRepo) BatchUpsert(ctx context.Context, settings []*setting.Setting) error {
	if m.batchUpsertError != nil {
		return m.batchUpsertError
	}
	m.batchUpsertSettings = settings
	return nil
}

type batchUpdateSettingsTestSetup struct {
	handler     *BatchUpdateSettingsHandler
	commandRepo *batchMockSettingCommandRepo
	queryRepo   *mockSettingQueryRepo
}

func setupBatchUpdateSettingsTest() *batchUpdateSettingsTestSetup {
	commandRepo := &batchMockSettingCommandRepo{mockSettingCommandRepo: *newMockSettingCommandRepo()}
	queryRepo := newMockSettingQueryRepo()
	handler := NewBatchUpdateSettingsHandler(commandRepo, queryRepo)

	return &batchUpdateSettingsTestSetup{
		handler:     handler,
		commandRepo: commandRepo,
		queryRepo:   queryRepo,
	}
}

func TestBatchUpdateSettingsHandler_Handle(t *testing.T) {
	ctx := context.Background()

	t.Run("成功批量更新设置", func(t *testing.T) {
		setup := setupBatchUpdateSettingsTest()
		setup.queryRepo.settings["site.name"] = &setting.Setting{
			ID:    1,
			Key:   "site.name",
			Value: "Old Site",
		}
		setup.queryRepo.settings["site.title"] = &setting.Setting{
			ID:    2,
			Key:   "site.title",
			Value: "Old Title",
		}

		cmd := BatchUpdateSettingsCommand{
			Settings: []SettingItem{
				{Key: "site.name", Value: "New Site"},
				{Key: "site.title", Value: "New Title"},
			},
		}

		err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		require.Len(t, setup.commandRepo.batchUpsertSettings, 2)
		assert.Equal(t, "New Site", setup.commandRepo.batchUpsertSettings[0].Value)
		assert.Equal(t, "New Title", setup.commandRepo.batchUpsertSettings[1].Value)
	})

	t.Run("空设置列表", func(t *testing.T) {
		setup := setupBatchUpdateSettingsTest()

		cmd := BatchUpdateSettingsCommand{
			Settings: []SettingItem{},
		}

		err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
	})

	t.Run("部分设置不存在", func(t *testing.T) {
		setup := setupBatchUpdateSettingsTest()
		setup.queryRepo.settings["existing.key"] = &setting.Setting{
			ID:  1,
			Key: "existing.key",
		}

		cmd := BatchUpdateSettingsCommand{
			Settings: []SettingItem{
				{Key: "existing.key", Value: "new value"},
				{Key: "missing.key", Value: "value"},
			},
		}

		err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "setting key missing.key does not exist")
	})

	t.Run("查询设置失败", func(t *testing.T) {
		setup := setupBatchUpdateSettingsTest()
		setup.queryRepo.findError = errors.New("database error")

		cmd := BatchUpdateSettingsCommand{
			Settings: []SettingItem{
				{Key: "any.key", Value: "value"},
			},
		}

		err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to find setting")
	})

	t.Run("批量更新失败", func(t *testing.T) {
		setup := setupBatchUpdateSettingsTest()
		setup.queryRepo.settings["test.key"] = &setting.Setting{
			ID:  1,
			Key: "test.key",
		}
		setup.commandRepo.batchUpsertError = errors.New("batch update failed")

		cmd := BatchUpdateSettingsCommand{
			Settings: []SettingItem{
				{Key: "test.key", Value: "new value"},
			},
		}

		err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to batch update settings")
	})

	t.Run("批量更新多个设置保留原有属性", func(t *testing.T) {
		setup := setupBatchUpdateSettingsTest()
		setup.queryRepo.settings["app.debug"] = &setting.Setting{
			ID:        1,
			Key:       "app.debug",
			Value:     "false",
			ValueType: setting.ValueTypeBoolean,
			Label:     "Debug Mode",
			Category:  setting.CategoryGeneral,
		}

		cmd := BatchUpdateSettingsCommand{
			Settings: []SettingItem{
				{Key: "app.debug", Value: "true"},
			},
		}

		err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		require.Len(t, setup.commandRepo.batchUpsertSettings, 1)
		updated := setup.commandRepo.batchUpsertSettings[0]
		assert.Equal(t, "true", updated.Value)
		assert.Equal(t, setting.ValueTypeBoolean, updated.ValueType)
		assert.Equal(t, "Debug Mode", updated.Label)
	})
}

func TestNewBatchUpdateSettingsHandler(t *testing.T) {
	commandRepo := newMockSettingCommandRepo()
	queryRepo := newMockSettingQueryRepo()

	handler := NewBatchUpdateSettingsHandler(commandRepo, queryRepo)

	assert.NotNil(t, handler)
}
