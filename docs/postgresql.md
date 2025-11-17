# PostgreSQL 实现说明

## 概述

本项目已实现完整的 PostgreSQL 集成（使用 GORM），包括：
- 数据库连接管理和连接池配置
- 自动迁移支持
- 用户领域模型（示例）
- 用户仓储接口和 GORM 实现
- 完整的用户 CRUD API
- 健康检查（包含数据库状态和连接池统计）

## 快速开始

### 1. 启动 PostgreSQL 服务

使用 Docker Compose：

```bash
docker-compose up -d postgres
```

或手动启动 PostgreSQL：

```bash
# 使用 Docker
docker run -d --name postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=myapp \
  -p 5432:5432 \
  postgres:16-alpine
```

### 2. 配置数据库连接

方式一：配置文件（`config.yaml`）
```yaml
data:
  pgsql_url: "postgresql://postgres:postgres@localhost:5432/myapp?sslmode=disable"
```

方式二：环境变量
```bash
export APP_DATA_PGSQL_URL="postgresql://postgres:postgres@localhost:5432/myapp?sslmode=disable"
```

### 3. 运行应用

```bash
# 启动服务（会自动迁移数据库表）
task go:run -- api

# 或直接运行
.local/bin/bd-vmalert api
```

应用启动时会自动执行数据库迁移，创建 `users` 表。

## API 端点

### 健康检查（包含数据库状态）

```bash
curl http://localhost:8080/health
```

响应示例：
```json
{
  "status": "ok",
  "checks": {
    "database": {
      "status": "healthy",
      "stats": {
        "idle": 0,
        "in_use": 1,
        "max_idle_closed": 0,
        "max_lifetime_closed": 0,
        "max_open_connections": 25,
        "open_connections": 1,
        "wait_count": 0,
        "wait_duration": "0s"
      }
    },
    "redis": {
      "status": "healthy"
    }
  }
}
```

### 创建用户

```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "zhangsan",
    "email": "zhangsan@example.com",
    "password": "password123",
    "full_name": "张三"
  }'
```

### 获取用户列表（分页）

```bash
curl "http://localhost:8080/api/users?page=1&limit=10"
```

响应示例：
```json
{
  "data": [
    {
      "id": 1,
      "username": "zhangsan",
      "email": "zhangsan@example.com",
      "full_name": "张三",
      "avatar": "",
      "bio": "",
      "status": "active",
      "created_at": "2025-01-18T00:00:00Z",
      "updated_at": "2025-01-18T00:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 1
  }
}
```

### 获取用户详情

```bash
curl http://localhost:8080/api/users/1
```

### 更新用户

```bash
curl -X PUT http://localhost:8080/api/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "张三丰",
    "bio": "武当派创始人",
    "status": "active"
  }'
```

### 删除用户（软删除）

```bash
curl -X DELETE http://localhost:8080/api/users/1
```

## 代码结构

```
internal/
├── domain/
│   └── user/
│       ├── model.go       # 用户领域模型和 DTO
│       └── repository.go  # 用户仓储接口
├── infrastructure/
│   ├── database/
│   │   ├── connection.go  # 数据库连接管理
│   │   └── migrator.go    # 数据库迁移
│   └── persistence/
│       └── user_repository.go  # 用户仓储 GORM 实现
└── adapters/
    └── http/
        ├── handler_user.go     # 用户 HTTP 处理器
        └── handler_health.go   # 健康检查（包含数据库状态）
```

## 领域模型

### User 模型

