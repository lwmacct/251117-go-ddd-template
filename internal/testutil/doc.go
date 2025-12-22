// Package testutil 提供测试工具和公共 Mock 实现。
//
// 本包集中管理测试中常用的工具函数和 Mock 结构，避免在各测试文件中重复定义。
//
// # Mock 仓储
//
// 每个领域模块的 Repository 接口都有对应的 Mock 实现：
//   - MockUserCommandRepo / MockUserQueryRepo
//   - MockRoleCommandRepo / MockRoleQueryRepo
//   - MockMenuCommandRepo / MockMenuQueryRepo
//   - MockSettingCommandRepo / MockSettingQueryRepo
//
// # 使用示例
//
//	func TestSomething(t *testing.T) {
//	    userRepo := testutil.NewMockUserQueryRepo()
//	    userRepo.ExistingUsernames["admin"] = true
//	    // 使用 userRepo 进行测试...
//	}
//
// # 辅助函数
//
//   - PtrUint(v uint) *uint: 创建 uint 指针
//   - PtrString(v string) *string: 创建 string 指针
//   - PtrTime(v time.Time) *time.Time: 创建 time.Time 指针
package testutil
