package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/response"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/twofa"
)

// TwoFAHandler 2FA 处理器
type TwoFAHandler struct {
	setupHandler        *twofa.SetupHandler
	verifyEnableHandler *twofa.VerifyEnableHandler
	disableHandler      *twofa.DisableHandler
	getStatusHandler    *twofa.GetStatusHandler
}

// NewTwoFAHandler 创建 2FA 处理器
func NewTwoFAHandler(
	setupHandler *twofa.SetupHandler,
	verifyEnableHandler *twofa.VerifyEnableHandler,
	disableHandler *twofa.DisableHandler,
	getStatusHandler *twofa.GetStatusHandler,
) *TwoFAHandler {
	return &TwoFAHandler{
		setupHandler:        setupHandler,
		verifyEnableHandler: verifyEnableHandler,
		disableHandler:      disableHandler,
		getStatusHandler:    getStatusHandler,
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
// @Success      200 {object} response.DataResponse[twofa.SetupDTO] "2FA初始化成功"
// @Failure      400 {object} response.ErrorResponse "设置失败"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Router       /api/auth/2fa/setup [post]
// @x-permission {"scope":"user:2fa:setup"}
func (h *TwoFAHandler) Setup(c *gin.Context) {
	ctx := c.Request.Context()

	userID, ok := getUserID(c)
	if !ok {
		return
	}

	result, err := h.setupHandler.Handle(ctx, twofa.SetupCommand{
		UserID: userID,
	})
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 转换为 DTO
	resp := twofa.SetupDTO{
		Secret:    result.Secret,
		QRCodeURL: result.QRCodeURL,
		QRCodeImg: result.QRCodeImg,
	}
	response.OK(c, "2FA setup initiated", resp)
}

// VerifyAndEnable 验证 TOTP 代码并启用 2FA
//
// @Summary      验证并启用两步验证
// @Description  验证身份验证器应用生成的TOTP代码，成功后启用2FA并返回恢复代码
// @Tags         认证 - 两步验证 (Authentication - 2FA)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body twofa.VerifyDTO true "TOTP验证码"
// @Success      200 {object} response.DataResponse[twofa.EnableDTO] "2FA启用成功"
// @Failure      400 {object} response.ErrorResponse "验证码错误"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Router       /api/auth/2fa/verify [post]
// @x-permission {"scope":"user:2fa:setup"}
func (h *TwoFAHandler) VerifyAndEnable(c *gin.Context) {
	ctx := c.Request.Context()

	userID, ok := getUserID(c)
	if !ok {
		return
	}

	var req twofa.VerifyDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result, err := h.verifyEnableHandler.Handle(ctx, twofa.VerifyEnableCommand{
		UserID: userID,
		Code:   req.Code,
	})
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	resp := twofa.EnableDTO{
		RecoveryCodes: result.RecoveryCodes,
		Message:       "Please save these recovery codes in a safe place. You won't be able to see them again.",
	}
	response.OK(c, "2FA enabled successfully", resp)
}

// Disable 禁用 2FA
//
// @Summary      禁用两步验证
// @Description  禁用当前用户的两步验证功能
// @Tags         认证 - 两步验证 (Authentication - 2FA)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} response.MessageResponse "2FA禁用成功"
// @Failure      400 {object} response.ErrorResponse "禁用失败"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Router       /api/auth/2fa/disable [post]
// @x-permission {"scope":"user:2fa:disable"}
func (h *TwoFAHandler) Disable(c *gin.Context) {
	ctx := c.Request.Context()

	userID, ok := getUserID(c)
	if !ok {
		return
	}

	if err := h.disableHandler.Handle(ctx, twofa.DisableCommand{
		UserID: userID,
	}); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, "2FA disabled successfully", nil)
}

// GetStatus 获取 2FA 状态
//
// @Summary      获取两步验证状态
// @Description  获取当前用户的2FA启用状态和剩余恢复代码数量
// @Tags         认证 - 两步验证 (Authentication - 2FA)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} response.DataResponse[twofa.StatusDTO] "2FA状态"
// @Failure      400 {object} response.ErrorResponse "获取失败"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Router       /api/auth/2fa/status [get]
// @x-permission {"scope":"user:2fa:read"}
func (h *TwoFAHandler) GetStatus(c *gin.Context) {
	ctx := c.Request.Context()

	userID, ok := getUserID(c)
	if !ok {
		return
	}

	result, err := h.getStatusHandler.Handle(ctx, twofa.GetStatusQuery{
		UserID: userID,
	})
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	resp := twofa.StatusDTO{
		Enabled:            result.Enabled,
		RecoveryCodesCount: result.RecoveryCodesCount,
	}
	response.OK(c, "success", resp)
}

// getUserID 从上下文获取用户ID，并输出统一未认证响应
func getUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "Authentication required")
		return 0, false
	}

	id, ok := userID.(uint)
	if !ok {
		response.Unauthorized(c, "Invalid user context")
		return 0, false
	}

	return id, true
}
