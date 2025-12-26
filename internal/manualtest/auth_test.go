package manualtest

import (
	"fmt"
	"testing"
	"time"

	"github.com/lwmacct/251117-go-ddd-template/internal/application/auth"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/user"
	"github.com/lwmacct/251117-go-ddd-template/internal/manualtest/helper"
)

// TestLoginSuccess 测试使用正确凭证登录。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestLoginSuccess ./internal/manualtest/
func TestLoginSuccess(t *testing.T) {
	helper.SkipIfNotManual(t)

	c := helper.NewClient()

	t.Log("使用默认管理员账户登录...")
	resp, err := c.Login("admin", "admin123")
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}

	// 验证关键字段
	if resp.AccessToken == "" {
		t.Fatal("登录成功但未返回 access_token")
	}
	if resp.RefreshToken == "" {
		t.Fatal("登录成功但未返回 refresh_token")
	}
	if resp.User.UserID == 0 {
		t.Fatal("登录成功但未返回有效用户 ID")
	}
	if resp.User.Username == "" {
		t.Fatal("登录成功但未返回用户名")
	}

	t.Logf("登录成功!")
	t.Logf("  Access Token: %s...", resp.AccessToken[:50])
	t.Logf("  Refresh Token: %s...", resp.RefreshToken[:50])
	t.Logf("  用户 ID: %d", resp.User.UserID)
	t.Logf("  用户名: %s", resp.User.Username)
}

// TestLoginWrongPassword 测试使用错误密码登录。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestLoginWrongPassword ./internal/manualtest/
func TestLoginWrongPassword(t *testing.T) {
	helper.SkipIfNotManual(t)

	c := helper.NewClient()

	captcha, err := c.GetCaptcha()
	if err != nil {
		t.Fatalf("获取验证码失败: %v", err)
	}

	t.Log("使用错误密码登录...")
	req := auth.LoginDTO{
		Account:   "admin",
		Password:  "wrong_password",
		CaptchaID: captcha.ID,
		Captcha:   captcha.Code,
	}

	_, err = c.LoginWithCaptcha(req)
	if err == nil {
		t.Fatal("错误密码应该返回错误")
	}

	t.Logf("错误密码被正确拒绝: %v", err)
}

// TestLoginWrongCaptcha 测试使用错误验证码登录。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestLoginWrongCaptcha ./internal/manualtest/
func TestLoginWrongCaptcha(t *testing.T) {
	helper.SkipIfNotManual(t)

	c := helper.NewClient()

	captcha, err := c.GetCaptcha()
	if err != nil {
		t.Fatalf("获取验证码失败: %v", err)
	}

	t.Log("使用错误验证码登录...")
	req := auth.LoginDTO{
		Account:   "admin",
		Password:  "admin123",
		CaptchaID: captcha.ID,
		Captcha:   "0000",
	}

	_, err = c.LoginWithCaptcha(req)
	if err == nil {
		t.Fatal("错误验证码应该返回错误")
	}

	t.Logf("错误验证码被正确拒绝: %v", err)
}

// TestGetCaptcha 测试获取验证码。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestGetCaptcha ./internal/manualtest/
func TestGetCaptcha(t *testing.T) {
	helper.SkipIfNotManual(t)

	c := helper.NewClient()

	t.Log("获取验证码（开发模式）...")
	captcha, err := c.GetCaptcha()
	if err != nil {
		t.Fatalf("获取验证码失败: %v", err)
	}

	if captcha.ID == "" {
		t.Fatal("验证码 ID 为空")
	}
	if captcha.Code == "" {
		t.Fatal("验证码答案为空（开发模式应返回）")
	}

	t.Logf("验证码获取成功!")
	t.Logf("  ID: %s", captcha.ID)
	t.Logf("  Code: %s", captcha.Code)
}

