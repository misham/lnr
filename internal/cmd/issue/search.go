package issue

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/misham/linear-cli/internal/cmdutil"
	"github.com/misham/linear-cli/internal/ui"
)

func newSearchCmd(f *cmdutil.Factory) *cobra.Command {
	var (
		query string
		limit int
	)

	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search issues",
		RunE: func(cmd *cobra.Command, args []string) error {
			if query == "" {
				return fmt.Errorf("query is required (use --query)")
			}

			client, err := f.APIClient()
			if err != nil {
				return err
			}

			// Team ID is optional for search — boosts relevance but doesn't filter.
			var teamID string
			if f.TeamKey != "" {
				teamID, err = cmdutil.ResolveTeamID(cmd.Context(), f)
				if err != nil {
					return err
				}
			}

			result, err := client.SearchIssues(cmd.Context(), query, teamID, limit, "")
			if err != nil {
				return fmt.Errorf("search issues: %w", err)
			}

			return ui.PrintIssues(f.IO, result.Issues)
		},
	}

	cmd.Flags().StringVarP(&query, "query", "q", "", "Search text (required)")
	cmd.Flags().IntVar(&limit, "limit", 50, "Maximum number of results")

	return cmd
}
