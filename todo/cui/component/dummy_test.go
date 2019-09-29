package component

import (
	"testing"

	ui "github.com/gizak/termui/v3"
	"github.com/stretchr/testify/assert"
)

func TestDummyComponent(t *testing.T) {
	t.Parallel()
	var c InteractiveComponent = &DummyComponent{}
	assert.NotPanics(t, func() {
		assert.NoError(t, c.Update())
		assert.False(t, c.IsActivated())
		c.SetActivate(true)
		assert.False(t, c.IsActivated())
		c.SetActivate(false)
		assert.False(t, c.IsActivated())
		assert.NoError(t, c.HandleEvent(ui.Event{}))
	})
}
