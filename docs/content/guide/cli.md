# CLI 命令行工具

Go DDD Template 提供命令行工具来管理应用。

<!--TOC-->

## Table of Contents

- [命令概览](#命令概览) `:22+10`
- [核心命令](#核心命令) `:32+33`
  - [api - 启动 HTTP 服务](#api-启动-http-服务) `:34+8`
  - [migrate - 数据库迁移](#migrate-数据库迁移) `:42+8`
  - [seed - 数据填充](#seed-数据填充) `:50+9`
  - [worker - 后台任务处理](#worker-后台任务处理) `:59+6`
- [全局选项](#全局选项) `:65+8`
- [使用 Task 运行](#使用-task-运行) `:73+9`
- [开发模式](#开发模式) `:82+11`
- [故障排查](#故障排查) `:93+13`

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

```bash
./251117-go-ddd-template api
./251117-go-ddd-template api --addr :9000
./251117-go-ddd-template --config configs/prod.yaml api
```

### migrate - 数据库迁移

```bash
./251117-go-ddd-template migrate up       # 执行迁移
./251117-go-ddd-template migrate down     # 回滚最后一次
./251117-go-ddd-template migrate status   # 查看状态
```

### seed - 数据填充

```bash
./251117-go-ddd-template seed             # 填充种子数据
./251117-go-ddd-template seed --type users
```

预置账号：`admin / password123`

### worker - 后台任务处理

```bash
./251117-go-ddd-template worker           # 启动后台任务处理器
```

## 全局选项

| 选项        | 说明         | 示例                        |
| ----------- | ------------ | --------------------------- |
| --config    | 配置文件路径 | `--config configs/dev.yaml` |
| --log-level | 日志级别     | `--log-level debug`         |
| --help      | 显示帮助     | `--help`                    |

## 使用 Task 运行

```bash
task go:run -- api           # 运行 API 服务
task go:run -- migrate up    # 执行迁移
task go:run -- seed          # 填充数据
task go:run -- worker        # 启动 worker
```

## 开发模式

```bash
# 热重载运行
air

# 构建二进制
task go:build
# 输出: .local/bin/251117-go-ddd-template
```

## 故障排查

```bash
# 命令不存在
task go:build

# 权限错误
chmod +x .local/bin/251117-go-ddd-template

# 数据库连接失败
docker-compose ps
psql $APP_DATA_PGSQL_URL -c "SELECT 1"
```
