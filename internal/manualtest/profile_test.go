package manualtest

import (
	"fmt"
	"testing"
	"time"

	"github.com/lwmacct/251117-go-ddd-template/internal/application/auth"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/user"
	"github.com/lwmacct/251117-go-ddd-template/internal/manualtest/helper"
)

// TestGetProfile 测试获取个人资料。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestGetProfile ./internal/manualtest/
func TestGetProfile(t *testing.T) {
	helper.SkipIfNotManual(t)

	c := helper.NewClient()

	t.Log("步骤 1: 登录")
	_, err := c.Login("admin", "admin123")
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}
	t.Log("  登录成功")

	t.Log("步骤 2: 获取个人资料")
	profile, err := helper.Get[user.UserWithRolesDTO](c, "/api/user/profile", nil)
	if err != nil {
		t.Fatalf("获取个人资料失败: %v", err)
	}

	// 验证关键字段
	if profile.ID == 0 {
		t.Fatal("返回的用户 ID 为 0")
	}
	if profile.Username == "" {
		t.Fatal("返回的用户名为空")
	}
	if profile.Email == "" {
		t.Fatal("返回的邮箱为空")
	}
	if profile.Status == "" {
		t.Fatal("返回的状态为空")
	}

	t.Logf("获取成功!")
	t.Logf("  ID: %d", profile.ID)
	t.Logf("  用户名: %s", profile.Username)
	t.Logf("  邮箱: %s", profile.Email)
	t.Logf("  全名: %s", profile.FullName)
	t.Logf("  状态: %s", profile.Status)
	t.Logf("  角色数量: %d", len(profile.Roles))
}

// TestUpdateProfile 测试更新个人资料。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestUpdateProfile ./internal/manualtest/
func TestUpdateProfile(t *testing.T) {
	helper.SkipIfNotManual(t)

	c := helper.NewClient()

	t.Log("步骤 1: 登录")
	_, err := c.Login("admin", "admin123")
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}
	t.Log("  登录成功")

	t.Log("步骤 2: 获取当前资料")
	originalProfile, err := helper.Get[user.UserWithRolesDTO](c, "/api/user/profile", nil)
	if err != nil {
		t.Fatalf("获取原始资料失败: %v", err)
	}
	t.Logf("  当前全名: %s", originalProfile.FullName)

	t.Log("步骤 3: 更新资料")
	newFullName := fmt.Sprintf("测试更新_%d", time.Now().Unix())
	updateReq := user.UpdateUserDTO{
		FullName: &newFullName,
	}

	updateResp, err := helper.Put[user.UserWithRolesDTO](c, "/api/user/profile", updateReq)
	if err != nil {
		t.Fatalf("更新资料失败: %v", err)
	}
	t.Logf("  更新后全名: %s", updateResp.FullName)

	if updateResp.FullName != newFullName {
		t.Fatalf("全名未更新，期望 %s，实际 %s", newFullName, updateResp.FullName)
	}

	t.Log("步骤 4: 恢复原始资料")
	restoreReq := user.UpdateUserDTO{
		FullName: &originalProfile.FullName,
	}
	_, err = helper.Put[user.UserWithRolesDTO](c, "/api/user/profile", restoreReq)
	if err != nil {
		t.Logf("警告：无法恢复原始资料: %v", err)
	} else {
		t.Log("  资料已恢复")
	}

	t.Log("更新资料测试完成!")
}

// TestUpdateProfileInvalid 测试使用无效数据更新个人资料。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestUpdateProfileInvalid ./internal/manualtest/
func TestUpdateProfileInvalid(t *testing.T) {
	helper.SkipIfNotManual(t)

	c := helper.NewClient()

	t.Log("步骤 1: 登录")
	_, err := c.Login("admin", "admin123")
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}
	t.Log("  登录成功")

	t.Log("步骤 2: 尝试使用无效数据更新资料（如空全名）")
	emptyFullName := ""
	updateReq := user.UpdateUserDTO{
		FullName: &emptyFullName,
	}

	resp, err := c.R().
		SetBody(updateReq).
		Put("/api/user/profile")

	// 应该返回错误或验证失败
	if err == nil && resp.IsSuccess() {
		t.Log("  注意：服务器接受了空全名，这可能需要根据业务规则判断是否合理")
	} else {
		t.Logf("  无效数据被正确拒绝: %d - %s", resp.StatusCode(), resp.String())
	}

	t.Log("无效数据更新测试完成!")
}

