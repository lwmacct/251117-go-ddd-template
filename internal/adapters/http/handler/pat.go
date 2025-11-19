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
	getTokenHandler  *patQuery.GetTokenHandler
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
