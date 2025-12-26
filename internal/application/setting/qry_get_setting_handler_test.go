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

func TestGetSettingHandler_Handle_Success(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name    string
		query   GetSettingQuery
		setting *setting.Setting
		want    *SettingDTO
	}{
		{
			name: "成功获取字符串类型设置",
			query: GetSettingQuery{
				Key: "site.name",
			},
			setting: &setting.Setting{
				ID:        1,
				Key:       "site.name",
				Value:     "My Site",
				Category:  setting.CategoryGeneral,
				ValueType: setting.ValueTypeString,
				Label:     "Site Name",
				CreatedAt: now,
				UpdatedAt: now,
			},
			want: &SettingDTO{
				ID:        1,
				Key:       "site.name",
				Value:     "My Site",
				Category:  setting.CategoryGeneral,
				ValueType: setting.ValueTypeString,
				Label:     "Site Name",
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		{
			name: "获取布尔类型设置",
			query: GetSettingQuery{
				Key: "security.2fa_enabled",
			},
			setting: &setting.Setting{
				ID:        2,
				Key:       "security.2fa_enabled",
				Value:     "true",
				Category:  setting.CategorySecurity,
				ValueType: setting.ValueTypeBoolean,
				Label:     "Enable 2FA",
				CreatedAt: now,
				UpdatedAt: now,
			},
			want: &SettingDTO{
				ID:        2,
				Key:       "security.2fa_enabled",
				Value:     "true",
				Category:  setting.CategorySecurity,
				ValueType: setting.ValueTypeBoolean,
				Label:     "Enable 2FA",
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		{
			name: "获取数字类型设置",
			query: GetSettingQuery{
				Key: "backup.retention_days",
			},
			setting: &setting.Setting{
				ID:        3,
				Key:       "backup.retention_days",
				Value:     "30",
				Category:  setting.CategoryBackup,
				ValueType: setting.ValueTypeNumber,
				Label:     "Backup Retention Days",
				CreatedAt: now,
				UpdatedAt: now,
			},
			want: &SettingDTO{
				ID:        3,
				Key:       "backup.retention_days",
				Value:     "30",
				Category:  setting.CategoryBackup,
				ValueType: setting.ValueTypeNumber,
				Label:     "Backup Retention Days",
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockQryRepo := new(MockSettingQueryRepository)

			mockQryRepo.On("FindByKey", mock.Anything, tt.query.Key).Return(tt.setting, nil)

			handler := NewGetSettingHandler(mockQryRepo)

			// Act
			result, err := handler.Handle(context.Background(), tt.query)

			// Assert
			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tt.want.ID, result.ID)
			assert.Equal(t, tt.want.Key, result.Key)
			assert.Equal(t, tt.want.Value, result.Value)
			assert.Equal(t, tt.want.Category, result.Category)
			assert.Equal(t, tt.want.ValueType, result.ValueType)
			assert.Equal(t, tt.want.Label, result.Label)
			mockQryRepo.AssertExpectations(t)
		})
	}
}

func TestGetSettingHandler_Handle_Error(t *testing.T) {
	tests := []struct {
		name       string
		query      GetSettingQuery
		setupMocks func(*MockSettingQueryRepository)
		wantErr    string
	}{
		{
			name: "设置不存在",
			query: GetSettingQuery{
				Key: "nonexistent.key",
			},
			setupMocks: func(qryRepo *MockSettingQueryRepository) {
				qryRepo.On("FindByKey", mock.Anything, "nonexistent.key").Return(nil, nil)
			},
			wantErr: "setting not found",
		},
		{
			name: "查询设置时数据库错误",
			query: GetSettingQuery{
				Key: "test.key",
			},
			setupMocks: func(qryRepo *MockSettingQueryRepository) {
				qryRepo.On("FindByKey", mock.Anything, "test.key").Return(nil, errors.New("database error"))
			},
			wantErr: "failed to find setting",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockQryRepo := new(MockSettingQueryRepository)
			tt.setupMocks(mockQryRepo)

			handler := NewGetSettingHandler(mockQryRepo)

			// Act
			result, err := handler.Handle(context.Background(), tt.query)

			// Assert
			assert.Nil(t, result)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}
