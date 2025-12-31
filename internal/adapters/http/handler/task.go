package handler

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/response"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/task"
	taskDomain "github.com/lwmacct/251117-go-ddd-template/internal/domain/task"
)

// ListTasksQuery 任务列表查询参数
type ListTasksQuery struct {
	response.PaginationQueryDTO
}

// ToQuery 转换为 Application 层 Query 对象
func (q *ListTasksQuery) ToQuery(orgID, teamID uint) task.ListTasksQuery {
	return task.ListTasksQuery{
		OrgID:  orgID,
		TeamID: teamID,
		Offset: q.GetOffset(),
		Limit:  q.GetLimit(),
	}
}

// TaskHandler 团队任务管理 Handler
type TaskHandler struct {
	createHandler *task.CreateHandler
	updateHandler *task.UpdateHandler
	deleteHandler *task.DeleteHandler
	getHandler    *task.GetHandler
	listHandler   *task.ListHandler
}

// NewTaskHandler 创建团队任务管理 Handler
func NewTaskHandler(
	createHandler *task.CreateHandler,
	updateHandler *task.UpdateHandler,
	deleteHandler *task.DeleteHandler,
	getHandler *task.GetHandler,
	listHandler *task.ListHandler,
) *TaskHandler {
	return &TaskHandler{
		createHandler: createHandler,
		updateHandler: updateHandler,
		deleteHandler: deleteHandler,
		getHandler:    getHandler,
		listHandler:   listHandler,
	}
}

