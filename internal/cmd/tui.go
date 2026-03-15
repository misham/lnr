package cmd

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/misham/linear-cli/internal/cmdutil"
	"github.com/misham/linear-cli/internal/tui"
)

func newTuiCmd(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "tui",
		Short: "Launch interactive dashboard",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if f.IO.IsPlain() {
				return fmt.Errorf("lnr tui requires an interactive terminal")
			}

			client, err := f.APIClient()
			if err != nil {
				return err
			}

			teamID, err := cmdutil.ResolveTeamID(cmd.Context(), f)
			if err != nil {
				return err
			}

			teams, err := client.ListTeams(cmd.Context())
			if err != nil {
				return fmt.Errorf("listing teams: %w", err)
			}
			teamName := teamID
			for _, t := range teams {
				if t.ID == teamID {
					teamName = t.Name
					break
				}
			}

			ctx, cancel := context.WithCancel(cmd.Context())
			defer cancel()

			tabNames := []string{"Issues", "Cycles", "Projects", "Initiatives"}
			tabs := []tui.TabModel{
				tui.NewIssueTab(ctx, client, teamID),
				tui.NewCycleTab(ctx, client, teamID),
				tui.NewProjectTab(ctx, client, teamID),
				tui.NewInitiativeTab(ctx, client),
			}

			app := tui.NewApp(ctx, cancel, tabs, tabNames, teamName)
			p := tea.NewProgram(app, tea.WithAltScreen())
			_, err = p.Run()
			return err
		},
	}
}
