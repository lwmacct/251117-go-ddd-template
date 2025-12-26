---
paths:
  - "internal/**/*.go"
---

# Go Doc 文档规范

<!--TOC-->

## Table of Contents

- [语言选择](#语言选择) `:20+4`
- [包注释（doc.go）](#包注释docgo) `:24+24`
- [类型注释](#类型注释) `:48+11`
- [方法注释](#方法注释) `:59+17`
- [Go 1.19+ 文档特性](#go-119-文档特性) `:76+9`

<!--TOC-->

## 语言选择

**统一使用中文**编写文档注释，与项目整体风格保持一致。

## 包注释（doc.go）

每个 Domain 模块**必须**包含 `doc.go` 文件：

```go
// Package user 定义用户领域模型和仓储接口。
//
// 本包是用户管理的领域层核心，定义了：
//   - [User]: 用户实体（富领域模型）
//   - [CommandRepository]: 写仓储接口
//   - [QueryRepository]: 读仓储接口
//   - 用户领域错误（见 errors.go）
//
// 依赖倒置：
// 本包仅定义接口，实现位于 infrastructure/persistence 包。
package user
```

**要点**：

- 首行以 `// Package xxx` 开头，简述包职责
- 使用 `[TypeName]` 语法链接到同包类型（Go 1.19+）
- 列出包内关键类型和职责

## 类型注释

```go
// User 用户实体，包含用户基本信息和 RBAC 角色关联。
//
// 业务行为：
//   - [User.CanLogin]: 检查用户是否可登录
//   - [User.HasRole]: 检查用户是否拥有指定角色
type User struct { ... }
```

## 方法注释

```go
// HasRole 检查用户是否拥有指定角色。
func (u *User) HasRole(roleName string) bool { ... }

// CanLogin 报告用户是否可以登录。
// 当用户状态为 "active" 时返回 true。
func (u *User) CanLogin() bool { ... }
```

**要点**：

- 首句以方法名开头，使用动词描述功能
- 布尔方法使用 "报告..." 或 "检查..." 开头
- 可附加参数说明、返回值含义、错误条件

## Go 1.19+ 文档特性

| 特性         | 语法          | 示例                     |
| ------------ | ------------- | ------------------------ |
| **类型链接** | `[TypeName]`  | `参见 [User] 实体定义`   |
| **跨包链接** | `[pkg.Type]`  | `使用 [context.Context]` |
| **标题**     | `// # 标题`   | 需前后空行               |
| **列表**     | `//   - item` | 缩进 2-3 空格            |
| **代码块**   | 缩进 4 空格   | 不会被重新换行           |
