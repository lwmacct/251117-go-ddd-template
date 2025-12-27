---
paths:
  - "internal/infrastructure/eventbus/**/*.go"
  - "internal/infrastructure/eventhandler/**/*.go"
---

# Event Infrastructure 规范

## 核心职责

实现事件驱动架构：事件总线（发布/订阅）和事件处理器（业务响应）。

## 目录结构

```
internal/infrastructure/
├── eventbus/                    # 事件总线实现
│   ├── doc.go                   # 包文档（必需）
│   └── memory_bus.go            # 内存事件总线
└── eventhandler/                # 事件处理器
    ├── doc.go                   # 包文档（必需）
    ├── audit_log.go             # 审计日志处理器
    └── cache_invalidation.go    # 缓存失效处理器
```

## EventBus 原则

- **接口实现**：实现 `domain/event.EventBus` 接口
- **同步执行**：`Publish()` 同步执行所有匹配的处理器
- **通配符订阅**：支持 `user.*`（前缀匹配）和 `*`（全部匹配）

### 通配符匹配规则

| 模式           | 匹配                           |
| -------------- | ------------------------------ |
| `user.*`       | `user.created`, `user.deleted` |
| `*`            | 所有事件                       |
| `user.created` | 精确匹配                       |

## EventHandler 原则

- **接口实现**：实现 `domain/event.EventHandler` 接口
- **错误隔离**：审计日志写入失败**不应阻塞**业务流程
- **日志记录**：缓存失效失败应**记录日志**但不返回错误

## 依赖方向

```
domain/event.EventBus     (接口)
domain/event.EventHandler (接口)
              ↑
infrastructure/eventbus     (EventBus 实现)
infrastructure/eventhandler (EventHandler 实现)
```
