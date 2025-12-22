package command

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domainAuth "github.com/lwmacct/251117-go-ddd-template/internal/domain/auth"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/pat"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/role"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/auth"
)

// mockUserQueryRepoForPAT 是用户查询仓储的 Mock 实现
type mockUserQueryRepoForPAT struct {
	user     *user.User
	findErr  error
	existErr error
}

func newMockUserQueryRepoForPAT() *mockUserQueryRepoForPAT {
	return &mockUserQueryRepoForPAT{}
}

func (m *mockUserQueryRepoForPAT) GetByID(ctx context.Context, id uint) (*user.User, error) {
	return m.user, m.findErr
}
func (m *mockUserQueryRepoForPAT) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	return m.user, m.findErr
}
func (m *mockUserQueryRepoForPAT) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	return m.user, m.findErr
}
func (m *mockUserQueryRepoForPAT) GetByIDWithRoles(ctx context.Context, id uint) (*user.User, error) {
	return m.user, m.findErr
}
func (m *mockUserQueryRepoForPAT) GetByUsernameWithRoles(ctx context.Context, username string) (*user.User, error) {
	return m.user, m.findErr
}
func (m *mockUserQueryRepoForPAT) GetByEmailWithRoles(ctx context.Context, email string) (*user.User, error) {
	return m.user, m.findErr
}
func (m *mockUserQueryRepoForPAT) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	return m.user != nil, m.existErr
}
func (m *mockUserQueryRepoForPAT) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return m.user != nil, m.existErr
}
func (m *mockUserQueryRepoForPAT) Exists(ctx context.Context, id uint) (bool, error) {
	return m.user != nil, m.existErr
}
func (m *mockUserQueryRepoForPAT) List(ctx context.Context, offset, limit int) ([]*user.User, error) {
	return nil, nil
}
func (m *mockUserQueryRepoForPAT) Count(ctx context.Context) (int64, error) { return 0, nil }
func (m *mockUserQueryRepoForPAT) GetRoles(ctx context.Context, userID uint) ([]uint, error) {
	return nil, nil
}
func (m *mockUserQueryRepoForPAT) Search(ctx context.Context, keyword string, offset, limit int) ([]*user.User, error) {
	return nil, nil
}
func (m *mockUserQueryRepoForPAT) CountBySearch(ctx context.Context, keyword string) (int64, error) {
	return 0, nil
}
func (m *mockUserQueryRepoForPAT) GetUserIDsByRole(ctx context.Context, roleID uint) ([]uint, error) {
	return nil, nil
}

// mockPATCommandRepoForCreate PAT 命令仓储 mock
type mockPATCommandRepoForCreate struct {
	createError error
	createdPAT  *pat.PersonalAccessToken
}

func newMockPATCommandRepoForCreate() *mockPATCommandRepoForCreate {
	return &mockPATCommandRepoForCreate{}
}

func (m *mockPATCommandRepoForCreate) Create(ctx context.Context, p *pat.PersonalAccessToken) error {
	if m.createError != nil {
		return m.createError
	}
	p.ID = 1 // 模拟数据库分配 ID
	m.createdPAT = p
	return nil
}
func (m *mockPATCommandRepoForCreate) Update(ctx context.Context, p *pat.PersonalAccessToken) error {
	return nil
}
func (m *mockPATCommandRepoForCreate) Delete(ctx context.Context, id uint) error { return nil }
func (m *mockPATCommandRepoForCreate) Disable(ctx context.Context, id uint) error {
	return nil
}
func (m *mockPATCommandRepoForCreate) Enable(ctx context.Context, id uint) error { return nil }
func (m *mockPATCommandRepoForCreate) DeleteByUserID(ctx context.Context, userID uint) error {
	return nil
}
func (m *mockPATCommandRepoForCreate) CleanupExpired(ctx context.Context) error { return nil }

type createTokenTestSetup struct {
	handler        *CreateTokenHandler
	patCommandRepo *mockPATCommandRepoForCreate
	userQueryRepo  *mockUserQueryRepoForPAT
	tokenGenerator domainAuth.TokenGenerator
}

func setupCreateTokenTest() *createTokenTestSetup {
	patCommandRepo := newMockPATCommandRepoForCreate()
	userQueryRepo := newMockUserQueryRepoForPAT()
	tokenGenerator := auth.NewTokenGenerator()

	handler := NewCreateTokenHandler(patCommandRepo, userQueryRepo, tokenGenerator)

	return &createTokenTestSetup{
		handler:        handler,
		patCommandRepo: patCommandRepo,
		userQueryRepo:  userQueryRepo,
		tokenGenerator: tokenGenerator,
	}
}

// createUserWithPermissions 创建带权限的测试用户
func createUserWithPermissions(perms []string) *user.User {
	u := &user.User{
		ID:       100,
		Username: "testuser",
		Email:    "test@example.com",
		Status:   "active",
		Roles:    []role.Role{},
	}

	// 创建包含权限的角色
	permissions := make([]role.Permission, len(perms))
	for i, code := range perms {
		permissions[i] = role.Permission{
			ID:   uint(i + 1), //nolint:gosec // test data with small indices
			Code: code,
		}
	}

	u.Roles = append(u.Roles, role.Role{
		ID:          1,
		Name:        "admin",
		Permissions: permissions,
	})

	return u
}

