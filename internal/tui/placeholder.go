package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type placeholderTab struct {
	name string
	size tea.WindowSizeMsg
}

func newPlaceholderTab(name string) *placeholderTab {
	return &placeholderTab{name: name}
}

func (p *placeholderTab) Init() tea.Cmd {
	return nil
}

func (p *placeholderTab) Update(msg tea.Msg) (TabModel, tea.Cmd) {
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		p.size = msg
	}
	return p, nil
}

func (p *placeholderTab) View() string {
	style := lipgloss.NewStyle().
		Width(p.size.Width).
		Height(p.size.Height).
		Align(lipgloss.Center, lipgloss.Center)
	return style.Render(p.name)
}

func (p *placeholderTab) ShortHelp() []key.Binding {
	return nil
}
