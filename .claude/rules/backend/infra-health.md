---
paths:
  - "internal/infrastructure/health/**/*.go"
---

# Health Infrastructure 规范

系统健康检查。

## 文件命名

| 文件类型 | 命名规范             |
| -------- | -------------------- |
| 包文档   | `doc.go`（**必需**） |
| 检查器   | `checker.go`         |

## 设计原则

- 并行检查多个组件
- 状态聚合：Healthy / Degraded / Unhealthy

## HTTP 端点

| 路径            | 用途                 |
| --------------- | -------------------- |
| `/health`       | 完整健康报告         |
| `/health/live`  | Kubernetes Liveness  |
| `/health/ready` | Kubernetes Readiness |
