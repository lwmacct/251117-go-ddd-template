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
	users       []*user.User
	createError error
	idCounter   uint
}

func newMockUserCommandRepo() *mockUserCommandRepo {
	return &mockUserCommandRepo{
		users:     make([]*user.User, 0),
		idCounter: 0,
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

// mockUserQueryRepoForRegister 是用于注册测试的用户查询仓储 Mock 实现。
type mockUserQueryRepoForRegister struct {
	existingUsernames map[string]bool
	existingEmails    map[string]bool
	checkError        error
}

func newMockUserQueryRepoForRegister() *mockUserQueryRepoForRegister {
	return &mockUserQueryRepoForRegister{
		existingUsernames: make(map[string]bool),
		existingEmails:    make(map[string]bool),
	}
}

func (m *mockUserQueryRepoForRegister) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	if m.checkError != nil {
		return false, m.checkError
	}
	return m.existingUsernames[username], nil
}

func (m *mockUserQueryRepoForRegister) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	if m.checkError != nil {
		return false, m.checkError
	}
	return m.existingEmails[email], nil
}

// 实现其他接口方法（空实现）
func (m *mockUserQueryRepoForRegister) GetByID(ctx context.Context, id uint) (*user.User, error) {
	return nil, nil //nolint:nilnil // test mock returns nil for not found
}
func (m *mockUserQueryRepoForRegister) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	return nil, nil //nolint:nilnil // test mock returns nil for not found
}
func (m *mockUserQueryRepoForRegister) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	return nil, nil //nolint:nilnil // test mock returns nil for not found
}
func (m *mockUserQueryRepoForRegister) GetByIDWithRoles(ctx context.Context, id uint) (*user.User, error) {
	return nil, nil //nolint:nilnil // test mock returns nil for not found
}
func (m *mockUserQueryRepoForRegister) GetByUsernameWithRoles(ctx context.Context, username string) (*user.User, error) {
	return nil, nil //nolint:nilnil // test mock returns nil for not found
}
func (m *mockUserQueryRepoForRegister) GetByEmailWithRoles(ctx context.Context, email string) (*user.User, error) {
	return nil, nil //nolint:nilnil // test mock returns nil for not found
}
func (m *mockUserQueryRepoForRegister) List(ctx context.Context, offset, limit int) ([]*user.User, error) {
	return nil, nil
}
func (m *mockUserQueryRepoForRegister) Count(ctx context.Context) (int64, error) { return 0, nil }
func (m *mockUserQueryRepoForRegister) GetRoles(ctx context.Context, userID uint) ([]uint, error) {
	return nil, nil
}
func (m *mockUserQueryRepoForRegister) Search(ctx context.Context, keyword string, offset, limit int) ([]*user.User, error) {
	return nil, nil
}
func (m *mockUserQueryRepoForRegister) CountBySearch(ctx context.Context, keyword string) (int64, error) {
	return 0, nil
}
func (m *mockUserQueryRepoForRegister) Exists(ctx context.Context, id uint) (bool, error) {
	return false, nil
}
func (m *mockUserQueryRepoForRegister) GetUserIDsByRole(ctx context.Context, roleID uint) ([]uint, error) {
	return nil, nil
}

// mockAuthServiceForRegister 是用于注册测试的认证服务 Mock 实现。
type mockAuthServiceForRegister struct {
	weakPassword bool
	hashError    error
	tokenError   error
}

func newMockAuthServiceForRegister() *mockAuthServiceForRegister {
	return &mockAuthServiceForRegister{}
}

func (m *mockAuthServiceForRegister) ValidatePasswordPolicy(ctx context.Context, password string) error {
	if m.weakPassword || len(password) < 6 {
		return domainAuth.ErrWeakPassword
	}
	return nil
}

func (m *mockAuthServiceForRegister) GeneratePasswordHash(ctx context.Context, password string) (string, error) {
	if m.hashError != nil {
		return "", m.hashError
	}
	return "hashed_" + password, nil
}

func (m *mockAuthServiceForRegister) VerifyPassword(ctx context.Context, hashedPassword, plainPassword string) error {
	return nil
}

func (m *mockAuthServiceForRegister) GenerateAccessToken(ctx context.Context, userID uint, username string) (string, time.Time, error) {
	if m.tokenError != nil {
		return "", time.Time{}, m.tokenError
	}
	return "access_token_" + username, time.Now().Add(time.Hour), nil
}

func (m *mockAuthServiceForRegister) GenerateRefreshToken(ctx context.Context, userID uint) (string, time.Time, error) {
	if m.tokenError != nil {
		return "", time.Time{}, m.tokenError
	}
	return "refresh_token", time.Now().Add(24 * time.Hour), nil
}

func (m *mockAuthServiceForRegister) ValidateAccessToken(ctx context.Context, token string) (*domainAuth.TokenClaims, error) {
	return nil, nil //nolint:nilnil // test mock returns nil for not found
}

func (m *mockAuthServiceForRegister) ValidateRefreshToken(ctx context.Context, token string) (uint, error) {
	return 0, nil
}

func (m *mockAuthServiceForRegister) GeneratePATToken(ctx context.Context) (string, error) {
	return "pat", nil
}

func (m *mockAuthServiceForRegister) HashPATToken(ctx context.Context, token string) string {
	return "hashed"
}

type registerTestSetup struct {
	handler     *RegisterHandler
	commandRepo *mockUserCommandRepo
	queryRepo   *mockUserQueryRepoForRegister
	authService *mockAuthServiceForRegister
}

func setupRegisterTest() *registerTestSetup {
	commandRepo := newMockUserCommandRepo()
	queryRepo := newMockUserQueryRepoForRegister()
	authService := newMockAuthServiceForRegister()

	handler := NewRegisterHandler(commandRepo, queryRepo, authService)

	return &registerTestSetup{
		handler:     handler,
		commandRepo: commandRepo,
		queryRepo:   queryRepo,
		authService: authService,
	}
}