// TestAuthFlow 完整认证流程测试。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestAuthFlow ./internal/manualtest/
func TestAuthFlow(t *testing.T) {
	helper.SkipIfNotManual(t)

	c := helper.NewClient()

	t.Log("步骤 1: 获取验证码")
	captcha, err := c.GetCaptcha()
	if err != nil {
		t.Fatalf("获取验证码失败: %v", err)
	}
	t.Logf("  验证码 ID: %s", captcha.ID)

	t.Log("步骤 2: 登录")
	loginResp, err := c.Login("admin", "admin123")
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}
	if loginResp.AccessToken == "" {
		t.Fatal("登录成功但未返回 token")
	}
	t.Logf("  登录成功，获取到 token")

	t.Log("步骤 3: 访问用户列表（验证 token）")
	resp, err := c.R().
		SetQueryParams(map[string]string{"page": "1", "limit": "1"}).
		Get("/api/admin/users")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if resp.IsError() {
		t.Fatalf("预期状态码 200，实际 %d", resp.StatusCode())
	}
	t.Log("  Token 验证成功，可以访问受保护资源")

	t.Log("认证流程测试完成!")
}

// TestRegister 测试用户注册。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestRegister ./internal/manualtest/
func TestRegister(t *testing.T) {
	helper.SkipIfNotManual(t)

	c := helper.NewClient()

	// 生成唯一用户名
	testUsername := fmt.Sprintf("reguser_%d", time.Now().Unix())
	testEmail := testUsername + "@example.com"

	t.Log("测试用户注册...")
	t.Logf("  用户名: %s", testUsername)
	t.Logf("  邮箱: %s", testEmail)

	registerReq := auth.RegisterDTO{
		Username: testUsername,
		Email:    testEmail,
		Password: "password123",
		FullName: "注册测试用户",
	}

	resp, err := helper.Post[auth.RegisterResultDTO](c, "/api/auth/register", registerReq)
	if err != nil {
		t.Fatalf("注册失败: %v", err)
	}

	if resp.UserID == 0 {
		t.Fatal("注册成功但未返回 user_id")
	}
	if resp.AccessToken == "" {
		t.Fatal("注册成功但未返回 access_token")
	}

	t.Logf("注册成功!")
	t.Logf("  User ID: %d", resp.UserID)
	t.Logf("  Username: %s", resp.Username)
	t.Logf("  Access Token: %s...", resp.AccessToken[:50])

	// 清理：删除测试用户
	t.Log("清理：登录管理员删除测试用户...")
	adminClient := helper.NewClient()
	_, err = adminClient.Login("admin", "admin123")
	if err != nil {
		t.Logf("警告：无法登录管理员账户清理测试用户: %v", err)
		return
	}
	err = adminClient.Delete(fmt.Sprintf("/api/admin/users/%d", resp.UserID))
	if err != nil {
		t.Logf("警告：无法删除测试用户: %v", err)
	} else {
		t.Log("  测试用户已删除")
	}
}

// TestRegisterDuplicate 测试注册重复用户名。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestRegisterDuplicate ./internal/manualtest/
func TestRegisterDuplicate(t *testing.T) {
	helper.SkipIfNotManual(t)

	c := helper.NewClient()

	// 生成唯一用户名
	testUsername := fmt.Sprintf("dupuser_%d", time.Now().Unix())
	testEmail := testUsername + "@example.com"

	t.Log("步骤 1: 注册第一个用户")
	registerReq := auth.RegisterDTO{
		Username: testUsername,
		Email:    testEmail,
		Password: "password123",
		FullName: "重复测试用户",
	}

	firstResp, err := helper.Post[auth.RegisterResultDTO](c, "/api/auth/register", registerReq)
	if err != nil {
		t.Fatalf("首次注册失败: %v", err)
	}
	t.Logf("  首次注册成功，用户 ID: %d", firstResp.UserID)

	t.Log("步骤 2: 尝试注册同名用户")
	duplicateReq := auth.RegisterDTO{
		Username: testUsername, // 相同用户名
		Email:    "another@example.com",
		Password: "password456",
		FullName: "重复测试用户2",
	}

	_, err = helper.Post[auth.RegisterResultDTO](c, "/api/auth/register", duplicateReq)
	if err == nil {
		t.Fatal("重复用户名应该返回错误")
	}
	t.Logf("  重复用户名被正确拒绝: %v", err)

	// 清理
	t.Log("步骤 3: 清理测试用户")
	adminClient := helper.NewClient()
	_, err = adminClient.Login("admin", "admin123")
	if err != nil {
		t.Logf("警告：无法登录管理员账户: %v", err)
		return
	}
	err = adminClient.Delete(fmt.Sprintf("/api/admin/users/%d", firstResp.UserID))
	if err != nil {
		t.Logf("警告：无法删除测试用户: %v", err)
	} else {
		t.Log("  测试用户已删除")
	}
}

