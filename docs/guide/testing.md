# 测试指南

本文档说明如何为 Go DDD Template 项目编写单元测试和集成测试。

## 测试框架选择

本项目推荐使用以下测试工具：

| 工具                                             | 用途              | 安装                                             |
| ------------------------------------------------ | ----------------- | ------------------------------------------------ |
| Go 标准库 `testing`                              | 基础测试框架      | 内置                                             |
| [Testify](https://github.com/stretchr/testify)   | 断言和 Mock       | `go get github.com/stretchr/testify`             |
| [Mockery](https://github.com/vektra/mockery)     | Mock 生成         | `go install github.com/vektra/mockery/v2@latest` |
| [Dockertest](https://github.com/ory/dockertest)  | 集成测试 (Docker) | `go get github.com/ory/dockertest/v3`            |
| [httptest](https://pkg.go.dev/net/http/httptest) | HTTP 测试         | 内置                                             |

## 安装依赖

```bash
# 安装 Testify
go get github.com/stretchr/testify/assert
go get github.com/stretchr/testify/mock
go get github.com/stretchr/testify/suite

# 安装 Dockertest (集成测试)
go get github.com/ory/dockertest/v3

# 安装 Mockery (可选，用于生成 Mock)
go install github.com/vektra/mockery/v2@latest
```

## 测试目录结构

```
.
├── internal/
│   ├── domain/
│   │   └── user/
│   │       ├── entity_user.go
│   │       ├── entity_user_test.go    # 单元测试
│   │       ├── command_repository.go
│   │       ├── query_repository.go
│   │       └── repository_mock.go     # Command/Query Mock（通过 mockery 生成）
│   ├── infrastructure/
│   │   └── persistence/
│   │       ├── user_command_repository.go      # 写实现
│   │       ├── user_command_repository_test.go # 写仓储集成测试
│   │       ├── user_query_repository.go        # 读实现
│   │       └── user_query_repository_test.go   # 读仓储集成测试
│   └── adapters/
│       └── http/
│           └── handler/
│               ├── user_handler.go
│               └── user_handler_test.go # HTTP 测试
└── tests/
    ├── integration/                    # 端到端集成测试
    │   └── api_test.go
    └── testdata/                       # 测试数据
        └── fixtures.json
```

## 单元测试

### 1. 领域模型测试

测试领域模型的业务逻辑：

```go
// internal/domain/user/model_test.go
package user_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
)

func TestUser_Validate(t *testing.T) {
	tests := []struct {
		name    string
		user    *user.User
		wantErr bool
		errMsg  string
	}{
		{
			name: "有效用户",
			user: &user.User{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "无效邮箱",
			user: &user.User{
				Username: "testuser",
				Email:    "invalid-email",
				Password: "password123",
			},
			wantErr: true,
			errMsg:  "invalid email format",
		},
		{
			name: "用户名为空",
			user: &user.User{
				Username: "",
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: true,
			errMsg:  "username cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUser_HashPassword(t *testing.T) {
	u := &user.User{
		Password: "plainpassword",
	}

	err := u.HashPassword()
	assert.NoError(t, err)
	assert.NotEqual(t, "plainpassword", u.Password)
	assert.NotEmpty(t, u.Password)
}

func TestUser_ComparePassword(t *testing.T) {
	u := &user.User{
		Password: "plainpassword",
	}
	u.HashPassword()

	// 正确密码
	assert.True(t, u.ComparePassword("plainpassword"))

	// 错误密码
	assert.False(t, u.ComparePassword("wrongpassword"))
}
```

### 2. 使用 Mock 测试

创建 Mock 接口：

```go
// internal/domain/user/repository_mock.go
package user

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, user *User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockRepository) FindByID(ctx context.Context, id uint) (*User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockRepository) FindByUsername(ctx context.Context, username string) (*User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, user *User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
```

使用 Mock 测试服务层：

```go
// internal/application/user_service_test.go
package application_test

import (
	"context"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
	"github.com/lwmacct/251117-go-ddd-template/internal/application"
)

func TestUserService_CreateUser(t *testing.T) {
	// 准备 Mock
	mockRepo := new(user.MockRepository)
	service := application.NewUserService(mockRepo)

	newUser := &user.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	// 设置 Mock 期望
	mockRepo.On("FindByUsername", mock.Anything, "testuser").Return(nil, nil)
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)

	// 执行测试
	err := service.CreateUser(context.Background(), newUser)

	// 断言
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockRepo.AssertCalled(t, "FindByUsername", mock.Anything, "testuser")
	mockRepo.AssertCalled(t, "Create", mock.Anything, mock.AnythingOfType("*user.User"))
}

func TestUserService_CreateUser_DuplicateUsername(t *testing.T) {
	mockRepo := new(user.MockRepository)
	service := application.NewUserService(mockRepo)

	newUser := &user.User{
		Username: "existinguser",
		Email:    "test@example.com",
		Password: "password123",
	}

	existingUser := &user.User{
		ID:       1,
		Username: "existinguser",
	}

	// 模拟用户已存在
	mockRepo.On("FindByUsername", mock.Anything, "existinguser").Return(existingUser, nil)

	// 执行测试
	err := service.CreateUser(context.Background(), newUser)

	// 断言
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "username already exists")
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}
```

### 3. HTTP Handler 测试

测试 HTTP 处理器：

```go
// internal/adapters/http/handler/user_handler_test.go
package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/handler"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

func TestUserHandler_Register(t *testing.T) {
	router := setupTestRouter()
	mockRepo := new(user.MockRepository)
	h := handler.NewUserHandler(mockRepo)

	router.POST("/api/auth/register", h.Register)

	tests := []struct {
		name       string
		payload    map[string]interface{}
		setupMock  func()
		wantStatus int
		wantBody   map[string]interface{}
	}{
		{
			name: "成功注册",
			payload: map[string]interface{}{
				"username": "newuser",
				"email":    "new@example.com",
				"password": "password123",
			},
			setupMock: func() {
				mockRepo.On("FindByUsername", mock.Anything, "newuser").Return(nil, nil).Once()
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil).Once()
			},
			wantStatus: http.StatusCreated,
			wantBody: map[string]interface{}{
				"success": true,
			},
		},
		{
			name: "用户名已存在",
			payload: map[string]interface{}{
				"username": "existinguser",
				"email":    "test@example.com",
				"password": "password123",
			},
			setupMock: func() {
				existingUser := &user.User{ID: 1, Username: "existinguser"}
				mockRepo.On("FindByUsername", mock.Anything, "existinguser").Return(existingUser, nil).Once()
			},
			wantStatus: http.StatusConflict,
			wantBody: map[string]interface{}{
				"success": false,
				"error":   "username already exists",
			},
		},
		{
			name: "无效的邮箱",
			payload: map[string]interface{}{
				"username": "testuser",
				"email":    "invalid-email",
				"password": "password123",
			},
			wantStatus: http.StatusBadRequest,
			wantBody: map[string]interface{}{
				"success": false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
			}

			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			for key, value := range tt.wantBody {
				assert.Equal(t, value, response[key])
			}
		})
	}
}
```

## 集成测试

### 1. 数据库集成测试

使用 Dockertest 进行数据库集成测试：

```go
// internal/infrastructure/persistence/user_repository_test.go
package persistence_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/persistence"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	db       *gorm.DB
	pool     *dockertest.Pool
	resource *dockertest.Resource
	commandRepo user.CommandRepository
	queryRepo   user.QueryRepository
}

func (suite *UserRepositoryTestSuite) SetupSuite() {
	var err error

	// 创建 Docker 连接池
	suite.pool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// 启动 PostgreSQL 容器
	suite.resource, err = suite.pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "17-alpine",
		Env: []string{
			"POSTGRES_PASSWORD=password",
			"POSTGRES_USER=postgres",
			"POSTGRES_DB=testdb",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// 设置过期时间
	suite.resource.Expire(120)

	// 等待数据库就绪
	dsn := fmt.Sprintf("host=localhost port=%s user=postgres password=password dbname=testdb sslmode=disable",
		suite.resource.GetPort("5432/tcp"))

	if err := suite.pool.Retry(func() error {
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			return err
		}
		suite.db = db
		return nil
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	// 自动迁移
	suite.db.AutoMigrate(&user.User{})

	// 创建 CQRS 仓储实例
	suite.commandRepo = persistence.NewUserCommandRepository(suite.db)
	suite.queryRepo = persistence.NewUserQueryRepository(suite.db)
}

func (suite *UserRepositoryTestSuite) TearDownSuite() {
	if err := suite.pool.Purge(suite.resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}

func (suite *UserRepositoryTestSuite) SetupTest() {
	// 每个测试前清空表
	suite.db.Exec("DELETE FROM users")
}

func (suite *UserRepositoryTestSuite) TestCreate() {
	ctx := context.Background()
	newUser := &user.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}

	err := suite.commandRepo.Create(ctx, newUser)
	assert.NoError(suite.T(), err)
	assert.NotZero(suite.T(), newUser.ID)
	assert.NotZero(suite.T(), newUser.CreatedAt)
}

func (suite *UserRepositoryTestSuite) TestGetByID() {
	ctx := context.Background()

	// 创建用户
	newUser := &user.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}
	suite.commandRepo.Create(ctx, newUser)

	// 查找用户
	foundUser, err := suite.queryRepo.GetByID(ctx, newUser.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), newUser.Username, foundUser.Username)
	assert.Equal(suite.T(), newUser.Email, foundUser.Email)
}

func (suite *UserRepositoryTestSuite) TestUpdate() {
	ctx := context.Background()

	// 创建用户
	newUser := &user.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}
	suite.commandRepo.Create(ctx, newUser)

	// 更新用户
	newUser.Email = "updated@example.com"
	err := suite.commandRepo.Update(ctx, newUser)
	assert.NoError(suite.T(), err)

	// 验证更新
	foundUser, _ := suite.queryRepo.GetByID(ctx, newUser.ID)
	assert.Equal(suite.T(), "updated@example.com", foundUser.Email)
}

func (suite *UserRepositoryTestSuite) TestDelete() {
	ctx := context.Background()

	// 创建用户
	newUser := &user.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}
	suite.commandRepo.Create(ctx, newUser)

	// 删除用户 (软删除)
	err := suite.commandRepo.Delete(ctx, newUser.ID)
	assert.NoError(suite.T(), err)

	// 验证已删除 (软删除，仍能查到 DeletedAt)
	foundUser, err := suite.queryRepo.GetByID(ctx, newUser.ID)
	assert.Error(suite.T(), err) // 应该找不到 (因为软删除)
}

func TestUserRepositoryTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}
	suite.Run(t, new(UserRepositoryTestSuite))
}
```

### 2. 端到端 API 测试

```go
// tests/integration/api_test.go
package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/lwmacct/251117-go-ddd-template/internal/bootstrap"
)

func setupTestApp(t *testing.T) *gin.Engine {
	// 设置测试配置
	// 初始化容器和路由
	container := bootstrap.NewContainer()
	router := setupRouter(container)
	return router
}

func TestAuthFlow(t *testing.T) {
	router := setupTestApp(t)

	// 1. 注册用户
	registerPayload := map[string]string{
		"username": "testuser",
		"email":    "test@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(registerPayload)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// 2. 登录获取 Token
	loginPayload := map[string]string{
		"username": "testuser",
		"password": "password123",
	}
	body, _ = json.Marshal(loginPayload)
	req = httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var loginResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	accessToken := loginResponse["data"].(map[string]interface{})["access_token"].(string)
	assert.NotEmpty(t, accessToken)

	// 3. 使用 Token 访问受保护端点
	req = httptest.NewRequest(http.MethodGet, "/api/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
```

## 运行测试

### 基本命令

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./internal/domain/user

# 查看覆盖率
go test -cover ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# 运行详细输出
go test -v ./...

# 跳过集成测试 (仅单元测试)
go test -short ./...

# 运行特定测试
go test -run TestUser_Validate ./internal/domain/user

# 并行测试
go test -parallel 4 ./...

# 测试超时
go test -timeout 30s ./...
```

### 使用 Task 运行测试

在 `Taskfile.yml` 中添加测试任务：

```yaml
# .taskfile/go.yml
tasks:
  test:
    desc: Run unit tests
    cmds:
      - go test -short -v ./...

  test:integration:
    desc: Run integration tests
    cmds:
      - go test -v ./...

  test:coverage:
    desc: Run tests with coverage
    cmds:
      - go test -coverprofile=coverage.out ./...
      - go tool cover -html=coverage.out -o coverage.html
      - echo "Coverage report generated: coverage.html"

  test:bench:
    desc: Run benchmarks
    cmds:
      - go test -bench=. -benchmem ./...
```

运行：

```bash
task test                # 单元测试
task test:integration    # 集成测试
task test:coverage       # 覆盖率测试
task test:bench          # 基准测试
```

## CI/CD 集成

### GitHub Actions

创建 `.github/workflows/test.yml`：

```yaml
name: Tests

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest

    services:
      postgres:
        image: postgres:17-alpine
        env:
          POSTGRES_PASSWORD: postgres
          POSTGRES_DB: testdb
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432

      redis:
        image: redis:7-alpine
        options: >-
          --health-cmd "redis-cli ping"
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 6379:6379

    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: "1.25"

      - name: Install dependencies
        run: go mod download

      - name: Run unit tests
        run: go test -short -v ./...

      - name: Run integration tests
        run: go test -v ./...
        env:
          APP_DATA_PGSQL_URL: postgresql://postgres:postgres@localhost:5432/testdb
          APP_DATA_REDIS_URL: redis://localhost:6379/0

      - name: Generate coverage
        run: |
          go test -coverprofile=coverage.out ./...
          go tool cover -func=coverage.out

      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./coverage.out
```

## 测试覆盖率目标

| 层级         | 目标覆盖率 | 说明                     |
| ------------ | ---------- | ------------------------ |
| 领域模型     | 90%+       | 核心业务逻辑必须充分测试 |
| 仓储层       | 80%+       | 数据访问层需要集成测试   |
| 应用服务     | 85%+       | 业务流程编排需要测试     |
| HTTP Handler | 75%+       | API 端点需要测试         |
| 整体项目     | 80%+       | 整体覆盖率目标           |

## 最佳实践

### 1. 测试命名规范

```go
func TestFunction_Scenario_ExpectedBehavior(t *testing.T) {}

// 示例
func TestUserService_CreateUser_Success(t *testing.T) {}
func TestUserService_CreateUser_DuplicateUsername(t *testing.T) {}
func TestUserService_CreateUser_InvalidEmail(t *testing.T) {}
```

### 2. 表驱动测试

优先使用表驱动测试提高可维护性：

```go
tests := []struct {
    name    string
    input   string
    want    string
    wantErr bool
}{
    // 测试用例
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        // 测试逻辑
    })
}
```

### 3. 测试隔离

- 每个测试应该独立运行
- 使用 `SetupTest` 和 `TearDownTest` 清理状态
- 不要依赖测试执行顺序

### 4. Mock 最佳实践

- 只 Mock 外部依赖 (数据库、API)
- 不要 Mock 领域模型
- 使用接口而不是具体实现

### 5. 集成测试提示

- 使用 Dockertest 启动真实服务
- 使用 `-short` 标志跳过集成测试
- 在 CI/CD 中运行完整集成测试

## 常见问题

### Q: 如何跳过耗时的集成测试？

```go
func TestIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }
    // 测试逻辑
}
```

运行：`go test -short ./...`

### Q: 如何测试私有函数？

不要测试私有函数，应该通过公有接口测试。如果确实需要，将测试文件放在同一包内 (不使用 `_test` 后缀) 。

### Q: 如何处理时间依赖？

使用时间接口：

```go
type TimeProvider interface {
    Now() time.Time
}

// 测试时注入 Mock 时间
type MockTimeProvider struct {
    CurrentTime time.Time
}

func (m *MockTimeProvider) Now() time.Time {
    return m.CurrentTime
}
```

## 相关链接

- [Go 测试官方文档](https://pkg.go.dev/testing)
- [Testify GitHub](https://github.com/stretchr/testify)
- [Dockertest GitHub](https://github.com/ory/dockertest)
- [贡献指南](/guide/contributing) - 提交代码规范
- [CI/CD 配置](/guide/application-deployment#cicd-集成)
