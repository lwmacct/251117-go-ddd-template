---
paths:
  - "internal/**/*.go"
---

# DDD 核心原则

## 依赖方向

```
Adapters → Application → Domain ← Infrastructure
```

## 核心原则

1. **依赖倒置** - Domain 定义接口，Infrastructure 实现
2. **领域纯度** - Domain 无 ORM 依赖
3. **CQRS 分离** - Command/Query Repository 分离
4. **Use Case** - 业务逻辑在 Application Handler
5. **富领域模型** - 行为通过方法体现
6. **单一职责** - 各层职责明确
7. **依赖注入** - Container 统一注册
8. **统一响应** - 使用 response 包
9. **接口优先** - 先 Domain 接口后实现
10. **统一架构** - 禁止兼容层

## 禁止操作

- ❌ Handler 中编排业务或调用 Repository
- ❌ Application 依赖 Infrastructure 实现
- ❌ Domain 层 import 外层代码
- ❌ Command/Query Repository 混用
- ❌ Domain 实体使用 GORM Tag
