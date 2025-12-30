---
paths:
  - "internal/infrastructure/captcha/**/*.go"
---

# Captcha Infrastructure 规范

图形验证码生成和内存存储。

## 目录结构

```
internal/infrastructure/captcha/
├── doc.go        # 包文档（必需）
├── service.go    # 服务实现
└── repository.go # 仓储实现
```

## 内存存储原则

| 特性     | 说明                 |
| -------- | -------------------- |
| 并发安全 | `sync.RWMutex`       |
| 一次性   | 验证成功后自动删除   |
| 自动清理 | 后台定期清理过期条目 |
| LRU 淘汰 | 超容量时淘汰最早条目 |
