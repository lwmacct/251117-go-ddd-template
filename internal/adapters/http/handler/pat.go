package handler

import (
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/response"
	patdto "github.com/lwmacct/251117-go-ddd-template/internal/application/pat"
	patCommand "github.com/lwmacct/251117-go-ddd-template/internal/application/pat/command"
	patQuery "github.com/lwmacct/251117-go-ddd-template/internal/application/pat/query"
)

// PATHandler handles Personal Access Token operations (DDD+CQRS Use Case Pattern)
type PATHandler struct {
	// Command Handlers
	createTokenHandler  *patCommand.CreateTokenHandler
	deleteTokenHandler  *patCommand.DeleteTokenHandler
	disableTokenHandler *patCommand.DisableTokenHandler
	enableTokenHandler  *patCommand.EnableTokenHandler

	// Query Handlers
	getTokenHandler   *patQuery.GetTokenHandler
	listTokensHandler *patQuery.ListTokensHandler
}

// NewPATHandler creates a new PAT handler
func NewPATHandler(
	createTokenHandler *patCommand.CreateTokenHandler,
	deleteTokenHandler *patCommand.DeleteTokenHandler,
	disableTokenHandler *patCommand.DisableTokenHandler,
	enableTokenHandler *patCommand.EnableTokenHandler,
	getTokenHandler *patQuery.GetTokenHandler,
	listTokensHandler *patQuery.ListTokensHandler,
) *PATHandler {
	return &PATHandler{
		createTokenHandler:  createTokenHandler,
		deleteTokenHandler:  deleteTokenHandler,
		disableTokenHandler: disableTokenHandler,
		enableTokenHandler:  enableTokenHandler,
		getTokenHandler:     getTokenHandler,
		listTokensHandler:   listTokensHandler,
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
// @Router       /api/user/tokens [post]
// @x-permission {"scope":"user:tokens:create"}
func (h *PATHandler) CreateToken(c *gin.Context) {
	var req patdto.CreateTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body", err.Error())
		return
	}

	// Get user ID from context (set by Auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "unauthorized: user ID not found")
		return
	}

	uid, ok := userID.(uint)
	if !ok {
		response.Unauthorized(c, "unauthorized: invalid user ID type")
		return
	}

	expiresAt, err := parseExpiresAt(req.ExpiresAt, req.ExpiresIn)
	if err != nil {
		response.BadRequest(c, "invalid expiration date", err.Error())
		return
	}

	// 调用 Use Case Handler
	result, err := h.createTokenHandler.Handle(c.Request.Context(), patCommand.CreateTokenCommand{
		UserID:      uid,
		Name:        req.Name,
		Permissions: req.Permissions,
		ExpiresAt:   expiresAt,
		IPWhitelist: req.IPWhitelist,
		Description: req.Description,
	})

	if err != nil {
		response.BadRequest(c, "failed to create token", err.Error())
		return
	}

	response.Created(c, "token created successfully", patdto.ToCreateTokenResponse(result.Token, result.PlainToken))
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
// @Router       /api/user/tokens [get]
// @x-permission {"scope":"user:tokens:read"}
func (h *PATHandler) ListTokens(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "unauthorized: user ID not found")
		return
	}

	uid, ok := userID.(uint)
	if !ok {
		response.Unauthorized(c, "unauthorized: invalid user ID type")
		return
	}

	// 调用 Use Case Handler
	tokens, err := h.listTokensHandler.Handle(c.Request.Context(), patQuery.ListTokensQuery{
		UserID: uid,
	})

	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, "tokens retrieved successfully", tokens)
}

// DeleteToken deletes a specific token
//
// @Summary      删除个人访问令牌
// @Description  用户删除指定的个人访问令牌（不可恢复）
// @Tags         用户 - 个人访问令牌 (User - Personal Access Token)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "令牌ID" minimum(1)
// @Success      200 {object} response.Response{message=string} "令牌删除成功"
// @Failure      400 {object} response.ErrorResponse "无效的令牌ID"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Failure      404 {object} response.ErrorResponse "令牌不存在"
// @Router       /api/user/tokens/{id} [delete]
// @x-permission {"scope":"user:tokens:delete"}
func (h *PATHandler) DeleteToken(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "unauthorized: user ID not found")
		return
	}

	uid, ok := userID.(uint)
	if !ok {
		response.Unauthorized(c, "unauthorized: invalid user ID type")
		return
	}

	tokenIDStr := c.Param("id")
	tokenID, err := parseTokenID(tokenIDStr)
	if err != nil {
		response.BadRequest(c, "invalid token ID", err.Error())
		return
	}

	err = h.deleteTokenHandler.Handle(c.Request.Context(), patCommand.DeleteTokenCommand{
		UserID:  uid,
		TokenID: tokenID,
	})
	if err != nil {
		response.BadRequest(c, "failed to delete token", err.Error())
		return
	}

	response.OK(c, "token deleted successfully", nil)
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
// @Router       /api/user/tokens/{id} [get]
// @x-permission {"scope":"user:tokens:read"}
func (h *PATHandler) GetToken(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "unauthorized: user ID not found")
		return
	}

	uid, ok := userID.(uint)
	if !ok {
		response.Unauthorized(c, "unauthorized: invalid user ID type")
		return
	}

	tokenIDStr := c.Param("id")
	tokenID, err := strconv.ParseUint(tokenIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid token ID", nil)
		return
	}

	// 调用 Use Case Handler
	token, err := h.getTokenHandler.Handle(c.Request.Context(), patQuery.GetTokenQuery{
		UserID:  uid,
		TokenID: uint(tokenID),
	})

	if err != nil {
		response.NotFound(c, "token")
		return
	}

	response.OK(c, "token retrieved successfully", token)
}

