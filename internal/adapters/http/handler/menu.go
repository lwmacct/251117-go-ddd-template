package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/response"
	menuCommand "github.com/lwmacct/251117-go-ddd-template/internal/application/menu/command"
	menuQuery "github.com/lwmacct/251117-go-ddd-template/internal/application/menu/query"
)

// MenuHandler handles menu management operations (DDD+CQRS Use Case Pattern)
type MenuHandler struct {
	// Command Handlers
	createMenuHandler   *menuCommand.CreateMenuHandler
	updateMenuHandler   *menuCommand.UpdateMenuHandler
	deleteMenuHandler   *menuCommand.DeleteMenuHandler
	reorderMenusHandler *menuCommand.ReorderMenusHandler

	// Query Handlers
	getMenuHandler   *menuQuery.GetMenuHandler
	listMenusHandler *menuQuery.ListMenusHandler
}

// NewMenuHandler creates a new MenuHandler instance
func NewMenuHandler(
	createMenuHandler *menuCommand.CreateMenuHandler,
	updateMenuHandler *menuCommand.UpdateMenuHandler,
	deleteMenuHandler *menuCommand.DeleteMenuHandler,
	reorderMenusHandler *menuCommand.ReorderMenusHandler,
	getMenuHandler *menuQuery.GetMenuHandler,
	listMenusHandler *menuQuery.ListMenusHandler,
) *MenuHandler {
	return &MenuHandler{
		createMenuHandler:   createMenuHandler,
		updateMenuHandler:   updateMenuHandler,
		deleteMenuHandler:   deleteMenuHandler,
		reorderMenusHandler: reorderMenusHandler,
		getMenuHandler:      getMenuHandler,
		listMenusHandler:    listMenusHandler,
	}
}

// CreateMenuRequest 创建菜单请求
type CreateMenuRequest struct {
	Title    string `json:"title" binding:"required,min=1,max=100" example:"系统管理"`
	Path     string `json:"path" binding:"required,max=255" example:"/system"`
	Icon     string `json:"icon" binding:"omitempty,max=100" example:"setting"`
	ParentID *uint  `json:"parent_id" example:"0"`
	Order    int    `json:"order" example:"1"`
	Visible  *bool  `json:"visible" example:"true"`
}

// UpdateMenuRequest 更新菜单请求
type UpdateMenuRequest struct {
	Title    *string `json:"title" binding:"omitempty,min=1,max=100"`
	Path     *string `json:"path" binding:"omitempty,max=255"`
	Icon     *string `json:"icon" binding:"omitempty,max=100"`
	ParentID *uint   `json:"parent_id"`
	Order    *int    `json:"order"`
	Visible  *bool   `json:"visible"`
}

// ReorderMenusRequest 批量更新排序请求
type ReorderMenusRequest struct {
	Menus []struct {
		ID       uint  `json:"id" binding:"required"`
		Order    int   `json:"order"`
		ParentID *uint `json:"parent_id"`
	} `json:"menus" binding:"required,dive"`
}

