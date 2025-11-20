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
//
// @Summary      初始化两步验证
// @Description  为当前用户生成2FA密钥和二维码，用于配置身份验证器应用
// @Tags         认证 - 两步验证 (Authentication - 2FA)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} response.Response{data=object{secret=string,qr_code=string}} "2FA初始化成功"
// @Failure      400 {object} response.ErrorResponse "设置失败"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Router       /api/auth/2fa/setup [post]
// @x-permission {"scope":"user:2fa:setup"}
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
	Code string `json:"code" binding:"required" example:"123456"` // TOTP 验证码
}

// VerifyAndEnable 验证 TOTP 代码并启用 2FA
//
// @Summary      验证并启用两步验证
// @Description  验证身份验证器应用生成的TOTP代码，成功后启用2FA并返回恢复代码
// @Tags         认证 - 两步验证 (Authentication - 2FA)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body VerifyAndEnableRequest true "TOTP验证码"
// @Success      200 {object} response.Response{data=object{recovery_codes=[]string,message=string}} "2FA启用成功"
// @Failure      400 {object} response.ErrorResponse "验证码错误"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Router       /api/auth/2fa/verify [post]
// @x-permission {"scope":"user:2fa:setup"}
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
//
// @Summary      禁用两步验证
// @Description  禁用当前用户的两步验证功能
// @Tags         认证 - 两步验证 (Authentication - 2FA)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} response.Response{message=string} "2FA禁用成功"
// @Failure      400 {object} response.ErrorResponse "禁用失败"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Router       /api/auth/2fa/disable [post]
// @x-permission {"scope":"user:2fa:disable"}
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
//
// @Summary      获取两步验证状态
// @Description  获取当前用户的2FA启用状态和剩余恢复代码数量
// @Tags         认证 - 两步验证 (Authentication - 2FA)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} response.Response{data=object{enabled=bool,recovery_codes_count=int}} "2FA状态"
// @Failure      400 {object} response.ErrorResponse "获取失败"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Router       /api/auth/2fa/status [get]
// @x-permission {"scope":"user:2fa:read"}
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
