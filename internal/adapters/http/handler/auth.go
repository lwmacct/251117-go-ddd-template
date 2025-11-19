package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/auth"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService *auth.Service
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(authService *auth.Service) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register 用户注册
// POST /api/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req auth.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{
		"message": "user registered successfully",
		"data":    resp,
	})
}

// Login 用户登录
// POST /api/auth/login
// 支持两种登录流程：
// 1. 第一次登录：login + password + captcha_id + captcha
// 2. 第二次登录（2FA）：session_token + two_factor_code
func (h *AuthHandler) Login(c *gin.Context) {
	var req auth.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		// 检查是否是2FA验证需求错误
		if _, ok := err.(*auth.TwoFactorRequiredError); ok {
			// 返回session token，前端需要显示2FA验证页面
			c.JSON(200, gin.H{
				"message": "Two factor authentication required",
				"data":    resp,
			})
			return
		}
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "login successful",
		"data":    resp,
	})
}

// RefreshToken 刷新访问令牌
// POST /api/auth/refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "token refreshed successfully",
		"data":    resp,
	})
}

// Me 获取当前用户信息
// GET /api/auth/me
func (h *AuthHandler) Me(c *gin.Context) {
	// 从上下文获取用户信息 (由 JWT 中间件设置)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	username, _ := c.Get("username")
	email, _ := c.Get("email")

	c.JSON(200, gin.H{
		"data": gin.H{
			"user_id":  userID,
			"username": username,
			"email":    email,
		},
	})
}
