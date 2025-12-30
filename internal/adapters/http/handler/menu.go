package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/response"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/menu"
)

// MenuHandler handles menu management operations (DDD+CQRS Use Case Pattern)
type MenuHandler struct {
	// Command Handlers
	createHandler  *menu.CreateHandler
	updateHandler  *menu.UpdateHandler
	deleteHandler  *menu.DeleteHandler
	reorderHandler *menu.ReorderHandler

	// Query Handlers
	getHandler  *menu.GetHandler
	listHandler *menu.ListHandler
}

// NewMenuHandler creates a new MenuHandler instance
func NewMenuHandler(
	createHandler *menu.CreateHandler,
	updateHandler *menu.UpdateHandler,
	deleteHandler *menu.DeleteHandler,
	reorderHandler *menu.ReorderHandler,
	getHandler *menu.GetHandler,
	listHandler *menu.ListHandler,
) *MenuHandler {
	return &MenuHandler{
		createHandler:  createHandler,
		updateHandler:  updateHandler,
		deleteHandler:  deleteHandler,
		reorderHandler: reorderHandler,
		getHandler:     getHandler,
		listHandler:    listHandler,
	}
}

// Create 创建菜单
//
// @Summary      创建菜单
// @Description  管理员创建新的系统菜单项
// @Tags         Admin - Menu Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body menu.CreateDTO true "菜单信息"
// @Success      201 {object} response.DataResponse[menu.MenuDTO] "菜单创建成功"
// @Failure      400 {object} response.ErrorResponse "参数错误"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Failure      403 {object} response.ErrorResponse "权限不足"
// @Failure      500 {object} response.ErrorResponse "服务器内部错误"
// @Router       /api/system/menus [post]
func (h *MenuHandler) Create(c *gin.Context) {
	var req menu.CreateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	visible := true
	if req.Visible != nil {
		visible = *req.Visible
	}

	// 调用 Use Case Handler
	result, err := h.createHandler.Handle(c.Request.Context(), menu.CreateCommand{
		Title:    req.Title,
		Path:     req.Path,
		Icon:     req.Icon,
		ParentID: req.ParentID,
		Order:    req.Order,
		Visible:  visible,
	})

	if err != nil {
		response.InternalError(c, "Failed to create menu")
		return
	}

	response.Created(c, "menu created successfully", result)
}

// List 获取菜单列表（树形结构）
//
// @Summary      菜单列表
// @Description  获取所有菜单的树形结构（包含父子关系）
// @Tags         Admin - Menu Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} response.DataResponse[[]menu.MenuDTO] "菜单树"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Failure      403 {object} response.ErrorResponse "权限不足"
// @Failure      500 {object} response.ErrorResponse "服务器内部错误"
// @Router       /api/system/menus [get]
func (h *MenuHandler) List(c *gin.Context) {
	// 调用 Use Case Handler
	menus, err := h.listHandler.Handle(c.Request.Context(), menu.ListQuery{})

	if err != nil {
		response.InternalError(c, "Failed to fetch menus")
		return
	}

	response.OK(c, "success", menus)
}

// Get 获取菜单详情
//
// @Summary      菜单详情
// @Description  根据菜单ID获取菜单详细信息
// @Tags         Admin - Menu Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "菜单ID" minimum(1)
// @Success      200 {object} response.DataResponse[menu.MenuDTO] "菜单详情"
// @Failure      400 {object} response.ErrorResponse "无效的菜单ID"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Failure      403 {object} response.ErrorResponse "权限不足"
// @Failure      404 {object} response.ErrorResponse "菜单不存在"
// @Router       /api/system/menus/{id} [get]
func (h *MenuHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid menu ID")
		return
	}

	// 调用 Use Case Handler
	menu, err := h.getHandler.Handle(c.Request.Context(), menu.GetQuery{
		MenuID: uint(id),
	})

	if err != nil {
		response.NotFound(c, "Menu")
		return
	}

	response.OK(c, "success", menu)
}

// Update 更新菜单
//
// @Summary      更新菜单
// @Description  管理员更新菜单的标题、路径、图标等信息
// @Tags         Admin - Menu Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "菜单ID" minimum(1)
// @Param        request body menu.UpdateDTO true "更新信息"
// @Success      200 {object} response.DataResponse[menu.MenuDTO] "菜单更新成功"
// @Failure      400 {object} response.ErrorResponse "无效的菜单ID或参数错误"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Failure      403 {object} response.ErrorResponse "权限不足"
// @Failure      404 {object} response.ErrorResponse "菜单不存在"
// @Failure      500 {object} response.ErrorResponse "服务器内部错误"
// @Router       /api/system/menus/{id} [put]
func (h *MenuHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid menu ID")
		return
	}

	var req menu.UpdateDTO
	if err = c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 调用 Use Case Handler
	menu, err := h.updateHandler.Handle(c.Request.Context(), menu.UpdateCommand{
		MenuID:   uint(id),
		Title:    req.Title,
		Path:     req.Path,
		Icon:     req.Icon,
		ParentID: req.ParentID,
		Order:    req.Order,
		Visible:  req.Visible,
	})

	if err != nil {
		response.InternalError(c, "Failed to update menu")
		return
	}

	response.OK(c, "menu updated successfully", menu)
}

// Delete 删除菜单
//
// @Summary      删除菜单
// @Description  管理员删除指定菜单（如果有子菜单，可能会失败）
// @Tags         Admin - Menu Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "菜单ID" minimum(1)
// @Success      204 "菜单删除成功"
// @Failure      400 {object} response.ErrorResponse "无效的菜单ID"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Failure      403 {object} response.ErrorResponse "权限不足"
// @Failure      404 {object} response.ErrorResponse "菜单不存在"
// @Failure      500 {object} response.ErrorResponse "服务器内部错误或菜单有子项"
// @Router       /api/system/menus/{id} [delete]
func (h *MenuHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid menu ID")
		return
	}

	// 调用 Use Case Handler
	err = h.deleteHandler.Handle(c.Request.Context(), menu.DeleteCommand{
		MenuID: uint(id),
	})

	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.NoContent(c)
}

// Reorder 批量更新菜单排序
//
// @Summary      重排菜单
// @Description  管理员批量更新菜单的排序和父级关系
// @Tags         Admin - Menu Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body menu.ReorderDTO true "菜单排序信息"
// @Success      204 "菜单排序更新成功"
// @Failure      400 {object} response.ErrorResponse "参数错误"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Failure      403 {object} response.ErrorResponse "权限不足"
// @Failure      500 {object} response.ErrorResponse "服务器内部错误"
// @Router       /api/system/menus/reorder [post]
func (h *MenuHandler) Reorder(c *gin.Context) {
	var req menu.ReorderDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 转换为 Command
	menus := make([]menu.MenuItemCommand, len(req.Menus))
	for i, m := range req.Menus {
		menus[i].ID = m.ID
		menus[i].Order = m.Order
		menus[i].ParentID = m.ParentID
	}

	// 调用 Use Case Handler
	err := h.reorderHandler.Handle(c.Request.Context(), menu.ReorderCommand{
		Menus: menus,
	})

	if err != nil {
		response.InternalError(c, "Failed to update menu order")
		return
	}

	response.NoContent(c)
}
