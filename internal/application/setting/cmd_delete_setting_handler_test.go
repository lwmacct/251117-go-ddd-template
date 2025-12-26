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

func TestDeleteSettingHandler_Handle_Success(t *testing.T) {
	tests := []struct {
		name            string
		cmd             DeleteSettingCommand
		existingSetting *setting.Setting
	}{
		{
			name: "成功删除设置",
			cmd: DeleteSettingCommand{
				Key: "test.key",
			},
			existingSetting: &setting.Setting{
				ID:       1,
				Key:      "test.key",
				Value:    "test value",
				Category: setting.CategoryGeneral,
			},
		},
		{
			name: "删除安全类设置",
			cmd: DeleteSettingCommand{
				Key: "security.setting",
			},
			existingSetting: &setting.Setting{
				ID:       2,
				Key:      "security.setting",
				Value:    "secure value",
				Category: setting.CategorySecurity,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockCmdRepo := new(MockSettingCommandRepository)
			mockQryRepo := new(MockSettingQueryRepository)

			mockQryRepo.On("FindByKey", mock.Anything, tt.cmd.Key).Return(tt.existingSetting, nil)
			mockCmdRepo.On("Delete", mock.Anything, tt.existingSetting.ID).Return(nil)

			handler := NewDeleteSettingHandler(mockCmdRepo, mockQryRepo)

			// Act
			err := handler.Handle(context.Background(), tt.cmd)

			// Assert
			require.NoError(t, err)
			mockCmdRepo.AssertExpectations(t)
			mockQryRepo.AssertExpectations(t)
		})
	}
}

func TestDeleteSettingHandler_Handle_Error(t *testing.T) {
	tests := []struct {
		name       string
		cmd        DeleteSettingCommand
		setupMocks func(*MockSettingCommandRepository, *MockSettingQueryRepository)
		wantErr    string
	}{
		{
			name: "设置不存在",
			cmd: DeleteSettingCommand{
				Key: "nonexistent.key",
			},
			setupMocks: func(cmdRepo *MockSettingCommandRepository, qryRepo *MockSettingQueryRepository) {
				qryRepo.On("FindByKey", mock.Anything, "nonexistent.key").Return(nil, nil)
			},
			wantErr: "setting not found",
		},
		{
			name: "查询设置时数据库错误",
			cmd: DeleteSettingCommand{
				Key: "test.key",
			},
			setupMocks: func(cmdRepo *MockSettingCommandRepository, qryRepo *MockSettingQueryRepository) {
				qryRepo.On("FindByKey", mock.Anything, "test.key").Return(nil, errors.New("database error"))
			},
			wantErr: "failed to find setting",
		},
		{
			name: "删除设置时数据库错误",
			cmd: DeleteSettingCommand{
				Key: "test.key",
			},
			setupMocks: func(cmdRepo *MockSettingCommandRepository, qryRepo *MockSettingQueryRepository) {
				qryRepo.On("FindByKey", mock.Anything, "test.key").Return(&setting.Setting{
					ID:  1,
					Key: "test.key",
				}, nil)
				cmdRepo.On("Delete", mock.Anything, uint(1)).Return(errors.New("database error"))
			},
			wantErr: "failed to delete setting",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockCmdRepo := new(MockSettingCommandRepository)
			mockQryRepo := new(MockSettingQueryRepository)
			tt.setupMocks(mockCmdRepo, mockQryRepo)

			handler := NewDeleteSettingHandler(mockCmdRepo, mockQryRepo)

			// Act
			err := handler.Handle(context.Background(), tt.cmd)

			// Assert
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}
