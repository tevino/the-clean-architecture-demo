package model

import "time"

// Task is the response of a task to user.
type Task struct {
	ID          int64
	Title       string
	Type        TaskType
	State       TaskState
	Due         time.Time
	Description string
	Order       uint64
}
