package component

import ui "github.com/gizak/termui/v3"

// Component contains interfaces every Component should implement.
type Component interface {
	updater
}

// InteractiveComponent is a component that interacts to user input.
type InteractiveComponent interface {
	Component
	activator
	handleEventer
}

type updater interface {
	Update() error
}

type activator interface {
	SetActivate(bool)
	IsActivated() bool
}

type handleEventer interface {
	HandleEvent(ui.Event) error
}
