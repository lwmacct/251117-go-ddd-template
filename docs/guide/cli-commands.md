# CLI 命令

本项目提供多个 CLI 命令来管理应用的不同方面。所有命令都通过统一的入口运行。

## 命令概览

```bash
go run main.go [command] [options]
```

可用命令：

- `api` - 启动 REST API 服务器
- `migrate` - 数据库迁移管理
- `seed` - 填充数据库种子数据
- `worker` - 启动后台任务处理器

## API 服务器

启动 REST API 服务器。

### 基本用法

```bash
go run main.go api
```

或使用别名：

```bash
go run main.go serve
go run main.go server
```

### 选项

| 选项       | 说明         | 默认值         |
| ---------- | ------------ | -------------- |
| `--addr`   | 监听地址     | `0.0.0.0:8080` |
| `--static` | 静态资源目录 | `web/dist`     |

### 示例

```bash
# 在端口 9000 启动
go run main.go api --addr :9000

# 指定静态资源目录
go run main.go api --static ./dist

# 使用环境变量和命令行组合
APP_JWT_SECRET="secret" go run main.go api --addr :9000
```

### 优雅关闭

服务器支持优雅关闭。接收到 `SIGINT` 或 `SIGTERM` 信号时：

1. 停止接受新请求
2. 等待现有请求处理完成（最多 30 秒）
3. 关闭数据库和 Redis 连接
4. 退出进程

---

## 数据库迁移

管理数据库表结构的迁移。

### 子命令

#### migrate up

执行数据库迁移，创建或更新所有表结构。

```bash
go run main.go migrate up
```

**功能：**

- 自动创建 `migrations` 表记录迁移历史
- 执行所有已注册模型的迁移
- 记录迁移版本和时间

**输出：**

```
INFO Running database migration...
INFO Database migration completed
```

#### migrate status

查看已执行的迁移记录。

```bash
go run main.go migrate status
```

**输出示例：**

```
INFO Migration history:

  ID | Version        | Name          | Applied At
  ---|----------------|---------------|----------------------------
  1  | 20251118142000 | auto_migrate  | 2025-11-18 14:20:00
  2  | 20251118143000 | auto_migrate  | 2025-11-18 14:30:00
```

#### migrate fresh

删除所有表并重新执行迁移。

::: danger 危险操作
此命令会删除所有数据！仅适用于开发环境。
:::

```bash
go run main.go migrate fresh
```

**交互式确认：**

```
⚠️  WARNING: This will delete ALL data in the database!
Are you sure you want to continue? (yes/no):
```

输入 `yes` 确认执行。

**跳过确认（生产环境需要 --force）：**

```bash
# 开发环境可以直接强制执行
go run main.go migrate fresh --force

# 生产环境即使使用 --force 也会被拒绝
APP_SERVER_ENV=production go run main.go migrate fresh --force
# Error: Cannot run fresh migration in production environment
```

### 添加新模型到迁移

1. 在 `internal/domain/` 定义模型
2. 在 `internal/bootstrap/container.go` 的 `GetAllModels()` 函数中注册：

```go
func GetAllModels() []any {
    return []any{
        &user.User{},
        &product.Product{},  // 新增模型
    }
}
```

3. 执行迁移：

```bash
go run main.go migrate up
```

---

## 数据库种子

填充开发和测试环境的示例数据。

### 基本用法

```bash
go run main.go seed
```

### 功能

- 填充示例用户数据
- 自动跳过已存在的记录（幂等性）
- 使用 bcrypt 加密密码

### 默认用户

| 用户名     | 邮箱              | 密码        | 角色     |
| ---------- | ----------------- | ----------- | -------- |
| `admin`    | admin@example.com | password123 | 管理员   |
| `testuser` | test@example.com  | password123 | 普通用户 |
| `demo`     | demo@example.com  | password123 | 演示用户 |

### 添加自定义种子

1. 在 `internal/infrastructure/database/seeds/` 创建种子文件：

```go
// product_seeder.go
package seeds

import (
    "context"
    "gorm.io/gorm"
    "github.com/lwmacct/251117-go-ddd-template/internal/domain/product"
)

type ProductSeeder struct{}

func (s *ProductSeeder) Seed(ctx context.Context, db *gorm.DB) error {
    products := []product.Product{
        {Name: "Product 1", Price: 100},
        {Name: "Product 2", Price: 200},
    }

    for _, p := range products {
        var existing product.Product
        if err := db.Where("name = ?", p.Name).First(&existing).Error; err == gorm.ErrRecordNotFound {
            if err := db.Create(&p).Error; err != nil {
                return err
            }
        }
    }

    return nil
}
```

