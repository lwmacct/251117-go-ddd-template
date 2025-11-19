package handler

import (
	"github.com/gin-gonic/gin"

	twofaService "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/twofa"
)

// TwoFAHandler 2FA 处理器
type TwoFAHandler struct {
	twofaService *twofaService.Service
}

// NewTwoFAHandler 创建 2FA 处理器
func NewTwoFAHandler(twofaService *twofaService.Service) *TwoFAHandler {
	return &TwoFAHandler{
		twofaService: twofaService,
	}
}

// Setup 设置 2FA（生成密钥和二维码）
// POST /api/auth/2fa/setup
func (h *TwoFAHandler) Setup(c *gin.Context) {
	ctx := c.Request.Context()

	// 从 JWT 中间件获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	result, err := h.twofaService.Setup(ctx, userID.(uint))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "2FA setup initiated",
		"data":    result,
	})
}

// VerifyAndEnableRequest 验证并启用 2FA 请求
type VerifyAndEnableRequest struct {
	Code string `json:"code" binding:"required"` // TOTP 验证码
}

// VerifyAndEnable 验证 TOTP 代码并启用 2FA
// POST /api/auth/2fa/verify
func (h *TwoFAHandler) VerifyAndEnable(c *gin.Context) {
	ctx := c.Request.Context()

	// 从 JWT 中间件获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	var req VerifyAndEnableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request: " + err.Error()})
		return
	}

	recoveryCodes, err := h.twofaService.VerifyAndEnable(ctx, userID.(uint), req.Code)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "2FA enabled successfully",
		"data": gin.H{
			"recovery_codes": recoveryCodes,
			"message":        "Please save these recovery codes in a safe place. You won't be able to see them again.",
		},
	})
}

// Disable 禁用 2FA
// POST /api/auth/2fa/disable
func (h *TwoFAHandler) Disable(c *gin.Context) {
	ctx := c.Request.Context()

	// 从 JWT 中间件获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.twofaService.Disable(ctx, userID.(uint)); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "2FA disabled successfully",
	})
}

// GetStatus 获取 2FA 状态
// GET /api/auth/2fa/status
func (h *TwoFAHandler) GetStatus(c *gin.Context) {
	ctx := c.Request.Context()

	// 从 JWT 中间件获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	enabled, recoveryCodesCount, err := h.twofaService.GetStatus(ctx, userID.(uint))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "success",
		"data": gin.H{
			"enabled":              enabled,
			"recovery_codes_count": recoveryCodesCount,
		},
	})
}
