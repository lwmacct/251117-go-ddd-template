//nolint:forcetypeassert // 测试中的类型断言是可控的
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

func TestCreateUserHandler_Handle_Success(t *testing.T) {
	tests := []struct {
		name string
		cmd  CreateUserCommand
		want *CreateUserResultDTO
	}{
		{
			name: "成功创建用户",
			cmd: CreateUserCommand{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
				FullName: "Test User",
			},
			want: &CreateUserResultDTO{
				UserID:   1,
				Username: "testuser",
				Email:    "test@example.com",
			},
		},
		{
			name: "创建带角色的用户",
			cmd: CreateUserCommand{
				Username: "admin",
				Email:    "admin@example.com",
				Password: "admin123",
				FullName: "Admin User",
				RoleIDs:  []uint{1, 2},
			},
			want: &CreateUserResultDTO{
				UserID:   1,
				Username: "admin",
				Email:    "admin@example.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockCmdRepo := new(MockUserCommandRepository)
			mockQryRepo := new(MockUserQueryRepository)
			mockAuthService := new(MockAuthService)

			mockAuthService.On("ValidatePasswordPolicy", mock.Anything, tt.cmd.Password).Return(nil)
			mockQryRepo.On("ExistsByUsername", mock.Anything, tt.cmd.Username).Return(false, nil)
			mockQryRepo.On("ExistsByEmail", mock.Anything, tt.cmd.Email).Return(false, nil)
			mockAuthService.On("GeneratePasswordHash", mock.Anything, tt.cmd.Password).Return("hashed_password", nil)
			mockCmdRepo.On("Create", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)

			if len(tt.cmd.RoleIDs) > 0 {
				mockCmdRepo.On("AssignRoles", mock.Anything, uint(1), tt.cmd.RoleIDs).Return(nil)
			}

			handler := NewCreateUserHandler(mockCmdRepo, mockQryRepo, mockAuthService)

			// Act
			result, err := handler.Handle(context.Background(), tt.cmd)

			// Assert
			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tt.want.Username, result.Username)
			assert.Equal(t, tt.want.Email, result.Email)
			mockCmdRepo.AssertExpectations(t)
			mockQryRepo.AssertExpectations(t)
			mockAuthService.AssertExpectations(t)
		})
	}
}

