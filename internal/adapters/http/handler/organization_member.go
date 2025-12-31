package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/response"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/organization"
)

// ListMembersQuery 成员列表查询参数
type ListMembersQuery struct {
	response.PaginationQueryDTO
}

// ToQuery 转换为 Application 层 Query 对象
func (q *ListMembersQuery) ToQuery(orgID uint) organization.ListMembersQuery {
	return organization.ListMembersQuery{
		OrgID:  orgID,
		Offset: q.GetOffset(),
		Limit:  q.GetLimit(),
	}
}

// OrgMemberHandler 组织成员管理 Handler
type OrgMemberHandler struct {
	// Command Handlers
	addMemberHandler        *organization.MemberAddHandler
	removeMemberHandler     *organization.MemberRemoveHandler
	updateMemberRoleHandler *organization.MemberUpdateRoleHandler

	// Query Handlers
	listMembersHandler *organization.MemberListHandler
}

// NewOrgMemberHandler 创建组织成员管理 Handler
func NewOrgMemberHandler(
	addMemberHandler *organization.MemberAddHandler,
	removeMemberHandler *organization.MemberRemoveHandler,
	updateMemberRoleHandler *organization.MemberUpdateRoleHandler,
	listMembersHandler *organization.MemberListHandler,
) *OrgMemberHandler {
	return &OrgMemberHandler{
		addMemberHandler:        addMemberHandler,
		removeMemberHandler:     removeMemberHandler,
		updateMemberRoleHandler: updateMemberRoleHandler,
		listMembersHandler:      listMembersHandler,
	}
}

// List 获取组织成员列表
//
//	@Summary		成员列表
//	@Description	分页获取组织成员列表
//	@Tags			Organization - Member Management
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			org_id	path		int												true	"组织ID"	minimum(1)
//	@Param			params	query		handler.ListMembersQuery						false	"查询参数"
//	@Success		200		{object}	response.PagedResponse[organization.MemberDTO]	"成员列表"
//	@Failure		401		{object}	response.ErrorResponse							"未授权"
//	@Failure		403		{object}	response.ErrorResponse							"权限不足"
//	@Failure		500		{object}	response.ErrorResponse							"服务器内部错误"
//	@Router			/api/org/{org_id}/members [get]
func (h *OrgMemberHandler) List(c *gin.Context) {
	orgID, err := strconv.ParseUint(c.Param("org_id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid organization ID")
		return
	}

	var q ListMembersQuery
	if bindErr := c.ShouldBindQuery(&q); bindErr != nil {
		response.ValidationError(c, bindErr.Error())
		return
	}

	result, err := h.listMembersHandler.Handle(c.Request.Context(), q.ToQuery(uint(orgID)))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	meta := response.NewPaginationMeta(int(result.Total), q.GetPage(), q.GetLimit())
	response.List(c, "success", result.Items, meta)
}

// Add 添加组织成员
//
//	@Summary		添加成员
//	@Description	添加用户到组织
//	@Tags			Organization - Member Management
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			org_id	path		int												true	"组织ID"	minimum(1)
//	@Param			request	body		organization.AddMemberDTO						true	"成员信息"
//	@Success		201		{object}	response.DataResponse[organization.MemberDTO]	"成员添加成功"
//	@Failure		400		{object}	response.ErrorResponse							"参数错误或成员已存在"
//	@Failure		401		{object}	response.ErrorResponse							"未授权"
//	@Failure		403		{object}	response.ErrorResponse							"权限不足"
//	@Failure		500		{object}	response.ErrorResponse							"服务器内部错误"
//	@Router			/api/org/{org_id}/members [post]
func (h *OrgMemberHandler) Add(c *gin.Context) {
	orgID, err := strconv.ParseUint(c.Param("org_id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid organization ID")
		return
	}

	var req organization.AddMemberDTO
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		response.ValidationError(c, bindErr.Error())
		return
	}

	result, err := h.addMemberHandler.Handle(c.Request.Context(), organization.AddMemberCommand{
		OrgID:  uint(orgID),
		UserID: req.UserID,
		Role:   req.Role,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, "member added successfully", result)
}

// Remove 移除组织成员
//
//	@Summary		移除成员
//	@Description	从组织中移除成员
//	@Tags			Organization - Member Management
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			org_id	path		int							true	"组织ID"	minimum(1)
//	@Param			user_id	path		int							true	"用户ID"	minimum(1)
//	@Success		200		{object}	response.MessageResponse	"成员移除成功"
//	@Failure		400		{object}	response.ErrorResponse		"无效的ID或无法移除最后的所有者"
//	@Failure		401		{object}	response.ErrorResponse		"未授权"
//	@Failure		403		{object}	response.ErrorResponse		"权限不足"
//	@Failure		404		{object}	response.ErrorResponse		"成员不存在"
//	@Failure		500		{object}	response.ErrorResponse		"服务器内部错误"
//	@Router			/api/org/{org_id}/members/{user_id} [delete]
func (h *OrgMemberHandler) Remove(c *gin.Context) {
	orgID, err := strconv.ParseUint(c.Param("org_id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid organization ID")
		return
	}

	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid user ID")
		return
	}

	if err = h.removeMemberHandler.Handle(c.Request.Context(), organization.RemoveMemberCommand{
		OrgID:  uint(orgID),
		UserID: uint(userID),
	}); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, "member removed successfully", nil)
}

// UpdateRole 更新成员角色
//
//	@Summary		更新成员角色
//	@Description	更新组织成员的角色
//	@Tags			Organization - Member Management
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			org_id	path		int												true	"组织ID"	minimum(1)
//	@Param			user_id	path		int												true	"用户ID"	minimum(1)
//	@Param			request	body		organization.UpdateMemberRoleDTO				true	"角色信息"
//	@Success		200		{object}	response.DataResponse[organization.MemberDTO]	"角色更新成功"
//	@Failure		400		{object}	response.ErrorResponse							"参数错误或无法降级最后的所有者"
//	@Failure		401		{object}	response.ErrorResponse							"未授权"
//	@Failure		403		{object}	response.ErrorResponse							"权限不足"
//	@Failure		404		{object}	response.ErrorResponse							"成员不存在"
//	@Failure		500		{object}	response.ErrorResponse							"服务器内部错误"
//	@Router			/api/org/{org_id}/members/{user_id}/role [put]
func (h *OrgMemberHandler) UpdateRole(c *gin.Context) {
	orgID, err := strconv.ParseUint(c.Param("org_id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid organization ID")
		return
	}

	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid user ID")
		return
	}

	var req organization.UpdateMemberRoleDTO
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		response.ValidationError(c, bindErr.Error())
		return
	}

	if err = h.updateMemberRoleHandler.Handle(c.Request.Context(), organization.UpdateMemberRoleCommand{
		OrgID:  uint(orgID),
		UserID: uint(userID),
		Role:   req.Role,
	}); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, "member role updated successfully", nil)
}
