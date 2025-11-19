package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/pat"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/auth"
)

// PATHandler handles Personal Access Token operations
type PATHandler struct {
	patService *auth.PATService
}

// NewPATHandler creates a new PAT handler
func NewPATHandler(patService *auth.PATService) *PATHandler {
	return &PATHandler{
		patService: patService,
	}
}

// CreateToken creates a new Personal Access Token
//
//	@Summary		Create personal access token
//	@Description	Create a new PAT with selected permissions
//	@Tags			tokens
//	@Accept			json
//	@Produce		json
//	@Param			request	body		pat.CreateTokenRequest	true	"Token creation request"
//	@Success		201		{object}	pat.TokenResponse
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		401		{object}	map[string]interface{}
//	@Router			/user/tokens [post]
func (h *PATHandler) CreateToken(c *gin.Context) {
	var req pat.CreateTokenRequest

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

	req.UserID = userID.(uint)

	// Create token
	resp, err := h.patService.CreateToken(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to create token: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "token created successfully",
		"data":    resp,
		"warning": "Please save this token now. You won't be able to see it again!",
	})
}

// ListTokens lists all tokens for the current user
//
//	@Summary		List personal access tokens
//	@Description	Get all PATs for the current user (without sensitive data)
//	@Tags			tokens
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}
//	@Failure		401	{object}	map[string]interface{}
//	@Router			/user/tokens [get]
func (h *PATHandler) ListTokens(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized: user ID not found",
		})
		return
	}

	tokens, err := h.patService.ListTokens(c.Request.Context(), userID.(uint))
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
//	@Summary		Revoke personal access token
//	@Description	Revoke a PAT by ID
//	@Tags			tokens
//	@Produce		json
//	@Param			id	path		int	true	"Token ID"
//	@Success		200	{object}	map[string]interface{}
//	@Failure		400	{object}	map[string]interface{}
//	@Failure		401	{object}	map[string]interface{}
//	@Failure		404	{object}	map[string]interface{}
//	@Router			/user/tokens/{id} [delete]
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

	if err := h.patService.RevokeToken(c.Request.Context(), userID.(uint), uint(tokenID)); err != nil {
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
//	@Summary		Get personal access token details
//	@Description	Get details of a PAT by ID (without sensitive data)
//	@Tags			tokens
//	@Produce		json
//	@Param			id	path		int	true	"Token ID"
//	@Success		200	{object}	map[string]interface{}
//	@Failure		400	{object}	map[string]interface{}
//	@Failure		401	{object}	map[string]interface{}
//	@Failure		404	{object}	map[string]interface{}
//	@Router			/user/tokens/{id} [get]
func (h *PATHandler) GetToken(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized: user ID not found",
		})
		return
	}

	// Get all user's tokens and find the specific one
	tokens, err := h.patService.ListTokens(c.Request.Context(), userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to retrieve token: " + err.Error(),
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

	// Find the token
	for _, token := range tokens {
		if token.ID == uint(tokenID) {
			c.JSON(http.StatusOK, gin.H{
				"message": "token retrieved successfully",
				"data":    token,
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"error": "token not found",
	})
}
