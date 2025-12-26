package role

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/role"
)

func TestUpdateRoleHandler_Handle_Success(t *testing.T) {
	displayName := "更新后的显示名"
	description := "更新后的描述"

	tests := []struct {
		name         string
		cmd          UpdateRoleCommand
		existingRole *role.Role
		want         *RoleDTO
	}{
		{
			name: "更新显示名和描述",
			cmd: UpdateRoleCommand{
				RoleID:      1,
				DisplayName: &displayName,
				Description: &description,
			},
			existingRole: &role.Role{
				ID:          1,
				Name:        "developer",
				DisplayName: "开发者",
				Description: "原描述",
				IsSystem:    false,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			want: &RoleDTO{
				ID:          1,
				Name:        "developer",
				DisplayName: "更新后的显示名",
				Description: "更新后的描述",
				IsSystem:    false,
			},
		},
		{
			name: "仅更新显示名",
			cmd: UpdateRoleCommand{
				RoleID:      2,
				DisplayName: &displayName,
			},
			existingRole: &role.Role{
				ID:          2,
				Name:        "viewer",
				DisplayName: "查看者",
				Description: "原描述",
				IsSystem:    false,
			},
			want: &RoleDTO{
				ID:          2,
				Name:        "viewer",
				DisplayName: "更新后的显示名",
				Description: "原描述",
				IsSystem:    false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockCmdRepo := new(MockRoleCommandRepository)
			mockQryRepo := new(MockRoleQueryRepository)

			mockQryRepo.On("FindByID", mock.Anything, tt.cmd.RoleID).Return(tt.existingRole, nil)
			mockCmdRepo.On("Update", mock.Anything, mock.AnythingOfType("*role.Role")).Return(nil)

			handler := NewUpdateRoleHandler(mockCmdRepo, mockQryRepo)

			// Act
			result, err := handler.Handle(context.Background(), tt.cmd)

			// Assert
			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tt.want.DisplayName, result.DisplayName)
			mockCmdRepo.AssertExpectations(t)
			mockQryRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateRoleHandler_Handle_Error(t *testing.T) {
	displayName := "新名称"

	tests := []struct {
		name       string
		cmd        UpdateRoleCommand
		setupMocks func(*MockRoleCommandRepository, *MockRoleQueryRepository)
		wantErr    string
	}{
		{
			name: "角色不存在",
			cmd:  UpdateRoleCommand{RoleID: 999, DisplayName: &displayName},
			setupMocks: func(cmdRepo *MockRoleCommandRepository, qryRepo *MockRoleQueryRepository) {
				qryRepo.On("FindByID", mock.Anything, uint(999)).Return(nil, nil)
			},
			wantErr: "role not found",
		},
		{
			name: "系统角色不可修改",
			cmd:  UpdateRoleCommand{RoleID: 1, DisplayName: &displayName},
			setupMocks: func(cmdRepo *MockRoleCommandRepository, qryRepo *MockRoleQueryRepository) {
				qryRepo.On("FindByID", mock.Anything, uint(1)).Return(&role.Role{
					ID:       1,
					Name:     "admin",
					IsSystem: true,
				}, nil)
			},
			wantErr: "cannot modify system role",
		},
		{
			name: "查询角色时数据库错误",
			cmd:  UpdateRoleCommand{RoleID: 1, DisplayName: &displayName},
			setupMocks: func(cmdRepo *MockRoleCommandRepository, qryRepo *MockRoleQueryRepository) {
				qryRepo.On("FindByID", mock.Anything, uint(1)).Return(nil, errors.New("database error"))
			},
			wantErr: "failed to find role",
		},
		{
			name: "更新角色时数据库错误",
			cmd:  UpdateRoleCommand{RoleID: 1, DisplayName: &displayName},
			setupMocks: func(cmdRepo *MockRoleCommandRepository, qryRepo *MockRoleQueryRepository) {
				qryRepo.On("FindByID", mock.Anything, uint(1)).Return(&role.Role{
					ID:       1,
					Name:     "custom",
					IsSystem: false,
				}, nil)
				cmdRepo.On("Update", mock.Anything, mock.AnythingOfType("*role.Role")).
					Return(errors.New("update failed"))
			},
			wantErr: "failed to update role",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockCmdRepo := new(MockRoleCommandRepository)
			mockQryRepo := new(MockRoleQueryRepository)
			tt.setupMocks(mockCmdRepo, mockQryRepo)

			handler := NewUpdateRoleHandler(mockCmdRepo, mockQryRepo)

			// Act
			result, err := handler.Handle(context.Background(), tt.cmd)

			// Assert
			assert.Nil(t, result)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
			mockCmdRepo.AssertExpectations(t)
			mockQryRepo.AssertExpectations(t)
		})
	}
}
