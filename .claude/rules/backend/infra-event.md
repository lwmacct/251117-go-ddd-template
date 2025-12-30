---
paths:
  - "internal/infrastructure/eventbus/**/*.go"
  - "internal/infrastructure/eventhandler/**/*.go"
---

# Event Infrastructure 规范

事件驱动架构：事件总线（发布/订阅）+ 事件处理器。

## 目录结构

```
internal/infrastructure/
├── eventbus/
│   └── {type}_bus.go   # 事件总线实现
└── eventhandler/
    └── {feature}.go    # 事件处理器
```

**命名约定**:

- `{type}` 为总线类型（如 `sync`、`async`）
- `{feature}` 为功能名（如 `audit_log`、`notification`）

## EventBus 原则

- 实现 Domain 层 EventBus 接口
- `Publish()` 同步执行处理器
- 支持通配符：`user.*`（前缀）、`*`（全部）

## EventHandler 原则

- 实现 Domain 层 EventHandler 接口
- 错误隔离：失败不阻塞业务流程
- 失败记录日志但不返回错误
