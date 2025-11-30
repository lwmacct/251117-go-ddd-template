package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/response"
	appAuth "github.com/lwmacct/251117-go-ddd-template/internal/application/auth"
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

// Register 用户注册
//
// @Summary      用户注册
// @Description  创建新用户账号，注册成功后自动登录并返回访问令牌
// @Tags         认证 (Authentication)
// @Accept       json
// @Produce      json
// @Param        request body appAuth.RegisterDTO true "注册信息"
// @Success      201 {object} response.Response{data=object{user_id=uint,username=string,email=string,access_token=string,refresh_token=string,token_type=string,expires_in=int}} "注册成功"
// @Failure      400 {object} response.ErrorResponse "参数错误或用户名/邮箱已存在"
// @Router       /api/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req appAuth.RegisterDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
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
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, "user registered successfully", gin.H{
		"user_id":       result.UserID,
		"username":      result.Username,
		"email":         result.Email,
		"access_token":  result.AccessToken,
		"refresh_token": result.RefreshToken,
		"token_type":    result.TokenType,
		"expires_in":    result.ExpiresIn,
	})
}

// Login 用户登录
//
// @Summary      用户登录
// @Description  使用手机号/用户名/邮箱和密码登录系统，需要提供图形验证码。如果启用了2FA，返回session_token用于后续2FA验证
// @Tags         认证 (Authentication)
// @Accept       json
// @Produce      json
// @Param        request body appAuth.LoginDTO true "登录凭证"
// @Success      200 {object} response.Response{data=object{access_token=string,refresh_token=string,token_type=string,expires_in=int,user=object{user_id=uint,username=string}}} "登录成功"
// @Success      200 {object} response.Response{data=object{requires_2fa=bool,session_token=string}} "需要2FA验证"
// @Failure      401 {object} response.ErrorResponse "登录失败：凭证无效、验证码错误或账户被禁用"
// @Router       /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req appAuth.LoginDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	// 调用 Use Case Handler
	result, err := h.loginHandler.Handle(c.Request.Context(), authCommand.LoginCommand{
		Account:   req.Account, // 前端使用 account，内部使用 login
		Password:  req.Password,
		CaptchaID: req.CaptchaID,
		Captcha:   req.Captcha,
	})

	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	// 检查是否需要 2FA
	if result.Requires2FA {
		response.OK(c, "Two factor authentication required", gin.H{
			"requires_2fa":  true,
			"session_token": result.SessionToken,
		})
		return
	}

	// 正常登录成功
	response.OK(c, "login successful", gin.H{
		"access_token":  result.AccessToken,
		"refresh_token": result.RefreshToken,
		"token_type":    result.TokenType,
		"expires_in":    result.ExpiresIn,
		"user": gin.H{
			"user_id":  result.UserID,
			"username": result.Username,
		},
	})
}

// RefreshToken 刷新访问令牌
//
// @Summary      刷新访问令牌
// @Description  使用refresh_token获取新的access_token和refresh_token，延长会话有效期
// @Tags         认证 (Authentication)
// @Accept       json
// @Produce      json
// @Param        request body appAuth.RefreshTokenDTO true "刷新令牌"
// @Success      200 {object} response.Response{data=object{access_token=string,refresh_token=string,token_type=string,expires_in=int}} "令牌刷新成功"
// @Failure      401 {object} response.ErrorResponse "刷新令牌无效或已过期"
// @Router       /api/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req appAuth.RefreshTokenDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	// 调用 Use Case Handler
	result, err := h.refreshTokenHandler.Handle(c.Request.Context(), authCommand.RefreshTokenCommand{
		RefreshToken: req.RefreshToken,
	})

	if err != nil {
		// 处理特定错误
		if errors.Is(err, domainAuth.ErrInvalidToken) || errors.Is(err, domainAuth.ErrTokenExpired) {
			response.Unauthorized(c, "invalid or expired token")
			return
		}
		response.Unauthorized(c, err.Error())
		return
	}

	response.OK(c, "token refreshed successfully", gin.H{
		"access_token":  result.AccessToken,
		"refresh_token": result.RefreshToken,
		"token_type":    result.TokenType,
		"expires_in":    result.ExpiresIn,
	})
}

// Me 获取当前用户信息
//
// @Summary      获取当前用户信息
// @Description  获取当前登录用户的详细信息，包括角色和权限
// @Tags         认证 (Authentication)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} response.Response{data=object{id=uint,username=string,email=string,full_name=string,status=string,roles=array}} "用户信息"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Failure      404 {object} response.ErrorResponse "用户不存在"
// @Router       /api/auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	// 从上下文获取用户信息 (由 JWT 中间件设置)
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "unauthorized")
		return
	}

	userID, ok := userIDRaw.(uint)
	if !ok {
		response.InternalError(c, "invalid user_id type")
		return
	}

	// 调用 Query Handler 获取完整用户信息
	user, err := h.getUserHandler.Handle(c.Request.Context(), userQuery.GetUserQuery{
		UserID:    userID,
		WithRoles: true, // 包含角色信息
	})

	if err != nil {
		response.NotFound(c, "user")
		return
	}

	response.OK(c, "success", user)
}
