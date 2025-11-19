package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
	"golang.org/x/crypto/bcrypt"
)

// AdminUserHandler handles admin user management operations
type AdminUserHandler struct {
	userRepo user.Repository
}

// NewAdminUserHandler creates a new AdminUserHandler instance
func NewAdminUserHandler(userRepo user.Repository) *AdminUserHandler {
	return &AdminUserHandler{
		userRepo: userRepo,
	}
}

// CreateUser creates a new user (admin only)
func (h *AdminUserHandler) CreateUser(c *gin.Context) {
	var dto user.UserCreateDTO
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

	if err := h.userRepo.Create(c.Request.Context(), newUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user created successfully",
		"user":    newUser.ToResponse(),
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
	users, err := h.userRepo.List(c.Request.Context(), offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list users"})
		return
	}

	total, err := h.userRepo.Count(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count users"})
		return
	}

	responses := make([]*user.UserResponse, 0, len(users))
	for _, u := range users {
		responses = append(responses, u.ToResponse())
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

	u, err := h.userRepo.GetByIDWithRoles(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, u.ToResponse())
}

// UpdateUser updates a user (admin only)
func (h *AdminUserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var dto user.UserUpdateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := h.userRepo.GetByID(c.Request.Context(), uint(id))
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

	if err := h.userRepo.Update(c.Request.Context(), u); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user updated successfully",
		"user":    u.ToResponse(),
	})
}

// DeleteUser deletes a user (admin only)
func (h *AdminUserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	if err := h.userRepo.Delete(c.Request.Context(), uint(id)); err != nil {
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

	if err := h.userRepo.AssignRoles(c.Request.Context(), uint(id), req.RoleIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to assign roles"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "roles assigned successfully"})
}
