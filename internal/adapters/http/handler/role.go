package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/response"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/role"
)

// RoleHandler handles role management operations (DDD+CQRS Use Case Pattern)
type RoleHandler struct {
	// Command Handlers
	createRoleHandler     *role.CreateRoleHandler
	updateRoleHandler     *role.UpdateRoleHandler
	deleteRoleHandler     *role.DeleteRoleHandler
	setPermissionsHandler *role.SetPermissionsHandler

	// Query Handlers
	getRoleHandler         *role.GetRoleHandler
	listRolesHandler       *role.ListRolesHandler
	listPermissionsHandler *role.ListPermissionsHandler
}

// NewRoleHandler creates a new RoleHandler instance
func NewRoleHandler(
	createRoleHandler *role.CreateRoleHandler,
	updateRoleHandler *role.UpdateRoleHandler,
	deleteRoleHandler *role.DeleteRoleHandler,
	setPermissionsHandler *role.SetPermissionsHandler,
	getRoleHandler *role.GetRoleHandler,
	listRolesHandler *role.ListRolesHandler,
	listPermissionsHandler *role.ListPermissionsHandler,
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

// CreateRole creates a new role
//
// @Summary      创建角色
// @Description  管理员创建新的系统角色
// @Tags         管理员 - 角色管理 (Admin - Role Management)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body role.CreateRoleDTO true "角色信息"
// @Success      201 {object} response.DataResponse[role.CreateRoleResultDTO] "角色创建成功"
// @Failure      400 {object} response.ErrorResponse "参数错误或角色名已存在"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Failure      403 {object} response.ErrorResponse "权限不足"
// @Failure      500 {object} response.ErrorResponse "服务器内部错误"
// @Router       /api/admin/roles [post]
// @x-permission {"scope":"admin:roles:create"}
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req role.CreateRoleDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	// 调用 Use Case Handler
	result, err := h.createRoleHandler.Handle(c.Request.Context(), role.CreateRoleCommand(req))

	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	// 转换为 DTO 响应
	resp := role.CreateRoleResultDTO{
		RoleID:      result.RoleID,
		Name:        result.Name,
		DisplayName: result.DisplayName,
	}
	response.Created(c, "role created successfully", resp)
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
// @Success      200 {object} response.PagedResponse[role.RoleDTO] "角色列表"
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
	result, err := h.listRolesHandler.Handle(c.Request.Context(), role.ListRolesQuery{
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
// @Success      200 {object} response.DataResponse[role.RoleDTO] "角色详情"
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
	result, err := h.getRoleHandler.Handle(c.Request.Context(), role.GetRoleQuery{
		RoleID: uint(id),
	})

	if err != nil {
		response.NotFound(c, "role")
		return
	}

	response.OK(c, "success", result)
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
// @Param        request body role.UpdateRoleDTO true "更新信息"
// @Success      200 {object} response.DataResponse[role.RoleDTO] "角色更新成功"
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

	var req role.UpdateRoleDTO

	if err = c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	// 调用 Use Case Handler
	result, err := h.updateRoleHandler.Handle(c.Request.Context(), role.UpdateRoleCommand{
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
// @Success      200 {object} response.MessageResponse "角色删除成功"
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
	err = h.deleteRoleHandler.Handle(c.Request.Context(), role.DeleteRoleCommand{
		RoleID: uint(id),
	})

	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, "role deleted successfully", nil)
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
// @Param        request body role.SetPermissionsDTO true "权限ID列表"
// @Success      200 {object} response.MessageResponse "权限设置成功"
// @Failure      400 {object} response.ErrorResponse "无效的角色ID或参数错误"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Failure      403 {object} response.ErrorResponse "权限不足"
// @Failure      404 {object} response.ErrorResponse "角色不存在"
// @Failure      500 {object} response.ErrorResponse "服务器内部错误"
// @Router       /api/admin/roles/{id}/permissions [put]
// @x-permission {"scope":"admin:roles:update"}
func (h *RoleHandler) SetPermissions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid role ID")
		return
	}

	var req role.SetPermissionsDTO

	if err = c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	// 调用 Use Case Handler
	err = h.setPermissionsHandler.Handle(c.Request.Context(), role.SetPermissionsCommand{
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
// @Success      200 {object} response.PagedResponse[role.PermissionDTO] "权限列表"
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
	result, err := h.listPermissionsHandler.Handle(c.Request.Context(), role.ListPermissionsQuery{
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
