package tui

import "charm.land/bubbles/v2/key"

// ListKeyMap defines the key bindings for the list view.
type ListKeyMap struct {
	Up      key.Binding
	Down    key.Binding
	SortName key.Binding
	SortKind key.Binding
	Select  key.Binding
	Refresh key.Binding
	Quit    key.Binding
}

// OverlayKeyMap defines the key bindings for overlays.
type OverlayKeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Left   key.Binding
	Right  key.Binding
	Select key.Binding
	Close  key.Binding
}

var listKeys = ListKeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	SortName: key.NewBinding(
		key.WithKeys("left"),
		key.WithHelp("←", "sort by name"),
	),
	SortKind: key.NewBinding(
		key.WithKeys("right"),
		key.WithHelp("→", "sort by kind"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Refresh: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "refresh"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

var overlayKeys = OverlayKeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left"),
		key.WithHelp("←", "decrease"),
	),
	Right: key.NewBinding(
		key.WithKeys("right"),
		key.WithHelp("→", "increase"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "execute"),
	),
	Close: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "close"),
	),
}
