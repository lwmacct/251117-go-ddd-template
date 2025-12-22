package command

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

// mockSettingCommandRepo 是设置命令仓储的 Mock 实现。
type mockSettingCommandRepo struct {
	updateError error
}

func newMockSettingCommandRepo() *mockSettingCommandRepo {
	return &mockSettingCommandRepo{}
}

func (m *mockSettingCommandRepo) Create(ctx context.Context, s *setting.Setting) error {
	return nil
}

func (m *mockSettingCommandRepo) Update(ctx context.Context, s *setting.Setting) error {
	if m.updateError != nil {
		return m.updateError
	}
	return nil
}

func (m *mockSettingCommandRepo) Delete(ctx context.Context, id uint) error {
	return nil
}

func (m *mockSettingCommandRepo) BatchUpsert(ctx context.Context, settings []*setting.Setting) error {
	return nil
}

// mockSettingQueryRepo 是设置查询仓储的 Mock 实现。
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

type updateSettingTestSetup struct {
	handler     *UpdateSettingHandler
	commandRepo *mockSettingCommandRepo
	queryRepo   *mockSettingQueryRepo
}

func setupUpdateSettingTest() *updateSettingTestSetup {
	commandRepo := newMockSettingCommandRepo()
	queryRepo := newMockSettingQueryRepo()
	handler := NewUpdateSettingHandler(commandRepo, queryRepo)

	return &updateSettingTestSetup{
		handler:     handler,
		commandRepo: commandRepo,
		queryRepo:   queryRepo,
	}
}

func TestUpdateSettingHandler_Handle(t *testing.T) {
	ctx := context.Background()

	t.Run("成功更新设置值", func(t *testing.T) {
		setup := setupUpdateSettingTest()

		// 准备现有设置
		setup.queryRepo.settings["site.name"] = &setting.Setting{
			ID:        1,
			Key:       "site.name",
			Value:     "Old Site Name",
			Category:  setting.CategoryGeneral,
			ValueType: setting.ValueTypeString,
			Label:     "Site Name",
		}

		cmd := UpdateSettingCommand{
			Key:   "site.name",
			Value: "New Site Name",
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "New Site Name", result.Value)
		assert.Equal(t, "site.name", result.Key)
	})

	t.Run("更新设置的多个字段", func(t *testing.T) {
		setup := setupUpdateSettingTest()

		setup.queryRepo.settings["app.debug"] = &setting.Setting{
			ID:        2,
			Key:       "app.debug",
			Value:     "false",
			ValueType: setting.ValueTypeBoolean,
			Label:     "Debug Mode",
		}

		cmd := UpdateSettingCommand{
			Key:       "app.debug",
			Value:     "true",
			ValueType: setting.ValueTypeBoolean,
			Label:     "Enable Debug Mode",
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "true", result.Value)
		assert.Equal(t, setting.ValueTypeBoolean, result.ValueType)
		assert.Equal(t, "Enable Debug Mode", result.Label)
	})

	t.Run("设置不存在", func(t *testing.T) {
		setup := setupUpdateSettingTest()

		cmd := UpdateSettingCommand{
			Key:   "nonexistent.key",
			Value: "some value",
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "setting not found")
		assert.Nil(t, result)
	})

	t.Run("查询设置失败", func(t *testing.T) {
		setup := setupUpdateSettingTest()
		setup.queryRepo.findError = errors.New("database error")

		cmd := UpdateSettingCommand{
			Key:   "any.key",
			Value: "any value",
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to find setting")
		assert.Nil(t, result)
	})

	t.Run("更新设置失败", func(t *testing.T) {
		setup := setupUpdateSettingTest()
		setup.commandRepo.updateError = errors.New("update failed")

		setup.queryRepo.settings["test.key"] = &setting.Setting{
			ID:    1,
			Key:   "test.key",
			Value: "old",
		}

		cmd := UpdateSettingCommand{
			Key:   "test.key",
			Value: "new",
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to update setting")
		assert.Nil(t, result)
	})

	t.Run("只更新值，不更新类型和标签", func(t *testing.T) {
		setup := setupUpdateSettingTest()

		setup.queryRepo.settings["partial.update"] = &setting.Setting{
			ID:        1,
			Key:       "partial.update",
			Value:     "old value",
			ValueType: setting.ValueTypeString,
			Label:     "Original Label",
		}

		cmd := UpdateSettingCommand{
			Key:       "partial.update",
			Value:     "new value",
			ValueType: "", // 空字符串，不更新
			Label:     "", // 空字符串，不更新
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "new value", result.Value)
		assert.Equal(t, setting.ValueTypeString, result.ValueType)
		assert.Equal(t, "Original Label", result.Label)
	})
}

func TestUpdateSettingHandler_ValidationScenarios(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		cmd       UpdateSettingCommand
		setupFunc func(*updateSettingTestSetup)
		wantErr   bool
		errMsg    string
	}{
		{
			name: "有效更新",
			cmd: UpdateSettingCommand{
				Key:   "valid.key",
				Value: "valid value",
			},
			setupFunc: func(s *updateSettingTestSetup) {
				s.queryRepo.settings["valid.key"] = &setting.Setting{
					ID:  1,
					Key: "valid.key",
				}
			},
			wantErr: false,
		},
		{
			name: "更新不存在的设置",
			cmd: UpdateSettingCommand{
				Key:   "missing.key",
				Value: "value",
			},
			wantErr: true,
			errMsg:  "setting not found",
		},
		{
			name: "更新 JSON 类型的设置",
			cmd: UpdateSettingCommand{
				Key:       "json.config",
				Value:     `{"enabled": true, "limit": 100}`,
				ValueType: setting.ValueTypeJSON,
			},
			setupFunc: func(s *updateSettingTestSetup) {
				s.queryRepo.settings["json.config"] = &setting.Setting{
					ID:        1,
					Key:       "json.config",
					ValueType: setting.ValueTypeJSON,
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup := setupUpdateSettingTest()
			if tt.setupFunc != nil {
				tt.setupFunc(setup)
			}

			result, err := setup.handler.Handle(ctx, tt.cmd)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestNewUpdateSettingHandler(t *testing.T) {
	commandRepo := newMockSettingCommandRepo()
	queryRepo := newMockSettingQueryRepo()

	handler := NewUpdateSettingHandler(commandRepo, queryRepo)

	assert.NotNil(t, handler)
}
