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
│   └── memory_bus.go            # 内存事件总线
└── eventhandler/                # 事件处理器
    ├── audit_log.go             # 审计日志处理器
    └── cache_invalidation.go    # 缓存失效处理器
```

## EventBus 规范

```go
// eventbus/memory_bus.go - 实现 domain/event.EventBus 接口
type InMemoryEventBus struct {
    handlers map[string][]event.EventHandler
    mu       sync.RWMutex
}

func NewInMemoryEventBus() *InMemoryEventBus

// 发布事件（同步执行所有匹配的处理器）
func (b *InMemoryEventBus) Publish(ctx context.Context, events ...event.Event) error

// 订阅事件（支持通配符：user.* 或 *）
func (b *InMemoryEventBus) Subscribe(eventName string, handler event.EventHandler)
```

### 通配符匹配规则

| 模式           | 匹配                           |
| -------------- | ------------------------------ |
| `user.*`       | `user.created`, `user.deleted` |
| `*`            | 所有事件                       |
| `user.created` | 精确匹配                       |

## EventHandler 规范

```go
// eventhandler/xxx.go - 实现 domain/event.EventHandler 接口
type AuditLogHandler struct {
    auditLogRepo auditlog.CommandRepository
    logger       *slog.Logger
}

func NewAuditLogHandler(auditLogRepo auditlog.CommandRepository) *AuditLogHandler

// Handle 根据事件类型分发处理
func (h *AuditLogHandler) Handle(ctx context.Context, e event.Event) error
```

## 错误处理原则

- 审计日志写入失败**不应阻塞**业务流程
- 缓存失效失败应**记录日志**但不返回错误

## 依赖方向

```
domain/event.EventBus     (接口)
domain/event.EventHandler (接口)
              ↑
infrastructure/eventbus     (EventBus 实现)
infrastructure/eventhandler (EventHandler 实现)
```
