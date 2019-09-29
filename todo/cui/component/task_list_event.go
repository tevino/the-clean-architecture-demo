package component

type TaskListEventType int

const (
	TaskListEventAfterUpdate TaskListEventType = iota
	EventChangeTaskState
	EventInsertTaskWithOrder
)

type TaskListEvent struct {
	Type  TaskListEventType
	Order uint64
}
