package task

import "errors"

var (
	// ErrTaskNotFound 任务不存在。
	ErrTaskNotFound = errors.New("task not found")
	// ErrInvalidStatusTransition 无效的状态转换。
	ErrInvalidStatusTransition = errors.New("invalid status transition")
)
