// Package redis 提供 Redis 客户端和缓存服务。
//
// # 客户端管理
//
// [NewClient] 创建 Redis 客户端连接：
//   - 支持单节点和集群模式
//   - 连接池管理
//   - 自动重连
//
// 使用示例：
//
//	client, err := redis.NewClient(cfg.Redis)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer client.Close()
//
// # 缓存仓储
//
// 遵循 CQRS 模式提供读写分离的仓储：
//   - [NewCacheCommandRepository]: 写仓储（Set/Delete/SetNX）
//   - [NewCacheQueryRepository]: 读仓储（Get/Exists）
//
// 使用示例：
//
//	cmdRepo := redis.NewCacheCommandRepository(client, "myapp:")
//	queryRepo := redis.NewCacheQueryRepository(client, "myapp:")
//	err := cmdRepo.Set(ctx, "user:1", user, time.Hour)
//	var cached User
//	err = queryRepo.Get(ctx, "user:1", &cached)
//
// # 键前缀约定
//
// 推荐的键命名规范：
//   - 用户相关：user:{id}
//   - 权限缓存：perm:{user_id}
//   - 会话相关：session:{token}
//   - 验证码：captcha:{id}
//   - 令牌刷新：refresh:{token}
//
// # 配置项
//
// 通过 config.Redis 配置：
//   - Host: Redis 主机
//   - Port: Redis 端口
//   - Password: 密码（可选）
//   - DB: 数据库编号（默认 0）
//   - PoolSize: 连接池大小
//   - MinIdleConns: 最小空闲连接数
package redis
