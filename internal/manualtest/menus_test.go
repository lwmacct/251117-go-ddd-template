package manualtest

import (
	"fmt"
	"testing"
	"time"

	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/handler"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/menu"
	"github.com/lwmacct/251117-go-ddd-template/internal/manualtest/helper"
)

// 使用 time 包生成唯一标识符
var _ = time.Now

// TestMenusFlow 菜单管理完整流程测试。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestMenusFlow ./internal/manualtest/
func TestMenusFlow(t *testing.T) {
	helper.SkipIfNotManual(t)

	c := helper.NewClient()

	t.Log("准备工作: 登录管理员账户")
	_, err := c.Login("admin", "admin123")
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}
	t.Log("  登录成功")

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
	if err != nil {
		t.Fatalf("获取菜单列表失败: %v", err)
	}
	initialCount := len(menus)
	t.Logf("  菜单数量: %d", initialCount)

	// 测试 2: 创建菜单
	t.Log("\n测试 2: 创建菜单")
	menuName := fmt.Sprintf("testmenu_%d", time.Now().UnixNano())
	menuPath := "/test/" + menuName
	visible := true
	createReq := handler.CreateMenuRequest{
		Title:   menuName,
		Path:    menuPath,
		Icon:    "test-icon",
		Order:   99,
		Visible: &visible,
	}
	t.Logf("  创建菜单: %s", menuName)

	created, err := helper.Post[menu.CreateMenuResultDTO](c, "/api/admin/menus", createReq)
	if err != nil {
		t.Fatalf("创建菜单失败: %v", err)
	}
	createdMenuID = created.ID

	// 验证创建结果（创建 API 只返回 ID）
	if created.ID == 0 {
		t.Fatal("创建菜单后 ID 不应为 0")
	}
	t.Logf("  创建成功! 菜单 ID: %d", created.ID)

	// 测试 3: 获取菜单详情并验证字段
	t.Log("\n测试 3: 获取菜单详情并验证字段")
	detail, err := helper.Get[menu.MenuDTO](c, fmt.Sprintf("/api/admin/menus/%d", created.ID), nil)
	if err != nil {
		t.Fatalf("获取菜单详情失败: %v", err)
	}
	if detail.ID != created.ID {
		t.Errorf("详情 ID 不匹配: 期望 %d, 实际 %d", created.ID, detail.ID)
	}
	if detail.Title != menuName {
		t.Errorf("标题不匹配: 期望 %q, 实际 %q", menuName, detail.Title)
	}
	if detail.Path != menuPath {
		t.Errorf("路径不匹配: 期望 %q, 实际 %q", menuPath, detail.Path)
	}
	if detail.Icon != "test-icon" {
		t.Errorf("图标不匹配: 期望 %q, 实际 %q", "test-icon", detail.Icon)
	}
	if detail.Order != 99 {
		t.Errorf("排序不匹配: 期望 %d, 实际 %d", 99, detail.Order)
	}
	if !detail.Visible {
		t.Error("可见性应为 true")
	}
	t.Logf("  标题: %s, 路径: %s", detail.Title, detail.Path)
	t.Logf("  图标: %s, 可见: %v", detail.Icon, detail.Visible)

	// 测试 4: 更新菜单
	t.Log("\n测试 4: 更新菜单")
	newTitle := menuName + "_updated"
	newOrder := 100
	updateReq := handler.UpdateMenuRequest{
		Title: &newTitle,
		Order: &newOrder,
	}
	updated, err := helper.Put[menu.MenuDTO](c, fmt.Sprintf("/api/admin/menus/%d", created.ID), updateReq)
	if err != nil {
		t.Fatalf("更新菜单失败: %v", err)
	}
	if updated.Title != newTitle {
		t.Errorf("更新后标题不匹配: 期望 %q, 实际 %q", newTitle, updated.Title)
	}
	if updated.Order != newOrder {
		t.Errorf("更新后排序不匹配: 期望 %d, 实际 %d", newOrder, updated.Order)
	}
	// 验证未更新的字段保持不变
	if updated.Path != menuPath {
		t.Errorf("未更新的路径不应改变: 期望 %q, 实际 %q", menuPath, updated.Path)
	}
	t.Logf("  更新成功! 新标题: %s", updated.Title)

	// 测试 5: 创建子菜单
	t.Log("\n测试 5: 创建子菜单")
	childReq := handler.CreateMenuRequest{
		Title:    menuName + "_child",
		Path:     "/test/" + menuName + "/child",
		Icon:     "child-icon",
		ParentID: &created.ID,
		Order:    1,
		Visible:  &visible,
	}
	childResult, err := helper.Post[menu.CreateMenuResultDTO](c, "/api/admin/menus", childReq)
	if err != nil {
		t.Fatalf("创建子菜单失败: %v", err)
	}
	childMenuID = childResult.ID
	if childResult.ID == 0 {
		t.Fatal("子菜单 ID 不应为 0")
	}
	t.Logf("  子菜单创建成功! ID: %d", childResult.ID)

	// 测试 6: 验证列表数量增加
	t.Log("\n测试 6: 验证菜单列表数量")
	menusAfter, _, err := helper.GetList[menu.MenuDTO](c, "/api/admin/menus", nil)
	if err != nil {
		t.Fatalf("获取菜单列表失败: %v", err)
	}
	if len(menusAfter) < initialCount+2 {
		t.Errorf("菜单数量应至少增加 2: 初始 %d, 现在 %d", initialCount, len(menusAfter))
	}
	t.Logf("  菜单数量: %d (增加了 %d)", len(menusAfter), len(menusAfter)-initialCount)

	// 测试 7: 删除子菜单
	t.Log("\n测试 7: 删除子菜单")
	err = c.Delete(fmt.Sprintf("/api/admin/menus/%d", childResult.ID))
	if err != nil {
		t.Fatalf("删除子菜单失败: %v", err)
	}
	childMenuID = 0 // 已删除，清理时不需要再删
	t.Log("  子菜单删除成功!")

	// 测试 8: 删除父菜单
	t.Log("\n测试 8: 删除父菜单")
	err = c.Delete(fmt.Sprintf("/api/admin/menus/%d", created.ID))
	if err != nil {
		t.Fatalf("删除菜单失败: %v", err)
	}
	createdMenuID = 0 // 已删除，清理时不需要再删
	t.Log("  菜单删除成功!")

	// 测试 9: 验证删除后获取应失败
	t.Log("\n测试 9: 验证删除后获取应返回 404")
	_, err = helper.Get[menu.MenuDTO](c, fmt.Sprintf("/api/admin/menus/%d", created.ID), nil)
	if err == nil {
		t.Error("删除后获取应返回错误")
	} else {
		t.Logf("  正确返回错误: %v", err)
	}

	t.Log("\n菜单管理流程测试完成!")
}

// TestListMenus 测试获取菜单列表。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestListMenus ./internal/manualtest/
func TestListMenus(t *testing.T) {
	helper.SkipIfNotManual(t)

	c := helper.NewClient()

	t.Log("获取菜单列表...")
	_, err := c.Login("admin", "admin123")
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}

	menus, meta, err := helper.GetList[menu.MenuDTO](c, "/api/admin/menus", nil)
	if err != nil {
		t.Fatalf("获取菜单列表失败: %v", err)
	}

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
