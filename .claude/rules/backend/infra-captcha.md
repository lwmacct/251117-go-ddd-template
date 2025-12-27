---
paths:
  - "internal/infrastructure/captcha/**/*.go"
---

# Captcha Infrastructure 规范

## 核心职责

实现 Domain 层 `captcha.Service`、`captcha.CommandRepository` 和 `captcha.QueryRepository` 接口，提供图形验证码生成和存储功能。

**特点**：使用内存存储（非数据库），支持 LRU 淘汰和自动清理。

## 文件结构

| 文件            | 职责                               |
| --------------- | ---------------------------------- |
| `doc.go`        | 包文档（**必需**）                 |
| `service.go`    | 验证码生成服务，使用 base64Captcha |
| `repository.go` | 内存存储仓储，实现 Command/Query   |

## 设计原则

### Service 原则

- 实现 `domain/captcha.Service` 接口
- 输出 Base64 编码的 PNG 图片，便于前端直接展示
- 支持字符型验证码（a-z, 0-9）

### Repository 原则（内存存储）

| 特性         | 说明                            |
| ------------ | ------------------------------- |
| **并发安全** | 使用 `sync.RWMutex` 保护        |
| **一次性**   | 验证成功后自动删除              |
| **自动清理** | 后台 goroutine 定期清理过期条目 |
| **LRU 淘汰** | 超过容量上限时淘汰最早条目      |
| **大小写**   | 验证时忽略大小写                |

### 配置建议

- 最大存储条目：10000
- 默认过期时间：5 分钟
- 清理间隔：10 分钟

## 依赖方向

```
domain/captcha.Service          (接口)
domain/captcha.CommandRepository (接口)
domain/captcha.QueryRepository   (接口)
              ↑
infrastructure/captcha (实现)
```
