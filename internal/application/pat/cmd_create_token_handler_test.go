package pat

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	domainPAT "github.com/lwmacct/251117-go-ddd-template/internal/domain/pat"
	domainRole "github.com/lwmacct/251117-go-ddd-template/internal/domain/role"
	domainUser "github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
)

func TestCreateTokenHandler_Handle_Success(t *testing.T) {
	// Arrange
	mockPATCmdRepo := new(MockPATCommandRepository)
	mockUserQryRepo := new(MockUserQueryRepository)
	mockTokenGen := new(MockTokenGenerator)

	expiresAt := time.Now().Add(30 * 24 * time.Hour)

	user := &domainUser.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Status:   "active",
		Roles: []domainRole.Role{
			{
				ID:   1,
				Name: "admin",
				Permissions: []domainRole.Permission{
					{ID: 1, Code: "user:read", Description: "Read Users"},
					{ID: 2, Code: "user:write", Description: "Write Users"},
				},
			},
		},
	}

	mockUserQryRepo.On("GetByIDWithRoles", mock.Anything, uint(1)).Return(user, nil)
	mockTokenGen.On("GeneratePAT").Return("pat_ABC12_plaintoken", "hashed_token", "pat_ABC12", nil)
	mockPATCmdRepo.On("Create", mock.Anything, mock.AnythingOfType("*pat.PersonalAccessToken")).Return(nil)

	handler := NewCreateTokenHandler(mockPATCmdRepo, mockUserQryRepo, mockTokenGen)

	// Act
	result, err := handler.Handle(context.Background(), CreateTokenCommand{
		UserID:      1,
		Name:        "Test Token",
		Permissions: []string{"user:read"},
		ExpiresAt:   &expiresAt,
		IPWhitelist: []string{"127.0.0.1"},
		Description: "Test token description",
	})

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "pat_ABC12_plaintoken", result.PlainToken)
	assert.NotNil(t, result.Token)
	assert.Equal(t, uint(1), result.Token.ID)
	assert.Equal(t, uint(1), result.Token.UserID)
	assert.Equal(t, "Test Token", result.Token.Name)
	assert.Equal(t, "hashed_token", result.Token.Token)
	assert.Equal(t, "pat_ABC12", result.Token.TokenPrefix)
	assert.Equal(t, []string{"user:read"}, []string(result.Token.Permissions))
	assert.Equal(t, "active", result.Token.Status)

	mockUserQryRepo.AssertExpectations(t)
	mockTokenGen.AssertExpectations(t)
	mockPATCmdRepo.AssertExpectations(t)
}

func TestCreateTokenHandler_Handle_InheritAllPermissions(t *testing.T) {
	// Arrange
	mockPATCmdRepo := new(MockPATCommandRepository)
	mockUserQryRepo := new(MockUserQueryRepository)
	mockTokenGen := new(MockTokenGenerator)

	user := &domainUser.User{
		ID:       1,
		Username: "testuser",
		Status:   "active",
		Roles: []domainRole.Role{
			{
				ID:   1,
				Name: "admin",
				Permissions: []domainRole.Permission{
					{ID: 1, Code: "user:read", Description: "Read Users"},
					{ID: 2, Code: "user:write", Description: "Write Users"},
					{ID: 3, Code: "role:read", Description: "Read Roles"},
				},
			},
		},
	}

	mockUserQryRepo.On("GetByIDWithRoles", mock.Anything, uint(1)).Return(user, nil)
	mockTokenGen.On("GeneratePAT").Return("pat_XYZ99_token", "hashed", "pat_XYZ99", nil)
	mockPATCmdRepo.On("Create", mock.Anything, mock.MatchedBy(func(pat *domainPAT.PersonalAccessToken) bool {
		// 验证继承了所有权限
		return len(pat.Permissions) == 3 &&
			pat.Permissions[0] == "user:read" &&
			pat.Permissions[1] == "user:write" &&
			pat.Permissions[2] == "role:read"
	})).Return(nil)

	handler := NewCreateTokenHandler(mockPATCmdRepo, mockUserQryRepo, mockTokenGen)

	// Act - 不传 Permissions，应继承所有权限
	result, err := handler.Handle(context.Background(), CreateTokenCommand{
		UserID: 1,
		Name:   "Full Access Token",
	})

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Token.Permissions, 3)

	mockUserQryRepo.AssertExpectations(t)
	mockTokenGen.AssertExpectations(t)
	mockPATCmdRepo.AssertExpectations(t)
}

