# 开发环境设置

本指南介绍如何设置 Go DDD Template 的开发环境。

<!--TOC-->

## Table of Contents

- [系统要求](#系统要求) `:25+8`
- [环境准备](#环境准备) `:33+16`
- [项目设置](#项目设置) `:49+50`
  - [1. 克隆代码](#1-克隆代码) `:51+7`
  - [2. 安装依赖](#2-安装依赖) `:58+13`
  - [3. 环境配置](#3-环境配置) `:71+8`
  - [4. 启动基础服务](#4-启动基础服务) `:79+10`
  - [5. 初始化数据库](#5-初始化数据库) `:89+10`
- [开发工作流](#开发工作流) `:99+19`
- [故障排查](#故障排查) `:118+22`
  - [端口占用](#端口占用) `:120+7`
  - [依赖问题](#依赖问题) `:127+7`
  - [Docker 问题](#docker-问题) `:134+6`

<!--TOC-->

## 系统要求

- **Go** 1.25.4+
- **Docker** 20.10+ & Docker Compose 2.0+
- **Git** 2.0+
- **Task** - 任务自动化工具
- **Air** - Go 热重载工具（可选）

## 环境准备

```bash
# 安装 Task (任务运行器)
go install github.com/go-task/task/v3/cmd/task@latest

# 安装 Air (热重载)
go install github.com/air-verse/air@latest

# 安装 golangci-lint (代码检查)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 安装 swag (Swagger 文档生成)
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

```bash
cp configs/config.example.yaml configs/config.yaml
```

编辑 `configs/config.yaml` 配置数据库连接等参数。

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

## 开发工作流

```bash
# 启动后端（热重载）
air

# 启动前端
cd web && npm run dev

# 启动文档服务
cd docs && npm run dev

# 运行测试
go test ./...

# 代码检查
golangci-lint run
```

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
