package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"

	"github.com/misham/linear-cli/internal/api"
)

// PrintStates formats and prints a list of workflow states.
func PrintStates(ios *IOStreams, states []api.WorkflowState) error {
	if len(states) == 0 {
		_, err := fmt.Fprintln(ios.Out, "No workflow states found")
		return err
	}
	if ios.IsPlain() {
		return printStatesPlain(ios, states)
	}
	return printStatesStyled(ios, states)
}

func printStatesPlain(ios *IOStreams, states []api.WorkflowState) error {
	for _, s := range states {
		if _, err := fmt.Fprintf(ios.Out, "%s\t%s\t%s\t%.0f\n",
			s.Name,
			s.Type,
			s.Color,
			s.Position,
		); err != nil {
			return err
		}
	}
	return nil
}

func printStatesStyled(ios *IOStreams, states []api.WorkflowState) error {
	nameW, catW, colorW := len("NAME"), len("CATEGORY"), len("COLOR")
	for _, s := range states {
		if l := len(s.Name); l > nameW {
			nameW = l
		}
		if l := len(s.Type); l > catW {
			catW = l
		}
		if l := len(s.Color); l > colorW {
			colorW = l
		}
	}

	const gap = "  "
	ew := &errWriter{w: ios.Out}
	ew.printf("%s%s%s%s%s%s%s\n",
		padRight(headerStyle.Render("NAME"), len("NAME"), nameW), gap,
		padRight(headerStyle.Render("CATEGORY"), len("CATEGORY"), catW), gap,
		padRight(headerStyle.Render("COLOR"), len("COLOR"), colorW), gap,
		headerStyle.Render("POSITION"),
	)
	for _, s := range states {
		posStr := fmt.Sprintf("%.0f", s.Position)
		colorStyled := lipgloss.NewStyle().Foreground(lipgloss.Color(s.Color)).Render(s.Color)
		ew.printf("%s%s%s%s%s%s%s\n",
			padRight(s.Name, len(s.Name), nameW), gap,
			padRight(s.Type, len(s.Type), catW), gap,
			padRight(colorStyled, len(s.Color), colorW), gap,
			posStr,
		)
	}
	return ew.err
}