2. 在 `internal/commands/seed/seed.go` 的 `getAllSeeders()` 中注册：

```go
func getAllSeeders() []database.Seeder {
    return []database.Seeder{
        &seeds.UserSeeder{},
        &seeds.ProductSeeder{},  // 新增
    }
}
```

3. 执行种子：

```bash
go run main.go seed
```

---

## 后台任务处理器

启动 Worker 进程处理队列中的后台任务。

### 基本用法

```bash
go run main.go worker
```

### 选项

| 选项            | 短选项 | 说明       | 默认值    |
| --------------- | ------ | ---------- | --------- |
| `--queue`       | `-q`   | 队列名称   | `default` |
| `--concurrency` | `-c`   | 并发处理数 | `5`       |

### 示例

```bash
# 默认配置启动
go run main.go worker

# 指定队列和并发数
go run main.go worker --queue jobs --concurrency 10

# 使用短选项
go run main.go worker -q jobs -c 10
```

### 功能特性

**并发处理**：

- 支持多个 worker 并发处理任务
- 每个 worker 独立处理任务，互不干扰

**优雅关闭**：

- 接收到终止信号后停止接受新任务
- 等待当前处理中的任务完成
- 清理资源后退出

**错误处理**：

- 任务处理失败会记录错误日志
- 可扩展实现重试逻辑或死信队列

### 任务入队示例

在代码中添加任务到队列：

```go
import "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/queue"

// 创建队列
q := queue.NewRedisQueue(redisClient, "default")

// 添加任务
job := map[string]interface{}{
    "type": "send_email",
    "to": "user@example.com",
    "subject": "Welcome",
}

if err := q.Enqueue(ctx, job); err != nil {
    log.Fatal(err)
}
```

### 自定义任务处理器

创建自定义的任务处理器：

```go
// email_handler.go
package handlers

import (
    "context"
    "encoding/json"
    "log/slog"
)

type EmailHandler struct{}

func (h *EmailHandler) Handle(ctx context.Context, data []byte) error {
    var job struct {
        Type    string `json:"type"`
        To      string `json:"to"`
        Subject string `json:"subject"`
    }

    if err := json.Unmarshal(data, &job); err != nil {
        return err
    }

    slog.Info("Sending email", "to", job.To, "subject", job.Subject)

    // 实际发送邮件逻辑
    // ...

    return nil
}
```

在 `internal/commands/worker/worker.go` 中使用：

```go
// 替换默认 handler
handler := &handlers.EmailHandler{}
processor := queue.NewProcessor(q, handler, concurrency)
```

---

## 开发工作流

### 初始化项目

```bash
# 1. 启动依赖服务
docker-compose up -d

# 2. 执行数据库迁移
go run main.go migrate up

# 3. 填充种子数据
go run main.go seed

# 4. 启动 API 服务器
go run main.go api
```

### 开发环境（自动迁移）

如果希望在开发环境自动执行迁移，可以在配置中启用：

```yaml
# config.yaml
data:
  auto_migrate: true # 仅开发环境推荐
```

或使用环境变量：

```bash
APP_DATA_AUTO_MIGRATE=true go run main.go api
```

### 生产部署流程

```bash
# 1. 执行迁移（生产环境应该独立执行，不要自动迁移）
./app migrate up

# 2. 启动 API 服务器
./app api

# 3. 启动 Worker（如果需要）
./app worker --queue jobs --concurrency 20
```

---

## 故障排查

### 迁移失败

**问题**：数据库连接失败

**解决**：

1. 检查数据库是否运行：`docker-compose ps`
2. 验证连接字符串：检查环境变量或配置文件
3. 检查网络连接

**问题**：表已存在

**解决**：

- GORM 的 AutoMigrate 是安全的，会自动更新表结构而不是覆盖
- 如需重置，使用 `migrate fresh`（仅开发环境）

### Worker 无法启动

**问题**：Redis 连接失败

**解决**：

1. 确保 Redis 正在运行
2. 检查连接字符串格式：`redis://localhost:6379/0`
3. 验证网络连接

### 性能优化

**数据库连接池**：

- 在 `internal/infrastructure/database/connection.go` 配置连接池大小

**Worker 并发数**：

- 根据 CPU 核心数和任务类型调整 `--concurrency`
- I/O 密集型任务可以设置更高的并发数
- CPU 密集型任务建议不超过 CPU 核心数

---

## 相关链接

- [项目架构](/guide/architecture)
- [配置系统](/guide/configuration)
- [快速开始](/guide/getting-started)
