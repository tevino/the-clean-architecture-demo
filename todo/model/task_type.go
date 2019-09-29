package model

// TaskType indicates the type of the task.
type TaskType int

// All task types.
const (
	TaskTypeCategory TaskType = iota
	TaskTypeTask
)
