---
paths:
  - "internal/infrastructure/stats/**/*.go"
---

# Stats Infrastructure 规范

## 核心职责

跨域聚合查询（只读）。

## 文件命名

| 文件类型 | 命名规范              |
| -------- | --------------------- |
| 包文档   | `doc.go`（**必需**）  |
| 查询仓储 | `query_repository.go` |

## 设计原则

- 只读设计，无 CommandRepository
- 直接使用 GORM `Table()` 查询
- 聚合多表统计数据

## 性能建议

- 高频调用添加缓存
- 大数据量分页或采样