func TestCreateUserHandler_Handle_Error(t *testing.T) {
	tests := []struct {
		name       string
		cmd        CreateUserCommand
		setupMocks func(*MockUserCommandRepository, *MockUserQueryRepository, *MockAuthService)
		wantErr    string
	}{
		{
			name: "密码策略验证失败",
			cmd:  CreateUserCommand{Username: "test", Email: "test@example.com", Password: "123"},
			setupMocks: func(cmdRepo *MockUserCommandRepository, qryRepo *MockUserQueryRepository, authService *MockAuthService) {
				authService.On("ValidatePasswordPolicy", mock.Anything, "123").Return(errors.New("password too short"))
			},
			wantErr: "password too short",
		},
		{
			name: "用户名已存在",
			cmd:  CreateUserCommand{Username: "existing", Email: "test@example.com", Password: "password123"},
			setupMocks: func(cmdRepo *MockUserCommandRepository, qryRepo *MockUserQueryRepository, authService *MockAuthService) {
				authService.On("ValidatePasswordPolicy", mock.Anything, "password123").Return(nil)
				qryRepo.On("ExistsByUsername", mock.Anything, "existing").Return(true, nil)
			},
			wantErr: "username already exists",
		},
		{
			name: "邮箱已存在",
			cmd:  CreateUserCommand{Username: "newuser", Email: "existing@example.com", Password: "password123"},
			setupMocks: func(cmdRepo *MockUserCommandRepository, qryRepo *MockUserQueryRepository, authService *MockAuthService) {
				authService.On("ValidatePasswordPolicy", mock.Anything, "password123").Return(nil)
				qryRepo.On("ExistsByUsername", mock.Anything, "newuser").Return(false, nil)
				qryRepo.On("ExistsByEmail", mock.Anything, "existing@example.com").Return(true, nil)
			},
			wantErr: "email already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockCmdRepo := new(MockUserCommandRepository)
			mockQryRepo := new(MockUserQueryRepository)
			mockAuthService := new(MockAuthService)
			tt.setupMocks(mockCmdRepo, mockQryRepo, mockAuthService)

			handler := NewCreateUserHandler(mockCmdRepo, mockQryRepo, mockAuthService)

			// Act
			result, err := handler.Handle(context.Background(), tt.cmd)

			// Assert
			assert.Nil(t, result)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestCreateUserHandler_Handle_DefaultStatus(t *testing.T) {
	// 验证默认状态为 active
	mockCmdRepo := new(MockUserCommandRepository)
	mockQryRepo := new(MockUserQueryRepository)
	mockAuthService := new(MockAuthService)

	var capturedUser *user.User
	mockAuthService.On("ValidatePasswordPolicy", mock.Anything, mock.Anything).Return(nil)
	mockQryRepo.On("ExistsByUsername", mock.Anything, mock.Anything).Return(false, nil)
	mockQryRepo.On("ExistsByEmail", mock.Anything, mock.Anything).Return(false, nil)
	mockAuthService.On("GeneratePasswordHash", mock.Anything, mock.Anything).Return("hashed", nil)
	mockCmdRepo.On("Create", mock.Anything, mock.AnythingOfType("*user.User")).
		Run(func(args mock.Arguments) {
			capturedUser = args.Get(1).(*user.User)
		}).Return(nil)

	handler := NewCreateUserHandler(mockCmdRepo, mockQryRepo, mockAuthService)
	_, err := handler.Handle(context.Background(), CreateUserCommand{
		Username: "test",
		Email:    "test@example.com",
		Password: "password123",
	})

	require.NoError(t, err)
	assert.NotNil(t, capturedUser)
	assert.Equal(t, "active", capturedUser.Status)
}

// ============================================================
// 补充覆盖率测试
// ============================================================

func TestCreateUserHandler_Handle_CheckUsernameError(t *testing.T) {
	mockCmdRepo := new(MockUserCommandRepository)
	mockQryRepo := new(MockUserQueryRepository)
	mockAuthService := new(MockAuthService)

	mockAuthService.On("ValidatePasswordPolicy", mock.Anything, "password123").Return(nil)
	mockQryRepo.On("ExistsByUsername", mock.Anything, "test").Return(false, errors.New("db error"))

	handler := NewCreateUserHandler(mockCmdRepo, mockQryRepo, mockAuthService)
	result, err := handler.Handle(context.Background(), CreateUserCommand{
		Username: "test",
		Email:    "test@example.com",
		Password: "password123",
	})

	assert.Nil(t, result)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to check username existence")
}

func TestCreateUserHandler_Handle_CheckEmailError(t *testing.T) {
	mockCmdRepo := new(MockUserCommandRepository)
	mockQryRepo := new(MockUserQueryRepository)
	mockAuthService := new(MockAuthService)

	mockAuthService.On("ValidatePasswordPolicy", mock.Anything, "password123").Return(nil)
	mockQryRepo.On("ExistsByUsername", mock.Anything, "test").Return(false, nil)
	mockQryRepo.On("ExistsByEmail", mock.Anything, "test@example.com").Return(false, errors.New("db error"))

	handler := NewCreateUserHandler(mockCmdRepo, mockQryRepo, mockAuthService)
	result, err := handler.Handle(context.Background(), CreateUserCommand{
		Username: "test",
		Email:    "test@example.com",
		Password: "password123",
	})

	assert.Nil(t, result)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to check email existence")
}

func TestCreateUserHandler_Handle_HashPasswordError(t *testing.T) {
	mockCmdRepo := new(MockUserCommandRepository)
	mockQryRepo := new(MockUserQueryRepository)
	mockAuthService := new(MockAuthService)

	mockAuthService.On("ValidatePasswordPolicy", mock.Anything, "password123").Return(nil)
	mockQryRepo.On("ExistsByUsername", mock.Anything, "test").Return(false, nil)
	mockQryRepo.On("ExistsByEmail", mock.Anything, "test@example.com").Return(false, nil)
	mockAuthService.On("GeneratePasswordHash", mock.Anything, "password123").Return("", errors.New("hash error"))

	handler := NewCreateUserHandler(mockCmdRepo, mockQryRepo, mockAuthService)
	result, err := handler.Handle(context.Background(), CreateUserCommand{
		Username: "test",
		Email:    "test@example.com",
		Password: "password123",
	})

	assert.Nil(t, result)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to hash password")
}

func TestCreateUserHandler_Handle_CreateUserError(t *testing.T) {
	mockCmdRepo := new(MockUserCommandRepository)
	mockQryRepo := new(MockUserQueryRepository)
	mockAuthService := new(MockAuthService)

	mockAuthService.On("ValidatePasswordPolicy", mock.Anything, "password123").Return(nil)
	mockQryRepo.On("ExistsByUsername", mock.Anything, "test").Return(false, nil)
	mockQryRepo.On("ExistsByEmail", mock.Anything, "test@example.com").Return(false, nil)
	mockAuthService.On("GeneratePasswordHash", mock.Anything, "password123").Return("hashed", nil)
	mockCmdRepo.On("Create", mock.Anything, mock.AnythingOfType("*user.User")).Return(errors.New("db error"))

	handler := NewCreateUserHandler(mockCmdRepo, mockQryRepo, mockAuthService)
	result, err := handler.Handle(context.Background(), CreateUserCommand{
		Username: "test",
		Email:    "test@example.com",
		Password: "password123",
	})

	assert.Nil(t, result)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create user")
}

func TestCreateUserHandler_Handle_AssignRolesError(t *testing.T) {
	mockCmdRepo := new(MockUserCommandRepository)
	mockQryRepo := new(MockUserQueryRepository)
	mockAuthService := new(MockAuthService)

	mockAuthService.On("ValidatePasswordPolicy", mock.Anything, "password123").Return(nil)
	mockQryRepo.On("ExistsByUsername", mock.Anything, "test").Return(false, nil)
	mockQryRepo.On("ExistsByEmail", mock.Anything, "test@example.com").Return(false, nil)
	mockAuthService.On("GeneratePasswordHash", mock.Anything, "password123").Return("hashed", nil)
	mockCmdRepo.On("Create", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)
	mockCmdRepo.On("AssignRoles", mock.Anything, uint(1), []uint{1, 2}).Return(errors.New("role error"))

	handler := NewCreateUserHandler(mockCmdRepo, mockQryRepo, mockAuthService)
	result, err := handler.Handle(context.Background(), CreateUserCommand{
		Username: "test",
		Email:    "test@example.com",
		Password: "password123",
		RoleIDs:  []uint{1, 2},
	})

	assert.Nil(t, result)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to assign roles")
}

func TestCreateUserHandler_Handle_WithStatus(t *testing.T) {
	mockCmdRepo := new(MockUserCommandRepository)
	mockQryRepo := new(MockUserQueryRepository)
	mockAuthService := new(MockAuthService)

	status := "inactive"
	var capturedUser *user.User

	mockAuthService.On("ValidatePasswordPolicy", mock.Anything, mock.Anything).Return(nil)
	mockQryRepo.On("ExistsByUsername", mock.Anything, mock.Anything).Return(false, nil)
	mockQryRepo.On("ExistsByEmail", mock.Anything, mock.Anything).Return(false, nil)
	mockAuthService.On("GeneratePasswordHash", mock.Anything, mock.Anything).Return("hashed", nil)
	mockCmdRepo.On("Create", mock.Anything, mock.AnythingOfType("*user.User")).
		Run(func(args mock.Arguments) {
			capturedUser = args.Get(1).(*user.User)
		}).Return(nil)

	handler := NewCreateUserHandler(mockCmdRepo, mockQryRepo, mockAuthService)
	_, err := handler.Handle(context.Background(), CreateUserCommand{
		Username: "test",
		Email:    "test@example.com",
		Password: "password123",
		Status:   &status,
	})

	require.NoError(t, err)
	assert.NotNil(t, capturedUser)
	assert.Equal(t, "inactive", capturedUser.Status)
}
