package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	authCommand "github.com/lwmacct/251117-go-ddd-template/internal/application/auth/command"
	userQuery "github.com/lwmacct/251117-go-ddd-template/internal/application/user/query"
	domainAuth "github.com/lwmacct/251117-go-ddd-template/internal/domain/auth"
)

// AuthHandler 认证处理器（新架构）
type AuthHandler struct {
	loginHandler        *authCommand.LoginHandler
	registerHandler     *authCommand.RegisterHandler
	refreshTokenHandler *authCommand.RefreshTokenHandler
	getUserHandler      *userQuery.GetUserHandler
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(
	loginHandler *authCommand.LoginHandler,
	registerHandler *authCommand.RegisterHandler,
	refreshTokenHandler *authCommand.RefreshTokenHandler,
	getUserHandler *userQuery.GetUserHandler,
) *AuthHandler {
	return &AuthHandler{
		loginHandler:        loginHandler,
		registerHandler:     registerHandler,
		refreshTokenHandler: refreshTokenHandler,
		getUserHandler:      getUserHandler,
	}
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"max=100"`
}

// Register 用户注册
// POST /api/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 调用 Use Case Handler
	result, err := h.registerHandler.Handle(c.Request.Context(), authCommand.RegisterCommand{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		FullName: req.FullName,
	})

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{
		"message": "user registered successfully",
		"data": gin.H{
			"user_id":       result.UserID,
			"username":      result.Username,
			"email":         result.Email,
			"access_token":  result.AccessToken,
			"refresh_token": result.RefreshToken,
			"token_type":    result.TokenType,
			"expires_in":    result.ExpiresIn,
		},
	})
}

// LoginRequest 登录请求
type LoginRequest struct {
	Login     string `json:"login" binding:"required"`           // 用户名或邮箱
	Password  string `json:"password" binding:"required"`        // 密码
	CaptchaID string `json:"captcha_id" binding:"required"`      // 验证码ID
	Captcha   string `json:"captcha" binding:"required"`         // 验证码
}

// Login 用户登录
// POST /api/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 调用 Use Case Handler
	result, err := h.loginHandler.Handle(c.Request.Context(), authCommand.LoginCommand{
		Login:     req.Login,
		Password:  req.Password,
		CaptchaID: req.CaptchaID,
		Captcha:   req.Captcha,
	})

	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	// 检查是否需要 2FA
	if result.Requires2FA {
		c.JSON(200, gin.H{
			"message": "Two factor authentication required",
			"data": gin.H{
				"requires_2fa":  true,
				"session_token": result.SessionToken,
			},
		})
		return
	}

	// 正常登录成功
	c.JSON(200, gin.H{
		"message": "login successful",
		"data": gin.H{
			"access_token":  result.AccessToken,
			"refresh_token": result.RefreshToken,
			"token_type":    result.TokenType,
			"expires_in":    result.ExpiresIn,
			"user": gin.H{
				"user_id":  result.UserID,
				"username": result.Username,
			},
		},
	})
}

// RefreshTokenRequest 刷新令牌请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshToken 刷新访问令牌
// POST /api/auth/refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 调用 Use Case Handler
	result, err := h.refreshTokenHandler.Handle(c.Request.Context(), authCommand.RefreshTokenCommand{
		RefreshToken: req.RefreshToken,
	})

	if err != nil {
		// 处理特定错误
		if errors.Is(err, domainAuth.ErrInvalidToken) || errors.Is(err, domainAuth.ErrTokenExpired) {
			c.JSON(401, gin.H{"error": "invalid or expired token"})
			return
		}
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "token refreshed successfully",
		"data": gin.H{
			"access_token":  result.AccessToken,
			"refresh_token": result.RefreshToken,
			"token_type":    result.TokenType,
			"expires_in":    result.ExpiresIn,
		},
	})
}

// Me 获取当前用户信息
// GET /api/auth/me
func (h *AuthHandler) Me(c *gin.Context) {
	// 从上下文获取用户信息 (由 JWT 中间件设置)
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	userID, ok := userIDRaw.(uint)
	if !ok {
		c.JSON(500, gin.H{"error": "invalid user_id type"})
		return
	}

	// 调用 Query Handler 获取完整用户信息
	user, err := h.getUserHandler.Handle(c.Request.Context(), userQuery.GetUserQuery{
		UserID:    userID,
		WithRoles: true, // 包含角色信息
	})

	if err != nil {
		c.JSON(404, gin.H{"error": "user not found"})
		return
	}

	c.JSON(200, gin.H{
		"data": user,
	})
}
