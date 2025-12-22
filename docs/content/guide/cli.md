# CLI 命令行工具

Go DDD Template 提供了丰富的命令行工具来管理应用的各个方面。

<!--TOC-->

## Table of Contents

- [命令概览](#命令概览) `:35+10`
- [核心命令](#核心命令) `:45+77`
  - [api - 启动 HTTP 服务](#api-启动-http-服务) `:47+20`
  - [migrate - 数据库迁移](#migrate-数据库迁移) `:67+21`
  - [seed - 数据填充](#seed-数据填充) `:88+19`
  - [validate - 配置验证](#validate-配置验证) `:107+15`
- [全局选项](#全局选项) `:122+11`
- [使用 Task 运行](#使用-task-运行) `:133+15`
- [开发模式](#开发模式) `:148+27`
  - [热重载运行](#热重载运行) `:150+15`
  - [构建命令](#构建命令) `:165+10`
- [生产部署](#生产部署) `:175+29`
  - [Docker 运行](#docker-运行) `:177+13`
  - [Systemd 服务](#systemd-服务) `:190+14`
- [命令示例](#命令示例) `:204+41`
  - [完整初始化流程](#完整初始化流程) `:206+16`
  - [数据库重置](#数据库重置) `:222+10`
  - [生产环境检查](#生产环境检查) `:232+13`
- [故障排查](#故障排查) `:245+34`
  - [命令不存在](#命令不存在) `:247+8`
  - [权限错误](#权限错误) `:255+8`
  - [配置错误](#配置错误) `:263+8`
  - [数据库连接失败](#数据库连接失败) `:271+8`

<!--TOC-->

## 命令概览

```bash
# 查看所有可用命令
./251117-go-ddd-template --help

# 命令格式
./251117-go-ddd-template [全局选项] 命令 [命令选项]
```

## 核心命令

### api - 启动 HTTP 服务

启动 Web API 服务器：

```bash
# 使用默认配置
./251117-go-ddd-template api

# 指定端口
./251117-go-ddd-template api --addr :9000

# 指定配置文件
./251117-go-ddd-template --config configs/prod.yaml api
```

选项：

- `--addr` - 服务监听地址（默认：`:8080`）
- `--env` - 运行环境（development/production）

### migrate - 数据库迁移

管理数据库结构变更：

```bash
# 执行所有未运行的迁移
./251117-go-ddd-template migrate up

# 回滚最后一次迁移
./251117-go-ddd-template migrate down

# 查看迁移状态
./251117-go-ddd-template migrate status

# 回滚到指定版本
./251117-go-ddd-template migrate down --to 20240101120000

# 强制设置版本（危险操作）
./251117-go-ddd-template migrate force 20240101120000
```

### seed - 数据填充

填充开发测试数据：

```bash
# 填充所有种子数据
./251117-go-ddd-template seed

# 填充特定类型数据
./251117-go-ddd-template seed --type users
./251117-go-ddd-template seed --type roles
```

预置数据：

- 管理员账号：`admin / password123`
- 测试账号：`testuser / password123`
- 示例角色和权限

### validate - 配置验证

验证配置文件和环境：

```bash
# 验证当前配置
./251117-go-ddd-template validate

# 验证指定配置文件
./251117-go-ddd-template --config configs/test.yaml validate

# 详细输出
./251117-go-ddd-template validate --verbose
```

## 全局选项

所有命令都支持以下全局选项：

| 选项         | 说明         | 示例                        |
| ------------ | ------------ | --------------------------- |
| --config     | 配置文件路径 | `--config configs/dev.yaml` |
| --log-level  | 日志级别     | `--log-level debug`         |
| --log-format | 日志格式     | `--log-format json`         |
| --help       | 显示帮助     | `--help`                    |

## 使用 Task 运行

推荐使用 Task 来运行命令：

```bash
# 运行 API 服务
task go:run -- api

# 执行迁移
task go:run -- migrate up

# 填充数据
task go:run -- seed
```

## 开发模式

### 热重载运行

使用 Air 实现代码变更自动重启：

```bash
# 安装 air
go install github.com/air-verse/air@latest

# 启动热重载
air

# 或使用 Task
task dev
```

### 构建命令

```bash
# 构建二进制文件
task go:build

# 输出位置
.local/bin/251117-go-ddd-template
```

## 生产部署

### Docker 运行

```bash
# 构建镜像
docker build -t myapp .

# 运行容器
docker run -p 8080:8080 \
  -e APP_JWT_SECRET=secret \
  -e APP_DATA_PGSQL_URL=postgresql://... \
  myapp api
```

### Systemd 服务

```bash
# 复制二进制文件
sudo cp .local/bin/251117-go-ddd-template /usr/local/bin/

# 创建服务文件
sudo vi /etc/systemd/system/myapp.service

# 启动服务
sudo systemctl start myapp
sudo systemctl enable myapp
```

## 命令示例

### 完整初始化流程

```bash
# 1. 启动依赖服务
docker-compose up -d

# 2. 执行数据库迁移
task go:run -- migrate up

# 3. 填充种子数据
task go:run -- seed

# 4. 启动应用
task go:run -- api
```

### 数据库重置

```bash
# 停止应用
# 重置数据库
task go:run -- migrate down --all
task go:run -- migrate up
task go:run -- seed
```

### 生产环境检查

```bash
# 验证配置
APP_SERVER_ENV=production ./app validate

# 检查迁移状态
./app migrate status

# 健康检查
curl http://localhost:8080/health
```

## 故障排查

### 命令不存在

确保已构建二进制文件：

```bash
task go:build
```

### 权限错误

添加执行权限：

```bash
chmod +x .local/bin/251117-go-ddd-template
```

### 配置错误

验证配置文件：

```bash
./app validate --verbose
```

### 数据库连接失败

检查连接字符串和数据库状态：

```bash
docker-compose ps
psql $APP_DATA_PGSQL_URL -c "SELECT 1"
```
