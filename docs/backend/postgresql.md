# PostgreSQL 集成

本项目使用 GORM 实现了完整的 PostgreSQL 数据库集成，提供连接管理、自动迁移、用户 CRUD 等功能。

## 功能特性

- ✅ 数据库连接管理和连接池配置
- ✅ 自动迁移支持
- ✅ 用户领域模型 (示例)
- ✅ 用户仓储接口和 GORM 实现
- ✅ 完整的用户 CRUD API
- ✅ 软删除支持
- ✅ 健康检查 (包含连接池统计)
- ✅ bcrypt 密码加密

## 快速开始

### 1. 启动 PostgreSQL

使用 Docker Compose (推荐) ：

```bash
docker-compose up -d postgres
```

或手动启动：

```bash
docker run -d --name postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=app \
  -p 5432:5432 \
  postgres:16-alpine
```

### 2. 配置数据库连接

**配置文件方式** (`config.yaml`):

```yaml
data:
  pgsql:
    url: "postgresql://postgres:postgres@localhost:5432/app?sslmode=disable"
```

**环境变量方式**:

```bash
export APP_DATA_PGSQL_URL="postgresql://postgres:postgres@localhost:5432/app?sslmode=disable"
```

### 3. 运行应用

```bash
# 使用 Task
task go:run -- api

# 或直接运行
.local/bin/251117-go-ddd-template api
```

应用启动时会**自动执行数据库迁移**，创建必要的表结构。

## 架构设计

### 代码结构

```
internal/
├── domain/
│   └── user/
│       ├── entity_user.go          # 用户领域模型
│       ├── command_repository.go   # 写接口（Create/Update/Delete...）
│       └── query_repository.go     # 读接口（Get/List/Search...）
├── infrastructure/
│   ├── database/
│   │   ├── connection.go           # 数据库连接管理
│   │   └── migrator.go             # 数据库迁移
│   └── persistence/
│       ├── user_command_repository.go # GORM 写实现
│       └── user_query_repository.go   # GORM 读实现
└── adapters/
    └── http/
        └── handler/
            ├── user.go             # 用户 HTTP 处理器
            └── health.go           # 健康检查
```

### 用户领域模型

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
    Status   string `gorm:"size:20;default:'active'" json:"status"`
}
```

**模型特点：**

- `ID` - 主键，自增
- `CreatedAt/UpdatedAt` - GORM 自动维护
- `DeletedAt` - 软删除标记
- `Username/Email` - 唯一索引
- `Password` - 响应时自动隐藏

### 仓储接口

```go
type Repository interface {
    // 创建用户
    Create(ctx context.Context, user *User) error

    // 查询用户
    FindByID(ctx context.Context, id uint) (*User, error)
    FindByUsername(ctx context.Context, username string) (*User, error)
    FindByEmail(ctx context.Context, email string) (*User, error)

    // 列表查询 (分页)
    List(ctx context.Context, page, pageSize int) (*PaginatedUsers, error)

    // 更新用户
    Update(ctx context.Context, user *User) error

    // 删除用户 (软删除)
    Delete(ctx context.Context, id uint) error

    // 统计数量
    Count(ctx context.Context) (int64, error)
}
```

## 连接池配置

在 `internal/infrastructure/database/connection.go` 中配置：

```go
// 连接池参数
sqlDB.SetMaxOpenConns(25)           // 最大打开连接数
sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数
sqlDB.SetConnMaxLifetime(5 * time.Minute)  // 连接最大生命周期
```

**推荐配置：**

- **开发环境**: MaxOpenConns=10, MaxIdleConns=5
- **生产环境**: 根据负载调整，建议 MaxOpenConns=25-100

## API 端点

### 健康检查

```bash
curl http://localhost:8080/health
```

**响应示例：**

```json
{
  "status": "ok",
  "database": "connected",
  "redis": "connected",
  "db_stats": {
    "max_open_connections": 25,
    "open_connections": 1,
    "in_use": 1,
    "idle": 0
  }
}
```

### 用户 CRUD (需要认证)

详细的 API 文档请参考 [用户接口文档](/api/users)。

**创建用户：**

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "zhangsan",
    "email": "zhangsan@example.com",
    "password": "password123",
    "full_name": "张三"
  }'
```

