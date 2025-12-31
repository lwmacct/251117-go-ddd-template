// Package manualtest 提供手动集成测试辅助工具。
//
// 本包提供 HTTP 客户端和工厂函数，用于针对运行中的服务进行 API 集成测试。
// 测试需要手动触发，通过 MANUAL 环境变量控制执行。
//
// 运行方式：
//
//	# 运行所有测试
//	MANUAL=1 go test -v ./internal/manualtest/...
//
//	# 运行单个模块测试
//	MANUAL=1 go test -v ./internal/manualtest/auth/
//
//	# 运行单个测试函数
//	MANUAL=1 go test -v -run TestLoginScenarios ./internal/manualtest/auth/
//
// 核心类型：
//   - [Client]: HTTP 测试客户端，封装 resty 库
//   - [Get], [Post], [Put], [Delete]: 泛型 HTTP 方法
//   - [CreateTestUser], [CreateTestRole]: 资源工厂函数
//   - [LoginAsAdmin], [LoginAs]: 登录辅助函数
//
// 各领域模块测试位于对应子包：
//   - auth/: 认证测试
//   - user/: 用户管理测试
//   - role/: 角色管理测试
//   - profile/: 个人资料测试
//   - pat/: 个人访问令牌测试
//   - twofa/: 双因素认证测试
//   - setting/: 系统配置测试
//   - usersetting/: 用户设置测试
//   - cache/: 缓存管理测试
//   - task/: 任务管理测试
//   - system/: 系统状态测试
//   - auditlog/: 审计日志测试
//   - errors/: 错误处理测试
package manualtest
