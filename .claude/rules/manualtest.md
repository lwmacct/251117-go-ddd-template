---
paths:
  - "internal/manualtest/**/*.go"
---

# manualtest 手动测试规范

本包提供针对 HTTP API 的集成测试，需要服务运行时手动执行。

## 目录结构

```
internal/manualtest/
├── doc.go        # 包文档
├── client.go     # Client, HTTP 方法
├── factory.go    # 资源工厂函数
├── assert.go     # 提取函数
├── helper.go     # 测试辅助函数
└── {module}/     # 测试子包，每个模块一个目录
    └── {module}_test.go
```

**命名约定**：

- 根目录包名 `manualtest`，子目录包名 `{module}_test`（外部测试包模式）
- 每个 API 模块对应一个测试子包

## 运行方式

```bash
# 运行所有测试
MANUAL=1 go test -v ./internal/manualtest/...

# 串行执行（服务端压力大时）
MANUAL=1 go test -v -p 1 ./internal/manualtest/...

# 运行单个测试
MANUAL=1 go test -v -run TestLoginScenarios ./internal/manualtest/auth/
```

## DTO 使用原则

- **类型来源**：`internal/application/*/dto.go`
- **禁止定义任何 DTO，只消费 Application 层的类型。**

### 示例

```go
import (
    "github.com/lwmacct/251117-go-ddd-template/internal/application/auth"
    "github.com/lwmacct/251117-go-ddd-template/internal/application/user"
    "github.com/lwmacct/251117-go-ddd-template/internal/manualtest"
)

// 使用 Application 层 DTO 解析响应
result, err := manualtest.Post[auth.LoginResponseDTO](c, "/api/auth/login", req)
profile, err := manualtest.Get[user.UserWithRolesDTO](c, "/api/user/profile", nil)

// 创建用户后直接获取 DTO
createResp, err := manualtest.Post[user.UserWithRolesDTO](c, "/api/admin/users", req)
userID := createResp.ID  // 直接访问字段
```

## 设计反思原则

**测试困难是设计问题的信号。**

如果发现以下情况，说明 Application 层设计需要检视:

| 症状                     | 可能的设计问题                  |
| ------------------------ | ------------------------------- |
| 需要在测试中定义 DTO     | Application 层 DTO 缺失或不完整 |
| 响应结构难以断言         | Handler 响应格式与 DTO 不一致   |
| 需要复杂的类型转换       | DTO 设计不符合使用场景          |
| 测试代码比业务代码还复杂 | API 设计过于复杂                |

**正确的修复方向**：

1. 检查 Application 层 DTO 是否缺少 JSON tags
2. 检查 Handler 响应格式是否与 DTO 结构匹配
3. 在 Application 层补充缺失的 DTO

**❌ 禁止在 manualtest 中临时定义 DTO 来"绕过"问题。**

## 资源清理规范

**必须使用 `t.Cleanup()` 在创建资源后立即注册清理，禁止在测试末尾手动清理。**

```go
// ✅ 正确：创建后立即注册清理
createResp, _ := manualtest.Post[user.UserWithRolesDTO](c, "/api/admin/users", req)
t.Cleanup(func() {
    _ = c.Delete(fmt.Sprintf("/api/admin/users/%d", createResp.ID))
})

// ❌ 禁止：清理在测试末尾（中途失败不会执行）
```

**当测试本身包含删除操作时**：删除成功后将 ID 置为 0，Cleanup 中检查 `if id > 0`。

## Testify 断言规范

### require vs assert 选择

| 场景                       | 使用      | 理由               |
| -------------------------- | --------- | ------------------ |
| 前置条件（登录、创建资源） | `require` | 失败则后续无意义   |
| 业务验证（字段值比较）     | `assert`  | 收集所有错误后报告 |

```go
// ✅ 正确：登录失败直接停止
_, err := c.Login("admin", "admin123")
require.NoError(t, err, "登录失败")

// ✅ 正确：字段验证使用 assert
assert.Equal(t, expected, actual, "用户名不匹配")
```

### 集合验证

**禁止手动循环查找，使用 testify 断言：**

```go
// ❌ 禁止：手动循环
permIDMap := make(map[uint]bool)
for _, p := range perms { permIDMap[p.ID] = true }
assert.True(t, permIDMap[expectedID], "未找到")

// ✅ 正确：使用 assert.Contains
ids := manualtest.ExtractIDs(perms, func(p role.PermissionDTO) uint { return p.ID })
assert.Contains(t, ids, expectedID, "未找到权限 ID")

// ✅ 正确：批量验证使用 ElementsMatch
assert.ElementsMatch(t, expectedIDs, actualIDs)
```

### 数值范围验证

```go
// ✅ 使用语义化断言
assert.GreaterOrEqual(t, count, 0, "计数不应为负")
assert.Positive(t, total, "总数应为正数")
```

## Helper 函数规范

### 登录辅助

```go
// ✅ 推荐：使用 manualtest 封装
c := manualtest.LoginAsAdmin(t)

// ❌ 避免：重复的登录代码
c := manualtest.NewClient()
_, err := c.Login("admin", "admin123")
require.NoError(t, err)
```

### 资源工厂

```go
// ✅ 推荐：使用工厂函数（自动清理）
user := manualtest.CreateTestUser(t, c, "testprefix")

// ❌ 避免：手动创建 + 手动 Cleanup
```

### 可用的 Helper 函数

| 函数                                         | 说明                           |
| -------------------------------------------- | ------------------------------ |
| `LoginAsAdmin(t) *Client`                    | 登录管理员，返回已认证客户端   |
| `LoginAs(t, account, password) *Client`      | 指定账户登录                   |
| `CreateTestUser(t, c, prefix)`               | 创建测试用户，自动注册 Cleanup |
| `CreateTestRole(t, c, prefix)`               | 创建测试角色，自动注册 Cleanup |
| `ExtractIDs(items, getter) []uint`           | 从结构体切片提取 ID            |
| `Get[T](c, path, query) (*T, error)`         | HTTP GET 请求                  |
| `GetList[T](c, path, query) ([]T, *, error)` | HTTP GET 列表请求              |
| `Post[T](c, path, body) (*T, error)`         | HTTP POST 请求                 |
| `Put[T](c, path, body) (*T, error)`          | HTTP PUT 请求                  |
