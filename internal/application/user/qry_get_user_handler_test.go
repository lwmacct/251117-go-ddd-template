package user

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	domainRole "github.com/lwmacct/251117-go-ddd-template/internal/domain/role"
	domainUser "github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
)

func TestGetUserHandler_Handle_Success(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name  string
		query GetUserQuery
		user  *domainUser.User
	}{
		{
			name:  "获取用户基本信息",
			query: GetUserQuery{UserID: 1, WithRoles: false},
			user: &domainUser.User{
				ID:        1,
				Username:  "testuser",
				Email:     "test@example.com",
				FullName:  "Test User",
				Status:    "active",
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		{
			name:  "获取用户包含角色信息",
			query: GetUserQuery{UserID: 1, WithRoles: true},
			user: &domainUser.User{
				ID:        1,
				Username:  "admin",
				Email:     "admin@example.com",
				Status:    "active",
				CreatedAt: now,
				UpdatedAt: now,
				Roles: []domainRole.Role{
					{ID: 1, Name: "admin", DisplayName: "管理员"},
					{ID: 2, Name: "developer", DisplayName: "开发者"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockQryRepo := new(MockUserQueryRepository)

			if tt.query.WithRoles {
				mockQryRepo.On("GetByIDWithRoles", mock.Anything, tt.query.UserID).Return(tt.user, nil)
			} else {
				mockQryRepo.On("GetByID", mock.Anything, tt.query.UserID).Return(tt.user, nil)
			}

			handler := NewGetUserHandler(mockQryRepo)

			// Act
			result, err := handler.Handle(context.Background(), tt.query)

			// Assert
			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tt.user.ID, result.ID)
			assert.Equal(t, tt.user.Username, result.Username)
			assert.Equal(t, tt.user.Email, result.Email)

			if tt.query.WithRoles {
				assert.Len(t, result.Roles, len(tt.user.Roles))
			}

			mockQryRepo.AssertExpectations(t)
		})
	}
}

func TestGetUserHandler_Handle_Error(t *testing.T) {
	tests := []struct {
		name       string
		query      GetUserQuery
		setupMocks func(*MockUserQueryRepository)
		wantErr    string
	}{
		{
			name:  "用户不存在（不含角色）",
			query: GetUserQuery{UserID: 999, WithRoles: false},
			setupMocks: func(qryRepo *MockUserQueryRepository) {
				qryRepo.On("GetByID", mock.Anything, uint(999)).Return(nil, errors.New("user not found"))
			},
			wantErr: "user not found",
		},
		{
			name:  "用户不存在（含角色）",
			query: GetUserQuery{UserID: 999, WithRoles: true},
			setupMocks: func(qryRepo *MockUserQueryRepository) {
				qryRepo.On("GetByIDWithRoles", mock.Anything, uint(999)).Return(nil, errors.New("user not found"))
			},
			wantErr: "user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockQryRepo := new(MockUserQueryRepository)
			tt.setupMocks(mockQryRepo)

			handler := NewGetUserHandler(mockQryRepo)

			// Act
			result, err := handler.Handle(context.Background(), tt.query)

			// Assert
			assert.Nil(t, result)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
			mockQryRepo.AssertExpectations(t)
		})
	}
}

func TestGetUserHandler_Handle_RoleConversion(t *testing.T) {
	// 验证角色正确转换为 DTO
	now := time.Now()
	mockQryRepo := new(MockUserQueryRepository)

	user := &domainUser.User{
		ID:        1,
		Username:  "admin",
		Email:     "admin@example.com",
		Status:    "active",
		CreatedAt: now,
		UpdatedAt: now,
		Roles: []domainRole.Role{
			{
				ID:          1,
				Name:        "admin",
				DisplayName: "管理员",
				Description: "系统管理员",
			},
		},
	}

	mockQryRepo.On("GetByIDWithRoles", mock.Anything, uint(1)).Return(user, nil)

	handler := NewGetUserHandler(mockQryRepo)
	result, err := handler.Handle(context.Background(), GetUserQuery{UserID: 1, WithRoles: true})

	require.NoError(t, err)
	assert.Len(t, result.Roles, 1)

	roleDTO := result.Roles[0]
	assert.Equal(t, uint(1), roleDTO.ID)
	assert.Equal(t, "admin", roleDTO.Name)
	assert.Equal(t, "管理员", roleDTO.DisplayName)
	assert.Equal(t, "系统管理员", roleDTO.Description)
}