// TestRefreshToken 测试刷新访问令牌。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestRefreshToken ./internal/manualtest/
func TestRefreshToken(t *testing.T) {
	helper.SkipIfNotManual(t)

	c := helper.NewClient()

	t.Log("步骤 1: 登录获取 token")
	loginResp, err := c.Login("admin", "admin123")
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}
	if loginResp.RefreshToken == "" {
		t.Fatal("登录成功但未返回 refresh_token")
	}
	t.Logf("  获取到 refresh_token: %s...", loginResp.RefreshToken[:50])

	t.Log("步骤 2: 使用 refresh_token 刷新")
	refreshReq := auth.RefreshTokenDTO{
		RefreshToken: loginResp.RefreshToken,
	}

	newTokens, err := helper.Post[auth.RefreshTokenResultDTO](c, "/api/auth/refresh", refreshReq)
	if err != nil {
		t.Fatalf("刷新 token 失败: %v", err)
	}

	if newTokens.AccessToken == "" {
		t.Fatal("刷新成功但未返回新 access_token")
	}

	t.Logf("Token 刷新成功!")
	t.Logf("  新 Access Token: %s...", newTokens.AccessToken[:50])
	t.Logf("  Token 类型: %s", newTokens.TokenType)
	t.Logf("  过期时间: %d 秒", newTokens.ExpiresIn)
}

// TestRefreshTokenInvalid 测试使用无效的 refresh token。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestRefreshTokenInvalid ./internal/manualtest/
func TestRefreshTokenInvalid(t *testing.T) {
	helper.SkipIfNotManual(t)

	c := helper.NewClient()

	t.Log("使用无效的 refresh_token 刷新")
	refreshReq := auth.RefreshTokenDTO{
		RefreshToken: "invalid_token_string",
	}

	_, err := helper.Post[auth.RefreshTokenResultDTO](c, "/api/auth/refresh", refreshReq)
	if err == nil {
		t.Fatal("无效的 refresh_token 应该返回错误")
	}

	t.Logf("无效 token 被正确拒绝: %v", err)
}

// TestGetCurrentUser 测试获取当前登录用户信息。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestGetCurrentUser ./internal/manualtest/
func TestGetCurrentUser(t *testing.T) {
	helper.SkipIfNotManual(t)

	c := helper.NewClient()

	t.Log("步骤 1: 登录")
	_, err := c.Login("admin", "admin123")
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}
	t.Log("  登录成功")

	t.Log("步骤 2: 获取当前用户信息")
	me, err := helper.Get[user.UserWithRolesDTO](c, "/api/user/profile", nil)
	if err != nil {
		t.Fatalf("获取当前用户失败: %v", err)
	}

	if me.ID == 0 {
		t.Fatal("返回的用户 ID 为 0")
	}
	if me.Username == "" {
		t.Fatal("返回的用户名为空")
	}

	t.Logf("获取成功!")
	t.Logf("  ID: %d", me.ID)
	t.Logf("  用户名: %s", me.Username)
	t.Logf("  邮箱: %s", me.Email)
	t.Logf("  角色数量: %d", len(me.Roles))
	for _, role := range me.Roles {
		t.Logf("    - %s (%s)", role.DisplayName, role.Name)
	}
}
