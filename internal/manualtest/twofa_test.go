package manualtest

import (
	"fmt"
	"testing"
	"time"

	"github.com/pquerna/otp/totp"

	"github.com/lwmacct/251117-go-ddd-template/internal/application/twofa"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/user"
	"github.com/lwmacct/251117-go-ddd-template/internal/manualtest/helper"
)

// TestGetTwoFAStatus 测试获取 2FA 状态。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestGetTwoFAStatus ./internal/manualtest/
func TestGetTwoFAStatus(t *testing.T) {
	helper.SkipIfNotManual(t)

	c := helper.NewClient()

	t.Log("步骤 1: 登录")
	_, err := c.Login("admin", "admin123")
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}
	t.Log("  登录成功")

	t.Log("步骤 2: 获取 2FA 状态")
	status, err := helper.Get[twofa.StatusDTO](c, "/api/auth/2fa/status", nil)
	if err != nil {
		t.Fatalf("获取 2FA 状态失败: %v", err)
	}

	t.Logf("2FA 状态获取成功!")
	t.Logf("  启用状态: %v", status.Enabled)
	t.Logf("  剩余恢复码数量: %d", status.RecoveryCodesCount)
}

// TestTwoFAFlow 2FA 完整流程测试（设置 → 验证启用 → 状态检查 → 禁用）。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestTwoFAFlow ./internal/manualtest/
func TestTwoFAFlow(t *testing.T) {
	helper.SkipIfNotManual(t)

	// 创建测试用户
	adminClient := helper.NewClient()
	_, err := adminClient.Login("admin", "admin123")
	if err != nil {
		t.Fatalf("管理员登录失败: %v", err)
	}

	testUsername := fmt.Sprintf("twofa_test_%d", time.Now().Unix())
	testPassword := "password123"

	t.Log("步骤 1: 创建测试用户（带 user 角色）")
	createReq := user.CreateUserDTO{
		Username: testUsername,
		Email:    testUsername + "@example.com",
		Password: testPassword,
		FullName: "2FA 测试用户",
		RoleIDs:  []uint{2}, // user 角色 ID
	}

	createResp, err := helper.Post[user.UserWithRolesDTO](adminClient, "/api/admin/users", createReq)
	if err != nil {
		t.Fatalf("创建测试用户失败: %v", err)
	}
	testUserID := createResp.ID
	t.Logf("  创建成功，用户 ID: %d", testUserID)

	// 用测试用户登录
	t.Log("步骤 2: 测试用户登录")
	testClient := helper.NewClient()
	_, err = testClient.Login(testUsername, testPassword)
	if err != nil {
		t.Fatalf("测试用户登录失败: %v", err)
	}
	t.Log("  登录成功")

	// 设置 2FA
	t.Log("步骤 3: 设置 2FA")
	setup, err := helper.Post[twofa.SetupDTO](testClient, "/api/auth/2fa/setup", nil)
	if err != nil {
		t.Fatalf("设置 2FA 失败: %v", err)
	}
	if setup.Secret == "" {
		t.Fatal("2FA 密钥为空")
	}
	t.Logf("  密钥: %s", setup.Secret)
	t.Logf("  二维码 URL: %s", setup.QRCodeURL[:50]+"...")
	if setup.QRCodeImg != "" {
		t.Log("  二维码图片: [已生成]")
	}

	// 使用密钥生成 TOTP 代码
	t.Log("步骤 4: 生成并验证 TOTP 代码")
	code, err := totp.GenerateCode(setup.Secret, time.Now())
	if err != nil {
		t.Fatalf("生成 TOTP 代码失败: %v", err)
	}
	t.Logf("  生成的 TOTP 代码: %s", code)

	// 验证并启用 2FA
	verifyReq := map[string]string{"code": code}
	enableResp, err := helper.Post[twofa.EnableDTO](testClient, "/api/auth/2fa/verify", verifyReq)
	if err != nil {
		t.Fatalf("验证并启用 2FA 失败: %v", err)
	}
	t.Logf("  2FA 启用成功!")
	t.Logf("  恢复码数量: %d", len(enableResp.RecoveryCodes))
	if len(enableResp.RecoveryCodes) > 0 {
		t.Logf("  第一个恢复码: %s", enableResp.RecoveryCodes[0])
	}

	// 检查 2FA 状态
	t.Log("步骤 5: 检查 2FA 状态")
	status, err := helper.Get[twofa.StatusDTO](testClient, "/api/auth/2fa/status", nil)
	if err != nil {
		t.Fatalf("获取 2FA 状态失败: %v", err)
	}
	if !status.Enabled {
		t.Fatal("2FA 应该已启用")
	}
	t.Logf("  2FA 已启用，剩余恢复码: %d", status.RecoveryCodesCount)

	// 禁用 2FA
	t.Log("步骤 6: 禁用 2FA")
	resp, err := testClient.R().Post("/api/auth/2fa/disable")
	if err != nil {
		t.Fatalf("禁用 2FA 请求失败: %v", err)
	}
	if resp.IsError() {
		t.Fatalf("禁用 2FA 失败，状态码: %d", resp.StatusCode())
	}
	t.Log("  2FA 已禁用")

	// 验证 2FA 已禁用
	status2, err := helper.Get[twofa.StatusDTO](testClient, "/api/auth/2fa/status", nil)
	if err != nil {
		t.Fatalf("获取 2FA 状态失败: %v", err)
	}
	if status2.Enabled {
		t.Fatal("2FA 应该已禁用")
	}
	t.Log("  验证：2FA 已禁用")

	// 清理
	t.Log("步骤 7: 清理测试用户")
	err = adminClient.Delete(fmt.Sprintf("/api/admin/users/%d", testUserID))
	if err != nil {
		t.Logf("警告：无法删除测试用户: %v", err)
	} else {
		t.Log("  测试用户已删除")
	}

	t.Log("2FA 完整流程测试完成!")
}

