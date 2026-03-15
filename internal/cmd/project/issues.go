package project

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/misham/linear-cli/internal/api"
	"github.com/misham/linear-cli/internal/cmdutil"
	"github.com/misham/linear-cli/internal/ui"
)

func newIssuesCmd(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "issues <name-id-or-slug>",
		Short: "List project issues",
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

			var all []api.Issue
			cursor := ""
			for {
				result, err := client.ListProjectIssues(cmd.Context(), project.ID, 50, cursor)
				if err != nil {
					return fmt.Errorf("list project issues: %w", err)
				}
				all = append(all, result.Issues...)
				if !result.PageInfo.HasNextPage {
					break
				}
				cursor = result.PageInfo.EndCursor
			}

			return ui.PrintIssues(f.IO, all)
		},
	}
}