// Create 创建菜单
//
// @Summary      创建菜单
// @Description  管理员创建新的系统菜单项
// @Tags         管理员 - 菜单管理 (Admin - Menu Management)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body CreateMenuRequest true "菜单信息"
// @Success      201 {object} response.Response{data=object{id=uint,title=string,path=string}} "菜单创建成功"
// @Failure      400 {object} response.ErrorResponse "参数错误"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Failure      403 {object} response.ErrorResponse "权限不足"
// @Failure      500 {object} response.ErrorResponse "服务器内部错误"
// @Router       /api/admin/menus [post]
// @x-permission {"scope":"admin:menus:create"}
func (h *MenuHandler) Create(c *gin.Context) {
	var req CreateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	visible := true
	if req.Visible != nil {
		visible = *req.Visible
	}

	// 调用 Use Case Handler
	result, err := h.createMenuHandler.Handle(c.Request.Context(), menuCommand.CreateMenuCommand{
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
// @Summary      获取菜单列表
// @Description  获取所有菜单的树形结构（包含父子关系）
// @Tags         管理员 - 菜单管理 (Admin - Menu Management)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} response.Response{data=[]object{id=uint,title=string,path=string,icon=string,order=int,visible=bool,children=[]object}} "菜单树"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Failure      403 {object} response.ErrorResponse "权限不足"
// @Failure      500 {object} response.ErrorResponse "服务器内部错误"
// @Router       /api/admin/menus [get]
// @x-permission {"scope":"admin:menus:read"}
func (h *MenuHandler) List(c *gin.Context) {
	// 调用 Use Case Handler
	menus, err := h.listMenusHandler.Handle(c.Request.Context(), menuQuery.ListMenusQuery{})

	if err != nil {
		response.InternalError(c, "Failed to fetch menus")
		return
	}

	response.OK(c, "success", menus)
}

// Get 获取菜单详情
//
// @Summary      获取菜单详情
// @Description  根据菜单ID获取菜单详细信息
// @Tags         管理员 - 菜单管理 (Admin - Menu Management)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "菜单ID" minimum(1)
// @Success      200 {object} response.Response{data=object{id=uint,title=string,path=string,icon=string,parent_id=uint,order=int,visible=bool}} "菜单详情"
// @Failure      400 {object} response.ErrorResponse "无效的菜单ID"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Failure      403 {object} response.ErrorResponse "权限不足"
// @Failure      404 {object} response.ErrorResponse "菜单不存在"
// @Router       /api/admin/menus/{id} [get]
// @x-permission {"scope":"admin:menus:read"}
func (h *MenuHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid menu ID")
		return
	}

	// 调用 Use Case Handler
	menu, err := h.getMenuHandler.Handle(c.Request.Context(), menuQuery.GetMenuQuery{
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
// @Summary      更新菜单信息
// @Description  管理员更新菜单的标题、路径、图标等信息
// @Tags         管理员 - 菜单管理 (Admin - Menu Management)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "菜单ID" minimum(1)
// @Param        request body UpdateMenuRequest true "更新信息"
// @Success      200 {object} response.Response{data=object{id=uint,title=string,path=string}} "菜单更新成功"
// @Failure      400 {object} response.ErrorResponse "无效的菜单ID或参数错误"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Failure      403 {object} response.ErrorResponse "权限不足"
// @Failure      404 {object} response.ErrorResponse "菜单不存在"
// @Failure      500 {object} response.ErrorResponse "服务器内部错误"
// @Router       /api/admin/menus/{id} [put]
// @x-permission {"scope":"admin:menus:update"}
func (h *MenuHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid menu ID")
		return
	}

	var req UpdateMenuRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 调用 Use Case Handler
	menu, err := h.updateMenuHandler.Handle(c.Request.Context(), menuCommand.UpdateMenuCommand{
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
// @Tags         管理员 - 菜单管理 (Admin - Menu Management)
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
// @Router       /api/admin/menus/{id} [delete]
// @x-permission {"scope":"admin:menus:delete"}
func (h *MenuHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid menu ID")
		return
	}

	// 调用 Use Case Handler
	err = h.deleteMenuHandler.Handle(c.Request.Context(), menuCommand.DeleteMenuCommand{
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
// @Summary      批量更新菜单排序
// @Description  管理员批量更新菜单的排序和父级关系
// @Tags         管理员 - 菜单管理 (Admin - Menu Management)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body ReorderMenusRequest true "菜单排序信息"
// @Success      204 "菜单排序更新成功"
// @Failure      400 {object} response.ErrorResponse "参数错误"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Failure      403 {object} response.ErrorResponse "权限不足"
// @Failure      500 {object} response.ErrorResponse "服务器内部错误"
// @Router       /api/admin/menus/reorder [post]
// @x-permission {"scope":"admin:menus:update"}
func (h *MenuHandler) Reorder(c *gin.Context) {
	var req ReorderMenusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 转换为 Command
	menus := make([]menuCommand.MenuItem, len(req.Menus))
	for i, m := range req.Menus {
		menus[i].ID = m.ID
		menus[i].Order = m.Order
		menus[i].ParentID = m.ParentID
	}

	// 调用 Use Case Handler
	err := h.reorderMenusHandler.Handle(c.Request.Context(), menuCommand.ReorderMenusCommand{
		Menus: menus,
	})

	if err != nil {
		response.InternalError(c, "Failed to update menu order")
		return
	}

	response.NoContent(c)
}
