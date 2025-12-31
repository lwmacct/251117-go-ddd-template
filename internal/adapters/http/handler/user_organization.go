package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/response"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/organization"
)

// UserOrganizationHandler 用户视角的组织/团队 Handler
type UserOrganizationHandler struct {
	userOrgsHandler  *organization.UserOrgsHandler
	userTeamsHandler *organization.UserTeamsHandler
}

// NewUserOrganizationHandler 创建用户视角组织 Handler
func NewUserOrganizationHandler(
	userOrgsHandler *organization.UserOrgsHandler,
	userTeamsHandler *organization.UserTeamsHandler,
) *UserOrganizationHandler {
	return &UserOrganizationHandler{
		userOrgsHandler:  userOrgsHandler,
		userTeamsHandler: userTeamsHandler,
	}
}

// ListMyOrganizations 获取我加入的组织列表
//
//	@Summary		我的组织
//	@Description	获取当前用户加入的所有组织
//	@Tags			User - Organization
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response.DataResponse[[]organization.UserOrgDTO]	"组织列表"
//	@Failure		401	{object}	response.ErrorResponse								"未授权"
//	@Failure		500	{object}	response.ErrorResponse								"服务器内部错误"
//	@Router			/api/user/organizations [get]
func (h *UserOrganizationHandler) ListMyOrganizations(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "user not authenticated")
		return
	}
	userID, ok := userIDVal.(uint)
	if !ok {
		response.InternalError(c, "invalid user ID type")
		return
	}

	result, err := h.userOrgsHandler.Handle(c.Request.Context(), organization.ListUserOrgsQuery{
		UserID: userID,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, "success", result)
}

// ListUserTeamsQuery 用户团队列表查询参数
type ListUserTeamsQuery struct {
	OrgID uint `form:"org_id" binding:"omitempty,min=1"`
}

// ListMyTeams 获取我加入的团队列表
//
//	@Summary		我的团队
//	@Description	获取当前用户加入的所有团队（可按组织筛选）
//	@Tags			User - Organization
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			params	query		handler.ListUserTeamsQuery							false	"查询参数"
//	@Success		200		{object}	response.DataResponse[[]organization.UserTeamDTO]	"团队列表"
//	@Failure		401		{object}	response.ErrorResponse								"未授权"
//	@Failure		500		{object}	response.ErrorResponse								"服务器内部错误"
//	@Router			/api/user/teams [get]
func (h *UserOrganizationHandler) ListMyTeams(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "user not authenticated")
		return
	}
	userID, ok := userIDVal.(uint)
	if !ok {
		response.InternalError(c, "invalid user ID type")
		return
	}

	// 解析可选的 org_id 参数
	var orgID uint
	if orgIDStr := c.Query("org_id"); orgIDStr != "" {
		id, err := strconv.ParseUint(orgIDStr, 10, 32)
		if err != nil {
			response.BadRequest(c, "invalid org_id")
			return
		}
		orgID = uint(id)
	}

	result, err := h.userTeamsHandler.Handle(c.Request.Context(), organization.ListUserTeamsQuery{
		UserID: userID,
		OrgID:  orgID,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, "success", result)
}
