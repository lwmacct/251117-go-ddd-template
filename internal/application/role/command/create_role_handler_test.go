package command

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/role"
)

// mockRoleCommandRepo 是角色命令仓储的 Mock 实现。
type mockRoleCommandRepo struct {
	roles       []*role.Role
	createError error
	idCounter   uint
}

func newMockRoleCommandRepo() *mockRoleCommandRepo {
	return &mockRoleCommandRepo{
		roles:     make([]*role.Role, 0),
		idCounter: 0,
	}
}

func (m *mockRoleCommandRepo) Create(ctx context.Context, r *role.Role) error {
	if m.createError != nil {
		return m.createError
	}
	m.idCounter++
	r.ID = m.idCounter
	m.roles = append(m.roles, r)
	return nil
}

func (m *mockRoleCommandRepo) Update(ctx context.Context, r *role.Role) error {
	return nil
}

func (m *mockRoleCommandRepo) Delete(ctx context.Context, id uint) error {
	return nil
}

func (m *mockRoleCommandRepo) SetPermissions(ctx context.Context, roleID uint, permissionIDs []uint) error {
	return nil
}

func (m *mockRoleCommandRepo) AddPermissions(ctx context.Context, roleID uint, permissionIDs []uint) error {
	return nil
}

func (m *mockRoleCommandRepo) RemovePermissions(ctx context.Context, roleID uint, permissionIDs []uint) error {
	return nil
}

// mockRoleQueryRepo 是角色查询仓储的 Mock 实现。
type mockRoleQueryRepo struct {
	existingNames map[string]bool
	roles         map[uint]*role.Role
	checkError    error
}

func newMockRoleQueryRepo() *mockRoleQueryRepo {
	return &mockRoleQueryRepo{
		existingNames: make(map[string]bool),
		roles:         make(map[uint]*role.Role),
	}
}

func (m *mockRoleQueryRepo) ExistsByName(ctx context.Context, name string) (bool, error) {
	if m.checkError != nil {
		return false, m.checkError
	}
	return m.existingNames[name], nil
}

func (m *mockRoleQueryRepo) FindByID(ctx context.Context, id uint) (*role.Role, error) {
	if r, ok := m.roles[id]; ok {
		return r, nil
	}
	return nil, errors.New("role not found")
}

func (m *mockRoleQueryRepo) FindByName(ctx context.Context, name string) (*role.Role, error) {
	return nil, nil //nolint:nilnil // test mock returns nil for not found
}

func (m *mockRoleQueryRepo) FindByIDWithPermissions(ctx context.Context, id uint) (*role.Role, error) {
	return m.FindByID(ctx, id)
}

func (m *mockRoleQueryRepo) List(ctx context.Context, page, limit int) ([]role.Role, int64, error) {
	return nil, 0, nil
}

func (m *mockRoleQueryRepo) GetPermissions(ctx context.Context, roleID uint) ([]role.Permission, error) {
	return nil, nil
}

func (m *mockRoleQueryRepo) Exists(ctx context.Context, id uint) (bool, error) {
	_, ok := m.roles[id]
	return ok, nil
}

type createRoleTestSetup struct {
	handler     *CreateRoleHandler
	commandRepo *mockRoleCommandRepo
	queryRepo   *mockRoleQueryRepo
}

func setupCreateRoleTest() *createRoleTestSetup {
	commandRepo := newMockRoleCommandRepo()
	queryRepo := newMockRoleQueryRepo()
	handler := NewCreateRoleHandler(commandRepo, queryRepo)

	return &createRoleTestSetup{
		handler:     handler,
		commandRepo: commandRepo,
		queryRepo:   queryRepo,
	}
}

func TestCreateRoleHandler_Handle(t *testing.T) {
	ctx := context.Background()

	t.Run("成功创建角色", func(t *testing.T) {
		setup := setupCreateRoleTest()

		cmd := CreateRoleCommand{
			Name:        "editor",
			DisplayName: "Editor",
			Description: "Can edit content",
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		require.NotNil(t, result)

		assert.Equal(t, uint(1), result.RoleID)
		assert.Equal(t, "editor", result.Name)
		assert.Equal(t, "Editor", result.DisplayName)

		// 验证角色已创建
		require.Len(t, setup.commandRepo.roles, 1)
		assert.False(t, setup.commandRepo.roles[0].IsSystem, "用户创建的角色不应该是系统角色")
	})

	t.Run("角色名已存在", func(t *testing.T) {
		setup := setupCreateRoleTest()
		setup.queryRepo.existingNames["admin"] = true

		cmd := CreateRoleCommand{
			Name:        "admin",
			DisplayName: "Administrator",
			Description: "Full access",
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "role name already exists")
		assert.Nil(t, result)
	})

	t.Run("检查角色名存在性失败", func(t *testing.T) {
		setup := setupCreateRoleTest()
		setup.queryRepo.checkError = errors.New("database error")

		cmd := CreateRoleCommand{
			Name:        "editor",
			DisplayName: "Editor",
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to check role name existence")
		assert.Nil(t, result)
	})

	t.Run("创建角色失败", func(t *testing.T) {
		setup := setupCreateRoleTest()
		setup.commandRepo.createError = errors.New("database error")

		cmd := CreateRoleCommand{
			Name:        "editor",
			DisplayName: "Editor",
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create role")
		assert.Nil(t, result)
	})
}

func TestCreateRoleHandler_ValidationScenarios(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		cmd       CreateRoleCommand
		setupFunc func(*createRoleTestSetup)
		wantErr   bool
		errMsg    string
	}{
		{
			name: "有效输入",
			cmd: CreateRoleCommand{
				Name:        "moderator",
				DisplayName: "Moderator",
				Description: "Can moderate content",
			},
			wantErr: false,
		},
		{
			name: "角色名冲突",
			cmd: CreateRoleCommand{
				Name:        "existing",
				DisplayName: "Existing Role",
			},
			setupFunc: func(s *createRoleTestSetup) {
				s.queryRepo.existingNames["existing"] = true
			},
			wantErr: true,
			errMsg:  "already exists",
		},
		{
			name: "空显示名称也可创建",
			cmd: CreateRoleCommand{
				Name:        "minimal",
				DisplayName: "",
				Description: "",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup := setupCreateRoleTest()
			if tt.setupFunc != nil {
				tt.setupFunc(setup)
			}

			result, err := setup.handler.Handle(ctx, tt.cmd)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestNewCreateRoleHandler(t *testing.T) {
	commandRepo := newMockRoleCommandRepo()
	queryRepo := newMockRoleQueryRepo()

	handler := NewCreateRoleHandler(commandRepo, queryRepo)

	assert.NotNil(t, handler)
}
