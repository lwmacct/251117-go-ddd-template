package handler

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/database"
	redisinfra "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/redis"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// HealthHandler 健康检查处理器
type HealthHandler struct {
	db          *gorm.DB
	redisClient *redis.Client
}

// NewHealthHandler 创建健康检查处理器
func NewHealthHandler(db *gorm.DB, redisClient *redis.Client) *HealthHandler {
	return &HealthHandler{
		db:          db,
		redisClient: redisClient,
	}
}

// Check 执行健康检查
func (h *HealthHandler) Check(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	response := gin.H{
		"status": "ok",
		"checks": gin.H{},
	}

	statusCode := 200
	allHealthy := true

	// 检查数据库连接
	if err := database.HealthCheck(ctx, h.db); err != nil {
		response["checks"].(gin.H)["database"] = gin.H{
			"status": "unhealthy",
			"error":  err.Error(),
		}
		allHealthy = false
	} else {
		// 获取数据库连接池统计
		stats, _ := database.GetStats(h.db)
		response["checks"].(gin.H)["database"] = gin.H{
			"status": "healthy",
			"stats":  stats,
		}
	}

	// 检查 Redis 连接
	if err := redisinfra.HealthCheck(ctx, h.redisClient); err != nil {
		response["checks"].(gin.H)["redis"] = gin.H{
			"status": "unhealthy",
			"error":  err.Error(),
		}
		allHealthy = false
	} else {
		response["checks"].(gin.H)["redis"] = gin.H{
			"status": "healthy",
		}
	}

	if !allHealthy {
		statusCode = 503
		response["status"] = "degraded"
	}

	c.JSON(statusCode, response)
}