// DisableToken 暂停令牌
//
// @Summary      禁用个人访问令牌
// @Description  暂停指定令牌的使用（可再次启用）
// @Tags         用户 - 个人访问令牌 (User - Personal Access Token)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "令牌ID" minimum(1)
// @Success      200 {object} response.Response{message=string} "令牌已禁用"
// @Failure      400 {object} response.ErrorResponse "无效的令牌ID"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Router       /api/user/tokens/{id}/disable [patch]
// @x-permission {"scope":"user:tokens:disable"}
func (h *PATHandler) DisableToken(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "unauthorized: user ID not found")
		return
	}

	uid, ok := userID.(uint)
	if !ok {
		response.Unauthorized(c, "unauthorized: invalid user ID type")
		return
	}

	tokenID, err := parseTokenID(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid token ID", err.Error())
		return
	}

	if err := h.disableTokenHandler.Handle(c.Request.Context(), patCommand.DisableTokenCommand{
		UserID:  uid,
		TokenID: tokenID,
	}); err != nil {
		response.BadRequest(c, "failed to disable token", err.Error())
		return
	}

	response.OK(c, "token disabled successfully", nil)
}

// EnableToken 启用令牌
//
// @Summary      启用个人访问令牌
// @Description  重新启用已禁用的令牌
// @Tags         用户 - 个人访问令牌 (User - Personal Access Token)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "令牌ID" minimum(1)
// @Success      200 {object} response.Response{message=string} "令牌已启用"
// @Failure      400 {object} response.ErrorResponse "无效的令牌ID"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Router       /api/user/tokens/{id}/enable [patch]
// @x-permission {"scope":"user:tokens:enable"}
func (h *PATHandler) EnableToken(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "unauthorized: user ID not found")
		return
	}

	uid, ok := userID.(uint)
	if !ok {
		response.Unauthorized(c, "unauthorized: invalid user ID type")
		return
	}

	tokenID, err := parseTokenID(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid token ID", err.Error())
		return
	}

	if err := h.enableTokenHandler.Handle(c.Request.Context(), patCommand.EnableTokenCommand{
		UserID:  uid,
		TokenID: tokenID,
	}); err != nil {
		response.BadRequest(c, "failed to enable token", err.Error())
		return
	}

	response.OK(c, "token enabled successfully", nil)
}

func parseTokenID(raw string) (uint, error) {
	id, err := strconv.ParseUint(raw, 10, 32)
	return uint(id), err
}

// parseExpiresAt parses expire parameters from request into a timestamp pointer.
func parseExpiresAt(expiresAt *string, expiresIn *int) (*time.Time, error) {
	// 优先处理 expiresAt 字符串
	if expiresAt != nil && *expiresAt != "" {
		return parseExpiresAtString(*expiresAt)
	}

	// 次之处理 expiresIn 天数
	if expiresIn != nil && *expiresIn > 0 {
		t := time.Now().Add(time.Duration(*expiresIn) * 24 * time.Hour).UTC()
		return &t, nil
	}

	return nil, nil //nolint:nilnil // returns nil for not found, valid pattern
}

// parseExpiresAtString 解析过期时间字符串
func parseExpiresAtString(expiresAt string) (*time.Time, error) {
	// 尝试 RFC3339 格式
	parsed, err := time.Parse(time.RFC3339, expiresAt)
	if err == nil {
		utc := parsed.UTC()
		return &utc, nil
	}

	// 尝试本地时区格式（前端 datetime-local）
	localTime, localErr := time.ParseInLocation("2006-01-02T15:04", expiresAt, time.Local)
	if localErr == nil {
		utc := localTime.UTC()
		return &utc, nil
	}

	return nil, errors.New("expires_at must be RFC3339 or yyyy-MM-ddTHH:mm")
}
