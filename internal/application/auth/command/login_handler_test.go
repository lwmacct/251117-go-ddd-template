package command

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domainAuth "github.com/lwmacct/251117-go-ddd-template/internal/domain/auth"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/role"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/twofa"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
)

// mockUserQueryRepo 是用户查询仓储的 Mock 实现。
type mockUserQueryRepo struct {
	users           map[string]*user.User // key: username or email
	getUserByIDFunc func(ctx context.Context, id uint) (*user.User, error)
}

func newMockUserQueryRepo() *mockUserQueryRepo {
	return &mockUserQueryRepo{
		users: make(map[string]*user.User),
	}
}

func (m *mockUserQueryRepo) GetByID(ctx context.Context, id uint) (*user.User, error) {
	if m.getUserByIDFunc != nil {
		return m.getUserByIDFunc(ctx, id)
	}
	return nil, user.ErrUserNotFound
}

func (m *mockUserQueryRepo) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	if u, ok := m.users[username]; ok {
		return u, nil
	}
	return nil, user.ErrUserNotFound
}

func (m *mockUserQueryRepo) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	if u, ok := m.users[email]; ok {
		return u, nil
	}
	return nil, user.ErrUserNotFound
}

func (m *mockUserQueryRepo) GetByIDWithRoles(ctx context.Context, id uint) (*user.User, error) {
	return m.GetByID(ctx, id)
}

func (m *mockUserQueryRepo) GetByUsernameWithRoles(ctx context.Context, username string) (*user.User, error) {
	return m.GetByUsername(ctx, username)
}

func (m *mockUserQueryRepo) GetByEmailWithRoles(ctx context.Context, email string) (*user.User, error) {
	return m.GetByEmail(ctx, email)
}

func (m *mockUserQueryRepo) List(ctx context.Context, offset, limit int) ([]*user.User, error) {
	return nil, nil
}

func (m *mockUserQueryRepo) Count(ctx context.Context) (int64, error) {
	return 0, nil
}

func (m *mockUserQueryRepo) GetRoles(ctx context.Context, userID uint) ([]uint, error) {
	return nil, nil
}

func (m *mockUserQueryRepo) Search(ctx context.Context, keyword string, offset, limit int) ([]*user.User, error) {
	return nil, nil
}

func (m *mockUserQueryRepo) CountBySearch(ctx context.Context, keyword string) (int64, error) {
	return 0, nil
}

func (m *mockUserQueryRepo) Exists(ctx context.Context, id uint) (bool, error) {
	return false, nil
}

func (m *mockUserQueryRepo) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	_, ok := m.users[username]
	return ok, nil
}

func (m *mockUserQueryRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	_, ok := m.users[email]
	return ok, nil
}

func (m *mockUserQueryRepo) GetUserIDsByRole(ctx context.Context, roleID uint) ([]uint, error) {
	return nil, nil
}

func (m *mockUserQueryRepo) addUser(u *user.User) {
	m.users[u.Username] = u
	m.users[u.Email] = u
}

// mockCaptchaCommandRepo 是验证码命令仓储的 Mock 实现。
type mockCaptchaCommandRepo struct {
	validCaptchas map[string]string // captchaID -> code
	verifyError   error
}

func newMockCaptchaCommandRepo() *mockCaptchaCommandRepo {
	return &mockCaptchaCommandRepo{
		validCaptchas: make(map[string]string),
	}
}

func (m *mockCaptchaCommandRepo) Create(ctx context.Context, captchaID string, code string, expiration time.Duration) error {
	m.validCaptchas[captchaID] = code
	return nil
}

func (m *mockCaptchaCommandRepo) Verify(ctx context.Context, captchaID string, code string) (bool, error) {
	if m.verifyError != nil {
		return false, m.verifyError
	}
	if stored, ok := m.validCaptchas[captchaID]; ok {
		return stored == code, nil
	}
	return false, nil
}

func (m *mockCaptchaCommandRepo) Delete(ctx context.Context, captchaID string) error {
	delete(m.validCaptchas, captchaID)
	return nil
}

// mockTwofaQueryRepo 是 2FA 查询仓储的 Mock 实现。
type mockTwofaQueryRepo struct {
	configs map[uint]*twofa.TwoFA // userID -> TwoFA config
}

func newMockTwofaQueryRepo() *mockTwofaQueryRepo {
	return &mockTwofaQueryRepo{
		configs: make(map[uint]*twofa.TwoFA),
	}
}

