---
paths:
  - ".claude/rules/design/module-creation.md"
---

# DDD 模块创建流程

按 DDD 四层架构 + CQRS 模式创建新模块的规划步骤。

> 各层实现细节见 `.claude/rules/backend/*.md`

---

## 规划步骤总览

```
Phase 1: Domain      → 定义领域边界和接口契约
Phase 2: Persistence → 实现数据持久化
Phase 3: Application → 实现业务用例
Phase 4: Adapters    → 暴露 HTTP 接口
Phase 5: Container   → 依赖注入装配
Phase 6: 收尾       → 种子数据和验证
```

---

## Phase 1: Domain 层

**目标**：定义领域边界，确定实体和仓储契约。

**产出物**：

- `doc.go` - 领域边界描述
- `entity.go` - 核心实体
- `repository.go` - Command/Query 接口
- `errors.go` - 领域错误

**规划要点**：

- 识别核心实体和关联实体
- 确定读写操作边界（CQRS 分离）
- 多实体时按实体拆分接口文件

---

## Phase 2: Persistence 层

**目标**：实现 Domain 层定义的仓储接口。

**产出物**：

- `{module}_model.go` - GORM Model
- `{module}_command_repository.go`
- `{module}_query_repository.go`
- `{module}_repositories.go` - 聚合（多仓储时）

**规划要点**：

- Model 与 Entity 的映射策略
- 索引设计（注册到 `db/action.go`）
- 多仓储时创建 Repositories 聚合结构体

---

## Phase 3: Application 层

**目标**：实现业务用例，定义 API 契约。

**产出物**：

- `commands.go` / `queries.go` - 输入定义
- `cmd_*.go` / `qry_*.go` - Handler 实现
- `dto.go` - API 响应格式
- `mapper.go` - 实体到 DTO 转换

**规划要点**：

- 识别所有用例（CRUD + 特殊操作）
- Handler 依赖 Domain 接口（依赖倒置）

---

## Phase 4: Adapters 层

**目标**：暴露 HTTP 接口，定义路由。

**产出物**：

- `permission/constants.go` - Operation 常量
- `routes/{module}.go` - 路由元数据
- `handler/{module}.go` - HTTP Handler

**规划要点**：

- Operation 命名遵循 URN 格式（`scope:resource:action`）
- 路由分组（admin/user/public）
- 是否需要自定义中间件

---

## Phase 5: Container 层

**目标**：依赖注入装配，连接所有层。

**修改文件**：

- `repository.go` - 注册 Repositories
- `usecase.go` - 创建 UseCases 聚合
- `http.go` - 注入 Handler 和路由参数

**装配顺序**：

```
Repository → UseCase → Handler → Router
```

---

## Phase 6: 收尾工作

**产出物**：

- `seeds/{module}_seeder.go` - 种子数据
- `seeds/registry.go` - 注册 Seeder

**验证步骤**：

1. `go build -o /dev/null ./...` - 编译检查
2. `db reset` - 迁移和种子验证
3. `MANUAL=1 go test ./internal/manualtest/` - 集成测试

---

## 检查清单

| Phase       | 检查项                                           |
| ----------- | ------------------------------------------------ |
| Domain      | `doc.go` 存在、实体无 GORM Tag、接口分离         |
| Persistence | `TableName()` 方法、`Create()` 回写 ID、索引注册 |
| Application | Handler 依赖 Domain 接口、DTO 有 json tags       |
| Adapters    | Operation 常量、路由元数据、`AllRouteBindings()` |
| Container   | Repository/UseCase/Handler 注入完整              |
| 收尾        | Seeder 幂等、执行顺序正确、测试通过              |

---

## 常见陷阱

| 陷阱                 | 症状         | 解决方案                    |
| -------------------- | ------------ | --------------------------- |
| 返回类型不一致       | 编译错误     | Repository 返回 Domain 类型 |
| 缺少 Repository 参数 | nil panic    | 检查 UseCase 构造函数       |
| 依赖倒置违反         | import cycle | Handler 依赖 Domain 接口    |
| 路由未绑定           | 404          | 检查 `AllRouteBindings()`   |
| Seeder 顺序错误      | 外键约束失败 | 按依赖关系排序              |
