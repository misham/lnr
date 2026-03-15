package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

func renderStatusBar(width int, helpKeys []key.Binding, info string) string {
	var parts []string
	for _, k := range helpKeys {
		if k.Enabled() {
			parts = append(parts, k.Help().Key+" "+k.Help().Desc)
		}
	}
	left := strings.Join(parts, "  ")
	right := info

	gap := max(width-lipgloss.Width(left)-lipgloss.Width(right), 1)

	line := left + strings.Repeat(" ", gap) + right
	return statusBarStyle.Width(width).Render(line)
}