func (m *mockTwofaQueryRepo) FindByUserID(ctx context.Context, userID uint) (*twofa.TwoFA, error) {
	if config, ok := m.configs[userID]; ok {
		return config, nil
	}
	return nil, errors.New("2fa not found")
}

func (m *mockTwofaQueryRepo) IsEnabled(ctx context.Context, userID uint) (bool, error) {
	if config, ok := m.configs[userID]; ok {
		return config.Enabled, nil
	}
	return false, nil
}

// mockAuthService 是认证服务的 Mock 实现。
type mockAuthService struct {
	passwordHashes map[string]string // plainPassword -> hashedPassword
	verifyError    error
	tokenCounter   int
}

func newMockAuthService() *mockAuthService {
	return &mockAuthService{
		passwordHashes: make(map[string]string),
	}
}

func (m *mockAuthService) VerifyPassword(ctx context.Context, hashedPassword, plainPassword string) error {
	if m.verifyError != nil {
		return m.verifyError
	}
	// 简化验证：检查 plain -> hashed 映射
	if expected, ok := m.passwordHashes[plainPassword]; ok && expected == hashedPassword {
		return nil
	}
	return domainAuth.ErrPasswordMismatch
}

func (m *mockAuthService) GeneratePasswordHash(ctx context.Context, password string) (string, error) {
	hash := "hashed_" + password
	m.passwordHashes[password] = hash
	return hash, nil
}

func (m *mockAuthService) ValidatePasswordPolicy(ctx context.Context, password string) error {
	if len(password) < 6 {
		return domainAuth.ErrWeakPassword
	}
	return nil
}

func (m *mockAuthService) GenerateAccessToken(ctx context.Context, userID uint, username string) (string, time.Time, error) {
	m.tokenCounter++
	token := "access_token_" + username
	expiresAt := time.Now().Add(time.Hour)
	return token, expiresAt, nil
}

func (m *mockAuthService) GenerateRefreshToken(ctx context.Context, userID uint) (string, time.Time, error) {
	m.tokenCounter++
	token := "refresh_token"
	expiresAt := time.Now().Add(24 * time.Hour)
	return token, expiresAt, nil
}

func (m *mockAuthService) ValidateAccessToken(ctx context.Context, token string) (*domainAuth.TokenClaims, error) {
	return nil, nil //nolint:nilnil // test mock returns nil for not found
}

func (m *mockAuthService) ValidateRefreshToken(ctx context.Context, token string) (uint, error) {
	return 0, nil
}

func (m *mockAuthService) GeneratePATToken(ctx context.Context) (string, error) {
	return "pat_token", nil
}

func (m *mockAuthService) HashPATToken(ctx context.Context, token string) string {
	return "hashed_" + token
}

//nolint:unparam // test helper - plain is used to construct hash
func (m *mockAuthService) setPassword(plain, hashed string) {
	m.passwordHashes[plain] = hashed
}

//nolint:unparam // test helper - id is used in test variations
func newTestUser(id uint, username, email, password, status string) *user.User {
	return &user.User{
		ID:       id,
		Username: username,
		Email:    email,
		Password: "hashed_" + password, // 与 mockAuthService 一致
		Status:   status,
		Roles: []role.Role{
			{ID: 1, Name: "user"},
		},
	}
}

type loginTestSetup struct {
	handler     *LoginHandler
	userRepo    *mockUserQueryRepo
	captchaRepo *mockCaptchaCommandRepo
	twofaRepo   *mockTwofaQueryRepo
	authService *mockAuthService
}

func setupLoginTest() *loginTestSetup {
	userRepo := newMockUserQueryRepo()
	captchaRepo := newMockCaptchaCommandRepo()
	twofaRepo := newMockTwofaQueryRepo()
	authService := newMockAuthService()

	handler := NewLoginHandler(userRepo, captchaRepo, twofaRepo, authService)

	return &loginTestSetup{
		handler:     handler,
		userRepo:    userRepo,
		captchaRepo: captchaRepo,
		twofaRepo:   twofaRepo,
		authService: authService,
	}
}

