---
paths:
  - "internal/infrastructure/captcha/**/*.go"
---

# Captcha Infrastructure 规范

图形验证码生成和内存存储。

## 文件命名

| 文件类型 | 命名规范             |
| -------- | -------------------- |
| 包文档   | `doc.go`（**必需**） |
| 服务实现 | `service.go`         |
| 仓储实现 | `repository.go`      |

## 内存存储原则

| 特性     | 说明                 |
| -------- | -------------------- |
| 并发安全 | `sync.RWMutex`       |
| 一次性   | 验证成功后自动删除   |
| 自动清理 | 后台定期清理过期条目 |
| LRU 淘汰 | 超容量时淘汰最早条目 |
