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

func TestGetRoleHandler_Handle_Success(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name  string
		query GetRoleQuery
		role  *role.Role
	}{
		{
			name:  "获取带权限的角色",
			query: GetRoleQuery{RoleID: 1},
			role: &role.Role{
				ID:          1,
				Name:        "developer",
				DisplayName: "开发者",
				Description: "开发人员角色",
				IsSystem:    false,
				CreatedAt:   now,
				UpdatedAt:   now,
				Permissions: []role.Permission{
					{ID: 1, Code: "user:read", Description: "读取用户"},
					{ID: 2, Code: "user:write", Description: "写入用户"},
				},
			},
		},
		{
			name:  "获取无权限的角色",
			query: GetRoleQuery{RoleID: 2},
			role: &role.Role{
				ID:          2,
				Name:        "viewer",
				DisplayName: "查看者",
				IsSystem:    false,
				CreatedAt:   now,
				UpdatedAt:   now,
				Permissions: []role.Permission{},
			},
		},
		{
			name:  "获取系统角色",
			query: GetRoleQuery{RoleID: 3},
			role: &role.Role{
				ID:          3,
				Name:        "admin",
				DisplayName: "管理员",
				IsSystem:    true,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockQryRepo := new(MockRoleQueryRepository)
			mockQryRepo.On("FindByIDWithPermissions", mock.Anything, tt.query.RoleID).Return(tt.role, nil)

			handler := NewGetRoleHandler(mockQryRepo)

			// Act
			result, err := handler.Handle(context.Background(), tt.query)

			// Assert
			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tt.role.ID, result.ID)
			assert.Equal(t, tt.role.Name, result.Name)
			assert.Equal(t, tt.role.DisplayName, result.DisplayName)
			assert.Equal(t, tt.role.IsSystem, result.IsSystem)
			assert.Len(t, result.Permissions, len(tt.role.Permissions))
			mockQryRepo.AssertExpectations(t)
		})
	}
}

func TestGetRoleHandler_Handle_Error(t *testing.T) {
	tests := []struct {
		name       string
		query      GetRoleQuery
		setupMocks func(*MockRoleQueryRepository)
		wantErr    string
	}{
		{
			name:  "角色不存在",
			query: GetRoleQuery{RoleID: 999},
			setupMocks: func(qryRepo *MockRoleQueryRepository) {
				qryRepo.On("FindByIDWithPermissions", mock.Anything, uint(999)).Return(nil, nil)
			},
			wantErr: "role not found",
		},
		{
			name:  "数据库错误",
			query: GetRoleQuery{RoleID: 1},
			setupMocks: func(qryRepo *MockRoleQueryRepository) {
				qryRepo.On("FindByIDWithPermissions", mock.Anything, uint(1)).
					Return(nil, errors.New("database error"))
			},
			wantErr: "failed to find role",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockQryRepo := new(MockRoleQueryRepository)
			tt.setupMocks(mockQryRepo)

			handler := NewGetRoleHandler(mockQryRepo)

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
