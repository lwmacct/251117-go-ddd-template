package command

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domainAuth "github.com/lwmacct/251117-go-ddd-template/internal/domain/auth"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
)

// mockUserCommandRepo 是用户命令仓储的 Mock 实现。
type mockUserCommandRepo struct {
	users            []*user.User
	createError      error
	assignRolesError error
	idCounter        uint
	assignedRoles    map[uint][]uint // userID -> roleIDs
}

func newMockUserCommandRepo() *mockUserCommandRepo {
	return &mockUserCommandRepo{
		users:         make([]*user.User, 0),
		idCounter:     0,
		assignedRoles: make(map[uint][]uint),
	}
}

func (m *mockUserCommandRepo) Create(ctx context.Context, u *user.User) error {
	if m.createError != nil {
		return m.createError
	}
	m.idCounter++
	u.ID = m.idCounter
	m.users = append(m.users, u)
	return nil
}

func (m *mockUserCommandRepo) Update(ctx context.Context, u *user.User) error {
	return nil
}

func (m *mockUserCommandRepo) Delete(ctx context.Context, id uint) error {
	return nil
}

func (m *mockUserCommandRepo) AssignRoles(ctx context.Context, userID uint, roleIDs []uint) error {
	if m.assignRolesError != nil {
		return m.assignRolesError
	}
	m.assignedRoles[userID] = roleIDs
	return nil
}

func (m *mockUserCommandRepo) RemoveRoles(ctx context.Context, userID uint, roleIDs []uint) error {
	return nil
}

func (m *mockUserCommandRepo) UpdatePassword(ctx context.Context, userID uint, hashedPassword string) error {
	return nil
}

func (m *mockUserCommandRepo) UpdateStatus(ctx context.Context, userID uint, status string) error {
	return nil
}

// mockUserQueryRepo 是用户查询仓储的 Mock 实现。
type mockUserQueryRepo struct {
	existingUsernames map[string]bool
	existingEmails    map[string]bool
	checkError        error
}

func newMockUserQueryRepo() *mockUserQueryRepo {
	return &mockUserQueryRepo{
		existingUsernames: make(map[string]bool),
		existingEmails:    make(map[string]bool),
	}
}

func (m *mockUserQueryRepo) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	if m.checkError != nil {
		return false, m.checkError
	}
	return m.existingUsernames[username], nil
}

func (m *mockUserQueryRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	if m.checkError != nil {
		return false, m.checkError
	}
	return m.existingEmails[email], nil
}

// 实现其他接口方法（空实现）
func (m *mockUserQueryRepo) GetByID(ctx context.Context, id uint) (*user.User, error) {
	return nil, nil //nolint:nilnil // test mock returns nil for not found
}
func (m *mockUserQueryRepo) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	return nil, nil //nolint:nilnil // test mock returns nil for not found
}
func (m *mockUserQueryRepo) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	return nil, nil //nolint:nilnil // test mock returns nil for not found
}
func (m *mockUserQueryRepo) GetByIDWithRoles(ctx context.Context, id uint) (*user.User, error) {
	return nil, nil //nolint:nilnil // test mock returns nil for not found
}
func (m *mockUserQueryRepo) GetByUsernameWithRoles(ctx context.Context, username string) (*user.User, error) {
	return nil, nil //nolint:nilnil // test mock returns nil for not found
}
func (m *mockUserQueryRepo) GetByEmailWithRoles(ctx context.Context, email string) (*user.User, error) {
	return nil, nil //nolint:nilnil // test mock returns nil for not found
}
func (m *mockUserQueryRepo) List(ctx context.Context, offset, limit int) ([]*user.User, error) {
	return nil, nil
}
func (m *mockUserQueryRepo) Count(ctx context.Context) (int64, error) { return 0, nil }
func (m *mockUserQueryRepo) GetRoles(ctx context.Context, userID uint) ([]uint, error) {
	return nil, nil
}
func (m *mockUserQueryRepo) Search(ctx context.Context, keyword string, offset, limit int) ([]*user.User, error) {
	return nil, nil
}
func (m *mockUserQueryRepo) CountBySearch(ctx context.Context, keyword string) (int64, error) {
	return 0, nil
}
func (m *mockUserQueryRepo) Exists(ctx context.Context, id uint) (bool, error) { return false, nil }
func (m *mockUserQueryRepo) GetUserIDsByRole(ctx context.Context, roleID uint) ([]uint, error) {
	return nil, nil
}

// mockAuthService 是认证服务的 Mock 实现。
type mockAuthService struct {
	weakPassword bool
	hashError    error
}

func newMockAuthService() *mockAuthService {
	return &mockAuthService{}
}

