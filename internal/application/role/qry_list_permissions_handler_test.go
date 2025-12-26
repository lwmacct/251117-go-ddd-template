//nolint:dupl // 与 qry_list_roles_handler_test.go 测试结构相似，但测试不同 Handler
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

func TestListPermissionsHandler_Handle_Success(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name        string
		query       ListPermissionsQuery
		permissions []role.Permission
		total       int64
	}{
		{
			name:  "列出权限",
			query: ListPermissionsQuery{Page: 1, Limit: 10},
			permissions: []role.Permission{
				{ID: 1, Code: "user:read", Description: "读取用户", Resource: "user", Action: "read", CreatedAt: now, UpdatedAt: now},
				{ID: 2, Code: "user:write", Description: "写入用户", Resource: "user", Action: "write", CreatedAt: now, UpdatedAt: now},
				{ID: 3, Code: "role:manage", Description: "管理角色", Resource: "role", Action: "manage", CreatedAt: now, UpdatedAt: now},
			},
			total: 3,
		},
		{
			name:        "空权限列表",
			query:       ListPermissionsQuery{Page: 1, Limit: 10},
			permissions: []role.Permission{},
			total:       0,
		},
		{
			name:  "分页查询",
			query: ListPermissionsQuery{Page: 2, Limit: 5},
			permissions: []role.Permission{
				{ID: 6, Code: "menu:read", Description: "读取菜单", CreatedAt: now, UpdatedAt: now},
			},
			total: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockPermRepo := new(MockPermissionQueryRepository)
			mockPermRepo.On("List", mock.Anything, tt.query.Page, tt.query.Limit).
				Return(tt.permissions, tt.total, nil)

			handler := NewListPermissionsHandler(mockPermRepo)

			// Act
			result, err := handler.Handle(context.Background(), tt.query)

			// Assert
			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tt.total, result.Total)
			assert.Equal(t, tt.query.Page, result.Page)
			assert.Equal(t, tt.query.Limit, result.Limit)
			assert.Len(t, result.Permissions, len(tt.permissions))
			mockPermRepo.AssertExpectations(t)
		})
	}
}

func TestListPermissionsHandler_Handle_Error(t *testing.T) {
	tests := []struct {
		name       string
		query      ListPermissionsQuery
		setupMocks func(*MockPermissionQueryRepository)
		wantErr    string
	}{
		{
			name:  "数据库错误",
			query: ListPermissionsQuery{Page: 1, Limit: 10},
			setupMocks: func(permRepo *MockPermissionQueryRepository) {
				permRepo.On("List", mock.Anything, 1, 10).
					Return(nil, int64(0), errors.New("database error"))
			},
			wantErr: "failed to list permissions",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockPermRepo := new(MockPermissionQueryRepository)
			tt.setupMocks(mockPermRepo)

			handler := NewListPermissionsHandler(mockPermRepo)

			// Act
			result, err := handler.Handle(context.Background(), tt.query)

			// Assert
			assert.Nil(t, result)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
			mockPermRepo.AssertExpectations(t)
		})
	}
}

func TestListPermissionsHandler_Handle_PermissionConversion(t *testing.T) {
	// 验证权限正确转换为 DTO
	now := time.Now()
	mockPermRepo := new(MockPermissionQueryRepository)

	permissions := []role.Permission{
		{
			ID:          1,
			Code:        "admin:user:create",
			Description: "创建新用户权限",
			Resource:    "user",
			Action:      "create",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	mockPermRepo.On("List", mock.Anything, 1, 10).Return(permissions, int64(1), nil)

	handler := NewListPermissionsHandler(mockPermRepo)
	result, err := handler.Handle(context.Background(), ListPermissionsQuery{Page: 1, Limit: 10})

	require.NoError(t, err)
	assert.Len(t, result.Permissions, 1)

	permDTO := result.Permissions[0]
	assert.Equal(t, uint(1), permDTO.ID)
	assert.Equal(t, "admin:user:create", permDTO.Code)
	assert.Equal(t, "创建新用户权限", permDTO.Description)
	assert.Equal(t, "user", permDTO.Resource)
	assert.Equal(t, "create", permDTO.Action)
}
