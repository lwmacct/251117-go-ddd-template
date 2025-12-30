---
paths:
  - ".claude/rules/**/*.md"
---

# 规范文件编写指南

## 文件结构

```markdown
---
paths:
  - "{匹配路径}/**/*.{ext}"
---

# {规范名称}

{一句话描述}

## 目录结构

## 命名规范（如有 precommit 测试）

## 禁止事项（可选）

## {具体规范}...
```

## 目录结构格式

````markdown
```
{path}/
├── file.go         # 说明
├── {template}.go   # 占位符用 {}
└── doc.go          # （必需）/（可选）标注
```

**命名约定**:

- `{template}` 占位符说明
````

## 命名规范格式

仅当有 precommit 测试时添加：

```markdown
| 文件      | 类型   | 规则          | 示例      |
| --------- | ------ | ------------- | --------- |
| `file.go` | struct | 以 `Xxx` 结尾 | `UserXxx` |

> 规范由 `internal/precommit/xxx_test.go` 自动检查
```

## 代码示例格式

```go
// ✅ 正确
response.OK(c, "success", data)

// ❌ 禁止
c.JSON(200, gin.H{...})
```
