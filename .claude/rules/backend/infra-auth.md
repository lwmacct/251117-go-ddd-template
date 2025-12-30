---
paths:
  - "internal/infrastructure/auth/**/*.go"
---

# Auth Infrastructure 规范

实现认证相关技术细节（JWT、密码哈希、会话管理）。

## 目录结构

```
internal/infrastructure/auth/
├── doc.go               # 包文档（必需）
├── {feature}_service.go # 服务实现
├── jwt.go               # JWT 处理
└── token_generator.go   # Token 生成器
```

**命名约定**:

- `{feature}` 为功能名（如 `password`、`session`）

## 设计原则

- 实现 Domain 层定义的接口
- 依赖注入技术组件
- 密码使用 bcrypt 哈希
