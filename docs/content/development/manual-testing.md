# 手动测试框架

`internal/manualtest/` 提供针对真实服务端的 API 集成测试框架，用于验证完整的请求-响应流程。

<!--TOC-->

## Table of Contents

- [概述](#概述) `:28+18`
- [快速开始](#快速开始) `:46+34`
  - [1. 启动服务端](#1-启动服务端) `:48+10`
  - [2. 运行测试](#2-运行测试) `:58+13`
  - [3. 使用自定义配置](#3-使用自定义配置) `:71+9`
- [环境变量](#环境变量) `:80+8`
- [Helper API 参考](#helper-api-参考) `:88+84`
  - [测试控制](#测试控制) `:90+7`
  - [客户端创建](#客户端创建) `:97+7`
  - [Client 方法](#client-方法) `:104+40`
  - [响应类型](#响应类型) `:144+28`
- [编写测试用例](#编写测试用例) `:172+72`
  - [基本模板](#基本模板) `:174+29`
  - [处理包装响应](#处理包装响应) `:203+17`
  - [资源清理](#资源清理) `:220+24`
- [最佳实践](#最佳实践) `:244+32`

<!--TOC-->

## 概述

与单元测试不同，手动测试需要真实的服务端运行。通过 `MANUAL=1` 环境变量控制测试执行，避免在 CI 环境意外运行。

**目录结构**：

```
internal/manualtest/
├── helper/           # 测试辅助库
│   ├── helper.go     # 环境检查、客户端工厂
│   └── client.go     # HTTP 客户端、泛型请求方法
├── auth_test.go      # 认证相关测试
├── profile_test.go   # 用户资料测试
├── roles_test.go     # 角色权限测试
├── twofa_test.go     # 2FA 测试
└── users_test.go     # 用户管理测试
```

## 快速开始

### 1. 启动服务端

```bash
# 使用 air 热重载
air

# 或直接运行
go run ./cmd/server api
```

### 2. 运行测试

```bash
# 运行所有手动测试
MANUAL=1 go test -v ./internal/manualtest/...

# 运行特定测试
MANUAL=1 go test -v -run TestLoginSuccess ./internal/manualtest/

# 运行某一类测试
MANUAL=1 go test -v -run "TestRole.*" ./internal/manualtest/
```

### 3. 使用自定义配置

```bash
# 指定 API 地址和开发密钥
API_BASE_URL=http://localhost:8080 \
DEV_SECRET=your-dev-secret \
MANUAL=1 go test -v ./internal/manualtest/...
```

## 环境变量

| 变量           | 默认值                   | 说明                  |
| -------------- | ------------------------ | --------------------- |
| `MANUAL`       | (空)                     | 设为 `1` 启用手动测试 |
| `API_BASE_URL` | `http://localhost:40012` | API 服务地址          |
| `DEV_SECRET`   | `dev-secret-change-me`   | 开发模式验证码密钥    |

## Helper API 参考

### 测试控制

```go
// SkipIfNotManual 在非手动模式下跳过测试
func SkipIfNotManual(t *testing.T)
```

### 客户端创建

```go
// NewClient 创建 HTTP 测试客户端（从环境变量读取配置）
func NewClient() *Client
```

### Client 方法

#### 认证相关

```go
// GetCaptcha 获取验证码（开发模式自动填充）
func (c *Client) GetCaptcha() (*captcha.CaptchaResponse, error)

// Login 登录并自动设置 Token
func (c *Client) Login(account, password string) (*auth.LoginResponse, error)

// LoginWithCaptcha 使用指定验证码登录
func (c *Client) LoginWithCaptcha(req auth.LoginDTO) (*auth.LoginResponse, error)

// SetToken 手动设置访问令牌
func (c *Client) SetToken(token string)
```

#### 泛型请求方法

```go
// Get 发送 GET 请求，自动解析 {"data": T} 响应
func Get[T any](c *Client, path string, queryParams map[string]string) (*T, error)

// GetList 发送 GET 请求，解析分页响应 {"data": []T, "meta": {...}}
func GetList[T any](c *Client, path string, queryParams map[string]string) ([]T, *response.PaginationMeta, error)

// Post 发送 POST 请求
func Post[T any](c *Client, path string, body any) (*T, error)

// Put 发送 PUT 请求
func Put[T any](c *Client, path string, body any) (*T, error)

// Delete 发送 DELETE 请求
func (c *Client) Delete(path string) error

// R 返回底层 resty.Request，用于自定义请求
func (c *Client) R() *resty.Request
```

### 响应类型

API 响应遵循统一格式，泛型方法自动解析：

```go
// 单对象响应
type DataResponse[T any] struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Data    T      `json:"data"`
}

// 分页响应
type PagedResponse[T any] struct {
    Code    int             `json:"code"`
    Message string          `json:"message"`
    Data    []T             `json:"data"`
    Meta    *PaginationMeta `json:"meta"`
}

type PaginationMeta struct {
    Total      int64 `json:"total"`
    Page       int   `json:"page"`
    PageSize   int   `json:"page_size"`
    TotalPages int   `json:"total_pages"`
}
```

## 编写测试用例

### 基本模板

```go
func TestXxxFlow(t *testing.T) {
    helper.SkipIfNotManual(t)

    c := helper.NewClient()

    // 1. 登录
    _, err := c.Login("admin", "admin123")
    if err != nil {
        t.Fatalf("登录失败: %v", err)
    }

    // 2. 执行测试操作
    result, err := helper.Get[YourType](c, "/api/xxx", nil)
    if err != nil {
        t.Fatalf("请求失败: %v", err)
    }

    // 3. 验证结果
    if result.ID == 0 {
        t.Fatal("期望获取有效 ID")
    }

    t.Logf("测试成功: %+v", result)
}
```

### 处理包装响应

部分 API 返回嵌套结构，需定义包装类型：

```go
// API 返回 {"data": {"user": {...}}}
type createUserResponse struct {
    User user.UserWithRolesResponse `json:"user"`
}

func TestCreateUser(t *testing.T) {
    // ...
    resp, err := helper.Post[createUserResponse](c, "/api/admin/users", req)
    userID := resp.User.ID  // 从包装中取值
}
```

### 资源清理

创建测试资源后应清理：

```go
func TestWithCleanup(t *testing.T) {
    helper.SkipIfNotManual(t)

    c := helper.NewClient()
    c.Login("admin", "admin123")

    // 创建资源
    resp, _ := helper.Post[createResp](c, "/api/admin/users", req)
    testUserID := resp.User.ID

    // 使用 t.Cleanup 确保清理执行
    t.Cleanup(func() {
        c.Delete(fmt.Sprintf("/api/admin/users/%d", testUserID))
    })

    // 执行测试...
}
```

## 最佳实践

1. **独立测试数据**：每个测试创建独立用户/资源，使用时间戳确保唯一性

   ```go
   username := fmt.Sprintf("test_%d", time.Now().Unix())
   ```

2. **权限隔离**：测试普通用户功能时，创建新用户并分配适当角色

   ```go
   createReq := user.CreateUserDTO{
       Username: testUsername,
       RoleIDs:  []uint{2}, // user 角色
   }
   ```

3. **详细日志**：使用 `t.Log()` 记录关键步骤，便于调试

   ```go
   t.Log("步骤 1: 创建测试用户")
   t.Logf("  用户 ID: %d", resp.User.ID)
   ```

4. **错误断言**：测试错误场景时验证预期行为
   ```go
   _, err := c.Login("admin", "wrong-password")
   if err == nil {
       t.Fatal("期望登录失败")
   }
   t.Logf("预期的错误: %v", err)
   ```
