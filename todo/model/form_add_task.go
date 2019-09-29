package model

import "time"

// FormAddTask represents the input from user while adding a new task.
type FormAddTask struct {
	Title       string
	Due         time.Time
	Description string
	Type        TaskType
	ParentID    int64
	Order       uint64
}
