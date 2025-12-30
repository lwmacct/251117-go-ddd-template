---
paths:
  - "internal/infrastructure/validation/**/*.go"
---

# Validation Infrastructure 规范

基于 JSON Logic 的动态规则验证。

## 目录结构

```
internal/infrastructure/validation/
├── doc.go                  # 包文档（必需）
└── jsonlogic_validator.go  # 验证器
```

## 规则格式

### JSON Logic（推荐）

```json
{">=": [{"var": "value"}, 6]}
{"and": [{">=": [...]}, {"<=": [...]}]}
```

### 简单格式

| 字段        | 说明     |
| ----------- | -------- |
| `required`  | 是否必填 |
| `min`/`max` | 数值范围 |
| `minLength` | 最小长度 |

## 设计原则

- 规则存储在数据库，支持运行时修改
- 支持跨字段验证
