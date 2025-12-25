package manualtest

import (
	"fmt"
	"testing"
	"time"

	"github.com/lwmacct/251117-go-ddd-template/internal/application/pat"
	"github.com/lwmacct/251117-go-ddd-template/internal/manualtest/helper"
)

// 使用 time 包生成唯一标识符
var _ = time.Now

// TestPATFlow PAT 令牌完整流程测试。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestPATFlow ./internal/manualtest/
func TestPATFlow(t *testing.T) {
	helper.SkipIfNotManual(t)

	c := helper.NewClient()

	t.Log("准备工作: 登录管理员账户")
	_, err := c.Login("admin", "admin123")
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}
	t.Log("  登录成功")

	// 测试 1: 获取 PAT 列表
	t.Log("\n测试 1: 获取 PAT 列表")
	tokens, _, err := helper.GetList[pat.TokenDTO](c, "/api/user/tokens", nil)
	if err != nil {
		t.Fatalf("获取 PAT 列表失败: %v", err)
	}
	t.Logf("  现有 PAT 数量: %d", len(tokens))

	// 测试 2: 创建 PAT
	t.Log("\n测试 2: 创建 PAT")
	tokenName := fmt.Sprintf("test_pat_%d", time.Now().Unix())
	expiresIn := 30 // 30 天
	createReq := pat.CreateTokenDTO{
		Name:        tokenName,
		Permissions: []string{"user:profile:read"},
		ExpiresIn:   &expiresIn,
		Description: "测试用 PAT",
	}
	t.Logf("  创建 PAT: %s", tokenName)

	created, err := helper.Post[pat.CreateTokenResultDTO](c, "/api/user/tokens", createReq)
	if err != nil {
		t.Fatalf("创建 PAT 失败: %v", err)
	}
	t.Logf("  创建成功! PAT ID: %d", created.Token.ID)
	t.Logf("  明文令牌: %s... (仅显示一次)", created.PlainToken[:20])
	t.Logf("  状态: %s", created.Token.Status)

	tokenID := created.Token.ID

	// 测试 3: 获取 PAT 详情
	t.Log("\n测试 3: 获取 PAT 详情")
	detail, err := helper.Get[pat.TokenDTO](c, fmt.Sprintf("/api/user/tokens/%d", tokenID), nil)
	if err != nil {
		t.Fatalf("获取 PAT 详情失败: %v", err)
	}
	t.Logf("  名称: %s", detail.Name)
	t.Logf("  前缀: %s", detail.TokenPrefix)
	t.Logf("  权限: %v", detail.Permissions)
	t.Logf("  状态: %s", detail.Status)
	if detail.ExpiresAt != nil {
		t.Logf("  过期时间: %s", detail.ExpiresAt.Format("2006-01-02 15:04:05"))
	}

	// 测试 4: 禁用 PAT
	t.Log("\n测试 4: 禁用 PAT")
	resp, err := c.R().Patch(fmt.Sprintf("/api/user/tokens/%d/disable", tokenID))
	if err != nil {
		t.Fatalf("禁用 PAT 失败: %v", err)
	}
	if resp.IsError() {
		t.Fatalf("禁用 PAT 失败: 状态码 %d", resp.StatusCode())
	}
	t.Log("  禁用成功!")

	// 验证状态
	disabled, err := helper.Get[pat.TokenDTO](c, fmt.Sprintf("/api/user/tokens/%d", tokenID), nil)
	if err != nil {
		t.Fatalf("获取 PAT 详情失败: %v", err)
	}
	t.Logf("  当前状态: %s", disabled.Status)

	// 测试 5: 启用 PAT
	t.Log("\n测试 5: 启用 PAT")
	resp, err = c.R().Patch(fmt.Sprintf("/api/user/tokens/%d/enable", tokenID))
	if err != nil {
		t.Fatalf("启用 PAT 失败: %v", err)
	}
	if resp.IsError() {
		t.Fatalf("启用 PAT 失败: 状态码 %d", resp.StatusCode())
	}
	t.Log("  启用成功!")

	// 验证状态
	enabled, err := helper.Get[pat.TokenDTO](c, fmt.Sprintf("/api/user/tokens/%d", tokenID), nil)
	if err != nil {
		t.Fatalf("获取 PAT 详情失败: %v", err)
	}
	t.Logf("  当前状态: %s", enabled.Status)

	// 测试 6: 删除 PAT
	t.Log("\n测试 6: 删除 PAT")
	err = c.Delete(fmt.Sprintf("/api/user/tokens/%d", tokenID))
	if err != nil {
		t.Fatalf("删除 PAT 失败: %v", err)
	}
	t.Log("  删除成功!")

	t.Log("\nPAT 令牌流程测试完成!")
}

// TestListPATs 测试获取 PAT 列表。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestListPATs ./internal/manualtest/
func TestListPATs(t *testing.T) {
	helper.SkipIfNotManual(t)

	c := helper.NewClient()

	t.Log("获取 PAT 列表...")
	_, err := c.Login("admin", "admin123")
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}

	tokens, meta, err := helper.GetList[pat.TokenDTO](c, "/api/user/tokens", nil)
	if err != nil {
		t.Fatalf("获取 PAT 列表失败: %v", err)
	}

	t.Logf("PAT 数量: %d", len(tokens))
	if meta != nil {
		t.Logf("总数: %d", meta.Total)
	}

	for _, token := range tokens {
		statusIcon := "✓"
		if token.Status != "active" {
			statusIcon = "✗"
		}
		t.Logf("  [%s] [%d] %s (%s)", statusIcon, token.ID, token.Name, token.TokenPrefix)
	}
}

// TestPATWithPermissions 测试创建带特定权限的 PAT。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestPATWithPermissions ./internal/manualtest/
func TestPATWithPermissions(t *testing.T) {
	helper.SkipIfNotManual(t)

	c := helper.NewClient()

	t.Log("准备工作: 登录管理员账户")
	_, err := c.Login("admin", "admin123")
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}
	t.Log("  登录成功")

	// 创建带限制权限的 PAT
	t.Log("\n创建带限制权限的 PAT...")
	tokenName := fmt.Sprintf("limited_pat_%d", time.Now().Unix())
	createReq := pat.CreateTokenDTO{
		Name: tokenName,
		Permissions: []string{
			"user:profile:read",
			"user:profile:update",
		},
		Description: "仅限读写个人资料",
	}

	created, err := helper.Post[pat.CreateTokenResultDTO](c, "/api/user/tokens", createReq)
	if err != nil {
		t.Fatalf("创建 PAT 失败: %v", err)
	}
	t.Logf("  创建成功! PAT ID: %d", created.Token.ID)
	t.Logf("  权限: %v", created.Token.Permissions)

	// 清理
	t.Log("\n清理测试 PAT...")
	err = c.Delete(fmt.Sprintf("/api/user/tokens/%d", created.Token.ID))
	if err != nil {
		t.Logf("  警告：删除失败: %v", err)
	} else {
		t.Log("  清理完成")
	}
}
