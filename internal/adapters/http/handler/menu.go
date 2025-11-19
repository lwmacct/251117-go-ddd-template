package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/response"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/menu"
)

type MenuHandler struct {
	menuRepo menu.Repository
}

func NewMenuHandler(menuRepo menu.Repository) *MenuHandler {
	return &MenuHandler{
		menuRepo: menuRepo,
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

	m := &menu.Menu{
		Title:    req.Title,
		Path:     req.Path,
		Icon:     req.Icon,
		ParentID: req.ParentID,
		Order:    req.Order,
		Visible:  visible,
	}

	if err := h.menuRepo.Create(c.Request.Context(), m); err != nil {
		response.InternalError(c, "Failed to create menu")
		return
	}

	response.Created(c, m)
}

// List 获取菜单列表（树形结构）
func (h *MenuHandler) List(c *gin.Context) {
	menus, err := h.menuRepo.FindAll(c.Request.Context())
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

	m, err := h.menuRepo.FindByID(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, "Menu")
		return
	}

	response.OK(c, m)
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

	m, err := h.menuRepo.FindByID(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, "Menu")
		return
	}

	// 更新字段
	if req.Title != nil {
		m.Title = *req.Title
	}
	if req.Path != nil {
		m.Path = *req.Path
	}
	if req.Icon != nil {
		m.Icon = *req.Icon
	}
	if req.ParentID != nil {
		m.ParentID = req.ParentID
	}
	if req.Order != nil {
		m.Order = *req.Order
	}
	if req.Visible != nil {
		m.Visible = *req.Visible
	}

	if err := h.menuRepo.Update(c.Request.Context(), m); err != nil {
		response.InternalError(c, "Failed to update menu")
		return
	}

	response.OK(c, m)
}

// Delete 删除菜单
func (h *MenuHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid menu ID")
		return
	}

	// 检查是否有子菜单
	children, err := h.menuRepo.FindByParentID(c.Request.Context(), ptrUint(uint(id)))
	if err != nil {
		response.InternalError(c, "Failed to check children")
		return
	}

	if len(children) > 0 {
		response.BadRequest(c, "Cannot delete menu with children")
		return
	}

	if err := h.menuRepo.Delete(c.Request.Context(), uint(id)); err != nil {
		response.InternalError(c, "Failed to delete menu")
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

	// 转换格式
	menus := make([]struct {
		ID       uint
		Order    int
		ParentID *uint
	}, len(req.Menus))

	for i, m := range req.Menus {
		menus[i].ID = m.ID
		menus[i].Order = m.Order
		menus[i].ParentID = m.ParentID
	}

	if err := h.menuRepo.UpdateOrder(c.Request.Context(), menus); err != nil {
		response.InternalError(c, "Failed to update menu order")
		return
	}

	response.NoContent(c)
}

// ptrUint 辅助函数：创建 uint 指针
func ptrUint(u uint) *uint {
	return &u
}
