package component

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

//go:generate mockgen -destination mock_component/text_mock.go github.com/tevino/the-clean-architecture-demo/todo/cui/component Text

// Text represents a readonly text element.
type Text interface {
	Plain(string)
	Info(string)
	Warn(error)
}

type TextComponent struct {
	*widgets.Paragraph
}

func NewTextComponent(title string) *TextComponent {
	p := widgets.NewParagraph()
	p.Title = title
	p.TitleStyle.Modifier = ui.ModifierBold
	return &TextComponent{
		Paragraph: p,
	}
}

func (t *TextComponent) Plain(s string) {
	t.TextStyle = ui.StyleClear
	t.Text = s
}

func (t *TextComponent) Info(s string) {
	t.TextStyle.Fg = ui.ColorWhite
	t.BorderStyle.Fg = ui.ColorCyan
	t.Text = s
}

func (t *TextComponent) Warn(e error) {
	t.TextStyle.Fg = ui.ColorWhite
	t.BorderStyle.Fg = ui.ColorCyan
	t.Text = e.Error()
}
