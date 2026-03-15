package project

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/misham/linear-cli/internal/api"
	"github.com/misham/linear-cli/internal/cmdutil"
	"github.com/misham/linear-cli/internal/ui"
)

var validProjectStatuses = []string{"backlog", "planned", "started", "paused", "completed", "canceled"}

func newListCmd(f *cmdutil.Factory) *cobra.Command {
	var statusFilter string

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List projects",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if statusFilter != "" {
				valid := false
				for _, s := range validProjectStatuses {
					if s == statusFilter {
						valid = true
						break
					}
				}
				if !valid {
					return fmt.Errorf("invalid status %q: must be one of %v", statusFilter, validProjectStatuses)
				}
			}

			teamID, err := cmdutil.ResolveTeamID(cmd.Context(), f)
			if err != nil {
				return err
			}

			client, err := f.APIClient()
			if err != nil {
				return err
			}

			var all []api.Project
			cursor := ""
			for {
				result, err := client.ListProjects(cmd.Context(), teamID, statusFilter, 50, cursor)
				if err != nil {
					return fmt.Errorf("list projects: %w", err)
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

	cmd.Flags().StringVar(&statusFilter, "status", "", "Filter by status type (backlog, planned, started, paused, completed, canceled)")

	return cmd
}