func TestRegisterHandler_Handle(t *testing.T) {
	ctx := context.Background()

	t.Run("成功注册新用户", func(t *testing.T) {
		setup := setupRegisterTest()

		cmd := RegisterCommand{
			Username: "newuser",
			Email:    "new@example.com",
			Password: "password123",
			FullName: "New User",
		}

		result, err := setup.handler.Handle(ctx, cmd)

		// 使用 require 验证前置条件
		require.NoError(t, err, "注册应该成功")
		require.NotNil(t, result, "结果不应为空")

		// 使用 assert 验证结果字段
		assert.Equal(t, uint(1), result.UserID, "UserID 应该是 1")
		assert.Equal(t, "newuser", result.Username)
		assert.Equal(t, "new@example.com", result.Email)
		assert.NotEmpty(t, result.AccessToken, "AccessToken 不应为空")
		assert.NotEmpty(t, result.RefreshToken, "RefreshToken 不应为空")
		assert.Equal(t, "Bearer", result.TokenType)

		// 验证用户已创建
		assert.Len(t, setup.commandRepo.users, 1, "应该创建了 1 个用户")
		assert.Equal(t, "hashed_password123", setup.commandRepo.users[0].Password, "密码应该被哈希")
		assert.Equal(t, "active", setup.commandRepo.users[0].Status, "新用户应该是 active 状态")
	})

	t.Run("密码策略验证失败 - 密码太短", func(t *testing.T) {
		setup := setupRegisterTest()

		cmd := RegisterCommand{
			Username: "newuser",
			Email:    "new@example.com",
			Password: "123", // 太短
			FullName: "New User",
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.ErrorIs(t, err, domainAuth.ErrWeakPassword, "应该返回密码太弱错误")
		assert.Nil(t, result, "结果应为空")
	})

	t.Run("用户名已存在", func(t *testing.T) {
		setup := setupRegisterTest()
		setup.queryRepo.existingUsernames["existinguser"] = true

		cmd := RegisterCommand{
			Username: "existinguser",
			Email:    "new@example.com",
			Password: "password123",
			FullName: "New User",
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.ErrorIs(t, err, user.ErrUsernameAlreadyExists, "应该返回用户名已存在错误")
		assert.Nil(t, result)
	})

	t.Run("邮箱已存在", func(t *testing.T) {
		setup := setupRegisterTest()
		setup.queryRepo.existingEmails["existing@example.com"] = true

		cmd := RegisterCommand{
			Username: "newuser",
			Email:    "existing@example.com",
			Password: "password123",
			FullName: "New User",
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.ErrorIs(t, err, user.ErrEmailAlreadyExists, "应该返回邮箱已存在错误")
		assert.Nil(t, result)
	})

	t.Run("检查用户名存在性时发生错误", func(t *testing.T) {
		setup := setupRegisterTest()
		setup.queryRepo.checkError = errors.New("database error")

		cmd := RegisterCommand{
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
		setup := setupRegisterTest()
		setup.authService.hashError = errors.New("hash error")

		cmd := RegisterCommand{
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
		setup := setupRegisterTest()
		setup.commandRepo.createError = errors.New("database error")

		cmd := RegisterCommand{
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

	t.Run("生成访问令牌失败", func(t *testing.T) {
		setup := setupRegisterTest()
		setup.authService.tokenError = errors.New("token error")

		cmd := RegisterCommand{
			Username: "newuser",
			Email:    "new@example.com",
			Password: "password123",
			FullName: "New User",
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to generate access token")
		assert.Nil(t, result)
	})
}

func TestNewRegisterHandler(t *testing.T) {
	t.Run("创建 RegisterHandler", func(t *testing.T) {
		commandRepo := newMockUserCommandRepo()
		queryRepo := newMockUserQueryRepoForRegister()
		authService := newMockAuthServiceForRegister()

		handler := NewRegisterHandler(commandRepo, queryRepo, authService)

		assert.NotNil(t, handler, "Handler 不应为空")
	})
}

func TestRegisterHandler_PasswordValidation(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		password    string
		wantErr     bool
		expectedErr error
	}{
		{
			name:        "有效密码",
			password:    "validPassword123",
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:        "密码太短",
			password:    "123",
			wantErr:     true,
			expectedErr: domainAuth.ErrWeakPassword,
		},
		{
			name:        "空密码",
			password:    "",
			wantErr:     true,
			expectedErr: domainAuth.ErrWeakPassword,
		},
		{
			name:        "刚好 6 位",
			password:    "123456",
			wantErr:     false,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup := setupRegisterTest()

			cmd := RegisterCommand{
				Username: "testuser",
				Email:    "test@example.com",
				Password: tt.password,
				FullName: "Test User",
			}

			_, err := setup.handler.Handle(ctx, cmd)

			if tt.wantErr {
				require.Error(t, err, "应该返回错误")
				if tt.expectedErr != nil {
					require.ErrorIs(t, err, tt.expectedErr, "错误类型应该匹配")
				}
			} else {
				assert.NoError(t, err, "不应该返回错误")
			}
		})
	}
}

func BenchmarkRegisterHandler_Handle(b *testing.B) {
	ctx := context.Background()
	setup := setupRegisterTest()

	for i := 0; b.Loop(); i++ {
		cmd := RegisterCommand{
			Username: "user" + string(rune('0'+i%10)),
			Email:    "user" + string(rune('0'+i%10)) + "@example.com",
			Password: "password123",
			FullName: "Test User",
		}
		_, _ = setup.handler.Handle(ctx, cmd)
	}
}
