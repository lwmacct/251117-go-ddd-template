package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/response"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/org"
)

// ListOrgsQuery 组织列表查询参数
type ListOrgsQuery struct {
	response.PaginationQueryDTO

	Status string `form:"status" binding:"omitempty,oneof=active inactive"`
}

// ToQuery 转换为 Application 层 Query 对象
func (q *ListOrgsQuery) ToQuery() org.ListOrgsQuery {
	return org.ListOrgsQuery{
		Offset: q.GetOffset(),
		Limit:  q.GetLimit(),
	}
}

// OrgHandler 组织管理 Handler（系统管理域）
type OrgHandler struct {
	// Command Handlers
	createHandler *org.CreateHandler
	updateHandler *org.UpdateHandler
	deleteHandler *org.DeleteHandler

	// Query Handlers
	getHandler  *org.GetHandler
	listHandler *org.ListHandler
}

// NewOrgHandler 创建组织管理 Handler
func NewOrgHandler(
	createHandler *org.CreateHandler,
	updateHandler *org.UpdateHandler,
	deleteHandler *org.DeleteHandler,
	getHandler *org.GetHandler,
	listHandler *org.ListHandler,
) *OrgHandler {
	return &OrgHandler{
		createHandler: createHandler,
		updateHandler: updateHandler,
		deleteHandler: deleteHandler,
		getHandler:    getHandler,
		listHandler:   listHandler,
	}
}

// Create 创建组织
//
//	@Summary		创建组织
//	@Description	系统管理员创建新组织
//	@Tags			Admin - Organization Management
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		org.CreateOrgDTO					true	"组织信息"
//	@Success		201		{object}	response.DataResponse[org.OrgDTO]	"组织创建成功"
//	@Failure		400		{object}	response.ErrorResponse						"参数错误或组织名已存在"
//	@Failure		401		{object}	response.ErrorResponse						"未授权"
//	@Failure		403		{object}	response.ErrorResponse						"权限不足"
//	@Failure		500		{object}	response.ErrorResponse						"服务器内部错误"
//	@Router			/api/system/orgs [post]
func (h *OrgHandler) Create(c *gin.Context) {
	var req org.CreateOrgDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result, err := h.createHandler.Handle(c.Request.Context(), org.CreateOrgCommand{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		Avatar:      req.Avatar,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, "organization created successfully", result)
}

// List 获取组织列表
//
//	@Summary		组织列表
//	@Description	分页获取所有组织
//	@Tags			Admin - Organization Management
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			params	query		handler.ListOrgsQuery							false	"查询参数"
//	@Success		200		{object}	response.PagedResponse[org.OrgDTO]		"组织列表"
//	@Failure		401		{object}	response.ErrorResponse							"未授权"
//	@Failure		403		{object}	response.ErrorResponse							"权限不足"
//	@Failure		500		{object}	response.ErrorResponse							"服务器内部错误"
//	@Router			/api/system/orgs [get]
func (h *OrgHandler) List(c *gin.Context) {
	var q ListOrgsQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result, err := h.listHandler.Handle(c.Request.Context(), q.ToQuery())
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	meta := response.NewPaginationMeta(int(result.Total), q.GetPage(), q.GetLimit())
	response.List(c, "success", result.Items, meta)
}

// Get 获取组织详情
//
//	@Summary		组织详情
//	@Description	根据 ID 获取组织详情
//	@Tags			Admin - Organization Management
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int											true	"组织ID"	minimum(1)
//	@Success		200	{object}	response.DataResponse[org.OrgDTO]	"组织详情"
//	@Failure		400	{object}	response.ErrorResponse						"无效的组织ID"
//	@Failure		401	{object}	response.ErrorResponse						"未授权"
//	@Failure		403	{object}	response.ErrorResponse						"权限不足"
//	@Failure		404	{object}	response.ErrorResponse						"组织不存在"
//	@Router			/api/system/orgs/{id} [get]
func (h *OrgHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid organization ID")
		return
	}

	result, err := h.getHandler.Handle(c.Request.Context(), org.GetOrgQuery{
		OrgID: uint(id),
	})
	if err != nil {
		response.NotFound(c, "organization")
		return
	}

	response.OK(c, "success", result)
}

// Update 更新组织
//
//	@Summary		更新组织
//	@Description	更新组织信息
//	@Tags			Admin - Organization Management
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int											true	"组织ID"	minimum(1)
//	@Param			request	body		org.UpdateOrgDTO					true	"更新信息"
//	@Success		200		{object}	response.DataResponse[org.OrgDTO]	"组织更新成功"
//	@Failure		400		{object}	response.ErrorResponse						"参数错误"
//	@Failure		401		{object}	response.ErrorResponse						"未授权"
//	@Failure		403		{object}	response.ErrorResponse						"权限不足"
//	@Failure		404		{object}	response.ErrorResponse						"组织不存在"
//	@Failure		500		{object}	response.ErrorResponse						"服务器内部错误"
//	@Router			/api/system/orgs/{id} [put]
func (h *OrgHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid organization ID")
		return
	}

	var req org.UpdateOrgDTO
	if err = c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result, err := h.updateHandler.Handle(c.Request.Context(), org.UpdateOrgCommand{
		OrgID:       uint(id),
		DisplayName: req.DisplayName,
		Description: req.Description,
		Avatar:      req.Avatar,
		Status:      req.Status,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, "organization updated successfully", result)
}

// Delete 删除组织
//
//	@Summary		删除组织
//	@Description	软删除组织
//	@Tags			Admin - Organization Management
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int							true	"组织ID"	minimum(1)
//	@Success		200	{object}	response.MessageResponse	"组织删除成功"
//	@Failure		400	{object}	response.ErrorResponse		"无效的组织ID"
//	@Failure		401	{object}	response.ErrorResponse		"未授权"
//	@Failure		403	{object}	response.ErrorResponse		"权限不足"
//	@Failure		404	{object}	response.ErrorResponse		"组织不存在"
//	@Failure		500	{object}	response.ErrorResponse		"服务器内部错误"
//	@Router			/api/system/orgs/{id} [delete]
func (h *OrgHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid organization ID")
		return
	}

	if err = h.deleteHandler.Handle(c.Request.Context(), org.DeleteOrgCommand{
		OrgID: uint(id),
	}); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, "organization deleted successfully", nil)
}
