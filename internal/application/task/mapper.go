package task

import "github.com/lwmacct/251117-go-ddd-template/internal/domain/task"

// ToTaskDTO 将实体转换为 DTO。
func ToTaskDTO(t *task.Task) *TaskDTO {
	if t == nil {
		return nil
	}
	return &TaskDTO{
		ID:          t.ID,
		OrgID:       t.OrgID,
		TeamID:      t.TeamID,
		Title:       t.Title,
		Description: t.Description,
		Status:      string(t.Status),
		AssigneeID:  t.AssigneeID,
		CreatedBy:   t.CreatedBy,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

// ToTaskDTOs 将实体切片转换为 DTO 切片。
func ToTaskDTOs(tasks []*task.Task) []*TaskDTO {
	if len(tasks) == 0 {
		return nil
	}
	dtos := make([]*TaskDTO, 0, len(tasks))
	for _, t := range tasks {
		if dto := ToTaskDTO(t); dto != nil {
			dtos = append(dtos, dto)
		}
	}
	return dtos
}