func TestCreateTokenHandler_Handle_Error(t *testing.T) {
	tests := []struct {
		name       string
		cmd        CreateTokenCommand
		setupMocks func(*MockPATCommandRepository, *MockUserQueryRepository, *MockTokenGenerator)
		wantErr    string
	}{
		{
			name:    "用户 ID 为空",
			cmd:     CreateTokenCommand{Name: "Test"},
			wantErr: "user ID is required",
			setupMocks: func(patCmd *MockPATCommandRepository, userQry *MockUserQueryRepository, tokenGen *MockTokenGenerator) {
				// 不需要设置 mock
			},
		},
		{
			name: "用户不存在",
			cmd:  CreateTokenCommand{UserID: 999, Name: "Test"},
			setupMocks: func(patCmd *MockPATCommandRepository, userQry *MockUserQueryRepository, tokenGen *MockTokenGenerator) {
				userQry.On("GetByIDWithRoles", mock.Anything, uint(999)).Return(nil, errors.New("not found"))
			},
			wantErr: "user not found",
		},
		{
			name: "用户没有权限",
			cmd:  CreateTokenCommand{UserID: 1, Name: "Test"},
			setupMocks: func(patCmd *MockPATCommandRepository, userQry *MockUserQueryRepository, tokenGen *MockTokenGenerator) {
				user := &domainUser.User{
					ID:       1,
					Username: "noPermUser",
					Status:   "active",
					Roles:    []domainRole.Role{}, // 没有角色，没有权限
				}
				userQry.On("GetByIDWithRoles", mock.Anything, uint(1)).Return(user, nil)
			},
			wantErr: "user has no permissions",
		},
		{
			name: "请求的权限不在用户权限范围内",
			cmd: CreateTokenCommand{
				UserID:      1,
				Name:        "Test",
				Permissions: []string{"admin:super"}, // 用户没有此权限
			},
			setupMocks: func(patCmd *MockPATCommandRepository, userQry *MockUserQueryRepository, tokenGen *MockTokenGenerator) {
				user := &domainUser.User{
					ID:       1,
					Username: "testuser",
					Status:   "active",
					Roles: []domainRole.Role{
						{
							ID:   1,
							Name: "user",
							Permissions: []domainRole.Permission{
								{ID: 1, Code: "user:read", Description: "Read Users"},
							},
						},
					},
				}
				userQry.On("GetByIDWithRoles", mock.Anything, uint(1)).Return(user, nil)
			},
			wantErr: "permission 'admin:super' is not granted to user",
		},
		{
			name: "生成 Token 失败",
			cmd: CreateTokenCommand{
				UserID:      1,
				Name:        "Test",
				Permissions: []string{"user:read"},
			},
			setupMocks: func(patCmd *MockPATCommandRepository, userQry *MockUserQueryRepository, tokenGen *MockTokenGenerator) {
				user := &domainUser.User{
					ID:       1,
					Username: "testuser",
					Status:   "active",
					Roles: []domainRole.Role{
						{
							ID:   1,
							Name: "user",
							Permissions: []domainRole.Permission{
								{ID: 1, Code: "user:read", Description: "Read Users"},
							},
						},
					},
				}
				userQry.On("GetByIDWithRoles", mock.Anything, uint(1)).Return(user, nil)
				tokenGen.On("GeneratePAT").Return("", "", "", errors.New("crypto error"))
			},
			wantErr: "failed to generate token",
		},
		{
			name: "创建 Token 失败",
			cmd: CreateTokenCommand{
				UserID:      1,
				Name:        "Test",
				Permissions: []string{"user:read"},
			},
			setupMocks: func(patCmd *MockPATCommandRepository, userQry *MockUserQueryRepository, tokenGen *MockTokenGenerator) {
				user := &domainUser.User{
					ID:       1,
					Username: "testuser",
					Status:   "active",
					Roles: []domainRole.Role{
						{
							ID:   1,
							Name: "user",
							Permissions: []domainRole.Permission{
								{ID: 1, Code: "user:read", Description: "Read Users"},
							},
						},
					},
				}
				userQry.On("GetByIDWithRoles", mock.Anything, uint(1)).Return(user, nil)
				tokenGen.On("GeneratePAT").Return("plain", "hashed", "prefix", nil)
				patCmd.On("Create", mock.Anything, mock.AnythingOfType("*pat.PersonalAccessToken")).Return(errors.New("db error"))
			},
			wantErr: "failed to create token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockPATCmdRepo := new(MockPATCommandRepository)
			mockUserQryRepo := new(MockUserQueryRepository)
			mockTokenGen := new(MockTokenGenerator)

			tt.setupMocks(mockPATCmdRepo, mockUserQryRepo, mockTokenGen)

			handler := NewCreateTokenHandler(mockPATCmdRepo, mockUserQryRepo, mockTokenGen)

			// Act
			result, err := handler.Handle(context.Background(), tt.cmd)

			// Assert
			assert.Nil(t, result)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestValidatePermissions(t *testing.T) {
	tests := []struct {
		name        string
		requested   []string
		userPerms   []string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "验证通过 - 所有请求权限都在用户权限内",
			requested:   []string{"user:read", "user:write"},
			userPerms:   []string{"user:read", "user:write", "role:read"},
			expectError: false,
		},
		{
			name:        "验证通过 - 请求单个权限",
			requested:   []string{"user:read"},
			userPerms:   []string{"user:read", "user:write"},
			expectError: false,
		},
		{
			name:        "验证失败 - 请求权限为空",
			requested:   []string{},
			userPerms:   []string{"user:read"},
			expectError: true,
			errorMsg:    "at least one permission is required",
		},
		{
			name:        "验证失败 - 请求权限不在用户权限内",
			requested:   []string{"admin:super"},
			userPerms:   []string{"user:read", "user:write"},
			expectError: true,
			errorMsg:    "permission 'admin:super' is not granted to user",
		},
		{
			name:        "验证失败 - 部分权限不在用户权限内",
			requested:   []string{"user:read", "admin:super"},
			userPerms:   []string{"user:read", "user:write"},
			expectError: true,
			errorMsg:    "permission 'admin:super' is not granted to user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePermissions(tt.requested, tt.userPerms)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
