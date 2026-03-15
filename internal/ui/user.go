package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"

	"github.com/misham/linear-cli/internal/api"
)

var (
	labelStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("243"))
	valueStyle = lipgloss.NewStyle()
)

// PrintUser formats and prints user details to the IOStreams output.
func PrintUser(ios *IOStreams, user *api.User) error {
	if ios.IsPlain() {
		return printUserPlain(ios, user)
	}
	return printUserStyled(ios, user)
}

func printUserPlain(ios *IOStreams, user *api.User) error {
	_, err := fmt.Fprintf(ios.Out, "%s\t%s\t%s\t%s\t%t\n",
		user.ID, user.Name, user.DisplayName, user.Email, user.Active)
	return err
}

func printUserStyled(ios *IOStreams, user *api.User) error {
	fields := []struct {
		label string
		value string
	}{
		{"Name", user.Name},
		{"Display Name", user.DisplayName},
		{"Email", user.Email},
		{"ID", user.ID},
		{"Active", fmt.Sprintf("%t", user.Active)},
	}

	for _, f := range fields {
		if _, err := fmt.Fprintf(ios.Out, "%s %s\n",
			labelStyle.Render(f.label+":"),
			valueStyle.Render(f.value),
		); err != nil {
			return err
		}
	}
	return nil
}
