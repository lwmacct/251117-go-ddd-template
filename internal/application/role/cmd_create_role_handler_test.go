//nolint:forcetypeassert // 测试中的类型断言是可控的
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

func TestCreateRoleHandler_Handle_Success(t *testing.T) {
	tests := []struct {
		name string
		cmd  CreateRoleCommand
		want *CreateRoleResultDTO
	}{
		{
			name: "成功创建角色",
			cmd: CreateRoleCommand{
				Name:        "developer",
				DisplayName: "开发者",
				Description: "系统开发人员角色",
			},
			want: &CreateRoleResultDTO{
				RoleID:      1,
				Name:        "developer",
				DisplayName: "开发者",
			},
		},
		{
			name: "创建无描述的角色",
			cmd: CreateRoleCommand{
				Name:        "viewer",
				DisplayName: "查看者",
				Description: "",
			},
			want: &CreateRoleResultDTO{
				RoleID:      1,
				Name:        "viewer",
				DisplayName: "查看者",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockCmdRepo := new(MockRoleCommandRepository)
			mockQryRepo := new(MockRoleQueryRepository)

			mockQryRepo.On("ExistsByName", mock.Anything, tt.cmd.Name).Return(false, nil)
			mockCmdRepo.On("Create", mock.Anything, mock.AnythingOfType("*role.Role")).Return(nil)

			handler := NewCreateRoleHandler(mockCmdRepo, mockQryRepo)

			// Act
			result, err := handler.Handle(context.Background(), tt.cmd)

			// Assert
			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tt.want.Name, result.Name)
			assert.Equal(t, tt.want.DisplayName, result.DisplayName)
			mockCmdRepo.AssertExpectations(t)
			mockQryRepo.AssertExpectations(t)
		})
	}
}

func TestCreateRoleHandler_Handle_Error(t *testing.T) {
	tests := []struct {
		name       string
		cmd        CreateRoleCommand
		setupMocks func(*MockRoleCommandRepository, *MockRoleQueryRepository)
		wantErr    string
	}{
		{
			name: "角色名已存在",
			cmd: CreateRoleCommand{
				Name:        "admin",
				DisplayName: "管理员",
			},
			setupMocks: func(cmdRepo *MockRoleCommandRepository, qryRepo *MockRoleQueryRepository) {
				qryRepo.On("ExistsByName", mock.Anything, "admin").Return(true, nil)
			},
			wantErr: "role name already exists",
		},
		{
			name: "检查角色名时数据库错误",
			cmd: CreateRoleCommand{
				Name:        "test",
				DisplayName: "测试",
			},
			setupMocks: func(cmdRepo *MockRoleCommandRepository, qryRepo *MockRoleQueryRepository) {
				qryRepo.On("ExistsByName", mock.Anything, "test").Return(false, errors.New("database error"))
			},
			wantErr: "failed to check role name existence",
		},
		{
			name: "创建角色时数据库错误",
			cmd: CreateRoleCommand{
				Name:        "newrole",
				DisplayName: "新角色",
			},
			setupMocks: func(cmdRepo *MockRoleCommandRepository, qryRepo *MockRoleQueryRepository) {
				qryRepo.On("ExistsByName", mock.Anything, "newrole").Return(false, nil)
				cmdRepo.On("Create", mock.Anything, mock.AnythingOfType("*role.Role")).
					Return(errors.New("insert failed"))
			},
			wantErr: "failed to create role",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockCmdRepo := new(MockRoleCommandRepository)
			mockQryRepo := new(MockRoleQueryRepository)
			tt.setupMocks(mockCmdRepo, mockQryRepo)

			handler := NewCreateRoleHandler(mockCmdRepo, mockQryRepo)

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

func TestCreateRoleHandler_NotSystemRole(t *testing.T) {
	// 验证用户创建的角色不是系统角色
	mockCmdRepo := new(MockRoleCommandRepository)
	mockQryRepo := new(MockRoleQueryRepository)

	var capturedRole *role.Role
	mockQryRepo.On("ExistsByName", mock.Anything, "custom").Return(false, nil)
	mockCmdRepo.On("Create", mock.Anything, mock.AnythingOfType("*role.Role")).
		Run(func(args mock.Arguments) {
			capturedRole = args.Get(1).(*role.Role)
		}).Return(nil)

	handler := NewCreateRoleHandler(mockCmdRepo, mockQryRepo)
	_, err := handler.Handle(context.Background(), CreateRoleCommand{
		Name:        "custom",
		DisplayName: "自定义角色",
	})

	require.NoError(t, err)
	assert.NotNil(t, capturedRole)
	assert.False(t, capturedRole.IsSystem, "用户创建的角色不应是系统角色")
}
