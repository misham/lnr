package project

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/misham/linear-cli/internal/cmdutil"
	"github.com/misham/linear-cli/internal/ui"
)

func newViewCmd(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "view <name-id-or-slug>",
		Short: "View project details",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				_ = cmd.Usage()
				return fmt.Errorf("requires exactly 1 arg: project name, id, or slug")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			teamID, err := cmdutil.ResolveTeamID(cmd.Context(), f)
			if err != nil {
				return err
			}

			client, err := f.APIClient()
			if err != nil {
				return err
			}

			project, err := resolveProject(cmd.Context(), client, teamID, args[0])
			if err != nil {
				return err
			}
			return ui.PrintProjectDetail(f.IO, project)
		},
	}
}
