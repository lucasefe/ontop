package tui

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines the keybindings for the TUI
type KeyMap struct {
	Up            key.Binding
	Down          key.Binding
	Left          key.Binding
	Right         key.Binding
	Select        key.Binding
	Back          key.Binding
	Quit          key.Binding
	Help          key.Binding
	Move          key.Binding
	Archive       key.Binding
	Delete        key.Binding
	Refresh       key.Binding
	New           key.Binding
	Edit          key.Binding
	Tab           key.Binding
	Sort          key.Binding
	ToggleArchive key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Select, k.Back, k.New, k.Edit},
		{k.Move, k.Archive, k.Delete, k.Refresh},
		{k.Sort, k.ToggleArchive, k.Help, k.Quit},
	}
}

// DefaultKeyMap returns the default keybindings
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up: key.NewBinding(
			key.WithKeys("k", "up"),
			key.WithHelp("k/↑", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("j", "down"),
			key.WithHelp("j/↓", "down"),
		),
		Left: key.NewBinding(
			key.WithKeys("h", "left"),
			key.WithHelp("h/←", "left"),
		),
		Right: key.NewBinding(
			key.WithKeys("l", "right"),
			key.WithHelp("l/→", "right"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select/details"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back/cancel"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		Move: key.NewBinding(
			key.WithKeys("m"),
			key.WithHelp("m", "move task"),
		),
		Archive: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "archive/unarchive"),
		),
		Delete: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete"),
		),
		Refresh: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "refresh"),
		),
		New: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "new task"),
		),
		Edit: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit"),
		),
		Tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next field"),
		),
		Sort: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "cycle sort"),
		),
		ToggleArchive: key.NewBinding(
			key.WithKeys("z"),
			key.WithHelp("z", "toggle archive view"),
		),
	}
}
