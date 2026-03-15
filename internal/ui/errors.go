package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var errorPrefixStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("196"))

// PrintError formats an error for the terminal.
func (s *IOStreams) PrintError(err error) {
	if s.IsPlain() {
		_, _ = fmt.Fprintf(s.ErrOut, "error: %v\n", err)
		return
	}
	_, _ = fmt.Fprintf(s.ErrOut, "%s %v\n", errorPrefixStyle.Render("Error:"), err)
}
