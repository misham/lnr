package cycle

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/misham/linear-cli/internal/cmdutil"
)

func newRemoveIssueCmd(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "remove-issue <issue-id>",
		Short: "Remove an issue from its cycle",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			issueID := args[0]
			err = client.RemoveIssueCycle(cmd.Context(), issueID)
			if err != nil {
				return fmt.Errorf("remove issue from cycle: %w", err)
			}

			_, err = fmt.Fprintf(f.IO.Out, "Removed issue %s from its cycle\n", issueID)
			return err
		},
	}
}
