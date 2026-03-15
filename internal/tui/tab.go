package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// TabModel defines the interface that each tab must implement.
type TabModel interface {
	Init() tea.Cmd
	Update(tea.Msg) (TabModel, tea.Cmd)
	View() string
	ShortHelp() []key.Binding
}
