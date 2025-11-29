package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	redisinfra "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/redis"
	"github.com/redis/go-redis/v9"
)

// CacheHandler 缓存处理器示例
type CacheHandler struct {
	cacheRepo redisinfra.CacheRepository
}

// NewCacheHandler 创建缓存处理器
// keyPrefix: Redis key 前缀，例如 "myapp:"
func NewCacheHandler(redisClient *redis.Client, keyPrefix string) *CacheHandler {
	return &CacheHandler{
		cacheRepo: redisinfra.NewCacheRepository(redisClient, keyPrefix),
	}
}

// SetCache 设置缓存示例
// POST /api/cache
// Body: {"key": "test", "value": "hello", "ttl": 60}
func (h *CacheHandler) SetCache(c *gin.Context) {
	var req struct {
		Key   string `json:"key" binding:"required"`
		Value any    `json:"value" binding:"required"`
		TTL   int    `json:"ttl"` // 秒，默认 60
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 默认 TTL 60 秒
	if req.TTL == 0 {
		req.TTL = 60
	}

	if err := h.cacheRepo.Set(c.Request.Context(), req.Key, req.Value, time.Duration(req.TTL)*time.Second); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "cache set successfully",
		"key":     req.Key,
		"ttl":     req.TTL,
	})
}

// GetCache 获取缓存示例
// GET /api/cache/:key
func (h *CacheHandler) GetCache(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.JSON(400, gin.H{"error": "key is required"})
		return
	}

	var value any
	if err := h.cacheRepo.Get(c.Request.Context(), key, &value); err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"key":   key,
		"value": value,
	})
}

// DeleteCache 删除缓存示例
// DELETE /api/cache/:key
func (h *CacheHandler) DeleteCache(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.JSON(400, gin.H{"error": "key is required"})
		return
	}

	if err := h.cacheRepo.Delete(c.Request.Context(), key); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "cache deleted successfully",
		"key":     key,
	})
}
