package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	patdto "github.com/lwmacct/251117-go-ddd-template/internal/application/pat"
	patCommand "github.com/lwmacct/251117-go-ddd-template/internal/application/pat/command"
	patQuery "github.com/lwmacct/251117-go-ddd-template/internal/application/pat/query"
)

// PATHandler handles Personal Access Token operations (DDD+CQRS Use Case Pattern)
type PATHandler struct {
	// Command Handlers
	createTokenHandler *patCommand.CreateTokenHandler
	revokeTokenHandler *patCommand.RevokeTokenHandler

	// Query Handlers
	getTokenHandler   *patQuery.GetTokenHandler
	listTokensHandler *patQuery.ListTokensHandler
}

// NewPATHandler creates a new PAT handler
func NewPATHandler(
	createTokenHandler *patCommand.CreateTokenHandler,
	revokeTokenHandler *patCommand.RevokeTokenHandler,
	getTokenHandler *patQuery.GetTokenHandler,
	listTokensHandler *patQuery.ListTokensHandler,
) *PATHandler {
	return &PATHandler{
		createTokenHandler: createTokenHandler,
		revokeTokenHandler: revokeTokenHandler,
		getTokenHandler:    getTokenHandler,
		listTokensHandler:  listTokensHandler,
	}
}

// CreateToken creates a new Personal Access Token
//
// @Summary      创建个人访问令牌
// @Description  用户创建新的个人访问令牌(PAT)，用于API访问。令牌仅在创建时显示一次
// @Tags         用户 - 个人访问令牌 (User - Personal Access Token)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body patdto.CreateTokenRequest true "令牌信息"
// @Success      201 {object} response.Response{data=object{id=uint,name=string,token=string,expires_at=string},warning=string} "令牌创建成功"
// @Failure      400 {object} response.ErrorResponse "参数错误"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Router       /api/user/pats [post]
// @x-permission {"scope":"user:pats:create"}
func (h *PATHandler) CreateToken(c *gin.Context) {
	var req patdto.CreateTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request: " + err.Error(),
		})
		return
	}

	// Get user ID from context (set by Auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized: user ID not found",
		})
		return
	}

	// 调用 Use Case Handler
	result, err := h.createTokenHandler.Handle(c.Request.Context(), patCommand.CreateTokenCommand{
		UserID:      userID.(uint),
		Name:        req.Name,
		Permissions: req.Permissions,
		ExpiresAt:   nil, // TODO: 处理过期时间
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to create token: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "token created successfully",
		"data":    result,
		"warning": "Please save this token now. You won't be able to see it again!",
	})
}

// ListTokens lists all tokens for the current user
//
// @Summary      获取个人访问令牌列表
// @Description  获取当前用户的所有个人访问令牌（不包含令牌值）
// @Tags         用户 - 个人访问令牌 (User - Personal Access Token)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} response.Response{data=[]object{id=uint,name=string,last_used_at=string,expires_at=string,created_at=string},count=int} "令牌列表"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Failure      500 {object} response.ErrorResponse "服务器内部错误"
// @Router       /api/user/pats [get]
// @x-permission {"scope":"user:pats:read"}
func (h *PATHandler) ListTokens(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized: user ID not found",
		})
		return
	}

	// 调用 Use Case Handler
	tokens, err := h.listTokensHandler.Handle(c.Request.Context(), patQuery.ListTokensQuery{
		UserID: userID.(uint),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to list tokens: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "tokens retrieved successfully",
		"data":    tokens,
		"count":   len(tokens),
	})
}

// RevokeToken revokes a specific token
//
// @Summary      撤销个人访问令牌
// @Description  用户撤销（删除）指定的个人访问令牌
// @Tags         用户 - 个人访问令牌 (User - Personal Access Token)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "令牌ID" minimum(1)
// @Success      200 {object} response.Response{message=string} "令牌撤销成功"
// @Failure      400 {object} response.ErrorResponse "无效的令牌ID"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Failure      404 {object} response.ErrorResponse "令牌不存在"
// @Router       /api/user/pats/{id} [delete]
// @x-permission {"scope":"user:pats:delete"}
func (h *PATHandler) RevokeToken(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized: user ID not found",
		})
		return
	}

	tokenIDStr := c.Param("id")
	tokenID, err := strconv.ParseUint(tokenIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid token ID",
		})
		return
	}

	// 调用 Use Case Handler
	err = h.revokeTokenHandler.Handle(c.Request.Context(), patCommand.RevokeTokenCommand{
		UserID:  userID.(uint),
		TokenID: uint(tokenID),
	})

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "failed to revoke token: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "token revoked successfully",
	})
}

// GetToken retrieves details of a specific token
//
// @Summary      获取个人访问令牌详情
// @Description  获取指定个人访问令牌的详细信息（不包含令牌值）
// @Tags         用户 - 个人访问令牌 (User - Personal Access Token)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "令牌ID" minimum(1)
// @Success      200 {object} response.Response{data=object{id=uint,name=string,last_used_at=string,expires_at=string,created_at=string}} "令牌详情"
// @Failure      400 {object} response.ErrorResponse "无效的令牌ID"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Failure      404 {object} response.ErrorResponse "令牌不存在"
// @Router       /api/user/pats/{id} [get]
// @x-permission {"scope":"user:pats:read"}
func (h *PATHandler) GetToken(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized: user ID not found",
		})
		return
	}

	tokenIDStr := c.Param("id")
	tokenID, err := strconv.ParseUint(tokenIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid token ID",
		})
		return
	}

	// 调用 Use Case Handler
	token, err := h.getTokenHandler.Handle(c.Request.Context(), patQuery.GetTokenQuery{
		UserID:  userID.(uint),
		TokenID: uint(tokenID),
	})

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "token not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "token retrieved successfully",
		"data":    token,
	})
}
