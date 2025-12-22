package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/response"
	roleCommand "github.com/lwmacct/251117-go-ddd-template/internal/application/role/command"
	roleQuery "github.com/lwmacct/251117-go-ddd-template/internal/application/role/query"
)

// RoleHandler handles role management operations (DDD+CQRS Use Case Pattern)
type RoleHandler struct {
	// Command Handlers
	createRoleHandler     *roleCommand.CreateRoleHandler
	updateRoleHandler     *roleCommand.UpdateRoleHandler
	deleteRoleHandler     *roleCommand.DeleteRoleHandler
	setPermissionsHandler *roleCommand.SetPermissionsHandler

	// Query Handlers
	getRoleHandler         *roleQuery.GetRoleHandler
	listRolesHandler       *roleQuery.ListRolesHandler
	listPermissionsHandler *roleQuery.ListPermissionsHandler
}

// NewRoleHandler creates a new RoleHandler instance
func NewRoleHandler(
	createRoleHandler *roleCommand.CreateRoleHandler,
	updateRoleHandler *roleCommand.UpdateRoleHandler,
	deleteRoleHandler *roleCommand.DeleteRoleHandler,
	setPermissionsHandler *roleCommand.SetPermissionsHandler,
	getRoleHandler *roleQuery.GetRoleHandler,
	listRolesHandler *roleQuery.ListRolesHandler,
	listPermissionsHandler *roleQuery.ListPermissionsHandler,
) *RoleHandler {
	return &RoleHandler{
		createRoleHandler:      createRoleHandler,
		updateRoleHandler:      updateRoleHandler,
		deleteRoleHandler:      deleteRoleHandler,
		setPermissionsHandler:  setPermissionsHandler,
		getRoleHandler:         getRoleHandler,
		listRolesHandler:       listRolesHandler,
		listPermissionsHandler: listPermissionsHandler,
	}
}

// CreateRoleRequest 创建角色请求
type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=50" example:"developer"`
	DisplayName string `json:"display_name" binding:"required,max=100" example:"开发者"`
	Description string `json:"description" binding:"max=255" example:"系统开发人员角色"`
}

// CreateRole creates a new role
//
// @Summary      创建角色
// @Description  管理员创建新的系统角色
// @Tags         管理员 - 角色管理 (Admin - Role Management)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body CreateRoleRequest true "角色信息"
// @Success      201 {object} response.Response{data=object{role_id=uint,name=string,display_name=string}} "角色创建成功"
// @Failure      400 {object} response.ErrorResponse "参数错误或角色名已存在"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Failure      403 {object} response.ErrorResponse "权限不足"
// @Failure      500 {object} response.ErrorResponse "服务器内部错误"
// @Router       /api/admin/roles [post]
// @x-permission {"scope":"admin:roles:create"}
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req CreateRoleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	// 调用 Use Case Handler
	result, err := h.createRoleHandler.Handle(c.Request.Context(), roleCommand.CreateRoleCommand{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
	})

	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, "role created successfully", result)
}

// ListRoles lists all roles
//
// @Summary      获取角色列表
// @Description  分页获取所有系统角色
// @Tags         管理员 - 角色管理 (Admin - Role Management)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page query int false "页码" default(1) minimum(1)
// @Param        limit query int false "每页数量" default(20) minimum(1) maximum(100)
// @Success      200 {object} response.ListResponse{data=[]object{id=uint,name=string,display_name=string,description=string},meta=response.PaginationMeta} "角色列表"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Failure      403 {object} response.ErrorResponse "权限不足"
// @Failure      500 {object} response.ErrorResponse "服务器内部错误"
// @Router       /api/admin/roles [get]
// @x-permission {"scope":"admin:roles:read"}
func (h *RoleHandler) ListRoles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	// 调用 Use Case Handler
	result, err := h.listRolesHandler.Handle(c.Request.Context(), roleQuery.ListRolesQuery{
		Page:  page,
		Limit: limit,
	})

	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	meta := response.NewPaginationMeta(int(result.Total), page, limit)
	response.List(c, "success", result.Roles, meta)
}

// GetRole gets a role by ID
//
// @Summary      获取角色详情
// @Description  根据角色ID获取角色详细信息（包含权限列表）
// @Tags         管理员 - 角色管理 (Admin - Role Management)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "角色ID" minimum(1)
// @Success      200 {object} response.Response{data=object{id=uint,name=string,display_name=string,description=string,permissions=[]object}} "角色详情"
// @Failure      400 {object} response.ErrorResponse "无效的角色ID"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Failure      403 {object} response.ErrorResponse "权限不足"
// @Failure      404 {object} response.ErrorResponse "角色不存在"
// @Router       /api/admin/roles/{id} [get]
// @x-permission {"scope":"admin:roles:read"}
func (h *RoleHandler) GetRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid role ID")
		return
	}

	// 调用 Use Case Handler
	result, err := h.getRoleHandler.Handle(c.Request.Context(), roleQuery.GetRoleQuery{
		RoleID: uint(id),
	})

	if err != nil {
		response.NotFound(c, "role")
		return
	}

	response.OK(c, "success", result)
}

