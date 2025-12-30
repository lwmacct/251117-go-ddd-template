---
paths:
  - "internal/infrastructure/telemetry/**/*.go"
---

# Telemetry Infrastructure 规范

OpenTelemetry 分布式追踪（横切关注点）。

## 目录结构

```
internal/infrastructure/telemetry/
├── doc.go   # 包文档（必需）
└── otel.go  # 初始化
```

## 设计原则

- **最先初始化**：在其他基础设施之前
- **优雅关闭**：确保 span 导出完成
- **自动传播**：W3C TraceContext + Baggage

## Exporter 类型

| 类型     | 用途     |
| -------- | -------- |
| `otlp`   | 生产环境 |
| `stdout` | 开发调试 |
| `none`   | 禁用     |