// TestSetup2FA 测试设置 2FA。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestSetup2FA ./internal/manualtest/
func TestSetup2FA(t *testing.T) {
	helper.SkipIfNotManual(t)

	// 创建测试用户
	adminClient := helper.NewClient()
	_, err := adminClient.Login("admin", "admin123")
	if err != nil {
		t.Fatalf("管理员登录失败: %v", err)
	}

	testUsername := fmt.Sprintf("setup2fa_%d", time.Now().Unix())

	t.Log("步骤 1: 创建测试用户")
	createReq := user.CreateUserDTO{
		Username: testUsername,
		Email:    testUsername + "@example.com",
		Password: "password123",
		FullName: "2FA Setup 测试用户",
		RoleIDs:  []uint{2},
	}

	createResp, err := helper.Post[user.UserWithRolesDTO](adminClient, "/api/admin/users", createReq)
	if err != nil {
		t.Fatalf("创建测试用户失败: %v", err)
	}
	testUserID := createResp.ID
	t.Logf("  创建成功，用户 ID: %d", testUserID)

	// 用测试用户登录
	t.Log("步骤 2: 测试用户登录")
	c := helper.NewClient()
	_, err = c.Login(testUsername, "password123")
	if err != nil {
		t.Fatalf("测试用户登录失败: %v", err)
	}
	t.Log("  登录成功")

	t.Log("步骤 3: 设置 2FA")
	setup, err := helper.Post[twofa.SetupDTO](c, "/api/auth/2fa/setup", nil)
	if err != nil {
		t.Fatalf("设置 2FA 失败: %v", err)
	}

	if setup.Secret == "" {
		t.Fatal("2FA 密钥为空")
	}
	if setup.QRCodeURL == "" {
		t.Fatal("二维码 URL 为空")
	}
	if setup.QRCodeImg == "" {
		t.Fatal("二维码图片为空")
	}

	t.Logf("2FA 设置成功!")
	t.Logf("  密钥: %s", setup.Secret)
	t.Logf("  二维码 URL 长度: %d", len(setup.QRCodeURL))
	t.Logf("  二维码图片大小: %d bytes", len(setup.QRCodeImg))

	// 清理
	t.Log("步骤 4: 清理测试用户")
	err = adminClient.Delete(fmt.Sprintf("/api/admin/users/%d", testUserID))
	if err != nil {
		t.Logf("警告：无法删除测试用户: %v", err)
	} else {
		t.Log("  测试用户已删除")
	}
}

// TestDisable2FA 测试禁用 2FA。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestDisable2FA ./internal/manualtest/
func TestDisable2FA(t *testing.T) {
	helper.SkipIfNotManual(t)

	c := helper.NewClient()

	t.Log("步骤 1: 登录")
	_, err := c.Login("admin", "admin123")
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}
	t.Log("  登录成功")

	t.Log("步骤 2: 禁用 2FA")
	resp, err := c.R().Post("/api/auth/2fa/disable")
	if err != nil {
		t.Fatalf("禁用 2FA 请求失败: %v", err)
	}

	// 即使 2FA 未启用，禁用也应该成功（幂等操作）
	if resp.IsError() {
		t.Logf("禁用 2FA 返回状态码: %d，响应: %s", resp.StatusCode(), resp.String())
	} else {
		t.Log("  禁用 2FA 成功")
	}

	// 检查状态
	t.Log("步骤 3: 检查 2FA 状态")
	status, err := helper.Get[twofa.StatusDTO](c, "/api/auth/2fa/status", nil)
	if err != nil {
		t.Fatalf("获取 2FA 状态失败: %v", err)
	}

	t.Logf("  2FA 状态: 启用=%v, 恢复码=%d", status.Enabled, status.RecoveryCodesCount)
}
