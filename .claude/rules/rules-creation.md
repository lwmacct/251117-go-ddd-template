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

## Frontmatter 规范

```yaml
---
paths:
  - "internal/application/**/*.go" # 单路径
  - "internal/infrastructure/**/*.go" # 多路径
---
```

- `paths` 定义规范文件的触发路径（编辑匹配路径时自动加载）
- 使用 glob 模式匹配

## 目录结构章节

### 格式

````markdown
## 目录结构

```
internal/{layer}/{module}/
├── file1.go        # 说明
├── {template}.go   # 模板文件（占位符用 {}）
└── doc.go          # 包文档
```

**命名约定**:

- `{template}` 占位符说明
- 多实体时的变体规则
````

### 要点

- 使用树形图展示目录层级
- 注释内联在文件名后（`#` 分隔）
- 可选文件标注 `（可选）`
- 必需文件标注 `（必需）`

### 命名约定

仅当目录结构包含 `{占位符}` 模板时添加：

```markdown
**命名约定**:

- `{module}` 为领域模块名（如 `user`、`role`）
- `{action}` 为操作名（如 `create`、`update`）
```

## 命名规范章节

**仅当存在 precommit 测试检查时添加**。

### 格式

```markdown
## 命名规范

| 文件          | 类型   | 规则                   | 示例                |
| ------------- | ------ | ---------------------- | ------------------- |
| `commands.go` | struct | 以 `Command` 结尾      | `CreateUserCommand` |
| `dto.go`      | struct | 以 `DTO` 结尾          | `UserDTO`           |
| `mapper.go`   | func   | `To` 开头 + `DTO` 结尾 | `ToUserDTO`         |

> 规范由 `internal/precommit/xxx_test.go` 自动检查
```

### 要点

- 表格列：文件、类型（struct/func）、规则、示例
- 底部引用 precommit 测试文件路径
- 让开发者知道规则有自动化检查

## 代码示例格式

### 正确/禁止对比

````markdown
```go
// ✅ 正确
response.OK(c, "success", data)

// ❌ 禁止
c.JSON(200, gin.H{...})
```
````

````

### 禁止事项列表

```markdown
## 禁止事项

- ❌ 在 Handler 中编排业务逻辑
- ❌ 直接调用 Repository
````

## 表格格式

使用 Markdown 表格，列对齐由 linter 自动处理：

```markdown
| 列1 | 列2 | 列3 |
| --- | --- | --- |
| 值1 | 值2 | 值3 |
```

## 检查清单

创建规范文件时确认：

- [ ] Frontmatter `paths` 正确匹配目标文件
- [ ] `## 目录结构` 使用树形图格式
- [ ] 模板占位符有 `**命名约定**` 说明
- [ ] 如有 precommit 测试，添加 `## 命名规范` 表格
- [ ] 代码示例使用 ✅/❌ 对比格式
