package query

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

// listMockSettingQueryRepo 用于列表测试的查询仓储 Mock
type listMockSettingQueryRepo struct {
	mockSettingQueryRepo

	allSettings      []*setting.Setting
	categorySettings []*setting.Setting
	findAllError     error
	findCategoryErr  error
}

func newListMockSettingQueryRepo() *listMockSettingQueryRepo {
	return &listMockSettingQueryRepo{
		mockSettingQueryRepo: *newMockSettingQueryRepo(),
	}
}

func (m *listMockSettingQueryRepo) FindAll(ctx context.Context) ([]*setting.Setting, error) {
	if m.findAllError != nil {
		return nil, m.findAllError
	}
	return m.allSettings, nil
}

func (m *listMockSettingQueryRepo) FindByCategory(ctx context.Context, category string) ([]*setting.Setting, error) {
	if m.findCategoryErr != nil {
		return nil, m.findCategoryErr
	}
	return m.categorySettings, nil
}

type listSettingsTestSetup struct {
	handler   *ListSettingsHandler
	queryRepo *listMockSettingQueryRepo
}

func setupListSettingsTest() *listSettingsTestSetup {
	queryRepo := newListMockSettingQueryRepo()
	handler := NewListSettingsHandler(queryRepo)

	return &listSettingsTestSetup{
		handler:   handler,
		queryRepo: queryRepo,
	}
}

func TestListSettingsHandler_Handle(t *testing.T) {
	ctx := context.Background()

	t.Run("成功获取所有设置", func(t *testing.T) {
		setup := setupListSettingsTest()
		setup.queryRepo.allSettings = []*setting.Setting{
			{ID: 1, Key: "site.name", Value: "My Site", Category: setting.CategoryGeneral},
			{ID: 2, Key: "site.title", Value: "My Title", Category: setting.CategoryGeneral},
			{ID: 3, Key: "app.debug", Value: "false", Category: setting.CategoryGeneral},
		}

		query := ListSettingsQuery{
			Category: "", // 不指定分类，获取所有
		}

		result, err := setup.handler.Handle(ctx, query)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result, 3)
	})

	t.Run("按分类获取设置", func(t *testing.T) {
		setup := setupListSettingsTest()
		setup.queryRepo.categorySettings = []*setting.Setting{
			{ID: 1, Key: "site.name", Value: "My Site", Category: setting.CategoryGeneral},
			{ID: 2, Key: "site.title", Value: "My Title", Category: setting.CategoryGeneral},
		}

		query := ListSettingsQuery{
			Category: setting.CategoryGeneral,
		}

		result, err := setup.handler.Handle(ctx, query)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result, 2)
	})

	t.Run("空设置列表", func(t *testing.T) {
		setup := setupListSettingsTest()
		setup.queryRepo.allSettings = []*setting.Setting{}

		query := ListSettingsQuery{
			Category: "",
		}

		result, err := setup.handler.Handle(ctx, query)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Empty(t, result)
	})

	t.Run("获取所有设置失败", func(t *testing.T) {
		setup := setupListSettingsTest()
		setup.queryRepo.findAllError = errors.New("database error")

		query := ListSettingsQuery{
			Category: "",
		}

		result, err := setup.handler.Handle(ctx, query)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to fetch settings")
		assert.Nil(t, result)
	})

	t.Run("按分类获取设置失败", func(t *testing.T) {
		setup := setupListSettingsTest()
		setup.queryRepo.findCategoryErr = errors.New("database error")

		query := ListSettingsQuery{
			Category: setting.CategoryGeneral,
		}

		result, err := setup.handler.Handle(ctx, query)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to fetch settings")
		assert.Nil(t, result)
	})

	t.Run("返回设置包含所有字段", func(t *testing.T) {
		setup := setupListSettingsTest()
		setup.queryRepo.allSettings = []*setting.Setting{
			{
				ID:        1,
				Key:       "app.config",
				Value:     `{"enabled": true}`,
				Category:  setting.CategoryGeneral,
				ValueType: setting.ValueTypeJSON,
				Label:     "App Configuration",
			},
		}

		query := ListSettingsQuery{
			Category: "",
		}

		result, err := setup.handler.Handle(ctx, query)

		require.NoError(t, err)
		require.Len(t, result, 1)
		assert.Equal(t, uint(1), result[0].ID)
		assert.Equal(t, "app.config", result[0].Key)
		assert.Equal(t, setting.ValueTypeJSON, result[0].ValueType)
		assert.Equal(t, "App Configuration", result[0].Label)
	})
}

func TestNewListSettingsHandler(t *testing.T) {
	queryRepo := newListMockSettingQueryRepo()

	handler := NewListSettingsHandler(queryRepo)

	assert.NotNil(t, handler)
}
