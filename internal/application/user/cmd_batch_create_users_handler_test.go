package user

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestBatchCreateUsersHandler_Handle_AllSuccess(t *testing.T) {
	// Arrange
	mockCmdRepo := new(MockUserCommandRepository)
	mockQryRepo := new(MockUserQueryRepository)
	mockAuthService := new(MockAuthService)

	users := []BatchUserItemDTO{
		{Username: "user1", Email: "user1@example.com", Password: "password1"},
		{Username: "user2", Email: "user2@example.com", Password: "password2"},
	}

	for _, u := range users {
		mockAuthService.On("ValidatePasswordPolicy", mock.Anything, u.Password).Return(nil).Once()
		mockQryRepo.On("ExistsByUsername", mock.Anything, u.Username).Return(false, nil).Once()
		mockQryRepo.On("ExistsByEmail", mock.Anything, u.Email).Return(false, nil).Once()
		mockAuthService.On("GeneratePasswordHash", mock.Anything, u.Password).Return("hashed", nil).Once()
		mockCmdRepo.On("Create", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil).Once()
	}

	handler := NewBatchCreateUsersHandler(mockCmdRepo, mockQryRepo, mockAuthService)

	// Act
	result, err := handler.Handle(context.Background(), BatchCreateUsersCommand{Users: users})

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, result.Total)
	assert.Equal(t, 2, result.Success)
	assert.Equal(t, 0, result.Failed)
	assert.Empty(t, result.Errors)
}

func TestBatchCreateUsersHandler_Handle_PartialFailure(t *testing.T) {
	// Arrange
	mockCmdRepo := new(MockUserCommandRepository)
	mockQryRepo := new(MockUserQueryRepository)
	mockAuthService := new(MockAuthService)

	users := []BatchUserItemDTO{
		{Username: "user1", Email: "user1@example.com", Password: "password1"},
		{Username: "existing", Email: "existing@example.com", Password: "password2"},
	}

	// 第一个用户成功
	mockAuthService.On("ValidatePasswordPolicy", mock.Anything, "password1").Return(nil).Once()
	mockQryRepo.On("ExistsByUsername", mock.Anything, "user1").Return(false, nil).Once()
	mockQryRepo.On("ExistsByEmail", mock.Anything, "user1@example.com").Return(false, nil).Once()
	mockAuthService.On("GeneratePasswordHash", mock.Anything, "password1").Return("hashed", nil).Once()
	mockCmdRepo.On("Create", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil).Once()

	// 第二个用户失败（用户名已存在）
	mockAuthService.On("ValidatePasswordPolicy", mock.Anything, "password2").Return(nil).Once()
	mockQryRepo.On("ExistsByUsername", mock.Anything, "existing").Return(true, nil).Once()

	handler := NewBatchCreateUsersHandler(mockCmdRepo, mockQryRepo, mockAuthService)

	// Act
	result, err := handler.Handle(context.Background(), BatchCreateUsersCommand{Users: users})

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, result.Total)
	assert.Equal(t, 1, result.Success)
	assert.Equal(t, 1, result.Failed)
	assert.Len(t, result.Errors, 1)
	assert.Equal(t, 1, result.Errors[0].Index)
	assert.Equal(t, "existing", result.Errors[0].Username)
	assert.Contains(t, result.Errors[0].Error, "用户名已存在")
}

func TestBatchCreateUsersHandler_Handle_ValidationError(t *testing.T) {
	// Arrange
	mockCmdRepo := new(MockUserCommandRepository)
	mockQryRepo := new(MockUserQueryRepository)
	mockAuthService := new(MockAuthService)

	users := []BatchUserItemDTO{
		{Username: "ab", Email: "invalid", Password: ""},
	}

	handler := NewBatchCreateUsersHandler(mockCmdRepo, mockQryRepo, mockAuthService)

	// Act
	result, err := handler.Handle(context.Background(), BatchCreateUsersCommand{Users: users})

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.Total)
	assert.Equal(t, 0, result.Success)
	assert.Equal(t, 1, result.Failed)
	assert.Len(t, result.Errors, 1)
}

