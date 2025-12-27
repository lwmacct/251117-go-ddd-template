//nolint:forcetypeassert // 测试中的类型断言是可控的
package setting

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

func TestBatchUpdateHandler_Handle_Success(t *testing.T) {
	tests := []struct {
		name             string
		cmd              BatchUpdateCommand
		existingSettings []*setting.Setting
	}{
		{
			name: "成功批量更新多个设置",
			cmd: BatchUpdateCommand{
				Settings: []SettingItemCommand{
					{Key: "site.name", Value: "New Site"},
					{Key: "site.logo", Value: "new-logo.png"},
					{Key: "site.description", Value: "New Description"},
				},
			},
			existingSettings: []*setting.Setting{
				{ID: 1, Key: "site.name", Value: "Old Site", Category: setting.CategoryGeneral},
				{ID: 2, Key: "site.logo", Value: "old-logo.png", Category: setting.CategoryGeneral},
				{ID: 3, Key: "site.description", Value: "Old Description", Category: setting.CategoryGeneral},
			},
		},
		{
			name: "批量更新单个设置",
			cmd: BatchUpdateCommand{
				Settings: []SettingItemCommand{
					{Key: "test.key", Value: "updated value"},
				},
			},
			existingSettings: []*setting.Setting{
				{ID: 1, Key: "test.key", Value: "old value", Category: setting.CategoryGeneral},
			},
		},
		{
			name: "空设置列表不执行更新",
			cmd: BatchUpdateCommand{
				Settings: []SettingItemCommand{},
			},
			existingSettings: []*setting.Setting{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockCmdRepo := new(MockSettingCommandRepository)
			mockQryRepo := new(MockSettingQueryRepository)

			if len(tt.cmd.Settings) > 0 {
				keys := make([]string, len(tt.cmd.Settings))
				for i, item := range tt.cmd.Settings {
					keys[i] = item.Key
				}
				mockQryRepo.On("FindByKeys", mock.Anything, keys).Return(tt.existingSettings, nil)
				mockCmdRepo.On("BatchUpsert", mock.Anything, mock.AnythingOfType("[]*setting.Setting")).Return(nil)
			}

			handler := NewBatchUpdateHandler(mockCmdRepo, mockQryRepo, nil)

			// Act
			err := handler.Handle(context.Background(), tt.cmd)

			// Assert
			require.NoError(t, err)
			mockCmdRepo.AssertExpectations(t)
			mockQryRepo.AssertExpectations(t)
		})
	}
}

