package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	userdto "github.com/lwmacct/251117-go-ddd-template/internal/application/user"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
	"golang.org/x/crypto/bcrypt"
)

// AdminUserHandler handles admin user management operations
type AdminUserHandler struct {
	userCommandRepo user.CommandRepository
	userQueryRepo   user.QueryRepository
}

// NewAdminUserHandler creates a new AdminUserHandler instance
func NewAdminUserHandler(userCommandRepo user.CommandRepository, userQueryRepo user.QueryRepository) *AdminUserHandler {
	return &AdminUserHandler{
		userCommandRepo: userCommandRepo,
		userQueryRepo:   userQueryRepo,
	}
}

// CreateUser creates a new user (admin only)
func (h *AdminUserHandler) CreateUser(c *gin.Context) {
	var dto userdto.CreateUserDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	newUser := &user.User{
		Username: dto.Username,
		Email:    dto.Email,
		Password: string(hashedPassword),
		FullName: dto.FullName,
		Status:   "active",
	}

	if err := h.userCommandRepo.Create(c.Request.Context(), newUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user created successfully",
		"user":    userdto.ToUserResponse(newUser),
	})
}

// ListUsers lists all users with pagination (admin only)
func (h *AdminUserHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit
	users, err := h.userQueryRepo.List(c.Request.Context(), offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list users"})
		return
	}

	total, err := h.userQueryRepo.Count(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count users"})
		return
	}

	responses := make([]*userdto.UserResponse, 0, len(users))
	for _, u := range users {
		responses = append(responses, userdto.ToUserResponse(u))
	}

	c.JSON(http.StatusOK, gin.H{
		"users": responses,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// GetUser gets a user by ID (admin only)
func (h *AdminUserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	u, err := h.userQueryRepo.GetByIDWithRoles(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, userdto.ToUserWithRolesResponse(u))
}

// UpdateUser updates a user (admin only)
func (h *AdminUserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var dto userdto.UpdateUserDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := h.userQueryRepo.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if dto.FullName != nil {
		u.FullName = *dto.FullName
	}
	if dto.Avatar != nil {
		u.Avatar = *dto.Avatar
	}
	if dto.Bio != nil {
		u.Bio = *dto.Bio
	}
	if dto.Status != nil {
		u.Status = *dto.Status
	}

	if err := h.userCommandRepo.Update(c.Request.Context(), u); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user updated successfully",
		"user":    userdto.ToUserResponse(u),
	})
}

// DeleteUser deletes a user (admin only)
func (h *AdminUserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	if err := h.userCommandRepo.Delete(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}

// AssignRoles assigns roles to a user (admin only)
func (h *AdminUserHandler) AssignRoles(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var req struct {
		RoleIDs []uint `json:"role_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.userCommandRepo.AssignRoles(c.Request.Context(), uint(id), req.RoleIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to assign roles"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "roles assigned successfully"})
}
