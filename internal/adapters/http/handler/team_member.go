package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/response"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/organization"
)

// ListTeamMembersQuery 团队成员列表查询参数
type ListTeamMembersQuery struct {
	response.PaginationQueryDTO
}

// ToQuery 转换为 Application 层 Query 对象
func (q *ListTeamMembersQuery) ToQuery(orgID, teamID uint) organization.ListTeamMembersQuery {
	return organization.ListTeamMembersQuery{
		OrgID:  orgID,
		TeamID: teamID,
		Offset: q.GetOffset(),
		Limit:  q.GetLimit(),
	}
}

// TeamMemberHandler 团队成员管理 Handler
type TeamMemberHandler struct {
	// Command Handlers
	addMemberHandler    *organization.TeamMemberAddHandler
	removeMemberHandler *organization.TeamMemberRemoveHandler

	// Query Handlers
	listMembersHandler *organization.TeamMemberListHandler
}

// NewTeamMemberHandler 创建团队成员管理 Handler
func NewTeamMemberHandler(
	addMemberHandler *organization.TeamMemberAddHandler,
	removeMemberHandler *organization.TeamMemberRemoveHandler,
	listMembersHandler *organization.TeamMemberListHandler,
) *TeamMemberHandler {
	return &TeamMemberHandler{
		addMemberHandler:    addMemberHandler,
		removeMemberHandler: removeMemberHandler,
		listMembersHandler:  listMembersHandler,
	}
}

// List 获取团队成员列表
//
//	@Summary		团队成员列表
//	@Description	分页获取团队成员列表
//	@Tags			Organization - Team Member Management
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			org_id	path		int														true	"组织ID"	minimum(1)
//	@Param			team_id	path		int														true	"团队ID"	minimum(1)
//	@Param			params	query		handler.ListTeamMembersQuery							false	"查询参数"
//	@Success		200		{object}	response.PagedResponse[organization.TeamMemberDTO]		"成员列表"
//	@Failure		400		{object}	response.ErrorResponse									"无效的ID"
//	@Failure		401		{object}	response.ErrorResponse									"未授权"
//	@Failure		403		{object}	response.ErrorResponse									"权限不足"
//	@Failure		500		{object}	response.ErrorResponse									"服务器内部错误"
//	@Router			/api/org/{org_id}/teams/{team_id}/members [get]
func (h *TeamMemberHandler) List(c *gin.Context) {
	orgID, err := strconv.ParseUint(c.Param("org_id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid organization ID")
		return
	}

	teamID, err := strconv.ParseUint(c.Param("team_id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid team ID")
		return
	}

	var q ListTeamMembersQuery
	if bindErr := c.ShouldBindQuery(&q); bindErr != nil {
		response.ValidationError(c, bindErr.Error())
		return
	}

	result, err := h.listMembersHandler.Handle(c.Request.Context(), q.ToQuery(uint(orgID), uint(teamID)))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	meta := response.NewPaginationMeta(int(result.Total), q.GetPage(), q.GetLimit())
	response.List(c, "success", result.Items, meta)
}

// Add 添加团队成员
//
//	@Summary		添加团队成员
//	@Description	添加用户到团队（用户必须先是组织成员）
//	@Tags			Organization - Team Member Management
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			org_id	path		int													true	"组织ID"	minimum(1)
//	@Param			team_id	path		int													true	"团队ID"	minimum(1)
//	@Param			request	body		organization.AddTeamMemberDTO						true	"成员信息"
//	@Success		201		{object}	response.DataResponse[organization.TeamMemberDTO]	"成员添加成功"
//	@Failure		400		{object}	response.ErrorResponse								"参数错误、成员已存在或用户非组织成员"
//	@Failure		401		{object}	response.ErrorResponse								"未授权"
//	@Failure		403		{object}	response.ErrorResponse								"权限不足"
//	@Failure		500		{object}	response.ErrorResponse								"服务器内部错误"
//	@Router			/api/org/{org_id}/teams/{team_id}/members [post]
//
//nolint:dupl // Similar structure to TeamHandler.Update is intentional
func (h *TeamMemberHandler) Add(c *gin.Context) {
	orgID, err := strconv.ParseUint(c.Param("org_id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid organization ID")
		return
	}

	teamID, err := strconv.ParseUint(c.Param("team_id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid team ID")
		return
	}

	var req organization.AddTeamMemberDTO
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		response.ValidationError(c, bindErr.Error())
		return
	}

	result, err := h.addMemberHandler.Handle(c.Request.Context(), organization.AddTeamMemberCommand{
		OrgID:  uint(orgID),
		TeamID: uint(teamID),
		UserID: req.UserID,
		Role:   req.Role,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, "team member added successfully", result)
}

// Remove 移除团队成员
//
//	@Summary		移除团队成员
//	@Description	从团队中移除成员
//	@Tags			Organization - Team Member Management
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			org_id	path		int							true	"组织ID"	minimum(1)
//	@Param			team_id	path		int							true	"团队ID"	minimum(1)
//	@Param			user_id	path		int							true	"用户ID"	minimum(1)
//	@Success		200		{object}	response.MessageResponse	"成员移除成功"
//	@Failure		400		{object}	response.ErrorResponse		"无效的ID"
//	@Failure		401		{object}	response.ErrorResponse		"未授权"
//	@Failure		403		{object}	response.ErrorResponse		"权限不足"
//	@Failure		404		{object}	response.ErrorResponse		"成员不存在"
//	@Failure		500		{object}	response.ErrorResponse		"服务器内部错误"
//	@Router			/api/org/{org_id}/teams/{team_id}/members/{user_id} [delete]
func (h *TeamMemberHandler) Remove(c *gin.Context) {
	orgID, err := strconv.ParseUint(c.Param("org_id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid organization ID")
		return
	}

	teamID, err := strconv.ParseUint(c.Param("team_id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid team ID")
		return
	}

	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid user ID")
		return
	}

	if err = h.removeMemberHandler.Handle(c.Request.Context(), organization.RemoveTeamMemberCommand{
		OrgID:  uint(orgID),
		TeamID: uint(teamID),
		UserID: uint(userID),
	}); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, "team member removed successfully", nil)
}
