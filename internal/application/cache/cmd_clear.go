package cache

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

// DeleteHandler 缓存删除 Handler（类似 redis-cli DEL）。
type DeleteHandler struct {
	client    *redis.Client
	keyPrefix string
}

// NewDeleteHandler 创建缓存删除 Handler。
func NewDeleteHandler(client *redis.Client, keyPrefix string) *DeleteHandler {
	return &DeleteHandler{
		client:    client,
		keyPrefix: keyPrefix,
	}
}

// DeleteKey 删除单个 Key。
//
// key 参数是完整的 key 名称（含应用前缀）。
func (h *DeleteHandler) DeleteKey(ctx context.Context, key string) (*DeleteResultDTO, error) {
	deleted, err := h.client.Del(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("delete key failed: %w", err)
	}

	slog.Info("cache key deleted", "key", key)

	return &DeleteResultDTO{
		DeletedCount: deleted,
	}, nil
}

// DeleteByPattern 按 pattern 批量删除 Keys。
//
// pattern 参数不含应用前缀（如 "setting:*"）。
func (h *DeleteHandler) DeleteByPattern(ctx context.Context, pattern string) (*DeleteResultDTO, error) {
	fullPattern := h.keyPrefix + pattern

	var deleted int64
	var cursor uint64

	for {
		keys, nextCursor, err := h.client.Scan(ctx, cursor, fullPattern, 100).Result()
		if err != nil {
			return nil, fmt.Errorf("scan keys failed: %w", err)
		}

		if len(keys) > 0 {
			n, err := h.client.Del(ctx, keys...).Result()
			if err != nil {
				slog.Warn("delete keys failed", "pattern", fullPattern, "err", err)
			}
			deleted += n
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	slog.Info("cache keys deleted by pattern", "pattern", fullPattern, "count", deleted)

	return &DeleteResultDTO{
		DeletedCount: deleted,
	}, nil
}
