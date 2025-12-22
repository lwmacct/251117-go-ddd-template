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
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
)

// refreshTokenAuthService 刷新令牌专用 Mock
type refreshTokenAuthService struct {
	validateRefreshError error
	validateRefreshUID   uint
	generateAccessError  error
	generateRefreshError error
}

func newRefreshTokenAuthService() *refreshTokenAuthService {
	return &refreshTokenAuthService{}
}

func (m *refreshTokenAuthService) VerifyPassword(ctx context.Context, hashedPassword, plainPassword string) error {
	return nil
}

func (m *refreshTokenAuthService) GeneratePasswordHash(ctx context.Context, password string) (string, error) {
	return "hashed_" + password, nil
}

func (m *refreshTokenAuthService) ValidatePasswordPolicy(ctx context.Context, password string) error {
	return nil
}

func (m *refreshTokenAuthService) GenerateAccessToken(ctx context.Context, userID uint, username string) (string, time.Time, error) {
	if m.generateAccessError != nil {
		return "", time.Time{}, m.generateAccessError
	}
	return "new_access_token", time.Now().Add(time.Hour), nil
}

func (m *refreshTokenAuthService) GenerateRefreshToken(ctx context.Context, userID uint) (string, time.Time, error) {
	if m.generateRefreshError != nil {
		return "", time.Time{}, m.generateRefreshError
	}
	return "new_refresh_token", time.Now().Add(24 * time.Hour), nil
}

func (m *refreshTokenAuthService) ValidateAccessToken(ctx context.Context, token string) (*domainAuth.TokenClaims, error) {
	return nil, nil //nolint:nilnil // test mock returns nil for not found
}

func (m *refreshTokenAuthService) ValidateRefreshToken(ctx context.Context, token string) (uint, error) {
	if m.validateRefreshError != nil {
		return 0, m.validateRefreshError
	}
	return m.validateRefreshUID, nil
}

func (m *refreshTokenAuthService) GeneratePATToken(ctx context.Context) (string, error) {
	return "pat_token", nil
}

func (m *refreshTokenAuthService) HashPATToken(ctx context.Context, token string) string {
	return "hashed_" + token
}

// refreshTokenUserQueryRepo 刷新令牌专用查询 Mock
type refreshTokenUserQueryRepo struct {
	users map[uint]*user.User
}

func newRefreshTokenUserQueryRepo() *refreshTokenUserQueryRepo {
	return &refreshTokenUserQueryRepo{
		users: make(map[uint]*user.User),
	}
}

func (m *refreshTokenUserQueryRepo) GetByID(ctx context.Context, id uint) (*user.User, error) {
	return nil, user.ErrUserNotFound
}

func (m *refreshTokenUserQueryRepo) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	return nil, user.ErrUserNotFound
}

func (m *refreshTokenUserQueryRepo) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	return nil, user.ErrUserNotFound
}

func (m *refreshTokenUserQueryRepo) GetByIDWithRoles(ctx context.Context, id uint) (*user.User, error) {
	if u, ok := m.users[id]; ok {
		return u, nil
	}
	return nil, user.ErrUserNotFound
}

func (m *refreshTokenUserQueryRepo) GetByUsernameWithRoles(ctx context.Context, username string) (*user.User, error) {
	return nil, user.ErrUserNotFound
}

func (m *refreshTokenUserQueryRepo) GetByEmailWithRoles(ctx context.Context, email string) (*user.User, error) {
	return nil, user.ErrUserNotFound
}

func (m *refreshTokenUserQueryRepo) List(ctx context.Context, offset, limit int) ([]*user.User, error) {
	return nil, nil
}

func (m *refreshTokenUserQueryRepo) Count(ctx context.Context) (int64, error) {
	return 0, nil
}

func (m *refreshTokenUserQueryRepo) GetRoles(ctx context.Context, userID uint) ([]uint, error) {
	return nil, nil
}

func (m *refreshTokenUserQueryRepo) Search(ctx context.Context, keyword string, offset, limit int) ([]*user.User, error) {
	return nil, nil
}

func (m *refreshTokenUserQueryRepo) CountBySearch(ctx context.Context, keyword string) (int64, error) {
	return 0, nil
}

func (m *refreshTokenUserQueryRepo) Exists(ctx context.Context, id uint) (bool, error) {
	return false, nil
}

func (m *refreshTokenUserQueryRepo) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	return false, nil
}

func (m *refreshTokenUserQueryRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return false, nil
}

