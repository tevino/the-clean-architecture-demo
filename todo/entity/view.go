package entity

// TaskView represents a set of conditions to filter tasks.
type TaskView struct {
	// TODO:
}

type ConditionTarget int

const (
	Title ConditionTarget = iota
	Description
	CreatedAt
	UpdatedAt
	ParentTaskID
)

type Condition struct {
	Type   ConditionType
	Target ConditionTarget
	Value  string
}

type Composition struct {
	Type       ComposeType
	Conditions []*Condition
}
type ComposeType int

const (
	And ComposeType = iota
	Or
)

type ConditionType int

const (
	Equal ConditionType = iota
	NotEqual
	Contains
	NotContain
)