func TestCreateTokenHandler_Handle(t *testing.T) {
	ctx := context.Background()

	t.Run("成功创建Token - 继承用户全部权限", func(t *testing.T) {
		setup := setupCreateTokenTest()
		setup.userQueryRepo.user = createUserWithPermissions([]string{"admin:users:read", "admin:users:write"})

		cmd := CreateTokenCommand{
			UserID:      100,
			Name:        "My API Token",
			ExpiresAt:   timePtr(time.Now().Add(24 * time.Hour)),
			Description: "Test token",
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotEmpty(t, result.PlainToken, "plain token should be returned")
		assert.NotNil(t, result.Token, "token entity should be returned")
		assert.Equal(t, "My API Token", result.Token.Name)
		assert.Equal(t, uint(100), result.Token.UserID)
		assert.Equal(t, "active", result.Token.Status)
		// 应该继承用户的全部权限
		assert.Len(t, result.Token.Permissions, 2)
	})

	t.Run("成功创建Token - 指定部分权限", func(t *testing.T) {
		setup := setupCreateTokenTest()
		setup.userQueryRepo.user = createUserWithPermissions([]string{"admin:users:read", "admin:users:write", "admin:roles:read"})

		cmd := CreateTokenCommand{
			UserID:      100,
			Name:        "Limited Token",
			Permissions: []string{"admin:users:read"}, // 只请求部分权限
			ExpiresAt:   timePtr(time.Now().Add(24 * time.Hour)),
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result.Token.Permissions, 1)
		assert.Contains(t, result.Token.Permissions, "admin:users:read")
	})

	t.Run("UserID为空", func(t *testing.T) {
		setup := setupCreateTokenTest()

		cmd := CreateTokenCommand{
			UserID: 0,
			Name:   "Test Token",
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "user ID is required")
		assert.Nil(t, result)
	})

	t.Run("用户不存在", func(t *testing.T) {
		setup := setupCreateTokenTest()
		setup.userQueryRepo.user = nil
		setup.userQueryRepo.findErr = errors.New("not found")

		cmd := CreateTokenCommand{
			UserID: 999,
			Name:   "Test Token",
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
		assert.Nil(t, result)
	})

	t.Run("用户没有权限", func(t *testing.T) {
		setup := setupCreateTokenTest()
		setup.userQueryRepo.user = &user.User{
			ID:       100,
			Username: "noperm",
			Status:   "active",
			Roles:    []role.Role{}, // 没有角色
		}

		cmd := CreateTokenCommand{
			UserID: 100,
			Name:   "Test Token",
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "user has no permissions")
		assert.Nil(t, result)
	})

	t.Run("请求的权限不在用户权限范围内", func(t *testing.T) {
		setup := setupCreateTokenTest()
		setup.userQueryRepo.user = createUserWithPermissions([]string{"admin:users:read"})

		cmd := CreateTokenCommand{
			UserID:      100,
			Name:        "Test Token",
			Permissions: []string{"admin:users:read", "admin:super:admin"}, // super:admin 不在用户权限中
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "is not granted to user")
		assert.Nil(t, result)
	})

	t.Run("数据库创建失败", func(t *testing.T) {
		setup := setupCreateTokenTest()
		setup.userQueryRepo.user = createUserWithPermissions([]string{"admin:users:read"})
		setup.patCommandRepo.createError = errors.New("database error")

		cmd := CreateTokenCommand{
			UserID: 100,
			Name:   "Test Token",
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create token")
		assert.Nil(t, result)
	})

	t.Run("带IP白名单的Token", func(t *testing.T) {
		setup := setupCreateTokenTest()
		setup.userQueryRepo.user = createUserWithPermissions([]string{"admin:users:read"})

		cmd := CreateTokenCommand{
			UserID:      100,
			Name:        "IP Limited Token",
			IPWhitelist: []string{"192.168.1.1", "10.0.0.0/24"},
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result.Token.IPWhitelist, 2)
		assert.Contains(t, []string(result.Token.IPWhitelist), "192.168.1.1")
		assert.Contains(t, []string(result.Token.IPWhitelist), "10.0.0.0/24")
	})
}

func TestNewCreateTokenHandler(t *testing.T) {
	patCommandRepo := newMockPATCommandRepoForCreate()
	userQueryRepo := newMockUserQueryRepoForPAT()
	tokenGenerator := auth.NewTokenGenerator()

	handler := NewCreateTokenHandler(patCommandRepo, userQueryRepo, tokenGenerator)

	assert.NotNil(t, handler)
}

func TestValidatePermissions(t *testing.T) {
	t.Run("有效的权限子集", func(t *testing.T) {
		userPerms := []string{"read", "write", "delete"}
		requested := []string{"read", "write"}

		err := validatePermissions(requested, userPerms)

		assert.NoError(t, err)
	})

	t.Run("请求全部权限", func(t *testing.T) {
		userPerms := []string{"read", "write"}
		requested := []string{"read", "write"}

		err := validatePermissions(requested, userPerms)

		assert.NoError(t, err)
	})

	t.Run("请求不存在的权限", func(t *testing.T) {
		userPerms := []string{"read", "write"}
		requested := []string{"read", "admin"}

		err := validatePermissions(requested, userPerms)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "admin")
	})

	t.Run("空权限请求", func(t *testing.T) {
		userPerms := []string{"read", "write"}
		requested := []string{}

		err := validatePermissions(requested, userPerms)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "at least one permission is required")
	})
}

// timePtr 返回 time.Time 指针
func timePtr(t time.Time) *time.Time {
	return &t
}
