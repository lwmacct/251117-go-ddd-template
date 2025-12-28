---
paths:
  - "internal/infrastructure/telemetry/**/*.go"
---

# Telemetry Infrastructure 规范

## 核心职责

OpenTelemetry 分布式追踪（横切关注点）。

## 文件命名

| 文件类型 | 命名规范             |
| -------- | -------------------- |
| 包文档   | `doc.go`（**必需**） |
| 初始化   | `otel.go`            |

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
