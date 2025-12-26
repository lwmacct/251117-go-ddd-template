package user

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDeleteUserHandler_Handle_Success(t *testing.T) {
	tests := []struct {
		name string
		cmd  DeleteUserCommand
	}{
		{
			name: "成功删除用户",
			cmd:  DeleteUserCommand{UserID: 1},
		},
		{
			name: "删除另一用户",
			cmd:  DeleteUserCommand{UserID: 100},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockCmdRepo := new(MockUserCommandRepository)
			mockQryRepo := new(MockUserQueryRepository)
			mockEventBus := new(MockEventBus)

			mockQryRepo.On("Exists", mock.Anything, tt.cmd.UserID).Return(true, nil)
			mockCmdRepo.On("Delete", mock.Anything, tt.cmd.UserID).Return(nil)
			mockEventBus.On("Publish", mock.Anything, mock.Anything).Return(nil)

			handler := NewDeleteUserHandler(mockCmdRepo, mockQryRepo, mockEventBus)

			// Act
			err := handler.Handle(context.Background(), tt.cmd)

			// Assert
			require.NoError(t, err)
			mockCmdRepo.AssertExpectations(t)
			mockQryRepo.AssertExpectations(t)
		})
	}
}

func TestDeleteUserHandler_Handle_Error(t *testing.T) {
	tests := []struct {
		name       string
		cmd        DeleteUserCommand
		setupMocks func(*MockUserCommandRepository, *MockUserQueryRepository, *MockEventBus)
		wantErr    string
	}{
		{
			name: "用户不存在",
			cmd:  DeleteUserCommand{UserID: 999},
			setupMocks: func(cmdRepo *MockUserCommandRepository, qryRepo *MockUserQueryRepository, eventBus *MockEventBus) {
				qryRepo.On("Exists", mock.Anything, uint(999)).Return(false, nil)
			},
			wantErr: "not found",
		},
		{
			name: "检查用户存在时数据库错误",
			cmd:  DeleteUserCommand{UserID: 1},
			setupMocks: func(cmdRepo *MockUserCommandRepository, qryRepo *MockUserQueryRepository, eventBus *MockEventBus) {
				qryRepo.On("Exists", mock.Anything, uint(1)).Return(false, errors.New("database error"))
			},
			wantErr: "failed to check user existence",
		},
		{
			name: "删除用户时数据库错误",
			cmd:  DeleteUserCommand{UserID: 1},
			setupMocks: func(cmdRepo *MockUserCommandRepository, qryRepo *MockUserQueryRepository, eventBus *MockEventBus) {
				qryRepo.On("Exists", mock.Anything, uint(1)).Return(true, nil)
				cmdRepo.On("Delete", mock.Anything, uint(1)).Return(errors.New("delete failed"))
			},
			wantErr: "failed to delete user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockCmdRepo := new(MockUserCommandRepository)
			mockQryRepo := new(MockUserQueryRepository)
			mockEventBus := new(MockEventBus)
			tt.setupMocks(mockCmdRepo, mockQryRepo, mockEventBus)

			handler := NewDeleteUserHandler(mockCmdRepo, mockQryRepo, mockEventBus)

			// Act
			err := handler.Handle(context.Background(), tt.cmd)

			// Assert
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestDeleteUserHandler_Handle_NilEventBus(t *testing.T) {
	// 验证事件总线为 nil 时不会 panic
	mockCmdRepo := new(MockUserCommandRepository)
	mockQryRepo := new(MockUserQueryRepository)

	mockQryRepo.On("Exists", mock.Anything, uint(1)).Return(true, nil)
	mockCmdRepo.On("Delete", mock.Anything, uint(1)).Return(nil)

	handler := NewDeleteUserHandler(mockCmdRepo, mockQryRepo, nil)
	err := handler.Handle(context.Background(), DeleteUserCommand{UserID: 1})

	require.NoError(t, err)
}
