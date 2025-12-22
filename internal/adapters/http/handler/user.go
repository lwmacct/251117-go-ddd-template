// Package handler 提供 HTTP 请求处理器。
//
// 本包是适配器层的核心组件，职责：
//   - 请求绑定：解析 HTTP 请求参数到 DTO
//   - 响应转换：将 Use Case 结果转换为 HTTP 响应
//   - 错误处理：统一的错误响应格式
//
// 设计原则：
//   - Handler 不包含业务逻辑，仅做 HTTP 适配
//   - 业务逻辑委托给 Application 层的 Command/Query Handler
//   - 使用 response 包提供统一的响应格式
//
// 文件命名：{module}.go（单数形式）
// 例如：user.go, role.go, auth.go
package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/response"
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
		response.ValidationError(c, err.Error())
		return
	}

	// 调用 Use Case Handler
	result, err := h.createUserHandler.Handle(c.Request.Context(), userCommand.CreateUserCommand{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		FullName: req.FullName,
		RoleIDs:  req.RoleIDs,
	})

	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, "user created successfully", gin.H{
		"user_id":  result.UserID,
		"username": result.Username,
		"email":    result.Email,
	})
}

// GetByID 获取用户详情
// GET /api/users/:id
func (h *UserHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	// 调用 Query Handler
	user, err := h.getUserHandler.Handle(c.Request.Context(), userQuery.GetUserQuery{
		UserID:    uint(id),
		WithRoles: true, // 包含角色信息
	})

	if err != nil {
		response.NotFound(c, "user")
		return
	}

	response.OK(c, "success", user)
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
		response.InternalError(c, err.Error())
		return
	}

	meta := response.NewPaginationMeta(int(result.Total), page, limit)
	response.List(c, "success", result.Users, meta)
}

// Update 更新用户
// PUT /api/users/:id
func (h *UserHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	var req appUserDTO.UpdateUserDTO
	if err = c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
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
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, "user updated successfully", nil)
}

// Delete 删除用户
// DELETE /api/users/:id
func (h *UserHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	// 调用 Command Handler
	err = h.deleteUserHandler.Handle(c.Request.Context(), userCommand.DeleteUserCommand{
		UserID: uint(id),
	})

	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, "user deleted successfully", nil)
}