func (m *refreshTokenUserQueryRepo) GetUserIDsByRole(ctx context.Context, roleID uint) ([]uint, error) {
	return nil, nil
}

type refreshTokenTestSetup struct {
	handler     *RefreshTokenHandler
	userRepo    *refreshTokenUserQueryRepo
	authService *refreshTokenAuthService
}

func setupRefreshTokenTest() *refreshTokenTestSetup {
	userRepo := newRefreshTokenUserQueryRepo()
	authService := newRefreshTokenAuthService()

	handler := NewRefreshTokenHandler(userRepo, authService)

	return &refreshTokenTestSetup{
		handler:     handler,
		userRepo:    userRepo,
		authService: authService,
	}
}

func TestRefreshTokenHandler_Handle(t *testing.T) {
	ctx := context.Background()

	t.Run("成功刷新令牌", func(t *testing.T) {
		setup := setupRefreshTokenTest()
		setup.authService.validateRefreshUID = 1
		setup.userRepo.users[1] = &user.User{
			ID:       1,
			Username: "testuser",
			Email:    "test@example.com",
			Status:   "active",
			Roles: []role.Role{
				{ID: 1, Name: "user"},
			},
		}

		cmd := RefreshTokenCommand{
			RefreshToken: "valid_refresh_token",
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "new_access_token", result.AccessToken)
		assert.Equal(t, "new_refresh_token", result.RefreshToken)
		assert.Equal(t, "Bearer", result.TokenType)
	})

	t.Run("无效的刷新令牌", func(t *testing.T) {
		setup := setupRefreshTokenTest()
		setup.authService.validateRefreshError = domainAuth.ErrInvalidToken

		cmd := RefreshTokenCommand{
			RefreshToken: "invalid_token",
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		require.ErrorIs(t, err, domainAuth.ErrInvalidToken)
		assert.Nil(t, result)
	})

	t.Run("用户不存在", func(t *testing.T) {
		setup := setupRefreshTokenTest()
		setup.authService.validateRefreshUID = 999 // 用户 999 不存在

		cmd := RefreshTokenCommand{
			RefreshToken: "valid_token_but_user_gone",
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		require.ErrorIs(t, err, domainAuth.ErrUserNotFound)
		assert.Nil(t, result)
	})

	t.Run("用户被禁用", func(t *testing.T) {
		setup := setupRefreshTokenTest()
		setup.authService.validateRefreshUID = 1
		setup.userRepo.users[1] = &user.User{
			ID:       1,
			Username: "banneduser",
			Status:   "banned",
		}

		cmd := RefreshTokenCommand{
			RefreshToken: "valid_token",
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		require.ErrorIs(t, err, domainAuth.ErrUserBanned)
		assert.Nil(t, result)
	})

	t.Run("用户未激活", func(t *testing.T) {
		setup := setupRefreshTokenTest()
		setup.authService.validateRefreshUID = 1
		setup.userRepo.users[1] = &user.User{
			ID:       1,
			Username: "inactiveuser",
			Status:   "inactive",
		}

		cmd := RefreshTokenCommand{
			RefreshToken: "valid_token",
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		require.ErrorIs(t, err, domainAuth.ErrUserInactive)
		assert.Nil(t, result)
	})

	t.Run("生成访问令牌失败", func(t *testing.T) {
		setup := setupRefreshTokenTest()
		setup.authService.validateRefreshUID = 1
		setup.authService.generateAccessError = errors.New("token generation failed")
		setup.userRepo.users[1] = &user.User{
			ID:       1,
			Username: "testuser",
			Status:   "active",
		}

		cmd := RefreshTokenCommand{
			RefreshToken: "valid_token",
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to generate access token")
		assert.Nil(t, result)
	})

	t.Run("生成刷新令牌失败", func(t *testing.T) {
		setup := setupRefreshTokenTest()
		setup.authService.validateRefreshUID = 1
		setup.authService.generateRefreshError = errors.New("refresh token generation failed")
		setup.userRepo.users[1] = &user.User{
			ID:       1,
			Username: "testuser",
			Status:   "active",
		}

		cmd := RefreshTokenCommand{
			RefreshToken: "valid_token",
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to generate refresh token")
		assert.Nil(t, result)
	})
}

func TestNewRefreshTokenHandler(t *testing.T) {
	userRepo := newRefreshTokenUserQueryRepo()
	authService := newRefreshTokenAuthService()

	handler := NewRefreshTokenHandler(userRepo, authService)

	assert.NotNil(t, handler)
}
