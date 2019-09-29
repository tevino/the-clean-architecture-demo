package component

import ui "github.com/gizak/termui/v3"

// DummyComponent is a Component that does nothing.
type DummyComponent struct{}

func (DummyComponent) Update() error              { return nil }
func (DummyComponent) IsActivated() bool          { return false }
func (DummyComponent) SetActivate(bool)           {}
func (DummyComponent) HandleEvent(ui.Event) error { return nil }
