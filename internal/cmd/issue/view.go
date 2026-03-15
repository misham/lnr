package issue

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/misham/linear-cli/internal/cmdutil"
	"github.com/misham/linear-cli/internal/ui"
)

func newViewCmd(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "view <identifier>",
		Short: "View issue details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}
			issue, err := client.GetIssue(cmd.Context(), args[0])
			if err != nil {
				return fmt.Errorf("get issue: %w", err)
			}
			return ui.PrintIssueDetail(f.IO, issue)
		},
	}
}
