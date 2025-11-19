package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

// CreateRole creates a new role
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required,min=2,max=50"`
		DisplayName string `json:"display_name" binding:"required,max=100"`
		Description string `json:"description" binding:"max=255"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用 Use Case Handler
	result, err := h.createRoleHandler.Handle(c.Request.Context(), roleCommand.CreateRoleCommand{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "role created successfully",
		"data":    result,
	})
}

// ListRoles lists all roles
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": result.Roles,
		"pagination": gin.H{
			"page":  result.Page,
			"limit": result.Limit,
			"total": result.Total,
		},
	})
}

// GetRole gets a role by ID
func (h *RoleHandler) GetRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role ID"})
		return
	}

	// 调用 Use Case Handler
	result, err := h.getRoleHandler.Handle(c.Request.Context(), roleQuery.GetRoleQuery{
		RoleID: uint(id),
	})

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

// UpdateRole updates a role
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role ID"})
		return
	}

	var req struct {
		DisplayName *string `json:"display_name" binding:"omitempty,max=100"`
		Description *string `json:"description" binding:"omitempty,max=255"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用 Use Case Handler
	result, err := h.updateRoleHandler.Handle(c.Request.Context(), roleCommand.UpdateRoleCommand{
		RoleID:      uint(id),
		DisplayName: req.DisplayName,
		Description: req.Description,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "role updated successfully",
		"data":    result,
	})
}

// DeleteRole deletes a role
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role ID"})
		return
	}

	// 调用 Use Case Handler
	err = h.deleteRoleHandler.Handle(c.Request.Context(), roleCommand.DeleteRoleCommand{
		RoleID: uint(id),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "role deleted successfully"})
}

// SetPermissions sets permissions for a role
func (h *RoleHandler) SetPermissions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role ID"})
		return
	}

	var req struct {
		PermissionIDs []uint `json:"permission_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用 Use Case Handler
	err = h.setPermissionsHandler.Handle(c.Request.Context(), roleCommand.SetPermissionsCommand{
		RoleID:        uint(id),
		PermissionIDs: req.PermissionIDs,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "permissions set successfully"})
}

// ListPermissions lists all permissions
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": result.Permissions,
		"pagination": gin.H{
			"page":  result.Page,
			"limit": result.Limit,
			"total": result.Total,
		},
	})
}
