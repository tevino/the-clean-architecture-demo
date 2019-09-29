package entity

// ItemType indicates the type of an item.
type ItemType int

const (
	// ItemTypeTask indicates a task.
	ItemTypeTask ItemType = iota
	// ItemTypeProject indicates a project usually contains multiple tasks.
	ItemTypeProject
	// ItemTypeCategory indicates a category to sort tasks or projects.
	ItemTypeCategory
)
