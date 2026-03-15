package initiative

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/misham/linear-cli/internal/api"
	"github.com/misham/linear-cli/internal/cmdutil"
	"github.com/misham/linear-cli/internal/ui"
)

func newProjectsCmd(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "projects <name-or-id>",
		Short: "List initiative projects",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				_ = cmd.Usage()
				return fmt.Errorf("requires exactly 1 arg: initiative name or id")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			initiative, err := resolveInitiative(cmd.Context(), client, args[0])
			if err != nil {
				return err
			}

			var all []api.Project
			cursor := ""
			for {
				result, err := client.ListInitiativeProjects(cmd.Context(), initiative.ID, 50, cursor)
				if err != nil {
					return fmt.Errorf("list initiative projects: %w", err)
				}
				all = append(all, result.Projects...)
				if !result.PageInfo.HasNextPage {
					break
				}
				cursor = result.PageInfo.EndCursor
			}

			return ui.PrintProjects(f.IO, all)
		},
	}
}
