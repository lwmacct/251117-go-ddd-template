---
paths:
  - "internal/infrastructure/health/**/*.go"
---

# Health Infrastructure 规范

## 核心职责

实现 Domain 层 `health.Checker` 接口，提供系统健康检查功能。

## 文件结构

| 文件         | 职责                       |
| ------------ | -------------------------- |
| `doc.go`     | 包文档（**必需**）         |
| `checker.go` | 健康检查器，检查各组件状态 |

## 设计原则

- **接口实现**：实现 `domain/health.Checker` 接口
- **并行检查**：同时检查多个组件，减少总耗时
- **状态聚合**：汇总各组件状态，返回整体健康报告

## 健康状态

| 状态        | 含义                     |
| ----------- | ------------------------ |
| `Healthy`   | 所有组件正常             |
| `Degraded`  | 部分组件异常但核心可用   |
| `Unhealthy` | 核心组件异常，服务不可用 |

## 检查项

| 组件       | 检查内容               |
| ---------- | ---------------------- |
| `database` | PostgreSQL 连接和 Ping |
| `redis`    | Redis 连接和 Ping      |

## HTTP 端点

| 路径            | 用途                 | 响应码    |
| --------------- | -------------------- | --------- |
| `/health`       | 完整健康报告         | 200 / 503 |
| `/health/live`  | Kubernetes Liveness  | 200       |
| `/health/ready` | Kubernetes Readiness | 200 / 503 |

## 依赖方向

```
domain/health.Checker (接口)
           ↑
infrastructure/health (实现)
```
