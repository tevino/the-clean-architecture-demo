package entity

// ItemState indicate the state of an item.
type ItemState int

const (
	// ItemStateNormal indicates the default state of an item.
	ItemStateNormal ItemState = iota
	// ItemStateCompleted indicates the item has been completed.
	ItemStateCompleted
)
