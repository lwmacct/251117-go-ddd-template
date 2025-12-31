package persistence

import (
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/task"
	"gorm.io/gorm"
)

// TaskRepositories 聚合任务读写仓储。
type TaskRepositories struct {
	Command task.CommandRepository
	Query   task.QueryRepository
}

// NewTaskRepositories 创建任务仓储聚合实例。
func NewTaskRepositories(db *gorm.DB) TaskRepositories {
	return TaskRepositories{
		Command: NewTaskCommandRepository(db),
		Query:   NewTaskQueryRepository(db),
	}
}