**获取用户列表 (分页) ：**

```bash
curl "http://localhost:8080/api/users?page=1&page_size=10" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**获取用户详情：**

```bash
curl http://localhost:8080/api/users/1 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**更新用户：**

```bash
curl -X PUT http://localhost:8080/api/users/1 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "张三丰",
    "status": "active"
  }'
```

**删除用户 (软删除) ：**

```bash
curl -X DELETE http://localhost:8080/api/users/1 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## 使用示例

### 在代码中使用仓储

```go
import (
    "context"
    "github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
)

// 从容器获取仓储
userRepo := container.UserRepository
ctx := context.Background()

// 创建用户
newUser := &user.User{
    Username: "lisi",
    Email:    "lisi@example.com",
    Password: hashedPassword,
    FullName: "李四",
    Status:   "active",
}
err := userRepo.Create(ctx, newUser)

// 查询用户
u, err := userRepo.FindByID(ctx, 1)
u, err = userRepo.FindByUsername(ctx, "lisi")
u, err = userRepo.FindByEmail(ctx, "lisi@example.com")

// 分页查询
result, err := userRepo.List(ctx, 1, 10)
fmt.Printf("Total: %d, Page: %d/%d\n",
    result.Total, result.Page, result.TotalPages)

// 更新用户
u.FullName = "李四丰"
err = userRepo.Update(ctx, u)

// 删除用户 (软删除)
err = userRepo.Delete(ctx, 1)

// 统计数量
count, err := userRepo.Count(ctx)
```

## 特性说明

### 1. 自动迁移

应用启动时自动执行迁移，创建或更新表结构：

```go
// 在 bootstrap.NewContainer() 中
migrator := database.NewMigrator(conn)
if err := migrator.AutoMigrate(); err != nil {
    log.Fatalf("Failed to migrate database: %v", err)
}
```

**迁移规则：**

- 自动创建表 (如果不存在)
- 自动添加缺失的列
- 自动添加索引和约束
- 不会删除已存在的列

### 2. 软删除

使用 GORM 的 `DeletedAt` 字段实现软删除：

```go
// 删除操作只设置 DeletedAt 时间戳
db.Delete(&user, 1)

// 查询时自动排除已删除记录
db.Find(&users)

// 包含已删除记录的查询
db.Unscoped().Find(&users)

// 永久删除
db.Unscoped().Delete(&user, 1)
```

### 3. 密码加密

用户密码自动使用 bcrypt 加密：

```go
// 在 User 模型中
func (u *User) HashPassword() error {
    hashedBytes, err := bcrypt.GenerateFromPassword(
        []byte(u.Password),
        bcrypt.DefaultCost,
    )
    if err != nil {
        return err
    }
    u.Password = string(hashedBytes)
    return nil
}

