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
	createMenuHandler  *menuCommand.CreateMenuHandler
	updateMenuHandler  *menuCommand.UpdateMenuHandler
	deleteMenuHandler  *menuCommand.DeleteMenuHandler
	reorderMenusHandler *menuCommand.ReorderMenusHandler

	// Query Handlers
	getMenuHandler  *menuQuery.GetMenuHandler
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
	Title    string `json:"title" binding:"required,min=1,max=100"`
	Path     string `json:"path" binding:"required,max=255"`
	Icon     string `json:"icon" binding:"omitempty,max=100"`
	ParentID *uint  `json:"parent_id"`
	Order    int    `json:"order"`
	Visible  *bool  `json:"visible"`
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

	response.Created(c, result)
}

// List 获取菜单列表（树形结构）
func (h *MenuHandler) List(c *gin.Context) {
	// 调用 Use Case Handler
	menus, err := h.listMenusHandler.Handle(c.Request.Context(), menuQuery.ListMenusQuery{})

	if err != nil {
		response.InternalError(c, "Failed to fetch menus")
		return
	}

	response.OK(c, menus)
}

// Get 获取菜单详情
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

	response.OK(c, menu)
}

// Update 更新菜单
func (h *MenuHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid menu ID")
		return
	}

	var req UpdateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
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

	response.OK(c, menu)
}

// Delete 删除菜单
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
