---
paths:
  - "internal/infrastructure/auth/**/*.go"
---

# Auth Infrastructure 规范

## 核心职责

实现认证相关技术细节（JWT、密码哈希、会话管理）。

## 文件命名

| 文件类型     | 命名规范             |
| ------------ | -------------------- |
| 包文档       | `doc.go`（**必需**） |
| 服务实现     | `{功能}_service.go`  |
| JWT 处理     | `jwt.go`             |
| Token 生成器 | `token_generator.go` |

## 设计原则

- 实现 Domain 层定义的接口
- 依赖注入技术组件
- 密码使用 bcrypt 哈希
