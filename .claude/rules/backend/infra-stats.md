---
paths:
  - "internal/infrastructure/stats/**/*.go"
---

# Stats Infrastructure 规范

## 核心职责

实现 Domain 层 `stats.QueryRepository` 接口，提供跨域聚合查询能力。

**特点**：

- **只读操作**，无写入
- 直接使用 GORM `Table()` 查询，无独立 Model
- 聚合多个表的统计数据

## 文件结构

| 文件                  | 职责               |
| --------------------- | ------------------ |
| `doc.go`              | 包文档（**必需**） |
| `query_repository.go` | 聚合查询仓储实现   |

## 设计原则

- **接口实现**：实现 `domain/stats.QueryRepository` 接口
- **跨表聚合**：单次调用聚合多个表的统计数据
- **只读设计**：无 CommandRepository，仅提供查询

## 查询方法

| 方法                    | 说明             |
| ----------------------- | ---------------- |
| `GetSystemStats(limit)` | 完整系统统计     |
| `GetUserCountByStatus`  | 按状态统计用户数 |
| `GetTotalUsers`         | 用户总数         |
| `GetTotalRoles`         | 角色总数         |
| `GetTotalPermissions`   | 权限总数         |
| `GetTotalMenus`         | 菜单总数         |
| `GetRecentAuditLogs`    | 最近审计日志     |

## 聚合数据源

| 表            | 统计内容                 |
| ------------- | ------------------------ |
| `users`       | 用户数（按 status 分类） |
| `roles`       | 角色总数                 |
| `permissions` | 权限总数                 |
| `menus`       | 菜单总数                 |
| `audit_logs`  | 最近审计日志摘要         |

## 性能建议

- 高频调用场景考虑添加缓存
- 大数据量场景考虑分页或采样

## 依赖方向

```
domain/stats.QueryRepository (接口)
              ↑
infrastructure/stats (实现)
```