// 验证密码
func (u *User) CheckPassword(password string) bool {
    err := bcrypt.CompareHashAndPassword(
        []byte(u.Password),
        []byte(password),
    )
    return err == nil
}
```

### 4. 数据验证

使用 Gin 的 binding 标签验证请求数据：

```go
type CreateUserRequest struct {
    Username string `json:"username" binding:"required,min=3,max=50"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
    FullName string `json:"full_name" binding:"max=100"`
}
```

## 连接字符串格式

```
postgresql://[用户名]:[密码]@[主机]:[端口]/[数据库]?[参数]
```

**示例：**

```bash
# 本地开发 (禁用 SSL)
postgresql://postgres:postgres@localhost:5432/app?sslmode=disable

# 生产环境 (启用 SSL)
postgresql://user:pass@db.example.com:5432/prod_db?sslmode=require

# 使用 Unix Socket
postgresql:///mydb?host=/var/run/postgresql

# 自定义连接池参数
postgresql://user:pass@localhost:5432/app?pool_max_conns=25&pool_min_conns=5
```

## 性能优化

### 1. 索引优化

```go
// 在模型中定义索引
type User struct {
    Username string `gorm:"uniqueIndex"`           // 唯一索引
    Email    string `gorm:"uniqueIndex"`           // 唯一索引
    Status   string `gorm:"index"`                 // 普通索引
    DeletedAt gorm.DeletedAt `gorm:"index"`       // 软删除索引
}
```

### 2. 批量操作

```go
// 批量创建
users := []*user.User{{...}, {...}}
db.CreateInBatches(users, 100) // 每批 100 条

// 批量更新
db.Model(&user.User{}).
    Where("status = ?", "inactive").
    Update("status", "active")
```

### 3. 预加载关联

```go
// 预加载关联数据
db.Preload("Orders").Find(&users)
db.Preload("Orders.Items").Find(&users)
```

### 4. 选择字段

```go
// 只查询需要的字段
db.Select("id", "username", "email").Find(&users)

// 排除某些字段
db.Omit("password").Find(&users)
```

## 事务支持

```go
// 手动事务
err := db.Transaction(func(tx *gorm.DB) error {
    // 在事务中执行操作
    if err := tx.Create(&user1).Error; err != nil {
        return err
    }
    if err := tx.Create(&user2).Error; err != nil {
        return err
    }
    return nil
})

// 嵌套事务
db.Transaction(func(tx *gorm.DB) error {
    tx.Create(&user1)

    tx.Transaction(func(tx2 *gorm.DB) error {
        tx2.Create(&user2)
        return nil
    })

    return nil
})
```

## 故障排查

### 连接失败

```bash
# 检查 PostgreSQL 是否运行
docker ps | grep postgres

# 查看 PostgreSQL 日志
docker logs go-ddd-postgres

# 测试连接
docker exec -it go-ddd-postgres psql -U postgres -d app

# 使用 psql 客户端测试
psql "postgresql://postgres:postgres@localhost:5432/app"
```

### 迁移失败

**常见原因：**

1. 模型定义有误
2. 数据库权限不足
3. 表已存在且结构冲突

**解决方法：**

```bash
# 查看应用日志
tail -f logs/app.log

# 手动连接数据库检查
docker exec -it go-ddd-postgres psql -U postgres -d app
\dt          # 列出所有表
\d users     # 查看 users 表结构
```

### 性能问题

**检查连接池状态：**

```bash
curl http://localhost:8080/health | jq '.db_stats'
```

**优化建议：**

1. 添加合适的索引
2. 调整连接池参数
3. 使用 GORM 日志查看慢查询
4. 分页查询大数据集
5. 使用 `Select` 只查询需要的字段

**启用 GORM 日志：**

```go
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info),
})
```

## 最佳实践

1. **使用事务**：对于涉及多个操作的业务逻辑
2. **避免 N+1 查询**：使用 Preload 预加载关联
3. **合理使用索引**：为常用查询字段添加索引
4. **定期维护**：VACUUM、ANALYZE
5. **监控性能**：记录慢查询，优化 SQL
6. **备份策略**：定期备份数据库
7. **连接池管理**：根据负载调整连接池大小

## 扩展功能

### 1. 添加新的领域模型

```go
// 1. 定义模型
type Product struct {
    gorm.Model
    Name  string
    Price float64
}

// 2. 定义仓储接口
type ProductRepository interface {
    Create(ctx context.Context, p *Product) error
    FindByID(ctx context.Context, id uint) (*Product, error)
}

// 3. 实现仓储
type productRepository struct {
    db *gorm.DB
}

// 4. 在迁移中添加
migrator.AutoMigrate(&Product{})

// 5. 注入到容器
container.ProductRepository = NewProductRepository(conn.DB)
```

### 2. 全文搜索

```go
// 使用 PostgreSQL 的全文搜索
db.Where("to_tsvector('english', full_name) @@ to_tsquery('english', ?)",
    query).Find(&users)
```

### 3. 读写分离

```go
// 配置主从数据库
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
db.Use(dbresolver.Register(dbresolver.Config{
    Sources:  []gorm.Dialector{postgres.Open(sourceDSN)},
    Replicas: []gorm.Dialector{postgres.Open(replicaDSN)},
    Policy:   dbresolver.RandomPolicy{},
}))
```

## 下一步

- 了解 [认证授权](/backend/authentication)
- 学习 [Redis 缓存](/backend/redis)
- 查看 [用户 API 文档](/api/users)
- 探索 [项目架构](/backend/)
