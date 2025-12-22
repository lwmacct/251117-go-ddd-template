package command

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/auth"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
)

// changePwdMockUserQueryRepo 用于修改密码测试的查询仓储 Mock
type changePwdMockUserQueryRepo struct {
	mockUserQueryRepo

	user       *user.User
	getUserErr error
}

func (m *changePwdMockUserQueryRepo) GetByID(ctx context.Context, id uint) (*user.User, error) {
	if m.getUserErr != nil {
		return nil, m.getUserErr
	}
	if m.user != nil && m.user.ID == id {
		return m.user, nil
	}
	return nil, user.ErrUserNotFound
}

// changePwdMockUserCommandRepo 用于修改密码测试的命令仓储 Mock
type changePwdMockUserCommandRepo struct {
	mockUserCommandRepo

	updatePwdError   error
	updatedUserID    uint
	updatedPassword  string
	updatePwdCounter int
}

func (m *changePwdMockUserCommandRepo) UpdatePassword(ctx context.Context, userID uint, hashedPassword string) error {
	if m.updatePwdError != nil {
		return m.updatePwdError
	}
	m.updatePwdCounter++
	m.updatedUserID = userID
	m.updatedPassword = hashedPassword
	return nil
}

// changePwdMockAuthService 用于修改密码测试的认证服务 Mock
type changePwdMockAuthService struct {
	mockAuthService

	verifyPwdError error
	weakPassword   bool
	hashError      error
}

func (m *changePwdMockAuthService) VerifyPassword(ctx context.Context, hashedPassword, plainPassword string) error {
	if m.verifyPwdError != nil {
		return m.verifyPwdError
	}
	// 简单模拟：比较哈希前缀
	if hashedPassword == "hashed_"+plainPassword {
		return nil
	}
	return user.ErrInvalidPassword
}

func (m *changePwdMockAuthService) ValidatePasswordPolicy(ctx context.Context, password string) error {
	if m.weakPassword || len(password) < 6 {
		return auth.ErrWeakPassword
	}
	return nil
}

func (m *changePwdMockAuthService) GeneratePasswordHash(ctx context.Context, password string) (string, error) {
	if m.hashError != nil {
		return "", m.hashError
	}
	return "hashed_" + password, nil
}

type changePasswordTestSetup struct {
	handler     *ChangePasswordHandler
	commandRepo *changePwdMockUserCommandRepo
	queryRepo   *changePwdMockUserQueryRepo
	authService *changePwdMockAuthService
}

func setupChangePasswordTest() *changePasswordTestSetup {
	commandRepo := &changePwdMockUserCommandRepo{mockUserCommandRepo: *newMockUserCommandRepo()}
	queryRepo := &changePwdMockUserQueryRepo{mockUserQueryRepo: *newMockUserQueryRepo()}
	authService := &changePwdMockAuthService{mockAuthService: *newMockAuthService()}

	handler := NewChangePasswordHandler(commandRepo, queryRepo, authService)

	return &changePasswordTestSetup{
		handler:     handler,
		commandRepo: commandRepo,
		queryRepo:   queryRepo,
		authService: authService,
	}
}

func TestChangePasswordHandler_Handle(t *testing.T) {
	ctx := context.Background()

	t.Run("成功修改密码", func(t *testing.T) {
		setup := setupChangePasswordTest()
		setup.queryRepo.user = &user.User{
			ID:       1,
			Username: "testuser",
			Password: "hashed_oldpassword",
		}

		cmd := ChangePasswordCommand{
			UserID:      1,
			OldPassword: "oldpassword",
			NewPassword: "newpassword123",
		}

		err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		assert.Equal(t, uint(1), setup.commandRepo.updatedUserID)
		assert.Equal(t, "hashed_newpassword123", setup.commandRepo.updatedPassword)
	})

	t.Run("用户不存在", func(t *testing.T) {
		setup := setupChangePasswordTest()
		setup.queryRepo.getUserErr = user.ErrUserNotFound

		cmd := ChangePasswordCommand{
			UserID:      999,
			OldPassword: "oldpassword",
			NewPassword: "newpassword123",
		}

		err := setup.handler.Handle(ctx, cmd)

		assert.ErrorIs(t, err, user.ErrUserNotFound)
	})

	t.Run("旧密码错误", func(t *testing.T) {
		setup := setupChangePasswordTest()
		setup.queryRepo.user = &user.User{
			ID:       1,
			Username: "testuser",
			Password: "hashed_correctpassword",
		}

		cmd := ChangePasswordCommand{
			UserID:      1,
			OldPassword: "wrongpassword",
			NewPassword: "newpassword123",
		}

		err := setup.handler.Handle(ctx, cmd)

		assert.ErrorIs(t, err, user.ErrInvalidPassword)
	})

	t.Run("新密码不符合策略", func(t *testing.T) {
		setup := setupChangePasswordTest()
		setup.queryRepo.user = &user.User{
			ID:       1,
			Username: "testuser",
			Password: "hashed_oldpassword",
		}

		cmd := ChangePasswordCommand{
			UserID:      1,
			OldPassword: "oldpassword",
			NewPassword: "123", // 太短
		}

		err := setup.handler.Handle(ctx, cmd)

		assert.ErrorIs(t, err, auth.ErrWeakPassword)
	})

	t.Run("密码哈希生成失败", func(t *testing.T) {
		setup := setupChangePasswordTest()
		setup.queryRepo.user = &user.User{
			ID:       1,
			Username: "testuser",
			Password: "hashed_oldpassword",
		}
		setup.authService.hashError = errors.New("hash error")

		cmd := ChangePasswordCommand{
			UserID:      1,
			OldPassword: "oldpassword",
			NewPassword: "newpassword123",
		}

		err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to hash password")
	})

	t.Run("更新密码失败", func(t *testing.T) {
		setup := setupChangePasswordTest()
		setup.queryRepo.user = &user.User{
			ID:       1,
			Username: "testuser",
			Password: "hashed_oldpassword",
		}
		setup.commandRepo.updatePwdError = errors.New("database error")

		cmd := ChangePasswordCommand{
			UserID:      1,
			OldPassword: "oldpassword",
			NewPassword: "newpassword123",
		}

		err := setup.handler.Handle(ctx, cmd)

		assert.Error(t, err)
	})
}

func TestNewChangePasswordHandler(t *testing.T) {
	t.Run("创建 ChangePasswordHandler", func(t *testing.T) {
		commandRepo := newMockUserCommandRepo()
		queryRepo := newMockUserQueryRepo()
		authService := newMockAuthService()

		handler := NewChangePasswordHandler(commandRepo, queryRepo, authService)

		assert.NotNil(t, handler)
	})
}