func TestBatchCreateUsersHandler_Handle_WithRoles(t *testing.T) {
	// Arrange
	mockCmdRepo := new(MockUserCommandRepository)
	mockQryRepo := new(MockUserQueryRepository)
	mockAuthService := new(MockAuthService)

	users := []BatchUserItemDTO{
		{Username: "user1", Email: "user1@example.com", Password: "password1", RoleIDs: []uint{1, 2}},
	}

	mockAuthService.On("ValidatePasswordPolicy", mock.Anything, "password1").Return(nil)
	mockQryRepo.On("ExistsByUsername", mock.Anything, "user1").Return(false, nil)
	mockQryRepo.On("ExistsByEmail", mock.Anything, "user1@example.com").Return(false, nil)
	mockAuthService.On("GeneratePasswordHash", mock.Anything, "password1").Return("hashed", nil)
	mockCmdRepo.On("Create", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)
	mockCmdRepo.On("AssignRoles", mock.Anything, uint(1), []uint{1, 2}).Return(nil)

	handler := NewBatchCreateUsersHandler(mockCmdRepo, mockQryRepo, mockAuthService)

	// Act
	result, err := handler.Handle(context.Background(), BatchCreateUsersCommand{Users: users})

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.Success)
	mockCmdRepo.AssertExpectations(t)
}

func TestBatchCreateUsersHandler_Handle_RoleAssignmentFailure(t *testing.T) {
	// Arrange
	mockCmdRepo := new(MockUserCommandRepository)
	mockQryRepo := new(MockUserQueryRepository)
	mockAuthService := new(MockAuthService)

	users := []BatchUserItemDTO{
		{Username: "user1", Email: "user1@example.com", Password: "password1", RoleIDs: []uint{999}},
	}

	mockAuthService.On("ValidatePasswordPolicy", mock.Anything, "password1").Return(nil)
	mockQryRepo.On("ExistsByUsername", mock.Anything, "user1").Return(false, nil)
	mockQryRepo.On("ExistsByEmail", mock.Anything, "user1@example.com").Return(false, nil)
	mockAuthService.On("GeneratePasswordHash", mock.Anything, "password1").Return("hashed", nil)
	mockCmdRepo.On("Create", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)
	mockCmdRepo.On("AssignRoles", mock.Anything, uint(1), []uint{999}).Return(errors.New("role not found"))

	handler := NewBatchCreateUsersHandler(mockCmdRepo, mockQryRepo, mockAuthService)

	// Act
	result, err := handler.Handle(context.Background(), BatchCreateUsersCommand{Users: users})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, 1, result.Failed)
	assert.Contains(t, result.Errors[0].Error, "角色分配失败")
}

// ============================================================
// 补充覆盖率测试
// ============================================================

func TestBatchCreateUsersHandler_Handle_EmailValidationError(t *testing.T) {
	mockCmdRepo := new(MockUserCommandRepository)
	mockQryRepo := new(MockUserQueryRepository)
	mockAuthService := new(MockAuthService)

	tests := []struct {
		name    string
		users   []BatchUserItemDTO
		wantErr string
	}{
		{
			name:    "用户名为空",
			users:   []BatchUserItemDTO{{Username: "", Email: "test@example.com", Password: "password"}},
			wantErr: "用户名不能为空",
		},
		{
			name:    "用户名只有空格",
			users:   []BatchUserItemDTO{{Username: "   ", Email: "test@example.com", Password: "password"}},
			wantErr: "用户名不能为空",
		},
		{
			name:    "邮箱为空",
			users:   []BatchUserItemDTO{{Username: "validuser", Email: "", Password: "password"}},
			wantErr: "邮箱不能为空",
		},
		{
			name:    "邮箱格式无效",
			users:   []BatchUserItemDTO{{Username: "validuser", Email: "invalid-email", Password: "password"}},
			wantErr: "邮箱格式无效",
		},
		{
			name:    "密码为空",
			users:   []BatchUserItemDTO{{Username: "validuser", Email: "test@example.com", Password: ""}},
			wantErr: "密码不能为空",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewBatchCreateUsersHandler(mockCmdRepo, mockQryRepo, mockAuthService)

			result, err := handler.Handle(context.Background(), BatchCreateUsersCommand{Users: tt.users})

			require.NoError(t, err)
			assert.Equal(t, 1, result.Failed)
			assert.Contains(t, result.Errors[0].Error, tt.wantErr)
		})
	}
}

