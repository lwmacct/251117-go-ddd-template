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

func TestCreateSettingHandler_Handle_Success(t *testing.T) {
	tests := []struct {
		name string
		cmd  CreateSettingCommand
		want *CreateSettingResultDTO
	}{
		{
			name: "成功创建字符串类型设置",
			cmd: CreateSettingCommand{
				Key:       "site.name",
				Value:     "My Site",
				Category:  setting.CategoryGeneral,
				ValueType: setting.ValueTypeString,
				Label:     "Site Name",
			},
			want: &CreateSettingResultDTO{
				ID: 1,
			},
		},
		{
			name: "创建布尔类型设置",
			cmd: CreateSettingCommand{
				Key:       "security.2fa_required",
				Value:     "true",
				Category:  setting.CategorySecurity,
				ValueType: setting.ValueTypeBoolean,
				Label:     "Require 2FA",
			},
			want: &CreateSettingResultDTO{
				ID: 1,
			},
		},
		{
			name: "创建设置时未指定类型，使用默认值",
			cmd: CreateSettingCommand{
				Key:      "test.key",
				Value:    "test value",
				Category: setting.CategoryGeneral,
				Label:    "Test Setting",
			},
			want: &CreateSettingResultDTO{
				ID: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockCmdRepo := new(MockSettingCommandRepository)
			mockQryRepo := new(MockSettingQueryRepository)

			mockQryRepo.On("FindByKey", mock.Anything, tt.cmd.Key).Return(nil, nil)
			mockCmdRepo.On("Create", mock.Anything, mock.AnythingOfType("*setting.Setting")).Return(nil)

			handler := NewCreateSettingHandler(mockCmdRepo, mockQryRepo)

			// Act
			result, err := handler.Handle(context.Background(), tt.cmd)

			// Assert
			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tt.want.ID, result.ID)
			mockCmdRepo.AssertExpectations(t)
			mockQryRepo.AssertExpectations(t)
		})
	}
}

func TestCreateSettingHandler_Handle_Error(t *testing.T) {
	tests := []struct {
		name       string
		cmd        CreateSettingCommand
		setupMocks func(*MockSettingCommandRepository, *MockSettingQueryRepository)
		wantErr    string
	}{
		{
			name: "Key 已存在",
			cmd: CreateSettingCommand{
				Key:      "existing.key",
				Value:    "value",
				Category: setting.CategoryGeneral,
			},
			setupMocks: func(cmdRepo *MockSettingCommandRepository, qryRepo *MockSettingQueryRepository) {
				qryRepo.On("FindByKey", mock.Anything, "existing.key").Return(&setting.Setting{
					ID:  1,
					Key: "existing.key",
				}, nil)
			},
			wantErr: "setting key already exists",
		},
		{
			name: "查询 Key 时数据库错误",
			cmd: CreateSettingCommand{
				Key:      "test.key",
				Value:    "value",
				Category: setting.CategoryGeneral,
			},
			setupMocks: func(cmdRepo *MockSettingCommandRepository, qryRepo *MockSettingQueryRepository) {
				qryRepo.On("FindByKey", mock.Anything, "test.key").Return(nil, errors.New("database error"))
			},
			wantErr: "failed to check existing setting",
		},
		{
			name: "创建设置时数据库错误",
			cmd: CreateSettingCommand{
				Key:      "test.key",
				Value:    "value",
				Category: setting.CategoryGeneral,
			},
			setupMocks: func(cmdRepo *MockSettingCommandRepository, qryRepo *MockSettingQueryRepository) {
				qryRepo.On("FindByKey", mock.Anything, "test.key").Return(nil, nil)
				cmdRepo.On("Create", mock.Anything, mock.AnythingOfType("*setting.Setting")).Return(errors.New("database error"))
			},
			wantErr: "failed to create setting",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockCmdRepo := new(MockSettingCommandRepository)
			mockQryRepo := new(MockSettingQueryRepository)
			tt.setupMocks(mockCmdRepo, mockQryRepo)

			handler := NewCreateSettingHandler(mockCmdRepo, mockQryRepo)

			// Act
			result, err := handler.Handle(context.Background(), tt.cmd)

			// Assert
			assert.Nil(t, result)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestCreateSettingHandler_Handle_DefaultValueType(t *testing.T) {
	// 验证默认值类型为 string
	mockCmdRepo := new(MockSettingCommandRepository)
	mockQryRepo := new(MockSettingQueryRepository)

	var capturedSetting *setting.Setting
	mockQryRepo.On("FindByKey", mock.Anything, mock.Anything).Return(nil, nil)
	mockCmdRepo.On("Create", mock.Anything, mock.AnythingOfType("*setting.Setting")).
		Run(func(args mock.Arguments) {
			capturedSetting = args.Get(1).(*setting.Setting)
		}).Return(nil)

	handler := NewCreateSettingHandler(mockCmdRepo, mockQryRepo)
	_, err := handler.Handle(context.Background(), CreateSettingCommand{
		Key:      "test.key",
		Value:    "test value",
		Category: setting.CategoryGeneral,
	})

	require.NoError(t, err)
	assert.NotNil(t, capturedSetting)
	assert.Equal(t, setting.ValueTypeString, capturedSetting.ValueType)
}
