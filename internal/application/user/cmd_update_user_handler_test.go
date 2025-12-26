package user

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
)

func TestUpdateUserHandler_Handle_Success(t *testing.T) {
	now := time.Now()
	username := "newusername"
	email := "new@example.com"
	fullName := "New Name"
	status := "inactive"

	tests := []struct {
		name         string
		cmd          UpdateUserCommand
		existingUser *user.User
	}{
		{
			name: "更新用户名",
			cmd:  UpdateUserCommand{UserID: 1, Username: &username},
			existingUser: &user.User{
				ID:        1,
				Username:  "olduser",
				Email:     "old@example.com",
				Status:    "active",
				CreatedAt: now,
			},
		},
		{
			name: "更新多个字段",
			cmd:  UpdateUserCommand{UserID: 1, Email: &email, FullName: &fullName},
			existingUser: &user.User{
				ID:        1,
				Username:  "user1",
				Email:     "old@example.com",
				Status:    "active",
				CreatedAt: now,
			},
		},
		{
			name: "更新状态为 inactive",
			cmd:  UpdateUserCommand{UserID: 1, Status: &status},
			existingUser: &user.User{
				ID:        1,
				Username:  "user1",
				Email:     "user@example.com",
				Status:    "active",
				CreatedAt: now,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockCmdRepo := new(MockUserCommandRepository)
			mockQryRepo := new(MockUserQueryRepository)

			mockQryRepo.On("GetByID", mock.Anything, tt.cmd.UserID).Return(tt.existingUser, nil)
			if tt.cmd.Username != nil && *tt.cmd.Username != tt.existingUser.Username {
				mockQryRepo.On("ExistsByUsername", mock.Anything, *tt.cmd.Username).Return(false, nil)
			}
			if tt.cmd.Email != nil && *tt.cmd.Email != tt.existingUser.Email {
				mockQryRepo.On("ExistsByEmail", mock.Anything, *tt.cmd.Email).Return(false, nil)
			}
			mockCmdRepo.On("Update", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)

			handler := NewUpdateUserHandler(mockCmdRepo, mockQryRepo)

			// Act
			result, err := handler.Handle(context.Background(), tt.cmd)

			// Assert
			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tt.cmd.UserID, result.UserID)
			mockCmdRepo.AssertExpectations(t)
			mockQryRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateUserHandler_Handle_Error(t *testing.T) {
	username := "existinguser"
	email := "existing@example.com"
	invalidStatus := "unknown"

	tests := []struct {
		name       string
		cmd        UpdateUserCommand
		setupMocks func(*MockUserCommandRepository, *MockUserQueryRepository)
		wantErr    string
	}{
		{
			name: "用户不存在",
			cmd:  UpdateUserCommand{UserID: 999},
			setupMocks: func(cmdRepo *MockUserCommandRepository, qryRepo *MockUserQueryRepository) {
				qryRepo.On("GetByID", mock.Anything, uint(999)).Return(nil, errors.New("not found"))
			},
			wantErr: "not found",
		},
		{
			name: "用户名已被占用",
			cmd:  UpdateUserCommand{UserID: 1, Username: &username},
			setupMocks: func(cmdRepo *MockUserCommandRepository, qryRepo *MockUserQueryRepository) {
				qryRepo.On("GetByID", mock.Anything, uint(1)).Return(&user.User{
					ID:       1,
					Username: "olduser",
				}, nil)
				qryRepo.On("ExistsByUsername", mock.Anything, "existinguser").Return(true, nil)
			},
			wantErr: "username already exists",
		},
		{
			name: "检查用户名时数据库错误",
			cmd:  UpdateUserCommand{UserID: 1, Username: &username},
			setupMocks: func(cmdRepo *MockUserCommandRepository, qryRepo *MockUserQueryRepository) {
				qryRepo.On("GetByID", mock.Anything, uint(1)).Return(&user.User{
					ID:       1,
					Username: "olduser",
				}, nil)
				qryRepo.On("ExistsByUsername", mock.Anything, "existinguser").Return(false, errors.New("db error"))
			},
			wantErr: "failed to check username existence",
		},
		{
			name: "邮箱已被占用",
			cmd:  UpdateUserCommand{UserID: 1, Email: &email},
			setupMocks: func(cmdRepo *MockUserCommandRepository, qryRepo *MockUserQueryRepository) {
				qryRepo.On("GetByID", mock.Anything, uint(1)).Return(&user.User{
					ID:    1,
					Email: "old@example.com",
				}, nil)
				qryRepo.On("ExistsByEmail", mock.Anything, "existing@example.com").Return(true, nil)
			},
			wantErr: "email already exists",
		},
		{
			name: "检查邮箱时数据库错误",
			cmd:  UpdateUserCommand{UserID: 1, Email: &email},
			setupMocks: func(cmdRepo *MockUserCommandRepository, qryRepo *MockUserQueryRepository) {
				qryRepo.On("GetByID", mock.Anything, uint(1)).Return(&user.User{
					ID:    1,
					Email: "old@example.com",
				}, nil)
				qryRepo.On("ExistsByEmail", mock.Anything, "existing@example.com").Return(false, errors.New("db error"))
			},
			wantErr: "failed to check email existence",
		},
		{
			name: "无效的状态值",
			cmd:  UpdateUserCommand{UserID: 1, Status: &invalidStatus},
			setupMocks: func(cmdRepo *MockUserCommandRepository, qryRepo *MockUserQueryRepository) {
				qryRepo.On("GetByID", mock.Anything, uint(1)).Return(&user.User{
					ID:     1,
					Status: "active",
				}, nil)
			},
			wantErr: "invalid user status",
		},
		{
			name: "保存更新失败",
			cmd:  UpdateUserCommand{UserID: 1},
			setupMocks: func(cmdRepo *MockUserCommandRepository, qryRepo *MockUserQueryRepository) {
				qryRepo.On("GetByID", mock.Anything, uint(1)).Return(&user.User{
					ID:     1,
					Status: "active",
				}, nil)
				cmdRepo.On("Update", mock.Anything, mock.AnythingOfType("*user.User")).Return(errors.New("db error"))
			},
			wantErr: "failed to update user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockCmdRepo := new(MockUserCommandRepository)
			mockQryRepo := new(MockUserQueryRepository)
			tt.setupMocks(mockCmdRepo, mockQryRepo)

			handler := NewUpdateUserHandler(mockCmdRepo, mockQryRepo)

			// Act
			result, err := handler.Handle(context.Background(), tt.cmd)

			// Assert
			assert.Nil(t, result)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

// ============================================================
// 补充覆盖率测试
// ============================================================

func TestUpdateUserHandler_Handle_UpdateAllFields(t *testing.T) {
	now := time.Now()
	avatar := "https://example.com/avatar.png"
	bio := "This is my bio"
	statusActive := "active"
	statusBanned := "banned"

	tests := []struct {
		name         string
		cmd          UpdateUserCommand
		existingUser *user.User
		wantStatus   string
	}{
		{
			name: "更新 Avatar 和 Bio",
			cmd:  UpdateUserCommand{UserID: 1, Avatar: &avatar, Bio: &bio},
			existingUser: &user.User{
				ID:        1,
				Username:  "user1",
				Email:     "user@example.com",
				Status:    "active",
				CreatedAt: now,
			},
			wantStatus: "active",
		},
		{
			name: "更新状态为 banned",
			cmd:  UpdateUserCommand{UserID: 1, Status: &statusBanned},
			existingUser: &user.User{
				ID:        1,
				Username:  "user1",
				Email:     "user@example.com",
				Status:    "active",
				CreatedAt: now,
			},
			wantStatus: "banned",
		},
		{
			name: "更新状态为 active",
			cmd:  UpdateUserCommand{UserID: 1, Status: &statusActive},
			existingUser: &user.User{
				ID:        1,
				Username:  "user1",
				Email:     "user@example.com",
				Status:    "banned",
				CreatedAt: now,
			},
			wantStatus: "active",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCmdRepo := new(MockUserCommandRepository)
			mockQryRepo := new(MockUserQueryRepository)

			mockQryRepo.On("GetByID", mock.Anything, tt.cmd.UserID).Return(tt.existingUser, nil)
			mockCmdRepo.On("Update", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)

			handler := NewUpdateUserHandler(mockCmdRepo, mockQryRepo)

			result, err := handler.Handle(context.Background(), tt.cmd)

			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tt.wantStatus, tt.existingUser.Status)
			if tt.cmd.Avatar != nil {
				assert.Equal(t, *tt.cmd.Avatar, tt.existingUser.Avatar)
			}
			if tt.cmd.Bio != nil {
				assert.Equal(t, *tt.cmd.Bio, tt.existingUser.Bio)
			}
		})
	}
}
