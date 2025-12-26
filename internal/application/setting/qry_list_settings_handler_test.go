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

func TestListSettingsHandler_Handle_Success(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name     string
		query    ListSettingsQuery
		settings []*setting.Setting
		wantLen  int
	}{
		{
			name:  "获取所有设置",
			query: ListSettingsQuery{},
			settings: []*setting.Setting{
				{
					ID:        1,
					Key:       "site.name",
					Value:     "My Site",
					Category:  setting.CategoryGeneral,
					ValueType: setting.ValueTypeString,
					Label:     "Site Name",
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:        2,
					Key:       "security.2fa_enabled",
					Value:     "true",
					Category:  setting.CategorySecurity,
					ValueType: setting.ValueTypeBoolean,
					Label:     "Enable 2FA",
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
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
			wantLen: 3,
		},
		{
			name: "按分类获取设置",
			query: ListSettingsQuery{
				Category: setting.CategoryGeneral,
			},
			settings: []*setting.Setting{
				{
					ID:        1,
					Key:       "site.name",
					Value:     "My Site",
					Category:  setting.CategoryGeneral,
					ValueType: setting.ValueTypeString,
					Label:     "Site Name",
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:        2,
					Key:       "site.logo",
					Value:     "logo.png",
					Category:  setting.CategoryGeneral,
					ValueType: setting.ValueTypeString,
					Label:     "Site Logo",
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			wantLen: 2,
		},
		{
			name: "获取安全类设置",
			query: ListSettingsQuery{
				Category: setting.CategorySecurity,
			},
			settings: []*setting.Setting{
				{
					ID:        1,
					Key:       "security.2fa_enabled",
					Value:     "true",
					Category:  setting.CategorySecurity,
					ValueType: setting.ValueTypeBoolean,
					Label:     "Enable 2FA",
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			wantLen: 1,
		},
		{
			name:     "空结果集",
			query:    ListSettingsQuery{},
			settings: []*setting.Setting{},
			wantLen:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockQryRepo := new(MockSettingQueryRepository)

			if tt.query.Category != "" {
				mockQryRepo.On("FindByCategory", mock.Anything, tt.query.Category).Return(tt.settings, nil)
			} else {
				mockQryRepo.On("FindAll", mock.Anything).Return(tt.settings, nil)
			}

			handler := NewListSettingsHandler(mockQryRepo)

			// Act
			result, err := handler.Handle(context.Background(), tt.query)

			// Assert
			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.Len(t, result, tt.wantLen)

			// 验证 DTO 转换正确
			for i, dto := range result {
				assert.Equal(t, tt.settings[i].ID, dto.ID)
				assert.Equal(t, tt.settings[i].Key, dto.Key)
				assert.Equal(t, tt.settings[i].Value, dto.Value)
				assert.Equal(t, tt.settings[i].Category, dto.Category)
				assert.Equal(t, tt.settings[i].ValueType, dto.ValueType)
				assert.Equal(t, tt.settings[i].Label, dto.Label)
			}

			mockQryRepo.AssertExpectations(t)
		})
	}
}

func TestListSettingsHandler_Handle_Error(t *testing.T) {
	tests := []struct {
		name       string
		query      ListSettingsQuery
		setupMocks func(*MockSettingQueryRepository)
		wantErr    string
	}{
		{
			name:  "获取所有设置时数据库错误",
			query: ListSettingsQuery{},
			setupMocks: func(qryRepo *MockSettingQueryRepository) {
				qryRepo.On("FindAll", mock.Anything).Return(nil, errors.New("database error"))
			},
			wantErr: "failed to fetch settings",
		},
		{
			name: "按分类获取设置时数据库错误",
			query: ListSettingsQuery{
				Category: setting.CategoryGeneral,
			},
			setupMocks: func(qryRepo *MockSettingQueryRepository) {
				qryRepo.On("FindByCategory", mock.Anything, setting.CategoryGeneral).Return(nil, errors.New("database error"))
			},
			wantErr: "failed to fetch settings",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockQryRepo := new(MockSettingQueryRepository)
			tt.setupMocks(mockQryRepo)

			handler := NewListSettingsHandler(mockQryRepo)

			// Act
			result, err := handler.Handle(context.Background(), tt.query)

			// Assert
			assert.Nil(t, result)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestListSettingsHandler_Handle_EmptyCategory(t *testing.T) {
	// 验证空分类调用 FindAll
	mockQryRepo := new(MockSettingQueryRepository)

	mockQryRepo.On("FindAll", mock.Anything).Return([]*setting.Setting{}, nil)

	handler := NewListSettingsHandler(mockQryRepo)
	_, err := handler.Handle(context.Background(), ListSettingsQuery{Category: ""})

	require.NoError(t, err)
	mockQryRepo.AssertCalled(t, "FindAll", mock.Anything)
	mockQryRepo.AssertNotCalled(t, "FindByCategory", mock.Anything, mock.Anything)
}

func TestListSettingsHandler_Handle_WithCategory(t *testing.T) {
	// 验证指定分类调用 FindByCategory
	mockQryRepo := new(MockSettingQueryRepository)

	mockQryRepo.On("FindByCategory", mock.Anything, setting.CategorySecurity).Return([]*setting.Setting{}, nil)

	handler := NewListSettingsHandler(mockQryRepo)
	_, err := handler.Handle(context.Background(), ListSettingsQuery{Category: setting.CategorySecurity})

	require.NoError(t, err)
	mockQryRepo.AssertCalled(t, "FindByCategory", mock.Anything, setting.CategorySecurity)
	mockQryRepo.AssertNotCalled(t, "FindAll", mock.Anything)
}