func TestBatchCreateUsersHandler_Handle_PasswordPolicyError(t *testing.T) {
	mockCmdRepo := new(MockUserCommandRepository)
	mockQryRepo := new(MockUserQueryRepository)
	mockAuthService := new(MockAuthService)

	users := []BatchUserItemDTO{
		{Username: "validuser", Email: "test@example.com", Password: "weak"},
	}

	mockAuthService.On("ValidatePasswordPolicy", mock.Anything, "weak").Return(errors.New("密码太弱"))

	handler := NewBatchCreateUsersHandler(mockCmdRepo, mockQryRepo, mockAuthService)

	result, err := handler.Handle(context.Background(), BatchCreateUsersCommand{Users: users})

	require.NoError(t, err)
	assert.Equal(t, 1, result.Failed)
	assert.Contains(t, result.Errors[0].Error, "密码不符合策略")
}

func TestBatchCreateUsersHandler_Handle_CheckUsernameError(t *testing.T) {
	mockCmdRepo := new(MockUserCommandRepository)
	mockQryRepo := new(MockUserQueryRepository)
	mockAuthService := new(MockAuthService)

	users := []BatchUserItemDTO{
		{Username: "validuser", Email: "test@example.com", Password: "password"},
	}

	mockAuthService.On("ValidatePasswordPolicy", mock.Anything, "password").Return(nil)
	mockQryRepo.On("ExistsByUsername", mock.Anything, "validuser").Return(false, errors.New("db error"))

	handler := NewBatchCreateUsersHandler(mockCmdRepo, mockQryRepo, mockAuthService)

	result, err := handler.Handle(context.Background(), BatchCreateUsersCommand{Users: users})

	require.NoError(t, err)
	assert.Equal(t, 1, result.Failed)
	assert.Contains(t, result.Errors[0].Error, "检查用户名失败")
}

func TestBatchCreateUsersHandler_Handle_CheckEmailError(t *testing.T) {
	mockCmdRepo := new(MockUserCommandRepository)
	mockQryRepo := new(MockUserQueryRepository)
	mockAuthService := new(MockAuthService)

	users := []BatchUserItemDTO{
		{Username: "validuser", Email: "test@example.com", Password: "password"},
	}

	mockAuthService.On("ValidatePasswordPolicy", mock.Anything, "password").Return(nil)
	mockQryRepo.On("ExistsByUsername", mock.Anything, "validuser").Return(false, nil)
	mockQryRepo.On("ExistsByEmail", mock.Anything, "test@example.com").Return(false, errors.New("db error"))

	handler := NewBatchCreateUsersHandler(mockCmdRepo, mockQryRepo, mockAuthService)

	result, err := handler.Handle(context.Background(), BatchCreateUsersCommand{Users: users})

	require.NoError(t, err)
	assert.Equal(t, 1, result.Failed)
	assert.Contains(t, result.Errors[0].Error, "检查邮箱失败")
}

func TestBatchCreateUsersHandler_Handle_EmailExists(t *testing.T) {
	mockCmdRepo := new(MockUserCommandRepository)
	mockQryRepo := new(MockUserQueryRepository)
	mockAuthService := new(MockAuthService)

	users := []BatchUserItemDTO{
		{Username: "validuser", Email: "existing@example.com", Password: "password"},
	}

	mockAuthService.On("ValidatePasswordPolicy", mock.Anything, "password").Return(nil)
	mockQryRepo.On("ExistsByUsername", mock.Anything, "validuser").Return(false, nil)
	mockQryRepo.On("ExistsByEmail", mock.Anything, "existing@example.com").Return(true, nil)

	handler := NewBatchCreateUsersHandler(mockCmdRepo, mockQryRepo, mockAuthService)

	result, err := handler.Handle(context.Background(), BatchCreateUsersCommand{Users: users})

	require.NoError(t, err)
	assert.Equal(t, 1, result.Failed)
	assert.Contains(t, result.Errors[0].Error, "邮箱已存在")
}

