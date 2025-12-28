package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// ErrKeyNotFound key 不存在错误
var ErrKeyNotFound = errors.New("key not found")

// GetKeyHandler 获取单个 Key 值的 Handler（类似 redis-cli GET/JSON.GET）。
type GetKeyHandler struct {
	client    *redis.Client
	keyPrefix string
}

// NewGetKeyHandler 创建获取 Key 值的 Handler。
func NewGetKeyHandler(client *redis.Client, keyPrefix string) *GetKeyHandler {
	return &GetKeyHandler{
		client:    client,
		keyPrefix: keyPrefix,
	}
}

// Handle 获取单个 Key 的值。
//
// key 参数是完整的 key 名称（含应用前缀）。
func (h *GetKeyHandler) Handle(ctx context.Context, key string) (*CacheValueDTO, error) {
	// 检查 key 是否存在
	exists, err := h.client.Exists(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("check key exists failed: %w", err)
	}
	if exists == 0 {
		return nil, ErrKeyNotFound
	}

	result := &CacheValueDTO{
		CacheKeyDTO: CacheKeyDTO{Key: key},
	}

	// 获取类型
	keyType, err := h.client.Type(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("get key type failed: %w", err)
	}
	result.Type = keyType

	// 获取 TTL
	ttl, err := h.client.TTL(ctx, key).Result()
	if err == nil {
		result.TTL = int64(ttl.Seconds())
	}

	// 根据类型获取值
	var value any
	switch keyType {
	case "string":
		value, err = h.client.Get(ctx, key).Result()
	case "ReJSON-RL":
		// RedisJSON 类型
		jsonStr, jerr := h.client.JSONGet(ctx, key, "$").Result()
		if jerr == nil {
			// JSON.GET $ 返回数组包装
			value = json.RawMessage(jsonStr)
		}
		err = jerr
	case "hash":
		value, err = h.client.HGetAll(ctx, key).Result()
	case "list":
		value, err = h.client.LRange(ctx, key, 0, -1).Result()
	case "set":
		value, err = h.client.SMembers(ctx, key).Result()
	case "zset":
		value, err = h.client.ZRangeWithScores(ctx, key, 0, -1).Result()
	default:
		value = "unsupported type: " + keyType
	}

	if err != nil {
		return nil, fmt.Errorf("get key value failed: %w", err)
	}

	// 序列化值为 JSON
	valueBytes, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("marshal value failed: %w", err)
	}
	result.Value = valueBytes

	return result, nil
}