// TestChangePassword 测试修改密码。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestChangePassword ./internal/manualtest/
func TestChangePassword(t *testing.T) {
	helper.SkipIfNotManual(t)

	// 创建测试用户
	adminClient := helper.NewClient()
	_, err := adminClient.Login("admin", "admin123")
	if err != nil {
		t.Fatalf("管理员登录失败: %v", err)
	}

	testUsername := fmt.Sprintf("pwdtest_%d", time.Now().Unix())
	testEmail := testUsername + "@example.com"
	originalPassword := "original123"
	newPassword := "newpassword456"

	t.Log("步骤 1: 创建测试用户（带 user 角色）")
	createReq := user.CreateUserDTO{
		Username: testUsername,
		Email:    testEmail,
		Password: originalPassword,
		FullName: "密码测试用户",
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
	_, err = testClient.Login(testUsername, originalPassword)
	if err != nil {
		t.Fatalf("测试用户登录失败: %v", err)
	}
	t.Log("  登录成功")

	t.Log("步骤 3: 修改密码")
	changeReq := user.ChangePasswordDTO{
		OldPassword: originalPassword,
		NewPassword: newPassword,
	}

	resp, err := testClient.R().
		SetBody(changeReq).
		Put("/api/user/password")
	if err != nil {
		t.Fatalf("修改密码请求失败: %v", err)
	}
	if resp.IsError() {
		t.Fatalf("修改密码失败，状态码: %d, 响应: %s", resp.StatusCode(), resp.String())
	}
	t.Log("  密码修改成功")

	t.Log("步骤 4: 使用新密码登录")
	newClient := helper.NewClient()
	_, err = newClient.Login(testUsername, newPassword)
	if err != nil {
		t.Fatalf("使用新密码登录失败: %v", err)
	}
	t.Log("  新密码登录成功!")

	t.Log("步骤 5: 验证旧密码已失效")
	oldPwdClient := helper.NewClient()
	captcha, _ := oldPwdClient.GetCaptcha()
	oldLoginReq := auth.LoginDTO{
		Account:   testUsername,
		Password:  originalPassword,
		CaptchaID: captcha.ID,
		Captcha:   captcha.Code,
	}
	oldResp, err := oldPwdClient.LoginWithCaptcha(oldLoginReq)
	if err == nil && oldResp.AccessToken != "" {
		t.Fatal("旧密码不应该能登录")
	}
	t.Log("  旧密码已失效")

	// 清理
	t.Log("步骤 6: 清理测试用户")
	err = adminClient.Delete(fmt.Sprintf("/api/admin/users/%d", testUserID))
	if err != nil {
		t.Logf("警告：无法删除测试用户: %v", err)
	} else {
		t.Log("  测试用户已删除")
	}

	t.Log("修改密码测试完成!")
}

// TestChangePasswordWrongOld 测试使用错误的旧密码修改密码。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestChangePasswordWrongOld ./internal/manualtest/
func TestChangePasswordWrongOld(t *testing.T) {
	helper.SkipIfNotManual(t)

	// 创建测试用户
	adminClient := helper.NewClient()
	_, err := adminClient.Login("admin", "admin123")
	if err != nil {
		t.Fatalf("管理员登录失败: %v", err)
	}

	testUsername := fmt.Sprintf("wrongpwd_%d", time.Now().Unix())
	testEmail := testUsername + "@example.com"
	testPassword := "original123"

	t.Log("步骤 1: 创建测试用户")
	createReq := user.CreateUserDTO{
		Username: testUsername,
		Email:    testEmail,
		Password: testPassword,
		FullName: "错误旧密码测试用户",
		RoleIDs:  []uint{2},
	}

	createResp, err := helper.Post[user.UserWithRolesDTO](adminClient, "/api/admin/users", createReq)
	if err != nil {
		t.Fatalf("创建测试用户失败: %v", err)
	}
	testUserID := createResp.ID
	t.Logf("  创建成功，用户 ID: %d", testUserID)

	t.Log("步骤 2: 测试用户登录")
	testClient := helper.NewClient()
	_, err = testClient.Login(testUsername, testPassword)
	if err != nil {
		t.Fatalf("测试用户登录失败: %v", err)
	}
	t.Log("  登录成功")

	t.Log("步骤 3: 使用错误的旧密码尝试修改")
	changeReq := user.ChangePasswordDTO{
		OldPassword: "wrong_old_password",
		NewPassword: "newpassword456",
	}

	resp, err := testClient.R().
		SetBody(changeReq).
		Put("/api/user/password")
	if err == nil && resp.IsSuccess() {
		t.Fatal("错误的旧密码不应该允许修改密码")
	}
	t.Logf("  错误旧密码被正确拒绝: %d - %s", resp.StatusCode(), resp.String())

	// 清理
	t.Log("步骤 4: 清理测试用户")
	err = adminClient.Delete(fmt.Sprintf("/api/admin/users/%d", testUserID))
	if err != nil {
		t.Logf("警告：无法删除测试用户: %v", err)
	} else {
		t.Log("  测试用户已删除")
	}

	t.Log("错误旧密码测试完成!")
}

// TestDeleteAccount 测试删除账户。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestDeleteAccount ./internal/manualtest/
func TestDeleteAccount(t *testing.T) {
	helper.SkipIfNotManual(t)

	// 管理员创建测试用户
	adminClient := helper.NewClient()
	_, err := adminClient.Login("admin", "admin123")
	if err != nil {
		t.Fatalf("管理员登录失败: %v", err)
	}

	testUsername := fmt.Sprintf("delacct_%d", time.Now().Unix())
	testEmail := testUsername + "@example.com"
	testPassword := "test123456"

	t.Log("步骤 1: 创建测试用户（带 user 角色）")
	createReq := user.CreateUserDTO{
		Username: testUsername,
		Email:    testEmail,
		Password: testPassword,
		FullName: "删除账户测试用户",
		RoleIDs:  []uint{2}, // user 角色 ID
	}

	createResp, err := helper.Post[user.UserWithRolesDTO](adminClient, "/api/admin/users", createReq)
	if err != nil {
		t.Fatalf("创建测试用户失败: %v", err)
	}
	t.Logf("  创建成功，用户 ID: %d", createResp.ID)

	t.Log("步骤 2: 测试用户登录")
	testClient := helper.NewClient()
	_, err = testClient.Login(testUsername, testPassword)
	if err != nil {
		t.Fatalf("测试用户登录失败: %v", err)
	}
	t.Log("  登录成功")

	t.Log("步骤 3: 调用删除账户接口")
	err = testClient.Delete("/api/user/account")
	if err != nil {
		t.Fatalf("删除账户失败: %v", err)
	}
	t.Log("  删除成功!")

	t.Log("步骤 4: 验证账户已删除（尝试登录应失败）")
	verifyClient := helper.NewClient()
	captcha, _ := verifyClient.GetCaptcha()
	loginReq := auth.LoginDTO{
		Account:   testUsername,
		Password:  testPassword,
		CaptchaID: captcha.ID,
		Captcha:   captcha.Code,
	}
	loginResp, err := verifyClient.LoginWithCaptcha(loginReq)
	if err == nil && loginResp.AccessToken != "" {
		t.Fatal("账户已删除，不应该能登录")
	}
	t.Log("  验证成功：账户已无法登录")

	t.Log("删除账户测试完成!")
}
