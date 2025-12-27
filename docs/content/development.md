# 开发指南

本文档包含开发 Go DDD Template 所需的技术文档。

<!--TOC-->

## Table of Contents

- [环境准备](#环境准备) `:30+25`
  - [系统要求](#系统要求) `:32+7`
  - [安装开发工具](#安装开发工具) `:39+16`
- [项目设置](#项目设置) `:55+23`
- [开发工作流](#开发工作流) `:78+16`
- [测试策略](#测试策略) `:94+34`
  - [测试框架](#测试框架) `:96+5`
  - [测试类型](#测试类型) `:101+9`
  - [运行测试](#运行测试) `:110+9`
  - [测试覆盖率目标](#测试覆盖率目标) `:119+9`
- [手动测试](#手动测试) `:128+36`
  - [运行手动测试](#运行手动测试) `:132+13`
  - [环境变量](#环境变量) `:145+8`
  - [Helper API](#helper-api) `:153+11`
- [故障排查](#故障排查) `:164+22`
  - [端口占用](#端口占用) `:166+7`
  - [依赖问题](#依赖问题) `:173+7`
  - [Docker 问题](#docker-问题) `:180+6`

<!--TOC-->

## 环境准备

### 系统要求

- **Go** 1.25.4+
- **Docker** 20.10+ & Docker Compose 2.0+
- **Task** - 任务自动化工具
- **Air** - Go 热重载工具（可选）

### 安装开发工具

```bash
# Task (任务运行器)
go install github.com/go-task/task/v3/cmd/task@latest

# Air (热重载)
go install github.com/air-verse/air@latest

# golangci-lint v2 (代码检查，官方推荐二进制安装)
curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(go env GOPATH)/bin latest

# swag (Swagger 文档生成)
go install github.com/swaggo/swag/cmd/swag@latest
```

## 项目设置

```bash
# 1. 克隆代码
git clone https://github.com/lwmacct/251117-go-ddd-template.git
cd 251117-go-ddd-template

# 2. 安装依赖
go mod download
cd web && npm install && cd ..
cd docs && npm install && cd ..

# 3. 复制配置文件
cp configs/config.example.yaml configs/config.yaml

# 4. 启动基础服务
docker-compose up -d

# 5. 初始化数据库
task go:run -- migrate up
task go:run -- seed
```

## 开发工作流

```bash
# 启动后端（热重载）
air

# 启动前端
cd web && npm run dev

# 启动文档服务
cd docs && npm run dev

# 代码检查
golangci-lint run
```

## 测试策略

### 测试框架

- **testify**: 断言库和 Mock 框架
- **go test**: Go 标准测试工具

### 测试类型

| 类型            | 位置                                   | 说明                         |
| --------------- | -------------------------------------- | ---------------------------- |
| 单元测试        | `*_test.go` 同目录                     | 测试单个函数/方法            |
| Use Case 测试   | `application/*/command/*_test.go`      | Mock Repository 测试业务逻辑 |
| Repository 测试 | `infrastructure/persistence/*_test.go` | 真实数据库集成测试           |
| HTTP 测试       | `adapters/http/handler/*_test.go`      | Mock Handler 测试 HTTP 层    |

### 运行测试

```bash
go test ./...                              # 运行所有测试
go test -v ./...                           # 详细输出
go test ./internal/domain/user/...         # 特定目录
go test -coverprofile=coverage.out ./...   # 覆盖率报告
```

### 测试覆盖率目标

| 层级              | 目标 |
| ----------------- | ---- |
| Domain 层         | 90%+ |
| Application 层    | 80%+ |
| Infrastructure 层 | 70%+ |
| Adapters 层       | 60%+ |

## 手动测试

`internal/manualtest/` 提供针对真实服务端的 API 集成测试框架。

### 运行手动测试

```bash
# 1. 启动服务端
air

# 2. 运行测试
MANUAL=1 go test -v ./internal/manualtest/...

# 3. 运行特定测试
MANUAL=1 go test -v -run "TestRole.*" ./internal/manualtest/
```

### 环境变量

| 变量           | 默认值                   | 说明                  |
| -------------- | ------------------------ | --------------------- |
| `MANUAL`       | (空)                     | 设为 `1` 启用手动测试 |
| `API_BASE_URL` | `http://localhost:40012` | API 服务地址          |
| `DEV_SECRET`   | `dev-secret-change-me`   | 开发模式验证码密钥    |

### Helper API

**参考实现**: `internal/manualtest/helper/`

| 方法                 | 说明           |
| -------------------- | -------------- |
| `SkipIfNotManual(t)` | 非手动模式跳过 |
| `NewClient()`        | 创建测试客户端 |
| `Get[T]()`           | GET 请求       |
| `Post[T]()`          | POST 请求      |

## 故障排查

### 端口占用

```bash
lsof -i :8080
# 更改端口: APP_SERVER_ADDR=":9000" air
```

### 依赖问题

```bash
go clean -modcache
go mod download
```

### Docker 问题

```bash
docker-compose down -v
docker-compose up -d
```
