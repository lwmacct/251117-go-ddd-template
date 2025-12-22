package command

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/menu"
)

// reorderMenuCommandRepo 批量排序专用 Mock
type reorderMenuCommandRepo struct {
	updateOrderErr error
	updateCalled   bool
	updateData     []struct {
		ID       uint
		Order    int
		ParentID *uint
	}
}

func newReorderMenuCommandRepo() *reorderMenuCommandRepo {
	return &reorderMenuCommandRepo{}
}

func (m *reorderMenuCommandRepo) Create(ctx context.Context, menuItem *menu.Menu) error {
	return nil
}

func (m *reorderMenuCommandRepo) Update(ctx context.Context, menuItem *menu.Menu) error {
	return nil
}

func (m *reorderMenuCommandRepo) Delete(ctx context.Context, id uint) error {
	return nil
}

func (m *reorderMenuCommandRepo) UpdateOrder(ctx context.Context, menus []struct {
	ID       uint
	Order    int
	ParentID *uint
}) error {
	m.updateCalled = true
	m.updateData = menus
	return m.updateOrderErr
}

// reorderMenuQueryRepo 批量排序专用查询 Mock
type reorderMenuQueryRepo struct {
	existingMenus map[uint]*menu.Menu
	findError     error
}

func newReorderMenuQueryRepo() *reorderMenuQueryRepo {
	return &reorderMenuQueryRepo{
		existingMenus: make(map[uint]*menu.Menu),
	}
}

func (m *reorderMenuQueryRepo) FindByID(ctx context.Context, id uint) (*menu.Menu, error) {
	if m.findError != nil {
		return nil, m.findError
	}
	if menuItem, ok := m.existingMenus[id]; ok {
		return menuItem, nil
	}
	return nil, errors.New("menu not found")
}

func (m *reorderMenuQueryRepo) FindAll(ctx context.Context) ([]*menu.Menu, error) {
	result := make([]*menu.Menu, 0, len(m.existingMenus))
	for _, menuItem := range m.existingMenus {
		result = append(result, menuItem)
	}
	return result, nil
}

func (m *reorderMenuQueryRepo) FindByParentID(ctx context.Context, parentID *uint) ([]*menu.Menu, error) {
	return nil, nil
}

type reorderMenusTestSetup struct {
	handler     *ReorderMenusHandler
	commandRepo *reorderMenuCommandRepo
	queryRepo   *reorderMenuQueryRepo
}

func setupReorderMenusTest() *reorderMenusTestSetup {
	commandRepo := newReorderMenuCommandRepo()
	queryRepo := newReorderMenuQueryRepo()
	handler := NewReorderMenusHandler(commandRepo, queryRepo)

	return &reorderMenusTestSetup{
		handler:     handler,
		commandRepo: commandRepo,
		queryRepo:   queryRepo,
	}
}

func TestReorderMenusHandler_Handle(t *testing.T) {
	ctx := context.Background()

	t.Run("成功批量更新排序", func(t *testing.T) {
		setup := setupReorderMenusTest()
		setup.queryRepo.existingMenus[1] = &menu.Menu{ID: 1, Title: "Menu 1"}
		setup.queryRepo.existingMenus[2] = &menu.Menu{ID: 2, Title: "Menu 2"}
		setup.queryRepo.existingMenus[3] = &menu.Menu{ID: 3, Title: "Menu 3"}

		cmd := ReorderMenusCommand{
			Menus: []MenuItem{
				{ID: 1, Order: 3, ParentID: nil},
				{ID: 2, Order: 1, ParentID: nil},
				{ID: 3, Order: 2, ParentID: nil},
			},
		}

		err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		assert.True(t, setup.commandRepo.updateCalled)
		assert.Len(t, setup.commandRepo.updateData, 3)
	})

	t.Run("菜单不存在", func(t *testing.T) {
		setup := setupReorderMenusTest()
		setup.queryRepo.existingMenus[1] = &menu.Menu{ID: 1, Title: "Menu 1"}
		// 菜单 999 不存在

		cmd := ReorderMenusCommand{
			Menus: []MenuItem{
				{ID: 1, Order: 1, ParentID: nil},
				{ID: 999, Order: 2, ParentID: nil},
			},
		}

		err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "menu 999 not found")
	})

	t.Run("父菜单不存在", func(t *testing.T) {
		setup := setupReorderMenusTest()
		setup.queryRepo.existingMenus[1] = &menu.Menu{ID: 1, Title: "Menu 1"}
		// 父菜单 999 不存在

		parentID := uint(999)
		cmd := ReorderMenusCommand{
			Menus: []MenuItem{
				{ID: 1, Order: 1, ParentID: &parentID},
			},
		}

		err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "parent menu 999 not found")
	})

	t.Run("更新排序失败", func(t *testing.T) {
		setup := setupReorderMenusTest()
		setup.queryRepo.existingMenus[1] = &menu.Menu{ID: 1, Title: "Menu 1"}
		setup.commandRepo.updateOrderErr = errors.New("database error")

		cmd := ReorderMenusCommand{
			Menus: []MenuItem{
				{ID: 1, Order: 1, ParentID: nil},
			},
		}

		err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to update menu order")
	})

	t.Run("空菜单列表", func(t *testing.T) {
		setup := setupReorderMenusTest()

		cmd := ReorderMenusCommand{
			Menus: []MenuItem{},
		}

		err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		assert.True(t, setup.commandRepo.updateCalled)
		assert.Empty(t, setup.commandRepo.updateData)
	})

	t.Run("更新菜单父级", func(t *testing.T) {
		setup := setupReorderMenusTest()
		setup.queryRepo.existingMenus[1] = &menu.Menu{ID: 1, Title: "Parent"}
		setup.queryRepo.existingMenus[2] = &menu.Menu{ID: 2, Title: "Child"}

		parentID := uint(1)
		cmd := ReorderMenusCommand{
			Menus: []MenuItem{
				{ID: 2, Order: 1, ParentID: &parentID},
			},
		}

		err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		assert.True(t, setup.commandRepo.updateCalled)
		require.Len(t, setup.commandRepo.updateData, 1)
		assert.Equal(t, uint(2), setup.commandRepo.updateData[0].ID)
		assert.Equal(t, &parentID, setup.commandRepo.updateData[0].ParentID)
	})

	t.Run("父菜单ID为0时跳过验证", func(t *testing.T) {
		setup := setupReorderMenusTest()
		setup.queryRepo.existingMenus[1] = &menu.Menu{ID: 1, Title: "Menu 1"}

		// ParentID 为 0 时应该跳过父菜单验证（移动到顶级）
		zeroParentID := uint(0)
		cmd := ReorderMenusCommand{
			Menus: []MenuItem{
				{ID: 1, Order: 1, ParentID: &zeroParentID},
			},
		}

		err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		assert.True(t, setup.commandRepo.updateCalled)
	})
}

func TestNewReorderMenusHandler(t *testing.T) {
	commandRepo := newReorderMenuCommandRepo()
	queryRepo := newReorderMenuQueryRepo()

	handler := NewReorderMenusHandler(commandRepo, queryRepo)

	assert.NotNil(t, handler)
}
