package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Color palette.
	purpleColor = lipgloss.Color("#7D56F4")
	pinkColor   = lipgloss.Color("#F25D94")
	greenColor  = lipgloss.Color("#73F59F")
	cyanColor   = lipgloss.Color("#56B6C2")
	yellowColor = lipgloss.Color("#E5C07B")
	mutedColor  = lipgloss.Color("#555")
	dimColor    = lipgloss.Color("#333")
	textColor   = lipgloss.Color("#c9d1d9")
	brightColor = lipgloss.Color("#fff")

	// Tab bar styles — minimal underline style.
	activeTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(purpleColor).
			Padding(0, 2)

	inactiveTabStyle = lipgloss.NewStyle().
				Foreground(mutedColor).
				Padding(0, 2)

	// Status bar styles.
	statusBarStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Background(lipgloss.Color("#161b22"))
)