func TestBatchCreateUsersHandler_Handle_HashPasswordError(t *testing.T) {
	mockCmdRepo := new(MockUserCommandRepository)
	mockQryRepo := new(MockUserQueryRepository)
	mockAuthService := new(MockAuthService)

	users := []BatchUserItemDTO{
		{Username: "validuser", Email: "test@example.com", Password: "password"},
	}

	mockAuthService.On("ValidatePasswordPolicy", mock.Anything, "password").Return(nil)
	mockQryRepo.On("ExistsByUsername", mock.Anything, "validuser").Return(false, nil)
	mockQryRepo.On("ExistsByEmail", mock.Anything, "test@example.com").Return(false, nil)
	mockAuthService.On("GeneratePasswordHash", mock.Anything, "password").Return("", errors.New("hash error"))

	handler := NewBatchCreateUsersHandler(mockCmdRepo, mockQryRepo, mockAuthService)

	result, err := handler.Handle(context.Background(), BatchCreateUsersCommand{Users: users})

	require.NoError(t, err)
	assert.Equal(t, 1, result.Failed)
	assert.Contains(t, result.Errors[0].Error, "密码加密失败")
}

func TestBatchCreateUsersHandler_Handle_InvalidStatus(t *testing.T) {
	mockCmdRepo := new(MockUserCommandRepository)
	mockQryRepo := new(MockUserQueryRepository)
	mockAuthService := new(MockAuthService)

	users := []BatchUserItemDTO{
		{Username: "validuser", Email: "test@example.com", Password: "password", Status: "invalid_status"},
	}

	mockAuthService.On("ValidatePasswordPolicy", mock.Anything, "password").Return(nil)
	mockQryRepo.On("ExistsByUsername", mock.Anything, "validuser").Return(false, nil)
	mockQryRepo.On("ExistsByEmail", mock.Anything, "test@example.com").Return(false, nil)
	mockAuthService.On("GeneratePasswordHash", mock.Anything, "password").Return("hashed", nil)

	handler := NewBatchCreateUsersHandler(mockCmdRepo, mockQryRepo, mockAuthService)

	result, err := handler.Handle(context.Background(), BatchCreateUsersCommand{Users: users})

	require.NoError(t, err)
	assert.Equal(t, 1, result.Failed)
	assert.Contains(t, result.Errors[0].Error, "无效的状态值")
}

func TestBatchCreateUsersHandler_Handle_CreateUserError(t *testing.T) {
	mockCmdRepo := new(MockUserCommandRepository)
	mockQryRepo := new(MockUserQueryRepository)
	mockAuthService := new(MockAuthService)

	users := []BatchUserItemDTO{
		{Username: "validuser", Email: "test@example.com", Password: "password"},
	}

	mockAuthService.On("ValidatePasswordPolicy", mock.Anything, "password").Return(nil)
	mockQryRepo.On("ExistsByUsername", mock.Anything, "validuser").Return(false, nil)
	mockQryRepo.On("ExistsByEmail", mock.Anything, "test@example.com").Return(false, nil)
	mockAuthService.On("GeneratePasswordHash", mock.Anything, "password").Return("hashed", nil)
	mockCmdRepo.On("Create", mock.Anything, mock.AnythingOfType("*user.User")).Return(errors.New("db error"))

	handler := NewBatchCreateUsersHandler(mockCmdRepo, mockQryRepo, mockAuthService)

	result, err := handler.Handle(context.Background(), BatchCreateUsersCommand{Users: users})

	require.NoError(t, err)
	assert.Equal(t, 1, result.Failed)
	assert.Contains(t, result.Errors[0].Error, "创建用户失败")
}
