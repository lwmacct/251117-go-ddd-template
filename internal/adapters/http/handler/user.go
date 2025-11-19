package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	appUserDTO "github.com/lwmacct/251117-go-ddd-template/internal/application/user"
	userCommand "github.com/lwmacct/251117-go-ddd-template/internal/application/user/command"
	userQuery "github.com/lwmacct/251117-go-ddd-template/internal/application/user/query"
)

// UserHandler 用户处理器（新架构）
type UserHandler struct {
	createUserHandler *userCommand.CreateUserHandler
	updateUserHandler *userCommand.UpdateUserHandler
	deleteUserHandler *userCommand.DeleteUserHandler
	getUserHandler    *userQuery.GetUserHandler
	listUsersHandler  *userQuery.ListUsersHandler
}

// NewUserHandler 创建用户处理器
func NewUserHandler(
	createUserHandler *userCommand.CreateUserHandler,
	updateUserHandler *userCommand.UpdateUserHandler,
	deleteUserHandler *userCommand.DeleteUserHandler,
	getUserHandler *userQuery.GetUserHandler,
	listUsersHandler *userQuery.ListUsersHandler,
) *UserHandler {
	return &UserHandler{
		createUserHandler: createUserHandler,
		updateUserHandler: updateUserHandler,
		deleteUserHandler: deleteUserHandler,
		getUserHandler:    getUserHandler,
		listUsersHandler:  listUsersHandler,
	}
}

// Create 创建用户
// POST /api/users
func (h *UserHandler) Create(c *gin.Context) {
	var req appUserDTO.CreateUserDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 调用 Use Case Handler
	result, err := h.createUserHandler.Handle(c.Request.Context(), userCommand.CreateUserCommand{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		FullName: req.FullName,
	})

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{
		"message": "user created successfully",
		"data": gin.H{
			"user_id":  result.UserID,
			"username": result.Username,
			"email":    result.Email,
		},
	})
}

// GetByID 获取用户详情
// GET /api/users/:id
func (h *UserHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid user id"})
		return
	}

	// 调用 Query Handler
	user, err := h.getUserHandler.Handle(c.Request.Context(), userQuery.GetUserQuery{
		UserID:    uint(id),
		WithRoles: true, // 包含角色信息
	})

	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"data": user,
	})
}

// List 获取用户列表
// GET /api/users?page=1&limit=10
func (h *UserHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	// 调用 Query Handler
	result, err := h.listUsersHandler.Handle(c.Request.Context(), userQuery.ListUsersQuery{
		Offset: offset,
		Limit:  limit,
	})

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"data": gin.H{
			"users": result.Users,
			"total": result.Total,
			"page":  page,
			"limit": limit,
		},
	})
}

// Update 更新用户
// PUT /api/users/:id
func (h *UserHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid user id"})
		return
	}

	var req appUserDTO.UpdateUserDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 调用 Command Handler
	_, err = h.updateUserHandler.Handle(c.Request.Context(), userCommand.UpdateUserCommand{
		UserID:   uint(id),
		FullName: req.FullName,
		Avatar:   req.Avatar,
		Bio:      req.Bio,
		Status:   req.Status,
	})

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "user updated successfully",
	})
}

// Delete 删除用户
// DELETE /api/users/:id
func (h *UserHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid user id"})
		return
	}

	// 调用 Command Handler
	err = h.deleteUserHandler.Handle(c.Request.Context(), userCommand.DeleteUserCommand{
		UserID: uint(id),
	})

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "user deleted successfully",
	})
}
