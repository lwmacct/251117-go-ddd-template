---
paths:
  - "internal/manualtest/**/*.go"
---

# manualtest 手动测试规范

本包提供针对 HTTP API 的集成测试，需要服务运行时手动执行。

<!--TOC-->

## Table of Contents

- [运行方式](#运行方式) `:24+12`
- [DTO 使用原则](#dto-使用原则) `:36+44`
  - [正确做法](#正确做法) `:40+17`
  - [禁止做法](#禁止做法) `:57+8`
  - [类型来源](#类型来源) `:65+8`
  - [设计原因](#设计原因) `:73+7`
- [设计反思原则](#设计反思原则) `:80+20`

<!--TOC-->

## 运行方式

> **Note**: 服务端已使用 `air` 热重载运行，无需手动启动服务。

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
// ❌ 禁止在 manualtest 或 helper 包中定义任何 DTO
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

1. **单一职责** - DTO 定义属于 Application 层，测试代码只负责验证行为
2. **避免重复** - Application 层的 DTO 已有完整的 JSON tags，无需重复定义
3. **保持同步** - 直接使用 Application DTO 确保测试与实际 API 响应格式一致
4. **依赖方向** - `manualtest → application` 符合 DDD 依赖方向

## 设计反思原则

**测试困难是设计问题的信号。**

如果发现以下情况，说明 Application 层设计需要检视：

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
