---
paths:
  - "internal/bootstrap/**/*.go"
---

# Bootstrap 依赖注入规范

## Container 字段顺序

```go
type Container struct {
    // 1. 基础设施（DB, Redis, EventBus）
    // 2. Domain Services
    // 3. Repositories（按模块分组）
    // 4. Use Case Handlers（按模块分组）
    // 5. HTTP Handlers
}
```

## 初始化顺序

```go
func NewContainer(cfg *Config) (*Container, error) {
    // 1️⃣ 基础设施
    // 2️⃣ Domain Services
    // 3️⃣ Repositories
    // 4️⃣ Use Case Handlers
    // 5️⃣ HTTP Handlers
    // 6️⃣ Router
}
```

## 新增模块检查清单

1. Container 结构体添加字段
2. 创建 Command/Query Repository
3. 创建 Use Case Handlers
4. 创建 HTTP Handler
5. 注册路由
