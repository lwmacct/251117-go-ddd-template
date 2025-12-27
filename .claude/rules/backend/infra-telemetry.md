---
paths:
  - "internal/infrastructure/telemetry/**/*.go"
---

# Telemetry Infrastructure 规范

## 核心职责

提供 OpenTelemetry 分布式追踪初始化和配置。

**注意**：本包不实现 Domain 接口，是横切关注点（Cross-Cutting Concern）。

## 文件结构

| 文件      | 职责                           |
| --------- | ------------------------------ |
| `doc.go`  | 包文档（**必需**）             |
| `otel.go` | OpenTelemetry SDK 初始化和配置 |

## 配置项

| 字段           | 类型    | 说明                                |
| -------------- | ------- | ----------------------------------- |
| `ServiceName`  | string  | 必填，服务名称                      |
| `Enabled`      | bool    | 是否启用追踪                        |
| `ExporterType` | string  | `otlp` / `stdout` / `none`          |
| `OTLPEndpoint` | string  | OTLP gRPC 端点，如 `localhost:4317` |
| `SampleRate`   | float64 | 采样率 0.0-1.0                      |

## Exporter 类型

| 类型     | 用途     | 说明                   |
| -------- | -------- | ---------------------- |
| `otlp`   | 生产环境 | 导出到 Jaeger/Tempo 等 |
| `stdout` | 开发调试 | 输出到控制台           |
| `none`   | 禁用     | 空操作                 |

## 设计原则

- **最先初始化**：Telemetry 应在所有其他基础设施之前初始化
- **优雅关闭**：`shutdown()` 函数有 5 秒超时，确保 span 导出完成
- **自动传播**：配置 W3C TraceContext + Baggage 传播器

## 初始化顺序

```
Telemetry → Database → Redis → EventBus → Repositories → ...
```

## 与其他模块集成

| 模块     | 集成方式                          |
| -------- | --------------------------------- |
| Database | 配置 `EnableTracing` 参数         |
| Redis    | 构造函数传入 `enableTracing` 参数 |
| HTTP     | 使用 OTEL Gin 中间件              |
