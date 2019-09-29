package io

import ui "github.com/gizak/termui/v3"

//go:generate mockgen -destination mock_cui/io_mock.go github.com/tevino/the-clean-architecture-demo/todo/cui/io CUILib

// CUILib represents the CUI library.
type CUILib interface {
	Init() error
	Close()
	Clear()
	PollEvents() <-chan ui.Event
	Render(a ...ui.Drawable)
	TerminalDimensions() (int, int)
}

type TermUI struct{}

func (TermUI) Init() error {
	return ui.Init()
}

func (TermUI) Close() {
	ui.Close()
}

func (TermUI) Clear() {
	ui.Clear()
}

func (TermUI) PollEvents() <-chan ui.Event {
	return ui.PollEvents()
}

func (TermUI) Render(a ...ui.Drawable) {
	ui.Render(a...)
}

func (TermUI) TerminalDimensions() (int, int) {
	return ui.TerminalDimensions()
}
