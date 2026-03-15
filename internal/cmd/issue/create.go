package issue

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/misham/linear-cli/internal/api"
	"github.com/misham/linear-cli/internal/cmdutil"
	"github.com/misham/linear-cli/internal/ui"
)

func newCreateCmd(f *cmdutil.Factory) *cobra.Command {
	var (
		title       string
		description string
		priority    int
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new issue",
		RunE: func(cmd *cobra.Command, args []string) error {
			if title == "" {
				return fmt.Errorf("title is required (use --title)")
			}

			teamID, err := cmdutil.ResolveTeamID(cmd.Context(), f)
			if err != nil {
				return err
			}

			input := api.IssueCreateInput{
				Title:       title,
				TeamID:      teamID,
				Description: description,
				Priority:    priority,
			}

			client, err := f.APIClient()
			if err != nil {
				return err
			}

			issue, err := client.CreateIssue(cmd.Context(), input)
			if err != nil {
				return fmt.Errorf("create issue: %w", err)
			}

			return ui.PrintIssueCreated(f.IO, issue)
		},
	}

	cmd.Flags().StringVar(&title, "title", "", "Issue title (required)")
	cmd.Flags().StringVar(&description, "description", "", "Issue description")
	cmd.Flags().IntVar(&priority, "priority", 0, "Priority (0=none, 1=urgent, 2=high, 3=medium, 4=low)")

	return cmd
}
