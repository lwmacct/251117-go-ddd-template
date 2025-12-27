//nolint:forcetypeassert // 测试中的类型断言是可控的
package setting

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

func TestUpdateHandler_Handle_Success(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name            string
		cmd             UpdateCommand
		existingSetting *setting.Setting
		want            *SettingDTO
	}{
		{
			name: "成功更新设置值",
			cmd: UpdateCommand{
				Key:   "site.name",
				Value: "New Site Name",
			},
			existingSetting: &setting.Setting{
				ID:        1,
				Key:       "site.name",
				Value:     "Old Site Name",
				Category:  setting.CategoryGeneral,
				ValueType: setting.ValueTypeString,
				Label:     "Site Name",
				CreatedAt: now,
				UpdatedAt: now,
			},
			want: &SettingDTO{
				ID:        1,
				Key:       "site.name",
				Value:     "New Site Name",
				Category:  setting.CategoryGeneral,
				ValueType: setting.ValueTypeString,
				Label:     "Site Name",
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		{
			name: "更新值和类型",
			cmd: UpdateCommand{
				Key:       "test.setting",
				Value:     "true",
				ValueType: setting.ValueTypeBoolean,
			},
			existingSetting: &setting.Setting{
				ID:        2,
				Key:       "test.setting",
				Value:     "123",
				Category:  setting.CategoryGeneral,
				ValueType: setting.ValueTypeNumber,
				Label:     "Test Setting",
				CreatedAt: now,
				UpdatedAt: now,
			},
			want: &SettingDTO{
				ID:        2,
				Key:       "test.setting",
				Value:     "true",
				Category:  setting.CategoryGeneral,
				ValueType: setting.ValueTypeBoolean,
				Label:     "Test Setting",
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		{
			name: "更新值和标签",
			cmd: UpdateCommand{
				Key:   "test.key",
				Value: "new value",
				Label: "New Label",
			},
			existingSetting: &setting.Setting{
				ID:        3,
				Key:       "test.key",
				Value:     "old value",
				Category:  setting.CategoryGeneral,
				ValueType: setting.ValueTypeString,
				Label:     "Old Label",
				CreatedAt: now,
				UpdatedAt: now,
			},
			want: &SettingDTO{
				ID:        3,
				Key:       "test.key",
				Value:     "new value",
				Category:  setting.CategoryGeneral,
				ValueType: setting.ValueTypeString,
				Label:     "New Label",
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockCmdRepo := new(MockSettingCommandRepository)
			mockQryRepo := new(MockSettingQueryRepository)

			mockQryRepo.On("FindByKey", mock.Anything, tt.cmd.Key).Return(tt.existingSetting, nil)
			mockCmdRepo.On("Update", mock.Anything, mock.AnythingOfType("*setting.Setting")).Return(nil)

			handler := NewUpdateHandler(mockCmdRepo, mockQryRepo, nil)

			// Act
			result, err := handler.Handle(context.Background(), tt.cmd)

			// Assert
			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tt.want.Key, result.Key)
			assert.Equal(t, tt.want.Value, result.Value)
			assert.Equal(t, tt.want.ValueType, result.ValueType)
			assert.Equal(t, tt.want.Label, result.Label)
			mockCmdRepo.AssertExpectations(t)
			mockQryRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateHandler_Handle_Error(t *testing.T) {
	tests := []struct {
		name       string
		cmd        UpdateCommand
		setupMocks func(*MockSettingCommandRepository, *MockSettingQueryRepository)
		wantErr    string
	}{
		{
			name: "设置不存在",
			cmd: UpdateCommand{
				Key:   "nonexistent.key",
				Value: "value",
			},
			setupMocks: func(cmdRepo *MockSettingCommandRepository, qryRepo *MockSettingQueryRepository) {
				qryRepo.On("FindByKey", mock.Anything, "nonexistent.key").Return(nil, nil)
			},
			wantErr: "setting not found",
		},
		{
			name: "查询设置时数据库错误",
			cmd: UpdateCommand{
				Key:   "test.key",
				Value: "value",
			},
			setupMocks: func(cmdRepo *MockSettingCommandRepository, qryRepo *MockSettingQueryRepository) {
				qryRepo.On("FindByKey", mock.Anything, "test.key").Return(nil, errors.New("database error"))
			},
			wantErr: "failed to find setting",
		},
		{
			name: "更新设置时数据库错误",
			cmd: UpdateCommand{
				Key:   "test.key",
				Value: "new value",
			},
			setupMocks: func(cmdRepo *MockSettingCommandRepository, qryRepo *MockSettingQueryRepository) {
				qryRepo.On("FindByKey", mock.Anything, "test.key").Return(&setting.Setting{
					ID:  1,
					Key: "test.key",
				}, nil)
				cmdRepo.On("Update", mock.Anything, mock.AnythingOfType("*setting.Setting")).Return(errors.New("database error"))
			},
			wantErr: "failed to update setting",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockCmdRepo := new(MockSettingCommandRepository)
			mockQryRepo := new(MockSettingQueryRepository)
			tt.setupMocks(mockCmdRepo, mockQryRepo)

			handler := NewUpdateHandler(mockCmdRepo, mockQryRepo, nil)

			// Act
			result, err := handler.Handle(context.Background(), tt.cmd)

			// Assert
			assert.Nil(t, result)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestUpdateHandler_Handle_PartialUpdate(t *testing.T) {
	// 验证只更新提供的字段
	now := time.Now()
	mockCmdRepo := new(MockSettingCommandRepository)
	mockQryRepo := new(MockSettingQueryRepository)

	existingSetting := &setting.Setting{
		ID:        1,
		Key:       "test.key",
		Value:     "old value",
		Category:  setting.CategoryGeneral,
		ValueType: setting.ValueTypeString,
		Label:     "Old Label",
		CreatedAt: now,
		UpdatedAt: now,
	}

	var capturedSetting *setting.Setting
	mockQryRepo.On("FindByKey", mock.Anything, "test.key").Return(existingSetting, nil)
	mockCmdRepo.On("Update", mock.Anything, mock.AnythingOfType("*setting.Setting")).
		Run(func(args mock.Arguments) {
			capturedSetting = args.Get(1).(*setting.Setting)
		}).Return(nil)

	handler := NewUpdateHandler(mockCmdRepo, mockQryRepo, nil)

	// 只更新 Value，不更新 ValueType 和 Label
	_, err := handler.Handle(context.Background(), UpdateCommand{
		Key:   "test.key",
		Value: "new value",
	})

	require.NoError(t, err)
	assert.NotNil(t, capturedSetting)
	assert.Equal(t, "new value", capturedSetting.Value)
	assert.Equal(t, setting.ValueTypeString, capturedSetting.ValueType) // 保持不变
	assert.Equal(t, "Old Label", capturedSetting.Label)                 // 保持不变
}

func TestUpdateHandler_Handle_WithValidation(t *testing.T) {
	// 测试带验证器的更新（验证通过）
	mockCmdRepo := new(MockSettingCommandRepository)
	mockQryRepo := new(MockSettingQueryRepository)
	mockValidator := new(MockSettingValidator)

	existingSetting := &setting.Setting{
		ID:        1,
		Key:       "password.min_length",
		Value:     "6",
		ValueType: setting.ValueTypeNumber,
		UIConfig:  `{"input_type":"number","validation":{"min":6,"max":32}}`,
	}

	mockQryRepo.On("FindByKey", mock.Anything, "password.min_length").Return(existingSetting, nil)
	mockQryRepo.On("FindAll", mock.Anything).Return([]*setting.Setting{existingSetting}, nil)
	mockValidator.On("Validate", mock.Anything, mock.AnythingOfType("*setting.ValidationContext")).
		Return(&setting.ValidationResult{Valid: true}, nil)
	mockCmdRepo.On("Update", mock.Anything, mock.AnythingOfType("*setting.Setting")).Return(nil)

	handler := NewUpdateHandler(mockCmdRepo, mockQryRepo, mockValidator)
	result, err := handler.Handle(context.Background(), UpdateCommand{
		Key:   "password.min_length",
		Value: "8",
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
	mockValidator.AssertExpectations(t)
}

func TestUpdateHandler_Handle_ValidationFailed(t *testing.T) {
	// 测试验证失败场景
	mockCmdRepo := new(MockSettingCommandRepository)
	mockQryRepo := new(MockSettingQueryRepository)
	mockValidator := new(MockSettingValidator)

	existingSetting := &setting.Setting{
		ID:        1,
		Key:       "password.min_length",
		Value:     "8",
		ValueType: setting.ValueTypeNumber,
		UIConfig:  `{"input_type":"number","validation":{"min":6,"max":32}}`,
	}

	mockQryRepo.On("FindByKey", mock.Anything, "password.min_length").Return(existingSetting, nil)
	mockQryRepo.On("FindAll", mock.Anything).Return([]*setting.Setting{existingSetting}, nil)
	mockValidator.On("Validate", mock.Anything, mock.AnythingOfType("*setting.ValidationContext")).
		Return(&setting.ValidationResult{Valid: false, Message: "最小值为 6"}, nil)

	handler := NewUpdateHandler(mockCmdRepo, mockQryRepo, mockValidator)
	result, err := handler.Handle(context.Background(), UpdateCommand{
		Key:   "password.min_length",
		Value: "3", // 小于最小值
	})

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "validation failed")
	// 验证未调用 Update
	mockCmdRepo.AssertNotCalled(t, "Update")
}

func TestUpdateHandler_Handle_NoValidationRule(t *testing.T) {
	// 测试无验证规则时跳过验证
	mockCmdRepo := new(MockSettingCommandRepository)
	mockQryRepo := new(MockSettingQueryRepository)
	mockValidator := new(MockSettingValidator)

	existingSetting := &setting.Setting{
		ID:        1,
		Key:       "site.name",
		Value:     "Old Site",
		ValueType: setting.ValueTypeString,
		UIConfig:  `{"input_type":"text"}`, // 无验证规则
	}

	mockQryRepo.On("FindByKey", mock.Anything, "site.name").Return(existingSetting, nil)
	mockCmdRepo.On("Update", mock.Anything, mock.AnythingOfType("*setting.Setting")).Return(nil)

	handler := NewUpdateHandler(mockCmdRepo, mockQryRepo, mockValidator)
	result, err := handler.Handle(context.Background(), UpdateCommand{
		Key:   "site.name",
		Value: "New Site",
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
	// 验证未调用 Validate（无验证规则时跳过）
	mockValidator.AssertNotCalled(t, "Validate")
}