func (m *mockAuthService) ValidatePasswordPolicy(ctx context.Context, password string) error {
	if m.weakPassword || len(password) < 6 {
		return domainAuth.ErrWeakPassword
	}
	return nil
}

func (m *mockAuthService) GeneratePasswordHash(ctx context.Context, password string) (string, error) {
	if m.hashError != nil {
		return "", m.hashError
	}
	return "hashed_" + password, nil
}

func (m *mockAuthService) VerifyPassword(ctx context.Context, hashedPassword, plainPassword string) error {
	return nil
}

func (m *mockAuthService) GenerateAccessToken(ctx context.Context, userID uint, username string) (string, time.Time, error) {
	return "token", time.Now().Add(time.Hour), nil
}

func (m *mockAuthService) GenerateRefreshToken(ctx context.Context, userID uint) (string, time.Time, error) {
	return "refresh", time.Now().Add(24 * time.Hour), nil
}

func (m *mockAuthService) ValidateAccessToken(ctx context.Context, token string) (*domainAuth.TokenClaims, error) {
	return nil, nil //nolint:nilnil // test mock returns nil for not found
}

func (m *mockAuthService) ValidateRefreshToken(ctx context.Context, token string) (uint, error) {
	return 0, nil
}

func (m *mockAuthService) GeneratePATToken(ctx context.Context) (string, error)  { return "pat", nil }
func (m *mockAuthService) HashPATToken(ctx context.Context, token string) string { return "hashed" }

type createUserTestSetup struct {
	handler     *CreateUserHandler
	commandRepo *mockUserCommandRepo
	queryRepo   *mockUserQueryRepo
	authService *mockAuthService
}

func setupCreateUserTest() *createUserTestSetup {
	commandRepo := newMockUserCommandRepo()
	queryRepo := newMockUserQueryRepo()
	authService := newMockAuthService()

	handler := NewCreateUserHandler(commandRepo, queryRepo, authService)

	return &createUserTestSetup{
		handler:     handler,
		commandRepo: commandRepo,
		queryRepo:   queryRepo,
		authService: authService,
	}
}

