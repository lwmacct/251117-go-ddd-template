package cache

import (
	"context"
	"strconv"

	"github.com/redis/go-redis/v9"
)

// ScanKeysQuery SCAN 查询参数
type ScanKeysQuery struct {
	// Pattern 匹配模式（不含应用前缀，如 "setting:*"）
	Pattern string

	// Cursor 游标（用于分页），首次查询传 "0"
	Cursor string

	// Limit 每次返回的最大数量
	Limit int
}

// ScanKeysHandler 缓存 Key 扫描 Handler（类似 redis-cli SCAN）。
type ScanKeysHandler struct {
	client    *redis.Client
	keyPrefix string
}

// NewScanKeysHandler 创建缓存 Key 扫描 Handler。
func NewScanKeysHandler(client *redis.Client, keyPrefix string) *ScanKeysHandler {
	return &ScanKeysHandler{
		client:    client,
		keyPrefix: keyPrefix,
	}
}

// Handle 扫描缓存 Keys。
func (h *ScanKeysHandler) Handle(ctx context.Context, q ScanKeysQuery) (*ScanKeysResultDTO, error) {
	// 解析游标
	cursor, err := strconv.ParseUint(q.Cursor, 10, 64)
	if err != nil {
		cursor = 0
	}

	// 设置默认 limit
	limit := q.Limit
	if limit <= 0 {
		limit = 100
	}

	// 构建完整 pattern
	pattern := h.keyPrefix
	if q.Pattern != "" {
		pattern += q.Pattern
	} else {
		pattern += "*"
	}

	// 执行 SCAN
	keys, nextCursor, err := h.client.Scan(ctx, cursor, pattern, int64(limit)).Result()
	if err != nil {
		return nil, err
	}

	// 获取每个 key 的详细信息
	keyDTOs := make([]CacheKeyDTO, 0, len(keys))
	for _, k := range keys {
		dto := CacheKeyDTO{Key: k}

		// 获取 TTL
		ttl, err := h.client.TTL(ctx, k).Result()
		if err == nil {
			dto.TTL = int64(ttl.Seconds())
		}

		// 获取类型
		keyType, err := h.client.Type(ctx, k).Result()
		if err == nil {
			dto.Type = keyType
		}

		keyDTOs = append(keyDTOs, dto)
	}

	return &ScanKeysResultDTO{
		Keys:         keyDTOs,
		Cursor:       strconv.FormatUint(nextCursor, 10),
		TotalScanned: int64(len(keys)),
	}, nil
}
