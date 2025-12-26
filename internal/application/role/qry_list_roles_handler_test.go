//nolint:dupl // 与 qry_list_permissions_handler_test.go 测试结构相似，但测试不同 Handler
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

func TestListRolesHandler_Handle_Success(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name  string
		query ListRolesQuery
		roles []role.Role
		total int64
	}{
		{
			name:  "列出第一页角色",
			query: ListRolesQuery{Page: 1, Limit: 10},
			roles: []role.Role{
				{ID: 1, Name: "admin", DisplayName: "管理员", IsSystem: true, CreatedAt: now, UpdatedAt: now},
				{ID: 2, Name: "developer", DisplayName: "开发者", IsSystem: false, CreatedAt: now, UpdatedAt: now},
			},
			total: 2,
		},
		{
			name:  "空角色列表",
			query: ListRolesQuery{Page: 1, Limit: 10},
			roles: []role.Role{},
			total: 0,
		},
		{
			name:  "分页查询",
			query: ListRolesQuery{Page: 2, Limit: 5},
			roles: []role.Role{
				{ID: 6, Name: "role6", DisplayName: "角色6", CreatedAt: now, UpdatedAt: now},
			},
			total: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockQryRepo := new(MockRoleQueryRepository)
			mockQryRepo.On("List", mock.Anything, tt.query.Page, tt.query.Limit).
				Return(tt.roles, tt.total, nil)

			handler := NewListRolesHandler(mockQryRepo)

			// Act
			result, err := handler.Handle(context.Background(), tt.query)

			// Assert
			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tt.total, result.Total)
			assert.Equal(t, tt.query.Page, result.Page)
			assert.Equal(t, tt.query.Limit, result.Limit)
			assert.Len(t, result.Roles, len(tt.roles))
			mockQryRepo.AssertExpectations(t)
		})
	}
}

func TestListRolesHandler_Handle_Error(t *testing.T) {
	tests := []struct {
		name       string
		query      ListRolesQuery
		setupMocks func(*MockRoleQueryRepository)
		wantErr    string
	}{
		{
			name:  "数据库错误",
			query: ListRolesQuery{Page: 1, Limit: 10},
			setupMocks: func(qryRepo *MockRoleQueryRepository) {
				qryRepo.On("List", mock.Anything, 1, 10).Return(nil, int64(0), errors.New("database error"))
			},
			wantErr: "failed to list roles",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockQryRepo := new(MockRoleQueryRepository)
			tt.setupMocks(mockQryRepo)

			handler := NewListRolesHandler(mockQryRepo)

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

func TestListRolesHandler_Handle_RoleConversion(t *testing.T) {
	// 验证角色正确转换为 DTO
	now := time.Now()
	mockQryRepo := new(MockRoleQueryRepository)

	roles := []role.Role{
		{
			ID:          1,
			Name:        "admin",
			DisplayName: "管理员",
			Description: "系统管理员",
			IsSystem:    true,
			CreatedAt:   now,
			UpdatedAt:   now,
			Permissions: []role.Permission{
				{ID: 1, Code: "all:*:*"},
			},
		},
	}

	mockQryRepo.On("List", mock.Anything, 1, 10).Return(roles, int64(1), nil)

	handler := NewListRolesHandler(mockQryRepo)
	result, err := handler.Handle(context.Background(), ListRolesQuery{Page: 1, Limit: 10})

	require.NoError(t, err)
	assert.Len(t, result.Roles, 1)

	roleDTO := result.Roles[0]
	assert.Equal(t, uint(1), roleDTO.ID)
	assert.Equal(t, "admin", roleDTO.Name)
	assert.Equal(t, "管理员", roleDTO.DisplayName)
	assert.Equal(t, "系统管理员", roleDTO.Description)
	assert.True(t, roleDTO.IsSystem)
}
