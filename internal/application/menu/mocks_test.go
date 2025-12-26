//nolint:forcetypeassert,gosec // Mock 返回值类型在测试中总是已知的；测试数据不会溢出
package menu

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"

	domainMenu "github.com/lwmacct/251117-go-ddd-template/internal/domain/menu"
)

// ============================================================
// MockMenuCommandRepository
// ============================================================

type MockMenuCommandRepository struct {
	mock.Mock
}

func (m *MockMenuCommandRepository) Create(ctx context.Context, menu *domainMenu.Menu) error {
	args := m.Called(ctx, menu)
	// 模拟数据库分配 ID
	if menu.ID == 0 {
		menu.ID = 1
	}
	return args.Error(0)
}

func (m *MockMenuCommandRepository) Update(ctx context.Context, menu *domainMenu.Menu) error {
	args := m.Called(ctx, menu)
	return args.Error(0)
}

func (m *MockMenuCommandRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockMenuCommandRepository) UpdateOrder(ctx context.Context, menus []struct {
	ID       uint
	Order    int
	ParentID *uint
}) error {
	args := m.Called(ctx, menus)
	return args.Error(0)
}

// ============================================================
// MockMenuQueryRepository
// ============================================================

type MockMenuQueryRepository struct {
	mock.Mock
}

func (m *MockMenuQueryRepository) FindByID(ctx context.Context, id uint) (*domainMenu.Menu, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainMenu.Menu), args.Error(1)
}

func (m *MockMenuQueryRepository) FindAll(ctx context.Context) ([]*domainMenu.Menu, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domainMenu.Menu), args.Error(1)
}

func (m *MockMenuQueryRepository) FindByParentID(ctx context.Context, parentID *uint) ([]*domainMenu.Menu, error) {
	args := m.Called(ctx, parentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domainMenu.Menu), args.Error(1)
}

// ============================================================
// 测试辅助函数
// ============================================================

func newTestMenu(id uint, title string) *domainMenu.Menu {
	return &domainMenu.Menu{
		ID:        id,
		Title:     title,
		Path:      "/path/" + title,
		Icon:      "icon-" + title,
		ParentID:  nil,
		Order:     int(id),
		Visible:   true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func newTestMenuWithParent(id uint, title string, parentID uint) *domainMenu.Menu {
	m := newTestMenu(id, title)
	m.ParentID = &parentID
	return m
}

func ptrUint(v uint) *uint {
	return &v
}

func ptrString(v string) *string {
	return &v
}

func ptrInt(v int) *int {
	return &v
}

func ptrBool(v bool) *bool {
	return &v
}
