# 测试指南

本文档介绍项目的测试策略和最佳实践。

<!--TOC-->

## Table of Contents

- [测试框架](#测试框架) `:29+7`
- [测试类型](#测试类型) `:36+257`
  - [单元测试](#单元测试) `:38+34`
  - [Use Case Handler 测试](#use-case-handler-测试) `:72+105`
  - [Repository 测试](#repository-测试) `:177+58`
  - [HTTP Handler 测试](#http-handler-测试) `:235+58`
- [运行测试](#运行测试) `:293+32`
  - [命令行](#命令行) `:295+23`
  - [使用 Task](#使用-task) `:318+7`
- [测试规范](#测试规范) `:325+88`
  - [1. 使用 AAA 模式](#1-使用-aaa-模式) `:327+15`
  - [2. 表驱动测试](#2-表驱动测试) `:342+28`
  - [3. 断言选择](#3-断言选择) `:370+25`
  - [4. Mock 最佳实践](#4-mock-最佳实践) `:395+18`
- [测试覆盖率目标](#测试覆盖率目标) `:413+9`
- [测试文件组织](#测试文件组织) `:422+24`
- [相关文档](#相关文档) `:446+5`

<!--TOC-->

## 测试框架

项目使用以下测试工具：

- **testify**: 断言库和 Mock 框架
- **go test**: Go 标准测试工具

## 测试类型

### 单元测试

测试单个函数或方法的正确性：

```go
// internal/domain/user/entity_user_test.go
package user

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestUser_CanLogin(t *testing.T) {
    tests := []struct {
        name     string
        status   string
        expected bool
    }{
        {"active user can login", "active", true},
        {"banned user cannot login", "banned", false},
        {"inactive user cannot login", "inactive", false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            user := &User{Status: tt.status}
            assert.Equal(t, tt.expected, user.CanLogin())
        })
    }
}
```

### Use Case Handler 测试

测试应用层业务逻辑，使用 Mock Repository：

```go
// internal/application/user/command/create_user_handler_test.go
package command

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "your-project/internal/domain/user"
)

// MockUserCommandRepository Mock 实现
type MockUserCommandRepository struct {
    mock.Mock
}

func (m *MockUserCommandRepository) Create(ctx context.Context, u *user.User) error {
    args := m.Called(ctx, u)
    return args.Error(0)
}

func (m *MockUserCommandRepository) Update(ctx context.Context, u *user.User) error {
    args := m.Called(ctx, u)
    return args.Error(0)
}

func (m *MockUserCommandRepository) Delete(ctx context.Context, id uint) error {
    args := m.Called(ctx, id)
    return args.Error(0)
}

// MockUserQueryRepository Mock 实现
type MockUserQueryRepository struct {
    mock.Mock
}

func (m *MockUserQueryRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
    args := m.Called(ctx, username)
    return args.Bool(0), args.Error(1)
}

func (m *MockUserQueryRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
    args := m.Called(ctx, email)
    return args.Bool(0), args.Error(1)
}

func TestCreateUserHandler_Handle(t *testing.T) {
    t.Run("successfully creates user", func(t *testing.T) {
        // Arrange
        mockCommandRepo := new(MockUserCommandRepository)
        mockQueryRepo := new(MockUserQueryRepository)
        mockAuthService := new(MockAuthService)

        handler := NewCreateUserHandler(mockCommandRepo, mockQueryRepo, mockAuthService)

        mockQueryRepo.On("ExistsByUsername", mock.Anything, "testuser").Return(false, nil)
        mockQueryRepo.On("ExistsByEmail", mock.Anything, "test@example.com").Return(false, nil)
        mockAuthService.On("GeneratePasswordHash", mock.Anything, "password123").Return("hashed", nil)
        mockCommandRepo.On("Create", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)

        // Act
        result, err := handler.Handle(context.Background(), CreateUserCommand{
            Username: "testuser",
            Email:    "test@example.com",
            Password: "password123",
        })

        // Assert
        assert.NoError(t, err)
        assert.NotNil(t, result)
        mockCommandRepo.AssertExpectations(t)
        mockQueryRepo.AssertExpectations(t)
    })

    t.Run("fails when username exists", func(t *testing.T) {
        // Arrange
        mockCommandRepo := new(MockUserCommandRepository)
        mockQueryRepo := new(MockUserQueryRepository)
        mockAuthService := new(MockAuthService)

        handler := NewCreateUserHandler(mockCommandRepo, mockQueryRepo, mockAuthService)

        mockQueryRepo.On("ExistsByUsername", mock.Anything, "existinguser").Return(true, nil)

        // Act
        result, err := handler.Handle(context.Background(), CreateUserCommand{
            Username: "existinguser",
            Email:    "test@example.com",
            Password: "password123",
        })

        // Assert
        assert.Error(t, err)
        assert.Nil(t, result)
        assert.Contains(t, err.Error(), "username already exists")
    })
}
```

### Repository 测试

测试数据库访问层，使用真实数据库连接：

```go
// internal/infrastructure/persistence/user_repository_test.go
package persistence

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "your-project/internal/domain/user"
    "your-project/internal/infrastructure/database"
)

func setupTestDB(t *testing.T) *gorm.DB {
    db, err := database.NewTestConnection()
    require.NoError(t, err)

    // 自动迁移
    err = db.AutoMigrate(&UserModel{})
    require.NoError(t, err)

    // 测试结束清理数据
    t.Cleanup(func() {
        db.Exec("DELETE FROM users")
    })

    return db
}

func TestUserCommandRepository_Create(t *testing.T) {
    db := setupTestDB(t)
    repo := NewUserCommandRepository(db)

    ctx := context.Background()
    testUser := &user.User{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "hashedpassword",
    }

    // 执行创建
    err := repo.Create(ctx, testUser)
    assert.NoError(t, err)
    assert.NotZero(t, testUser.ID)

    // 验证数据库中的数据
    var model UserModel
    err = db.First(&model, testUser.ID).Error
    assert.NoError(t, err)
    assert.Equal(t, "testuser", model.Username)
}
```

### HTTP Handler 测试

测试 HTTP 请求和响应处理：

```go
// internal/adapters/http/handler/user_test.go
package handler

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestUserHandler_Create(t *testing.T) {
    gin.SetMode(gin.TestMode)

    t.Run("successfully creates user", func(t *testing.T) {
        // Setup
        mockHandler := new(MockCreateUserHandler)
        mockHandler.On("Handle", mock.Anything, mock.AnythingOfType("CreateUserCommand")).
            Return(&CreateUserResult{UserID: 1, Username: "testuser", Email: "test@example.com"}, nil)

        handler := &UserHandler{createUserHandler: mockHandler}

        router := gin.New()
        router.POST("/users", handler.Create)

        // Request
        body := map[string]interface{}{
            "username": "testuser",
            "email":    "test@example.com",
            "password": "password123",
        }
        jsonBody, _ := json.Marshal(body)

        req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonBody))
        req.Header.Set("Content-Type", "application/json")

        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)

        // Assertions
        assert.Equal(t, http.StatusCreated, w.Code)

        var response map[string]interface{}
        json.Unmarshal(w.Body.Bytes(), &response)
        assert.Equal(t, "user created successfully", response["message"])
    })
}
```

## 运行测试

### 命令行

```bash
# 运行所有测试
go test ./...

# 运行特定目录测试
go test ./internal/domain/user/...

# 详细输出
go test -v ./...

# 运行特定测试函数
go test -run TestUser_CanLogin ./internal/domain/user/

# 覆盖率统计
go test -cover ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### 使用 Task

```bash
task test           # 运行所有测试
task test:coverage  # 生成覆盖率报告
```

## 测试规范

### 1. 使用 AAA 模式

```go
func TestSomething(t *testing.T) {
    // Arrange - 准备测试数据
    user := &User{Status: "active"}

    // Act - 执行被测试操作
    result := user.CanLogin()

    // Assert - 验证结果
    assert.True(t, result)
}
```

### 2. 表驱动测试

```go
func TestPasswordValidation(t *testing.T) {
    tests := []struct {
        name     string
        password string
        wantErr  bool
    }{
        {"valid password", "Password123!", false},
        {"too short", "Pass1!", true},
        {"no uppercase", "password123!", true},
        {"no number", "Password!", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidatePassword(tt.password)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### 3. 断言选择

| 方法      | 用途               | 使用场景               |
| --------- | ------------------ | ---------------------- |
| `assert`  | 失败时**继续**执行 | 验证多个条件时使用     |
| `require` | 失败时**立即终止** | 前置条件失败应终止测试 |

```go
func TestUserCreation(t *testing.T) {
    // 前置条件使用 require
    db := setupTestDB(t)
    require.NotNil(t, db, "database connection is required")

    repo := NewUserRepository(db)
    require.NotNil(t, repo)

    // 验证结果使用 assert
    user := &User{Username: "test"}
    err := repo.Create(context.Background(), user)
    assert.NoError(t, err)
    assert.NotZero(t, user.ID)
    assert.Equal(t, "test", user.Username)
}
```

### 4. Mock 最佳实践

```go
// 定义清晰的 Mock 实现
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *User) error {
    args := m.Called(ctx, user)
    return args.Error(0)
}

// 在测试中使用
mockRepo := new(MockUserRepository)
mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*User")).Return(nil)
```

## 测试覆盖率目标

| 层级              | 覆盖率目标 |
| ----------------- | ---------- |
| Domain 层         | 90%+       |
| Application 层    | 80%+       |
| Infrastructure 层 | 70%+       |
| Adapters 层       | 60%+       |

## 测试文件组织

```
internal/
├── domain/
│   └── user/
│       ├── entity_user.go
│       └── entity_user_test.go      # 单元测试
├── application/
│   └── user/
│       └── command/
│           ├── create_user_handler.go
│           └── create_user_handler_test.go  # Use Case 测试
├── infrastructure/
│   └── persistence/
│       ├── user_command_repository.go
│       └── user_repository_test.go  # 集成测试
└── adapters/
    └── http/
        └── handler/
            ├── user.go
            └── user_test.go         # HTTP 测试
```

## 相关文档

- [错误处理](/development/error-handling)
- [日志记录](/development/logging)
- [代码规范](/development/code-style)
