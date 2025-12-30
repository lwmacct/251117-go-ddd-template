---
paths:
  - "internal/infrastructure/database/**/*.go"
  - "internal/command/db/**/*.go"
---

# Database Infrastructure 规范

管理数据库连接、迁移和种子数据。

## 开发环境

> [!TIP]
> 环境已预装 PostgreSQL，可直接使用 `psql` 命令, 无需配置。

## 目录结构

```
internal/infrastructure/database/
├── doc.go        # 包文档（必需）
├── connection.go # 连接管理
├── migrator.go   # 迁移执行器
└── seeder.go     # 种子执行器
```

## 迁移原则

使用 GORM AutoMigrate：

- Model 即 Schema（单一数据源）
- 记录迁移历史到 `migrations` 表

## CLI 命令

位于 `internal/command/db/`：

| 命令         | 说明                             |
| ------------ | -------------------------------- |
| `db migrate` | 执行迁移（只添加，不删除）       |
| `db reset`   | 重置数据库（删表+重建+种子数据） |

**所有 db 命令执行后会自动清空 Redis 缓存**，避免数据不一致。

## 索引管理

### 禁止临时 SQL

❌ 禁止：直接执行临时 SQL 创建索引

```sql
-- 禁止在命令行或脚本中执行
CREATE INDEX idx_xxx ON xxx(col);
```

✅ 正确：在迁移代码中定义索引

### Model 索引

定义在 `internal/command/db/action.go` 的 `getIndexMigrations()`：

```go
func getIndexMigrations() []database.IndexMigration {
    return []database.IndexMigration{
        {
            Model:   &persistence.SettingModel{},
            Indexes: []string{"idx_settings_category_sort"},
        },
    }
}
```

### 关联表索引

GORM `many2many` 关联表只创建复合主键，**不会为外键列创建单独索引**。

定义在 `internal/command/db/action.go` 的 `getJoinTableIndexes()`：

```go
func getJoinTableIndexes() []database.JoinTableIndex {
    return []database.JoinTableIndex{
        {Table: "user_roles", Name: "idx_user_roles_user_id", Columns: "user_id"},
        {Table: "user_roles", Name: "idx_user_roles_role_id", Columns: "role_id"},
        {Table: "role_permissions", Name: "idx_role_permissions_role_model_id", Columns: "role_model_id"},
        {Table: "role_permissions", Name: "idx_role_permissions_permission_model_id", Columns: "permission_model_id"},
    }
}
```

### 新增索引检查清单

1. [ ] 确定索引类型（Model 索引 / 关联表索引）
2. [ ] 在 `getIndexMigrations()` 或 `getJoinTableIndexes()` 添加配置
3. [ ] 运行 `db migrate` 或 `db reset` 应用索引

## 种子数据原则

- 实现 `Seed()` 和 `Name()` 方法
- 通过 registry 统一注册
- 幂等执行

## Seeder 性能优化

### 禁止 Association API 循环调用

❌ 禁止：循环中使用 GORM Association API

```go
// 每次 Append 触发 3 个 SQL：SELECT + UPDATE + INSERT
for _, perm := range permissions {
    db.Model(&role).Association("Permissions").Append(perm)
}
```

✅ 正确：直接批量插入关联表

```go
// 一次性批量插入，跳过 Association API 开销
records := make([]map[string]any, 0, len(perms))
for _, p := range perms {
    records = append(records, map[string]any{
        "role_model_id":       role.ID,
        "permission_model_id": p.ID,
    })
}
db.Table("role_permissions").Clauses(clause.OnConflict{DoNothing: true}).Create(&records)
```

### 优化原则

| 原则           | 说明                                              |
| -------------- | ------------------------------------------------- |
| 预加载分组     | 先查询所有数据，按需分组，避免循环内查询          |
| 直接操作关联表 | `db.Table("join_table").Create()` 跳过 ORM 开销   |
| 批量插入       | 单次 SQL 插入所有关联记录                         |
| 幂等 Upsert    | `clause.OnConflict{DoNothing: true}` 重复执行安全 |

### Association API 隐藏开销

GORM `Append()` 的实际行为：

1. `SELECT` - 加载当前关联
2. `UPDATE` - 更新主表 `updated_at`
3. `INSERT` - 插入关联记录

**循环 N 次 = 3N 个 SQL**，这是 Seeder 性能问题的常见根源。