func TestCreateUserHandler_Handle(t *testing.T) {
	ctx := context.Background()

	t.Run("成功创建用户 - 无角色", func(t *testing.T) {
		setup := setupCreateUserTest()

		cmd := CreateUserCommand{
			Username: "newuser",
			Email:    "new@example.com",
			Password: "password123",
			FullName: "New User",
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		require.NotNil(t, result)

		assert.Equal(t, uint(1), result.UserID)
		assert.Equal(t, "newuser", result.Username)
		assert.Equal(t, "new@example.com", result.Email)

		// 验证用户已创建
		require.Len(t, setup.commandRepo.users, 1)
		assert.Equal(t, "hashed_password123", setup.commandRepo.users[0].Password)
		assert.Equal(t, "active", setup.commandRepo.users[0].Status)
	})

	t.Run("成功创建用户 - 带角色", func(t *testing.T) {
		setup := setupCreateUserTest()

		cmd := CreateUserCommand{
			Username: "adminuser",
			Email:    "admin@example.com",
			Password: "password123",
			FullName: "Admin User",
			RoleIDs:  []uint{1, 2, 3}, // 分配 3 个角色
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		require.NotNil(t, result)

		// 验证角色分配
		assignedRoles, exists := setup.commandRepo.assignedRoles[result.UserID]
		assert.True(t, exists, "应该有角色分配记录")
		assert.Equal(t, []uint{1, 2, 3}, assignedRoles)
	})

	t.Run("密码策略验证失败", func(t *testing.T) {
		setup := setupCreateUserTest()

		cmd := CreateUserCommand{
			Username: "newuser",
			Email:    "new@example.com",
			Password: "123", // 太短
			FullName: "New User",
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.ErrorIs(t, err, domainAuth.ErrWeakPassword)
		assert.Nil(t, result)
	})

	t.Run("用户名已存在", func(t *testing.T) {
		setup := setupCreateUserTest()
		setup.queryRepo.existingUsernames["existinguser"] = true

		cmd := CreateUserCommand{
			Username: "existinguser",
			Email:    "new@example.com",
			Password: "password123",
			FullName: "New User",
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.ErrorIs(t, err, user.ErrUsernameAlreadyExists)
		assert.Nil(t, result)
	})

	t.Run("邮箱已存在", func(t *testing.T) {
		setup := setupCreateUserTest()
		setup.queryRepo.existingEmails["existing@example.com"] = true

		cmd := CreateUserCommand{
			Username: "newuser",
			Email:    "existing@example.com",
			Password: "password123",
			FullName: "New User",
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.ErrorIs(t, err, user.ErrEmailAlreadyExists)
		assert.Nil(t, result)
	})

	t.Run("检查用户名存在性失败", func(t *testing.T) {
		setup := setupCreateUserTest()
		setup.queryRepo.checkError = errors.New("database error")

		cmd := CreateUserCommand{
			Username: "newuser",
			Email:    "new@example.com",
			Password: "password123",
			FullName: "New User",
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to check username existence")
		assert.Nil(t, result)
	})

	t.Run("密码哈希生成失败", func(t *testing.T) {
		setup := setupCreateUserTest()
		setup.authService.hashError = errors.New("hash error")

		cmd := CreateUserCommand{
			Username: "newuser",
			Email:    "new@example.com",
			Password: "password123",
			FullName: "New User",
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to hash password")
		assert.Nil(t, result)
	})

	t.Run("创建用户失败", func(t *testing.T) {
		setup := setupCreateUserTest()
		setup.commandRepo.createError = errors.New("database error")

		cmd := CreateUserCommand{
			Username: "newuser",
			Email:    "new@example.com",
			Password: "password123",
			FullName: "New User",
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create user")
		assert.Nil(t, result)
	})

	t.Run("分配角色失败", func(t *testing.T) {
		setup := setupCreateUserTest()
		setup.commandRepo.assignRolesError = errors.New("role assignment error")

		cmd := CreateUserCommand{
			Username: "newuser",
			Email:    "new@example.com",
			Password: "password123",
			FullName: "New User",
			RoleIDs:  []uint{1, 2},
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to assign roles")
		assert.Nil(t, result)

		// 用户应该已创建（即使角色分配失败）
		assert.Len(t, setup.commandRepo.users, 1)
	})

	t.Run("空角色列表不调用分配", func(t *testing.T) {
		setup := setupCreateUserTest()

		cmd := CreateUserCommand{
			Username: "newuser",
			Email:    "new@example.com",
			Password: "password123",
			FullName: "New User",
			RoleIDs:  []uint{}, // 空角色列表
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		require.NotNil(t, result)

		// 不应该有角色分配记录
		_, exists := setup.commandRepo.assignedRoles[result.UserID]
		assert.False(t, exists, "空角色列表不应该触发角色分配")
	})
}

func TestNewCreateUserHandler(t *testing.T) {
	t.Run("创建 CreateUserHandler", func(t *testing.T) {
		commandRepo := newMockUserCommandRepo()
		queryRepo := newMockUserQueryRepo()
		authService := newMockAuthService()

		handler := NewCreateUserHandler(commandRepo, queryRepo, authService)

		assert.NotNil(t, handler)
	})
}

func TestCreateUserHandler_ValidationScenarios(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		cmd           CreateUserCommand
		setupFunc     func(*createUserTestSetup)
		wantErr       bool
		expectedError error
	}{
		{
			name: "有效输入",
			cmd: CreateUserCommand{
				Username: "validuser",
				Email:    "valid@example.com",
				Password: "password123",
				FullName: "Valid User",
			},
			wantErr: false,
		},
		{
			name: "密码太短",
			cmd: CreateUserCommand{
				Username: "user",
				Email:    "user@example.com",
				Password: "12345",
				FullName: "User",
			},
			wantErr:       true,
			expectedError: domainAuth.ErrWeakPassword,
		},
		{
			name: "用户名冲突",
			cmd: CreateUserCommand{
				Username: "taken",
				Email:    "new@example.com",
				Password: "password123",
				FullName: "User",
			},
			setupFunc: func(s *createUserTestSetup) {
				s.queryRepo.existingUsernames["taken"] = true
			},
			wantErr:       true,
			expectedError: user.ErrUsernameAlreadyExists,
		},
		{
			name: "邮箱冲突",
			cmd: CreateUserCommand{
				Username: "newuser",
				Email:    "taken@example.com",
				Password: "password123",
				FullName: "User",
			},
			setupFunc: func(s *createUserTestSetup) {
				s.queryRepo.existingEmails["taken@example.com"] = true
			},
			wantErr:       true,
			expectedError: user.ErrEmailAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup := setupCreateUserTest()
			if tt.setupFunc != nil {
				tt.setupFunc(setup)
			}

			result, err := setup.handler.Handle(ctx, tt.cmd)

			if tt.wantErr {
				require.Error(t, err)
				if tt.expectedError != nil {
					require.ErrorIs(t, err, tt.expectedError)
				}
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func BenchmarkCreateUserHandler_Handle(b *testing.B) {
	ctx := context.Background()
	setup := setupCreateUserTest()

	for b.Loop() {
		cmd := CreateUserCommand{
			Username: "user",
			Email:    "user@example.com",
			Password: "password123",
			FullName: "Test User",
		}
		_, _ = setup.handler.Handle(ctx, cmd)
	}
}
