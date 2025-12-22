package query

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/role"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
)

// mockUserQueryRepo 是用户查询仓储的 Mock 实现
type mockUserQueryRepo struct {
	users           map[uint]*user.User
	usersWithRoles  map[uint]*user.User
	getByIDErr      error
	getWithRolesErr error
}

func newMockUserQueryRepo() *mockUserQueryRepo {
	return &mockUserQueryRepo{
		users:          make(map[uint]*user.User),
		usersWithRoles: make(map[uint]*user.User),
	}
}

func (m *mockUserQueryRepo) GetByID(ctx context.Context, id uint) (*user.User, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}
	if u, ok := m.users[id]; ok {
		return u, nil
	}
	return nil, user.ErrUserNotFound
}

func (m *mockUserQueryRepo) GetByIDWithRoles(ctx context.Context, id uint) (*user.User, error) {
	if m.getWithRolesErr != nil {
		return nil, m.getWithRolesErr
	}
	if u, ok := m.usersWithRoles[id]; ok {
		return u, nil
	}
	return nil, user.ErrUserNotFound
}

// 空实现其他接口方法
func (m *mockUserQueryRepo) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	return nil, nil //nolint:nilnil // test mock returns nil for not found
}
func (m *mockUserQueryRepo) GetByEmail(ctx context.Context, email string) (*user.User, error) {
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
func (m *mockUserQueryRepo) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	return false, nil
}
func (m *mockUserQueryRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return false, nil
}
func (m *mockUserQueryRepo) GetUserIDsByRole(ctx context.Context, roleID uint) ([]uint, error) {
	return nil, nil
}

type getUserTestSetup struct {
	handler   *GetUserHandler
	queryRepo *mockUserQueryRepo
}

func setupGetUserTest() *getUserTestSetup {
	queryRepo := newMockUserQueryRepo()
	handler := NewGetUserHandler(queryRepo)

	return &getUserTestSetup{
		handler:   handler,
		queryRepo: queryRepo,
	}
}

func TestGetUserHandler_Handle(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	t.Run("成功获取用户 - 不含角色", func(t *testing.T) {
		setup := setupGetUserTest()
		setup.queryRepo.users[1] = &user.User{
			ID:        1,
			Username:  "testuser",
			Email:     "test@example.com",
			FullName:  "Test User",
			Avatar:    "avatar.jpg",
			Bio:       "Hello world",
			Status:    "active",
			CreatedAt: now,
			UpdatedAt: now,
		}

		query := GetUserQuery{
			UserID:    1,
			WithRoles: false,
		}

		result, err := setup.handler.Handle(ctx, query)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, uint(1), result.ID)
		assert.Equal(t, "testuser", result.Username)
		assert.Equal(t, "test@example.com", result.Email)
		assert.Equal(t, "Test User", result.FullName)
		assert.Empty(t, result.Roles) // 不含角色
	})

	t.Run("成功获取用户 - 含角色", func(t *testing.T) {
		setup := setupGetUserTest()
		setup.queryRepo.usersWithRoles[1] = &user.User{
			ID:        1,
			Username:  "adminuser",
			Email:     "admin@example.com",
			FullName:  "Admin User",
			Status:    "active",
			CreatedAt: now,
			UpdatedAt: now,
			Roles: []role.Role{
				{ID: 1, Name: "admin", DisplayName: "管理员", Description: "系统管理员"},
				{ID: 2, Name: "user", DisplayName: "普通用户", Description: "普通用户"},
			},
		}

		query := GetUserQuery{
			UserID:    1,
			WithRoles: true,
		}

		result, err := setup.handler.Handle(ctx, query)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, uint(1), result.ID)
		assert.Equal(t, "adminuser", result.Username)
		require.Len(t, result.Roles, 2)
		assert.Equal(t, "admin", result.Roles[0].Name)
		assert.Equal(t, "管理员", result.Roles[0].DisplayName)
		assert.Equal(t, "user", result.Roles[1].Name)
	})

	t.Run("用户不存在 - 不含角色", func(t *testing.T) {
		setup := setupGetUserTest()

		query := GetUserQuery{
			UserID:    999,
			WithRoles: false,
		}

		result, err := setup.handler.Handle(ctx, query)

		require.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("用户不存在 - 含角色", func(t *testing.T) {
		setup := setupGetUserTest()

		query := GetUserQuery{
			UserID:    999,
			WithRoles: true,
		}

		result, err := setup.handler.Handle(ctx, query)

		require.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("查询出错", func(t *testing.T) {
		setup := setupGetUserTest()
		setup.queryRepo.getByIDErr = errors.New("database error")

		query := GetUserQuery{
			UserID:    1,
			WithRoles: false,
		}

		result, err := setup.handler.Handle(ctx, query)

		require.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("用户没有角色", func(t *testing.T) {
		setup := setupGetUserTest()
		setup.queryRepo.usersWithRoles[1] = &user.User{
			ID:        1,
			Username:  "newuser",
			Email:     "new@example.com",
			Status:    "active",
			CreatedAt: now,
			UpdatedAt: now,
			Roles:     []role.Role{}, // 空角色列表
		}

		query := GetUserQuery{
			UserID:    1,
			WithRoles: true,
		}

		result, err := setup.handler.Handle(ctx, query)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Empty(t, result.Roles)
	})
}

func TestNewGetUserHandler(t *testing.T) {
	t.Run("创建 GetUserHandler", func(t *testing.T) {
		queryRepo := newMockUserQueryRepo()
		handler := NewGetUserHandler(queryRepo)

		assert.NotNil(t, handler)
	})
}
