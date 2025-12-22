package command

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/menu"
)

// deleteMockMenuQueryRepo 用于删除测试的查询仓储 Mock
type deleteMockMenuQueryRepo struct {
	mockMenuQueryRepo

	children        []*menu.Menu
	findChildrenErr error
}

func (m *deleteMockMenuQueryRepo) FindByParentID(ctx context.Context, parentID *uint) ([]*menu.Menu, error) {
	if m.findChildrenErr != nil {
		return nil, m.findChildrenErr
	}
	return m.children, nil
}

// deleteMockMenuCommandRepo 用于删除测试的命令仓储 Mock
type deleteMockMenuCommandRepo struct {
	mockMenuCommandRepo

	deleteError   error
	deletedMenuID uint
}

func (m *deleteMockMenuCommandRepo) Delete(ctx context.Context, id uint) error {
	if m.deleteError != nil {
		return m.deleteError
	}
	m.deletedMenuID = id
	return nil
}

type deleteMenuTestSetup struct {
	handler     *DeleteMenuHandler
	commandRepo *deleteMockMenuCommandRepo
	queryRepo   *deleteMockMenuQueryRepo
}

func setupDeleteMenuTest() *deleteMenuTestSetup {
	commandRepo := &deleteMockMenuCommandRepo{mockMenuCommandRepo: *newMockMenuCommandRepo()}
	queryRepo := &deleteMockMenuQueryRepo{mockMenuQueryRepo: *newMockMenuQueryRepo()}
	handler := NewDeleteMenuHandler(commandRepo, queryRepo)

	return &deleteMenuTestSetup{
		handler:     handler,
		commandRepo: commandRepo,
		queryRepo:   queryRepo,
	}
}

func TestDeleteMenuHandler_Handle(t *testing.T) {
	ctx := context.Background()

	t.Run("成功删除菜单", func(t *testing.T) {
		setup := setupDeleteMenuTest()
		setup.queryRepo.menus[1] = &menu.Menu{
			ID:    1,
			Title: "Menu to Delete",
		}
		setup.queryRepo.children = []*menu.Menu{} // 无子菜单

		cmd := DeleteMenuCommand{
			MenuID: 1,
		}

		err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		assert.Equal(t, uint(1), setup.commandRepo.deletedMenuID)
	})

	t.Run("菜单不存在", func(t *testing.T) {
		setup := setupDeleteMenuTest()

		cmd := DeleteMenuCommand{
			MenuID: 999,
		}

		err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "menu not found")
	})

	t.Run("有子菜单不能删除", func(t *testing.T) {
		setup := setupDeleteMenuTest()
		setup.queryRepo.menus[1] = &menu.Menu{
			ID:    1,
			Title: "Parent Menu",
		}
		setup.queryRepo.children = []*menu.Menu{
			{ID: 2, Title: "Child 1"},
			{ID: 3, Title: "Child 2"},
		}

		cmd := DeleteMenuCommand{
			MenuID: 1,
		}

		err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot delete menu with children")
	})

	t.Run("检查子菜单失败", func(t *testing.T) {
		setup := setupDeleteMenuTest()
		setup.queryRepo.menus[1] = &menu.Menu{
			ID:    1,
			Title: "Menu",
		}
		setup.queryRepo.findChildrenErr = errors.New("database error")

		cmd := DeleteMenuCommand{
			MenuID: 1,
		}

		err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to check children")
	})

	t.Run("删除操作失败", func(t *testing.T) {
		setup := setupDeleteMenuTest()
		setup.queryRepo.menus[1] = &menu.Menu{
			ID:    1,
			Title: "Menu",
		}
		setup.queryRepo.children = []*menu.Menu{}
		setup.commandRepo.deleteError = errors.New("database error")

		cmd := DeleteMenuCommand{
			MenuID: 1,
		}

		err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to delete menu")
	})
}

func TestNewDeleteMenuHandler(t *testing.T) {
	commandRepo := newMockMenuCommandRepo()
	queryRepo := newMockMenuQueryRepo()

	handler := NewDeleteMenuHandler(commandRepo, queryRepo)

	assert.NotNil(t, handler)
}
