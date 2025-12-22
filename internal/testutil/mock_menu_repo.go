package testutil

import (
	"context"
	"errors"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/menu"
)

// MockMenuCommandRepo 是菜单命令仓储的 Mock 实现。
type MockMenuCommandRepo struct {
	Menus       []*menu.Menu
	CreateError error
	UpdateError error
	DeleteError error
	IDCounter   uint
}

// NewMockMenuCommandRepo 创建 MockMenuCommandRepo 实例
func NewMockMenuCommandRepo() *MockMenuCommandRepo {
	return &MockMenuCommandRepo{
		Menus:     make([]*menu.Menu, 0),
		IDCounter: 0,
	}
}

func (m *MockMenuCommandRepo) Create(ctx context.Context, menuItem *menu.Menu) error {
	if m.CreateError != nil {
		return m.CreateError
	}
	m.IDCounter++
	menuItem.ID = m.IDCounter
	m.Menus = append(m.Menus, menuItem)
	return nil
}

func (m *MockMenuCommandRepo) Update(ctx context.Context, menuItem *menu.Menu) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}
	return nil
}

func (m *MockMenuCommandRepo) Delete(ctx context.Context, id uint) error {
	if m.DeleteError != nil {
		return m.DeleteError
	}
	return nil
}

func (m *MockMenuCommandRepo) UpdateOrder(ctx context.Context, menus []struct {
	ID       uint
	Order    int
	ParentID *uint
}) error {
	return nil
}

// MockMenuQueryRepo 是菜单查询仓储的 Mock 实现。
type MockMenuQueryRepo struct {
	Menus     map[uint]*menu.Menu
	FindError error
}

// NewMockMenuQueryRepo 创建 MockMenuQueryRepo 实例
func NewMockMenuQueryRepo() *MockMenuQueryRepo {
	return &MockMenuQueryRepo{
		Menus: make(map[uint]*menu.Menu),
	}
}

func (m *MockMenuQueryRepo) FindByID(ctx context.Context, id uint) (*menu.Menu, error) {
	if m.FindError != nil {
		return nil, m.FindError
	}
	if menuItem, ok := m.Menus[id]; ok {
		return menuItem, nil
	}
	return nil, errors.New("menu not found")
}

func (m *MockMenuQueryRepo) FindAll(ctx context.Context) ([]*menu.Menu, error) {
	result := make([]*menu.Menu, 0, len(m.Menus))
	for _, menuItem := range m.Menus {
		result = append(result, menuItem)
	}
	return result, nil
}

func (m *MockMenuQueryRepo) FindByParentID(ctx context.Context, parentID *uint) ([]*menu.Menu, error) {
	return nil, nil
}