func TestLoginHandler_Handle(t *testing.T) {
	ctx := context.Background()

	t.Run("成功登录 - 用户名", func(t *testing.T) {
		setup := setupLoginTest()

		// 准备测试数据
		testUser := newTestUser(1, "testuser", "test@example.com", "password123", "active")
		setup.userRepo.addUser(testUser)
		setup.authService.setPassword("password123", "hashed_password123")
		setup.captchaRepo.validCaptchas["captcha-id"] = "1234"

		// 执行
		cmd := LoginCommand{
			Account:   "testuser",
			Password:  "password123",
			CaptchaID: "captcha-id",
			Captcha:   "1234",
		}
		result, err := setup.handler.Handle(ctx, cmd)

		// 验证
		require.NoError(t, err, "Handle() 应该成功")
		require.NotNil(t, result, "结果不应为空")
		assert.NotEmpty(t, result.AccessToken, "AccessToken 不应为空")
		assert.NotEmpty(t, result.RefreshToken, "RefreshToken 不应为空")
		assert.Equal(t, testUser.ID, result.UserID, "UserID 应该匹配")
		assert.Equal(t, testUser.Username, result.Username, "Username 应该匹配")
		assert.False(t, result.Requires2FA, "Requires2FA 应该是 false")
	})

	t.Run("成功登录 - 邮箱", func(t *testing.T) {
		setup := setupLoginTest()

		testUser := newTestUser(1, "testuser", "test@example.com", "password123", "active")
		setup.userRepo.addUser(testUser)
		setup.authService.setPassword("password123", "hashed_password123")
		setup.captchaRepo.validCaptchas["captcha-id"] = "1234"

		cmd := LoginCommand{
			Account:   "test@example.com", // 使用邮箱
			Password:  "password123",
			CaptchaID: "captcha-id",
			Captcha:   "1234",
		}
		result, err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err, "Handle() 应该成功")
		assert.Equal(t, testUser.ID, result.UserID, "UserID 应该匹配")
	})

	t.Run("验证码错误", func(t *testing.T) {
		setup := setupLoginTest()

		testUser := newTestUser(1, "testuser", "test@example.com", "password123", "active")
		setup.userRepo.addUser(testUser)
		setup.captchaRepo.validCaptchas["captcha-id"] = "1234"

		cmd := LoginCommand{
			Account:   "testuser",
			Password:  "password123",
			CaptchaID: "captcha-id",
			Captcha:   "wrong", // 错误的验证码
		}
		_, err := setup.handler.Handle(ctx, cmd)

		assert.ErrorIs(t, err, domainAuth.ErrInvalidCaptcha, "应该返回 ErrInvalidCaptcha")
	})

	t.Run("验证码验证失败 - 系统错误", func(t *testing.T) {
		setup := setupLoginTest()
		setup.captchaRepo.verifyError = errors.New("redis connection error")

		cmd := LoginCommand{
			Account:   "testuser",
			Password:  "password123",
			CaptchaID: "captcha-id",
			Captcha:   "1234",
		}
		_, err := setup.handler.Handle(ctx, cmd)

		assert.Error(t, err, "应该返回错误")
	})

	t.Run("用户不存在", func(t *testing.T) {
		setup := setupLoginTest()
		setup.captchaRepo.validCaptchas["captcha-id"] = "1234"

		cmd := LoginCommand{
			Account:   "nonexistent",
			Password:  "password123",
			CaptchaID: "captcha-id",
			Captcha:   "1234",
		}
		_, err := setup.handler.Handle(ctx, cmd)

		assert.ErrorIs(t, err, domainAuth.ErrInvalidCredentials, "应该返回 ErrInvalidCredentials")
	})

	t.Run("用户被禁用", func(t *testing.T) {
		setup := setupLoginTest()

		testUser := newTestUser(1, "banneduser", "banned@example.com", "password123", "banned")
		setup.userRepo.addUser(testUser)
		setup.captchaRepo.validCaptchas["captcha-id"] = "1234"

		cmd := LoginCommand{
			Account:   "banneduser",
			Password:  "password123",
			CaptchaID: "captcha-id",
			Captcha:   "1234",
		}
		_, err := setup.handler.Handle(ctx, cmd)

		assert.ErrorIs(t, err, domainAuth.ErrUserBanned, "应该返回 ErrUserBanned")
	})

	t.Run("用户未激活", func(t *testing.T) {
		setup := setupLoginTest()

		testUser := newTestUser(1, "inactiveuser", "inactive@example.com", "password123", "inactive")
		setup.userRepo.addUser(testUser)
		setup.captchaRepo.validCaptchas["captcha-id"] = "1234"

		cmd := LoginCommand{
			Account:   "inactiveuser",
			Password:  "password123",
			CaptchaID: "captcha-id",
			Captcha:   "1234",
		}
		_, err := setup.handler.Handle(ctx, cmd)

		assert.ErrorIs(t, err, domainAuth.ErrUserInactive, "应该返回 ErrUserInactive")
	})

	t.Run("密码错误", func(t *testing.T) {
		setup := setupLoginTest()

		testUser := newTestUser(1, "testuser", "test@example.com", "password123", "active")
		setup.userRepo.addUser(testUser)
		setup.authService.setPassword("password123", "hashed_password123")
		setup.captchaRepo.validCaptchas["captcha-id"] = "1234"

		cmd := LoginCommand{
			Account:   "testuser",
			Password:  "wrongpassword", // 错误密码
			CaptchaID: "captcha-id",
			Captcha:   "1234",
		}
		_, err := setup.handler.Handle(ctx, cmd)

		assert.ErrorIs(t, err, domainAuth.ErrInvalidCredentials, "应该返回 ErrInvalidCredentials")
	})

	t.Run("启用 2FA 的用户登录", func(t *testing.T) {
		setup := setupLoginTest()

		testUser := newTestUser(1, "2fauser", "2fa@example.com", "password123", "active")
		setup.userRepo.addUser(testUser)
		setup.authService.setPassword("password123", "hashed_password123")
		setup.captchaRepo.validCaptchas["captcha-id"] = "1234"

		// 启用 2FA
		setup.twofaRepo.configs[testUser.ID] = &twofa.TwoFA{
			ID:      1,
			UserID:  testUser.ID,
			Enabled: true,
			Secret:  "TESTSECRET",
		}

		cmd := LoginCommand{
			Account:   "2fauser",
			Password:  "password123",
			CaptchaID: "captcha-id",
			Captcha:   "1234",
		}
		result, err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err, "Handle() 应该成功")
		assert.True(t, result.Requires2FA, "Requires2FA 应该是 true")
		assert.NotEmpty(t, result.SessionToken, "SessionToken 应该返回临时会话令牌")
		assert.Empty(t, result.AccessToken, "AccessToken 不应该返回，需要完成 2FA 验证")
	})

	t.Run("2FA 未启用时正常登录", func(t *testing.T) {
		setup := setupLoginTest()

		testUser := newTestUser(1, "normaluser", "normal@example.com", "password123", "active")
		setup.userRepo.addUser(testUser)
		setup.authService.setPassword("password123", "hashed_password123")
		setup.captchaRepo.validCaptchas["captcha-id"] = "1234"

		// 2FA 配置存在但未启用
		setup.twofaRepo.configs[testUser.ID] = &twofa.TwoFA{
			ID:      1,
			UserID:  testUser.ID,
			Enabled: false,
		}

		cmd := LoginCommand{
			Account:   "normaluser",
			Password:  "password123",
			CaptchaID: "captcha-id",
			Captcha:   "1234",
		}
		result, err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err, "Handle() 应该成功")
		assert.False(t, result.Requires2FA, "Requires2FA 应该是 false (2FA 未启用)")
		assert.NotEmpty(t, result.AccessToken, "AccessToken 应该正常返回令牌")
	})
}

func TestNewLoginHandler(t *testing.T) {
	t.Run("创建 LoginHandler", func(t *testing.T) {
		userRepo := newMockUserQueryRepo()
		captchaRepo := newMockCaptchaCommandRepo()
		twofaRepo := newMockTwofaQueryRepo()
		authService := newMockAuthService()

		handler := NewLoginHandler(userRepo, captchaRepo, twofaRepo, authService)

		require.NotNil(t, handler, "NewLoginHandler() 不应返回 nil")
	})
}

func BenchmarkLoginHandler_Handle(b *testing.B) {
	ctx := context.Background()
	setup := setupLoginTest()

	testUser := newTestUser(1, "benchuser", "bench@example.com", "password123", "active")
	setup.userRepo.addUser(testUser)
	setup.authService.setPassword("password123", "hashed_password123")

	cmd := LoginCommand{
		Account:   "benchuser",
		Password:  "password123",
		CaptchaID: "captcha-id",
		Captcha:   "1234",
	}

	for b.Loop() {
		setup.captchaRepo.validCaptchas["captcha-id"] = "1234" // 每次重置验证码
		_, _ = setup.handler.Handle(ctx, cmd)
	}
}
