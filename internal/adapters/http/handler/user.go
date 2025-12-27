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
	"github.com/lwmacct/251117-go-ddd-template/internal/application/user"
)

// UserHandler 用户处理器（新架构）
type UserHandler struct {
	createUserHandler *user.CreateHandler
	updateUserHandler *user.UpdateHandler
	deleteUserHandler *user.DeleteHandler
	getUserHandler    *user.GetHandler
	listUsersHandler  *user.ListHandler
}

// NewUserHandler 创建用户处理器
func NewUserHandler(
	createUserHandler *user.CreateHandler,
	updateUserHandler *user.UpdateHandler,
	deleteUserHandler *user.DeleteHandler,
	getUserHandler *user.GetHandler,
	listUsersHandler *user.ListHandler,
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
//
// @Summary      创建用户
// @Description  创建新用户账号
// @Tags         用户 (User)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body user.CreateDTO true "用户信息"
// @Success      201 {object} response.DataResponse[user.UserWithRolesDTO] "用户创建成功"
// @Failure      400 {object} response.ErrorResponse "参数错误"
// @Failure      500 {object} response.ErrorResponse "服务器内部错误"
// @Router       /api/users [post]
func (h *UserHandler) Create(c *gin.Context) {
	var req user.CreateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	// 调用 Use Case Handler
	result, err := h.createUserHandler.Handle(c.Request.Context(), user.CreateCommand(req))

	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, "user created successfully", result)
}

// GetByID 获取用户详情
//
// @Summary      获取用户详情
// @Description  根据用户ID获取用户详细信息（包含角色信息）
// @Tags         用户 (User)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "用户ID" minimum(1)
// @Success      200 {object} response.DataResponse[user.UserWithRolesDTO] "用户详情"
// @Failure      400 {object} response.ErrorResponse "无效的用户ID"
// @Failure      404 {object} response.ErrorResponse "用户不存在"
// @Router       /api/users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	// 调用 Query Handler
	user, err := h.getUserHandler.Handle(c.Request.Context(), user.GetQuery{
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
//
// @Summary      获取用户列表
// @Description  分页获取用户列表（包含角色信息）
// @Tags         用户 (User)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        params query handler.ListUsersQuery false "查询参数"
// @Success      200 {object} response.PagedResponse[user.UserWithRolesDTO] "用户列表"
// @Failure      400 {object} response.ErrorResponse "参数错误"
// @Failure      500 {object} response.ErrorResponse "服务器内部错误"
// @Router       /api/users [get]
func (h *UserHandler) List(c *gin.Context) {
	var q ListUsersQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result, err := h.listUsersHandler.Handle(c.Request.Context(), q.ToQuery())
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	meta := response.NewPaginationMeta(int(result.Total), q.GetPage(), q.GetLimit())
	response.List(c, "success", result.Users, meta)
}

// Update 更新用户
//
// @Summary      更新用户信息
// @Description  更新用户的基本信息
// @Tags         用户 (User)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "用户ID" minimum(1)
// @Param        request body user.UpdateDTO true "更新信息"
// @Success      200 {object} response.MessageResponse "用户更新成功"
// @Failure      400 {object} response.ErrorResponse "无效的用户ID或参数错误"
// @Failure      500 {object} response.ErrorResponse "服务器内部错误"
// @Router       /api/users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	var req user.UpdateDTO
	if err = c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	// 调用 Command Handler
	_, err = h.updateUserHandler.Handle(c.Request.Context(), user.UpdateCommand{
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
//
// @Summary      删除用户
// @Description  删除指定用户
// @Tags         用户 (User)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "用户ID" minimum(1)
// @Success      200 {object} response.MessageResponse "用户删除成功"
// @Failure      400 {object} response.ErrorResponse "无效的用户ID"
// @Failure      500 {object} response.ErrorResponse "服务器内部错误"
// @Router       /api/users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	// 调用 Command Handler
	err = h.deleteUserHandler.Handle(c.Request.Context(), user.DeleteCommand{
		UserID: uint(id),
	})

	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, "user deleted successfully", nil)
}
