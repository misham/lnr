package issue

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/misham/linear-cli/internal/api"
	"github.com/misham/linear-cli/internal/cmdutil"
	"github.com/misham/linear-cli/internal/ui"
)

func newUpdateCmd(f *cmdutil.Factory) *cobra.Command {
	var (
		title       string
		description string
		priority    int
	)

	cmd := &cobra.Command{
		Use:   "update <identifier>",
		Short: "Update an issue",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			input := api.IssueUpdateInput{}

			if cmd.Flags().Changed("title") {
				input.Title = &title
			}
			if cmd.Flags().Changed("description") {
				input.Description = &description
			}
			if cmd.Flags().Changed("priority") {
				input.Priority = &priority
			}

			client, err := f.APIClient()
			if err != nil {
				return err
			}

			issue, err := client.UpdateIssue(cmd.Context(), args[0], input)
			if err != nil {
				return fmt.Errorf("update issue: %w", err)
			}

			return ui.PrintIssueUpdated(f.IO, issue)
		},
	}

	cmd.Flags().StringVar(&title, "title", "", "Issue title")
	cmd.Flags().StringVar(&description, "description", "", "Issue description")
	cmd.Flags().IntVar(&priority, "priority", 0, "Priority (0=none, 1=urgent, 2=high, 3=medium, 4=low)")

	return cmd
}