```go
type User struct {
    ID        uint           `gorm:"primarykey" json:"id"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 软删除

    Username string `gorm:"uniqueIndex;size:50;not null" json:"username"`
    Email    string `gorm:"uniqueIndex;size:100;not null" json:"email"`
    Password string `gorm:"size:255;not null" json:"-"`
    FullName string `gorm:"size:100" json:"full_name"`
    Avatar   string `gorm:"size:255" json:"avatar"`
    Bio      string `gorm:"type:text" json:"bio"`
    Status   string `gorm:"size:20;default:'active'" json:"status"`
}
```

## 仓储接口

```go
type Repository interface {
    Create(ctx context.Context, user *User) error
    GetByID(ctx context.Context, id uint) (*User, error)
    GetByUsername(ctx context.Context, username string) (*User, error)
    GetByEmail(ctx context.Context, email string) (*User, error)
    List(ctx context.Context, offset, limit int) ([]*User, error)
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id uint) error
    Count(ctx context.Context) (int64, error)
}
```

## 特性

### 1. 数据库连接管理

- **连接池配置**：
  - 最大打开连接数：25
  - 最大空闲连接数：10
  - 连接最大生命周期：5 分钟
- **自动健康检查**：启动时和运行时
- **优雅关闭**：确保所有连接正确关闭

### 2. 自动迁移

- 应用启动时自动创建/更新表结构
- 支持 GORM 标签定义索引、约束等
- 禁用外键约束（在应用层处理）

### 3. 软删除

- 使用 GORM 的 `DeletedAt` 字段
- 删除操作不会真正删除记录
- 查询时自动过滤已删除记录

### 4. 密码加密

- 使用 bcrypt 加密存储密码
- 创建用户时自动加密
- 响应中自动隐藏密码字段

### 5. 数据验证

- 使用 Gin 的 binding 验证
- 支持必填、邮箱格式、长度限制等
- 自定义错误消息

## 使用示例

### 在代码中使用仓储

```go
// 从容器获取仓储
userRepo := container.UserRepository

ctx := context.Background()

// 创建用户
newUser := &user.User{
    Username: "lisi",
    Email:    "lisi@example.com",
    Password: hashedPassword,
    FullName: "李四",
}
err := userRepo.Create(ctx, newUser)

// 查询用户
u, err := userRepo.GetByID(ctx, 1)
u, err = userRepo.GetByUsername(ctx, "lisi")
u, err = userRepo.GetByEmail(ctx, "lisi@example.com")

// 更新用户
u.FullName = "李四丰"
err = userRepo.Update(ctx, u)

// 删除用户（软删除）
err = userRepo.Delete(ctx, 1)

// 分页查询
users, err := userRepo.List(ctx, 0, 10)

// 统计数量
count, err := userRepo.Count(ctx)
```

## 数据库配置

### 连接字符串格式

```
postgresql://[用户名]:[密码]@[主机]:[端口]/[数据库]?[参数]
```

示例：
```bash
# 本地开发
postgresql://postgres:postgres@localhost:5432/myapp?sslmode=disable

# 生产环境（启用 SSL）
postgresql://user:pass@db.example.com:5432/prod_db?sslmode=require

# 使用连接池参数
postgresql://user:pass@localhost:5432/myapp?sslmode=disable&pool_max_conns=25
```

## 注意事项

1. **密码安全**：密码使用 bcrypt 加密，成本因子为默认值（10）
2. **唯一索引**：username 和 email 字段有唯一索引
3. **软删除**：删除操作是软删除，数据仍在数据库中
4. **时区**：所有时间使用 UTC
5. **连接池**：生产环境需要根据实际负载调整连接池参数
6. **事务**：当前未使用事务，可根据需要添加

## 扩展建议

1. **添加事务支持**
   ```go
   func (r *userRepository) CreateWithTransaction(ctx context.Context, u *user.User) error {
       return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
           // 事务逻辑
           return tx.Create(u).Error
       })
   }
   ```

2. **添加更多查询方法**
   ```go
   func (r *userRepository) FindByStatus(ctx context.Context, status string) ([]*user.User, error)
   func (r *userRepository) Search(ctx context.Context, keyword string) ([]*user.User, error)
   ```

3. **实现缓存层**
   - 结合 Redis 实现查询缓存
   - 更新时清除缓存

4. **添加审计字段**
   - CreatedBy、UpdatedBy
   - IP 地址、用户代理

5. **实现完整的认证系统**
   - JWT token 生成
   - 登录、注册、密码重置
   - 权限管理

## 故障排查

### 连接失败

```bash
# 检查 PostgreSQL 是否运行
docker ps | grep postgres

# 查看 PostgreSQL 日志
docker logs go-ddd-postgres

# 测试连接
docker exec -it go-ddd-postgres psql -U postgres -d myapp
```

### 迁移失败

- 检查模型定义是否正确
- 查看应用日志中的错误信息
- 手动连接数据库检查表结构

### 性能问题

- 检查连接池统计（通过健康检查端点）
- 调整连接池参数
- 添加数据库索引
- 使用 GORM 的日志查看慢查询
