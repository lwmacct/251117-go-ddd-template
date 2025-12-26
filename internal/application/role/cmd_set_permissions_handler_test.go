package role

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSetPermissionsHandler_Handle_Success(t *testing.T) {
	tests := []struct {
		name string
		cmd  SetPermissionsCommand
	}{
		{
			name: "成功设置权限",
			cmd: SetPermissionsCommand{
				RoleID:        1,
				PermissionIDs: []uint{1, 2, 3},
			},
		},
		{
			name: "设置空权限列表",
			cmd: SetPermissionsCommand{
				RoleID:        2,
				PermissionIDs: []uint{},
			},
		},
		{
			name: "设置单个权限",
			cmd: SetPermissionsCommand{
				RoleID:        3,
				PermissionIDs: []uint{5},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockCmdRepo := new(MockRoleCommandRepository)
			mockQryRepo := new(MockRoleQueryRepository)
			mockPermRepo := new(MockPermissionQueryRepository)
			mockEventBus := new(MockEventBus)

			mockQryRepo.On("Exists", mock.Anything, tt.cmd.RoleID).Return(true, nil)
			for _, permID := range tt.cmd.PermissionIDs {
				mockPermRepo.On("Exists", mock.Anything, permID).Return(true, nil)
			}
			mockCmdRepo.On("SetPermissions", mock.Anything, tt.cmd.RoleID, tt.cmd.PermissionIDs).Return(nil)
			mockEventBus.On("Publish", mock.Anything, mock.Anything).Return(nil)

			handler := NewSetPermissionsHandler(mockCmdRepo, mockQryRepo, mockPermRepo, mockEventBus)

			// Act
			err := handler.Handle(context.Background(), tt.cmd)

			// Assert
			require.NoError(t, err)
			mockCmdRepo.AssertExpectations(t)
			mockQryRepo.AssertExpectations(t)
			mockPermRepo.AssertExpectations(t)
		})
	}
}

func TestSetPermissionsHandler_Handle_Error(t *testing.T) {
	tests := []struct {
		name       string
		cmd        SetPermissionsCommand
		setupMocks func(*MockRoleCommandRepository, *MockRoleQueryRepository, *MockPermissionQueryRepository, *MockEventBus)
		wantErr    string
	}{
		{
			name: "角色不存在",
			cmd:  SetPermissionsCommand{RoleID: 999, PermissionIDs: []uint{1}},
			setupMocks: func(cmdRepo *MockRoleCommandRepository, qryRepo *MockRoleQueryRepository, permRepo *MockPermissionQueryRepository, eventBus *MockEventBus) {
				qryRepo.On("Exists", mock.Anything, uint(999)).Return(false, nil)
			},
			wantErr: "role not found",
		},
		{
			name: "检查角色存在时数据库错误",
			cmd:  SetPermissionsCommand{RoleID: 1, PermissionIDs: []uint{1}},
			setupMocks: func(cmdRepo *MockRoleCommandRepository, qryRepo *MockRoleQueryRepository, permRepo *MockPermissionQueryRepository, eventBus *MockEventBus) {
				qryRepo.On("Exists", mock.Anything, uint(1)).Return(false, errors.New("database error"))
			},
			wantErr: "failed to check role existence",
		},
		{
			name: "权限不存在",
			cmd:  SetPermissionsCommand{RoleID: 1, PermissionIDs: []uint{1, 999}},
			setupMocks: func(cmdRepo *MockRoleCommandRepository, qryRepo *MockRoleQueryRepository, permRepo *MockPermissionQueryRepository, eventBus *MockEventBus) {
				qryRepo.On("Exists", mock.Anything, uint(1)).Return(true, nil)
				permRepo.On("Exists", mock.Anything, uint(1)).Return(true, nil)
				permRepo.On("Exists", mock.Anything, uint(999)).Return(false, nil)
			},
			wantErr: "permission not found",
		},
		{
			name: "检查权限存在时数据库错误",
			cmd:  SetPermissionsCommand{RoleID: 1, PermissionIDs: []uint{1}},
			setupMocks: func(cmdRepo *MockRoleCommandRepository, qryRepo *MockRoleQueryRepository, permRepo *MockPermissionQueryRepository, eventBus *MockEventBus) {
				qryRepo.On("Exists", mock.Anything, uint(1)).Return(true, nil)
				permRepo.On("Exists", mock.Anything, uint(1)).Return(false, errors.New("database error"))
			},
			wantErr: "failed to check permission existence",
		},
		{
			name: "设置权限时数据库错误",
			cmd:  SetPermissionsCommand{RoleID: 1, PermissionIDs: []uint{1}},
			setupMocks: func(cmdRepo *MockRoleCommandRepository, qryRepo *MockRoleQueryRepository, permRepo *MockPermissionQueryRepository, eventBus *MockEventBus) {
				qryRepo.On("Exists", mock.Anything, uint(1)).Return(true, nil)
				permRepo.On("Exists", mock.Anything, uint(1)).Return(true, nil)
				cmdRepo.On("SetPermissions", mock.Anything, uint(1), []uint{1}).Return(errors.New("set failed"))
			},
			wantErr: "failed to set permissions",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockCmdRepo := new(MockRoleCommandRepository)
			mockQryRepo := new(MockRoleQueryRepository)
			mockPermRepo := new(MockPermissionQueryRepository)
			mockEventBus := new(MockEventBus)
			tt.setupMocks(mockCmdRepo, mockQryRepo, mockPermRepo, mockEventBus)

			handler := NewSetPermissionsHandler(mockCmdRepo, mockQryRepo, mockPermRepo, mockEventBus)

			// Act
			err := handler.Handle(context.Background(), tt.cmd)

			// Assert
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
			mockCmdRepo.AssertExpectations(t)
			mockQryRepo.AssertExpectations(t)
			mockPermRepo.AssertExpectations(t)
		})
	}
}

func TestSetPermissionsHandler_Handle_NilEventBus(t *testing.T) {
	// 验证事件总线为 nil 时不会 panic
	mockCmdRepo := new(MockRoleCommandRepository)
	mockQryRepo := new(MockRoleQueryRepository)
	mockPermRepo := new(MockPermissionQueryRepository)

	mockQryRepo.On("Exists", mock.Anything, uint(1)).Return(true, nil)
	mockPermRepo.On("Exists", mock.Anything, uint(1)).Return(true, nil)
	mockCmdRepo.On("SetPermissions", mock.Anything, uint(1), []uint{1}).Return(nil)

	handler := NewSetPermissionsHandler(mockCmdRepo, mockQryRepo, mockPermRepo, nil)

	err := handler.Handle(context.Background(), SetPermissionsCommand{
		RoleID:        1,
		PermissionIDs: []uint{1},
	})

	require.NoError(t, err)
}
