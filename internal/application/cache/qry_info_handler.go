package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// InfoHandler 缓存信息查询 Handler（类似 redis-cli INFO）。
type InfoHandler struct {
	client    *redis.Client
	keyPrefix string
}

// NewInfoHandler 创建缓存信息查询 Handler。
func NewInfoHandler(client *redis.Client, keyPrefix string) *InfoHandler {
	return &InfoHandler{
		client:    client,
		keyPrefix: keyPrefix,
	}
}

// Handle 查询缓存信息。
func (h *InfoHandler) Handle(ctx context.Context) (*CacheInfoDTO, error) {
	result := &CacheInfoDTO{
		KeyPrefix: h.keyPrefix,
	}

	// 统计应用前缀下的 key 数量
	pattern := h.keyPrefix + "*"
	count, err := h.countKeys(ctx, pattern)
	if err != nil {
		return nil, fmt.Errorf("count keys failed: %w", err)
	}
	result.DBSize = count

	// 获取 Redis 服务器信息
	info, err := h.client.Info(ctx, "server").Result()
	if err == nil {
		result.RedisVersion = extractRedisVersion(info)
	}

	// 获取内存使用量
	memInfo, err := h.client.Info(ctx, "memory").Result()
	if err == nil {
		result.MemoryUsage = extractMemoryUsage(memInfo)
	}

	return result, nil
}

// countKeys 使用 SCAN 统计匹配的 key 数量
func (h *InfoHandler) countKeys(ctx context.Context, pattern string) (int64, error) {
	var count int64
	var cursor uint64

	for {
		_, nextCursor, err := h.client.Scan(ctx, cursor, pattern, 1000).Result()
		if err != nil {
			return count, err
		}

		// SCAN 返回的是本次扫描的 key 数量
		keys, _, _ := h.client.Scan(ctx, cursor, pattern, 1000).Result()
		count += int64(len(keys))

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return count, nil
}

// extractRedisVersion 从 INFO server 输出中提取 Redis 版本
func extractRedisVersion(info string) string {
	// INFO 输出格式：redis_version:7.2.0
	const prefix = "redis_version:"
	start := 0
	for i := range info {
		if info[i:i+1] == "\n" || i == 0 {
			if i > 0 {
				start = i + 1
			}
			if len(info) > start+len(prefix) && info[start:start+len(prefix)] == prefix {
				end := start + len(prefix)
				for end < len(info) && info[end] != '\r' && info[end] != '\n' {
					end++
				}
				return info[start+len(prefix) : end]
			}
		}
	}
	return ""
}

// extractMemoryUsage 从 INFO memory 输出中提取内存使用量
func extractMemoryUsage(info string) string {
	// INFO 输出格式：used_memory_human:1.23M
	const prefix = "used_memory_human:"
	start := 0
	for i := range info {
		if info[i:i+1] == "\n" || i == 0 {
			if i > 0 {
				start = i + 1
			}
			if len(info) > start+len(prefix) && info[start:start+len(prefix)] == prefix {
				end := start + len(prefix)
				for end < len(info) && info[end] != '\r' && info[end] != '\n' {
					end++
				}
				return info[start+len(prefix) : end]
			}
		}
	}
	return ""
}
