package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/response"
	appUserDTO "github.com/lwmacct/251117-go-ddd-template/internal/application/user"
	userCommand "github.com/lwmacct/251117-go-ddd-template/internal/application/user/command"
	userQuery "github.com/lwmacct/251117-go-ddd-template/internal/application/user/query"
)

// AdminUserHandler handles admin user management operations
type AdminUserHandler struct {
	createUserHandler  *userCommand.CreateUserHandler
	updateUserHandler  *userCommand.UpdateUserHandler
	deleteUserHandler  *userCommand.DeleteUserHandler
	assignRolesHandler *userCommand.AssignRolesHandler
	getUserHandler     *userQuery.GetUserHandler
	listUsersHandler   *userQuery.ListUsersHandler
}

// NewAdminUserHandler creates a new AdminUserHandler instance
func NewAdminUserHandler(
	createUserHandler *userCommand.CreateUserHandler,
	updateUserHandler *userCommand.UpdateUserHandler,
	deleteUserHandler *userCommand.DeleteUserHandler,
	assignRolesHandler *userCommand.AssignRolesHandler,
	getUserHandler *userQuery.GetUserHandler,
	listUsersHandler *userQuery.ListUsersHandler,
) *AdminUserHandler {
	return &AdminUserHandler{
		createUserHandler:  createUserHandler,
		updateUserHandler:  updateUserHandler,
		deleteUserHandler:  deleteUserHandler,
		assignRolesHandler: assignRolesHandler,
		getUserHandler:     getUserHandler,
		listUsersHandler:   listUsersHandler,
	}
}

// CreateUser creates a new user (admin only)
func (h *AdminUserHandler) CreateUser(c *gin.Context) {
	var dto appUserDTO.CreateUserDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result, err := h.createUserHandler.Handle(c.Request.Context(), userCommand.CreateUserCommand{
		Username: dto.Username,
		Email:    dto.Email,
		Password: dto.Password,
		FullName: dto.FullName,
		RoleIDs:  dto.RoleIDs,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	createdUser, err := h.getUserHandler.Handle(c.Request.Context(), userQuery.GetUserQuery{
		UserID:    result.UserID,
		WithRoles: true,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, gin.H{
		"message": "user created successfully",
		"user":    createdUser,
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

	result, err := h.listUsersHandler.Handle(c.Request.Context(), userQuery.ListUsersQuery{
		Offset: offset,
		Limit:  limit,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	meta := response.NewPaginationMeta(int(result.Total), page, limit)
	response.List(c, gin.H{
		"users": result.Users,
	}, meta)
}

// GetUser gets a user by ID (admin only)
func (h *AdminUserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid user ID")
		return
	}

	userResp, err := h.getUserHandler.Handle(c.Request.Context(), userQuery.GetUserQuery{
		UserID:    uint(id),
		WithRoles: true,
	})
	if err != nil {
		response.NotFound(c, "user")
		return
	}

	response.OK(c, userResp)
}

// UpdateUser updates a user (admin only)
func (h *AdminUserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid user ID")
		return
	}

	var dto appUserDTO.UpdateUserDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	_, err = h.updateUserHandler.Handle(c.Request.Context(), userCommand.UpdateUserCommand{
		UserID:   uint(id),
		FullName: dto.FullName,
		Avatar:   dto.Avatar,
		Bio:      dto.Bio,
		Status:   dto.Status,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	updatedUser, err := h.getUserHandler.Handle(c.Request.Context(), userQuery.GetUserQuery{
		UserID:    uint(id),
		WithRoles: true,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, gin.H{
		"message": "user updated successfully",
		"user":    updatedUser,
	})
}

// DeleteUser deletes a user (admin only)
func (h *AdminUserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid user ID")
		return
	}

	if err := h.deleteUserHandler.Handle(c.Request.Context(), userCommand.DeleteUserCommand{
		UserID: uint(id),
	}); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "user deleted successfully"})
}

// AssignRoles assigns roles to a user (admin only)
func (h *AdminUserHandler) AssignRoles(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid user ID")
		return
	}

	var req appUserDTO.AssignRolesDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	if err := h.assignRolesHandler.Handle(c.Request.Context(), userCommand.AssignRolesCommand{
		UserID:  uint(id),
		RoleIDs: req.RoleIDs,
	}); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "roles assigned successfully"})
}
