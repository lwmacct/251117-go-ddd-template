---
paths:
  - "internal/infrastructure/stats/**/*.go"
---

# Stats Infrastructure 规范

跨域聚合查询（只读）。

## 目录结构

```
internal/infrastructure/stats/
├── doc.go              # 包文档（必需）
└── query_repository.go # 查询仓储
```

## 设计原则

- 只读设计，无 CommandRepository
- 直接使用 GORM `Table()` 查询
- 聚合多表统计数据

## 性能建议

- 高频调用添加缓存
- 大数据量分页或采样
