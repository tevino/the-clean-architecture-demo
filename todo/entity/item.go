package entity

import "time"

// Item could be a task or a project contains multiple tasks.
type Item struct {
	ID           int64
	Title        string
	Description  string
	Type         ItemType
	State        ItemState
	Due          time.Time
	CompletedAt  time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	ParentItemID int64
	Order        uint64
}

// RootID is the ID of RootItem.
const RootID = 0

// RootItem is a virtual item, items on the top most level(without a real parent) should have its ParentItemID set to this.
var RootItem = &Item{ID: RootID}
