# 开发环境设置

本指南介绍如何设置 Go DDD Template 的开发环境。

<!--TOC-->

## Table of Contents

- [系统要求](#系统要求) `:42+16`
  - [必需软件](#必需软件) `:44+7`
  - [推荐软件](#推荐软件) `:51+7`
- [环境准备](#环境准备) `:58+47`
  - [1. 安装 Go](#1-安装-go) `:60+14`
  - [2. 安装 Docker](#2-安装-docker) `:74+15`
  - [3. 安装开发工具](#3-安装开发工具) `:89+16`
- [项目设置](#项目设置) `:105+72`
  - [1. 克隆代码](#1-克隆代码) `:107+7`
  - [2. 安装依赖](#2-安装依赖) `:114+13`
  - [3. 环境配置](#3-环境配置) `:127+30`
  - [4. 启动基础服务](#4-启动基础服务) `:157+10`
  - [5. 初始化数据库](#5-初始化数据库) `:167+10`
- [IDE 配置](#ide-配置) `:177+37`
  - [VSCode](#vscode) `:179+28`
  - [GoLand](#goland) `:207+7`
- [开发工作流](#开发工作流) `:214+60`
  - [1. 启动开发服务器](#1-启动开发服务器) `:216+13`
  - [2. 前端开发](#2-前端开发) `:229+8`
  - [3. 文档开发](#3-文档开发) `:237+8`
  - [4. 运行测试](#4-运行测试) `:245+16`
  - [5. 代码检查](#5-代码检查) `:261+13`
- [常用 Task 命令](#常用-task-命令) `:274+31`
- [调试配置](#调试配置) `:305+39`
  - [使用 Delve](#使用-delve) `:307+13`
  - [VSCode 调试配置](#vscode-调试配置) `:320+24`
- [故障排查](#故障排查) `:344+32`
  - [端口占用](#端口占用) `:346+10`
  - [依赖问题](#依赖问题) `:356+10`
  - [Docker 问题](#docker-问题) `:366+10`

<!--TOC-->

## 系统要求

### 必需软件

- **Go** 1.25.4+
- **Docker** 20.10+
- **Docker Compose** 2.0+
- **Git** 2.0+

### 推荐软件

- **Task** - 任务自动化工具
- **Air** - Go 热重载工具
- **Make** - 备用构建工具
- **VSCode** / **GoLand** - IDE

## 环境准备

### 1. 安装 Go

```bash
# macOS (使用 Homebrew)
brew install go

# Linux
wget https://go.dev/dl/go1.25.4.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.25.4.linux-amd64.tar.gz

# 验证安装
go version
```

### 2. 安装 Docker

```bash
# macOS
brew install --cask docker

# Linux (Ubuntu/Debian)
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker $USER

# 验证安装
docker --version
docker-compose --version
```

### 3. 安装开发工具

```bash
# Task (任务运行器)
go install github.com/go-task/task/v3/cmd/task@latest

# Air (热重载)
go install github.com/air-verse/air@latest

# golangci-lint (代码检查)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# swag (Swagger 文档生成)
go install github.com/swaggo/swag/cmd/swag@latest
```

## 项目设置

### 1. 克隆代码

```bash
git clone https://github.com/lwmacct/251117-go-ddd-template.git
cd 251117-go-ddd-template
```

### 2. 安装依赖

```bash
# Go 依赖
go mod download

# 前端依赖
cd web && npm install && cd ..

# 文档依赖
cd docs && npm install && cd ..
```

### 3. 环境配置

创建本地配置文件：

```bash
cp configs/config.example.yaml configs/config.yaml
```

编辑 `configs/config.yaml`：

```yaml
server:
  addr: ":8080"
  env: "development"

data:
  pgsql_url: "postgresql://postgres@localhost:5432/devdb?sslmode=disable"
  redis_url: "redis://localhost:6379/0"
  auto_migrate: true

jwt:
  secret: "dev-secret-change-me"
  access_token_expiry: "30m"
  refresh_token_expiry: "7d"

log:
  level: "debug"
  format: "text"
```

### 4. 启动基础服务

```bash
# 启动 PostgreSQL 和 Redis
docker-compose up -d

# 验证服务状态
docker-compose ps
```

### 5. 初始化数据库

```bash
# 运行迁移
task go:run -- migrate up

# 填充测试数据
task go:run -- seed
```

## IDE 配置

### VSCode

推荐扩展：

- Go (官方)
- Go Test Explorer
- Error Lens
- GitLens
- Docker
- Task

配置文件 `.vscode/settings.json`：

```json
{
  "go.useLanguageServer": true,
  "go.lintTool": "golangci-lint",
  "go.lintOnSave": "package",
  "go.formatTool": "goimports",
  "go.formatOnSave": true,
  "[go]": {
    "editor.codeActionsOnSave": {
      "source.organizeImports": true
    }
  }
}
```

### GoLand

1. 打开项目
2. 配置 Go Modules：`Settings > Go > Go Modules`
3. 启用格式化：`Settings > Tools > File Watchers`
4. 配置运行配置：`Run > Edit Configurations`

## 开发工作流

### 1. 启动开发服务器

```bash
# 使用 Air 热重载
air

# 或使用 Task
task dev

# 或直接运行
go run main.go api
```

### 2. 前端开发

```bash
# 启动前端开发服务器
cd web
npm run dev
```

### 3. 文档开发

```bash
# 启动文档服务器
cd docs
npm run dev
```

### 4. 运行测试

```bash
# 运行所有测试
task test

# 运行特定测试
go test ./internal/domain/user/...

# 运行集成测试
task test:integration

# 生成覆盖率报告
task test:coverage
```

### 5. 代码检查

```bash
# 运行 linter
task lint

# 修复可自动修复的问题
task lint:fix

# 格式化代码
task fmt
```

## 常用 Task 命令

```bash
# 查看所有任务
task --list

# 开发相关
task dev          # 启动开发服务器（热重载）
task build        # 构建所有组件
task clean        # 清理构建产物

# 测试相关
task test         # 运行测试
task test:unit    # 单元测试
task test:e2e     # 端到端测试

# 代码质量
task lint         # 代码检查
task fmt          # 代码格式化
task check        # 运行所有检查

# 文档
task docs:dev     # 启动文档开发服务器
task docs:build   # 构建文档

# Docker
task docker:build # 构建镜像
task docker:up    # 启动容器
task docker:down  # 停止容器
```

## 调试配置

### 使用 Delve

```bash
# 安装 delve
go install github.com/go-delve/delve/cmd/dlv@latest

# 调试运行
dlv debug main.go -- api

# 附加到运行中的进程
dlv attach <pid>
```

### VSCode 调试配置

`.vscode/launch.json`：

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Debug API",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}/main.go",
      "args": ["api"],
      "env": {
        "APP_SERVER_ADDR": ":8080",
        "APP_LOG_LEVEL": "debug"
      }
    }
  ]
}
```

## 故障排查

### 端口占用

```bash
# 查找占用端口的进程
lsof -i :8080

# 更改端口
APP_SERVER_ADDR=":9000" task dev
```

### 依赖问题

```bash
# 清理模块缓存
go clean -modcache

# 重新下载依赖
go mod download
```

### Docker 问题

```bash
# 重启 Docker 服务
docker-compose restart

# 清理容器和卷
docker-compose down -v
docker-compose up -d
```
