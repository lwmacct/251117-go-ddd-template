package manualtest

import (
	"fmt"
	"testing"
	"time"

	"github.com/lwmacct/251117-go-ddd-template/internal/application/menu"
	"github.com/lwmacct/251117-go-ddd-template/internal/manualtest/helper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 使用 time 包生成唯一标识符
var _ = time.Now

// TestMenusFlow 菜单管理完整流程测试。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestMenusFlow ./internal/manualtest/
func TestMenusFlow(t *testing.T) {
	c := helper.LoginAsAdmin(t)

	// 用于清理的变量
	var createdMenuID, childMenuID uint

	// 确保测试结束时清理资源
	t.Cleanup(func() {
		if childMenuID > 0 {
			_ = c.Delete(fmt.Sprintf("/api/admin/menus/%d", childMenuID))
		}
		if createdMenuID > 0 {
			_ = c.Delete(fmt.Sprintf("/api/admin/menus/%d", createdMenuID))
		}
	})

	// 测试 1: 获取菜单列表
	t.Log("\n测试 1: 获取菜单列表")
	menus, _, err := helper.GetList[menu.MenuDTO](c, "/api/admin/menus", nil)
	require.NoError(t, err, "获取菜单列表失败")
	initialCount := len(menus)
	t.Logf("  菜单数量: %d", initialCount)

	// 测试 2: 创建菜单
	t.Log("\n测试 2: 创建菜单")
	menuName := fmt.Sprintf("testmenu_%d", time.Now().UnixNano())
	menuPath := "/test/" + menuName
	visible := true
	createReq := menu.CreateDTO{
		Title:   menuName,
		Path:    menuPath,
		Icon:    "test-icon",
		Order:   99,
		Visible: &visible,
	}
	t.Logf("  创建菜单: %s", menuName)

	created, err := helper.Post[menu.CreateResultDTO](c, "/api/admin/menus", createReq)
	require.NoError(t, err, "创建菜单失败")
	createdMenuID = created.ID

	// 验证创建结果（创建 API 只返回 ID）
	require.NotZero(t, created.ID, "创建菜单后 ID 不应为 0")
	t.Logf("  创建成功! 菜单 ID: %d", created.ID)

	// 测试 3: 获取菜单详情并验证字段
	t.Log("\n测试 3: 获取菜单详情并验证字段")
	detail, err := helper.Get[menu.MenuDTO](c, fmt.Sprintf("/api/admin/menus/%d", created.ID), nil)
	require.NoError(t, err, "获取菜单详情失败")
	assert.Equal(t, created.ID, detail.ID, "详情 ID 不匹配")
	assert.Equal(t, menuName, detail.Title, "标题不匹配")
	assert.Equal(t, menuPath, detail.Path, "路径不匹配")
	assert.Equal(t, "test-icon", detail.Icon, "图标不匹配")
	assert.Equal(t, 99, detail.Order, "排序不匹配")
	assert.True(t, detail.Visible, "可见性应为 true")
	t.Logf("  标题: %s, 路径: %s", detail.Title, detail.Path)
	t.Logf("  图标: %s, 可见: %v", detail.Icon, detail.Visible)

	// 测试 4: 更新菜单
	t.Log("\n测试 4: 更新菜单")
	newTitle := menuName + "_updated"
	newOrder := 100
	updateReq := menu.UpdateDTO{
		Title: &newTitle,
		Order: &newOrder,
	}
	updated, err := helper.Put[menu.MenuDTO](c, fmt.Sprintf("/api/admin/menus/%d", created.ID), updateReq)
	require.NoError(t, err, "更新菜单失败")
	assert.Equal(t, newTitle, updated.Title, "更新后标题不匹配")
	assert.Equal(t, newOrder, updated.Order, "更新后排序不匹配")
	// 验证未更新的字段保持不变
	assert.Equal(t, menuPath, updated.Path, "未更新的路径不应改变")
	t.Logf("  更新成功! 新标题: %s", updated.Title)

	// 测试 5: 创建子菜单
	t.Log("\n测试 5: 创建子菜单")
	childReq := menu.CreateDTO{
		Title:    menuName + "_child",
		Path:     "/test/" + menuName + "/child",
		Icon:     "child-icon",
		ParentID: &created.ID,
		Order:    1,
		Visible:  &visible,
	}
	childResult, err := helper.Post[menu.CreateResultDTO](c, "/api/admin/menus", childReq)
	require.NoError(t, err, "创建子菜单失败")
	childMenuID = childResult.ID
	require.NotZero(t, childResult.ID, "子菜单 ID 不应为 0")
	t.Logf("  子菜单创建成功! ID: %d", childResult.ID)

	// 测试 6: 验证列表数量增加
	// 注意：菜单列表返回扁平结构，但只返回顶层菜单（parent_id 为 null 的）
	// 子菜单需要通过父菜单 ID 查询，所以顶层数量只增加 1
	t.Log("\n测试 6: 验证菜单列表数量")
	menusAfter, _, err := helper.GetList[menu.MenuDTO](c, "/api/admin/menus", nil)
	require.NoError(t, err, "获取菜单列表失败")
	assert.GreaterOrEqual(t, len(menusAfter), initialCount+1, "菜单数量应至少增加 1")
	t.Logf("  菜单数量: %d (增加了 %d)", len(menusAfter), len(menusAfter)-initialCount)

	// 测试 7: 删除子菜单
	t.Log("\n测试 7: 删除子菜单")
	err = c.Delete(fmt.Sprintf("/api/admin/menus/%d", childResult.ID))
	require.NoError(t, err, "删除子菜单失败")
	childMenuID = 0 // 已删除，清理时不需要再删
	t.Log("  子菜单删除成功!")

	// 测试 8: 删除父菜单
	t.Log("\n测试 8: 删除父菜单")
	err = c.Delete(fmt.Sprintf("/api/admin/menus/%d", created.ID))
	require.NoError(t, err, "删除菜单失败")
	createdMenuID = 0 // 已删除，清理时不需要再删
	t.Log("  菜单删除成功!")

	// 测试 9: 验证删除后获取应失败
	t.Log("\n测试 9: 验证删除后获取应返回 404")
	_, err = helper.Get[menu.MenuDTO](c, fmt.Sprintf("/api/admin/menus/%d", created.ID), nil)
	require.Error(t, err, "删除后获取应返回错误")

	t.Log("\n菜单管理流程测试完成!")
}

// TestListMenus 测试获取菜单列表。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestListMenus ./internal/manualtest/
func TestListMenus(t *testing.T) {
	c := helper.LoginAsAdmin(t)

	menus, meta, err := helper.GetList[menu.MenuDTO](c, "/api/admin/menus", nil)
	require.NoError(t, err, "获取菜单列表失败")

	t.Logf("菜单数量: %d", len(menus))
	if meta != nil {
		t.Logf("总数: %d, 总页数: %d", meta.Total, meta.TotalPages)
	}

	for _, m := range menus {
		parentInfo := "根菜单"
		if m.ParentID != nil {
			parentInfo = fmt.Sprintf("父菜单ID: %d", *m.ParentID)
		}
		visibleStr := "可见"
		if !m.Visible {
			visibleStr = "隐藏"
		}
		t.Logf("  - [%d] %s (%s) [%s] [%s]", m.ID, m.Title, m.Path, parentInfo, visibleStr)
	}
}

// TestMenuReorder 测试菜单重排序。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestMenuReorder ./internal/manualtest/
func TestMenuReorder(t *testing.T) {
	c := helper.LoginAsAdmin(t)

	// 用于清理的变量
	var menu1ID, menu2ID, menu3ID uint
	timestamp := time.Now().UnixNano()

	// 确保测试结束时清理资源
	t.Cleanup(func() {
		for _, id := range []uint{menu1ID, menu2ID, menu3ID} {
			if id > 0 {
				_ = c.Delete(fmt.Sprintf("/api/admin/menus/%d", id))
			}
		}
	})

	// 步骤 1: 创建三个测试菜单
	t.Log("\n步骤 1: 创建三个测试菜单")
	visible := true
	for i, order := range []int{1, 2, 3} {
		menuName := fmt.Sprintf("reorder_test_%d_%d", timestamp, i+1)
		createReq := menu.CreateDTO{
			Title:   menuName,
			Path:    "/test/" + menuName,
			Icon:    "test-icon",
			Order:   order,
			Visible: &visible,
		}
		created, createErr := helper.Post[menu.CreateResultDTO](c, "/api/admin/menus", createReq)
		require.NoError(t, createErr, "创建菜单失败")
		switch i {
		case 0:
			menu1ID = created.ID
		case 1:
			menu2ID = created.ID
		case 2:
			menu3ID = created.ID
		}
		t.Logf("  创建菜单 [%d] order=%d", created.ID, order)
	}

	// 步骤 2: 调用重排序接口交换顺序 (3, 1, 2 -> 将菜单3移到第一位)
	t.Log("\n步骤 2: 重排序菜单 (原: 1,2,3 -> 新: 3,1,2)")
	reorderReq := map[string]any{
		"menus": []map[string]any{
			{"id": menu3ID, "order": 1, "parent_id": nil},
			{"id": menu1ID, "order": 2, "parent_id": nil},
			{"id": menu2ID, "order": 3, "parent_id": nil},
		},
	}

	resp, err := c.R().
		SetBody(reorderReq).
		Post("/api/admin/menus/reorder")
	require.NoError(t, err, "重排序请求失败")
	// 重排序 API 返回 204 No Content
	require.False(t, resp.StatusCode() != 204 && resp.IsError(), "重排序失败，状态码: %d, 响应: %s", resp.StatusCode(), resp.String())
	t.Log("  重排序成功!")

	// 步骤 3: 验证顺序已更新
	t.Log("\n步骤 3: 验证顺序已更新")
	detail1, err := helper.Get[menu.MenuDTO](c, fmt.Sprintf("/api/admin/menus/%d", menu1ID), nil)
	require.NoError(t, err, "获取菜单 %d 详情失败", menu1ID)
	detail2, err := helper.Get[menu.MenuDTO](c, fmt.Sprintf("/api/admin/menus/%d", menu2ID), nil)
	require.NoError(t, err, "获取菜单 %d 详情失败", menu2ID)
	detail3, err := helper.Get[menu.MenuDTO](c, fmt.Sprintf("/api/admin/menus/%d", menu3ID), nil)
	require.NoError(t, err, "获取菜单 %d 详情失败", menu3ID)

	t.Logf("  菜单 %d: order=%d (期望 2)", menu1ID, detail1.Order)
	t.Logf("  菜单 %d: order=%d (期望 3)", menu2ID, detail2.Order)
	t.Logf("  菜单 %d: order=%d (期望 1)", menu3ID, detail3.Order)

	assert.Equal(t, 2, detail1.Order, "菜单 %d 顺序错误", menu1ID)
	assert.Equal(t, 3, detail2.Order, "菜单 %d 顺序错误", menu2ID)
	assert.Equal(t, 1, detail3.Order, "菜单 %d 顺序错误", menu3ID)

	t.Log("\n菜单重排序测试完成!")
}