func TestBatchUpdateHandler_Handle_Error(t *testing.T) {
	tests := []struct {
		name       string
		cmd        BatchUpdateCommand
		setupMocks func(*MockSettingCommandRepository, *MockSettingQueryRepository)
		wantErr    string
	}{
		{
			name: "查询设置时数据库错误",
			cmd: BatchUpdateCommand{
				Settings: []SettingItemCommand{
					{Key: "test.key", Value: "value"},
				},
			},
			setupMocks: func(cmdRepo *MockSettingCommandRepository, qryRepo *MockSettingQueryRepository) {
				qryRepo.On("FindByKeys", mock.Anything, []string{"test.key"}).Return(nil, errors.New("database error"))
			},
			wantErr: "failed to find settings",
		},
		{
			name: "部分设置 Key 不存在",
			cmd: BatchUpdateCommand{
				Settings: []SettingItemCommand{
					{Key: "existing.key", Value: "value1"},
					{Key: "nonexistent.key", Value: "value2"},
				},
			},
			setupMocks: func(cmdRepo *MockSettingCommandRepository, qryRepo *MockSettingQueryRepository) {
				qryRepo.On("FindByKeys", mock.Anything, []string{"existing.key", "nonexistent.key"}).
					Return([]*setting.Setting{
						{ID: 1, Key: "existing.key", Value: "old value"},
					}, nil)
			},
			wantErr: "setting key nonexistent.key does not exist",
		},
		{
			name: "批量更新时数据库错误",
			cmd: BatchUpdateCommand{
				Settings: []SettingItemCommand{
					{Key: "test.key", Value: "new value"},
				},
			},
			setupMocks: func(cmdRepo *MockSettingCommandRepository, qryRepo *MockSettingQueryRepository) {
				qryRepo.On("FindByKeys", mock.Anything, []string{"test.key"}).
					Return([]*setting.Setting{
						{ID: 1, Key: "test.key", Value: "old value"},
					}, nil)
				cmdRepo.On("BatchUpsert", mock.Anything, mock.AnythingOfType("[]*setting.Setting")).Return(errors.New("database error"))
			},
			wantErr: "failed to batch update settings",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockCmdRepo := new(MockSettingCommandRepository)
			mockQryRepo := new(MockSettingQueryRepository)
			tt.setupMocks(mockCmdRepo, mockQryRepo)

			handler := NewBatchUpdateHandler(mockCmdRepo, mockQryRepo, nil)

			// Act
			err := handler.Handle(context.Background(), tt.cmd)

			// Assert
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestBatchUpdateHandler_Handle_UpdatesCorrectValues(t *testing.T) {
	// 验证批量更新正确设置了新值
	mockCmdRepo := new(MockSettingCommandRepository)
	mockQryRepo := new(MockSettingQueryRepository)

	cmd := BatchUpdateCommand{
		Settings: []SettingItemCommand{
			{Key: "key1", Value: "new value 1"},
			{Key: "key2", Value: "new value 2"},
		},
	}

	existingSettings := []*setting.Setting{
		{ID: 1, Key: "key1", Value: "old value 1", Category: setting.CategoryGeneral},
		{ID: 2, Key: "key2", Value: "old value 2", Category: setting.CategoryGeneral},
	}

	var capturedSettings []*setting.Setting
	mockQryRepo.On("FindByKeys", mock.Anything, []string{"key1", "key2"}).Return(existingSettings, nil)
	mockCmdRepo.On("BatchUpsert", mock.Anything, mock.AnythingOfType("[]*setting.Setting")).
		Run(func(args mock.Arguments) {
			capturedSettings = args.Get(1).([]*setting.Setting)
		}).Return(nil)

	handler := NewBatchUpdateHandler(mockCmdRepo, mockQryRepo, nil)
	err := handler.Handle(context.Background(), cmd)

	require.NoError(t, err)
	assert.Len(t, capturedSettings, 2)

	// 验证值已更新
	for _, s := range capturedSettings {
		switch s.Key {
		case "key1":
			assert.Equal(t, "new value 1", s.Value)
		case "key2":
			assert.Equal(t, "new value 2", s.Value)
		}
	}
}

func TestBatchUpdateHandler_Handle_WithValidation(t *testing.T) {
	// 测试带验证器的批量更新
	mockCmdRepo := new(MockSettingCommandRepository)
	mockQryRepo := new(MockSettingQueryRepository)
	mockValidator := new(MockSettingValidator)

	cmd := BatchUpdateCommand{
		Settings: []SettingItemCommand{
			{Key: "password.min_length", Value: "8"},
		},
	}

	existingSettings := []*setting.Setting{
		{
			ID: 1, Key: "password.min_length", Value: "6",
			ValueType: setting.ValueTypeNumber,
			UIConfig:  `{"input_type":"number","validation":{"min":6,"max":32}}`,
		},
	}

	// Mock: 查询设置
	mockQryRepo.On("FindByKeys", mock.Anything, []string{"password.min_length"}).Return(existingSettings, nil)
	// Mock: 查询所有设置（用于跨字段验证）
	mockQryRepo.On("FindAll", mock.Anything).Return(existingSettings, nil)
	// Mock: 验证通过
	mockValidator.On("ValidateBatch", mock.Anything, mock.AnythingOfType("[]*setting.ValidationContext")).
		Return(map[string]string(nil), nil)
	// Mock: 批量更新
	mockCmdRepo.On("BatchUpsert", mock.Anything, mock.AnythingOfType("[]*setting.Setting")).Return(nil)

	handler := NewBatchUpdateHandler(mockCmdRepo, mockQryRepo, mockValidator)
	err := handler.Handle(context.Background(), cmd)

	require.NoError(t, err)
	mockValidator.AssertExpectations(t)
}

func TestBatchUpdateHandler_Handle_ValidationFailed(t *testing.T) {
	// 测试验证失败场景
	mockCmdRepo := new(MockSettingCommandRepository)
	mockQryRepo := new(MockSettingQueryRepository)
	mockValidator := new(MockSettingValidator)

	cmd := BatchUpdateCommand{
		Settings: []SettingItemCommand{
			{Key: "password.min_length", Value: "3"}, // 小于最小值 6
		},
	}

	existingSettings := []*setting.Setting{
		{
			ID: 1, Key: "password.min_length", Value: "8",
			ValueType: setting.ValueTypeNumber,
			UIConfig:  `{"input_type":"number","validation":{"min":6,"max":32}}`,
		},
	}

	// Mock: 查询设置
	mockQryRepo.On("FindByKeys", mock.Anything, []string{"password.min_length"}).Return(existingSettings, nil)
	// Mock: 查询所有设置
	mockQryRepo.On("FindAll", mock.Anything).Return(existingSettings, nil)
	// Mock: 验证失败
	mockValidator.On("ValidateBatch", mock.Anything, mock.AnythingOfType("[]*setting.ValidationContext")).
		Return(map[string]string{"password.min_length": "最小值为 6"}, nil)

	handler := NewBatchUpdateHandler(mockCmdRepo, mockQryRepo, mockValidator)
	err := handler.Handle(context.Background(), cmd)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")
	// 验证未调用 BatchUpsert
	mockCmdRepo.AssertNotCalled(t, "BatchUpsert")
}
