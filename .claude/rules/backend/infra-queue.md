---
paths:
  - "internal/infrastructure/queue/**/*.go"
---

# Queue Infrastructure 规范

基于 Redis 的异步任务队列（生产者-消费者模式）。

## 目录结构

```
internal/infrastructure/queue/
├── doc.go          # 包文档（必需）
├── redis_queue.go  # 队列实现
└── processor.go    # 处理器
```

## 队列特性

- FIFO（LPUSH + BRPOP）
- 阻塞消费
- 多 Worker 并行
- 优雅关闭

## 限制

- 无失败重试
- 无死信队列
- 无延迟队列

> 复杂场景考虑 RabbitMQ/Kafka。
