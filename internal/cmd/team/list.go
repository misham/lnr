package team

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/misham/linear-cli/internal/cmdutil"
	"github.com/misham/linear-cli/internal/ui"
)

func newListCmd(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List available teams",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}
			teams, err := client.ListTeams(cmd.Context())
			if err != nil {
				return fmt.Errorf("list teams: %w", err)
			}
			return ui.PrintTeams(f.IO, teams)
		},
	}
}