// UpdateRoleRequest 更新角色请求
type UpdateRoleRequest struct {
	DisplayName *string `json:"display_name" binding:"omitempty,max=100" example:"高级开发者"`
	Description *string `json:"description" binding:"omitempty,max=255" example:"具有更多权限的开发人员"`
}

// UpdateRole updates a role
//
// @Summary      更新角色信息
// @Description  管理员更新角色的显示名称和描述
// @Tags         管理员 - 角色管理 (Admin - Role Management)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "角色ID" minimum(1)
// @Param        request body UpdateRoleRequest true "更新信息"
// @Success      200 {object} response.Response{data=object{id=uint,name=string,display_name=string,description=string}} "角色更新成功"
// @Failure      400 {object} response.ErrorResponse "无效的角色ID或参数错误"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Failure      403 {object} response.ErrorResponse "权限不足"
// @Failure      404 {object} response.ErrorResponse "角色不存在"
// @Failure      500 {object} response.ErrorResponse "服务器内部错误"
// @Router       /api/admin/roles/{id} [put]
// @x-permission {"scope":"admin:roles:update"}
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid role ID")
		return
	}

	var req UpdateRoleRequest

	if err = c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	// 调用 Use Case Handler
	result, err := h.updateRoleHandler.Handle(c.Request.Context(), roleCommand.UpdateRoleCommand{
		RoleID:      uint(id),
		DisplayName: req.DisplayName,
		Description: req.Description,
	})

	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, "role updated successfully", result)
}

// DeleteRole deletes a role
//
// @Summary      删除角色
// @Description  管理员删除指定角色（如果角色被用户使用，可能会失败）
// @Tags         管理员 - 角色管理 (Admin - Role Management)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "角色ID" minimum(1)
// @Success      200 {object} response.Response{data=object{message=string}} "角色删除成功"
// @Failure      400 {object} response.ErrorResponse "无效的角色ID"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Failure      403 {object} response.ErrorResponse "权限不足"
// @Failure      404 {object} response.ErrorResponse "角色不存在"
// @Failure      500 {object} response.ErrorResponse "服务器内部错误或角色被使用中"
// @Router       /api/admin/roles/{id} [delete]
// @x-permission {"scope":"admin:roles:delete"}
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid role ID")
		return
	}

	// 调用 Use Case Handler
	err = h.deleteRoleHandler.Handle(c.Request.Context(), roleCommand.DeleteRoleCommand{
		RoleID: uint(id),
	})

	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, "role deleted successfully", nil)
}

// SetPermissionsRequest 设置权限请求
type SetPermissionsRequest struct {
	PermissionIDs []uint `json:"permission_ids" binding:"required" example:"1,2,3"`
}

// SetPermissions sets permissions for a role
//
// @Summary      设置角色权限
// @Description  管理员为指定角色设置权限（会覆盖现有权限）
// @Tags         管理员 - 角色管理 (Admin - Role Management)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "角色ID" minimum(1)
// @Param        request body SetPermissionsRequest true "权限ID列表"
// @Success      200 {object} response.Response{data=object{message=string}} "权限设置成功"
// @Failure      400 {object} response.ErrorResponse "无效的角色ID或参数错误"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Failure      403 {object} response.ErrorResponse "权限不足"
// @Failure      404 {object} response.ErrorResponse "角色不存在"
// @Failure      500 {object} response.ErrorResponse "服务器内部错误"
// @Router       /api/admin/roles/{id}/permissions [post]
// @x-permission {"scope":"admin:roles:update"}
func (h *RoleHandler) SetPermissions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid role ID")
		return
	}

	var req SetPermissionsRequest

	if err = c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	// 调用 Use Case Handler
	err = h.setPermissionsHandler.Handle(c.Request.Context(), roleCommand.SetPermissionsCommand{
		RoleID:        uint(id),
		PermissionIDs: req.PermissionIDs,
	})

	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, "permissions set successfully", nil)
}

// ListPermissions lists all permissions
//
// @Summary      获取权限列表
// @Description  分页获取所有系统权限
// @Tags         管理员 - 角色管理 (Admin - Role Management)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page query int false "页码" default(1) minimum(1)
// @Param        limit query int false "每页数量" default(50) minimum(1) maximum(100)
// @Success      200 {object} response.ListResponse{data=[]object{id=uint,name=string,display_name=string,resource=string,action=string},meta=response.PaginationMeta} "权限列表"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Failure      403 {object} response.ErrorResponse "权限不足"
// @Failure      500 {object} response.ErrorResponse "服务器内部错误"
// @Router       /api/admin/permissions [get]
// @x-permission {"scope":"admin:permissions:read"}
func (h *RoleHandler) ListPermissions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 50
	}

	// 调用 Use Case Handler
	result, err := h.listPermissionsHandler.Handle(c.Request.Context(), roleQuery.ListPermissionsQuery{
		Page:  page,
		Limit: limit,
	})

	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	meta := response.NewPaginationMeta(int(result.Total), page, limit)
	response.List(c, "success", result.Permissions, meta)
}
