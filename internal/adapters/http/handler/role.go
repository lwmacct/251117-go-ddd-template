package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/role"
)

// RoleHandler handles role management operations
type RoleHandler struct {
	roleRepo       role.RoleRepository
	permissionRepo role.PermissionRepository
}

// NewRoleHandler creates a new RoleHandler instance
func NewRoleHandler(roleRepo role.RoleRepository, permissionRepo role.PermissionRepository) *RoleHandler {
	return &RoleHandler{
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
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

	newRole := &role.Role{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		IsSystem:    false,
	}

	if err := h.roleRepo.Create(c.Request.Context(), newRole); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create role"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "role created successfully",
		"role":    newRole,
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

	roles, total, err := h.roleRepo.List(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list roles"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"roles": roles,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
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

	r, err := h.roleRepo.FindByID(c.Request.Context(), uint(id))
	if err != nil || r == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "role not found"})
		return
	}

	c.JSON(http.StatusOK, r)
}

// UpdateRole updates a role
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role ID"})
		return
	}

	r, err := h.roleRepo.FindByID(c.Request.Context(), uint(id))
	if err != nil || r == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "role not found"})
		return
	}

	if r.IsSystem {
		c.JSON(http.StatusForbidden, gin.H{"error": "cannot modify system role"})
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

	if req.DisplayName != nil {
		r.DisplayName = *req.DisplayName
	}
	if req.Description != nil {
		r.Description = *req.Description
	}

	if err := h.roleRepo.Update(c.Request.Context(), r); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "role updated successfully",
		"role":    r,
	})
}

// DeleteRole deletes a role
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role ID"})
		return
	}

	r, err := h.roleRepo.FindByID(c.Request.Context(), uint(id))
	if err != nil || r == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "role not found"})
		return
	}

	if r.IsSystem {
		c.JSON(http.StatusForbidden, gin.H{"error": "cannot delete system role"})
		return
	}

	if err := h.roleRepo.Delete(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete role"})
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

	if err := h.roleRepo.SetPermissions(c.Request.Context(), uint(id), req.PermissionIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set permissions"})
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

	permissions, total, err := h.permissionRepo.List(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list permissions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"permissions": permissions,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}
