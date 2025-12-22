package command

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

// createMockSettingCommandRepo 用于创建测试的命令仓储 Mock
type createMockSettingCommandRepo struct {
	mockSettingCommandRepo

	createError   error
	createdSeting *setting.Setting
}

func (m *createMockSettingCommandRepo) Create(ctx context.Context, s *setting.Setting) error {
	if m.createError != nil {
		return m.createError
	}
	s.ID = 1 // 模拟数据库生成 ID
	m.createdSeting = s
	return nil
}

type createSettingTestSetup struct {
	handler     *CreateSettingHandler
	commandRepo *createMockSettingCommandRepo
	queryRepo   *mockSettingQueryRepo
}

func setupCreateSettingTest() *createSettingTestSetup {
	commandRepo := &createMockSettingCommandRepo{mockSettingCommandRepo: *newMockSettingCommandRepo()}
	queryRepo := newMockSettingQueryRepo()
	handler := NewCreateSettingHandler(commandRepo, queryRepo)

	return &createSettingTestSetup{
		handler:     handler,
		commandRepo: commandRepo,
		queryRepo:   queryRepo,
	}
}

func TestCreateSettingHandler_Handle(t *testing.T) {
	ctx := context.Background()

	t.Run("成功创建设置", func(t *testing.T) {
		setup := setupCreateSettingTest()

		cmd := CreateSettingCommand{
			Key:       "site.name",
			Value:     "My Site",
			Category:  setting.CategoryGeneral,
			ValueType: setting.ValueTypeString,
			Label:     "Site Name",
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, uint(1), result.ID)
		assert.Equal(t, "site.name", setup.commandRepo.createdSeting.Key)
		assert.Equal(t, "My Site", setup.commandRepo.createdSeting.Value)
	})

	t.Run("Key已存在", func(t *testing.T) {
		setup := setupCreateSettingTest()
		setup.queryRepo.settings["existing.key"] = &setting.Setting{
			ID:  1,
			Key: "existing.key",
		}

		cmd := CreateSettingCommand{
			Key:   "existing.key",
			Value: "some value",
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "setting key already exists")
		assert.Nil(t, result)
	})

	t.Run("查询现有设置失败", func(t *testing.T) {
		setup := setupCreateSettingTest()
		setup.queryRepo.findError = errors.New("database error")

		cmd := CreateSettingCommand{
			Key:   "new.key",
			Value: "value",
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to check existing setting")
		assert.Nil(t, result)
	})

	t.Run("创建设置失败", func(t *testing.T) {
		setup := setupCreateSettingTest()
		setup.commandRepo.createError = errors.New("insert failed")

		cmd := CreateSettingCommand{
			Key:   "new.key",
			Value: "value",
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create setting")
		assert.Nil(t, result)
	})

	t.Run("默认值类型为String", func(t *testing.T) {
		setup := setupCreateSettingTest()

		cmd := CreateSettingCommand{
			Key:       "test.key",
			Value:     "test value",
			ValueType: "", // 不指定类型
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, setting.ValueTypeString, setup.commandRepo.createdSeting.ValueType)
	})

	t.Run("创建布尔类型设置", func(t *testing.T) {
		setup := setupCreateSettingTest()

		cmd := CreateSettingCommand{
			Key:       "app.debug",
			Value:     "true",
			Category:  setting.CategoryGeneral,
			ValueType: setting.ValueTypeBoolean,
			Label:     "Debug Mode",
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, setting.ValueTypeBoolean, setup.commandRepo.createdSeting.ValueType)
	})

	t.Run("创建JSON类型设置", func(t *testing.T) {
		setup := setupCreateSettingTest()

		cmd := CreateSettingCommand{
			Key:       "app.config",
			Value:     `{"enabled": true, "limit": 100}`,
			Category:  setting.CategoryGeneral,
			ValueType: setting.ValueTypeJSON,
			Label:     "App Configuration",
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, setting.ValueTypeJSON, setup.commandRepo.createdSeting.ValueType)
	})
}

func TestNewCreateSettingHandler(t *testing.T) {
	commandRepo := newMockSettingCommandRepo()
	queryRepo := newMockSettingQueryRepo()

	handler := NewCreateSettingHandler(commandRepo, queryRepo)

	assert.NotNil(t, handler)
}
