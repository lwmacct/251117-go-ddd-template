package command

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/menu"
)

// updateMockMenuCommandRepo 用于更新测试的命令仓储 Mock
type updateMockMenuCommandRepo struct {
	mockMenuCommandRepo

	updateError error
	updatedMenu *menu.Menu
}

func (m *updateMockMenuCommandRepo) Update(ctx context.Context, menuItem *menu.Menu) error {
	if m.updateError != nil {
		return m.updateError
	}
	m.updatedMenu = menuItem
	return nil
}

type updateMenuTestSetup struct {
	handler     *UpdateMenuHandler
	commandRepo *updateMockMenuCommandRepo
	queryRepo   *mockMenuQueryRepo
}

func setupUpdateMenuTest() *updateMenuTestSetup {
	commandRepo := &updateMockMenuCommandRepo{mockMenuCommandRepo: *newMockMenuCommandRepo()}
	queryRepo := newMockMenuQueryRepo()
	handler := NewUpdateMenuHandler(commandRepo, queryRepo)

	return &updateMenuTestSetup{
		handler:     handler,
		commandRepo: commandRepo,
		queryRepo:   queryRepo,
	}
}

func TestUpdateMenuHandler_Handle(t *testing.T) {
	ctx := context.Background()

	t.Run("成功更新菜单", func(t *testing.T) {
		setup := setupUpdateMenuTest()
		setup.queryRepo.menus[1] = &menu.Menu{
			ID:    1,
			Title: "Old Title",
			Path:  "/old",
			Icon:  "old-icon",
			Order: 1,
		}

		newTitle := "New Title"
		newPath := "/new"
		cmd := UpdateMenuCommand{
			MenuID: 1,
			Title:  &newTitle,
			Path:   &newPath,
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "New Title", result.Title)
		assert.Equal(t, "/new", result.Path)
	})

	t.Run("菜单不存在", func(t *testing.T) {
		setup := setupUpdateMenuTest()

		newTitle := "New Title"
		cmd := UpdateMenuCommand{
			MenuID: 999,
			Title:  &newTitle,
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "menu not found")
		assert.Nil(t, result)
	})

	t.Run("不能设置自己为父菜单", func(t *testing.T) {
		setup := setupUpdateMenuTest()
		setup.queryRepo.menus[1] = &menu.Menu{
			ID:    1,
			Title: "Menu",
		}

		selfID := uint(1)
		cmd := UpdateMenuCommand{
			MenuID:   1,
			ParentID: &selfID,
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot set menu as its own parent")
		assert.Nil(t, result)
	})

	t.Run("父菜单不存在", func(t *testing.T) {
		setup := setupUpdateMenuTest()
		setup.queryRepo.menus[1] = &menu.Menu{
			ID:    1,
			Title: "Menu",
		}

		nonExistentParent := uint(999)
		cmd := UpdateMenuCommand{
			MenuID:   1,
			ParentID: &nonExistentParent,
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "parent menu not found")
		assert.Nil(t, result)
	})

	t.Run("更新失败", func(t *testing.T) {
		setup := setupUpdateMenuTest()
		setup.queryRepo.menus[1] = &menu.Menu{
			ID:    1,
			Title: "Menu",
		}
		setup.commandRepo.updateError = errors.New("database error")

		newTitle := "New Title"
		cmd := UpdateMenuCommand{
			MenuID: 1,
			Title:  &newTitle,
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to update menu")
		assert.Nil(t, result)
	})

	t.Run("更新可见性", func(t *testing.T) {
		setup := setupUpdateMenuTest()
		setup.queryRepo.menus[1] = &menu.Menu{
			ID:      1,
			Title:   "Menu",
			Visible: true,
		}

		visible := false
		cmd := UpdateMenuCommand{
			MenuID:  1,
			Visible: &visible,
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.False(t, result.Visible)
	})
}

func TestNewUpdateMenuHandler(t *testing.T) {
	commandRepo := newMockMenuCommandRepo()
	queryRepo := newMockMenuQueryRepo()

	handler := NewUpdateMenuHandler(commandRepo, queryRepo)

	assert.NotNil(t, handler)
}
