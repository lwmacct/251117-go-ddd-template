package user

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
)

func TestAssignRolesHandler_Handle_Success(t *testing.T) {
	tests := []struct {
		name string
		cmd  AssignRolesCommand
	}{
		{
			name: "成功分配单个角色",
			cmd:  AssignRolesCommand{UserID: 1, RoleIDs: []uint{1}},
		},
		{
			name: "成功分配多个角色",
			cmd:  AssignRolesCommand{UserID: 1, RoleIDs: []uint{1, 2, 3}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockCmdRepo := new(MockUserCommandRepository)
			mockQryRepo := new(MockUserQueryRepository)
			mockEventBus := new(MockEventBus)

			mockQryRepo.On("GetByID", mock.Anything, tt.cmd.UserID).Return(&user.User{ID: tt.cmd.UserID}, nil)
			mockCmdRepo.On("AssignRoles", mock.Anything, tt.cmd.UserID, tt.cmd.RoleIDs).Return(nil)
			mockEventBus.On("Publish", mock.Anything, mock.Anything).Return(nil)

			handler := NewAssignRolesHandler(mockCmdRepo, mockQryRepo, mockEventBus)

			// Act
			err := handler.Handle(context.Background(), tt.cmd)

			// Assert
			require.NoError(t, err)
			mockCmdRepo.AssertExpectations(t)
			mockQryRepo.AssertExpectations(t)
		})
	}
}

func TestAssignRolesHandler_Handle_Error(t *testing.T) {
	tests := []struct {
		name       string
		cmd        AssignRolesCommand
		setupMocks func(*MockUserCommandRepository, *MockUserQueryRepository, *MockEventBus)
		wantErr    string
	}{
		{
			name: "用户不存在",
			cmd:  AssignRolesCommand{UserID: 999, RoleIDs: []uint{1}},
			setupMocks: func(cmdRepo *MockUserCommandRepository, qryRepo *MockUserQueryRepository, eventBus *MockEventBus) {
				qryRepo.On("GetByID", mock.Anything, uint(999)).Return(nil, errors.New("user not found"))
			},
			wantErr: "user not found",
		},
		{
			name: "分配角色失败",
			cmd:  AssignRolesCommand{UserID: 1, RoleIDs: []uint{1}},
			setupMocks: func(cmdRepo *MockUserCommandRepository, qryRepo *MockUserQueryRepository, eventBus *MockEventBus) {
				qryRepo.On("GetByID", mock.Anything, uint(1)).Return(&user.User{ID: 1}, nil)
				cmdRepo.On("AssignRoles", mock.Anything, uint(1), []uint{1}).Return(errors.New("role not found"))
			},
			wantErr: "role not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockCmdRepo := new(MockUserCommandRepository)
			mockQryRepo := new(MockUserQueryRepository)
			mockEventBus := new(MockEventBus)
			tt.setupMocks(mockCmdRepo, mockQryRepo, mockEventBus)

			handler := NewAssignRolesHandler(mockCmdRepo, mockQryRepo, mockEventBus)

			// Act
			err := handler.Handle(context.Background(), tt.cmd)

			// Assert
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestAssignRolesHandler_Handle_NilEventBus(t *testing.T) {
	mockCmdRepo := new(MockUserCommandRepository)
	mockQryRepo := new(MockUserQueryRepository)

	mockQryRepo.On("GetByID", mock.Anything, uint(1)).Return(&user.User{ID: 1}, nil)
	mockCmdRepo.On("AssignRoles", mock.Anything, uint(1), []uint{1}).Return(nil)

	handler := NewAssignRolesHandler(mockCmdRepo, mockQryRepo, nil)
	err := handler.Handle(context.Background(), AssignRolesCommand{UserID: 1, RoleIDs: []uint{1}})

	require.NoError(t, err)
}
