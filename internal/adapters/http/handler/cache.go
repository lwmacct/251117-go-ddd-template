package handler

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/response"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/cache"
)

// CacheHandler 缓存处理器示例
type CacheHandler struct {
	setHandler    *cache.SetHandler
	getHandler    *cache.GetHandler
	deleteHandler *cache.DeleteHandler
}

// NewCacheHandler 创建缓存处理器
func NewCacheHandler(
	setHandler *cache.SetHandler,
	getHandler *cache.GetHandler,
	deleteHandler *cache.DeleteHandler,
) *CacheHandler {
	return &CacheHandler{
		setHandler:    setHandler,
		getHandler:    getHandler,
		deleteHandler: deleteHandler,
	}
}

// SetCache 设置缓存示例
//
// @Summary      设置缓存
// @Description  设置缓存键值对（演示用，公开接口）
// @Tags         缓存示例 (Cache Demo)
// @Accept       json
// @Produce      json
// @Param        request body cache.SetDTO true "缓存数据"
// @Success      200 {object} response.DataResponse[cache.SetResultDTO] "设置成功"
// @Failure      400 {object} response.ErrorResponse "参数错误"
// @Failure      500 {object} response.ErrorResponse "服务器内部错误"
// @Router       /api/cache [post]
func (h *CacheHandler) SetCache(c *gin.Context) {
	var req cache.SetDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 默认 TTL 60 秒
	ttl := req.TTL
	if ttl == 0 {
		ttl = 60
	}

	err := h.setHandler.Handle(c.Request.Context(), cache.SetCommand{
		Key:   req.Key,
		Value: req.Value,
		TTL:   time.Duration(ttl) * time.Second,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, "cache set successfully", &cache.SetResultDTO{
		Key: req.Key,
		TTL: ttl,
	})
}

// GetCache 获取缓存示例
//
// @Summary      获取缓存
// @Description  根据 key 获取缓存值（演示用，公开接口）
// @Tags         缓存示例 (Cache Demo)
// @Accept       json
// @Produce      json
// @Param        key path string true "缓存键"
// @Success      200 {object} response.DataResponse[cache.GetResultDTO] "获取成功"
// @Failure      404 {object} response.ErrorResponse "缓存不存在"
// @Router       /api/cache/{key} [get]
func (h *CacheHandler) GetCache(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		response.BadRequest(c, "key is required")
		return
	}

	result, err := h.getHandler.Handle(c.Request.Context(), cache.GetQuery{
		Key: key,
	})
	if err != nil {
		response.NotFound(c, "cache key")
		return
	}

	response.OK(c, "success", result)
}

// DeleteCache 删除缓存示例
//
// @Summary      删除缓存
// @Description  根据 key 删除缓存（演示用，公开接口）
// @Tags         缓存示例 (Cache Demo)
// @Accept       json
// @Produce      json
// @Param        key path string true "缓存键"
// @Success      200 {object} response.DataResponse[cache.DeleteResultDTO] "删除成功"
// @Failure      500 {object} response.ErrorResponse "服务器内部错误"
// @Router       /api/cache/{key} [delete]
func (h *CacheHandler) DeleteCache(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		response.BadRequest(c, "key is required")
		return
	}

	err := h.deleteHandler.Handle(c.Request.Context(), cache.DeleteCommand{
		Key: key,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, "cache deleted successfully", &cache.DeleteResultDTO{
		Key: key,
	})
}
