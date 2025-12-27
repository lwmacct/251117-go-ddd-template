---
paths:
  - "internal/infrastructure/queue/**/*.go"
---

# Queue Infrastructure 规范

## 核心职责

提供基于 Redis 的异步任务队列，采用生产者-消费者模式。

**注意**：本包不实现 Domain 接口，是通用基础设施库。

## 文件结构

| 文件             | 职责                           |
| ---------------- | ------------------------------ |
| `doc.go`         | 包文档（**必需**）             |
| `redis_queue.go` | Redis List 队列（LPUSH/BRPOP） |
| `processor.go`   | 并发工作池处理器               |

## 组件职责

| 组件         | 职责                           |
| ------------ | ------------------------------ |
| `RedisQueue` | FIFO 队列，入队/出队操作       |
| `Processor`  | 并发工作池，管理多 Worker 消费 |
| `JobHandler` | 任务处理器接口（由使用方实现） |

## 队列特性

| 特性         | 说明                         |
| ------------ | ---------------------------- |
| **FIFO**     | 先进先出（LPUSH + BRPOP）    |
| **阻塞消费** | Worker 空闲时阻塞等待        |
| **并发处理** | 支持多 Worker 并行消费       |
| **优雅关闭** | 等待所有 Worker 完成当前任务 |

## 当前限制

- ❌ 无失败重试（需在 JobHandler 中自行实现）
- ❌ 无死信队列（失败任务直接丢弃）
- ❌ 无延迟队列（需使用 Redis Sorted Set 扩展）
- ❌ 无任务去重

> 对于复杂场景，请考虑使用专业消息队列（RabbitMQ、Kafka）。
