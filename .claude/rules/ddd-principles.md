---
paths:
  - "internal/**/*.go"
---

# DDD 核心原则与禁止操作

<!--TOC-->

## Table of Contents

- [依赖方向（严格单向）](#依赖方向严格单向) `:18+11`
- [10 条核心原则](#10-条核心原则) `:29+22`
- [禁止操作](#禁止操作) `:51+9`

<!--TOC-->

## 依赖方向（严格单向）

```
Adapters → Application → Domain ← Infrastructure
```

- **Adapters** 依赖 Application
- **Application** 依赖 Domain
- **Infrastructure** 实现 Domain 接口（依赖倒置）
- **Domain** 不依赖任何外层

## 10 条核心原则

1. **依赖倒置** - Domain 层定义接口，Infrastructure 层实现，Application 层依赖接口

2. **领域纯度** - Domain 模型仅承载业务语义，不得引用 GORM 或其它 ORM Tag；Infra 通过 `*_model.go` 负责映射

3. **CQRS 分离** - 写操作用 CommandRepository，读操作用 QueryRepository

4. **Use Case 模式** - 业务逻辑在 Application 层的 Handler 中处理，HTTP Handler 只做入参/出参绑定

5. **富领域模型** - 业务行为通过方法体现（如 `entity.Activate()`），禁止直接修改结构体字段

6. **单一职责** - Handler 仅做 HTTP 转换，Use Case Handler 编排业务，Repository 访问数据

7. **依赖注入** - 所有依赖在 `container.go` 中注册

8. **统一响应** - HTTP 响应使用 `adapters/http/response` 包

9. **接口优先** - 先定义 Domain 接口，再实现 Infrastructure

10. **统一架构** - 所有模块必须遵循最新 DDD+CQRS 约定，发现旧式实现立即拆分重构，禁止新增兼容层

## 禁止操作

- ❌ 在 HTTP Handler 中编排业务逻辑或直接调用 Repository
- ❌ 在 Application 层直接依赖 Infrastructure 实现（只能依赖 Domain 接口）
- ❌ Domain 层 import 外层代码（禁止 `gorm`/Infra 依赖）
- ❌ Command 和 Query Repository 混用，或复用旧的 `repository.go`
- ❌ 跳过 Use Case，直接从 Handler 或 Infra 操作数据库
- ❌ 在 Domain 实体中使用 GORM Tag
- ❌ 新增兼容层或过渡代码
