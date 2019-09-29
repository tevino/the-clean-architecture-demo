package component

import (
	"fmt"

	ui "github.com/gizak/termui/v3"
)

type GridComponent struct {
	*ui.Grid
	componentKeyMap    map[string]InteractiveComponent
	selectingComponent bool
	activatedComponent InteractiveComponent
}

// DefaultActivated is a special key used in the keymap while creating GridComponent to make a Component activated by default.
const DefaultActivated = "DEFAULT_ACTIVATED"

// NewGridComponent creates a GridComponent.
func NewGridComponent(keyMap map[string]InteractiveComponent, items ...ui.GridItem) *GridComponent {
	itemInterfaces := make([]interface{}, len(items))
	for i, it := range items {
		itemInterfaces[i] = it
	}
	grid := ui.NewGrid()
	grid.Set(itemInterfaces...)

	var activatedComponent InteractiveComponent = &DummyComponent{}
	if comp, ok := keyMap[DefaultActivated]; ok {
		comp.SetActivate(true)
		activatedComponent = comp
	}
	return &GridComponent{
		Grid:               grid,
		activatedComponent: activatedComponent,
		componentKeyMap:    keyMap,
	}
}

// HandleEvent handles Component activation.
func (g *GridComponent) HandleEvent(e ui.Event) error {
	switch e.ID {
	case "<C-w>":
		g.selectingComponent = true
	default:
		if g.selectingComponent {
			g.selectingComponent = false
			comp, keyMatch := g.componentKeyMap[e.ID]
			if keyMatch {
				g.activatedComponent.SetActivate(false)
				g.activatedComponent = comp
				g.activatedComponent.SetActivate(true)
			}
		} else {
			err := g.activatedComponent.HandleEvent(e)
			if err != nil {
				return fmt.Errorf("%s handling event %v: %w", g.activatedComponent, e, err)
			}
		}
	}
	return nil
}
