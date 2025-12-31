package task

import (
	"context"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/task"
)

// CreateHandler 创建任务处理器。
type CreateHandler struct {
	commandRepo task.CommandRepository
}

// NewCreateHandler 创建 CreateHandler 实例。
func NewCreateHandler(commandRepo task.CommandRepository) *CreateHandler {
	return &CreateHandler{
		commandRepo: commandRepo,
	}
}

// Handle 处理创建任务命令。
func (h *CreateHandler) Handle(ctx context.Context, cmd CreateTaskCommand) (*TaskDTO, error) {
	t := &task.Task{
		OrgID:       cmd.OrgID,
		TeamID:      cmd.TeamID,
		Title:       cmd.Title,
		Description: cmd.Description,
		Status:      task.StatusPending,
		AssigneeID:  cmd.AssigneeID,
		CreatedBy:   cmd.CreatedBy,
	}

	if err := h.commandRepo.Create(ctx, t); err != nil {
		return nil, err
	}

	return ToTaskDTO(t), nil
}
