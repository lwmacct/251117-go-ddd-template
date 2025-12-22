package query

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

// mockSettingQueryRepo 是设置查询仓储的 Mock 实现
type mockSettingQueryRepo struct {
	settings  map[string]*setting.Setting
	findError error
}

func newMockSettingQueryRepo() *mockSettingQueryRepo {
	return &mockSettingQueryRepo{
		settings: make(map[string]*setting.Setting),
	}
}

func (m *mockSettingQueryRepo) FindByID(ctx context.Context, id uint) (*setting.Setting, error) {
	return nil, nil //nolint:nilnil // test mock returns nil for not found
}

func (m *mockSettingQueryRepo) FindByKey(ctx context.Context, key string) (*setting.Setting, error) {
	if m.findError != nil {
		return nil, m.findError
	}
	if s, ok := m.settings[key]; ok {
		return s, nil
	}
	return nil, nil //nolint:nilnil // test mock returns nil for not found
}

func (m *mockSettingQueryRepo) FindByCategory(ctx context.Context, category string) ([]*setting.Setting, error) {
	return nil, nil
}

func (m *mockSettingQueryRepo) FindAll(ctx context.Context) ([]*setting.Setting, error) {
	return nil, nil
}

type getSettingTestSetup struct {
	handler   *GetSettingHandler
	queryRepo *mockSettingQueryRepo
}

func setupGetSettingTest() *getSettingTestSetup {
	queryRepo := newMockSettingQueryRepo()
	handler := NewGetSettingHandler(queryRepo)

	return &getSettingTestSetup{
		handler:   handler,
		queryRepo: queryRepo,
	}
}

func TestGetSettingHandler_Handle(t *testing.T) {
	ctx := context.Background()

	t.Run("成功获取设置", func(t *testing.T) {
		setup := setupGetSettingTest()
		setup.queryRepo.settings["site.name"] = &setting.Setting{
			ID:        1,
			Key:       "site.name",
			Value:     "My Site",
			Category:  setting.CategoryGeneral,
			ValueType: setting.ValueTypeString,
			Label:     "Site Name",
		}

		query := GetSettingQuery{
			Key: "site.name",
		}

		result, err := setup.handler.Handle(ctx, query)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, uint(1), result.ID)
		assert.Equal(t, "site.name", result.Key)
		assert.Equal(t, "My Site", result.Value)
		assert.Equal(t, "Site Name", result.Label)
	})

	t.Run("设置不存在", func(t *testing.T) {
		setup := setupGetSettingTest()

		query := GetSettingQuery{
			Key: "nonexistent.key",
		}

		result, err := setup.handler.Handle(ctx, query)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "setting not found")
		assert.Nil(t, result)
	})

	t.Run("查询失败", func(t *testing.T) {
		setup := setupGetSettingTest()
		setup.queryRepo.findError = errors.New("database error")

		query := GetSettingQuery{
			Key: "any.key",
		}

		result, err := setup.handler.Handle(ctx, query)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to find setting")
		assert.Nil(t, result)
	})

	t.Run("获取布尔类型设置", func(t *testing.T) {
		setup := setupGetSettingTest()
		setup.queryRepo.settings["app.debug"] = &setting.Setting{
			ID:        2,
			Key:       "app.debug",
			Value:     "true",
			ValueType: setting.ValueTypeBoolean,
			Label:     "Debug Mode",
		}

		query := GetSettingQuery{
			Key: "app.debug",
		}

		result, err := setup.handler.Handle(ctx, query)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "true", result.Value)
		assert.Equal(t, setting.ValueTypeBoolean, result.ValueType)
	})

	t.Run("获取JSON类型设置", func(t *testing.T) {
		setup := setupGetSettingTest()
		setup.queryRepo.settings["app.config"] = &setting.Setting{
			ID:        3,
			Key:       "app.config",
			Value:     `{"enabled": true, "limit": 100}`,
			ValueType: setting.ValueTypeJSON,
			Label:     "App Config",
		}

		query := GetSettingQuery{
			Key: "app.config",
		}

		result, err := setup.handler.Handle(ctx, query)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, setting.ValueTypeJSON, result.ValueType)
		assert.Contains(t, result.Value, "enabled")
	})
}

func TestNewGetSettingHandler(t *testing.T) {
	queryRepo := newMockSettingQueryRepo()

	handler := NewGetSettingHandler(queryRepo)

	assert.NotNil(t, handler)
}
