package role

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/role"
)

func TestDeleteRoleHandler_Handle_Success(t *testing.T) {
	tests := []struct {
		name         string
		cmd          DeleteRoleCommand
		existingRole *role.Role
	}{
		{
			name: "成功删除普通角色",
			cmd:  DeleteRoleCommand{RoleID: 1},
			existingRole: &role.Role{
				ID:       1,
				Name:     "developer",
				IsSystem: false,
			},
		},
		{
			name: "删除自定义角色",
			cmd:  DeleteRoleCommand{RoleID: 100},
			existingRole: &role.Role{
				ID:       100,
				Name:     "custom",
				IsSystem: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockCmdRepo := new(MockRoleCommandRepository)
			mockQryRepo := new(MockRoleQueryRepository)

			mockQryRepo.On("FindByID", mock.Anything, tt.cmd.RoleID).Return(tt.existingRole, nil)
			mockCmdRepo.On("Delete", mock.Anything, tt.cmd.RoleID).Return(nil)

			handler := NewDeleteRoleHandler(mockCmdRepo, mockQryRepo)

			// Act
			err := handler.Handle(context.Background(), tt.cmd)

			// Assert
			require.NoError(t, err)
			mockCmdRepo.AssertExpectations(t)
			mockQryRepo.AssertExpectations(t)
		})
	}
}

func TestDeleteRoleHandler_Handle_Error(t *testing.T) {
	tests := []struct {
		name       string
		cmd        DeleteRoleCommand
		setupMocks func(*MockRoleCommandRepository, *MockRoleQueryRepository)
		wantErr    string
	}{
		{
			name: "角色不存在",
			cmd:  DeleteRoleCommand{RoleID: 999},
			setupMocks: func(cmdRepo *MockRoleCommandRepository, qryRepo *MockRoleQueryRepository) {
				qryRepo.On("FindByID", mock.Anything, uint(999)).Return(nil, nil)
			},
			wantErr: "role not found",
		},
		{
			name: "系统角色不可删除",
			cmd:  DeleteRoleCommand{RoleID: 1},
			setupMocks: func(cmdRepo *MockRoleCommandRepository, qryRepo *MockRoleQueryRepository) {
				qryRepo.On("FindByID", mock.Anything, uint(1)).Return(&role.Role{
					ID:       1,
					Name:     "admin",
					IsSystem: true,
				}, nil)
			},
			wantErr: "cannot delete system role",
		},
		{
			name: "查询角色时数据库错误",
			cmd:  DeleteRoleCommand{RoleID: 1},
			setupMocks: func(cmdRepo *MockRoleCommandRepository, qryRepo *MockRoleQueryRepository) {
				qryRepo.On("FindByID", mock.Anything, uint(1)).Return(nil, errors.New("database error"))
			},
			wantErr: "failed to find role",
		},
		{
			name: "删除角色时数据库错误",
			cmd:  DeleteRoleCommand{RoleID: 1},
			setupMocks: func(cmdRepo *MockRoleCommandRepository, qryRepo *MockRoleQueryRepository) {
				qryRepo.On("FindByID", mock.Anything, uint(1)).Return(&role.Role{
					ID:       1,
					Name:     "custom",
					IsSystem: false,
				}, nil)
				cmdRepo.On("Delete", mock.Anything, uint(1)).Return(errors.New("delete failed"))
			},
			wantErr: "failed to delete role",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockCmdRepo := new(MockRoleCommandRepository)
			mockQryRepo := new(MockRoleQueryRepository)
			tt.setupMocks(mockCmdRepo, mockQryRepo)

			handler := NewDeleteRoleHandler(mockCmdRepo, mockQryRepo)

			// Act
			err := handler.Handle(context.Background(), tt.cmd)

			// Assert
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
			mockCmdRepo.AssertExpectations(t)
			mockQryRepo.AssertExpectations(t)
		})
	}
}
