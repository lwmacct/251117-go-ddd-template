# manualtest 手动测试包

本包提供针对 HTTP API 的集成测试，需要服务运行时手动执行。

<!--TOC-->

## Table of Contents

- [运行方式](#运行方式) `:19+10`
- [DTO 使用原则](#dto-使用原则) `:29+53`
  - [正确做法](#正确做法) `:33+17`
  - [禁止做法](#禁止做法) `:50+8`
  - [类型来源](#类型来源) `:58+8`
  - [设计原因](#设计原因) `:66+7`
  - [当需要新 DTO 时](#当需要新-dto-时) `:73+9`

<!--TOC-->

## 运行方式

```bash
# 运行所有测试
MANUAL=1 go test -v ./internal/manualtest/

# 运行单个测试
MANUAL=1 go test -v -run TestLoginSuccess ./internal/manualtest/
```

## DTO 使用原则

**核心原则：manualtest 不定义任何 DTO，只消费 Application 层的类型。**

### 正确做法

```go
import (
    "github.com/lwmacct/251117-go-ddd-template/internal/application/auth"
    "github.com/lwmacct/251117-go-ddd-template/internal/application/user"
)

// 使用 Application 层 DTO 解析响应
result, err := helper.Post[auth.LoginResponseDTO](c, "/api/auth/login", req)
profile, err := helper.Get[user.UserWithRolesDTO](c, "/api/user/profile", nil)

// 创建用户后直接获取 DTO
createResp, err := helper.Post[user.UserWithRolesDTO](c, "/api/admin/users", req)
userID := createResp.ID  // 直接访问字段
```

### 禁止做法

```go
// 禁止在 manualtest 或 helper 包中定义任何 DTO
type LoginResponse struct { ... }  // 禁止
type PATTokenDTO struct { ... }    // 禁止
```

### 类型来源

| 用途          | 来源                                          |
| ------------- | --------------------------------------------- |
| HTTP 响应解析 | `internal/application/*/dto.go`               |
| HTTP 请求构造 | `internal/application/*/dto.go`               |
| 通用响应包装  | `internal/adapters/http/response/response.go` |

### 设计原因

1. **单一职责**：DTO 定义属于 Application 层，测试代码只负责验证行为
2. **避免重复**：Application 层的 DTO 已有完整的 JSON tags，无需在测试中重复定义
3. **保持同步**：直接使用 Application DTO 可确保测试与实际 API 响应格式一致
4. **依赖方向**：`manualtest → application` 符合 DDD 依赖方向

### 当需要新 DTO 时

**如果发现 manualtest 需要定义 DTO，说明 Application 层设计有问题，应该：**

1. 检查 Application 层 DTO 是否缺少 JSON tags
2. 检查 Handler 响应格式是否与 DTO 结构不匹配（所有 API 统一使用 `response` 包封装）
3. 在 Application 层补充缺失的 DTO

**禁止在 manualtest 中临时定义 DTO 来"修复"问题。**
