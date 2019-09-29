package model

// TaskState indicates the state of a task.
type TaskState int

// All TaskState(s).
const (
	TaskStateNormal TaskState = iota
	TaskStateCompleted
)
