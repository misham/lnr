package ui

import (
	"fmt"
	"text/tabwriter"

	"github.com/charmbracelet/lipgloss"

	"github.com/misham/linear-cli/internal/api"
)

var (
	headerStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("243"))
	keyStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
)

// PrintTeams formats and prints teams to the IOStreams output.
func PrintTeams(ios *IOStreams, teams []api.Team) error {
	if len(teams) == 0 {
		_, err := fmt.Fprintln(ios.Out, "No teams found")
		return err
	}

	if ios.IsPlain() {
		return printTeamsPlain(ios, teams)
	}
	return printTeamsStyled(ios, teams)
}

func printTeamsPlain(ios *IOStreams, teams []api.Team) error {
	for _, team := range teams {
		if _, err := fmt.Fprintf(ios.Out, "%s\t%s\t%s\n", team.Key, team.Name, team.Description); err != nil {
			return err
		}
	}
	return nil
}

func printTeamsStyled(ios *IOStreams, teams []api.Team) error {
	w := tabwriter.NewWriter(ios.Out, 0, 0, 2, ' ', 0)

	header := fmt.Sprintf("%s\t%s\t%s",
		headerStyle.Render("KEY"),
		headerStyle.Render("NAME"),
		headerStyle.Render("DESCRIPTION"),
	)
	if _, err := fmt.Fprintln(w, header); err != nil {
		return err
	}

	for _, team := range teams {
		row := fmt.Sprintf("%s\t%s\t%s",
			keyStyle.Render(team.Key),
			team.Name,
			team.Description,
		)
		if _, err := fmt.Fprintln(w, row); err != nil {
			return err
		}
	}

	return w.Flush()
}