// Create 创建任务
//
//	@Summary		创建任务
//	@Description	在团队中创建新任务
//	@Tags			Organization - Team Task Management
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			org_id	path		int									true	"组织ID"	minimum(1)
//	@Param			team_id	path		int									true	"团队ID"	minimum(1)
//	@Param			request	body		task.CreateTaskDTO					true	"任务信息"
//	@Success		201		{object}	response.DataResponse[task.TaskDTO]	"任务创建成功"
//	@Failure		400		{object}	response.ErrorResponse				"参数错误"
//	@Failure		401		{object}	response.ErrorResponse				"未授权"
//	@Failure		403		{object}	response.ErrorResponse				"权限不足"
//	@Failure		500		{object}	response.ErrorResponse				"服务器内部错误"
//	@Router			/api/org/{org_id}/teams/{team_id}/tasks [post]
func (h *TaskHandler) Create(c *gin.Context) {
	orgID := c.GetUint("org_id")
	teamID := c.GetUint("team_id")
	userID := c.GetUint("user_id")

	var req task.CreateTaskDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result, err := h.createHandler.Handle(c.Request.Context(), task.CreateTaskCommand{
		OrgID:       orgID,
		TeamID:      teamID,
		Title:       req.Title,
		Description: req.Description,
		AssigneeID:  req.AssigneeID,
		CreatedBy:   userID,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, response.MsgCreated, result)
}

// List 任务列表
//
//	@Summary		任务列表
//	@Description	分页获取团队任务列表
//	@Tags			Organization - Team Task Management
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			org_id	path		int										true	"组织ID"	minimum(1)
//	@Param			team_id	path		int										true	"团队ID"	minimum(1)
//	@Param			params	query		ListTasksQuery							false	"查询参数"
//	@Success		200		{object}	response.PagedResponse[task.TaskDTO]	"任务列表"
//	@Failure		401		{object}	response.ErrorResponse					"未授权"
//	@Failure		403		{object}	response.ErrorResponse					"权限不足"
//	@Failure		500		{object}	response.ErrorResponse					"服务器内部错误"
//	@Router			/api/org/{org_id}/teams/{team_id}/tasks [get]
func (h *TaskHandler) List(c *gin.Context) {
	orgID := c.GetUint("org_id")
	teamID := c.GetUint("team_id")

	var query ListTasksQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result, err := h.listHandler.Handle(c.Request.Context(), query.ToQuery(orgID, teamID))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	meta := response.NewPaginationMeta(int(result.Total), query.GetPage(), query.GetLimit())
	response.List(c, response.MsgSuccess, result.Items, meta)
}

// Get 任务详情
//
//	@Summary		任务详情
//	@Description	获取任务详细信息
//	@Tags			Organization - Team Task Management
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			org_id	path		int									true	"组织ID"	minimum(1)
//	@Param			team_id	path		int									true	"团队ID"	minimum(1)
//	@Param			id		path		int									true	"任务ID"	minimum(1)
//	@Success		200		{object}	response.DataResponse[task.TaskDTO]	"任务详情"
//	@Failure		401		{object}	response.ErrorResponse				"未授权"
//	@Failure		403		{object}	response.ErrorResponse				"权限不足"
//	@Failure		404		{object}	response.ErrorResponse				"任务不存在"
//	@Failure		500		{object}	response.ErrorResponse				"服务器内部错误"
//	@Router			/api/org/{org_id}/teams/{team_id}/tasks/{id} [get]
func (h *TaskHandler) Get(c *gin.Context) {
	orgID := c.GetUint("org_id")
	teamID := c.GetUint("team_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ValidationError(c, "无效的任务ID")
		return
	}

	result, err := h.getHandler.Handle(c.Request.Context(), task.GetTaskQuery{
		OrgID:  orgID,
		TeamID: teamID,
		ID:     uint(id),
	})
	if err != nil {
		if errors.Is(err, taskDomain.ErrTaskNotFound) {
			response.NotFoundMessage(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, response.MsgSuccess, result)
}

// Update 更新任务
//
//	@Summary		更新任务
//	@Description	更新任务信息或状态
//	@Tags			Organization - Team Task Management
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			org_id	path		int									true	"组织ID"	minimum(1)
//	@Param			team_id	path		int									true	"团队ID"	minimum(1)
//	@Param			id		path		int									true	"任务ID"	minimum(1)
//	@Param			request	body		task.UpdateTaskDTO					true	"更新信息"
//	@Success		200		{object}	response.DataResponse[task.TaskDTO]	"更新成功"
//	@Failure		400		{object}	response.ErrorResponse				"参数错误或状态转换无效"
//	@Failure		401		{object}	response.ErrorResponse				"未授权"
//	@Failure		403		{object}	response.ErrorResponse				"权限不足"
//	@Failure		404		{object}	response.ErrorResponse				"任务不存在"
//	@Failure		500		{object}	response.ErrorResponse				"服务器内部错误"
//	@Router			/api/org/{org_id}/teams/{team_id}/tasks/{id} [put]
func (h *TaskHandler) Update(c *gin.Context) {
	orgID := c.GetUint("org_id")
	teamID := c.GetUint("team_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ValidationError(c, "无效的任务ID")
		return
	}

	var req task.UpdateTaskDTO
	if err = c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result, err := h.updateHandler.Handle(c.Request.Context(), task.UpdateTaskCommand{
		OrgID:       orgID,
		TeamID:      teamID,
		ID:          uint(id),
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		AssigneeID:  req.AssigneeID,
	})
	if err != nil {
		if errors.Is(err, taskDomain.ErrTaskNotFound) {
			response.NotFoundMessage(c, err.Error())
			return
		}
		if errors.Is(err, taskDomain.ErrInvalidStatusTransition) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, response.MsgUpdated, result)
}

// Delete 删除任务
//
//	@Summary		删除任务
//	@Description	删除任务（软删除）
//	@Tags			Organization - Team Task Management
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			org_id	path		int						true	"组织ID"	minimum(1)
//	@Param			team_id	path		int						true	"团队ID"	minimum(1)
//	@Param			id		path		int						true	"任务ID"	minimum(1)
//	@Success		200		{object}	response.MessageResponse	"删除成功"
//	@Failure		401		{object}	response.ErrorResponse	"未授权"
//	@Failure		403		{object}	response.ErrorResponse	"权限不足"
//	@Failure		404		{object}	response.ErrorResponse	"任务不存在"
//	@Failure		500		{object}	response.ErrorResponse	"服务器内部错误"
//	@Router			/api/org/{org_id}/teams/{team_id}/tasks/{id} [delete]
func (h *TaskHandler) Delete(c *gin.Context) {
	orgID := c.GetUint("org_id")
	teamID := c.GetUint("team_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ValidationError(c, "无效的任务ID")
		return
	}

	if err := h.deleteHandler.Handle(c.Request.Context(), task.DeleteTaskCommand{
		OrgID:  orgID,
		TeamID: teamID,
		ID:     uint(id),
	}); err != nil {
		if errors.Is(err, taskDomain.ErrTaskNotFound) {
			response.NotFoundMessage(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, response.MsgDeleted, nil)
}
